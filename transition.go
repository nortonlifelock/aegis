package jira

import (
	"fmt"
	"github.com/trivago/tgo/tcontainer"
	"github.com/nortonlifelock/domain"
	"strings"
	"time"
)

func (connector *ConnectorJira) getTransitionSeries(fromStatus, toStatus string) (transitionSeries []workflowTransition, err error) {
	fromStatus = strings.ToLower(fromStatus)
	toStatus = strings.ToLower(toStatus)

	//this method finds the series of transitions required to traverse the ticket fromStatus->toStatus
	if connector.TransitionMap[fromStatus] != nil {
		if len(connector.TransitionMap[fromStatus][toStatus]) > 0 {
			//if we've already calculated the series, just return the series of transitions
			transitionSeries = connector.TransitionMap[fromStatus][toStatus]
		} else {
			//otherwise calculate it and store it in memory
			transitionSeries, err = connector.calculateTransitionSeries(fromStatus, toStatus)
			if err == nil {
				connector.TransitionMap[fromStatus][toStatus] = transitionSeries
			}
		}
	} else {
		err = fmt.Errorf("could not find any transitions for the status %s - please check your spelling as it must be exact", fromStatus)
	}

	return transitionSeries, err
}

type bfsMetaInfo struct {
	CurrentStatus    string
	TransitionSeries []workflowTransition
}

//the list of transitions from fromStatus->toStatus is found VIA a breadth first search
func (connector *ConnectorJira) calculateTransitionSeries(fromStatus, toStatus string) (transitionSeries []workflowTransition, err error) {
	//openSet contains child nodes waiting to be visited
	openSet := []*bfsMetaInfo{{fromStatus, []workflowTransition{}}}
	//closedSet contains child nodes that have already been visited
	closedSet := []string{fromStatus}

	transitionSeries = make([]workflowTransition, 0)

	for len(openSet) > 0 {
		//if there are child nodes remaining in the queue, pull the front one
		child := openSet[0]
		//remove the child from the queue
		openSet = openSet[1:]

		if child.CurrentStatus == toStatus {
			//success! return the series of transitions
			transitionSeries = child.TransitionSeries
			break
		} else {
			//if not succeeding, add each child to the queue
			for statusCreatedByTransition, transition := range connector.TransitionMap[child.CurrentStatus] {
				//don't add the child to the queue if it's in the closed set
				if stringNotInCloseSet(statusCreatedByTransition, closedSet) {
					//only add the child to the open set if it's not already in there
					if stringNotInOpenSet(statusCreatedByTransition, openSet) {
						//add the child to the open set, include the most recent transition to the transition list
						openSet = append(openSet, &bfsMetaInfo{statusCreatedByTransition, append(child.TransitionSeries, transition...)})
					}
					//add the child to the closed set
					closedSet = append(closedSet, statusCreatedByTransition)
				}
			}
		}
	}

	//if algorithm finds it's way here, no successful path was found
	if len(transitionSeries) == 0 {
		err = fmt.Errorf("could not find transition from %s to %s", fromStatus, toStatus)
	}

	return transitionSeries, err
}

func executeTransition(transition workflowTransition, assignTo string, connector *ConnectorJira, ticket domain.Ticket, Comment string) (err error) {
	payload := TransitionPayload{
		ID:       transition.id,
		Assignee: &assignTo,
		Unknowns: make(tcontainer.MarshalMap),
	}

	if assignTo == Unassigned {
		payload.Assignee = nil
	}

	var assigneeLocation int
	var resolutionDateRequired bool
	assigneeLocation, resolutionDateRequired, err = connector.findAssigneeLocationAndSeeIfResDateIsRequired(ticket.Title(), payload)
	if err == nil {

		fields := handleAssigneeLocationAndResolutionDateInFields(resolutionDateRequired, ticket, &payload, connector, assigneeLocation)

		tpayload := createTransitionPayload{
			Transition:  payload,
			Fields:      fields,
			UpdateBlock: Update{Comment: []UpdateObjects{{AddBody{Comment}}}},
		}

		if strings.Index(strings.ToLower(transition.name), "reopen") >= 0 {
			if tpayload.Fields == nil {
				tpayload.Fields = &FieldStruct{}
			}
			tpayload.Fields.ReopenReason = Comment
		}

		resDate, found := payload.Unknowns.Value(connector.GetFieldMap(backendResolutionDate).getCreateID())
		if found && resolutionDateRequired {
			if timeVal, ok := resDate.(time.Time); ok {

				if tpayload.Fields == nil {
					tpayload.Fields = &FieldStruct{}
				}

				if !timeVal.IsZero() {
					tpayload.Fields.ResolutionDate = timeVal.UTC().Format("2006-01-02T15:04:05.000+0000")
				}
			}
		}

		_, err = connector.client.Issue.DoTransitionWithPayload(ticket.Title(), tpayload)
	} else {
		err = fmt.Errorf("error while finding the assignee location - %s", err.Error())
	}

	return err
}

// the assignee can be in two locations of the request. this method removes the assignee from the portion of the request that it is not required
// if the resolution date is required for this particular transition, it is set by this method
func handleAssigneeLocationAndResolutionDateInFields(resolutionDateRequired bool, Ticket domain.Ticket, payload *TransitionPayload, connector *ConnectorJira, assigneeLocation int) *FieldStruct {
	if resolutionDateRequired {
		if !tord(Ticket.ResolutionDate()).IsZero() {
			payload.Unknowns[connector.GetFieldMap(backendResolutionDate).getCreateID()] = tord(Ticket.ResolutionDate())
		}
	}

	var fields = &FieldStruct{}
	if payload.Assignee != nil {
		fields.Assignee = &Assignee{}
		fields.Assignee.Name = *payload.Assignee
	}

	//one of these will likely need to be nil
	if assigneeLocation == assigneeInField {
		payload.Assignee = nil
	} else if assigneeLocation == assigneeInTransition {
		fields = nil
	} else {
		// this occurs if the transition attempted isn't available, want the error JIRA returns in the API call that will be made at the end of this method
	}

	return fields
}

func buildTransitionIDToNameMap(workflow *workflow) (transitionIDToTransitionName map[string]string) {
	transitionIDToTransitionName = make(map[string]string)

	for _, status := range workflow.Statuses {
		transitionIDToTransitionName[status.StatusID] = status.StatusName
	}

	return transitionIDToTransitionName
}

func buildCommonActionMap(workflow *workflow) (actionIDToAction map[string]Action) {
	actionIDToAction = make(map[string]Action)

	for _, commonAction := range workflow.CommonActions {
		actionIDToAction[commonAction.ActionID] = commonAction
	}

	return actionIDToAction
}

func stringNotInOpenSet(input string, set []*bfsMetaInfo) bool {
	var haveNotSeenElement = true
	for _, element := range set {
		if input == element.CurrentStatus {
			haveNotSeenElement = false
			break
		}
	}
	return haveNotSeenElement
}

func stringNotInCloseSet(input string, set []string) bool {
	var haveNotSeenElement = true
	for _, element := range set {
		if input == element {
			haveNotSeenElement = false
			break
		}
	}
	return haveNotSeenElement
}
