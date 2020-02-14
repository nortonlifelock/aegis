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

// SourceConfig defines the struct that implements the SourceConfig interface
type SourceConfig struct {
	Addressvar        string
	AuthInfovar       string
	DBCreatedDatevar  time.Time
	DBUpdatedDatevar  *time.Time
	IDvar             string
	OrganizationIDvar string
	Payloadvar        *string
	Portvar           string
	Sourcevar         string
	SourceIDvar       string
}

//**********************************************************
// Struct Methods
//**********************************************************

// MarshalJSON marshals the struct by converting it to a map
func (mySourceConfig SourceConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Address":        mySourceConfig.Addressvar,
		"AuthInfo":       mySourceConfig.AuthInfovar,
		"DBCreatedDate":  mySourceConfig.DBCreatedDatevar,
		"DBUpdatedDate":  mySourceConfig.DBUpdatedDatevar,
		"ID":             mySourceConfig.IDvar,
		"OrganizationID": mySourceConfig.OrganizationIDvar,
		"Payload":        mySourceConfig.Payloadvar,
		"Port":           mySourceConfig.Portvar,
		"Source":         mySourceConfig.Sourcevar,
		"SourceID":       mySourceConfig.SourceIDvar,
	})
}

// Address returns the Address parameter from the SourceConfig struct
func (mySourceConfig *SourceConfig) Address() (param string) {
	return mySourceConfig.Addressvar
}

// AuthInfo returns the AuthInfo parameter from the SourceConfig struct
func (mySourceConfig *SourceConfig) AuthInfo() (param string) {
	return mySourceConfig.AuthInfovar
}

// DBCreatedDate returns the DBCreatedDate parameter from the SourceConfig struct
func (mySourceConfig *SourceConfig) DBCreatedDate() (param time.Time) {
	return mySourceConfig.DBCreatedDatevar
}

// DBUpdatedDate returns the DBUpdatedDate parameter from the SourceConfig struct
func (mySourceConfig *SourceConfig) DBUpdatedDate() (param *time.Time) {
	return mySourceConfig.DBUpdatedDatevar
}

// ID returns the ID parameter from the SourceConfig struct
func (mySourceConfig *SourceConfig) ID() (param string) {
	return mySourceConfig.IDvar
}

// OrganizationID returns the OrganizationID parameter from the SourceConfig struct
func (mySourceConfig *SourceConfig) OrganizationID() (param string) {
	return mySourceConfig.OrganizationIDvar
}

// Payload returns the Payload parameter from the SourceConfig struct
func (mySourceConfig *SourceConfig) Payload() (param *string) {
	return mySourceConfig.Payloadvar
}

// Port returns the Port parameter from the SourceConfig struct
func (mySourceConfig *SourceConfig) Port() (param string) {
	return mySourceConfig.Portvar
}

// Source returns the Source parameter from the SourceConfig struct
func (mySourceConfig *SourceConfig) Source() (param string) {
	return mySourceConfig.Sourcevar
}

// SourceID returns the SourceID parameter from the SourceConfig struct
func (mySourceConfig *SourceConfig) SourceID() (param string) {
	return mySourceConfig.SourceIDvar
}

// SetAddress sets the Address parameter from the SourceConfig struct
func (mySourceConfig *SourceConfig) SetAddress(val string) {
	mySourceConfig.Addressvar = val
}

// SetID sets the ID parameter from the SourceConfig struct
func (mySourceConfig *SourceConfig) SetID(val string) {
	mySourceConfig.IDvar = val
}

// SetOrganizationID sets the OrganizationID parameter from the SourceConfig struct
func (mySourceConfig *SourceConfig) SetOrganizationID(val string) {
	mySourceConfig.OrganizationIDvar = val
}

// SetPayload sets the Payload parameter from the SourceConfig struct
func (mySourceConfig *SourceConfig) SetPayload(val string) {
	mySourceConfig.Payloadvar = &val
}

// SetPort sets the Port parameter from the SourceConfig struct
func (mySourceConfig *SourceConfig) SetPort(val string) {
	mySourceConfig.Portvar = val
}

// SetSource sets the Source parameter from the SourceConfig struct
func (mySourceConfig *SourceConfig) SetSource(val string) {
	mySourceConfig.Sourcevar = val
}

// SetSourceID sets the SourceID parameter from the SourceConfig struct
func (mySourceConfig *SourceConfig) SetSourceID(val string) {
	mySourceConfig.SourceIDvar = val
}
