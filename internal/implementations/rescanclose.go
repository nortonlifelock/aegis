package implementations

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/nortonlifelock/aegis/internal/integrations"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/jira"
	"github.com/nortonlifelock/log"
)

// ScanCloseJob are created by rescan jobs, and do not have to be made by the user
type ScanCloseJob struct {
	Payload *ScanClosePayload

	id          string
	payloadJSON string
	ctx         context.Context
	db          domain.DatabaseConnection
	lstream     log.Logger
	appconfig   domain.Config
	config      domain.JobConfig
	insource    domain.SourceConfig
	outsource   domain.SourceConfig

	decommedDevices sync.Map
}

// ScanClosePayload is used to parse information from the job history Payload, which is generated automatically
type ScanClosePayload struct {
	RescanPayload
	Scan    interface{} `json:"scan"`
	Devices []string    `json:"devices"`
	Group   string      `json:"group"`
	ScanID  string      `json:"scan_id"`
}

// buildPayload parses the information from the Payload of the job history entry
func (job *ScanCloseJob) buildPayload(pjson string) (err error) {
	if len(pjson) > 0 {
		job.Payload = &ScanClosePayload{}
		err = json.Unmarshal([]byte(pjson), job.Payload)
	} else {
		err = fmt.Errorf("empty json string passed to ScanCloseJob")
	}
	return err
}

// Process loads and processes the results from the scanner. This includes updating the status of the associated JIRA ticket as well as creating exceptions in
// the ignore table if the asset is discovered to be decommissioned
func (job *ScanCloseJob) Process(ctx context.Context, id string, appconfig domain.Config, db domain.DatabaseConnection, lstream log.Logger, payload string, jobConfig domain.JobConfig, inSource []domain.SourceConfig, outSource []domain.SourceConfig) (err error) {

	var ok bool
	if job.ctx, job.id, job.appconfig, job.db, job.lstream, job.payloadJSON, job.config, job.insource, job.outsource, ok = validInputs(ctx, id, appconfig, db, lstream, payload, jobConfig, inSource, outSource); ok {

		if err = job.buildPayload(job.payloadJSON); err == nil {

			var vscanner integrations.Vscanner
			if vscanner, err = integrations.NewVulnScanner(job.ctx, job.insource.Source(), job.db, job.lstream, job.appconfig, job.insource); err == nil {

				job.lstream.Send(log.Debug("Scanning Engine Connection Initialized"))
				var engine integrations.TicketingEngine
				if engine, err = integrations.GetEngine(job.ctx, job.outsource.Source(), job.db, job.lstream, job.appconfig, job.outsource); err == nil {
					job.lstream.Send(log.Debug("Ticketing Engine Connection Initialized"))

					var tickets []domain.Ticket
					if tickets, err = loadTickets(job.lstream, engine, job.Payload.Tickets); err == nil {

						var additionalTickets []domain.Ticket
						if additionalTickets, err = job.loadAdditionalTickets(engine, tickets); err == nil {
							tickets = append(tickets, additionalTickets...)

							job.lstream.Send(log.Debugf("Loading scan [Id: %v] details.", job.Payload.ScanID))
							var scan domain.ScanSummary
							if scan, err = job.getScanFromPayload(); err == nil {
								job.processScanDetections(engine, vscanner, tickets, scan)
							} else {
								job.lstream.Send(log.Error("error while grabbing scan information from the database", err))
							}
						} else {
							job.lstream.Send(log.Errorf(err, "error while loading additional tickets for the scan"))
						}
					} else {
						job.lstream.Send(log.Errorf(err, "error while loading tickets for the scan"))
					}
				} else {
					job.lstream.Send(log.Errorf(err, "error while creating the ticketing connection"))
				}
			} else {
				job.lstream.Send(log.Errorf(err, "error while creating the vulnerability scanner"))
			}
		} else {
			err = fmt.Errorf("error while building payload - %s", err.Error())
		}
	} else {
		err = fmt.Errorf("input validation failed")
	}

	return err
}

func (job *ScanCloseJob) processScanDetections(engine integrations.TicketingEngine, vscanner integrations.Vscanner, tickets []domain.Ticket, scan domain.ScanSummary) {
	var err error
	var payload []byte
	if payload, err = scanClosePayloadToScanPayload(scan.ScanClosePayload()); err == nil {
		var detections <-chan domain.Detection
		var deadHostIPToProof <-chan domain.KeyValue
		if detections, deadHostIPToProof, err = vscanner.ScanResults(job.ctx, payload); err == nil {

			var wg sync.WaitGroup
			deviceIDToVulnIDToDetection, deadHostIPToProofMap := job.mapDetectionsAndDeadHosts(detections, deadHostIPToProof)
			for _, ticket := range tickets {
				wg.Add(1)
				go func(ticket domain.Ticket) {
					defer handleRoutinePanic(job.lstream)
					defer wg.Done()

					if sord(ticket.Status()) != engine.GetStatusMap(jira.StatusClosedRemediated) { // TODO: Ensure this is correct
						job.modifyJiraTicketAccordingToVulnerabilityStatus(
							engine,
							&lastCheckedTicket{ticket},
							scan,
							deadHostIPToProofMap,
							deviceIDToVulnIDToDetection,
						)
					}
				}(ticket)
			}
			wg.Wait()

		} else {
			job.lstream.Send(log.Errorf(err, "error while gathering scan information"))
		}
	} else {
		job.lstream.Send(log.Errorf(err, "error while pulling scan information from scan close payload"))
	}
}

func (job *ScanCloseJob) mapDetectionsAndDeadHosts(detections <-chan domain.Detection, deadHostIPToProof <-chan domain.KeyValue) (deviceIDToVulnIDToDetection map[string]map[string]domain.Detection, deadHostIPToProofMap map[string]string) {
	// first gather all detections from the scan results
	wg := &sync.WaitGroup{}
	deviceIDToVulnIDToDetection = make(map[string]map[string]domain.Detection)
	func() {
		lock := &sync.Mutex{}

		for {
			select {
			case <-job.ctx.Done():
				return
			case detection, ok := <-detections:
				if ok {
					wg.Add(1)
					go func(detection domain.Detection) {
						defer handleRoutinePanic(job.lstream)
						defer wg.Done()

						if dev, err := detection.Device(); err == nil {

							if len(sord(dev.SourceID())) > 0 {
								lock.Lock()
								if deviceIDToVulnIDToDetection[sord(dev.SourceID())] == nil {
									deviceIDToVulnIDToDetection[sord(dev.SourceID())] = make(map[string]domain.Detection)
								}

								deviceIDToVulnIDToDetection[sord(dev.SourceID())][detection.VulnerabilityID()] = detection
								lock.Unlock()
							} else {
								job.lstream.Send(log.Errorf(err, "empty device ID found for detection"))
							}
						} else {
							job.lstream.Send(log.Errorf(err, "error while loading detection %v", detection.VulnerabilityID()))
						}
					}(detection)
				} else {
					return
				}
			}
		}
	}()

	wg.Wait()

	deadHostIPToProofMap = make(map[string]string)
	func() {
		for {
			select {
			case <-job.ctx.Done():
				return
			case deadHost, ok := <-deadHostIPToProof:
				if ok {
					deadHostIPToProofMap[deadHost.Key()] = deadHost.Value()
				} else {
					return
				}
			}
		}
	}()
	return deviceIDToVulnIDToDetection, deadHostIPToProofMap
}

// transitions the status of the JIRA ticket if necessary
func (job *ScanCloseJob) modifyJiraTicketAccordingToVulnerabilityStatus(engine integrations.TicketingEngine, ticket domain.Ticket, scan domain.ScanSummary, deadHostIPToProofMap map[string]string, deviceIDToVulnIDToDetection map[string]map[string]domain.Detection) {
	var detectionsFoundForDevice = deviceIDToVulnIDToDetection[ticket.DeviceID()] != nil
	var deviceReportedAsDead = len(deadHostIPToProofMap[sord(ticket.IPAddress())]) > 0 && len(sord(ticket.IPAddress())) > 0

	var detection domain.Detection
	if detectionsFoundForDevice {
		detection = deviceIDToVulnIDToDetection[ticket.DeviceID()][ticket.VulnerabilityID()]
	}

	var err error
	var status string
	var inactiveKernel bool
	if detection != nil {
		status = detection.Status()
		inactiveKernel = iord(detection.ActiveKernel()) > 0
	}

	if detectionsFoundForDevice || deviceReportedAsDead {
		switch job.Payload.Type {
		case domain.RescanNormal:
			job.processTicketForNormalRescan(deadHostIPToProofMap, ticket, detection, err, engine, status, inactiveKernel, scan)
		case domain.RescanDecommission:
			job.processTicketForDecommRescan(deadHostIPToProofMap, ticket, detection, err, engine, scan, status)
		case domain.RescanExceptions:
			job.processTicketForPassiveOrExceptionRescan(deadHostIPToProofMap, ticket, detection, err, engine, status, inactiveKernel, scan)
		case domain.RescanPassive:
			job.processTicketForPassiveOrExceptionRescan(deadHostIPToProofMap, ticket, detection, err, engine, status, inactiveKernel, scan)
		default:
			job.lstream.Send(log.Critical(fmt.Sprintf("Unrecognized scan type [%s]", job.Payload.Type), nil))
		}
	} else {
		job.lstream.Send(log.Errorf(err, "scan [%s] did not seem to cover the device %v - scanner did not report any data for device", job.Payload.ScanID, ticket.DeviceID()))
		err = engine.Transition(ticket, engine.GetStatusMap(jira.StatusScanError), fmt.Sprintf("scan [%s] did not cover the device %v. Please make sure this asset is still in-scope and associated with an asset group. If this asset is out of scope, please move this ticket to NOTAVRR status or alert the vulnerability management team.", job.Payload.ScanID, ticket.DeviceID()), sord(ticket.AssignedTo()))
		if err != nil {
			job.lstream.Send(log.Errorf(err, "error while adding comment to ticket [%s]", ticket.Title()))
		}
	}

}

func (job *ScanCloseJob) processTicketForPassiveOrExceptionRescan(deadHostIPToProofMap map[string]string, ticket domain.Ticket, detection domain.Detection, err error, engine integrations.TicketingEngine, status string, inactiveKernel bool, scan domain.ScanSummary) {
	if (len(deadHostIPToProofMap[*ticket.IPAddress()]) > 0 && len(*ticket.IPAddress()) > 0) || (detection != nil && detection.Status() == domain.DeadHost) {
		job.lstream.Send(log.Infof("the device for %s seems to be dead, but this is not a decommission scan", ticket.Title()))
		err = engine.Transition(ticket, engine.GetStatusMap(jira.StatusResolvedDecom), fmt.Sprintf("The device could not be detected though a vulnerability rescan. It has been moved to a resolved decommission status and will be rescanned with another option profile to confirm\nPROOF:\n%s", deadHostIPToProofMap[sord(ticket.IPAddress())]), sord(ticket.AssignedTo()))
		if err != nil {
			job.lstream.Send(log.Errorf(err, "error while adding comment to ticket %s", ticket.Title()))
		}
	} else if detection == nil || status == domain.Fixed {
		// Non-decommission scan, the detection appears to be fixed, so close the ticket
		var closeReason = closeComment
		if inactiveKernel {
			closeReason = inactiveKernelComment
		}

		job.lstream.Send(log.Infof("Vulnerability NO LONGER EXISTS, closing Ticket [%s]", ticket.Title()))
		if err = job.closeTicket(engine, ticket, scan, engine.GetStatusMap(jira.StatusClosedRemediated), closeReason); err != nil {
			job.lstream.Send(log.Errorf(err, "Error while closing Ticket [%s]", ticket.Title()))
		}
	} else if status == domain.Vulnerable {
		// Passive/exception rescans don't reopen tickets - leave alone
	} else {
		job.lstream.Send(log.Errorf(fmt.Errorf("unrecognized status [%v]", status), "scan close job did not recognize vulnerability status"))
	}
}

func (job *ScanCloseJob) processTicketForDecommRescan(deadHostIPToProofMap map[string]string, ticket domain.Ticket, detection domain.Detection, err error, engine integrations.TicketingEngine, scan domain.ScanSummary, status string) {
	if (len(deadHostIPToProofMap[*ticket.IPAddress()]) > 0 && len(*ticket.IPAddress()) > 0) || (detection != nil && detection.Status() == domain.DeadHost) {
		err = job.markDeviceAsDecommissionedInDatabase(ticket)
		if err != nil {
			job.lstream.Send(log.Errorf(err, "error while marking ticket [%s] as decommissioned in the database", ticket.Title()))
		}

		job.lstream.Send(log.Infof("Device %v confirmed offline, closing Ticket [%s]", ticket.DeviceID(), ticket.Title()))

		closeReason := fmt.Sprintf("Device found to be dead by Scanner\n%v", detection.Proof())
		if err = job.closeTicket(engine, ticket, scan, engine.GetStatusMap(jira.StatusClosedDecommissioned), closeReason); err != nil {
			job.lstream.Send(log.Errorf(err, "Error while closing Ticket [%s]", ticket.Title()))
		}
	} else if detection == nil || status == domain.Fixed || status == domain.Vulnerable {
		// Decommission scan - this block hitting means that the host wasn't marked as dead, so we should reopen the associated tickets
		if job.shouldOpenTicket(engine, sord(ticket.Status()), job.Payload.Type) { // don't need to reopen tickets that are already opened (these tickets come from the additional ticket loading)
			job.lstream.Send(log.Info(fmt.Sprintf("Ticket [%s] reopened as host is alive", ticket.Title())))
			if err = job.openTicket(ticket, detection, scan, engine); err != nil {
				job.lstream.Send(log.Errorf(err, "Error while opening Ticket [%s]", ticket.Title()))
			}
		}
	} else {
		job.lstream.Send(log.Errorf(fmt.Errorf("unrecognized status [%v]", status), "scan close job did not recognize vulnerability status"))
	}
}

func (job *ScanCloseJob) processTicketForNormalRescan(deadHostIPToProofMap map[string]string, ticket domain.Ticket, detection domain.Detection, err error, engine integrations.TicketingEngine, status string, inactiveKernel bool, scan domain.ScanSummary) {
	if (len(deadHostIPToProofMap[*ticket.IPAddress()]) > 0 && len(*ticket.IPAddress()) > 0) || (detection != nil && detection.Status() == domain.DeadHost) {
		job.lstream.Send(log.Infof("the device for %s seems to be dead, but this is not a decommission scan", ticket.Title()))
		err = engine.Transition(ticket, engine.GetStatusMap(jira.StatusResolvedDecom), fmt.Sprintf("The device could not be detected though a vulnerability rescan. It has been moved to a resolved decommission status and will be rescanned with another option profile to confirm\nPROOF:\n%s", deadHostIPToProofMap[sord(ticket.IPAddress())]), sord(ticket.AssignedTo()))
		if err != nil {
			job.lstream.Send(log.Errorf(err, "error while adding comment to ticket %s", ticket.Title()))
		}
	} else if detection == nil || status == domain.Fixed {
		// Non-decommission scan, the detection appears to be fixed, so close the ticket
		var closeReason = closeComment
		if inactiveKernel {
			closeReason = inactiveKernelComment
		}

		job.lstream.Send(log.Infof("Vulnerability NO LONGER EXISTS, closing Ticket [%s]", ticket.Title()))
		if err = job.closeTicket(engine, ticket, scan, engine.GetStatusMap(jira.StatusClosedRemediated), closeReason); err != nil {
			job.lstream.Send(log.Errorf(err, "Error while closing Ticket [%s]", ticket.Title()))
		}
	} else if status == domain.Vulnerable {
		if job.shouldOpenTicket(engine, sord(ticket.Status()), job.Payload.Type) {
			job.lstream.Send(log.Infof("Vulnerability STILL EXISTS, re-opening Ticket [%s]", ticket.Title()))
			if err = job.openTicket(ticket, detection, scan, engine); err != nil {
				job.lstream.Send(log.Errorf(err, "Error while opening Ticket [%s]", ticket.Title()))
			}
		} else {
			job.lstream.Send(log.Infof("Vulnerability still exists, but the ticket it not in the proper resolved state - leaving Ticket [%s] in its current status", ticket.Title()))
		}
	} else {
		job.lstream.Send(log.Errorf(fmt.Errorf("unrecognized status [%v]", status), "scan close job did not recognize vulnerability status"))
	}
}

func (job *ScanCloseJob) markDeviceAsDecommissionedInDatabase(ticket domain.Ticket) (err error) {
	// we use a sync map to ensure we only decommission a device a single time. if loaded has a value of false, then this is the first time the device was found
	if _, loaded := job.decommedDevices.LoadOrStore(ticket.DeviceID(), true); !loaded {
		job.lstream.Send(log.Infof("Marking device for ticket [%s] as decommissioned in the database", ticket.Title()))
		if _, _, err = job.db.DeleteIgnoreForDevice(
			job.insource.SourceID(),
			ticket.DeviceID(),
			job.config.OrganizationID()); err == nil {

			if _, _, err = job.db.SaveIgnore(
				job.insource.SourceID(),
				job.config.OrganizationID(),
				domain.DecommAsset,
				"", // empty on purpose
				ticket.DeviceID(),
				time.Now(),
				"",
				true,
				""); err != nil {

				job.lstream.Send(log.Errorf(err, "Error while updating exception with asset Id  %s: %s", ticket.DeviceID(), err.Error()))
			}
		} else {
			job.lstream.Send(log.Errorf(err, "error while deleting old entries in the ignore table for the device associated with the ticket [%s]", ticket.Title()))
		}
	}

	return err
}

type lastCheckedTicket struct {
	domain.Ticket
}

func (t *lastCheckedTicket) LastChecked() *time.Time {
	val := time.Now()
	return &val
}

func (job *ScanCloseJob) openTicket(tix domain.Ticket, vuln domain.Detection, scan domain.ScanSummary, ticketing integrations.TicketingEngine) (err error) {
	// Still Exists
	var comment string
	if job.Payload.Type == domain.RescanDecommission {
		comment = fmt.Sprintf("Ticket reopened as scan [%s] has determined that the host on this IP is responsive. Please notify the Vulnerability Management team if the running asset is different from the one identified in this ticket as IP reuse may have occurred", sord(scan.SourceKey()))
	} else {

		if len(comment) == 0 {
			comment = fmt.Sprintf("%s\n\nPROOF:\n%s\n\nScan Id [Id: %v]", reopenComment, removeHTMLTags(vuln.Proof()), sord(scan.SourceKey()))
		}
	}

	job.lstream.Send(log.Debugf("TRANSITIONING Ticket [%s]", tix.Title()))
	err = ticketing.Transition(tix, ticketing.GetStatusMap(jira.StatusReopened), comment, "Unassigned")
	if err != nil {
		job.lstream.Send(log.Errorf(err, "Unable to transition ticket %s - %s",
			tix.Title(), err.Error()))
	}
	return err
}

func (job *ScanCloseJob) shouldOpenTicket(ticketing integrations.TicketingEngine, status, jobType string) (shouldOpen bool) {
	status = strings.ToLower(status)
	jobType = strings.ToLower(jobType)

	if status == strings.ToLower(ticketing.GetStatusMap(jira.StatusResolvedRemediated)) && jobType == strings.ToLower(domain.RescanNormal) {
		shouldOpen = true
	} else if status == strings.ToLower(ticketing.GetStatusMap(jira.StatusResolvedDecom)) && jobType == strings.ToLower(domain.RescanDecommission) {
		shouldOpen = true
	}

	return shouldOpen
}

func (job *ScanCloseJob) closeTicket(ticketing integrations.TicketingEngine, tix domain.Ticket, scan domain.ScanSummary, closeStatus string, closeReason string) (err error) {
	// Remediated
	job.lstream.Send(log.Debugf("TRANSITIONING Ticket [%s]", tix.Title()))

	err = ticketing.Transition(tix, closeStatus, fmt.Sprintf("%s\n\nScan [Id: %v]", closeReason, sord(scan.SourceKey())), "Unassigned") //should this be closed?
	if err == nil {
		job.lstream.Send(log.Infof("Ticket [%s] Closed Vulnerability No Longer Exists", tix.Title()))
	} else {
		job.lstream.Send(log.Errorf(err, "Unable to transition ticket %s", tix.Title()))
	}

	return err
}

// we ask the scanning engine to look for a series of vulnerabilities across a series of devices. If we search for a vulnerability on a single
// device in a scan, we check for that vulnerability on all the devices in the scan. here we check to see if the devices have any tickets
// for vulnerabilities that they are being checked for that might not have been included in the original scan
// we only do this check in normal rescans and decommission rescans
func (job *ScanCloseJob) loadAdditionalTickets(ticketing integrations.TicketingEngine, tickets []domain.Ticket) (additionalTickets []domain.Ticket, err error) {
	additionalTickets = make([]domain.Ticket, 0)

	var additionalTicketsChan <-chan domain.Ticket
	if strings.ToLower(job.Payload.Type) == strings.ToLower(domain.RescanDecommission) {
		additionalTicketsChan, err = ticketing.GetAdditionalTicketsForDecomDevices(tickets)
	} else if strings.ToLower(job.Payload.Type) == strings.ToLower(domain.RescanNormal) {
		additionalTicketsChan, err = ticketing.GetAdditionalTicketsForVulnPerDevice(tickets)
	}

	var seen = make(map[string]bool)
	for _, tic := range tickets {
		seen[tic.Title()] = true
	}

	if err == nil && additionalTicketsChan != nil {
		for {
			if ticket, ok := <-additionalTicketsChan; ok {
				if !seen[ticket.Title()] {
					seen[ticket.Title()] = true
					additionalTickets = append(additionalTickets, ticket)
				}
			} else {
				break
			}
		}
	}

	return additionalTickets, err
}

// uses the scan id from the job history Payload to gather the scan summary from the database
func (job *ScanCloseJob) getScanFromPayload() (scan domain.ScanSummary, err error) {
	if len(job.Payload.ScanID) > 0 {
		scan, err = job.db.GetScanSummaryBySourceKey(job.Payload.ScanID)
		if err == nil {
			if scan == nil {
				err = fmt.Errorf("could not find scan in database for %v", job.Payload.ScanID)
			}
		}
	} else {
		err = errors.New("invalid ScanID given to ScanCloseJob")
	}

	return scan, err
}

func scanClosePayloadToScanPayload(scanClosePayload string) (scanPayload []byte, err error) {
	if len(scanClosePayload) > 0 {
		scp := &ScanClosePayload{}
		if err = json.Unmarshal([]byte(scanClosePayload), scp); err == nil {
			if scp.Scan != nil {
				scanPayload, err = json.Marshal(scp.Scan)
			} else {
				err = fmt.Errorf("scan information not found in scan close payload")
			}
		}
	} else {
		err = fmt.Errorf("empty scan close payload")
	}

	return scanPayload, err
}
