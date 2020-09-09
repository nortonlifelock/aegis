package implementations

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/nortonlifelock/aegis/internal/integrations"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/log"
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

				var wg = sync.WaitGroup{}
				func() {

					permit := getPermitThread(100)
					for {

						select {
						case <-ctx.Done():
							return
						case inTicket, ok := <-tix:
							if ok {
								select {
								case <-permit:
								case <-job.ctx.Done():
									return
								}
								wg.Add(1)
								go func(ticket domain.Ticket) {
									defer handleRoutinePanic(job.lstream)
									defer wg.Done()
									defer func() {
										select {
										case permit <- true:
										case <-job.ctx.Done():
										}
									}()

									job.processExceptionOrFalsePositive(ticket)
								}(inTicket)
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

			job.updateCERFExpirationsInDB(eng)

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
func (job *ExceptionJob) processExceptionOrFalsePositive(ticket domain.Ticket) {
	var err error

	var deviceID = ticket.DeviceID()
	var vulnID = ticket.VulnerabilityID()
	var ignoreSaved bool

	if len(ticket.CERF()) > 0 && ticket.CERF() != "Empty" {

		// TODO: update the due date to be able to be passed as null to the sproc
		if ticket.CERFExpirationDate().After(time.Now()) {

			job.lstream.Send(log.Infof("Creating/updating EXCEPTION %s", ticket.Title()))

			if _, _, err = job.db.SaveIgnore(
				job.outsource.SourceID(),
				job.config.OrganizationID(),
				domain.Exception,
				vulnID,
				deviceID,
				ticket.CERFExpirationDate(),
				ticket.CERF(),
				true,
				sord(ticket.ServicePorts())); err == nil {
				ignoreSaved = true
			} else {
				job.lstream.Send(log.Errorf(err, "Error while updating ticket %s: %s", ticket.Title(), err.Error()))
			}
		} else {
			job.lstream.Send(log.Debugf("Skipping update for %s as it's CERF expired in the past (%s)", ticket.CERFExpirationDate().Format(time.RFC3339)))
		}
	} else {

		// TODO: update the due date to be able to be passed as null to the sproc
		job.lstream.Send(log.Infof("Creating/updating FALSE POSITIVE %s", ticket.Title()))
		t := time.Date(1111, 1, 1, 1, 1, 0, 1, time.UTC)
		if _, _, err = job.db.SaveIgnore(
			job.outsource.SourceID(),
			job.config.OrganizationID(),
			domain.FalsePositive,
			vulnID,
			deviceID,
			t,
			ticket.Title(),
			true,
			sord(ticket.ServicePorts())); err == nil {
			ignoreSaved = true
		} else {
			job.lstream.Send(log.Errorf(err, "Error while updating ticket %s: %s", ticket.Title(), err.Error()))
		}
	}

	_ = ignoreSaved
	//if ignoreSaved {
	//	var ignore domain.Ignore
	//	if ignore, err = job.db.HasIgnore(
	//		job.outsource.SourceID(),
	//		vulnID,
	//		deviceID,
	//		job.config.OrganizationID(),
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
	//						job.lstream.Send(log.Errorf(err, "failed to parse port [%s] as integer", portString))
	//					}
	//				} else {
	//					err = fmt.Errorf("port formatting error")
	//					job.lstream.Send(log.Errorf(err, "[%s] could not be broken into two", sord(ticket.ServicePorts())))
	//				}
	//			}
	//
	//			if err == nil {
	//				_, _, err = job.db.UpdateDetectionIgnore(deviceID, vulnID, portInt, protocol, ignore.ID())
	//				if err != nil {
	//					job.lstream.Send(log.Errorf(err, "error while updating ignore for [%s/%s]", deviceID, vulnID))
	//				}
	//
	//				job.lstream.Send(log.Infof("finished updating detection for %v/%v/%v/%v", deviceID, vulnID, portInt, protocol))
	//			}
	//
	//		} else {
	//			job.lstream.Send(log.Errorf(err, "failed to load ignore entry for [%s/%s]", deviceID, vulnID))
	//		}
	//	} else {
	//		job.lstream.Send(log.Errorf(err, "error while loading ignore entry for [%s/%s]", deviceID, vulnID))
	//	}
	//}
}

// This method updates the expiration date of the CERFs in the database that are past the date of the last job start
func (job *ExceptionJob) updateCERFExpirationsInDB(eng integrations.TicketingEngine) {
	var err error
	var cerfUpdates map[string]time.Time
	// Handle updated CERF tickets
	if cerfUpdates, err = eng.GetCERFExpirationUpdates(tord1970(job.config.LastJobStart())); err == nil {

		if len(cerfUpdates) > 0 {

			for key := range cerfUpdates {
				if !cerfUpdates[key].IsZero() {
					job.lstream.Send(log.Infof("Updating expiration date for exceptions with [%s] to [%s]", key, cerfUpdates[key].Format(time.RFC1123Z)))
					if _, _, err = job.db.UpdateExpirationDateByCERF(key, job.config.OrganizationID(), cerfUpdates[key]); err != nil {
						job.lstream.Send(log.Error("error while updating cerf expiration date", err))
					}
				} else {
					// TODO do we want to skip the update? or set the expiration date to a value in the past?
					job.lstream.Send(log.Infof("Skipping expiration date update for [%s] as it was not set in JIRA", key))
				}
			}

		}
	} else {
		job.lstream.Send(log.Error("Error when loading CERF expiration updates from JIRA", err))
	}
}
