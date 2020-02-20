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

// ExceptedDetection defines the struct that implements the ExceptedDetection interface
type ExceptedDetection struct {
	Approvalvar           *string
	AssignmentGroupvar    *string
	DueDatevar            *time.Time
	Hostnamevar           *string
	IPvar                 *string
	IgnoreIDvar           string
	IgnoreTypevar         int
	OSvar                 *string
	OSRegexvar            *string
	Titlevar              *string
	VulnerabilityIDvar    *string
	VulnerabilityTitlevar *string
}

//**********************************************************
// Struct Methods
//**********************************************************

// MarshalJSON marshals the struct by converting it to a map
func (myExceptedDetection ExceptedDetection) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Approval":           myExceptedDetection.Approvalvar,
		"AssignmentGroup":    myExceptedDetection.AssignmentGroupvar,
		"DueDate":            myExceptedDetection.DueDatevar,
		"Hostname":           myExceptedDetection.Hostnamevar,
		"IP":                 myExceptedDetection.IPvar,
		"IgnoreID":           myExceptedDetection.IgnoreIDvar,
		"IgnoreType":         myExceptedDetection.IgnoreTypevar,
		"OS":                 myExceptedDetection.OSvar,
		"OSRegex":            myExceptedDetection.OSRegexvar,
		"Title":              myExceptedDetection.Titlevar,
		"VulnerabilityID":    myExceptedDetection.VulnerabilityIDvar,
		"VulnerabilityTitle": myExceptedDetection.VulnerabilityTitlevar,
	})
}

// Approval returns the Approval parameter from the ExceptedDetection struct
func (myExceptedDetection *ExceptedDetection) Approval() (param *string) {
	return myExceptedDetection.Approvalvar
}

// AssignmentGroup returns the AssignmentGroup parameter from the ExceptedDetection struct
func (myExceptedDetection *ExceptedDetection) AssignmentGroup() (param *string) {
	return myExceptedDetection.AssignmentGroupvar
}

// DueDate returns the DueDate parameter from the ExceptedDetection struct
func (myExceptedDetection *ExceptedDetection) DueDate() (param *time.Time) {
	return myExceptedDetection.DueDatevar
}

// Hostname returns the Hostname parameter from the ExceptedDetection struct
func (myExceptedDetection *ExceptedDetection) Hostname() (param *string) {
	return myExceptedDetection.Hostnamevar
}

// IP returns the IP parameter from the ExceptedDetection struct
func (myExceptedDetection *ExceptedDetection) IP() (param *string) {
	return myExceptedDetection.IPvar
}

// IgnoreID returns the IgnoreID parameter from the ExceptedDetection struct
func (myExceptedDetection *ExceptedDetection) IgnoreID() (param string) {
	return myExceptedDetection.IgnoreIDvar
}

// IgnoreType returns the IgnoreType parameter from the ExceptedDetection struct
func (myExceptedDetection *ExceptedDetection) IgnoreType() (param int) {
	return myExceptedDetection.IgnoreTypevar
}

// OS returns the OS parameter from the ExceptedDetection struct
func (myExceptedDetection *ExceptedDetection) OS() (param *string) {
	return myExceptedDetection.OSvar
}

// OSRegex returns the OSRegex parameter from the ExceptedDetection struct
func (myExceptedDetection *ExceptedDetection) OSRegex() (param *string) {
	return myExceptedDetection.OSRegexvar
}

// Title returns the Title parameter from the ExceptedDetection struct
func (myExceptedDetection *ExceptedDetection) Title() (param *string) {
	return myExceptedDetection.Titlevar
}

// VulnerabilityID returns the VulnerabilityID parameter from the ExceptedDetection struct
func (myExceptedDetection *ExceptedDetection) VulnerabilityID() (param *string) {
	return myExceptedDetection.VulnerabilityIDvar
}

// VulnerabilityTitle returns the VulnerabilityTitle parameter from the ExceptedDetection struct
func (myExceptedDetection *ExceptedDetection) VulnerabilityTitle() (param *string) {
	return myExceptedDetection.VulnerabilityTitlevar
}

// SetApproval sets the Approval parameter from the ExceptedDetection struct
func (myExceptedDetection *ExceptedDetection) SetApproval(val string) {
	myExceptedDetection.Approvalvar = &val
}

// SetAssignmentGroup sets the AssignmentGroup parameter from the ExceptedDetection struct
func (myExceptedDetection *ExceptedDetection) SetAssignmentGroup(val string) {
	myExceptedDetection.AssignmentGroupvar = &val
}

// SetHostname sets the Hostname parameter from the ExceptedDetection struct
func (myExceptedDetection *ExceptedDetection) SetHostname(val string) {
	myExceptedDetection.Hostnamevar = &val
}

// SetIP sets the IP parameter from the ExceptedDetection struct
func (myExceptedDetection *ExceptedDetection) SetIP(val string) {
	myExceptedDetection.IPvar = &val
}

// SetIgnoreID sets the IgnoreID parameter from the ExceptedDetection struct
func (myExceptedDetection *ExceptedDetection) SetIgnoreID(val string) {
	myExceptedDetection.IgnoreIDvar = val
}

// SetOS sets the OS parameter from the ExceptedDetection struct
func (myExceptedDetection *ExceptedDetection) SetOS(val string) {
	myExceptedDetection.OSvar = &val
}

// SetOSRegex sets the OSRegex parameter from the ExceptedDetection struct
func (myExceptedDetection *ExceptedDetection) SetOSRegex(val string) {
	myExceptedDetection.OSRegexvar = &val
}

// SetTitle sets the Title parameter from the ExceptedDetection struct
func (myExceptedDetection *ExceptedDetection) SetTitle(val string) {
	myExceptedDetection.Titlevar = &val
}

// SetVulnerabilityID sets the VulnerabilityID parameter from the ExceptedDetection struct
func (myExceptedDetection *ExceptedDetection) SetVulnerabilityID(val string) {
	myExceptedDetection.VulnerabilityIDvar = &val
}

// SetVulnerabilityTitle sets the VulnerabilityTitle parameter from the ExceptedDetection struct
func (myExceptedDetection *ExceptedDetection) SetVulnerabilityTitle(val string) {
	myExceptedDetection.VulnerabilityTitlevar = &val
}
