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

// Organization defines the struct that implements the Organization interface
type Organization struct {
	Codevar           string
	Createdvar        time.Time
	Descriptionvar    *string
	EncryptionKeyvar  *string
	IDvar             string
	ParentOrgIDvar    *string
	Payloadvar        string
	TimeZoneOffsetvar float32
	Updatedvar        *time.Time
}

//**********************************************************
// Struct Methods
//**********************************************************

// MarshalJSON marshals the struct by converting it to a map
func (myOrganization Organization) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Code":           myOrganization.Codevar,
		"Created":        myOrganization.Createdvar,
		"Description":    myOrganization.Descriptionvar,
		"EncryptionKey":  myOrganization.EncryptionKeyvar,
		"ID":             myOrganization.IDvar,
		"ParentOrgID":    myOrganization.ParentOrgIDvar,
		"Payload":        myOrganization.Payloadvar,
		"TimeZoneOffset": myOrganization.TimeZoneOffsetvar,
		"Updated":        myOrganization.Updatedvar,
	})
}

// Code returns the Code parameter from the Organization struct
func (myOrganization *Organization) Code() (param string) {
	return myOrganization.Codevar
}

// Created returns the Created parameter from the Organization struct
func (myOrganization *Organization) Created() (param time.Time) {
	return myOrganization.Createdvar
}

// Description returns the Description parameter from the Organization struct
func (myOrganization *Organization) Description() (param *string) {
	return myOrganization.Descriptionvar
}

// EncryptionKey returns the EncryptionKey parameter from the Organization struct
func (myOrganization *Organization) EncryptionKey() (param *string) {
	return myOrganization.EncryptionKeyvar
}

// ID returns the ID parameter from the Organization struct
func (myOrganization *Organization) ID() (param string) {
	return myOrganization.IDvar
}

// ParentOrgID returns the ParentOrgID parameter from the Organization struct
func (myOrganization *Organization) ParentOrgID() (param *string) {
	return myOrganization.ParentOrgIDvar
}

// Payload returns the Payload parameter from the Organization struct
func (myOrganization *Organization) Payload() (param string) {
	return myOrganization.Payloadvar
}

// TimeZoneOffset returns the TimeZoneOffset parameter from the Organization struct
func (myOrganization *Organization) TimeZoneOffset() (param float32) {
	return myOrganization.TimeZoneOffsetvar
}

// Updated returns the Updated parameter from the Organization struct
func (myOrganization *Organization) Updated() (param *time.Time) {
	return myOrganization.Updatedvar
}

// SetCode sets the Code parameter from the Organization struct
func (myOrganization *Organization) SetCode(val string) {
	myOrganization.Codevar = val
}

// SetDescription sets the Description parameter from the Organization struct
func (myOrganization *Organization) SetDescription(val string) {
	myOrganization.Descriptionvar = &val
}

// SetID sets the ID parameter from the Organization struct
func (myOrganization *Organization) SetID(val string) {
	myOrganization.IDvar = val
}

// SetParentOrgID sets the ParentOrgID parameter from the Organization struct
func (myOrganization *Organization) SetParentOrgID(val string) {
	myOrganization.ParentOrgIDvar = &val
}
