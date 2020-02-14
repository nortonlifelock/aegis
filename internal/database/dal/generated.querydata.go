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

// QueryData defines the struct that implements the QueryData interface
type QueryData struct {
	Lengthvar int
}

//**********************************************************
// Struct Methods
//**********************************************************

// MarshalJSON marshals the struct by converting it to a map
func (myQueryData QueryData) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Length": myQueryData.Lengthvar,
	})
}

// Length returns the Length parameter from the QueryData struct
func (myQueryData *QueryData) Length() (param int) {
	return myQueryData.Lengthvar
}
