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

// JobConfig defines the struct that implements the JobConfig interface
type JobConfig struct {
	Activevar                bool
	AutoStartvar             bool
	Continuousvar            bool
	CreatedByvar             string
	CreatedDatevar           time.Time
	DataInSourceConfigIDvar  *string
	DataOutSourceConfigIDvar *string
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
func (myJobConfig JobConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Active":                myJobConfig.Activevar,
		"AutoStart":             myJobConfig.AutoStartvar,
		"Continuous":            myJobConfig.Continuousvar,
		"CreatedBy":             myJobConfig.CreatedByvar,
		"CreatedDate":           myJobConfig.CreatedDatevar,
		"DataInSourceConfigID":  myJobConfig.DataInSourceConfigIDvar,
		"DataOutSourceConfigID": myJobConfig.DataOutSourceConfigIDvar,
		"ID":                    myJobConfig.IDvar,
		"JobID":                 myJobConfig.JobIDvar,
		"LastJobStart":          myJobConfig.LastJobStartvar,
		"MaxInstances":          myJobConfig.MaxInstancesvar,
		"OrganizationID":        myJobConfig.OrganizationIDvar,
		"Payload":               myJobConfig.Payloadvar,
		"PriorityOverride":      myJobConfig.PriorityOverridevar,
		"UpdatedBy":             myJobConfig.UpdatedByvar,
		"UpdatedDate":           myJobConfig.UpdatedDatevar,
		"WaitInSeconds":         myJobConfig.WaitInSecondsvar,
	})
}

// Active returns the Active parameter from the JobConfig struct
func (myJobConfig *JobConfig) Active() (param bool) {
	return myJobConfig.Activevar
}

// AutoStart returns the AutoStart parameter from the JobConfig struct
func (myJobConfig *JobConfig) AutoStart() (param bool) {
	return myJobConfig.AutoStartvar
}

// Continuous returns the Continuous parameter from the JobConfig struct
func (myJobConfig *JobConfig) Continuous() (param bool) {
	return myJobConfig.Continuousvar
}

// CreatedBy returns the CreatedBy parameter from the JobConfig struct
func (myJobConfig *JobConfig) CreatedBy() (param string) {
	return myJobConfig.CreatedByvar
}

// CreatedDate returns the CreatedDate parameter from the JobConfig struct
func (myJobConfig *JobConfig) CreatedDate() (param time.Time) {
	return myJobConfig.CreatedDatevar
}

// DataInSourceConfigID returns the DataInSourceConfigID parameter from the JobConfig struct
func (myJobConfig *JobConfig) DataInSourceConfigID() (param *string) {
	return myJobConfig.DataInSourceConfigIDvar
}

// DataOutSourceConfigID returns the DataOutSourceConfigID parameter from the JobConfig struct
func (myJobConfig *JobConfig) DataOutSourceConfigID() (param *string) {
	return myJobConfig.DataOutSourceConfigIDvar
}

// ID returns the ID parameter from the JobConfig struct
func (myJobConfig *JobConfig) ID() (param string) {
	return myJobConfig.IDvar
}

// JobID returns the JobID parameter from the JobConfig struct
func (myJobConfig *JobConfig) JobID() (param int) {
	return myJobConfig.JobIDvar
}

// LastJobStart returns the LastJobStart parameter from the JobConfig struct
func (myJobConfig *JobConfig) LastJobStart() (param *time.Time) {
	return myJobConfig.LastJobStartvar
}

// MaxInstances returns the MaxInstances parameter from the JobConfig struct
func (myJobConfig *JobConfig) MaxInstances() (param int) {
	return myJobConfig.MaxInstancesvar
}

// OrganizationID returns the OrganizationID parameter from the JobConfig struct
func (myJobConfig *JobConfig) OrganizationID() (param string) {
	return myJobConfig.OrganizationIDvar
}

// Payload returns the Payload parameter from the JobConfig struct
func (myJobConfig *JobConfig) Payload() (param *string) {
	return myJobConfig.Payloadvar
}

// PriorityOverride returns the PriorityOverride parameter from the JobConfig struct
func (myJobConfig *JobConfig) PriorityOverride() (param *int) {
	return myJobConfig.PriorityOverridevar
}

// UpdatedBy returns the UpdatedBy parameter from the JobConfig struct
func (myJobConfig *JobConfig) UpdatedBy() (param *string) {
	return myJobConfig.UpdatedByvar
}

// UpdatedDate returns the UpdatedDate parameter from the JobConfig struct
func (myJobConfig *JobConfig) UpdatedDate() (param *time.Time) {
	return myJobConfig.UpdatedDatevar
}

// WaitInSeconds returns the WaitInSeconds parameter from the JobConfig struct
func (myJobConfig *JobConfig) WaitInSeconds() (param int) {
	return myJobConfig.WaitInSecondsvar
}

// SetCreatedBy sets the CreatedBy parameter from the JobConfig struct
func (myJobConfig *JobConfig) SetCreatedBy(val string) {
	myJobConfig.CreatedByvar = val
}

// SetDataInSourceConfigID sets the DataInSourceConfigID parameter from the JobConfig struct
func (myJobConfig *JobConfig) SetDataInSourceConfigID(val string) {
	myJobConfig.DataInSourceConfigIDvar = &val
}

// SetDataOutSourceConfigID sets the DataOutSourceConfigID parameter from the JobConfig struct
func (myJobConfig *JobConfig) SetDataOutSourceConfigID(val string) {
	myJobConfig.DataOutSourceConfigIDvar = &val
}

// SetID sets the ID parameter from the JobConfig struct
func (myJobConfig *JobConfig) SetID(val string) {
	myJobConfig.IDvar = val
}

// SetOrganizationID sets the OrganizationID parameter from the JobConfig struct
func (myJobConfig *JobConfig) SetOrganizationID(val string) {
	myJobConfig.OrganizationIDvar = val
}

// SetPayload sets the Payload parameter from the JobConfig struct
func (myJobConfig *JobConfig) SetPayload(val string) {
	myJobConfig.Payloadvar = &val
}

// SetUpdatedBy sets the UpdatedBy parameter from the JobConfig struct
func (myJobConfig *JobConfig) SetUpdatedBy(val string) {
	myJobConfig.UpdatedByvar = &val
}
