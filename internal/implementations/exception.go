package implementations

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/nortonlifelock/aegis/internal/integrations"
	"github.com/nortonlifelock/aegis/pkg/domain"
	"github.com/nortonlifelock/aegis/pkg/log"
)

// ExceptionJob is the struct used to run the job, which implements the IJob interface
type ExceptionJob struct {
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

// Process grabs closed tickets for an organization, and either creates an exception in the db if a valid CERF is associated with the ticket, or creates a false
func (job *ExceptionJob) Process(ctx context.Context, id string, appconfig domain.Config, db domain.DatabaseConnection, lstream log.Logger, payload string, jobConfig domain.JobConfig, inSource []domain.SourceConfig, outSource []domain.SourceConfig) (err error) {

	var ok bool
	if job.ctx, job.id, job.appconfig, job.db, job.lstream, job.payloadJSON, job.config, job.insource, job.outsource, ok = validInputs(ctx, id, appconfig, db, lstream, payload, jobConfig, inSource, outSource); ok {

		var eng integrations.TicketingEngine
		if eng, err = integrations.GetEngine(job.ctx, job.insource.Source(), job.db, job.lstream, job.appconfig, job.insource); err == nil {

			// Get organization information
			var orgCode string
			orgCode, err = job.pullOrgCodeFromDB() // the org code is needed to make sure we pull the correct tickets from JIRA
			if err == nil {

				methodOfDiscovery := job.outsource.Source()

				// kick off a thread that pushes closed tickets onto a channel
				var tix = eng.GetTicketsByClosedStatus(orgCode, methodOfDiscovery, tord1970(job.config.LastJobStart()).UTC())

				seen := make(map[string]bool)
				var wg = sync.WaitGroup{}
				func() {

					permit := getPermitThread(100)
					for {

						select {
						case <-ctx.Done():
							return
						case inTicket, ok := <-tix:
							if ok {

								if seen[inTicket.CERF()] {
									wg.Add(1)
									select {
									case <-permit:
									case <-job.ctx.Done():
										return
									}
									go func(ticket domain.Ticket) {
										defer handleRoutinePanic(job.lstream)
										defer wg.Done()
										defer func() {
											select {
											case permit <- true:
											case <-job.ctx.Done():
											}
										}()

										processExceptionOrFalsePositive(job.db, eng, job.lstream, job.config.OrganizationID(), job.outsource.SourceID(), ticket)
									}(inTicket)
								} else {
									seen[inTicket.CERF()] = true
									job.lstream.Send(log.Debugf("Slow loading %s", inTicket.CERF()))
									processExceptionOrFalsePositive(job.db, eng, job.lstream, job.config.OrganizationID(), job.outsource.SourceID(), inTicket)
									job.lstream.Send(log.Debugf("Done loading %s", inTicket.CERF()))
								}

							} else {
								return
							}
						}
					}
				}()

				wg.Wait()
			} else {
				job.lstream.Send(log.Error("error while gathering organization code from the database", err))
			}

			// TODO this requires further testing - this is only required for detections that are no longer synced which have expired ignores
			// TODO Synced detections have expired ignores removed
			_, _, err = job.db.RemoveExpiredIgnoreIDs(job.config.OrganizationID())
			if err != nil {
				job.lstream.Send(log.Errorf(err, "error while deleting outdated IgnoreIDs"))
			}
		} else {
			job.lstream.Send(log.Error("Error while creating ticketing connection", err))
		}

	} else {
		err = fmt.Errorf("input validation failed")
	}

	return err
}

// grabs the associated org code from the database using the organization id
func (job *ExceptionJob) pullOrgCodeFromDB() (orgcode string, err error) {
	if len(job.config.OrganizationID()) > 0 {

		// Get the organization from the database using the id in the ticket object
		var torg domain.Organization
		if torg, err = job.db.GetOrganizationByID(job.config.OrganizationID()); err == nil {
			orgcode = torg.Code()
		}
	}

	return orgcode, err
}

// This method creates an exception in the database if there is an associated CERF with the ticket that has not expired
// If there is not an associated CERF, a false positive entry in the database is created
func processExceptionOrFalsePositive(db domain.DatabaseConnection, engine integrations.TicketingEngine, lstream log.Logger, orgID, sourceID string, ticket domain.Ticket) {
	var err error

	var deviceID = ticket.DeviceID()
	var vulnID = ticket.VulnerabilityID()
	var ignoreSaved bool

	if sord(ticket.Status()) == engine.GetStatusMap(domain.StatusApprovedException) {
		if len(ticket.CERF()) > 0 && ticket.CERF() != "Empty" {
			if ticket.ExceptionExpiration().After(time.Now()) {

				lstream.Send(log.Infof("Creating/updating EXCEPTION %s", ticket.Title()))

				if _, _, err = db.SaveIgnore(
					sourceID,
					orgID,
					domain.Exception,
					vulnID,
					deviceID,
					ticket.ExceptionExpiration(),
					ticket.CERF(),
					true,
					sord(ticket.ServicePorts())); err == nil {
					ignoreSaved = true
				} else {
					lstream.Send(log.Errorf(err, "Error while updating ticket %s: %s", ticket.Title(), err.Error()))
				}
			} else {
				lstream.Send(log.Debugf("Skipping update for %s as it's CERF expired in the past (%s)", ticket.Title(), ticket.ExceptionExpiration().Format(time.RFC3339)))
			}
		} else {
			lstream.Send(log.Errorf(err, "empty CERF found on [%s]", ticket.Title()))
		}

	} else if sord(ticket.Status()) == engine.GetStatusMap(domain.StatusClosedFalsePositive) {

		// TODO: update the due date to be able to be passed as null to the sproc
		lstream.Send(log.Infof("Creating/updating FALSE POSITIVE %s", ticket.Title()))
		t := time.Date(1111, 1, 1, 1, 1, 0, 1, time.UTC)
		if _, _, err = db.SaveIgnore(
			sourceID,
			orgID,
			domain.FalsePositive,
			vulnID,
			deviceID,
			t,
			ticket.Title(),
			true,
			sord(ticket.ServicePorts())); err == nil {
			ignoreSaved = true
		} else {
			lstream.Send(log.Errorf(err, "Error while updating ticket %s: %s", ticket.Title(), err.Error()))
		}
	} else {
		lstream.Send(log.Errorf(err, "unrecognized status for exception processing [%s]", sord(ticket.Status())))
	}

	_ = ignoreSaved
	//if ignoreSaved {
	//	var ignore domain.Ignore
	//	if ignore, err = db.HasIgnore(
	//		sourceID,
	//		vulnID,
	//		deviceID,
	//		orgID,
	//		sord(ticket.ServicePorts()),
	//		tord1970(nil),
	//	); err == nil {
	//		if ignore != nil {
	//			var portString string
	//			var protocol string
	//			var portInt int
	//
	//			if len(sord(ticket.ServicePorts())) > 0 {
	//				var portProtocol = strings.Split(sord(ticket.ServicePorts()), " ")
	//				if len(portProtocol) == 2 {
	//					portString = portProtocol[0]
	//					protocol = portProtocol[1]
	//					if portInt, err = strconv.Atoi(portString); err != nil {
	//						lstream.Send(log.Errorf(err, "failed to parse port [%s] as integer", portString))
	//					}
	//				} else {
	//					err = fmt.Errorf("port formatting error")
	//					lstream.Send(log.Errorf(err, "[%s] could not be broken into two", sord(ticket.ServicePorts())))
	//				}
	//			}
	//
	//			if err == nil {
	//				_, _, err = db.UpdateDetectionIgnore(deviceID, vulnID, portInt, protocol, ignore.ID())
	//				if err != nil {
	//					lstream.Send(log.Errorf(err, "error while updating ignore for [%s/%s]", deviceID, vulnID))
	//				}
	//
	//				lstream.Send(log.Infof("finished updating detection for %v/%v/%v/%v", deviceID, vulnID, portInt, protocol))
	//			}
	//
	//		} else {
	//			lstream.Send(log.Errorf(err, "failed to load ignore entry for [%s/%s]", deviceID, vulnID))
	//		}
	//	} else {
	//		lstream.Send(log.Errorf(err, "error while loading ignore entry for [%s/%s]", deviceID, vulnID))
	//	}
	//}
}
