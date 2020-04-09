package jira

import "github.com/trivago/tgo/tcontainer"

// transitionResult is used to capture the information regarding a JIRA transition, which is useful for finding which fields must be included
// to execute the transition
type transitionResult struct {
	Transitions []Transition `json:"transitions" structs:"transitions"`
}

// Transition is used in transitionResult and must be exported in order to be marshaled
type Transition struct {
	ID     string                     `json:"id"     structs:"id"`
	Name   string                     `json:"name"   structs:"name"`
	To     Status                     `json:"to"     structs:"status"`
	Fields map[string]TransitionField `json:"fields" structs:"fields"`
}

// TransitionField is used in Transition and must be exported in order to be marshaled
type TransitionField struct {
	Required bool   `json:"required"`
	Name     string `json:"name"`
}

// this payload is required to transition the status of a JIRA ticket
type createTransitionPayload struct {
	Transition      TransitionPayload `json:"transition" structs:"transition"`
	Fields          *FieldStruct      `json:"fields,omitempty"`
	FieldsInterface interface{}       `json:"fields,omitempty"`
	UpdateBlock     Update            `json:"update,omitempty"`
}

// FieldStruct is used in createTransitionPayload and must be exported in order to be marshaled
type FieldStruct struct {
	ReopenReason   string    `json:"reopen_reason,omitempty"`
	ResolutionDate string    `json:"resolution_date,omitempty"`
	Assignee       *Assignee `json:"assignee,omitempty"`
}

// Assignee is used in FieldStruct and must be exported in order to be marshaled
type Assignee struct {
	Name string `json:"name,omitempty"`
}

// TransitionPayload is used in createTransitionPayload and must be exported in order to be marshaled
type TransitionPayload struct {
	ID       string  `json:"id" structs:"id"`
	Assignee *string `json:"assignee,omitempty"`
	Unknowns tcontainer.MarshalMap
}

// Update is used in createTransitionPayload and must be exported in order to be marshaled
type Update struct {
	Comment []UpdateObjects `json:"comment,omitempty" structs:"comment" `
}

// UpdateObjects is used in Update and must be exported in order to be marshaled
type UpdateObjects struct {
	Add AddBody `json:"add,omitempty" structs:"add" `
}

// AddBody is used in UpdateObjects and must be exported in order to be marshaled
type AddBody struct {
	Body string `json:"body" structs:"body"`
}
