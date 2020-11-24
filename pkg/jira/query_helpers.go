package jira

import (
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nortonlifelock/aegis/pkg/domain"
	"github.com/nortonlifelock/aegis/pkg/log"
	"github.com/pkg/errors"
)

func (connector *ConnectorJira) createStatusSlice(statuses map[string]bool) (statusStrings []string) {
	statusStrings = make([]string, 0)

	for status := range statuses {
		if len(status) > 0 {
			statusStrings = append(statusStrings, status)
		}
	}

	return statusStrings
}

func (connector *ConnectorJira) getTicketsUpdatedSince(since time.Time, orgCode string, methodOfDiscovery string) (issues <-chan domain.Ticket) {
	if !since.IsZero() {
		q := connector.queryStart().and().
			greaterThan(connector.GetFieldMap(backendUpdated), fmt.Sprintf("\"%s\"", since.Format(QueryDateTimeFormatJira))).
			and().
			contains(connector.GetFieldMap(backendOrg), orgCode).
			and().
			equals(connector.GetFieldMap(backendMOD), methodOfDiscovery).
			orderByAscend(connector.GetFieldMap(backendUpdated).Name)

		var errChan <-chan error
		issues, errChan = connector.getSearchResults(q)
		emptyErrChan(errChan) // the errors are logged but do not need to be processed by the caller
	} else {
		connector.lstream.Send(log.Error("zero date passed to getTicketsUpdatedSince", nil))
	}

	return issues
}

func (connector *ConnectorJira) getTicketsByClosedStatus(orgCode string, methodOfDiscovery string, startDate time.Time) (issues <-chan domain.Ticket) {
	q := connector.queryStart().
		and().
		contains(connector.GetFieldMap(backendOrg), orgCode).
		and().
		equals(connector.GetFieldMap(backendMOD), methodOfDiscovery).
		and().
		beginGroup().
		beginGroup().
		equals(connector.GetFieldMap(backendStatus), connector.GetStatusMap(domain.StatusClosedException)).
		and().
		notEmpty(connector.GetFieldMap(backendCERF)).
		endGroup().
		or().
		equals(connector.GetFieldMap(backendStatus), connector.GetStatusMap(domain.StatusClosedFalsePositive)).
		endGroup().
		and().
		greaterOrEquals(connector.GetFieldMap(backendUpdated), fmt.Sprintf("\"%s\"", startDate.Format(QueryDateTimeFormatJira))).
		orderByAscend(connector.GetFieldMap(backendUpdated).Name)

	var errChan <-chan error
	issues, errChan = connector.getSearchResults(q)
	emptyErrChan(errChan) // the errors are logged but do not need to be processed by the caller
	return issues
}

func (connector *ConnectorJira) getOpenTicketsByGroupID(statuses map[string]bool, methodOfDiscovery string, orgCode string, groupID string) (issues <-chan domain.Ticket, errChan <-chan error) {
	if len(statuses) > 0 {
		if len(methodOfDiscovery) > 0 {
			if len(orgCode) > 0 {
				if len(groupID) > 0 {

					q := connector.queryStart().
						and().
						equals(connector.GetFieldMap(backendMOD), methodOfDiscovery).
						and().
						contains(connector.GetFieldMap(backendOrg), orgCode).
						and().
						contains(connector.GetFieldMap(backendGroupID), groupID)

					q.and().beginGroup()

					var index = 0
					for status := range statuses {
						q.equals(connector.GetFieldMap(backendStatus), status)

						if index < (len(statuses) - 1) {
							q.or()
						}

						index++
					}

					q.endGroup()

					issues, errChan = connector.getSearchResults(q)
				} else {
					out := make(chan error, 1)
					out <- fmt.Errorf("zero length groupID sent to getOpenTicketsByGroupID")
					close(out)
					errChan = out
				}

			} else {
				out := make(chan error, 1)
				out <- fmt.Errorf("zero length orgCode sent to getOpenTicketsByGroupID")
				close(out)
				errChan = out
			}
		} else {
			out := make(chan error, 1)
			out <- fmt.Errorf("zero length methodOfDiscovery sent to getOpenTicketsByGroupID")
			close(out)
			errChan = out
		}
	} else {
		out := make(chan error, 1)
		out <- fmt.Errorf("zero length statuses sent to getOpenTicketsByGroupID")
		close(out)
		errChan = out
	}

	return issues, errChan
}

func (connector *ConnectorJira) getTicketsByDeviceIDVulnID(methodOfDiscovery string, orgCode string, deviceID string, vulnID string, statuses map[string]bool, port int, protocol string) (issues <-chan domain.Ticket, err error) {

	if len(deviceID) > 0 {
		if len(vulnID) > 0 {
			if statuses != nil {
				if len(statuses) > 0 {

					q := connector.queryStart().
						and().
						equals(connector.GetFieldMap(backendMOD), methodOfDiscovery). // Must filter on MOD in order to ensure no overlap
						and().
						contains(connector.GetFieldMap(backendOrg), orgCode).
						and().
						contains(connector.GetFieldMap(backendDeviceID), fmt.Sprintf("%v", deviceID)).
						and().
						contains(connector.GetFieldMap(backendVulnerabilityID), vulnID)

					if port >= 0 && port <= 65535 && len(protocol) > 0 {

						var portText string
						portText = strconv.Itoa(port)

						servicePort := fmt.Sprintf("%s %s", portText, protocol)
						q.and().contains(connector.GetFieldMap(backendServicePort), servicePort)

					}

					q.and().beginGroup()

					var index = 0
					for status := range statuses {
						q.equals(connector.GetFieldMap(backendStatus), status)

						if index < (len(statuses) - 1) {
							q.or()
						}

						index++
					}

					q.endGroup()

					var errChan <-chan error
					issues, errChan = connector.getSearchResults(q)
					emptyErrChan(errChan) // the errors are logged but do not need to be processed by the caller
				} else {
					err = errors.New("Must pass more than one status")
				}
			} else {
				err = errors.New("NIL Status object passed")
			}
		} else {
			err = errors.Errorf("Vulnerability id passed as [%s] to getTicketsByDeviceIDVulnID", vulnID)
		}
	} else {
		err = errors.New("Device id passed as [%d] to getTicketsByDeviceIDVulnID")
	}

	return issues, err
}

func (connector *ConnectorJira) getTicketsForRescan(cerfs []domain.CERF, groupID string, methodOfDiscovery string, orgCode string, algorithm string) (issues <-chan domain.Ticket, errChan <-chan error) {

	if len(methodOfDiscovery) > 0 {
		if len(orgCode) > 0 {
			switch algorithm {
			case domain.RescanExceptions:
				var status = make(map[string]bool)
				status[connector.GetStatusMap(domain.StatusClosedException)] = true

				issues, errChan = connector.getTicketsForExceptionRescan(cerfs, methodOfDiscovery, orgCode, status)
				break
			case domain.RescanPassive:
				var status = make(map[string]bool)
				status[connector.GetStatusMap(domain.StatusOpen)] = true
				status[connector.GetStatusMap(domain.StatusInProgress)] = true
				status[connector.GetStatusMap(domain.StatusReopened)] = true
				status[connector.GetStatusMap(domain.StatusResolvedException)] = true

				issues, errChan = connector.getTicketsForPassiveRescan(methodOfDiscovery, orgCode, status)
				break
			case domain.RescanNormal:
				var status = make(map[string]bool)
				status[connector.GetStatusMap(domain.StatusResolvedRemediated)] = true

				issues, errChan = connector.getTicketsByStatusDueDateAscending(groupID, methodOfDiscovery, orgCode, status)
				break
			case domain.RescanDecommission:
				var status = make(map[string]bool)
				status[connector.GetStatusMap(domain.StatusResolvedDecom)] = true

				issues, errChan = connector.getTicketsByStatusDueDateAscending(groupID, methodOfDiscovery, orgCode, status)
				break
			default:
				out := make(chan error, 1)
				out <- fmt.Errorf("error: unknown algorithm [%s]", algorithm)
				close(out)
				errChan = out
				break
			}
		} else {
			out := make(chan error, 1)
			out <- fmt.Errorf("error: empty organization code passed to getTicketsForRescan")
			close(out)
			errChan = out
		}
	} else {
		out := make(chan error, 1)
		out <- fmt.Errorf("error: empty method of discovery passed to getTicketsForRescan")
		close(out)
		errChan = out
	}

	return issues, errChan
}

func (connector *ConnectorJira) getTicketsByStatusDueDateAscending(groupID string, methodOfDiscovery string, orgCode string, statuses map[string]bool) (issues <-chan domain.Ticket, errChan <-chan error) {
	if statuses != nil {

		if len(statuses) > 0 {

			q := connector.queryStart().
				and().
				equals(connector.GetFieldMap(backendMOD), methodOfDiscovery). // Must filter on MOD in order to ensure no overlap
				and().
				contains(connector.GetFieldMap(backendOrg), orgCode).
				and()

			q.beginGroup()

			var index = 0
			for status := range statuses {
				q.equals(connector.GetFieldMap(backendStatus), status)

				if index < (len(statuses) - 1) {
					q.or()
				}

				index++
			}

			if len(groupID) > 0 {
				q.and().contains(connector.GetFieldMap(backendGroupID), groupID)
			}

			q.endGroup().
				orderByAscend("due")

			issues, errChan = connector.getSearchResults(q)
		} else {
			out := make(chan error, 1)
			out <- errors.New("zero length status slice passed to getTicketsByStatusDueDateAscending")
			close(out)
			errChan = out
		}
	} else {
		out := make(chan error, 1)
		out <- errors.New("nil status slice passed to getTicketsByStatusDueDateAscending")
		close(out)
		errChan = out
	}
	return issues, errChan
}

func (connector *ConnectorJira) getTicketsForPassiveRescan(methodOfDiscovery string, orgCode string, statuses map[string]bool) (issues <-chan domain.Ticket, errChan <-chan error) {
	if statuses != nil {
		if len(statuses) > 0 {

			// JQL -> project = Aegis and "Method of Discovery" = Nexpose and (status IN (Open, "In Progress", Reopened, Resolved-Exception) AND (due <= 15d or created <= -20d))
			q := connector.queryStart().
				and().
				equals(connector.GetFieldMap(backendMOD), methodOfDiscovery). // Must filter on MOD in order to ensure no overlap
				and().
				contains(connector.GetFieldMap(backendOrg), orgCode).
				and().
				beginGroup().
				in(connector.GetFieldMap(backendStatus), connector.createStatusSlice(statuses)). // status filter
				and().
				beginGroup().
				lessOrEquals("due", "15d"). // Due in the next 15 days
				or().
				lessOrEquals("created", "-20d"). // Created more than 20 days ago
				endGroup().
				endGroup().
				orderByAscend("created")

			issues, errChan = connector.getSearchResults(q)
		} else {
			out := make(chan error, 1)
			out <- errors.New("zero length status slice passed to getTicketsByStatusDueDateAscending")
			close(out)
			errChan = out
		}
	} else {
		out := make(chan error, 1)
		out <- errors.New("nil status slice passed to getTicketsByStatusDueDateAscending")
		close(out)
		errChan = out
	}
	return issues, errChan
}

func (connector *ConnectorJira) getTicketsForExceptionRescan(cerfs []domain.CERF, methodOfDiscovery string, orgCode string, statuses map[string]bool) (issues <-chan domain.Ticket, errChan <-chan error) {
	if statuses != nil {

		if len(statuses) > 0 {
			if len(cerfs) > 0 {

				q := connector.queryStart().
					and().
					equals(connector.GetFieldMap(backendMOD), methodOfDiscovery). // Must filter on MOD in order to ensure no overlap
					and().
					contains(connector.GetFieldMap(backendOrg), orgCode).
					and().
					in(connector.GetFieldMap(backendStatus), connector.createStatusSlice(statuses)). // status filter
					and()

				q.beginGroup()

				var index = 0
				for i := range cerfs {
					q.equals(connector.GetFieldMap(backendCERF), cerfs[i].CERForm())

					if index < (len(cerfs) - 1) {
						q.or()
					}
					index++
				}

				q.endGroup()

				q.orderByAscend("created")

				issues, errChan = connector.getSearchResults(q)
			}

		} else {
			out := make(chan error, 1)
			out <- errors.New("zero length status slice passed to getTicketsByStatusDueDateAscending")
			close(out)
			errChan = out
		}
	} else {
		out := make(chan error, 1)
		out <- errors.New("nil status slice passed to getTicketsByStatusDueDateAscending")
		close(out)
		errChan = out
	}

	return issues, errChan
}

func (connector *ConnectorJira) getByCustomJQL(JQL string) (issues <-chan domain.Ticket, err error) {
	var q = &Query{
		JQL:  JQL,
		Size: 1000,
	}
	var errChan <-chan error
	issues, errChan = connector.getSearchResults(q)
	emptyErrChan(errChan) // the errors are logged but do not need to be processed by the caller
	return issues, err
}

func (connector *ConnectorJira) getByCustomJQLChan(JQL string) (issues <-chan domain.Ticket) {
	var q = &Query{
		JQL:  JQL,
		Size: 1000,
	}
	var errChan <-chan error
	issues, errChan = connector.getSearchResults(q)
	emptyErrChan(errChan) // the errors are logged but do not need to be processed by the caller
	return issues
}

func (connector *ConnectorJira) DeleteProject(project string) (err error) {
	var replace = "{projectIdOrKey}"

	var endPoint = html.EscapeString(jproject)
	endPoint = strings.Replace(endPoint, replace, project, 1)

	var request *http.Request

	if request, err = connector.client.NewRequest(http.MethodDelete, endPoint, nil); err == nil {

		var response *http.Response
		if response, err = connector.funnelClient.Do(request); err == nil {
			if response != nil {
				defer response.Body.Close()
			}
		}
	}

	return err
}

// DeleteTicket delete the ticket in JIRA corresponding to the id parameter. Used by the JIRA tool
func (connector *ConnectorJira) DeleteTicket(id string) (err error) {
	var replace = "{issueIdOrKey}"

	var endPoint = html.EscapeString(jticket)
	endPoint = strings.Replace(endPoint, replace, id, 1)

	var request *http.Request

	if request, err = connector.client.NewRequest(http.MethodDelete, endPoint, nil); err == nil {

		var response *http.Response
		if response, err = connector.funnelClient.Do(request); err == nil {
			if response != nil {
				defer response.Body.Close()
			}
		}
	}

	return err
}

// GetFieldsForProject only returns the field allowed by the mappable_fields dictated in the jira payload. It is used by the API
// TODO should we just return the list of mappable fields in the connector payload? do we need to confirm with the API that the fields exist for the project?
func (connector *ConnectorJira) GetFieldsForProject(project string, includeNonMappables bool) (fields []string, err error) {
	var toReplace = "{project}"
	var endpoint = strings.Replace(jeditableprojectfields, toReplace, project, 1)

	var request *http.Request

	fields = make([]string, 0)

	if request, err = connector.client.NewRequest(http.MethodGet, endpoint, nil); err == nil {

		var response *http.Response
		if response, err = connector.funnelClient.Do(request); err == nil {
			if response != nil {
				defer response.Body.Close()

				var body []byte
				if body, err = ioutil.ReadAll(response.Body); err == nil {
					var splitByName = strings.Split(string(body), "\"name\":\"")
					if len(splitByName) > 0 {
						splitByName = splitByName[1:]

						var seenField = make(map[string]bool)
						for _, lineStartingWithName := range splitByName {
							indexOfEndOfFieldName := strings.Index(lineStartingWithName, "\"")

							var newField = lineStartingWithName[:indexOfEndOfFieldName]

							if !seenField[newField] {
								if connector.isMappableField(newField) || includeNonMappables {
									fields = append(fields, newField)
								}

								seenField[newField] = true
							}

						}
					}
				}
			}
		}
	}

	return fields, err
}

// getCountByJQL returns the amount of tickets returned by JIRA from a JQL
func (connector *ConnectorJira) getCountByJQL(q *Query) (count int, err error) {
	var ret interface{}
	ret, err = connector.getSearchResponse(q, newSearchResult(), 0)
	if err == nil {
		if searchResult, ok := ret.(*searchResult); ok {
			count = searchResult.Total
		} else {
			err = fmt.Errorf("return did not appear to be a search result")
		}
	}

	return count, err
}

// GetCountOfTicketsInStatus returns the amount of tickets in a given status. Used by the API to populate dashboard information
func (connector *ConnectorJira) GetCountOfTicketsInStatus(status string, orgCode string) (count int, err error) {
	if len(status) > 0 {
		var q = connector.queryStart().and().equals(connector.GetFieldMap(backendStatus), status)

		if len(orgCode) > 0 {
			q.and().contains(connector.GetFieldMap(backendOrg), orgCode)
		}

		count, err = connector.getCountByJQL(q)
	} else {
		err = fmt.Errorf("cannot pass empty status to GetCountOfTicketsInStatus")
	}

	return count, err
}

// getDeviceTicketsQueries retrieves all the tickets for the decommed devices
func (connector *ConnectorJira) getDeviceTicketsQueries(issues []domain.Ticket) (qs []Query, err error) {
	if issues != nil {
		//We getting all the statues except for this status
		var statuses = []string{
			connector.GetStatusMap(domain.StatusOpen), connector.GetStatusMap(domain.StatusReopened), connector.GetStatusMap(domain.StatusInProgress),
			connector.GetStatusMap(domain.StatusResolvedRemediated), connector.GetStatusMap(domain.StatusResolvedFalsePositive), connector.GetStatusMap(domain.StatusResolvedDecom), connector.GetStatusMap(domain.StatusResolvedException),
			connector.GetStatusMap(domain.StatusClosedException), connector.GetStatusMap(domain.StatusClosedFalsePositive),
		}

		//This map is created to not process the same device_ID again if exist in another issue
		var devices = make(map[string]bool)
		for index := range issues {
			deviceIDInSearch := issues[index].DeviceID()
			key := issues[index].Title()
			if len(deviceIDInSearch) > 0 {
				//I check if the device_id had processed before, if so exit and process another device_id
				if !devices[deviceIDInSearch] {
					//mark the device_id as processed once you entered this clause
					devices[deviceIDInSearch] = true
					q := connector.queryStart().
						and().
						in(connector.GetFieldMap(backendStatus), statuses).
						and().
						contains(connector.GetFieldMap(backendDeviceID), deviceIDInSearch).
						and().
						notEquals(connector.GetFieldMap(backendKey), key).
						and().
						contains(connector.GetFieldMap(backendOrg), sord(issues[index].OrgCode()))

					qs = append(qs, *q)
				}
			}
		}

	}
	return qs, err
}

// getDeviceVulnsQueries gathers the vulnerabilities for the device and maps them to the devices
func (connector *ConnectorJira) getDeviceVulnsQueries(issues []domain.Ticket) (qs []Query, err error) {
	if issues != nil {
		var statuses = []string{connector.GetStatusMap(domain.StatusOpen), connector.GetStatusMap(domain.StatusReopened),
			connector.GetStatusMap(domain.StatusResolvedException), connector.GetStatusMap(domain.StatusResolvedDecom),
			connector.GetStatusMap(domain.StatusInProgress)}

		//Get all the vulns and map them to the devices.
		var vulnsPerDevices map[string]map[string]bool
		vulnsPerDevices, err = connector.getVulnsFromIssues(issues)
		if vulnsPerDevices != nil {
			//This map is created to not process the same device_ID again if exist in another issue
			var devices = make(map[string]bool)
			for index := range issues {
				deviceIDInSearch := issues[index].DeviceID()
				if len(deviceIDInSearch) > 0 {
					//I check if the device_id had processed before, if so exit and process another device_id
					if !devices[deviceIDInSearch] {
						//mark the device_id as processed once you entered this clause
						devices[deviceIDInSearch] = true
						q := connector.queryStart().
							and().
							in(connector.GetFieldMap(backendStatus), statuses).
							and().
							contains(connector.GetFieldMap(backendDeviceID), deviceIDInSearch).
							and().
							contains(connector.GetFieldMap(backendOrg), sord(issues[index].OrgCode())).
							and().beginGroup() //start my query to look for the device_id

						var vulnCount = 0
						// We loop through vulns map(we dont have duplicate vulns to start with). The value contains a map of devices which donates that
						// this combination(device_id and vulnId) had a ticket already so dont pull the same ticket again.
						for key, value := range vulnsPerDevices {

							// This check helps us ignore existing tickets in the original issues list. This identify that this vuln exist in a ticket with this deviceid
							if !value[deviceIDInSearch] {
								if vulnCount > 0 {
									q.or()
								}
								//Here we create the query of all vulns for the device_id
								q.contains(connector.GetFieldMap(backendVulnerabilityID), key)
								// We count the vulns that was included in the query
								vulnCount++
							}
						}
						q.endGroup()
						//If we have at least one vulnerability added to the query, then we will include the query
						if vulnCount > 0 {
							qs = append(qs, *q)
						}
					}
				}
			}
		}
	}
	return qs, err
}

func (connector *ConnectorJira) runQueriesForScheduledScan(groupID string, methodOfDiscovery string, orgCode string) (relatedTickets <-chan domain.Ticket, err error) {
	var out = make(chan domain.Ticket)
	go func() {
		defer close(out)
		var errChan <-chan error
		var normalTickets, decommTickets <-chan domain.Ticket
		normalTickets, errChan = connector.getTicketsForRescan(nil, groupID, methodOfDiscovery, orgCode, domain.RescanNormal)
		if err = getFirstErrorFromChannel(errChan); err == nil {

			decommTickets, errChan = connector.getTicketsForRescan(nil, groupID, methodOfDiscovery, orgCode, domain.RescanDecommission)
			if err = getFirstErrorFromChannel(errChan); err == nil {

				for {
					if tic, ok := <-normalTickets; ok {
						out <- tic
					} else {
						break
					}
				}

				for {
					if tic, ok := <-decommTickets; ok {
						out <- tic
					} else {
						break
					}
				}
			} else {
				connector.lstream.Send(log.Errorf(err, "error while grabbing normal tickets for scheduled scans"))
			}
		} else {
			connector.lstream.Send(log.Errorf(err, "error while grabbing decommission tickets for scheduled scans"))
		}
	}()

	return out, err
}

// runQueriesForIssues runs all the queries that has been generated to pull the additional tickets per device/vuln
func (connector *ConnectorJira) runQueriesForIssues(qs []Query, tickets []domain.Ticket) <-chan domain.Ticket {
	var relatedTickets = make(chan domain.Ticket)

	go func() {
		defer close(relatedTickets)

		var err error

		var issuesMap = sync.Map{}
		for index := range tickets {
			key := tickets[index].Title()
			issuesMap.LoadOrStore(key, true)
		}

		var expectedDeviceIDs = make(map[string]bool)
		for _, ticket := range tickets {
			expectedDeviceIDs[ticket.DeviceID()] = true
		}

		if qs != nil {
			//Issues, err = this.mapResultsToIssues(q)
			for index := range qs {
				var errChan <-chan error
				var vulnIssues <-chan domain.Ticket
				if vulnIssues, errChan = connector.getSearchResults(&qs[index]); vulnIssues != nil {

					// TODO we query based on the device ID, which must be accomplished using a contains
					// this means an incorrect device ID may pass the query if it contains the other device ID
					// e.g. contains 1120 will also query a device with ID 11200
					for {
						if issue, ok := <-vulnIssues; ok {
							// ensure the exact device ID was seen
							if expectedDeviceIDs[issue.DeviceID()] {
								relatedTickets <- issue
							}
						} else {
							break
						}
					}
				} else {
					connector.lstream.Send(log.Errorf(err, "error while mapping result to issue [%v]", qs[index]))
				}

				emptyErrChan(errChan) // the errors are logged but do not need to be processed by the caller
			}
		}
	}()

	return relatedTickets
}

// Gets all vuls mapped to deviceids that exist in the tickets
// It helps us identify the vulns that we dont need to pull for a specif device ID when the jql is created
func (connector *ConnectorJira) getVulnsFromIssues(issues []domain.Ticket) (vulnPerDevices map[string]map[string]bool, err error) {
	if issues != nil {
		vulnPerDevices = make(map[string]map[string]bool)
		for index := range issues {
			deviceVuln := issues[index].VulnerabilityID()
			if len(deviceVuln) > 0 {
				deviceID := issues[index].DeviceID()
				if len(deviceID) > 0 {
					if vulnPerDevices[deviceVuln] == nil {
						vulnPerDevices[deviceVuln] = make(map[string]bool)
					}

					vulnPerDevices[deviceVuln][deviceID] = true
				}
			}
		}
	}

	return vulnPerDevices, err
}

func getFirstErrorFromChannel(errChan <-chan error) (err error) {
	for {
		var ok bool
		if err, ok = <-errChan; !ok || err != nil {
			break
		}
	}

	return err
}
