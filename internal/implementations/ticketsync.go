package implementations

import (
	"context"
	"fmt"
	"github.com/nortonlifelock/aegis/internal/integrations"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/log"
	"strconv"
	"strings"
	"sync"
)

// TicketSyncJob pulls ticket information from an engine and stores it in the database
type TicketSyncJob struct {
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

// Process pulls tickets from JIRA that have been updated since the last job run, and stores the updated information in the database
func (job *TicketSyncJob) Process(ctx context.Context, id string, appconfig domain.Config, db domain.DatabaseConnection, lstream log.Logger, payload string, jobConfig domain.JobConfig, inSource []domain.SourceConfig, outSource []domain.SourceConfig) (err error) {
	var ok bool
	if job.ctx, job.id, job.appconfig, job.db, job.lstream, job.payloadJSON, job.config, job.insource, job.outsource, ok = validInputs(ctx, id, appconfig, db, lstream, payload, jobConfig, inSource, outSource); ok {
		var engine integrations.TicketingEngine
		if engine, err = integrations.GetEngine(job.ctx, job.insource.Source(), job.db, job.lstream, job.appconfig, job.insource); err == nil {

			var org domain.Organization
			if org, err = job.db.GetOrganizationByID(job.config.OrganizationID()); err == nil {
				tics := engine.GetTicketsUpdatedSince(tord1970(job.config.LastJobStart()), org.Code(), job.outsource.Source())
				job.updateTicketsInDB(tics)
			} else {
				err = fmt.Errorf("error while loading organization information - %s", err.Error())
			}

		} else {
			err = fmt.Errorf("error while creating ticketing engine - %v", err.Error())
		}
	} else {
		err = fmt.Errorf("input validation failed")
	}

	return err
}

func (job *TicketSyncJob) updateTicketsInDB(tics <-chan domain.Ticket) {
	done := make(chan interface{})
	defer close(done)

	go func() {
		orgID := job.config.OrganizationID()
		wg := &sync.WaitGroup{}
		for {
			select {
			case <-job.ctx.Done():
				return
			case tic, ok := <-tics:
				if ok {
					wg.Add(1)
					go func(tic domain.Ticket) {
						defer wg.Done()
						job.processTicket(tic, orgID)
					}(tic)
				} else {
					wg.Wait()
					done <- true
					return
				}
			}
		}
	}()

	select {
	case <-job.ctx.Done():
	case <-done:
	}
}

func (job *TicketSyncJob) processTicket(tic domain.Ticket, orgID string) {
	existingDBTicket, err := job.db.GetTicketByTitle(tic.Title(), orgID)
	if err == nil {
		if existingDBTicket == nil {

			var portString string
			var protocol string
			var portInt int

			if len(sord(tic.ServicePorts())) > 0 {
				var portProtocol = strings.Split(sord(tic.ServicePorts()), " ")
				if len(portProtocol) == 2 {
					portString = portProtocol[0]
					protocol = portProtocol[1]
					if portInt, err = strconv.Atoi(portString); err != nil {
						job.lstream.Send(log.Errorf(err, "failed to parse port [%s] as integer", portString))
					}
				} else {
					err = fmt.Errorf("port formatting error")
					job.lstream.Send(log.Errorf(err, "[%s] could not be broken into two", sord(tic.ServicePorts())))
				}
			}

			if err == nil {
				var detection domain.Detection
				if detection, err = job.getDetection(tic.DeviceID(), tic.VulnerabilityID(), portInt, protocol); err == nil {

					_, _, err = job.db.CreateTicket(
						tic.Title(),
						sord(tic.Status()),
						detection.ID(),
						job.config.OrganizationID(),
						tord1970(tic.DueDate()),
						tord1970(tic.CreatedDate()),
						tord1970(tic.UpdatedDate()),
						tord1970(tic.ResolutionDate()),
						tord1970(nil), // used to set the resolution date to nil in the DB if the ticket doesn't have one
					)

					if err != nil {
						job.lstream.Send(log.Errorf(err, "error while creating database ticket for [%v]", tic.Title()))
					}
				} else {
					job.lstream.Send(log.Errorf(err, "error while loading detection for [%v]", tic.Title()))
				}
			}

		} else {
			_, _, err = job.db.UpdateTicket(
				tic.Title(),
				sord(tic.Status()),
				job.config.OrganizationID(),
				sord(tic.AssignmentGroup()),
				sord(tic.AssignedTo()),
				tord1970(tic.CreatedDate()),
				tord1970(tic.UpdatedDate()),
				tord1970(tic.ResolutionDate()),
				tord1970(nil), // used to set the resolution date to nil in the DB if the ticket doesn't have one
			)

			if err != nil {
				job.lstream.Send(log.Errorf(err, "error while updating database ticket [%v]", tic.Title()))
			}
		}
	} else {
		job.lstream.Send(log.Errorf(err, "error while loading existing ticket from database for [%v]", tic.Title()))
	}
}

func (job *TicketSyncJob) getDetection(deviceID string, vulnID string, port int, protocol string) (detection domain.Detection, err error) {
	if len(deviceID) > 0 && len(vulnID) > 0 {
		if detection, err = job.db.GetDetectionBySourceVulnID(deviceID, vulnID, port, protocol); err == nil {
			if detection == nil {
				err = fmt.Errorf("could not find detection for [%v|%v]", deviceID, vulnID)
			}
		}
	} else {
		err = fmt.Errorf("one of the following were not present: device id, vulnerability id")
	}

	return detection, err
}

// TODO convert tickets status to a status normalized for our database
func (job *TicketSyncJob) translateStatus(inStatus string) (outStatus string) {
	outStatus = inStatus
	switch strings.ToLower(inStatus) {

	}

	return outStatus
}
