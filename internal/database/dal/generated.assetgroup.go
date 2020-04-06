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

// AssetGroup defines the struct that implements the AssetGroup interface
type AssetGroup struct {
	CloudSourceIDvar         *string
	GroupIDvar               string
	LastTicketingvar         *time.Time
	OrganizationIDvar        string
	ScannerSourceConfigIDvar *string
	ScannerSourceIDvar       string
}

//**********************************************************
// Struct Methods
//**********************************************************

// MarshalJSON marshals the struct by converting it to a map
func (myAssetGroup AssetGroup) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"CloudSourceID":         myAssetGroup.CloudSourceIDvar,
		"GroupID":               myAssetGroup.GroupIDvar,
		"LastTicketing":         myAssetGroup.LastTicketingvar,
		"OrganizationID":        myAssetGroup.OrganizationIDvar,
		"ScannerSourceConfigID": myAssetGroup.ScannerSourceConfigIDvar,
		"ScannerSourceID":       myAssetGroup.ScannerSourceIDvar,
	})
}

// CloudSourceID returns the CloudSourceID parameter from the AssetGroup struct
func (myAssetGroup *AssetGroup) CloudSourceID() (param *string) {
	return myAssetGroup.CloudSourceIDvar
}

// GroupID returns the GroupID parameter from the AssetGroup struct
func (myAssetGroup *AssetGroup) GroupID() (param string) {
	return myAssetGroup.GroupIDvar
}

// LastTicketing returns the LastTicketing parameter from the AssetGroup struct
func (myAssetGroup *AssetGroup) LastTicketing() (param *time.Time) {
	return myAssetGroup.LastTicketingvar
}

// OrganizationID returns the OrganizationID parameter from the AssetGroup struct
func (myAssetGroup *AssetGroup) OrganizationID() (param string) {
	return myAssetGroup.OrganizationIDvar
}

// ScannerSourceConfigID returns the ScannerSourceConfigID parameter from the AssetGroup struct
func (myAssetGroup *AssetGroup) ScannerSourceConfigID() (param *string) {
	return myAssetGroup.ScannerSourceConfigIDvar
}

// ScannerSourceID returns the ScannerSourceID parameter from the AssetGroup struct
func (myAssetGroup *AssetGroup) ScannerSourceID() (param string) {
	return myAssetGroup.ScannerSourceIDvar
}

// SetCloudSourceID sets the CloudSourceID parameter from the AssetGroup struct
func (myAssetGroup *AssetGroup) SetCloudSourceID(val string) {
	myAssetGroup.CloudSourceIDvar = &val
}

// SetGroupID sets the GroupID parameter from the AssetGroup struct
func (myAssetGroup *AssetGroup) SetGroupID(val string) {
	myAssetGroup.GroupIDvar = val
}

// SetOrganizationID sets the OrganizationID parameter from the AssetGroup struct
func (myAssetGroup *AssetGroup) SetOrganizationID(val string) {
	myAssetGroup.OrganizationIDvar = val
}

// SetScannerSourceConfigID sets the ScannerSourceConfigID parameter from the AssetGroup struct
func (myAssetGroup *AssetGroup) SetScannerSourceConfigID(val string) {
	myAssetGroup.ScannerSourceConfigIDvar = &val
}

// SetScannerSourceID sets the ScannerSourceID parameter from the AssetGroup struct
func (myAssetGroup *AssetGroup) SetScannerSourceID(val string) {
	myAssetGroup.ScannerSourceIDvar = val
}
