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

// LogType defines the struct that implements the LogType interface
type LogType struct {
	IDvar      int
	LogTypevar string
	Namevar    string
}

//**********************************************************
// Struct Methods
//**********************************************************

// MarshalJSON marshals the struct by converting it to a map
func (myLogType LogType) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"ID":      myLogType.IDvar,
		"LogType": myLogType.LogTypevar,
		"Name":    myLogType.Namevar,
	})
}

// ID returns the ID parameter from the LogType struct
func (myLogType *LogType) ID() (param int) {
	return myLogType.IDvar
}

// LogType returns the LogType parameter from the LogType struct
func (myLogType *LogType) LogType() (param string) {
	return myLogType.LogTypevar
}

// Name returns the Name parameter from the LogType struct
func (myLogType *LogType) Name() (param string) {
	return myLogType.Namevar
}

// SetLogType sets the LogType parameter from the LogType struct
func (myLogType *LogType) SetLogType(val string) {
	myLogType.LogTypevar = val
}

// SetName sets the Name parameter from the LogType struct
func (myLogType *LogType) SetName(val string) {
	myLogType.Namevar = val
}
