package jira

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/andygrunwald/go-jira"
	"github.com/pkg/errors"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/log"
)

// ***********************************************************
// JIRA TICKET QUERIES
// ***********************************************************

// Unassigned holds the string value of the assignee of a ticket that has yet to be assigned
const Unassigned = "Unassigned"

// GetTicket returns a domain object containing information of the JIRA ticket relating to the SourceKey
func (connector *ConnectorJira) GetTicket(sourceKey string) (ticket domain.Ticket, err error) {
	if len(sourceKey) > 0 {
		var issue *jira.Issue
		issue, _, err = connector.client.Issue.Get(sourceKey, nil)

		if err == nil {
			if issue != nil {
				// Map the issue to the ticket
				ticket = &Issue{Issue: issue, connector: connector}
			} else {
				err = errors.New("Issue returned from JIRA is null prior to mapping to ticket in GetTicket")
			}
		}

	} else {
		err = errors.New("Invalid Source Key in GetTicket")
	}

	return ticket, err
}

// GetAdditionalTicketsForVulnPerDevice gets the additional tickets for the Vulns that have been scanned on all devices
func (connector *ConnectorJira) GetAdditionalTicketsForVulnPerDevice(tickets []domain.Ticket) (relatedTickets <-chan domain.Ticket, err error) {
	if len(tickets) > 0 {
		var queries []Query
		if queries, err = connector.getDeviceVulnsQueries(tickets); err == nil {
			relatedTickets = connector.runQueriesForIssues(queries, tickets)
		}
	}

	return relatedTickets, err
}

// GetAdditionalTicketsForDecomDevices gets the additional tickets by looking for all the tickets for the decommed devices(except those with "Closed-Remediated".
func (connector *ConnectorJira) GetAdditionalTicketsForDecomDevices(tickets []domain.Ticket) (relatedTickets <-chan domain.Ticket, err error) {
	if len(tickets) > 0 {
		var queries []Query
		if queries, err = connector.getDeviceTicketsQueries(tickets); err == nil {
			relatedTickets = connector.runQueriesForIssues(queries, tickets)
		}
	}

	return relatedTickets, err
}

// ***********************************************************
// JIRA TICKET CREATION METHODS
// ***********************************************************

// CreateTicket returns the source id and source key of the newly created ticket for logging in the database
func (connector *ConnectorJira) CreateTicket(ticket domain.Ticket) (sourceID int, sourceKey string, err error) {

	if ticket != nil {

		var ji *Issue
		if ji, err = connector.mapDalTicketToJiraIssue(ticket); err == nil {

			if ji.Issue != nil {

				ji.Issue.Fields.Project.Key = connector.project

				// Clear out these fields as they can't be sent
				// during create for vrr because they don't exist
				// in the create screen
				ji.Issue.Fields.Reporter = nil
				ji.Issue.Fields.Updated = jira.Time(time.Time{})
				ji.Issue.Fields.Created = jira.Time(time.Time{})
				ji.Issue.Fields.Status = nil

				var req *http.Request
				if req, err = connector.client.NewRequest(http.MethodPost, jcreateissue, ji.Issue); err == nil {
					var wg = sync.WaitGroup{}
					wg.Add(1)

					var response *http.Response
					if response, err = connector.funnelClient.Do(req); err == nil {
						if response != nil {
							defer response.Body.Close()

							var body []byte
							if body, err = ioutil.ReadAll(response.Body); err == nil {
								if err = json.Unmarshal(body, ji.Issue); err == nil {
									if sourceID, err = strconv.Atoi(ji.Issue.ID); err == nil {
										sourceKey = ji.Issue.Key
									}
								}
							}
						}
					}
				}
			} else {
				err = errors.New("JIRA Issue failed to map from Ticket Create Ticket")
			}
		}
	} else {
		err = errors.New("JIRA Ticket was passed null to Create Ticket")
	}

	return sourceID, sourceKey, err
}

// UpdateTicket takes a ticket of an input. The equivalent JIRA ticket is pulled. The difference between the two tickets is calculated, and
// a payload is generated that transforms the existing JIRA issue into the issue that was passed as a parameter
func (connector *ConnectorJira) UpdateTicket(ticket domain.Ticket, comment string) (SourceID int, SourceKey string, err error) {
	var oldToNewFieldName = make(map[string]string)

	if ticket != nil {
		var updateBlock = &updateBlock{}
		var existingTicket domain.Ticket
		existingTicket, err = connector.GetTicket(ticket.Title())
		var needToAddComment = true

		if existingTicket != nil {

			needToAddComment, err = connector.updateTicketStatusIfNecessary(ticket, existingTicket, comment)

			if ticket.HostName() != nil {
				field := connector.GetFieldMap(backendHostname)
				if field != nil {
					updateBlock.Fields.Hostname = ticket.HostName()
					oldToNewFieldName["hostname"] = field.ID
				}
			}

			if ticket.LastChecked() != nil && !ticket.LastChecked().IsZero() {
				field := connector.GetFieldMap(backendLastChecked)
				if field != nil {
					updateBlock.Fields.LastChecked = ticket.LastChecked()
					oldToNewFieldName["lastchecked"] = field.ID
				}
			}

			if ticket.ResolutionDate() != nil && existingTicket.ResolutionDate() != nil {
				if !ticket.ResolutionDate().IsZero() {
					var changeResolutionDate = false

					if !existingTicket.ResolutionDate().IsZero() {
						var firstTimeDiff = time.Since(*ticket.ResolutionDate())
						var secondTimeDiff = time.Since(*existingTicket.ResolutionDate())
						var timeDiff = firstTimeDiff - secondTimeDiff
						if timeDiff < 0 {
							timeDiff *= -1
						}
						if timeDiff > time.Hour*24 {
							changeResolutionDate = true
						}
					} else {
						changeResolutionDate = true
					}

					if changeResolutionDate {
						field := connector.GetFieldMap(backendResolutionDate)
						if field != nil {
							updateBlock.Fields.ResolutionDate = ticket.ResolutionDate().UTC().Format("2006-01-02T15:04:05.000+0000")
							oldToNewFieldName["resolutiondate"] = field.ID
						}
					}
				}
			}

			if ticket.Solution() != nil {
				field := connector.GetFieldMap(backendSolution)
				if field != nil {
					updateBlock.Fields.Solution = ticket.Solution()
					oldToNewFieldName["solution"] = field.ID
				}
			}

			if ticket.AssignedTo() != nil {
				updateBlock.Fields.AssignedTo = &Assignee{}
				updateBlock.Fields.AssignedTo.Name = *ticket.AssignedTo()
			}

			if ticket.AssignmentGroup() != nil {
				field := connector.GetFieldMap(backendAssignmentGroup)
				if field != nil {
					updateBlock.Fields.AssignmentGroup = &Assignee{}
					updateBlock.Fields.AssignmentGroup.Name = *ticket.AssignmentGroup()
					oldToNewFieldName["assignmentgroup"] = field.ID
				}
			}

			if ticket.CVEReferences() != nil {
				field := connector.GetFieldMap(backendCVEReferences)
				if field != nil {
					updateBlock.Fields.CveReferences = ticket.CVEReferences()
					oldToNewFieldName["cve_references"] = field.ID
				}
			}

			if existingTicket.CERF() != ticket.CERF() {
				field := connector.GetFieldMap(backendCERF)
				if field != nil {
					updateBlock.Fields.CerfLink = ticket.CERF()
					oldToNewFieldName["cerflink"] = field.ID
				}
			}

			if ticket.CVSS() != nil {
				field := connector.GetFieldMap(backendCVSS)
				if field != nil {
					updateBlock.Fields.CVSS = &ValueField{}
					updateBlock.Fields.CVSS.Value = fmt.Sprintf("%.1f", *ticket.CVSS())
					oldToNewFieldName["cvss"] = field.ID
				}
			}

			if ticket.Description() != nil {
				updateBlock.Fields.Description = ticket.Description()
			}

			//if ticket.OrgCode() != nil {
			field := connector.GetFieldMap(backendOrg)

			if field != nil {
				if existingTicket.OrgCode() != nil {
					if *existingTicket.OrgCode() != *ticket.OrgCode() {
						updateBlock.Fields.OrgCode = existingTicket.OrgCode()
						oldToNewFieldName["org"] = field.ID
					}
				} else {
					updateBlock.Fields.OrgCode = ticket.OrgCode()
					oldToNewFieldName["org"] = field.ID
				}
			}
			//}

			if ticket.DueDate() != nil {
				updateBlock.Fields.DueDate = ticket.DueDate()
			}

			if ticket.AlertDate() != nil {
				fields := connector.GetFieldMap(backendScanDate)
				if fields != nil {
					updateBlock.Fields.AlertDate = ticket.AlertDate()
					oldToNewFieldName["alertdate"] = fields.ID
				}
			}

			if ticket.VulnerabilityID() != existingTicket.VulnerabilityID() {
				field := connector.GetFieldMap(backendVulnerabilityID)
				if field != nil {
					updateBlock.Fields.VulnerabilityID = ticket.VulnerabilityID()
					oldToNewFieldName["vulnerabilityid"] = field.ID
				}
			}

			if ticket.IPAddress() != nil {
				field := connector.GetFieldMap(backendIPAddress)
				if field != nil {
					updateBlock.Fields.IPAddress = ticket.IPAddress()
					oldToNewFieldName["ipaddress"] = field.ID
				}
			}

			if ticket.VulnerabilityTitle() != nil {
				field := connector.GetFieldMap(backendVulnerability)
				if field != nil {
					updateBlock.Fields.VulnerabilityTitle = ticket.VulnerabilityTitle()
					oldToNewFieldName["vulnerabilitytitle"] = field.ID
				}
			}

			// TODO doesn't work
			//if ticket.ServicePorts() != nil {
			//	field := connector.GetFieldMap(backendServicePort)
			//	if field != nil {
			//		updateBlock.Fields.ServicePorts = ticket.ServicePorts()
			//		oldToNewFieldName["serviceports"] = field.ID
			//	}
			//}

			//TODO doesn't work - not on the edit screen
			//if ticket.Labels() != nil {
			//	if existingTicket.Labels() != nil {
			//		if *existingTicket.Labels() != *ticket.Labels() {
			//			updateBlock.Fields.Labels = &[]string{*ticket.Labels()}
			//		}
			//	} else {
			//		updateBlock.Fields.Labels = &[]string{*ticket.Labels()}
			//	}
			//}

			// TODO doesn't work
			//if ticket.MacAddress() != nil {
			//	field := connector.GetFieldMap(backendMACAddress)
			//	if field != nil {
			//		updateBlock.Fields.MACAddress = ticket.MacAddress()
			//		oldToNewFieldName["macaddress"] = field.ID
			//	}
			//}

			if ticket.Summary() != nil {
				updateBlock.Fields.Summary = ticket.Summary()
			}

			if err == nil {
				err = connector.executeUpdateAgainstJIRA(oldToNewFieldName, updateBlock, ticket, needToAddComment, comment)
			}
		} else {
			err = errors.Errorf("Could not find a ticket with the title %s", ticket.Title())
		}
	} else {
		err = errors.New("JIRA ticket was passed null to Update ticket")
	}

	return SourceID, SourceKey, err
}

func (connector *ConnectorJira) updateTicketStatusIfNecessary(ticket domain.Ticket, existingTicket domain.Ticket, comment string) (needToAddComment bool, err error) {
	needToAddComment = true

	if ticket.Status() != nil {
		if existingTicket.Status() != nil {
			if strings.ToLower(*existingTicket.Status()) != strings.ToLower(*ticket.Status()) {

				var assignedTo string
				if ticket.AssignedTo() != nil {
					assignedTo = *ticket.AssignedTo()
				} else {
					assignedTo = "Unassigned"
				}

				err = connector.Transition(ticket, sord(ticket.Status()), comment, assignedTo)

				if err != nil {
					if strings.Index(err.Error(), "not available") >= 0 {
						err = connector.getTransitionListForTicket(ticket, existingTicket)
					}
				} else {
					// TODO the comment is not being added on some transitions
					needToAddComment = true
				}
			}
		} else {
			//i expect each ticket to have a status
			err = errors.New("no status provided on ticket " + ticket.Title())
		}
	}

	return needToAddComment, err
}

func (connector *ConnectorJira) executeUpdateAgainstJIRA(oldToNewFieldName map[string]string, block *updateBlock, ticket domain.Ticket, needToAddComment bool, comment string) (err error) {
	var customFieldUpdateBlockBytes []byte
	var updateBlockWithCustomFieldNames interface{}
	customFieldUpdateBlockBytes, err = replaceJSONKey(oldToNewFieldName, block.Fields)

	if err == nil {
		var endString = fmt.Sprintf("{\"fields\":%s}", string(customFieldUpdateBlockBytes))
		//fmt.Println(endString)

		err = json.Unmarshal([]byte(endString), &updateBlockWithCustomFieldNames)
		if err == nil {
			var replace = "{issueIdOrKey}"
			var endPoint = html.EscapeString(jticket)
			endPoint = strings.Replace(endPoint, replace, ticket.Title(), 1)

			var request *http.Request
			if request, err = connector.client.NewRequest(http.MethodPut, endPoint, updateBlockWithCustomFieldNames); err == nil {

				var response *http.Response
				if response, err = connector.funnelClient.Do(request); err == nil {
					if response != nil {
						_ = response.Body.Close()

						if needToAddComment {
							err = connector.addComment(ticket, comment)
						}
					}
				}
			}
		} else {
			err = fmt.Errorf("error while trying unmarshall ticket update object - %s", err.Error())
		}

	}
	return err
}

// this method returns an error containing the transitions available from the status of the start ticket to the status of the end ticket
// this is an error as it is only called when an invalid transition was attempted
func (connector *ConnectorJira) getTransitionListForTicket(endTicket domain.Ticket, startTicket domain.Ticket) (err error) {
	var trans []jira.Transition
	var innerErr error
	trans, _, innerErr = connector.client.Issue.GetTransitions(endTicket.Title())
	if innerErr == nil {
		var availableTransitions = fmt.Sprintf("available transitions from %s to %s: ",
			*startTicket.Status(), *endTicket.Status())
		for index, tran := range trans {
			availableTransitions += tran.Name
			if index < len(trans)-1 {
				availableTransitions += ", "
			}
		}

		err = fmt.Errorf("%s", availableTransitions)
	} else {
		err = fmt.Errorf("%s: could not grab transitions", innerErr.Error())
	}
	return err
}

func replaceJSONKey(oldToNewKey map[string]string, object interface{}) (ret []byte, err error) {
	if object != nil {
		var err error
		var marshaled []byte
		marshaled, err = json.Marshal(object)
		if err == nil {
			var copyObject interface{}
			err = json.Unmarshal(marshaled, &copyObject)
			if err == nil {
				keys := copyObject.(map[string]interface{})
				for oldKey := range oldToNewKey {
					var newKey = oldToNewKey[oldKey]
					keys[newKey] = keys[oldKey]
					delete(keys, oldKey)
				}
				return json.Marshal(keys)
			}
		}
	} else {
		err = fmt.Errorf("empty object passed to replaceJSONKey")
	}

	return nil, err
}

// addComment adds a comment to a JIRA ticket
func (connector *ConnectorJira) addComment(ticket domain.Ticket, commentString string) (err error) {

	if len(commentString) > 0 {
		if ticket != nil {

			var issue *jira.Issue
			if issue, _, err = connector.client.Issue.Get(ticket.Title(), nil); err == nil {
				if issue != nil {

					var comment *jira.Comment
					comment = &jira.Comment{
						Body: commentString,
					}

					_, _, err = connector.client.Issue.AddComment(issue.ID, comment)
				}
			}

		} else {
			err = fmt.Errorf("empty ticket passed to addComment")
		}
	} else {
		// no error needs to be thrown
	}

	return err
}

// The assignee must be included in a different area of the body of the transition API call. These consts
// are used to track where in the body the assignee must be included
const assigneeNotFound = -1
const assigneeInField = 0
const assigneeInTransition = 1

func (connector *ConnectorJira) findAssigneeLocationAndSeeIfResDateIsRequired(ticketID string, payload TransitionPayload) (assigneeLocation int, resolutionDateRequired bool, err error) {
	assigneeLocation = assigneeNotFound

	var req *http.Request

	fieldEndpoint := fmt.Sprintf("/rest/api/2/issue/%s/transitions?expand=transitions.fields", ticketID)
	req, err = connector.client.NewRequest(http.MethodGet, fieldEndpoint, nil)
	if err == nil {

		var response *http.Response
		if response, err = connector.funnelClient.Do(req); err == nil {
			if response != nil {
				defer response.Body.Close()

				var body []byte
				if body, err = ioutil.ReadAll(response.Body); err == nil {

					obj := &transitionResult{}
					if err = json.Unmarshal(body, obj); err == nil {
						for index := range obj.Transitions {
							if obj.Transitions[index].ID == payload.ID {
								if len(obj.Transitions[index].Fields["assignee"].Name) > 0 {
									assigneeLocation = assigneeInField
								} else {
									assigneeLocation = assigneeInTransition
								}

								var resolutionDate = connector.GetFieldMap(backendResolutionDate).getCreateID()
								resolutionDateRequired = len(obj.Transitions[index].Fields[resolutionDate].Name) > 0
								// On some transitions required is returned false here, when it is actually required for a transition
								// maybe if the field is present in this map at all, it's required?
								break
							}
						}
					}
				}
			}
		}

	} else {
		err = fmt.Errorf("error while querying JIRA - %s", err.Error())
	}

	return assigneeLocation, resolutionDateRequired, err
}

// Transition changes the status of the ticket in JIRA (corresponding to the ticket parameter) to the parameter status
func (connector *ConnectorJira) Transition(ticket domain.Ticket, toStatus string, comment string, assignTo string) (err error) {
	if ticket != nil {

		var issue *jira.Issue
		if issue, _, err = connector.client.Issue.Get(ticket.Title(), nil); err == nil {
			if issue != nil {

				if strings.ToLower(issue.Fields.Status.Name) != strings.ToLower(toStatus) {
					var transitionSeries []workflowTransition
					transitionSeries, err = connector.getTransitionSeries(issue.Fields.Status.Name, toStatus)

					if err == nil {

						for transitionIndex := range transitionSeries {

							var trans []jira.Transition
							trans, _, err = connector.client.Issue.GetTransitions(ticket.Title())

							if err == nil {

								var apiHasTransition = false
								var transition = transitionSeries[transitionIndex]

								for checkIndex := range trans {
									if trans[checkIndex].ID == transition.id {
										apiHasTransition = true
										break
									}
								}

								if apiHasTransition {
									err = executeTransition(transition, assignTo, connector, ticket, comment)
								} else {
									err = errors.Errorf("Transition %s is not available for ticket %s. It is likely your JIRA account lacks a specific permission or your workflow is out of date", transition.name, ticket.Title())
								}

							} else {
								err = fmt.Errorf("error while gathering available transitions for ticket [%s] - %s", ticket.Title(), err.Error())
							}

							if err != nil {
								break
							}
						}

						if err == nil {
							err = connector.addComment(ticket, comment)
						}

					}
				}
			} else {
				err = errors.Errorf("Unable to find ticket %s", ticket.Title())
			}
		}

	} else {
		err = fmt.Errorf("ticket passed nil to Transition method")
	}

	return err
}

// AssignmentGroupExists checks to see if the input string is a valid assignment group in JIRA
func (connector *ConnectorJira) AssignmentGroupExists(groupName string) (exists bool, err error) {

	if connector.client != nil {

		var request *http.Request

		if request, err = connector.client.NewRequest(http.MethodGet, fmt.Sprintf("/rest/api/2/groups/picker?query=%s", groupName), nil); err == nil {

			var response *http.Response
			if response, err = connector.funnelClient.Do(request); err == nil {
				if response != nil {
					defer response.Body.Close()

					var body []byte
					if body, err = ioutil.ReadAll(response.Body); err == nil {

						var group *AssignmentGroupResponse
						var groupContainer = &assignmentGroupResponseWrapper{}
						if err = json.Unmarshal(body, groupContainer); err == nil {
							for index := range groupContainer.Groups {
								if groupContainer.Groups[index].Name == groupName {
									group = &groupContainer.Groups[index]
									break
								}
							}

							exists = group != nil
						}
					}
				}
			}
		}
	}

	return exists, err
}

// ***********************************************************
// JIRA SEARCH QUERIES
// ***********************************************************

// GetTicketsUpdatedSince get's all tickets that have been updated after the since date passed in
func (connector *ConnectorJira) GetTicketsUpdatedSince(since time.Time, orgCode string, methodOfDiscovery string) <-chan domain.Ticket {
	return connector.getTicketsUpdatedSince(since, orgCode, methodOfDiscovery)
}

// GetTicketsByClosedStatus retrieves closed exception and closed false positive tickets
func (connector *ConnectorJira) GetTicketsByClosedStatus(orgCode string, methodOfDiscovery string, startDate time.Time) (tix <-chan domain.Ticket) {
	return connector.getTicketsByClosedStatus(orgCode, methodOfDiscovery, startDate)
}

// GetCERFExpirationUpdates returns a map relating CERF tickets to their expiration date. It only grabs tickets that expire after the startDate parameter
func (connector *ConnectorJira) GetCERFExpirationUpdates(startDate time.Time) (cerfs map[string]time.Time, err error) {
	cerfs = make(map[string]time.Time)

	var issues <-chan domain.Ticket
	if issues, err = connector.getCERFExpirationUpdates(startDate); err == nil {

		for {

			if issue, ok := <-issues; ok {
				if len(issue.Title()) > 0 {
					cerfs[issue.Title()] = issue.CERFExpirationDate()
				}
			} else {
				break
			}
		}

	}

	return cerfs, err
}

// GetTicketsForRescan returns tickets for the rescan job. The type of rescan job is defined in Algorithm, and controls the tickets that are returned
func (connector *ConnectorJira) GetTicketsForRescan(cerfs []domain.CERF, MethodOfDiscovery string, OrgCode string, Algorithm string) (tickets <-chan domain.Ticket, err error) {
	tickets, err = connector.getTicketsForRescan(cerfs, MethodOfDiscovery, OrgCode, Algorithm)
	return tickets, err
}

// GetTicketsByDeviceIDVulnID returns tickets with the device and vulnerability id provided in the parameters
func (connector *ConnectorJira) GetTicketsByDeviceIDVulnID(methodOfDiscovery string, orgCode string, deviceID string, vulnID string, statuses map[string]bool, port int, protocol string) (tickets <-chan domain.Ticket, err error) {
	tickets, err = connector.getTicketsByDeviceIDVulnID(methodOfDiscovery, orgCode, deviceID, vulnID, statuses, port, protocol)
	return tickets, err
}

// GetOpenTicketsByGroupID returns tickets with an open status for an organization/method of discovery within a specified group
// For CIS tickets the entity ID is stored in the deviceID field, the ruleHash is stored in the vulnerabilityID field
// and the cloudAccountID is stored in the group ID
func (connector *ConnectorJira) GetOpenTicketsByGroupID(methodOfDiscovery string, orgCode string, groupID string) (tickets <-chan domain.Ticket, err error) {
	statuses := make(map[string]bool)
	statuses[connector.GetStatusMap(StatusOpen)] = true
	statuses[connector.GetStatusMap(StatusReopened)] = true
	statuses[connector.GetStatusMap(StatusResolvedRemediated)] = true
	statuses[connector.GetStatusMap(StatusResolvedDecom)] = true
	statuses[connector.GetStatusMap(StatusResolvedException)] = true
	statuses[connector.GetStatusMap(StatusResolvedFalsePositive)] = true

	tickets, err = connector.getOpenTicketsByGroupID(statuses, methodOfDiscovery, orgCode, groupID)
	return tickets, err
}

// GetStatusMap returns the mapped status to the status provided in the parameter. This is required as a JIRA project will have their own custom statuses
// that do not match the statuses defined in the application
func (connector *ConnectorJira) GetStatusMap(in string) string {
	var retVal = connector.statusMap[strings.ToLower(in)]
	if len(retVal) == 0 {
		message := fmt.Sprintf("[NO MAPPING FOUND FOR %s IN JIRA PAYLOAD]", in)
		connector.lstream.Send(log.Error("JIRA payload missing key elements", errors.New(message)))
		retVal = message
	}

	return retVal
}

// GetByCustomJQLChan returns the tickets that JIRA returns for a JQL statement onto a channel
func (connector *ConnectorJira) GetByCustomJQLChan(JQL string) (ticketChan <-chan domain.Ticket) {
	return connector.getByCustomJQLChan(JQL)
}
