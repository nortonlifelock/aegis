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

// DeviceGroup defines the struct that implements the DeviceGroup interface
type DeviceGroup struct {
	Descriptionvar      *string
	SourceIdentifiervar int
}

//**********************************************************
// Struct Methods
//**********************************************************

// MarshalJSON marshals the struct by converting it to a map
func (myDeviceGroup DeviceGroup) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Description":      myDeviceGroup.Descriptionvar,
		"SourceIdentifier": myDeviceGroup.SourceIdentifiervar,
	})
}

// Description returns the Description parameter from the DeviceGroup struct
func (myDeviceGroup *DeviceGroup) Description() (param *string) {
	return myDeviceGroup.Descriptionvar
}

// SourceIdentifier returns the SourceIdentifier parameter from the DeviceGroup struct
func (myDeviceGroup *DeviceGroup) SourceIdentifier() (param int) {
	return myDeviceGroup.SourceIdentifiervar
}
