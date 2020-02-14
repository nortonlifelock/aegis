package implementations

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/nortonlifelock/aegis/backend/domain"
	"github.com/nortonlifelock/aegis/backend/integrations"
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
func (scanClose *ScanCloseJob) buildPayload(pjson string) (err error) {
	if len(pjson) > 0 {
		scanClose.Payload = &ScanClosePayload{}
		err = json.Unmarshal([]byte(pjson), scanClose.Payload)
	} else {
		err = fmt.Errorf("empty json string passed to ScanCloseJob")
	}
	return err
}

// Process loads and processes the results from the scanner. This includes updating the status of the associated JIRA ticket as well as creating exceptions in
// the ignore table if the asset is discovered to be decommissioned
func (scanClose *ScanCloseJob) Process(ctx context.Context, id string, appconfig domain.Config, db domain.DatabaseConnection, lstream log.Logger, payload string, jobConfig domain.JobConfig, inSource []domain.SourceConfig, outSource []domain.SourceConfig) (err error) {

	var ok bool
	if scanClose.ctx, scanClose.id, scanClose.appconfig, scanClose.db, scanClose.lstream, scanClose.payloadJSON, scanClose.config, scanClose.insource, scanClose.outsource, ok = validInputs(ctx, id, appconfig, db, lstream, payload, jobConfig, inSource, outSource); ok {

		if err = scanClose.buildPayload(scanClose.payloadJSON); err == nil {

			var vscanner integrations.Vscanner
			if vscanner, err = integrations.NewVulnScanner(scanClose.ctx, scanClose.insource.Source(), scanClose.db, scanClose.lstream, scanClose.appconfig, scanClose.insource); err == nil {

				scanClose.lstream.Send(log.Debug("Scanning Engine Connection Initialized"))
				var engine integrations.TicketingEngine
				if engine, err = integrations.GetEngine(scanClose.ctx, scanClose.outsource.Source(), scanClose.db, scanClose.lstream, scanClose.appconfig, scanClose.outsource); err == nil {
					scanClose.lstream.Send(log.Debug("Ticketing Engine Connection Initialized"))

					var tickets []domain.Ticket
					if tickets, err = loadTickets(scanClose.lstream, engine, scanClose.Payload.Tickets); err == nil {

						var additionalTickets []domain.Ticket
						if additionalTickets, err = scanClose.loadAdditionalTickets(engine, tickets); err == nil {
							tickets = append(tickets, additionalTickets...)

							scanClose.lstream.Send(log.Debugf("Loading scan [Id: %v] details.", scanClose.Payload.ScanID))
							var scan domain.ScanSummary
							if scan, err = scanClose.getScanFromPayload(); err == nil {
								scanClose.processScanDetections(engine, vscanner, tickets, scan)
							} else {
								scanClose.lstream.Send(log.Error("error while grabbing scan information from the database", err))
							}
						} else {
							scanClose.lstream.Send(log.Errorf(err, "error while loading additional tickets for the scan"))
						}
					} else {
						scanClose.lstream.Send(log.Errorf(err, "error while loading tickets for the scan"))
					}
				} else {
					scanClose.lstream.Send(log.Errorf(err, "error while creating the ticketing connection"))
				}
			} else {
				scanClose.lstream.Send(log.Errorf(err, "error while creating the vulnerability scanner"))
			}
		} else {
			err = fmt.Errorf("error while building payload - %s", err.Error())
		}
	} else {
		err = fmt.Errorf("input validation failed")
	}

	return err
}

func (scanClose *ScanCloseJob) processScanDetections(engine integrations.TicketingEngine, vscanner integrations.Vscanner, tickets []domain.Ticket, scan domain.ScanSummary) {
	var err error
	var payload []byte
	if payload, err = scanClosePayloadToScanPayload(scan.ScanClosePayload()); err == nil {
		var detections <-chan domain.Detection
		var deadHostIPToProof <-chan domain.KeyValue
		if detections, deadHostIPToProof, err = vscanner.ScanResults(scanClose.ctx, payload); err == nil {

			// first gather all detections from the scan results
			wg := &sync.WaitGroup{}
			var devDetectionMap = make(map[string]map[string]domain.Detection)
			func() {
				lock := &sync.Mutex{}

				for {
					select {
					case <-scanClose.ctx.Done():
						return
					case detection, ok := <-detections:
						if ok {
							wg.Add(1)
							go func(detection domain.Detection) {
								defer handleRoutinePanic(scanClose.lstream)
								defer wg.Done()

								if dev, err := detection.Device(); err == nil {

									if len(sord(dev.SourceID())) > 0 {
										lock.Lock()
										if devDetectionMap[sord(dev.SourceID())] == nil {
											devDetectionMap[sord(dev.SourceID())] = make(map[string]domain.Detection)
										}

										devDetectionMap[sord(dev.SourceID())][detection.VulnerabilityID()] = detection
										lock.Unlock()
									} else {
										scanClose.lstream.Send(log.Errorf(err, "empty device ID found for detection"))
									}
								} else {
									scanClose.lstream.Send(log.Errorf(err, "error while loading detection %v", detection.VulnerabilityID()))
								}
							}(detection)
						} else {
							return
						}
					}
				}
			}()
			wg.Wait()

			var deadHostIPToProofMap = make(map[string]string)
			func() {
				for {
					select {
					case <-scanClose.ctx.Done():
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

			for _, ticket := range tickets {

				// make sure that we have detections for the device, and that the device isn't listed as dead from the scanner
				if devDetectionMap[ticket.DeviceID()] != nil && !(len(deadHostIPToProofMap[sord(ticket.IPAddress())]) > 0) {

					wg.Add(1)
					go func(ticket domain.Ticket, detection domain.Detection) {
						defer handleRoutinePanic(scanClose.lstream)
						defer wg.Done()

						if sord(ticket.Status()) != engine.GetStatusMap(jira.StatusClosedRemediated) { // TODO: Ensure this is correct
							scanClose.modifyJiraTicketAccordingToVulnerabilityStatus(
								engine,
								&lastCheckedTicket{ticket},
								detection,
								scan,
							)
						}
					}(ticket, devDetectionMap[ticket.DeviceID()][ticket.VulnerabilityID()])

				} else {
					// This block hits if the device did not show up in the scan results. This means the device is either dead, or there were no vulnerabilities found on the device/there was an error in the scan

					wg.Add(1)
					go func(ticket domain.Ticket) {
						defer handleRoutinePanic(scanClose.lstream)
						defer wg.Done()

						// The device appears to be dead
						if len(deadHostIPToProofMap[sord(ticket.IPAddress())]) > 0 {
							if len(sord(ticket.IPAddress())) > 0 {
								if strings.ToLower(scanClose.Payload.Type) == strings.ToLower(domain.RescanDecommission) {
									err = scanClose.closeTicket(engine, ticket, scan, engine.GetStatusMap(jira.StatusClosedDecommissioned), deadHostIPToProofMap[sord(ticket.IPAddress())])
									if err != nil {
										scanClose.lstream.Send(log.Errorf(err, "error while closing dead host ticket", ticket.Title()))
									}
								} else {
									scanClose.lstream.Send(log.Infof("the device for %s seems to be dead, but this is not a decommission scan", ticket.Title()))
									err = engine.Transition(ticket, engine.GetStatusMap(jira.StatusResolvedDecom), fmt.Sprintf("Device is offline. Moving it to a resolved decommission status so a scanner can confirm\nPROOF:\n%s", deadHostIPToProofMap[sord(ticket.IPAddress())]), sord(ticket.AssignedTo()))
									if err != nil {
										scanClose.lstream.Send(log.Errorf(err, "error while adding comment to ticket %s", ticket.Title()))
									}
								}
							}

							// There were no vulnerabilities found on the device or the scan failed to scan the device.
						} else {
							// TODO I believe this block will hit if no vulnerabilities are found on the device
							scanClose.lstream.Send(log.Errorf(err, "scan [%s] did not seem to cover the device %v", scanClose.Payload.ScanID, ticket.DeviceID()))

							//err = engine.Transition(ticket, engine.GetStatusMap(jira.StatusScanError), fmt.Sprintf("scan [%s] did not seem to cover the device %v", scanClose.Payload.ScanID, ticket.DeviceID()), sord(ticket.AssignedTo()))
							//if err != nil {
							//	scanClose.lstream.Send(log.Errorf(err, "error while adding comment to ticket [%s]", ticket.Title()))
							//}

							_, _, _ = engine.UpdateTicket(ticket, fmt.Sprintf("scan [%s] did not cover the device %v", scanClose.Payload.ScanID, ticket.DeviceID()))
						}
					}(ticket)
				}
			}
			wg.Wait()

		} else {
			scanClose.lstream.Send(log.Errorf(err, "error while gathering scan information"))
		}
	} else {
		scanClose.lstream.Send(log.Errorf(err, "error while pulling scan information from scan close payload"))
	}
}

// transitions the status of the JIRA ticket if necessary
func (scanClose *ScanCloseJob) modifyJiraTicketAccordingToVulnerabilityStatus(engine integrations.TicketingEngine, ticket domain.Ticket, detection domain.Detection, scan domain.ScanSummary) {
	var err error
	var status string
	var inactiveKernel bool
	if detection != nil {
		status = detection.Status()
		inactiveKernel = iord(detection.ActiveKernel()) > 0
	}

	// Close Ticket
	if detection == nil || status == domain.Fixed {
		// don't close fixed vulnerabilities during decommission rescans, only vulnerabilities on dead hosts
		if strings.ToLower(scanClose.Payload.Type) != strings.ToLower(domain.RescanDecommission) {
			// Non-decommission scan, the detection appears to be fixed, so close the ticket
			var closeReason = closeComment
			if inactiveKernel {
				closeReason = inactiveKernelComment
			}

			scanClose.lstream.Send(log.Infof("Vulnerability NO LONGER EXISTS, closing Ticket [%s]", ticket.Title()))
			if err = scanClose.closeTicket(engine, ticket, scan, engine.GetStatusMap(jira.StatusClosedRemediated), closeReason); err != nil {
				scanClose.lstream.Send(log.Errorf(err, "Error while closing Ticket [%s]", ticket.Title()))
			}
		} else {
			// Decommission scan - this block hitting means that the host wasn't marked as dead, so we should reopen the associated tickets
			if *ticket.Status() != engine.GetStatusMap(jira.StatusOpen) && *ticket.Status() != engine.GetStatusMap(jira.StatusReopened) { // don't need to reopen tickets that are already opened (these tickets come from the additional ticket loading)
				scanClose.lstream.Send(log.Info(fmt.Sprintf("Ticket [%s] reopened as host is alive", ticket.Title())))
				if err = scanClose.openTicket(ticket, detection, scan, engine); err != nil {
					scanClose.lstream.Send(log.Errorf(err, "Error while opening Ticket [%s]", ticket.Title()))
				}
			}
		}

		// Open ticket
	} else if status == domain.Vulnerable {

		if scanClose.shouldOpenTicket(engine, sord(ticket.Status()), scanClose.Payload.Type) {
			scanClose.lstream.Send(log.Infof("Vulnerability STILL EXISTS, re-opening Ticket [%s]", ticket.Title()))
			if err = scanClose.openTicket(ticket, detection, scan, engine); err != nil {
				scanClose.lstream.Send(log.Errorf(err, "Error while opening Ticket [%s]", ticket.Title()))
			}
		} else {
			scanClose.lstream.Send(log.Infof("Vulnerability still exists, but this is not a standard or decommission rescan - leaving Ticket [%s] in its current status", ticket.Title()))
		}

		// Decommission ticket
	} else if status == domain.DeadHost {
		if strings.ToLower(scanClose.Payload.Type) == strings.ToLower(domain.RescanDecommission) {
			err = scanClose.markDeviceAsDecommissionedInDatabase(ticket)
			if err != nil {
				scanClose.lstream.Send(log.Errorf(err, "error while marking ticket [%s] as decommissioned in the database", ticket.Title()))
			}

			scanClose.lstream.Send(log.Infof("Device %v confirmed offline, closing Ticket [%s]", ticket.DeviceID(), ticket.Title()))

			closeReason := fmt.Sprintf("Device found to be dead by Scanner\n%v", detection.Proof())
			if err = scanClose.closeTicket(engine, ticket, scan, engine.GetStatusMap(jira.StatusClosedDecommissioned), closeReason); err != nil {
				scanClose.lstream.Send(log.Errorf(err, "Error while closing Ticket [%s]", ticket.Title()))
			}
		} else {
			err = engine.Transition(ticket, engine.GetStatusMap(jira.StatusResolvedDecom), fmt.Sprintf("Device is offline. Moving it to a resolved decommission status so a scanner can confirm\nPROOF:\n%s", detection.Proof()), sord(ticket.AssignedTo()))
			if err != nil {
				scanClose.lstream.Send(log.Errorf(err, "error while adding comment to ticket %s", ticket.Title()))
			}
		}

	} else {
		scanClose.lstream.Send(log.Errorf(fmt.Errorf("unrecognized status [%v]", status), "scan close job did not recognize vulnerability status"))
	}
}

func (scanClose *ScanCloseJob) markDeviceAsDecommissionedInDatabase(ticket domain.Ticket) (err error) {
	// we use a sync map to ensure we only decommission a device a single time. if loaded has a value of false, then this is the first time the device was found
	if _, loaded := scanClose.decommedDevices.LoadOrStore(ticket.DeviceID(), true); !loaded {
		scanClose.lstream.Send(log.Infof("Marking device for ticket [%s] as decommissioned in the database", ticket.Title()))
		if _, _, err = scanClose.db.DeleteIgnoreForDevice(
			scanClose.insource.SourceID(),
			ticket.DeviceID(),
			scanClose.config.OrganizationID()); err == nil {

			if _, _, err = scanClose.db.SaveIgnore(
				scanClose.insource.SourceID(),
				scanClose.config.OrganizationID(),
				domain.DecommAsset,
				"", // empty on purpose
				ticket.DeviceID(),
				time.Now(),
				"",
				true,
				""); err != nil {

				scanClose.lstream.Send(log.Errorf(err, "Error while updating exception with asset Id  %s: %s", ticket.DeviceID(), err.Error()))
			}
		} else {
			scanClose.lstream.Send(log.Errorf(err, "error while deleting old entries in the ignore table for the device associated with the ticket [%s]", ticket.Title()))
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

func (scanClose *ScanCloseJob) openTicket(tix domain.Ticket, vuln domain.Detection, scan domain.ScanSummary, ticketing integrations.TicketingEngine) (err error) {
	// Still Exists
	var comment string
	if scanClose.Payload.Type == domain.RescanDecommission {
		comment = fmt.Sprintf("Host is alive according to scan [%s]", sord(scan.SourceKey()))
	} else {

		if len(comment) == 0 {
			comment = fmt.Sprintf("%s\n\nPROOF:\n%s\n\nScan Id [Id: %v]", reopenComment, removeHTMLTags(vuln.Proof()), sord(scan.SourceKey()))
		}
	}

	scanClose.lstream.Send(log.Debugf("TRANSITIONING Ticket [%s]", tix.Title()))
	err = ticketing.Transition(tix, ticketing.GetStatusMap(jira.StatusReopened), comment, "Unassigned")
	if err != nil {
		scanClose.lstream.Send(log.Errorf(err, "Unable to transition ticket %s - %s",
			tix.Title(), err.Error()))
	}
	return err
}

func (scanClose *ScanCloseJob) shouldOpenTicket(ticketing integrations.TicketingEngine, status, jobType string) (shouldOpen bool) {
	status = strings.ToLower(status)
	jobType = strings.ToLower(jobType)

	if status == strings.ToLower(ticketing.GetStatusMap(jira.StatusResolvedRemediated)) && jobType == strings.ToLower(domain.RescanNormal) {
		shouldOpen = true
	} else if status == strings.ToLower(ticketing.GetStatusMap(jira.StatusResolvedDecom)) && jobType == strings.ToLower(domain.RescanDecommission) {
		shouldOpen = true
	}

	return shouldOpen
}

func (scanClose *ScanCloseJob) closeTicket(ticketing integrations.TicketingEngine, tix domain.Ticket, scan domain.ScanSummary, closeStatus string, closeReason string) (err error) {
	// Remediated
	scanClose.lstream.Send(log.Debugf("TRANSITIONING Ticket [%s]", tix.Title()))

	err = ticketing.Transition(tix, closeStatus, fmt.Sprintf("%s\n\nScan [Id: %v]", closeReason, sord(scan.SourceKey())), "Unassigned") //should this be closed?
	if err == nil {
		scanClose.lstream.Send(log.Infof("Ticket [%s] Closed Vulnerability No Longer Exists", tix.Title()))
	} else {
		scanClose.lstream.Send(log.Errorf(err, "Unable to transition ticket %s", tix.Title()))
	}

	return err
}

// we ask the scanning engine to look for a series of vulnerabilities across a series of devices. If we search for a vulnerability on a single
// device in a scan, we check for that vulnerability on all the devices in the scan. here we check to see if the devices have any tickets
// for vulnerabilities that they are being checked for that might not have been included in the original scan
// we only do this check in normal rescans and decommission rescans
func (scanClose *ScanCloseJob) loadAdditionalTickets(ticketing integrations.TicketingEngine, tickets []domain.Ticket) (additionalTickets []domain.Ticket, err error) {
	additionalTickets = make([]domain.Ticket, 0)

	var additionalTicketsChan <-chan domain.Ticket
	if strings.ToLower(scanClose.Payload.Type) == strings.ToLower(domain.RescanDecommission) {
		additionalTicketsChan, err = ticketing.GetAdditionalTicketsForDecomDevices(tickets)
	} else if strings.ToLower(scanClose.Payload.Type) == strings.ToLower(domain.RescanNormal) {
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
func (scanClose *ScanCloseJob) getScanFromPayload() (scan domain.ScanSummary, err error) {
	if len(scanClose.Payload.ScanID) > 0 {
		scan, err = scanClose.db.GetScanSummaryBySourceKey(scanClose.Payload.ScanID)
		if err == nil {
			if scan == nil {
				err = fmt.Errorf("could not find scan in database for %v", scanClose.Payload.ScanID)
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
