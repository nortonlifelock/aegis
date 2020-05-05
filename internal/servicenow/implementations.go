package servicenow

import (
	"fmt"
	"github.com/nortonlifelock/domain"
	"github.com/pkg/errors"
	"net/http"
	"strings"
	"sync"
	"time"
)

// API Endpoints
const (
	//SEARCH_API         = "/rest/api/latest/search"
	tableName         = "sn_vul_vulnerable_item"
	orgCodeField      = "skills="
	stateField        = "state="
	orderByDesc       = "ORDERBYDESC"
	sysUpdated        = "sys_updated>"
	dueDate           = "due_date"
	methodOfDiscField = "user_input="
	vulnIDField       = "additional_assignee_list="
	devIDField        = "qualys_ticket="
)

// CreateTicket takes a domain.Ticket and creates a ticket using it's information in Service Now
func (connector *SvcNowConnector) CreateTicket(Ticket domain.Ticket) (SourceID int, SourceKey string, err error) {
	if Ticket != nil {
		var svcNowTicketBody *SvcNowRequest
		if svcNowTicketBody, err = mapDalTicketToSvcNowIssue(connector, Ticket); err == nil {
			if svcNowTicketBody != nil {
				var req *http.Request
				if req, err = connector.client.newSvcNowRequest(tableName, http.MethodPost, "", "", nil, 0, svcNowTicketBody); err == nil {
					var svcNowTicket SvcNowTicketDto
					var echeck Err

					var response *http.Response
					if response, err = connector.client.client.Do(req); err == nil {
						if response != nil {
							defer response.Body.Close()
						}

						_, _ = svcNowTicket, echeck

						//ticket := response.(*SvcNowTicketDto)
						//SourceKey = ticket.Result.SysID
					}
				}
			} else {
				err = errors.New("SeviceNow Issue failed to map from Ticket Create Ticket")
			}
		}
	} else {
		err = errors.New("ServiceNow Ticket was passed null to Create Ticket")
	}

	return SourceID, SourceKey, err
}

// UpdateTicket yet to be implemented
func (connector *SvcNowConnector) UpdateTicket(Ticket domain.Ticket, comment string) (SourceID int, SourceKey string, err error) {
	//it will be needed for the jiratool
	err = errors.New("Method not implemented")
	return SourceID, SourceKey, err
}

// Transition changes a ticket's status as defined in the parameters
func (connector *SvcNowConnector) Transition(ticket domain.Ticket, status string, comment string, Assignee string) (err error) {
	if ticket != nil {
		//Temporary until we get the status workflow implemented in svcnow
		if status == "Closed-Remediated" {
			status = "New"
		} else if status == "Reopened" {
			status = "New"
		} else {
			status = "some other status"
		}
		var issue domain.Ticket
		if issue, err = connector.GetTicket(ticket.Title()); err == nil {
			if issue != nil {
				body := &SvcNowRequest{
					WorkNotes: comment, //adding the comment to this field because "comments" field is not adding the value
					State:     status,
				}
				var req *http.Request
				if req, err = connector.client.newSvcNowRequest(tableName, http.MethodPut, issue.Title(), "", nil, 0, body); err == nil {
					var svcNowTicket SvcNowTicketDto
					var echeck Err
					var response *http.Response
					if response, err = connector.client.client.Do(req); err == nil {
						if response != nil {
							defer response.Body.Close()
						}

						_, _ = svcNowTicket, echeck
					}
				}
			} else {
				err = errors.New(fmt.Sprintf("Unable to find ticket %s", ticket.Title()))
			}
		}

	} else {
		err = errors.New("Ticket passed null for transition")
	}
	return err
}

// GetTicket returns the ticket with a corresponding SourceKey found in the parameter
func (connector *SvcNowConnector) GetTicket(SourceKey string) (ticket domain.Ticket, err error) {
	if len(SourceKey) > 0 {
		var req *http.Request
		if req, err = connector.client.newSvcNowRequest(tableName, http.MethodGet, SourceKey, "", nil, 0, nil); err == nil {
			var svcNowTicket SvcNowTicketDto
			var echeck Err

			var response *http.Response
			if response, err = connector.client.client.Do(req); err == nil {
				if response != nil {
					defer response.Body.Close()
				}

				_, _ = svcNowTicket, echeck

				//issue := response.(*SvcNowTicketDto)
				//ticket, err = mapSvcNowIssueToDalTicket(connector, issue.Result)
			}
		}

	} else {
		err = errors.New("Invalid Source Key in GetTicket")
	}

	return ticket, err
}

// GetTicketsUpdatedSince not implemented
func (connector *SvcNowConnector) GetTicketsUpdatedSince(since time.Time) <-chan domain.Ticket {
	return nil
}

// GetTicketsByClosedStatus gets all closed tickets and pushes them on the tix channel
func (connector *SvcNowConnector) GetTicketsByClosedStatus(orgCode string, startDate time.Time) <-chan domain.Ticket {
	tix := make(chan domain.Ticket)

	go func(tix chan<- domain.Ticket) {
		defer close(tix)

		var issues = &SvcNowTickets{}
		var err error
		err = connector.getTicketsByClosedStatus(orgCode, startDate, issues)

		if err == nil {

			if issues.Results != nil {
				var wg = sync.WaitGroup{}

				for index := range issues.Results {
					wg.Add(1)

					go func(svcIssue *Result) {
						defer wg.Done()

						var ticket domain.Ticket
						if ticket, err = mapSvcNowIssueToDalTicket(connector, svcIssue); err == nil {
							tix <- ticket
						} else {
							fmt.Println(err)
						}

					}(issues.Results[index])

				}
				wg.Wait()
			}
		}
	}(tix)

	return tix
}

func (connector *SvcNowConnector) getTicketsByClosedStatus(orgCode string, startDate time.Time, issues *SvcNowTickets) (err error) {

	query := createQuery([]string{orgCodeField, stateField, sysUpdated}, orgCode, "6", startDate.Format(svcNowDateFormat))
	err = connector.client.performFor(tableName, http.MethodGet, "", query, nil, 0, nil, issues)

	return err
}

// GetCERFExpirationUpdates not implemented
func (connector *SvcNowConnector) GetCERFExpirationUpdates(startDate time.Time) (cerfs map[string]time.Time, err error) {
	err = errors.New("Method not implemented")
	return cerfs, err
}

// GetTicketsByDeviceIDVulnID returns all tickets matched the device id and the vulnerability id
func (connector *SvcNowConnector) GetTicketsByDeviceIDVulnID(MethodOfDiscovery string, OrgCode string, DeviceID string, VulnID string, Statuses map[string]bool, port int, protocol string) (Issues <-chan domain.Ticket, err error) {
	if len(DeviceID) > 0 {
		if len(VulnID) > 0 {
			//Make sure the fields name and fields values are in same order respectively
			query := createQuery([]string{methodOfDiscField, orgCodeField, devIDField, vulnIDField}, MethodOfDiscovery, OrgCode, DeviceID, VulnID)
			var req *http.Request
			if req, err = connector.client.newSvcNowRequest(tableName, http.MethodGet, "", query, nil, 0, nil); err == nil {
				var svcNowTicket SvcNowTickets
				var echeck Err

				var response *http.Response
				if response, err = connector.client.client.Do(req); err == nil {
					if response != nil {
						defer response.Body.Close()
					}

					_, _ = svcNowTicket, echeck

					//tickets := response.(*SvcNowTickets)
					//Issues = mapToDalTickets(connector, *tickets)
				}
			}

		} else {
			err = errors.New("Vulnerability id passed as 0 to getTicketsByDeviceIdVulnId")
		}
	} else {
		err = errors.New("Device id passed as 0 to getTicketsByDeviceIdVulnId")
	}

	return Issues, err
}

func (connector *SvcNowConnector) getSysIDFor(field string, value string, table string) (sysID string, err error) {
	if field != "" && table != "" {
		field += "="
		query := createQuery([]string{field}, value)
		var req *http.Request
		if req, err = connector.client.newSvcNowRequest(table, http.MethodGet, "", query, nil, 0, nil); err == nil {
			var svcNowTicket SvcNowTickets
			var echeck Err

			var response *http.Response
			if response, err = connector.client.client.Do(req); err == nil {
				if response != nil {
					defer response.Body.Close()
				}

				_, _ = svcNowTicket, echeck
				//tickets := response.(*SvcNowTickets)
				//if len(tickets.Results) > 0 {
				//	sysID = tickets.Results[0].SysID
				//}
			}
		}
	}
	return
}

// GetOpenTicketsDueSoon not implemented
func (connector *SvcNowConnector) GetOpenTicketsDueSoon(MethodOfDiscovery string, OrgCode string, Statuses *map[string]bool) (Issues <-chan domain.Ticket, err error) {
	err = errors.New("Method not implemented")
	return Issues, err
}

// GetAdditionalTicketsForVulnPerDevice not implemented
func (connector *SvcNowConnector) GetAdditionalTicketsForVulnPerDevice(tickets []domain.Ticket) (Issues <-chan domain.Ticket, err error) {
	err = errors.New("Method not implemented")
	return Issues, err
}

// GetAdditionalTicketsForDecomDevices not implemented
func (connector *SvcNowConnector) GetAdditionalTicketsForDecomDevices(tickets []domain.Ticket) (Issues <-chan domain.Ticket, err error) {
	err = errors.New("Method not implemented")
	return Issues, err
}

// GetTicketsForRescan returns a list of tickets in a particular status depending on the type of rescan
func (connector *SvcNowConnector) GetTicketsForRescan(MethodOfDiscovery string, OrgCode string, Algorithm string) (Issues <-chan domain.Ticket, err error) {
	if len(MethodOfDiscovery) > 0 {
		if len(OrgCode) > 0 {
			switch Algorithm {
			case domain.RescanExceptions:
				var states = []string{"3", "4"}
				Issues, err = connector.getTicketsForExceptionRescan(MethodOfDiscovery, OrgCode, states)
				break
			case domain.RescanPassive:
				var states = []string{"7", "8"}
				Issues, err = connector.getTicketsForPassiveRescan(MethodOfDiscovery, OrgCode, states)
				break
			case domain.RescanNormal:
				Issues, err = connector.getTicketsByStatusDueDateAscending(MethodOfDiscovery, OrgCode, nil)
				break
			default:
				// TODO: Error
				break
			}
		} else {
			// TODO: Error
		}
	} // TODO: Error

	return Issues, err
}

func (connector *SvcNowConnector) getTicketsByStatusDueDateAscending(MethodOfDiscovery string, OrgCode string, Statuses map[string]bool) (Issues <-chan domain.Ticket, err error) {
	if len(MethodOfDiscovery) > 0 {
		if len(OrgCode) > 0 {
			//Make sure the fields name and fields vaules are in same order respectively
			// I am hardcoding the status until we implement the Status workflow.
			query := createQuery([]string{methodOfDiscField, orgCodeField, stateField, orderByDesc}, MethodOfDiscovery, OrgCode, "2", dueDate)
			var req *http.Request
			if req, err = connector.client.newSvcNowRequest(tableName, http.MethodGet, "", query, nil, 0, nil); err == nil {
				var svcNowTicket SvcNowTickets
				var echeck Err

				var response *http.Response
				if response, err = connector.client.client.Do(req); err == nil {
					if response != nil {
						defer response.Body.Close()
					}

					_, _ = svcNowTicket, echeck

					//tickets := response.(*SvcNowTickets)
					//Issues = mapToDalTickets(connector, *tickets)
				}
			}

		} else {
			err = errors.New("No orgID is provided")
		}
	} else {
		err = errors.New("No Method of discovery is provided")
	}
	return Issues, err
}

// addComment creates a new comment on a service now ticket
func (connector *SvcNowConnector) addComment(ticket domain.Ticket, comment string) (err error) {
	if ticket != nil {
		var issue = &SvcNowTicketDto{}
		if err = connector.client.performFor(tableName, http.MethodGet, ticket.Title(), "", nil, 0, nil, issue); err == nil {

			body := &SvcNowRequest{
				WorkNotes: comment, //adding the comment to this field because "comments" field is not adding the value
			}
			err = connector.client.performFor(tableName, http.MethodPut, issue.Result.SysID, "", nil, 0, body, &issue)

		}

	} else {
		err = fmt.Errorf("empty ticket passed to addComment")
	}
	return err
}

// AssignmentGroupExists not implemented
func (connector *SvcNowConnector) AssignmentGroupExists(groupName string) (exists bool, err error) {
	err = errors.New("Method not implemented")
	return exists, err
}

func createQuery(fieldsNames []string, fieldsValues ...string) (modifiedQuery string) {
	var valuePair string
	for index := range fieldsNames {

		if _, err := time.Parse("2006-01-02T15:04:05.999-0700", fieldsValues[index]); err == nil {
			values := strings.Split(fieldsValues[index], "T")
			comparedDate := values[0]
			comparedTime := values[1][0:8]

			valuePair = fmt.Sprintf("%sjavascript:gs.dateGenerate('%s','%s')", fieldsNames[index], comparedDate, comparedTime)

		} else {
			valuePair = fmt.Sprintf("%s%s", fieldsNames[index], fieldsValues[index])
		}

		if index == 0 {
			modifiedQuery = fmt.Sprintf("%s%s", modifiedQuery, valuePair)
		} else {
			modifiedQuery = fmt.Sprintf("%s^%s", modifiedQuery, valuePair)
		}

	}

	return
}

func (connector *SvcNowConnector) GetOpenTicketsByGroupID(methodOfDiscovery string, orgCode string, cloudAccountID string) (tickets <-chan domain.Ticket, err error) {
	err = fmt.Errorf("not implemented")
	return tickets, err
}

func (connector *SvcNowConnector) getTicketsForPassiveRescan(MethodOfDiscovery string, OrgCode string, Statuses []string) (Issues <-chan domain.Ticket, err error) {
	//if Statuses != nil {
	//	if len(*Statuses) > 0 {
	//
	//		// JQL -> project = vrr and "Method of Discovery" = Nexpose and (status IN (Open, "In Progress", Reopened, Resolved-Exception) AND (due <= 15d or created <= -20d))
	//		q := this.queryStart().
	//			and().
	//			equals(this.Fields[MOD], MethodOfDiscovery). // Must filter on MOD in order to ensure no overlap
	//			and().
	//			contains(this.Fields[ORG], OrgCode).
	//			and().
	//			beginGroup().
	//			In(this.Fields[STATUS], this.createStatusSlice(Statuses)). // status filter
	//			and().
	//			beginGroup().
	//			lessOrEquals("due", "15d"). // Due in the next 15 days
	//			or().
	//			lessOrEquals("created", "-20d"). // Created more than 20 days ago
	//			endGroup().
	//			endGroup().
	//			orderByAscend("created")
	//
	//		Issues, err = this.mapResultsToIssues(q)
	//	} else {
	//		err = errors.New("zero length status slice passed to getTicketsByStatusDueDateAscending")
	//	}
	//} else {
	//	err = errors.New("nil status slice passed to getTicketsByStatusDueDateAscending")
	//}
	return Issues, err
}

func (connector *SvcNowConnector) getTicketsForExceptionRescan(MethodOfDiscovery string, OrgCode string, Statuses []string) (Issues <-chan domain.Ticket, err error) {
	//if Statuses != nil {
	//
	//	if len(*Statuses) > 0 {
	//
	//		var cerfs []domain.CERF
	//		if cerfs, err = domain.GetExceptionsDueNext30Days(this.appconfig.DB); err == nil {
	//
	//			if len(cerfs) > 0 {
	//
	//				q := this.queryStart().
	//					and().
	//					equals(this.Fields[MOD], MethodOfDiscovery). // Must filter on MOD in order to ensure no overlap
	//					and().
	//					contains(this.Fields[ORG], OrgCode).
	//					and().
	//					In(this.Fields[STATUS], this.createStatusSlice(Statuses)). // status filter
	//					and()
	//
	//				q.beginGroup()
	//
	//				var index int = 0
	//				for i, _ := range cerfs {
	//					q.equals(this.Fields[CERF], cerfs[i].GetCERF())
	//
	//					if index < (len(cerfs) - 1) {
	//						q.or()
	//					}
	//					index++
	//				}
	//
	//				q.endGroup()
	//
	//				q.orderByAscend("created")
	//
	//				fmt.Println(q.JQL)
	//
	//				Issues, err = this.mapResultsToIssues(q)
	//			}
	//		} else {
	//			// TODO:
	//		}
	//	} else {
	//		err = errors.New("zero length status slice passed to getTicketsByStatusDueDateAscending")
	//	}
	//} else {
	//	err = errors.New("nil status slice passed to getTicketsByStatusDueDateAscending")
	//}
	return Issues, err
}

// GetStatusMap not implemented
func (connector *SvcNowConnector) GetStatusMap(in string) string {
	return ""
}
