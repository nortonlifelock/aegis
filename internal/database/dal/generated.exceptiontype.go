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

// ExceptionType defines the struct that implements the ExceptionType interface
type ExceptionType struct {
	IDvar   int
	Namevar string
	Typevar string
}

//**********************************************************
// Struct Methods
//**********************************************************

// MarshalJSON marshals the struct by converting it to a map
func (myExceptionType ExceptionType) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"ID":   myExceptionType.IDvar,
		"Name": myExceptionType.Namevar,
		"Type": myExceptionType.Typevar,
	})
}

// ID returns the ID parameter from the ExceptionType struct
func (myExceptionType *ExceptionType) ID() (param int) {
	return myExceptionType.IDvar
}

// Name returns the Name parameter from the ExceptionType struct
func (myExceptionType *ExceptionType) Name() (param string) {
	return myExceptionType.Namevar
}

// Type returns the Type parameter from the ExceptionType struct
func (myExceptionType *ExceptionType) Type() (param string) {
	return myExceptionType.Typevar
}

// SetName sets the Name parameter from the ExceptionType struct
func (myExceptionType *ExceptionType) SetName(val string) {
	myExceptionType.Namevar = val
}

// SetType sets the Type parameter from the ExceptionType struct
func (myExceptionType *ExceptionType) SetType(val string) {
	myExceptionType.Typevar = val
}
