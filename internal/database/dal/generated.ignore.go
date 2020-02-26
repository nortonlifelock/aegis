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

// Ignore defines the struct that implements the Ignore interface
type Ignore struct {
	Activevar          bool
	Approvalvar        string
	CreatedByvar       *string
	DBCreatedDatevar   time.Time
	DBUpdatedDatevar   *time.Time
	DeviceIDvar        string
	DueDatevar         *time.Time
	IDvar              string
	OSRegexvar         *string
	OrganizationIDvar  string
	Portvar            string
	SourceIDvar        string
	TypeIDvar          int
	UpdatedByvar       *string
	VulnerabilityIDvar string
}

//**********************************************************
// Struct Methods
//**********************************************************

// MarshalJSON marshals the struct by converting it to a map
func (myIgnore Ignore) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Active":          myIgnore.Activevar,
		"Approval":        myIgnore.Approvalvar,
		"CreatedBy":       myIgnore.CreatedByvar,
		"DBCreatedDate":   myIgnore.DBCreatedDatevar,
		"DBUpdatedDate":   myIgnore.DBUpdatedDatevar,
		"DeviceID":        myIgnore.DeviceIDvar,
		"DueDate":         myIgnore.DueDatevar,
		"ID":              myIgnore.IDvar,
		"OSRegex":         myIgnore.OSRegexvar,
		"OrganizationID":  myIgnore.OrganizationIDvar,
		"Port":            myIgnore.Portvar,
		"SourceID":        myIgnore.SourceIDvar,
		"TypeID":          myIgnore.TypeIDvar,
		"UpdatedBy":       myIgnore.UpdatedByvar,
		"VulnerabilityID": myIgnore.VulnerabilityIDvar,
	})
}

// Active returns the Active parameter from the Ignore struct
func (myIgnore *Ignore) Active() (param bool) {
	return myIgnore.Activevar
}

// Approval returns the Approval parameter from the Ignore struct
func (myIgnore *Ignore) Approval() (param string) {
	return myIgnore.Approvalvar
}

// CreatedBy returns the CreatedBy parameter from the Ignore struct
func (myIgnore *Ignore) CreatedBy() (param *string) {
	return myIgnore.CreatedByvar
}

// DBCreatedDate returns the DBCreatedDate parameter from the Ignore struct
func (myIgnore *Ignore) DBCreatedDate() (param time.Time) {
	return myIgnore.DBCreatedDatevar
}

// DBUpdatedDate returns the DBUpdatedDate parameter from the Ignore struct
func (myIgnore *Ignore) DBUpdatedDate() (param *time.Time) {
	return myIgnore.DBUpdatedDatevar
}

// DeviceID returns the DeviceID parameter from the Ignore struct
func (myIgnore *Ignore) DeviceID() (param string) {
	return myIgnore.DeviceIDvar
}

// DueDate returns the DueDate parameter from the Ignore struct
func (myIgnore *Ignore) DueDate() (param *time.Time) {
	return myIgnore.DueDatevar
}

// ID returns the ID parameter from the Ignore struct
func (myIgnore *Ignore) ID() (param string) {
	return myIgnore.IDvar
}

// OSRegex returns the OSRegex parameter from the Ignore struct
func (myIgnore *Ignore) OSRegex() (param *string) {
	return myIgnore.OSRegexvar
}

// OrganizationID returns the OrganizationID parameter from the Ignore struct
func (myIgnore *Ignore) OrganizationID() (param string) {
	return myIgnore.OrganizationIDvar
}

// Port returns the Port parameter from the Ignore struct
func (myIgnore *Ignore) Port() (param string) {
	return myIgnore.Portvar
}

// SourceID returns the SourceID parameter from the Ignore struct
func (myIgnore *Ignore) SourceID() (param string) {
	return myIgnore.SourceIDvar
}

// TypeID returns the TypeID parameter from the Ignore struct
func (myIgnore *Ignore) TypeID() (param int) {
	return myIgnore.TypeIDvar
}

// UpdatedBy returns the UpdatedBy parameter from the Ignore struct
func (myIgnore *Ignore) UpdatedBy() (param *string) {
	return myIgnore.UpdatedByvar
}

// VulnerabilityID returns the VulnerabilityID parameter from the Ignore struct
func (myIgnore *Ignore) VulnerabilityID() (param string) {
	return myIgnore.VulnerabilityIDvar
}

// SetApproval sets the Approval parameter from the Ignore struct
func (myIgnore *Ignore) SetApproval(val string) {
	myIgnore.Approvalvar = val
}

// SetCreatedBy sets the CreatedBy parameter from the Ignore struct
func (myIgnore *Ignore) SetCreatedBy(val string) {
	myIgnore.CreatedByvar = &val
}

// SetDeviceID sets the DeviceID parameter from the Ignore struct
func (myIgnore *Ignore) SetDeviceID(val string) {
	myIgnore.DeviceIDvar = val
}

// SetID sets the ID parameter from the Ignore struct
func (myIgnore *Ignore) SetID(val string) {
	myIgnore.IDvar = val
}

// SetOSRegex sets the OSRegex parameter from the Ignore struct
func (myIgnore *Ignore) SetOSRegex(val string) {
	myIgnore.OSRegexvar = &val
}

// SetOrganizationID sets the OrganizationID parameter from the Ignore struct
func (myIgnore *Ignore) SetOrganizationID(val string) {
	myIgnore.OrganizationIDvar = val
}

// SetPort sets the Port parameter from the Ignore struct
func (myIgnore *Ignore) SetPort(val string) {
	myIgnore.Portvar = val
}

// SetSourceID sets the SourceID parameter from the Ignore struct
func (myIgnore *Ignore) SetSourceID(val string) {
	myIgnore.SourceIDvar = val
}

// SetUpdatedBy sets the UpdatedBy parameter from the Ignore struct
func (myIgnore *Ignore) SetUpdatedBy(val string) {
	myIgnore.UpdatedByvar = &val
}

// SetVulnerabilityID sets the VulnerabilityID parameter from the Ignore struct
func (myIgnore *Ignore) SetVulnerabilityID(val string) {
	myIgnore.VulnerabilityIDvar = val
}
