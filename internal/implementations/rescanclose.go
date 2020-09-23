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
	Devices []string    `json:"devices"` // this field isn't used in processing, but is useful historical/debugging data
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

			ipsForCloudDecommissionScan := func() <-chan string {
				var ipsForCloudDecommissionScan = make(chan string)

				go func() {
					defer close(ipsForCloudDecommissionScan)

					var wg sync.WaitGroup
					deviceIDToVulnIDToDetection, deadHostIPToProofMap := job.mapDetectionsAndDeadHosts(detections, deadHostIPToProof)

					if len(deviceIDToVulnIDToDetection) > 0 || len(deadHostIPToProof) > 0 || job.Payload.Type == domain.RescanDecommission {
						for _, ticket := range tickets {
							wg.Add(1)
							go func(ticket domain.Ticket) {
								defer handleRoutinePanic(job.lstream)
								defer wg.Done()

								job.modifyJiraTicketAccordingToVulnerabilityStatus(
									engine,
									&lastCheckedTicket{ticket},
									scan,
									deadHostIPToProofMap,
									deviceIDToVulnIDToDetection,
									ipsForCloudDecommissionScan,
								)
							}(ticket)
						}
						wg.Wait()

						job.lstream.Send(log.Infof("Updating detection information in database"))
						job.updateDetectionInformationInDB(deviceIDToVulnIDToDetection)
						job.lstream.Send(log.Infof("Finished updating detection information in database"))
					} else {
						job.lstream.Send(log.Criticalf(err, "NO detections returned from the API, could not process scan results"))
					}

				}()

				return ipsForCloudDecommissionScan
			}()

			var seen = make(map[string]bool)
			var uniques = make([]string, 0)

			func() {
				for {
					select {
					case <-job.ctx.Done():
						return
					case ip, ok := <-ipsForCloudDecommissionScan:
						if ok {
							if !seen[ip] {
								seen[ip] = true
								uniques = append(uniques, ip)
							}
						} else {
							return
						}
					}
				}
			}()

			if len(uniques) > 0 {
				// unique IPs that were found without detections. The channel that populates this slice is only populated during scheduled scans

				// this block of code executes a cloud decommission scan, as scheduled scans are assumed to only run against cloud agents
				// instead of doing traditional decommission scans against cloud agents, we use a CloudDecommissionJob to check the live cloud inventory
				// for active IPs, and close out the tickets for the IPs that don't show in the live asset inventory
				job.lstream.Send(log.Infof("Creating CloudDecommissionJob for [%s]", strings.Join(uniques, ",")))
				createCloudDecommissionJob(job.id, job.db, job.lstream, job.config.OrganizationID(), job.Payload.Group, uniques)
			}
		} else {
			job.lstream.Send(log.Errorf(err, "error while gathering scan information"))
		}
	} else {
		job.lstream.Send(log.Errorf(err, "error while pulling scan information from scan close payload"))
	}
}

func (job *ScanCloseJob) updateDetectionInformationInDB(deviceIDToVulnIDToDetection map[string]map[string]domain.Detection) {
	var vulnCache sync.Map
	var wg sync.WaitGroup

	var permitThread = getPermitThread(10)

	if detectionStatuses, err := job.db.GetDetectionStatuses(); err == nil {
		for deviceID, vulnIDToDetection := range deviceIDToVulnIDToDetection {

			if device, err := job.db.GetDeviceByAssetOrgID(deviceID, job.config.OrganizationID()); err == nil {

				var decommIgnoreID string
				if decommIgnore, err := job.db.HasDecommissioned(deviceID, job.insource.SourceID(), job.config.OrganizationID()); err == nil {
					if decommIgnore != nil {
						decommIgnoreID = decommIgnore.ID()
					}
				}

				for vulnIDAndPort, detection := range vulnIDToDetection {

					vulnID := strings.Split(vulnIDAndPort, domain.VulnPathConcatenator)[0]

					var vulnInfo domain.VulnerabilityInfo
					if vulnInfoInterface, ok := vulnCache.Load(vulnID); ok {
						if vulnInfo, ok = vulnInfoInterface.(domain.VulnerabilityInfo); !ok {
							job.lstream.Send(log.Errorf(fmt.Errorf("cache error while loading vulnerability info"), "cache error"))
						}
					} else {
						vulnInfo, err = job.db.GetVulnInfoBySourceVulnID(vulnID)
						if err == nil && vulnInfo != nil {
							vulnCache.Store(vulnID, vulnInfo)
						} else {
							job.lstream.Send(log.Errorf(err, "cache error during db transaction for [%s]", vulnID))
						}
					}

					if device != nil && vulnInfo != nil && detection != nil {
						select {
						case <-job.ctx.Done():
							return
						case <-permitThread:
						}

						wg.Add(1)
						go func(deviceID string, device domain.Device, vulnInfo domain.VulnerabilityInfo, detection domain.Detection, decommIgnoreID string) {
							defer handleRoutinePanic(job.lstream)
							defer wg.Done()
							defer func() {
								permitThread <- true
							}()

							createOrUpdateDetection(
								job.db, job.lstream,
								detectionStatuses,
								device,
								vulnInfo,
								detection,
								deviceID,
								decommIgnoreID,
								job.config.OrganizationID(),
								job.insource.SourceID(),
								job.getExceptionID,
							)
						}(deviceID, device, vulnInfo, detection, decommIgnoreID)
					} else {
						var missingReason = make([]string, 0)
						if device == nil {
							missingReason = append(missingReason, fmt.Sprintf("failed to load device"))
						}
						if vulnInfo == nil {
							missingReason = append(missingReason, fmt.Sprintf("failed to load vulnerability info"))
						}
						if detection == nil {
							missingReason = append(missingReason, fmt.Sprintf("failed to load detection"))
						}

						job.lstream.Send(log.Errorf(err, "failed to load device|vuln|detection [%v|%v|%v] - [%s]", deviceID, vulnID, detection, strings.Join(missingReason, ",")))
					}
				}
			} else {
				job.lstream.Send(log.Errorf(err, "error while loading device information from db - could not update detection information in db"))
			}
		}
	} else {
		job.lstream.Send(log.Errorf(err, "error while loading detection statuses - could not update detections in db"))
	}

	wg.Wait()
}

func (job *ScanCloseJob) getExceptionID(assetID string, deviceInDb domain.Device, port string, vulnInfo domain.VulnerabilityInfo, detectionFromScanner domain.Detection) (exceptionID string, dontUpdateExceptionID bool) {
	// for the sake of efficiency the rescan close doesn't need to check for exceptions - the asset sync can take care of that
	dontUpdateExceptionID = true

	return exceptionID, dontUpdateExceptionID
}

func createCloudDecommissionJob(parentJobID string, db domain.DatabaseConnection, lstream log.Logger, orgID string, groupID string, ips []string) {
	if assetGroup, err := db.GetAssetGroupForOrgNoScanner(orgID, groupID); err == nil && assetGroup != nil {
		if len(sord(assetGroup.CloudSourceID())) > 0 {
			if scs, err := db.GetSourceConfigBySourceID(orgID, sord(assetGroup.CloudSourceID())); err == nil && len(scs) > 0 {
				var cloudSourceConfig = scs[0]

				if jobRegistration, err := db.GetJobsByStruct(cloudDecomJob); err == nil && jobRegistration != nil {
					if jobConfig, err := db.GetJobConfigByOrgIDAndJobIDWithSC(orgID, jobRegistration.ID(), cloudSourceConfig.ID()); err == nil && len(jobConfig) > 0 {

						var priority = jobRegistration.Priority()
						if jobConfig[0].PriorityOverride() != nil {
							priority = iord(jobConfig[0].PriorityOverride())
						}

						payload := &CloudDecommissionPayload{OnlyCheckIPs: ips}
						if payloadBody, err := json.Marshal(payload); err == nil {
							_, _, err = db.CreateJobHistoryWithParentID(
								jobRegistration.ID(),
								jobConfig[0].ID(),
								domain.JobStatusPending,
								priority,
								"",
								0,
								string(payloadBody),
								"",
								time.Now().UTC(),
								"",
								parentJobID,
							)

							if err == nil {
								lstream.Send(log.Infof("queued a cloud decommission scan for ips [%v]", ips))
							} else {
								lstream.Send(log.Errorf(err, "error while queueing cloud decommission scan for ips [%v]", ips))
							}
						} else {
							lstream.Send(log.Errorf(err, "error while creating payload for CloudDecommissionJob"))
						}
					} else {
						lstream.Send(log.Errorf(err, "error while loading job config for the CloudDecommissionJob"))
					}
				} else {
					lstream.Send(log.Errorf(err, "error while loading CloudDecommissionJob registration from database"))
				}
			} else {
				lstream.Send(log.Errorf(err, "failed to load cloud source config for cloud source ID [%s]", sord(assetGroup.CloudSourceID())))
			}
		} else {
			lstream.Send(log.Errorf(err, "wanted to create a cloud decommission scan for [%s], but it did not have the cloud source ID set", groupID))
		}
	} else {
		lstream.Send(log.Errorf(err, "error while loading asset group information for [org|group|source] [%s|%s]", orgID, groupID))
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

								deviceIDToVulnIDToDetection[sord(dev.SourceID())][combineVulnerabilityIDAndServicePortDetection(detection)] = detection
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

type lastFoundTicket struct {
	domain.Ticket
	lastFound time.Time
}

func (l *lastFoundTicket) AlertDate() *time.Time {
	return &l.lastFound
}

// transitions the status of the JIRA ticket if necessary
func (job *ScanCloseJob) modifyJiraTicketAccordingToVulnerabilityStatus(engine integrations.TicketingEngine, ticket domain.Ticket, scan domain.ScanSummary, deadHostIPToProofMap map[string]string, deviceIDToVulnIDToDetection map[string]map[string]domain.Detection, ipsForCloudDecommissionScan chan<- string) {
	// the scanner found detections for the device
	var detectionsFoundForDevice = deviceIDToVulnIDToDetection[ticket.DeviceID()] != nil

	// device was explicitly reported dead by the scanner
	var deviceReportedAsDead = len(deadHostIPToProofMap[sord(ticket.IPAddress())]) > 0 && len(sord(ticket.IPAddress())) > 0

	// certain types of assets have their host data purged when they are decommissioned (EC2/Agent), so if no detections are found for
	// such devices, it's a good sign that they can be considered dead hosts
	var deviceWithoutDetectionsLikelyDead bool

	var trackingMethod string

	var detection domain.Detection
	if detectionsFoundForDevice {
		detection = deviceIDToVulnIDToDetection[ticket.DeviceID()][combineVulnerabilityIDAndServicePortTicket(ticket)]
	}

	if device, err := job.db.GetDeviceByAssetOrgID(ticket.DeviceID(), job.config.OrganizationID()); err == nil && device != nil {
		trackingMethod = sord(device.TrackingMethod())

		if !detectionsFoundForDevice {
			if job.Payload.Type == domain.RescanScheduled || job.Payload.Type == domain.RescanDecommission { // TODO do we only want to only do this for these two rescan types?
				if trackingMethod == EC2Device || trackingMethod == AgentDevice {
					deviceWithoutDetectionsLikelyDead = true
				}
			}
		}
	} else {
		job.lstream.Send(log.Errorf(err, "error while loading device for %s [%v|%v]", ticket.Title(), err, device))
	}

	var status string
	var inactiveKernel bool
	if detection != nil {
		status = detection.Status()
		inactiveKernel = iord(detection.ActiveKernel()) == 0

		if detection.LastFound() != nil && !detection.LastFound().IsZero() && detection.LastFound().After(tord1970(ticket.AlertDate())) {
			ticket = &lastFoundTicket{
				Ticket:    ticket,
				lastFound: *detection.LastFound(),
			}

			_, _, err := engine.UpdateTicket(ticket, fmt.Sprintf("Updating LastFound field to [%s]", detection.LastFound().String()))
			if err != nil {
				job.lstream.Send(log.Errorf(err, "error while setting last found date of %s to [%s]", ticket.Title(), detection.LastFound().String()))
			}
		}
	}

	if detectionsFoundForDevice || deviceReportedAsDead || deviceWithoutDetectionsLikelyDead {
		switch job.Payload.Type {
		case domain.RescanNormal:
			job.processTicketForNormalRescan(deadHostIPToProofMap, ticket, detection, engine, status, inactiveKernel, scan, trackingMethod)
		case domain.RescanDecommission:
			job.processTicketForDecommRescan(deadHostIPToProofMap, ticket, detection, engine, scan, status, deviceWithoutDetectionsLikelyDead, deviceReportedAsDead, ipsForCloudDecommissionScan, trackingMethod)
		case domain.RescanScheduled:
			job.processTicketForScheduledScan(ticket, detection, engine, status, inactiveKernel, scan, deviceWithoutDetectionsLikelyDead, ipsForCloudDecommissionScan)
		case domain.RescanExceptions, domain.RescanPassive:
			job.processTicketForPassiveOrExceptionRescan(deadHostIPToProofMap, ticket, detection, engine, status, inactiveKernel, scan)
		default:
			job.lstream.Send(log.Critical(fmt.Sprintf("Unrecognized scan type [%s]", job.Payload.Type), nil))
		}
	} else {
		if !statusIsAClosedStatus(engine, sord(ticket.Status())) {
			job.lstream.Send(log.Errorf(nil, "scan [%s] did not seem to cover the device %v - scanner did not report any data for device", job.Payload.ScanID, ticket.DeviceID()))
			err := engine.Transition(ticket, engine.GetStatusMap(domain.StatusScanError), fmt.Sprintf("scan [%s] did not cover the device %v. Please make sure this asset is still in-	scope and associated with an asset group. If this asset is out of scope, please move this ticket to NOTAVRR status or alert the vulnerability management team.", job.Payload.ScanID, ticket.DeviceID()), sord(ticket.AssignedTo()))
			if err != nil {
				job.lstream.Send(log.Errorf(err, "error while marking ticket as ScanError [%s]", ticket.Title()))
			}
		} else {
			job.lstream.Send(log.Debugf("Not setting %s to Scan-Error as it is in a closed status [%s]", ticket.Title(), sord(ticket.Status())))
		}
	}
}

func (job *ScanCloseJob) processTicketForPassiveOrExceptionRescan(deadHostIPToProofMap map[string]string, ticket domain.Ticket, detection domain.Detection, engine integrations.TicketingEngine, status string, inactiveKernel bool, scan domain.ScanSummary) {
	var err error
	if (len(deadHostIPToProofMap[*ticket.IPAddress()]) > 0 && len(*ticket.IPAddress()) > 0) || (detection != nil && detection.Status() == domain.DeadHost) {

		if !statusIsAClosedStatus(engine, sord(ticket.Status())) {
			job.lstream.Send(log.Infof("the device for %s seems to be dead, but this is not a decommission scan", ticket.Title()))
			err = engine.Transition(ticket, engine.GetStatusMap(domain.StatusResolvedDecom), fmt.Sprintf("The device could not be detected though a vulnerability rescan. It has been moved to a resolved decommission status and will be rescanned with another option profile to confirm\nPROOF:\n%s", deadHostIPToProofMap[sord(ticket.IPAddress())]), sord(ticket.AssignedTo()))
			if err != nil {
				job.lstream.Send(log.Errorf(err, "error while marking ticket as ResolvedDecomm %s", ticket.Title()))
			}
		}

	} else if detection == nil || status == domain.Fixed || status == domain.Potential {
		// Non-decommission scan, the detection appears to be fixed, so close the ticket
		var closeReason = closeComment
		if inactiveKernel {
			closeReason = inactiveKernelComment
		} else if status == domain.Potential {
			closeReason = potentialComment
		}

		job.lstream.Send(log.Infof("Vulnerability NO LONGER EXISTS, closing Ticket [%s]", ticket.Title()))
		if err = job.closeTicket(engine, ticket, scan, engine.GetStatusMap(domain.StatusClosedRemediated), closeReason); err != nil {
			job.lstream.Send(log.Errorf(err, "Error while closing Ticket [%s]", ticket.Title()))
		}
	} else if status == domain.Vulnerable {
		// Passive/exception rescans don't reopen tickets - leave alone
	} else {
		job.lstream.Send(log.Errorf(fmt.Errorf("unrecognized status [%v]", status), "scan close job did not recognize vulnerability status"))
	}
}

// TODO informationals to prevent false decoms?
func (job *ScanCloseJob) processTicketForDecommRescan(deadHostIPToProofMap map[string]string, ticket domain.Ticket, detection domain.Detection, engine integrations.TicketingEngine, scan domain.ScanSummary, status string, deviceWithoutDetectionsLikelyDead bool, deviceReportedAsDead bool, ipsForCloudDecommissionScan chan<- string, trackingMethod string) {
	var err error
	if deviceWithoutDetectionsLikelyDead || deviceReportedAsDead || (detection != nil && detection.Status() == domain.DeadHost) {
		err = job.closeTicketAccordingToDeviceType(ticket, detection, deadHostIPToProofMap, engine, scan, ipsForCloudDecommissionScan, trackingMethod)
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

const (
	AgentDevice = "AGENT"
	EC2Device   = "EC2"
	IPDevice    = "IP"
)

func (job *ScanCloseJob) closeTicketAccordingToDeviceType(ticket domain.Ticket, detection domain.Detection, deadHostIPToProofMap map[string]string, engine integrations.TicketingEngine, scan domain.ScanSummary, ipsForCloudDecommissionScan chan<- string, trackingMethod string) (err error) {
	if trackingMethod != EC2Device {
		err = job.markDeviceAsDecommissionedInDatabase(ticket)
		if err != nil {
			job.lstream.Send(log.Errorf(err, "error while marking ticket [%s] as decommissioned in the database", ticket.Title()))
		}

		job.lstream.Send(log.Infof("Device %v confirmed offline, closing Ticket [%s]", ticket.DeviceID(), ticket.Title()))
		var closeReason string
		if detection != nil {
			closeReason = fmt.Sprintf("Device found to be dead by Scanner\n%v", detection.Proof())
		} else {
			closeReason = fmt.Sprintf("Device found to be dead by Scanner\n%v", deadHostIPToProofMap[*ticket.IPAddress()])
		}
		if err = job.closeTicket(engine, ticket, scan, engine.GetStatusMap(domain.StatusClosedDecommissioned), closeReason); err != nil {
			job.lstream.Send(log.Errorf(err, "Error while closing Ticket [%s]", ticket.Title()))
		}
	} else {
		if len(sord(ticket.IPAddress())) > 0 {
			select {
			case <-job.ctx.Done():
				return
			case ipsForCloudDecommissionScan <- sord(ticket.IPAddress()):
				job.lstream.Send(log.Infof("Queueing %s for a cloud decommission scan", ticket.Title()))
			}
		} else {
			job.lstream.Send(log.Errorf(nil, "empty IP on ticket %s", ticket.Title()))
		}
	}

	return err
}

func (job *ScanCloseJob) processTicketForScheduledScan(ticket domain.Ticket, detection domain.Detection, engine integrations.TicketingEngine, status string, inactiveKernel bool, scan domain.ScanSummary, deviceWithoutDetectionsLikelyDead bool, ipsForCloudDecommissionScan chan<- string) {
	var err error

	if deviceWithoutDetectionsLikelyDead {
		// this block assumes that all scheduled scans that Aegis is told to track are returning the results of LIVE CLOUD DEVICES
		// If this block hits, we've found a cloud IP that didn't have any results returned, so we check and see they are live in the cloud service
		// This is done because cloud agent data is wiped after it is decommissioned
		if len(sord(ticket.IPAddress())) > 0 {
			select {
			case <-job.ctx.Done():
				return
			case ipsForCloudDecommissionScan <- sord(ticket.IPAddress()):
				job.lstream.Send(log.Infof("Queueing %s for a cloud decommission job", ticket.Title()))
			}
		} else {
			job.lstream.Send(log.Errorf(err, "empty IP on ticket %s", ticket.Title()))
		}
	} else if detection == nil || status == domain.Fixed {
		// Non-decommission scan, the detection appears to be fixed, so close the ticket
		var closeReason = closeComment
		if inactiveKernel {
			closeReason = inactiveKernelComment
		} else if status == domain.Potential {
			closeReason = potentialComment
		}

		job.lstream.Send(log.Infof("Vulnerability NO LONGER EXISTS, closing Ticket [%s]", ticket.Title()))
		if err = job.closeTicket(engine, ticket, scan, engine.GetStatusMap(domain.StatusClosedRemediated), closeReason); err != nil {
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

func (job *ScanCloseJob) processTicketForNormalRescan(deadHostIPToProofMap map[string]string, ticket domain.Ticket, detection domain.Detection, engine integrations.TicketingEngine, status string, inactiveKernel bool, scan domain.ScanSummary, trackingMethod string) {
	var err error
	if (len(deadHostIPToProofMap[*ticket.IPAddress()]) > 0 && len(*ticket.IPAddress()) > 0) || (detection != nil && detection.Status() == domain.DeadHost) {

		if !statusIsAClosedStatus(engine, sord(ticket.Status())) {
			job.lstream.Send(log.Infof("the device for %s seems to be dead, but this is not a decommission scan", ticket.Title()))
			err = engine.Transition(ticket, engine.GetStatusMap(domain.StatusResolvedDecom), fmt.Sprintf("The device could not be detected though a vulnerability rescan. It has been moved to a resolved decommission status and will be rescanned with another option profile to confirm\nPROOF:\n%s", deadHostIPToProofMap[sord(ticket.IPAddress())]), sord(ticket.AssignedTo()))
			if err != nil {
				job.lstream.Send(log.Errorf(err, "error while marking ticket as ResolvedDecomm %s", ticket.Title()))
			}
		}

	} else if detection != nil && detection.LastUpdated() != nil && detection.LastUpdated().Before(scan.CreatedDate()) && !detection.LastUpdated().IsZero() && !scan.CreatedDate().IsZero() && trackingMethod != AgentDevice {
		if !statusIsAClosedStatus(engine, sord(ticket.Status())) {
			job.lstream.Send(log.Infof("the scan didn't check %s for vulnerability %s [%s before %s]", ticket.Title(), ticket.VulnerabilityID(), detection.LastUpdated().Format(time.RFC822), scan.CreatedDate().Format(time.RFC822)))
			err = engine.Transition(ticket, engine.GetStatusMap(domain.StatusScanError), fmt.Sprintf("The scan did not check the device for the vulnerability likely due to authentication issues"), sord(ticket.AssignedTo()))
			if err != nil {
				job.lstream.Send(log.Errorf(err, "error while transitioning ticket %s", ticket.Title()))
			}
		} else {
			job.lstream.Send(log.Debugf("Not setting %s to Scan-Error as it is in a closed status [%s]", ticket.Title(), sord(ticket.Status())))
		}
	} else if detection == nil || status == domain.Fixed || status == domain.Potential {
		// Non-decommission scan, the detection appears to be fixed, so close the ticket
		var closeReason = closeComment
		if inactiveKernel {
			closeReason = inactiveKernelComment
		} else if status == domain.Potential {
			closeReason = potentialComment
		}

		job.lstream.Send(log.Infof("Vulnerability NO LONGER EXISTS, closing Ticket [%s]", ticket.Title()))
		if err = job.closeTicket(engine, ticket, scan, engine.GetStatusMap(domain.StatusClosedRemediated), closeReason); err != nil {
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
	err = ticketing.Transition(tix, ticketing.GetStatusMap(domain.StatusReopened), comment, "Unassigned")
	if err != nil {
		job.lstream.Send(log.Errorf(err, "Unable to transition ticket %s - %s",
			tix.Title(), err.Error()))
	}
	return err
}

func (job *ScanCloseJob) shouldOpenTicket(ticketing integrations.TicketingEngine, status, jobType string) (shouldOpen bool) {
	status = strings.ToLower(status)
	jobType = strings.ToLower(jobType)

	if status == strings.ToLower(ticketing.GetStatusMap(domain.StatusResolvedRemediated)) && (jobType == strings.ToLower(domain.RescanNormal) || jobType == strings.ToLower(domain.RescanScheduled)) {
		shouldOpen = true
	} else if status == strings.ToLower(ticketing.GetStatusMap(domain.StatusResolvedDecom)) && jobType == strings.ToLower(domain.RescanDecommission) {
		shouldOpen = true
	}

	return shouldOpen
}

func (job *ScanCloseJob) closeTicket(ticketing integrations.TicketingEngine, tix domain.Ticket, scan domain.ScanSummary, closeStatus string, closeReason string) (err error) {
	if !statusIsAClosedStatus(ticketing, sord(tix.Status())) {
		// Remediated
		job.lstream.Send(log.Debugf("TRANSITIONING Ticket [%s]", tix.Title()))

		err = ticketing.Transition(tix, closeStatus, fmt.Sprintf("%s\n\nScan [Id: %v]", closeReason, sord(scan.SourceKey())), "Unassigned") //should this be closed?
		if err == nil {
			job.lstream.Send(log.Infof("Ticket [%s] Closed Vulnerability No Longer Exists", tix.Title()))
		} else {
			job.lstream.Send(log.Errorf(err, "Unable to transition ticket %s", tix.Title()))
		}
	}

	return err
}

// we ask the scanning engine to look for a series of vulnerabilities across a series of devices. If we search for a vulnerability on a single
// device in a scan, we check for that vulnerability on all the devices in the scan. here we check to see if the devices have any tickets
// for vulnerabilities that they are being checked for that might not have been included in the original scan
// we only do this check in normal rescans and decommission rescans
func (job *ScanCloseJob) loadAdditionalTickets(ticketing integrations.TicketingEngine, tickets []domain.Ticket) (additionalTickets []domain.Ticket, err error) {
	additionalTickets = make([]domain.Ticket, 0)

	var orgcode string
	// Get the organization from the database using the id in the ticket object
	var torg domain.Organization
	if torg, err = job.db.GetOrganizationByID(job.config.OrganizationID()); err == nil {
		orgcode = torg.Code()

		var additionalTicketsChan <-chan domain.Ticket
		additionalTicketsChan, err = ticketing.GetRelatedTicketsForRescan(tickets, job.Payload.Group, job.insource.Source(), orgcode, job.Payload.Type)

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

func combineVulnerabilityIDAndServicePortTicket(tic domain.Ticket) (result string) {
	return fmt.Sprintf("%s;%s", tic.VulnerabilityID(), sord(tic.ServicePorts()))
}

func combineVulnerabilityIDAndServicePortDetection(det domain.Detection) (result string) {
	var servicePort string
	if det.Port() > 0 || len(det.Protocol()) > 0 {
		servicePort = fmt.Sprintf("%d %s", det.Port(), det.Protocol())
	}
	return fmt.Sprintf("%s;%s", det.VulnerabilityID(), servicePort)
}

func statusIsAClosedStatus(engine integrations.TicketingEngine, status string) (isClosedStatus bool) {
	var closedStatuses = []string{
		domain.StatusClosedRemediated,
		domain.StatusClosedFalsePositive,
		domain.StatusClosedDecommissioned,
		domain.StatusClosedCerf,
		domain.StatusClosedError,
	}

	for _, closedStatus := range closedStatuses {
		if engine.GetStatusMap(closedStatus) == status {
			isClosedStatus = true
			break
		}
	}

	return isClosedStatus
}
