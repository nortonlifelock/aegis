package integrations

import (
	"context"
	"fmt"
	"time"

	"github.com/nortonlifelock/aegis/pkg/crypto"
	"github.com/nortonlifelock/aegis/pkg/domain"
	"github.com/nortonlifelock/aegis/pkg/jira"
	"github.com/pkg/errors"
)

// TicketingEngine defines the methods required for the application to interact with a ticketing website, such as JIRA or service now
type TicketingEngine interface {
	CreateTicket(ticket domain.Ticket) (sourceID int, sourceKey string, err error)
	UpdateTicket(ticket domain.Ticket, comment string) (sourceID int, sourceKey string, err error)
	Transition(ticket domain.Ticket, status string, comment string, Assignee string) (err error)

	GetTicket(sourceKey string) (ticket domain.Ticket, err error)
	GetTicketsByClosedStatus(orgCode string, methodOfDiscovery string, startDate time.Time) (tix <-chan domain.Ticket)
	GetTicketsUpdatedSince(since time.Time, orgCode string, methodOfDiscovery string) <-chan domain.Ticket

	GetTicketsForRescan(cerfs []domain.CERF, methodOfDiscovery string, orgCode string, algorithm string) (issues <-chan domain.Ticket, errChan <-chan error)
	GetOpenTicketsByGroupID(methodOfDiscovery string, orgCode string, groupID string) (tickets <-chan domain.Ticket, errChan <-chan error)
	GetRelatedTicketsForRescan(tickets []domain.Ticket, groupID string, methodOfDiscovery string, orgCode string, rescanType string) (issues <-chan domain.Ticket, err error)

	AssignmentGroupExists(groupName string) (exists bool, err error)
	GetStatusMap(backendStatus string) (equivalentTicketStatus string)
	LinkIssues(ticketTitle, linkedTicketTitle string) (err error)
}

const (
	// JIRA identifies the connection as a JIRA connection
	JIRA = "JIRA"
)

// GetEngine returns a struct that implements the TicketingEngine interface
func GetEngine(ctx context.Context, engineID string, db domain.DatabaseConnection, lstream logger, appconfig vulnScannerConfig, config domain.SourceConfig) (eng TicketingEngine, err error) {
	var decryptedConfig domain.SourceConfig
	decryptedConfig, err = crypto.DecryptSourceConfig(db, config, appconfig)

	if err == nil {
		if len(engineID) > 0 {
			switch engineID {

			case JIRA:
				eng, _, err = jira.NewJiraConnector(ctx, lstream, decryptedConfig)
				break
			default:
				err = errors.Errorf("Unknown engine type %s", engineID)
			}
		} else {
			err = errors.New("empty engine id passed to GetEngine")
		}
	} else {
		err = fmt.Errorf("error while decrypting the source config - %s", err.Error())
	}

	return eng, err
}
