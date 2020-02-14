package dal

//**********************************************************
// GENERATED CODE - DO NOT CHANGE
// This file is generated using scaffolding. Any changes to
// this file will be overwritten on the next build
//**********************************************************

import (
	"encoding/json"
	"time"
)

//**********************************************************
// Struct Declaration
//**********************************************************

// JobHistory defines the struct that implements the JobHistory interface
type JobHistory struct {
	ConfigIDvar         string
	CreatedDatevar      time.Time
	CurrentIterationvar *int
	IDvar               string
	Identifiervar       *string
	JobIDvar            int
	MaxInstancesvar     int
	ParentJobIDvar      *string
	Payloadvar          string
	Priorityvar         int
	PulseDatevar        *time.Time
	StatusIDvar         int
	ThreadIDvar         *string
	UpdatedDatevar      *time.Time
}

//**********************************************************
// Struct Methods
//**********************************************************

// MarshalJSON marshals the struct by converting it to a map
func (myJobHistory JobHistory) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"ConfigID":         myJobHistory.ConfigIDvar,
		"CreatedDate":      myJobHistory.CreatedDatevar,
		"CurrentIteration": myJobHistory.CurrentIterationvar,
		"ID":               myJobHistory.IDvar,
		"Identifier":       myJobHistory.Identifiervar,
		"JobID":            myJobHistory.JobIDvar,
		"MaxInstances":     myJobHistory.MaxInstancesvar,
		"ParentJobID":      myJobHistory.ParentJobIDvar,
		"Payload":          myJobHistory.Payloadvar,
		"Priority":         myJobHistory.Priorityvar,
		"PulseDate":        myJobHistory.PulseDatevar,
		"StatusID":         myJobHistory.StatusIDvar,
		"ThreadID":         myJobHistory.ThreadIDvar,
		"UpdatedDate":      myJobHistory.UpdatedDatevar,
	})
}

// ConfigID returns the ConfigID parameter from the JobHistory struct
func (myJobHistory *JobHistory) ConfigID() (param string) {
	return myJobHistory.ConfigIDvar
}

// CreatedDate returns the CreatedDate parameter from the JobHistory struct
func (myJobHistory *JobHistory) CreatedDate() (param time.Time) {
	return myJobHistory.CreatedDatevar
}

// CurrentIteration returns the CurrentIteration parameter from the JobHistory struct
func (myJobHistory *JobHistory) CurrentIteration() (param *int) {
	return myJobHistory.CurrentIterationvar
}

// ID returns the ID parameter from the JobHistory struct
func (myJobHistory *JobHistory) ID() (param string) {
	return myJobHistory.IDvar
}

// Identifier returns the Identifier parameter from the JobHistory struct
func (myJobHistory *JobHistory) Identifier() (param *string) {
	return myJobHistory.Identifiervar
}

// JobID returns the JobID parameter from the JobHistory struct
func (myJobHistory *JobHistory) JobID() (param int) {
	return myJobHistory.JobIDvar
}

// MaxInstances returns the MaxInstances parameter from the JobHistory struct
func (myJobHistory *JobHistory) MaxInstances() (param int) {
	return myJobHistory.MaxInstancesvar
}

// ParentJobID returns the ParentJobID parameter from the JobHistory struct
func (myJobHistory *JobHistory) ParentJobID() (param *string) {
	return myJobHistory.ParentJobIDvar
}

// Payload returns the Payload parameter from the JobHistory struct
func (myJobHistory *JobHistory) Payload() (param string) {
	return myJobHistory.Payloadvar
}

// Priority returns the Priority parameter from the JobHistory struct
func (myJobHistory *JobHistory) Priority() (param int) {
	return myJobHistory.Priorityvar
}

// PulseDate returns the PulseDate parameter from the JobHistory struct
func (myJobHistory *JobHistory) PulseDate() (param *time.Time) {
	return myJobHistory.PulseDatevar
}

// StatusID returns the StatusID parameter from the JobHistory struct
func (myJobHistory *JobHistory) StatusID() (param int) {
	return myJobHistory.StatusIDvar
}

// ThreadID returns the ThreadID parameter from the JobHistory struct
func (myJobHistory *JobHistory) ThreadID() (param *string) {
	return myJobHistory.ThreadIDvar
}

// UpdatedDate returns the UpdatedDate parameter from the JobHistory struct
func (myJobHistory *JobHistory) UpdatedDate() (param *time.Time) {
	return myJobHistory.UpdatedDatevar
}

// SetConfigID sets the ConfigID parameter from the JobHistory struct
func (myJobHistory *JobHistory) SetConfigID(val string) {
	myJobHistory.ConfigIDvar = val
}

// SetID sets the ID parameter from the JobHistory struct
func (myJobHistory *JobHistory) SetID(val string) {
	myJobHistory.IDvar = val
}

// SetIdentifier sets the Identifier parameter from the JobHistory struct
func (myJobHistory *JobHistory) SetIdentifier(val string) {
	myJobHistory.Identifiervar = &val
}

// SetParentJobID sets the ParentJobID parameter from the JobHistory struct
func (myJobHistory *JobHistory) SetParentJobID(val string) {
	myJobHistory.ParentJobIDvar = &val
}

// SetThreadID sets the ThreadID parameter from the JobHistory struct
func (myJobHistory *JobHistory) SetThreadID(val string) {
	myJobHistory.ThreadIDvar = &val
}
