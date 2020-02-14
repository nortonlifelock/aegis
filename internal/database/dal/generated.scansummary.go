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

// ScanSummary defines the struct that implements the ScanSummary interface
type ScanSummary struct {
	CreatedDatevar      time.Time
	OrgIDvar            string
	ParentJobIDvar      string
	ScanClosePayloadvar string
	ScanStatusvar       string
	Sourcevar           string
	SourceIDvar         string
	SourceKeyvar        *string
	TemplateIDvar       *string
	UpdatedDatevar      *time.Time
}

//**********************************************************
// Struct Methods
//**********************************************************

// MarshalJSON marshals the struct by converting it to a map
func (myScanSummary ScanSummary) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"CreatedDate":      myScanSummary.CreatedDatevar,
		"OrgID":            myScanSummary.OrgIDvar,
		"ParentJobID":      myScanSummary.ParentJobIDvar,
		"ScanClosePayload": myScanSummary.ScanClosePayloadvar,
		"ScanStatus":       myScanSummary.ScanStatusvar,
		"Source":           myScanSummary.Sourcevar,
		"SourceID":         myScanSummary.SourceIDvar,
		"SourceKey":        myScanSummary.SourceKeyvar,
		"TemplateID":       myScanSummary.TemplateIDvar,
		"UpdatedDate":      myScanSummary.UpdatedDatevar,
	})
}

// CreatedDate returns the CreatedDate parameter from the ScanSummary struct
func (myScanSummary *ScanSummary) CreatedDate() (param time.Time) {
	return myScanSummary.CreatedDatevar
}

// OrgID returns the OrgID parameter from the ScanSummary struct
func (myScanSummary *ScanSummary) OrgID() (param string) {
	return myScanSummary.OrgIDvar
}

// ParentJobID returns the ParentJobID parameter from the ScanSummary struct
func (myScanSummary *ScanSummary) ParentJobID() (param string) {
	return myScanSummary.ParentJobIDvar
}

// ScanClosePayload returns the ScanClosePayload parameter from the ScanSummary struct
func (myScanSummary *ScanSummary) ScanClosePayload() (param string) {
	return myScanSummary.ScanClosePayloadvar
}

// ScanStatus returns the ScanStatus parameter from the ScanSummary struct
func (myScanSummary *ScanSummary) ScanStatus() (param string) {
	return myScanSummary.ScanStatusvar
}

// Source returns the Source parameter from the ScanSummary struct
func (myScanSummary *ScanSummary) Source() (param string) {
	return myScanSummary.Sourcevar
}

// SourceID returns the SourceID parameter from the ScanSummary struct
func (myScanSummary *ScanSummary) SourceID() (param string) {
	return myScanSummary.SourceIDvar
}

// SourceKey returns the SourceKey parameter from the ScanSummary struct
func (myScanSummary *ScanSummary) SourceKey() (param *string) {
	return myScanSummary.SourceKeyvar
}

// TemplateID returns the TemplateID parameter from the ScanSummary struct
func (myScanSummary *ScanSummary) TemplateID() (param *string) {
	return myScanSummary.TemplateIDvar
}

// UpdatedDate returns the UpdatedDate parameter from the ScanSummary struct
func (myScanSummary *ScanSummary) UpdatedDate() (param *time.Time) {
	return myScanSummary.UpdatedDatevar
}

// SetOrgID sets the OrgID parameter from the ScanSummary struct
func (myScanSummary *ScanSummary) SetOrgID(val string) {
	myScanSummary.OrgIDvar = val
}

// SetParentJobID sets the ParentJobID parameter from the ScanSummary struct
func (myScanSummary *ScanSummary) SetParentJobID(val string) {
	myScanSummary.ParentJobIDvar = val
}

// SetScanStatus sets the ScanStatus parameter from the ScanSummary struct
func (myScanSummary *ScanSummary) SetScanStatus(val string) {
	myScanSummary.ScanStatusvar = val
}

// SetSource sets the Source parameter from the ScanSummary struct
func (myScanSummary *ScanSummary) SetSource(val string) {
	myScanSummary.Sourcevar = val
}

// SetSourceID sets the SourceID parameter from the ScanSummary struct
func (myScanSummary *ScanSummary) SetSourceID(val string) {
	myScanSummary.SourceIDvar = val
}

// SetSourceKey sets the SourceKey parameter from the ScanSummary struct
func (myScanSummary *ScanSummary) SetSourceKey(val string) {
	myScanSummary.SourceKeyvar = &val
}

// SetTemplateID sets the TemplateID parameter from the ScanSummary struct
func (myScanSummary *ScanSummary) SetTemplateID(val string) {
	myScanSummary.TemplateIDvar = &val
}
