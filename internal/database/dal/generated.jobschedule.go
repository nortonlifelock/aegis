package dal

//**********************************************************
// GENERATED CODE - DO NOT CHANGE
// This file is generated using scaffolding. Any changes to
// this file will be overwritten on the next build
//**********************************************************

import (
	"encoding/json"
)

//**********************************************************
// Struct Declaration
//**********************************************************

// JobSchedule defines the struct that implements the JobSchedule interface
type JobSchedule struct {
	ConfigIDvar string
	IDvar       string
	Payloadvar  *string
}

//**********************************************************
// Struct Methods
//**********************************************************

// MarshalJSON marshals the struct by converting it to a map
func (myJobSchedule JobSchedule) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"ConfigID": myJobSchedule.ConfigIDvar,
		"ID":       myJobSchedule.IDvar,
		"Payload":  myJobSchedule.Payloadvar,
	})
}

// ConfigID returns the ConfigID parameter from the JobSchedule struct
func (myJobSchedule *JobSchedule) ConfigID() (param string) {
	return myJobSchedule.ConfigIDvar
}

// ID returns the ID parameter from the JobSchedule struct
func (myJobSchedule *JobSchedule) ID() (param string) {
	return myJobSchedule.IDvar
}

// Payload returns the Payload parameter from the JobSchedule struct
func (myJobSchedule *JobSchedule) Payload() (param *string) {
	return myJobSchedule.Payloadvar
}

// SetConfigID sets the ConfigID parameter from the JobSchedule struct
func (myJobSchedule *JobSchedule) SetConfigID(val string) {
	myJobSchedule.ConfigIDvar = val
}

// SetID sets the ID parameter from the JobSchedule struct
func (myJobSchedule *JobSchedule) SetID(val string) {
	myJobSchedule.IDvar = val
}
