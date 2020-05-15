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

// DeviceInfo defines the struct that implements the DeviceInfo interface
type DeviceInfo struct {
	GroupIDvar         *string
	HostNamevar        string
	IDvar              string
	IPvar              string
	InstanceIDvar      *string
	MACvar             string
	OSvar              string
	Regionvar          *string
	ScannerSourceIDvar *string
	SourceIDvar        *string
	Statevar           *string
	TrackingMethodvar  *string
}

//**********************************************************
// Struct Methods
//**********************************************************

// MarshalJSON marshals the struct by converting it to a map
func (myDeviceInfo DeviceInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"GroupID":         myDeviceInfo.GroupIDvar,
		"HostName":        myDeviceInfo.HostNamevar,
		"ID":              myDeviceInfo.IDvar,
		"IP":              myDeviceInfo.IPvar,
		"InstanceID":      myDeviceInfo.InstanceIDvar,
		"MAC":             myDeviceInfo.MACvar,
		"OS":              myDeviceInfo.OSvar,
		"Region":          myDeviceInfo.Regionvar,
		"ScannerSourceID": myDeviceInfo.ScannerSourceIDvar,
		"SourceID":        myDeviceInfo.SourceIDvar,
		"State":           myDeviceInfo.Statevar,
		"TrackingMethod":  myDeviceInfo.TrackingMethodvar,
	})
}

// GroupID returns the GroupID parameter from the DeviceInfo struct
func (myDeviceInfo *DeviceInfo) GroupID() (param *string) {
	return myDeviceInfo.GroupIDvar
}

// HostName returns the HostName parameter from the DeviceInfo struct
func (myDeviceInfo *DeviceInfo) HostName() (param string) {
	return myDeviceInfo.HostNamevar
}

// ID returns the ID parameter from the DeviceInfo struct
func (myDeviceInfo *DeviceInfo) ID() (param string) {
	return myDeviceInfo.IDvar
}

// IP returns the IP parameter from the DeviceInfo struct
func (myDeviceInfo *DeviceInfo) IP() (param string) {
	return myDeviceInfo.IPvar
}

// InstanceID returns the InstanceID parameter from the DeviceInfo struct
func (myDeviceInfo *DeviceInfo) InstanceID() (param *string) {
	return myDeviceInfo.InstanceIDvar
}

// MAC returns the MAC parameter from the DeviceInfo struct
func (myDeviceInfo *DeviceInfo) MAC() (param string) {
	return myDeviceInfo.MACvar
}

// OS returns the OS parameter from the DeviceInfo struct
func (myDeviceInfo *DeviceInfo) OS() (param string) {
	return myDeviceInfo.OSvar
}

// Region returns the Region parameter from the DeviceInfo struct
func (myDeviceInfo *DeviceInfo) Region() (param *string) {
	return myDeviceInfo.Regionvar
}

// ScannerSourceID returns the ScannerSourceID parameter from the DeviceInfo struct
func (myDeviceInfo *DeviceInfo) ScannerSourceID() (param *string) {
	return myDeviceInfo.ScannerSourceIDvar
}

// SourceID returns the SourceID parameter from the DeviceInfo struct
func (myDeviceInfo *DeviceInfo) SourceID() (param *string) {
	return myDeviceInfo.SourceIDvar
}

// State returns the State parameter from the DeviceInfo struct
func (myDeviceInfo *DeviceInfo) State() (param *string) {
	return myDeviceInfo.Statevar
}

// TrackingMethod returns the TrackingMethod parameter from the DeviceInfo struct
func (myDeviceInfo *DeviceInfo) TrackingMethod() (param *string) {
	return myDeviceInfo.TrackingMethodvar
}

// SetGroupID sets the GroupID parameter from the DeviceInfo struct
func (myDeviceInfo *DeviceInfo) SetGroupID(val string) {
	myDeviceInfo.GroupIDvar = &val
}

// SetHostName sets the HostName parameter from the DeviceInfo struct
func (myDeviceInfo *DeviceInfo) SetHostName(val string) {
	myDeviceInfo.HostNamevar = val
}

// SetID sets the ID parameter from the DeviceInfo struct
func (myDeviceInfo *DeviceInfo) SetID(val string) {
	myDeviceInfo.IDvar = val
}

// SetIP sets the IP parameter from the DeviceInfo struct
func (myDeviceInfo *DeviceInfo) SetIP(val string) {
	myDeviceInfo.IPvar = val
}

// SetInstanceID sets the InstanceID parameter from the DeviceInfo struct
func (myDeviceInfo *DeviceInfo) SetInstanceID(val string) {
	myDeviceInfo.InstanceIDvar = &val
}

// SetMAC sets the MAC parameter from the DeviceInfo struct
func (myDeviceInfo *DeviceInfo) SetMAC(val string) {
	myDeviceInfo.MACvar = val
}

// SetOS sets the OS parameter from the DeviceInfo struct
func (myDeviceInfo *DeviceInfo) SetOS(val string) {
	myDeviceInfo.OSvar = val
}

// SetRegion sets the Region parameter from the DeviceInfo struct
func (myDeviceInfo *DeviceInfo) SetRegion(val string) {
	myDeviceInfo.Regionvar = &val
}

// SetScannerSourceID sets the ScannerSourceID parameter from the DeviceInfo struct
func (myDeviceInfo *DeviceInfo) SetScannerSourceID(val string) {
	myDeviceInfo.ScannerSourceIDvar = &val
}

// SetSourceID sets the SourceID parameter from the DeviceInfo struct
func (myDeviceInfo *DeviceInfo) SetSourceID(val string) {
	myDeviceInfo.SourceIDvar = &val
}

// SetState sets the State parameter from the DeviceInfo struct
func (myDeviceInfo *DeviceInfo) SetState(val string) {
	myDeviceInfo.Statevar = &val
}

// SetTrackingMethod sets the TrackingMethod parameter from the DeviceInfo struct
func (myDeviceInfo *DeviceInfo) SetTrackingMethod(val string) {
	myDeviceInfo.TrackingMethodvar = &val
}
