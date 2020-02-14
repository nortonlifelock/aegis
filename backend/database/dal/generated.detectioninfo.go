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

// DetectionInfo defines the struct that implements the DetectionInfo interface
type DetectionInfo struct {
	ActiveKernelvar      *int
	AlertDatevar         time.Time
	DetectionStatusIDvar int
	DeviceIDvar          string
	IDvar                string
	IgnoreIDvar          *string
	OrganizationIDvar    string
	Portvar              int
	Proofvar             string
	Protocolvar          string
	SourceIDvar          string
	TimesSeenvar         int
	Updatedvar           time.Time
	VulnerabilityIDvar   string
}

//**********************************************************
// Struct Methods
//**********************************************************

// MarshalJSON marshals the struct by converting it to a map
func (myDetectionInfo DetectionInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"ActiveKernel":      myDetectionInfo.ActiveKernelvar,
		"AlertDate":         myDetectionInfo.AlertDatevar,
		"DetectionStatusID": myDetectionInfo.DetectionStatusIDvar,
		"DeviceID":          myDetectionInfo.DeviceIDvar,
		"ID":                myDetectionInfo.IDvar,
		"IgnoreID":          myDetectionInfo.IgnoreIDvar,
		"OrganizationID":    myDetectionInfo.OrganizationIDvar,
		"Port":              myDetectionInfo.Portvar,
		"Proof":             myDetectionInfo.Proofvar,
		"Protocol":          myDetectionInfo.Protocolvar,
		"SourceID":          myDetectionInfo.SourceIDvar,
		"TimesSeen":         myDetectionInfo.TimesSeenvar,
		"Updated":           myDetectionInfo.Updatedvar,
		"VulnerabilityID":   myDetectionInfo.VulnerabilityIDvar,
	})
}

// ActiveKernel returns the ActiveKernel parameter from the DetectionInfo struct
func (myDetectionInfo *DetectionInfo) ActiveKernel() (param *int) {
	return myDetectionInfo.ActiveKernelvar
}

// AlertDate returns the AlertDate parameter from the DetectionInfo struct
func (myDetectionInfo *DetectionInfo) AlertDate() (param time.Time) {
	return myDetectionInfo.AlertDatevar
}

// DetectionStatusID returns the DetectionStatusID parameter from the DetectionInfo struct
func (myDetectionInfo *DetectionInfo) DetectionStatusID() (param int) {
	return myDetectionInfo.DetectionStatusIDvar
}

// DeviceID returns the DeviceID parameter from the DetectionInfo struct
func (myDetectionInfo *DetectionInfo) DeviceID() (param string) {
	return myDetectionInfo.DeviceIDvar
}

// ID returns the ID parameter from the DetectionInfo struct
func (myDetectionInfo *DetectionInfo) ID() (param string) {
	return myDetectionInfo.IDvar
}

// IgnoreID returns the IgnoreID parameter from the DetectionInfo struct
func (myDetectionInfo *DetectionInfo) IgnoreID() (param *string) {
	return myDetectionInfo.IgnoreIDvar
}

// OrganizationID returns the OrganizationID parameter from the DetectionInfo struct
func (myDetectionInfo *DetectionInfo) OrganizationID() (param string) {
	return myDetectionInfo.OrganizationIDvar
}

// Port returns the Port parameter from the DetectionInfo struct
func (myDetectionInfo *DetectionInfo) Port() (param int) {
	return myDetectionInfo.Portvar
}

// Proof returns the Proof parameter from the DetectionInfo struct
func (myDetectionInfo *DetectionInfo) Proof() (param string) {
	return myDetectionInfo.Proofvar
}

// Protocol returns the Protocol parameter from the DetectionInfo struct
func (myDetectionInfo *DetectionInfo) Protocol() (param string) {
	return myDetectionInfo.Protocolvar
}

// SourceID returns the SourceID parameter from the DetectionInfo struct
func (myDetectionInfo *DetectionInfo) SourceID() (param string) {
	return myDetectionInfo.SourceIDvar
}

// TimesSeen returns the TimesSeen parameter from the DetectionInfo struct
func (myDetectionInfo *DetectionInfo) TimesSeen() (param int) {
	return myDetectionInfo.TimesSeenvar
}

// Updated returns the Updated parameter from the DetectionInfo struct
func (myDetectionInfo *DetectionInfo) Updated() (param time.Time) {
	return myDetectionInfo.Updatedvar
}

// VulnerabilityID returns the VulnerabilityID parameter from the DetectionInfo struct
func (myDetectionInfo *DetectionInfo) VulnerabilityID() (param string) {
	return myDetectionInfo.VulnerabilityIDvar
}

// SetDeviceID sets the DeviceID parameter from the DetectionInfo struct
func (myDetectionInfo *DetectionInfo) SetDeviceID(val string) {
	myDetectionInfo.DeviceIDvar = val
}

// SetID sets the ID parameter from the DetectionInfo struct
func (myDetectionInfo *DetectionInfo) SetID(val string) {
	myDetectionInfo.IDvar = val
}

// SetIgnoreID sets the IgnoreID parameter from the DetectionInfo struct
func (myDetectionInfo *DetectionInfo) SetIgnoreID(val string) {
	myDetectionInfo.IgnoreIDvar = &val
}

// SetOrganizationID sets the OrganizationID parameter from the DetectionInfo struct
func (myDetectionInfo *DetectionInfo) SetOrganizationID(val string) {
	myDetectionInfo.OrganizationIDvar = val
}

// SetProof sets the Proof parameter from the DetectionInfo struct
func (myDetectionInfo *DetectionInfo) SetProof(val string) {
	myDetectionInfo.Proofvar = val
}

// SetProtocol sets the Protocol parameter from the DetectionInfo struct
func (myDetectionInfo *DetectionInfo) SetProtocol(val string) {
	myDetectionInfo.Protocolvar = val
}

// SetSourceID sets the SourceID parameter from the DetectionInfo struct
func (myDetectionInfo *DetectionInfo) SetSourceID(val string) {
	myDetectionInfo.SourceIDvar = val
}

// SetVulnerabilityID sets the VulnerabilityID parameter from the DetectionInfo struct
func (myDetectionInfo *DetectionInfo) SetVulnerabilityID(val string) {
	myDetectionInfo.VulnerabilityIDvar = val
}
