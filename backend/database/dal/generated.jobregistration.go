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

// JobRegistration defines the struct that implements the JobRegistration interface
type JobRegistration struct {
	CreatedByvar   string
	CreatedDatevar time.Time
	GoStructvar    string
	IDvar          int
	Priorityvar    int
	UpdatedByvar   *string
	UpdatedDatevar *time.Time
}

//**********************************************************
// Struct Methods
//**********************************************************

// MarshalJSON marshals the struct by converting it to a map
func (myJobRegistration JobRegistration) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"CreatedBy":   myJobRegistration.CreatedByvar,
		"CreatedDate": myJobRegistration.CreatedDatevar,
		"GoStruct":    myJobRegistration.GoStructvar,
		"ID":          myJobRegistration.IDvar,
		"Priority":    myJobRegistration.Priorityvar,
		"UpdatedBy":   myJobRegistration.UpdatedByvar,
		"UpdatedDate": myJobRegistration.UpdatedDatevar,
	})
}

// CreatedBy returns the CreatedBy parameter from the JobRegistration struct
func (myJobRegistration *JobRegistration) CreatedBy() (param string) {
	return myJobRegistration.CreatedByvar
}

// CreatedDate returns the CreatedDate parameter from the JobRegistration struct
func (myJobRegistration *JobRegistration) CreatedDate() (param time.Time) {
	return myJobRegistration.CreatedDatevar
}

// GoStruct returns the GoStruct parameter from the JobRegistration struct
func (myJobRegistration *JobRegistration) GoStruct() (param string) {
	return myJobRegistration.GoStructvar
}

// ID returns the ID parameter from the JobRegistration struct
func (myJobRegistration *JobRegistration) ID() (param int) {
	return myJobRegistration.IDvar
}

// Priority returns the Priority parameter from the JobRegistration struct
func (myJobRegistration *JobRegistration) Priority() (param int) {
	return myJobRegistration.Priorityvar
}

// UpdatedBy returns the UpdatedBy parameter from the JobRegistration struct
func (myJobRegistration *JobRegistration) UpdatedBy() (param *string) {
	return myJobRegistration.UpdatedByvar
}

// UpdatedDate returns the UpdatedDate parameter from the JobRegistration struct
func (myJobRegistration *JobRegistration) UpdatedDate() (param *time.Time) {
	return myJobRegistration.UpdatedDatevar
}

// SetCreatedBy sets the CreatedBy parameter from the JobRegistration struct
func (myJobRegistration *JobRegistration) SetCreatedBy(val string) {
	myJobRegistration.CreatedByvar = val
}

// SetGoStruct sets the GoStruct parameter from the JobRegistration struct
func (myJobRegistration *JobRegistration) SetGoStruct(val string) {
	myJobRegistration.GoStructvar = val
}

// SetUpdatedBy sets the UpdatedBy parameter from the JobRegistration struct
func (myJobRegistration *JobRegistration) SetUpdatedBy(val string) {
	myJobRegistration.UpdatedByvar = &val
}
