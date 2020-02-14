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

// JobConfigAudit defines the struct that implements the JobConfigAudit interface
type JobConfigAudit struct {
	Activevar                bool
	AutoStartvar             bool
	Continuousvar            bool
	CreatedByvar             string
	CreatedDatevar           time.Time
	DataInSourceConfigIDvar  *string
	DataOutSourceConfigIDvar *string
	EventDatevar             time.Time
	EventTypevar             string
	IDvar                    string
	JobIDvar                 int
	LastJobStartvar          *time.Time
	MaxInstancesvar          int
	OrganizationIDvar        string
	Payloadvar               *string
	PriorityOverridevar      *int
	UpdatedByvar             *string
	UpdatedDatevar           *time.Time
	WaitInSecondsvar         int
}

//**********************************************************
// Struct Methods
//**********************************************************

// MarshalJSON marshals the struct by converting it to a map
func (myJobConfigAudit JobConfigAudit) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Active":                myJobConfigAudit.Activevar,
		"AutoStart":             myJobConfigAudit.AutoStartvar,
		"Continuous":            myJobConfigAudit.Continuousvar,
		"CreatedBy":             myJobConfigAudit.CreatedByvar,
		"CreatedDate":           myJobConfigAudit.CreatedDatevar,
		"DataInSourceConfigID":  myJobConfigAudit.DataInSourceConfigIDvar,
		"DataOutSourceConfigID": myJobConfigAudit.DataOutSourceConfigIDvar,
		"EventDate":             myJobConfigAudit.EventDatevar,
		"EventType":             myJobConfigAudit.EventTypevar,
		"ID":                    myJobConfigAudit.IDvar,
		"JobID":                 myJobConfigAudit.JobIDvar,
		"LastJobStart":          myJobConfigAudit.LastJobStartvar,
		"MaxInstances":          myJobConfigAudit.MaxInstancesvar,
		"OrganizationID":        myJobConfigAudit.OrganizationIDvar,
		"Payload":               myJobConfigAudit.Payloadvar,
		"PriorityOverride":      myJobConfigAudit.PriorityOverridevar,
		"UpdatedBy":             myJobConfigAudit.UpdatedByvar,
		"UpdatedDate":           myJobConfigAudit.UpdatedDatevar,
		"WaitInSeconds":         myJobConfigAudit.WaitInSecondsvar,
	})
}

// Active returns the Active parameter from the JobConfigAudit struct
func (myJobConfigAudit *JobConfigAudit) Active() (param bool) {
	return myJobConfigAudit.Activevar
}

// AutoStart returns the AutoStart parameter from the JobConfigAudit struct
func (myJobConfigAudit *JobConfigAudit) AutoStart() (param bool) {
	return myJobConfigAudit.AutoStartvar
}

// Continuous returns the Continuous parameter from the JobConfigAudit struct
func (myJobConfigAudit *JobConfigAudit) Continuous() (param bool) {
	return myJobConfigAudit.Continuousvar
}

// CreatedBy returns the CreatedBy parameter from the JobConfigAudit struct
func (myJobConfigAudit *JobConfigAudit) CreatedBy() (param string) {
	return myJobConfigAudit.CreatedByvar
}

// CreatedDate returns the CreatedDate parameter from the JobConfigAudit struct
func (myJobConfigAudit *JobConfigAudit) CreatedDate() (param time.Time) {
	return myJobConfigAudit.CreatedDatevar
}

// DataInSourceConfigID returns the DataInSourceConfigID parameter from the JobConfigAudit struct
func (myJobConfigAudit *JobConfigAudit) DataInSourceConfigID() (param *string) {
	return myJobConfigAudit.DataInSourceConfigIDvar
}

// DataOutSourceConfigID returns the DataOutSourceConfigID parameter from the JobConfigAudit struct
func (myJobConfigAudit *JobConfigAudit) DataOutSourceConfigID() (param *string) {
	return myJobConfigAudit.DataOutSourceConfigIDvar
}

// EventDate returns the EventDate parameter from the JobConfigAudit struct
func (myJobConfigAudit *JobConfigAudit) EventDate() (param time.Time) {
	return myJobConfigAudit.EventDatevar
}

// EventType returns the EventType parameter from the JobConfigAudit struct
func (myJobConfigAudit *JobConfigAudit) EventType() (param string) {
	return myJobConfigAudit.EventTypevar
}

// ID returns the ID parameter from the JobConfigAudit struct
func (myJobConfigAudit *JobConfigAudit) ID() (param string) {
	return myJobConfigAudit.IDvar
}

// JobID returns the JobID parameter from the JobConfigAudit struct
func (myJobConfigAudit *JobConfigAudit) JobID() (param int) {
	return myJobConfigAudit.JobIDvar
}

// LastJobStart returns the LastJobStart parameter from the JobConfigAudit struct
func (myJobConfigAudit *JobConfigAudit) LastJobStart() (param *time.Time) {
	return myJobConfigAudit.LastJobStartvar
}

// MaxInstances returns the MaxInstances parameter from the JobConfigAudit struct
func (myJobConfigAudit *JobConfigAudit) MaxInstances() (param int) {
	return myJobConfigAudit.MaxInstancesvar
}

// OrganizationID returns the OrganizationID parameter from the JobConfigAudit struct
func (myJobConfigAudit *JobConfigAudit) OrganizationID() (param string) {
	return myJobConfigAudit.OrganizationIDvar
}

// Payload returns the Payload parameter from the JobConfigAudit struct
func (myJobConfigAudit *JobConfigAudit) Payload() (param *string) {
	return myJobConfigAudit.Payloadvar
}

// PriorityOverride returns the PriorityOverride parameter from the JobConfigAudit struct
func (myJobConfigAudit *JobConfigAudit) PriorityOverride() (param *int) {
	return myJobConfigAudit.PriorityOverridevar
}

// UpdatedBy returns the UpdatedBy parameter from the JobConfigAudit struct
func (myJobConfigAudit *JobConfigAudit) UpdatedBy() (param *string) {
	return myJobConfigAudit.UpdatedByvar
}

// UpdatedDate returns the UpdatedDate parameter from the JobConfigAudit struct
func (myJobConfigAudit *JobConfigAudit) UpdatedDate() (param *time.Time) {
	return myJobConfigAudit.UpdatedDatevar
}

// WaitInSeconds returns the WaitInSeconds parameter from the JobConfigAudit struct
func (myJobConfigAudit *JobConfigAudit) WaitInSeconds() (param int) {
	return myJobConfigAudit.WaitInSecondsvar
}

// SetCreatedBy sets the CreatedBy parameter from the JobConfigAudit struct
func (myJobConfigAudit *JobConfigAudit) SetCreatedBy(val string) {
	myJobConfigAudit.CreatedByvar = val
}

// SetDataInSourceConfigID sets the DataInSourceConfigID parameter from the JobConfigAudit struct
func (myJobConfigAudit *JobConfigAudit) SetDataInSourceConfigID(val string) {
	myJobConfigAudit.DataInSourceConfigIDvar = &val
}

// SetDataOutSourceConfigID sets the DataOutSourceConfigID parameter from the JobConfigAudit struct
func (myJobConfigAudit *JobConfigAudit) SetDataOutSourceConfigID(val string) {
	myJobConfigAudit.DataOutSourceConfigIDvar = &val
}

// SetEventType sets the EventType parameter from the JobConfigAudit struct
func (myJobConfigAudit *JobConfigAudit) SetEventType(val string) {
	myJobConfigAudit.EventTypevar = val
}

// SetID sets the ID parameter from the JobConfigAudit struct
func (myJobConfigAudit *JobConfigAudit) SetID(val string) {
	myJobConfigAudit.IDvar = val
}

// SetOrganizationID sets the OrganizationID parameter from the JobConfigAudit struct
func (myJobConfigAudit *JobConfigAudit) SetOrganizationID(val string) {
	myJobConfigAudit.OrganizationIDvar = val
}

// SetUpdatedBy sets the UpdatedBy parameter from the JobConfigAudit struct
func (myJobConfigAudit *JobConfigAudit) SetUpdatedBy(val string) {
	myJobConfigAudit.UpdatedByvar = &val
}
