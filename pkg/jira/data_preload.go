package jira

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"io/ioutil"
	"net/http"
	"strings"
)

// ***********************************************************
// This file holds the methods which gathers the information
// required for the JIRA connector
// ***********************************************************

// Get the complete list of Fields for the JIRA system and store them in a map
func (connector *ConnectorJira) getFields() (fields map[string]*Field, err error) {

	if connector.client != nil {

		// Init the map
		fields = make(map[string]*Field)

		var request *http.Request

		// Create a jira api request for the custom Fields in JIRA
		if request, err = connector.client.NewRequest(http.MethodGet, jfield, nil); err == nil {

			// Create a new slice of Custom Fields
			jFields := new([]*Field)

			var response *http.Response
			if response, err = connector.funnelClient.Do(request); err == nil {
				if response != nil {
					defer response.Body.Close()

					var body []byte
					if body, err = ioutil.ReadAll(response.Body); err == nil {
						if err = json.Unmarshal(body, jFields); err == nil {
							// Loop through the custom Fields and assign those Fields to the custom field map
							for _, field := range *jFields {
								if field != nil {
									// Map the Fields
									fields[field.Name] = field
								}
							}
						}
					}
				}
			}
		}
	}

	return fields, err
}

// Get the complete list of Resolutions for the JIRA system and store them in a map
func (connector *ConnectorJira) getResolutions() (resolutions map[string]*jira.Resolution, err error) {

	if connector.client != nil {

		// Init the map
		resolutions = make(map[string]*jira.Resolution)

		var request *http.Request

		// Create a jira api request for the resolutions in JIRA
		if request, err = connector.client.NewRequest(http.MethodGet, jresolution, nil); err == nil {

			// Create a new slice of resolution
			jResolutions := new([]*jira.Resolution)

			var response *http.Response
			if response, err = connector.funnelClient.Do(request); err == nil {
				if response != nil {
					defer response.Body.Close()

					var body []byte
					if body, err = ioutil.ReadAll(response.Body); err == nil {
						if err = json.Unmarshal(body, jResolutions); err == nil {
							// Loop through the resolution and assign those resolutions to the map
							for _, resolution := range *jResolutions {
								if resolution != nil {
									// Map the resolutions
									resolutions[resolution.Name] = resolution
								}
							}
						}
					}
				}
			}
		}
	}

	return resolutions, err
}

// TurnXMLWorkFlowToJSON takes in the JIRA workflow (describes the transitions between JIRA statuses) and turns it into JSON so it can be stored in the Payload of the Source Config for the JIRA connection
// How to find the JIRA workflow in your JIRA instance: go to JIRA Administration (cog) -> Projects -> Workflows (on left) -> Select "actions" button next to desired workflow -> Export as XML
func TurnXMLWorkFlowToMap(xmlWorkFlow string) (fromToTransition map[string]map[string][]workflowTransition, err error) {
	fromToTransition = make(map[string]map[string][]workflowTransition)
	workflow := &workflow{}

	err = xml.Unmarshal([]byte(xmlWorkFlow), workflow)
	var transitionIDToName = buildTransitionIDToNameMap(workflow)
	var commonActionMap = buildCommonActionMap(workflow)

	if err == nil {
		for _, status := range workflow.Statuses {
			var fromStatus = strings.ToLower(status.StatusName)

			if fromToTransition[fromStatus] == nil {
				fromToTransition[fromStatus] = make(map[string][]workflowTransition)
			}

			for _, action := range status.Actions {
				var idForStatusArrivedAtByTakingTransition = action.TransitionDetails.DestinationStatusID
				var toStatus = strings.ToLower(transitionIDToName[idForStatusArrivedAtByTakingTransition])
				var singleTrans = workflowTransition{
					ID:   action.ActionID,
					Name: action.TransitionName,
				}
				fromToTransition[fromStatus][toStatus] = []workflowTransition{singleTrans}
			}

			for _, commonActionID := range status.CommonActions {
				var action = commonActionMap[commonActionID.ID]
				var idForStatusArrivedAtByTakingTransition = action.TransitionDetails.DestinationStatusID
				var toStatus = strings.ToLower(transitionIDToName[idForStatusArrivedAtByTakingTransition])
				var singleTrans = workflowTransition{
					ID:   commonActionMap[commonActionID.ID].ActionID,
					Name: commonActionMap[commonActionID.ID].TransitionName,
				}
				fromToTransition[fromStatus][toStatus] = []workflowTransition{singleTrans}
			}
		}

	} else {
		err = fmt.Errorf("error while unmarshalling JIRA workflow - %s", err.Error())
	}

	return fromToTransition, err
}

// Get the complete list of Statuses for the JIRA system and store them in a map
func (connector *ConnectorJira) getStatuses() (statuses map[string]*jira.Status, err error) {

	if connector.client != nil {

		// Init the map
		statuses = make(map[string]*jira.Status)

		var request *http.Request

		// Create a jira api request for the statuses in JIRA
		if request, err = connector.client.NewRequest(http.MethodGet, jstatus, nil); err == nil {

			// Create a new slice of statuses
			jStatuses := new([]*jira.Status)

			var response *http.Response
			if response, err = connector.funnelClient.Do(request); err == nil {
				if response != nil {
					defer response.Body.Close()

					var body []byte
					if body, err = ioutil.ReadAll(response.Body); err == nil {
						if err = json.Unmarshal(body, jStatuses); err == nil {
							// Loop through the statuses and assign those status to the map
							for _, status := range *jStatuses {
								if status != nil {
									// Map the status
									statuses[status.Name] = status
								}
							}
						}
					}
				}
			}
		}
	}

	return statuses, err
}

// Get the complete list of Issue Types for the JIRA system and store them in a map
func (connector *ConnectorJira) getIssueTypes() (issueTypes map[string]jira.IssueType, err error) {

	if connector.client != nil {

		// Init the map
		issueTypes = make(map[string]jira.IssueType)

		var request *http.Request

		// Create a jira api request for the issuetypes in JIRA
		if request, err = connector.client.NewRequest(http.MethodGet, jissuetype, nil); err == nil {

			// Create a new slice of issuetypes
			jIssueType := new([]jira.IssueType)

			var response *http.Response
			if response, err = connector.funnelClient.Do(request); err == nil {
				if response != nil {
					defer response.Body.Close()

					var body []byte
					if body, err = ioutil.ReadAll(response.Body); err == nil {
						if err = json.Unmarshal(body, jIssueType); err == nil {
							// Loop through the issuetypes and assign those issuetype to the map
							for _, issuetype := range *jIssueType {

								// Map the issuetype
								issueTypes[issuetype.Name] = issuetype
							}
						}
					}
				}
			}
		}
	}

	return issueTypes, err
}
