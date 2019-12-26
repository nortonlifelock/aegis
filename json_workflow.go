package jira

import "encoding/xml"

type workflowTransition struct {
	id   string
	name string
}

// workflow is used to parse the JIRA.workflow file so the code can find a series of transitions to move from one status to another
type workflow struct {
	XMLName       xml.Name `xml:"workflow"`
	CommonActions []Action `xml:"common-actions>action"`
	Statuses      []Status `xml:"steps>step"`
}

// Status is a member of workflow and must be exported in order to be marshaled
type Status struct {
	StatusName    string     `xml:"name,attr"`
	StatusID      string     `xml:"id,attr"`
	Actions       []Action   `xml:"actions>action"`
	CommonActions []ActionID `xml:"actions>common-action"`
}

// Action is a member of Status and must be exported in order to be marshaled
type Action struct {
	ActionID          string  `xml:"id,attr"`
	TransitionName    string  `xml:"name,attr"`
	TransitionDetails Results `xml:"results>unconditional-result"`
}

// ActionID is a member of Action and must be exported in order to be marshaled
type ActionID struct {
	ID string `xml:"id,attr"`
}

// Results in a member of Action and must be exported in order to be marshaled
type Results struct {
	DestinationStatusID string `xml:"step,attr"`
}
