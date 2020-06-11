package implementations

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nortonlifelock/aegis/internal/integrations"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/log"
	"strings"
	"sync"
	"time"
)

// RescanQueuePayload is used to parse the Payload from the job history table. The type defines whether it kicks off normal rescans, exception rescans,
// or decommission rescans
type RescanQueuePayload struct {
	Type string `json:"type"`
}

// RescanQueueJob implements the Job interface required to run the job
type RescanQueueJob struct {
	Payload *RescanQueuePayload

	id          string
	payloadJSON string
	ctx         context.Context
	db          domain.DatabaseConnection
	lstream     log.Logger
	appconfig   domain.Config
	config      domain.JobConfig
	insource    domain.SourceConfig
	outsource   domain.SourceConfig
}

// buildPayload parses the information from the Payload of the job history entry
func (job *RescanQueueJob) buildPayload(pjson string) (err error) {

	// Default to a standard rescan queue job
	job.Payload = &RescanQueuePayload{
		Type: domain.RescanNormal,
	}

	// Parse json to RescanPayload
	// Verify pJson length > 0
	if len(pjson) > 0 {
		err = json.Unmarshal([]byte(pjson), job.Payload)
	}

	return err
}

// Process takes tickets that are ready for rescan, grabs their associated groups, and creates job histories for rescans to process those tickets
func (job *RescanQueueJob) Process(ctx context.Context, id string, appconfig domain.Config, db domain.DatabaseConnection, lstream log.Logger, payload string, jobConfig domain.JobConfig, inSource []domain.SourceConfig, outSource []domain.SourceConfig) (err error) {

	var ok bool
	if job.ctx, job.id, job.appconfig, job.db, job.lstream, job.payloadJSON, job.config, job.insource, job.outsource, ok = validInputs(ctx, id, appconfig, db, lstream, payload, jobConfig, inSource, outSource); ok {

		if err = job.buildPayload(job.payloadJSON); err == nil {

			job.lstream = &wrapLogger{
				logger:          job.lstream,
				rescanQueueType: job.Payload.Type,
			}

			//the engine is closed at the end of every loop so each session is ended, but the defer statement ensures a logout in case of panic
			job.lstream.Send(log.Debug("Connecting to Ticketing Engine"))

			var eng integrations.TicketingEngine
			if eng, err = integrations.GetEngine(job.ctx, job.insource.Source(), job.db, job.lstream, job.appconfig, job.insource); err == nil { //the engine is what makes the API calls to JIRA

				var issues <-chan domain.Ticket
				job.lstream.Send(log.Debug("Connection Established to Ticketing Engine"))

				var orgcode string
				if len(job.config.OrganizationID()) > 0 {

					// Get the organization from the database using the id in the ticket object
					var torg domain.Organization
					if torg, err = job.db.GetOrganizationByID(job.config.OrganizationID()); err == nil {
						orgcode = torg.Code()
					}
				}

				job.lstream.Send(log.Debug("Loading tickets for Rescan"))
				var cerfs []domain.CERF
				if cerfs, err = job.db.GetExceptionsDueNext30Days(); err == nil {
					if issues, err = eng.GetTicketsForRescan(cerfs, job.outsource.Source(), orgcode, job.Payload.Type); err == nil { //pulls all resolved remediated tickets from jira
						job.lstream.Send(log.Debugf("[%v] Tickets Loaded for Rescan", len(issues)))

						// exclude issues that are already being processed
						var cleanedIssues <-chan domain.Ticket
						if cleanedIssues, err = job.cleanTickets(issues); err == nil {
							job.processCleanedIssues(cleanedIssues)
						} else {
							job.lstream.Send(log.Error("Error occurred while sorting tickets to re-scan", err))
						}
					} else {
						job.lstream.Send(log.Error("while retrieving tickets", err))
					}
				} else {
					job.lstream.Send(log.Errorf(err, "error while gathering exceptions"))
				}

			} else {
				job.lstream.Send(log.Errorf(err, "Error creating ticketing engine"))
			}
		} else {
			err = fmt.Errorf("error while building payload - %s", err.Error())
		}
	} else {
		err = fmt.Errorf("input validation failed")
	}

	return err
}

func (job *RescanQueueJob) getAssetGroups() (assetGroupBelongsToThisOrgAndScanner map[string]bool, err error) {
	assetGroupBelongsToThisOrgAndScanner = make(map[string]bool)

	var assetGroups []domain.AssetGroup
	if assetGroups, err = job.db.GetAssetGroupForOrg(job.outsource.ID(), job.config.OrganizationID()); err == nil {
		for _, assetGroup := range assetGroups {
			assetGroupBelongsToThisOrgAndScanner[assetGroup.GroupID()] = true
		}
	} else {
		err = fmt.Errorf("error while loading asset groups - %s", err.Error())
	}

	return assetGroupBelongsToThisOrgAndScanner, err
}

func (job *RescanQueueJob) processCleanedIssues(issues <-chan domain.Ticket) {
	var groupIDToTickets = make(map[string][]domain.Ticket)
	var assetGroupBelongsToThisOrgAndScanner map[string]bool
	var err error
	assetGroupBelongsToThisOrgAndScanner, err = job.getAssetGroups()
	if err != nil {
		job.lstream.Send(log.Errorf(err, "error while loading asset groups"))
		return
	}

	var dbWG sync.WaitGroup
	func() {
		for {
			select {
			case <-job.ctx.Done():
				return
			case ticket, ok := <-issues:
				if ok {

					dbWG.Add(1)
					go func(ticket domain.Ticket) {
						defer handleRoutinePanic(job.lstream)
						defer dbWG.Done()

						if deviceInfo, err := job.db.GetDeviceInfoByAssetOrgID(ticket.DeviceID(), job.config.OrganizationID()); err == nil {

							if deviceInfo != nil {
								if deviceInfo.GroupID() != nil {

									groupID := *deviceInfo.GroupID()

									if assetGroupBelongsToThisOrgAndScanner[*deviceInfo.GroupID()] {
										if groupIDToTickets[groupID] == nil {
											groupIDToTickets[groupID] = make([]domain.Ticket, 0)
										}

										groupIDToTickets[groupID] = append(groupIDToTickets[groupID], ticket)
									} else {
										job.lstream.Send(log.Info(fmt.Sprintf("Skipping %s because it belongs to group %s which does not belong to the org/scanner souce config [%s|%s]",
											ticket.Title(), groupID, job.config.OrganizationID(), job.outsource.ID())))
									}

								} else {
									job.lstream.Send(log.Errorf(err, "no group ID found for device info with ID %s in Org %s", ticket.DeviceID(), job.config.OrganizationID()))
								}
							} else {
								job.lstream.Send(log.Errorf(err, "could not find device info with ID %s in Org %s", ticket.DeviceID(), job.config.OrganizationID()))
							}
						} else {
							job.lstream.Send(log.Errorf(err, "error while finding device info with ID %s in Org %s", ticket.DeviceID(), job.config.OrganizationID()))
						}
					}(ticket)

				} else {
					dbWG.Wait()
					return
				}
			}
		}
	}()

	var wg sync.WaitGroup

	for groupID, tickets := range groupIDToTickets {
		select {
		case <-job.ctx.Done():
			return
		default:
			wg.Add(1)
			go func(groupID string, tickets []domain.Ticket) {
				defer handleRoutinePanic(job.lstream)
				defer wg.Done()
				job.processGroup(groupID, tickets)
			}(groupID, tickets)
		}
	}

	wg.Wait()
}

// grabs the hosts from the group, and creates a rescan for each group of 4000 tickets
func (job *RescanQueueJob) processGroup(groupID string, tickets []domain.Ticket) {
	job.lstream.Send(log.Debugf("Processing group %v", groupID))

	var ticketTitles = make([]string, 0)
	for _, ticket := range tickets {
		ticketTitles = append(ticketTitles, ticket.Title())
	}

	job.queueRescan(groupID, ticketTitles)
}

// creates the Payload for the rescan job, and creates a job history for a rescan job to kick off a scan on the provided tickets
func (job *RescanQueueJob) queueRescan(groupID string, tickets []string) {
	defer handleRoutinePanic(job.lstream)

	var payload = &RescanPayload{
		Group:   groupID,
		Tickets: tickets,
		Type:    job.Payload.Type,
	}

	if strPayload, err := json.Marshal(payload); err == nil {

		job.lstream.Send(log.Debugf("Queuing Rescan of Tickets [%s]", strings.Join(tickets, ",")))

		var baseJob domain.JobRegistration
		if baseJob, err = job.db.GetJobsByStruct(rescanJob); err == nil {
			if baseJob != nil {
				job.createJobHistoryForRescanJob(baseJob, strPayload, tickets)
			} else {
				job.lstream.Send(log.Errorf(err, "Empty list of base jobs returned for tickets [%s]", strings.Join(tickets, ",")))
			}
		} else {
			job.lstream.Send(log.Errorf(err, "Error while loading rescan this struct from db for tickets [%s]", strings.Join(tickets, ",")))
		}
	} else {
		job.lstream.Send(log.Errorf(err, "Error occurred while creating rescan Payload for tickets [%s]", strings.Join(tickets, ",")))
	}
}

// creates an entry in the JobHistory table for a rescan job to process the tickets
func (job *RescanQueueJob) createJobHistoryForRescanJob(bjob domain.JobRegistration, strPayload []byte, tickets []string) {
	var err error
	if bjob != nil {

		var configs []domain.JobConfig
		if configs, err = job.db.GetJobConfigByOrgIDAndJobIDWithSC(job.config.OrganizationID(), bjob.ID(), job.outsource.ID()); err == nil { // here we pass the scanner source id so the spawned rescan job uses the same scanner (for the case of an organization using multiple scanners)

			if configs != nil && len(configs) > 0 && configs[0] != nil {

				config := configs[0]

				var priority = bjob.Priority()
				if config.PriorityOverride() != nil {
					priority = iord(config.PriorityOverride())
				}

				_, _, err = job.db.CreateJobHistory(
					bjob.ID(),
					config.ID(),
					domain.JobStatusPending,
					priority,
					"",
					0,
					string(strPayload),
					"",
					time.Now().UTC(),
					"RESCAN QUEUE JOB",
				)

				if err == nil {
					job.lstream.Send(log.Infof("Rescan Queued for Tickets [%s]", strings.Join(tickets, ",")))
				} else {
					job.lstream.Send(log.Errorf(err, "error while queueing rescan for tickets [%s]", strings.Join(tickets, ",")))
				}

			} else {
				job.lstream.Send(log.Errorf(err, "Invalid Config loaded for creating rescan this for tickets [%s]", strings.Join(tickets, ",")))
			}
		} else {
			job.lstream.Send(log.Errorf(err, "Error while loading config from database for tickets [%s]", strings.Join(tickets, ",")))
		}
	} else {
		job.lstream.Send(log.Errorf(err, "Base Rescan Job returned null for tickets [%s]", strings.Join(tickets, ",")))
	}
}

// removes tickets from the input slice that are already being rescanned by another job (either in a Rescan Job, a ScanSync job, or a ScanClose job)
func (job *RescanQueueJob) cleanTickets(tickets <-chan domain.Ticket) (<-chan domain.Ticket, error) {
	var err error
	cleanedTickets := make(chan domain.Ticket)

	var groupIDToGroup = make(map[string]domain.AssetGroup)
	var assetGroups []domain.AssetGroup
	if assetGroups, err = job.db.GetAssetGroupsForOrg(job.config.OrganizationID()); err == nil {
		for _, assetGroup := range assetGroups {
			groupIDToGroup[assetGroup.GroupID()] = assetGroup
		}
	} else {
		err = fmt.Errorf("error while getting asset groups for org [%s]", job.config.OrganizationID())
	}

	var tickMap map[string]bool
	if tickMap, err = job.loadRescans(); err == nil {

		go func() {
			defer handleRoutinePanic(job.lstream)
			defer close(cleanedTickets)

			for {
				if ticket, ok := <-tickets; ok {

					var skipRescanQueue bool
					if len(ticket.GroupID()) > 0 {
						if groupIDToGroup[ticket.GroupID()] != nil {
							skipRescanQueue = groupIDToGroup[ticket.GroupID()].RescanQueueSkip()
						}
					}

					if !tickMap[ticket.Title()] && !skipRescanQueue && job.ticketIsReadyForRescan(ticket) {
						cleanedTickets <- ticket
					} else if skipRescanQueue {
						job.lstream.Send(log.Debugf("skipping queuing of [%s] as group [%s] in the AssetGroup table is marked to skip the RSQ", ticket.Title(), ticket.GroupID()))
					}
				} else {
					break
				}
			}
		}()

	}

	return cleanedTickets, err
}

var ticketTitleOrgIDToUpdatedTime = &sync.Map{}

func (job *RescanQueueJob) ticketIsReadyForRescan(ticket domain.Ticket) (readyForRescan bool) {
	readyForRescan = true // if we can't discern the tracking method, rescan the ticket by default

	if trackingMethod, err := job.db.GetTicketTrackingMethod(ticket.Title(), job.config.OrganizationID()); err == nil && trackingMethod != nil {
		if trackingMethod.Value() == AgentDevice && ticket.UpdatedDate() != nil {

			key := fmt.Sprintf("%s;%s", ticket.Title(), job.config.OrganizationID())

			updatedDate, _ := ticketTitleOrgIDToUpdatedTime.LoadOrStore(key, ticket.UpdatedDate())
			if updatedDateVal, ok := updatedDate.(*time.Time); ok {
				if time.Since(*updatedDateVal) < time.Hour*4 {
					readyForRescan = false

					job.lstream.Send(log.Infof("Skipping rescan of [%s], waiting until [%s] as it is an agent ticket",
						ticket.Title(),
						updatedDateVal.Add(time.Hour*4).Format(time.RFC822)),
					)
				} else {
					// if this path takes, the agent is ready to be rescanned so we delete the timer that was tracking it
					// the next time the ticket is marked for rescan, a new updated date should be written
					ticketTitleOrgIDToUpdatedTime.Delete(key)
				}
			} else {
				job.lstream.Send(log.Errorf(err, "failed to load the updated date [%v] for [%s]", updatedDate, ticket.Title()))
			}
		}
	} else {
		job.lstream.Send(log.Errorf(err, "error while loading tracking method for [%s]", ticket.Title()))
	}

	return readyForRescan
}

// loads tickets that are currently being processed by another job so we don't rescan a ticket that is in the process of being scanned
func (job *RescanQueueJob) loadRescans() (tickets map[string]bool, err error) {
	tickets = make(map[string]bool)
	job.lstream.Send(log.Debugf("Looking for existing rescans for Organization [%v]", job.config.OrganizationID()))

	err = job.checkPendingRescans(tickets)
	if err == nil {
		err = job.checkUnfinishedScans(tickets)
	}

	return tickets, err
}

// discovers tickets that are currently being processed by pending rescan jobs
func (job *RescanQueueJob) checkPendingRescans(tickets map[string]bool) (err error) {
	var jobs []domain.JobHistory
	if jobs, err = job.db.GetPendingActiveRescanJob(job.config.OrganizationID()); err == nil {

		for jid := range jobs {
			if len(jobs[jid].Payload()) > 0 {

				var rp = RescanPayload{}
				if err = json.Unmarshal([]byte(jobs[jid].Payload()), &rp); err == nil {

					for tid := range rp.Tickets {
						tickets[rp.Tickets[tid]] = true
					}
				} else {
					job.lstream.Send(log.Error("error while parsing rescan Payload", err))
				}
			}
		}

	} else {
		err = fmt.Errorf("error occurred while loading pending and active jobs for organization [%v] | [%s]", job.config.OrganizationID(), err)
	}

	return err
}

// checks tickets that are being processed by unfinished scans
func (job *RescanQueueJob) checkUnfinishedScans(tickets map[string]bool) (err error) {
	var scans []domain.ScanSummary
	scans, err = job.db.GetUnfinishedScanSummariesBySourceOrgID(job.outsource.SourceID(), job.config.OrganizationID())
	if err == nil {
		var rp = RescanPayload{}
		for scanIndex := range scans {
			if err = json.Unmarshal([]byte(scans[scanIndex].ScanClosePayload()), &rp); err == nil {
				for tid := range rp.Tickets {
					tickets[rp.Tickets[tid]] = true
				}
			} else {
				job.lstream.Send(log.Error("error while parsing rescan Payload", err))
			}
		}
	} else {
		err = fmt.Errorf("error occurred while loading unfinished scans for organization [%v] | [%s]", job.config.OrganizationID(), err)
	}

	return err
}

type wrapLogger struct {
	logger          log.Logger
	rescanQueueType string
}

// Send wraps the logger so we can add the rescan queue type to all logs
func (l *wrapLogger) Send(log log.Log) {
	if l.rescanQueueType != domain.RescanNormal {
		log.Text = fmt.Sprintf("%s - %s", l.rescanQueueType, log.Text)
	}

	l.logger.Send(log)
}
