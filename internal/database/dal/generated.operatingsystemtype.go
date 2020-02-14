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

// OperatingSystemType defines the struct that implements the OperatingSystemType interface
type OperatingSystemType struct {
	IDvar       int
	Matchvar    string
	Priorityvar int
	Typevar     string
}

//**********************************************************
// Struct Methods
//**********************************************************

// MarshalJSON marshals the struct by converting it to a map
func (myOperatingSystemType OperatingSystemType) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"ID":       myOperatingSystemType.IDvar,
		"Match":    myOperatingSystemType.Matchvar,
		"Priority": myOperatingSystemType.Priorityvar,
		"Type":     myOperatingSystemType.Typevar,
	})
}

// ID returns the ID parameter from the OperatingSystemType struct
func (myOperatingSystemType *OperatingSystemType) ID() (param int) {
	return myOperatingSystemType.IDvar
}

// Match returns the Match parameter from the OperatingSystemType struct
func (myOperatingSystemType *OperatingSystemType) Match() (param string) {
	return myOperatingSystemType.Matchvar
}

// Priority returns the Priority parameter from the OperatingSystemType struct
func (myOperatingSystemType *OperatingSystemType) Priority() (param int) {
	return myOperatingSystemType.Priorityvar
}

// Type returns the Type parameter from the OperatingSystemType struct
func (myOperatingSystemType *OperatingSystemType) Type() (param string) {
	return myOperatingSystemType.Typevar
}

// SetMatch sets the Match parameter from the OperatingSystemType struct
func (myOperatingSystemType *OperatingSystemType) SetMatch(val string) {
	myOperatingSystemType.Matchvar = val
}

// SetType sets the Type parameter from the OperatingSystemType struct
func (myOperatingSystemType *OperatingSystemType) SetType(val string) {
	myOperatingSystemType.Typevar = val
}
