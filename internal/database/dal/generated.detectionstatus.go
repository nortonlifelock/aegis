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

// DetectionStatus defines the struct that implements the DetectionStatus interface
type DetectionStatus struct {
	IDvar     int
	Namevar   string
	Statusvar string
}

//**********************************************************
// Struct Methods
//**********************************************************

// MarshalJSON marshals the struct by converting it to a map
func (myDetectionStatus DetectionStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"ID":     myDetectionStatus.IDvar,
		"Name":   myDetectionStatus.Namevar,
		"Status": myDetectionStatus.Statusvar,
	})
}

// ID returns the ID parameter from the DetectionStatus struct
func (myDetectionStatus *DetectionStatus) ID() (param int) {
	return myDetectionStatus.IDvar
}

// Name returns the Name parameter from the DetectionStatus struct
func (myDetectionStatus *DetectionStatus) Name() (param string) {
	return myDetectionStatus.Namevar
}

// Status returns the Status parameter from the DetectionStatus struct
func (myDetectionStatus *DetectionStatus) Status() (param string) {
	return myDetectionStatus.Statusvar
}

// SetName sets the Name parameter from the DetectionStatus struct
func (myDetectionStatus *DetectionStatus) SetName(val string) {
	myDetectionStatus.Namevar = val
}

// SetStatus sets the Status parameter from the DetectionStatus struct
func (myDetectionStatus *DetectionStatus) SetStatus(val string) {
	myDetectionStatus.Statusvar = val
}
