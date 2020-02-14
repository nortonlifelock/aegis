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

// KeyValue defines the struct that implements the KeyValue interface
type KeyValue struct {
	Keyvar   string
	Valuevar string
}

//**********************************************************
// Struct Methods
//**********************************************************

// MarshalJSON marshals the struct by converting it to a map
func (myKeyValue KeyValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Key":   myKeyValue.Keyvar,
		"Value": myKeyValue.Valuevar,
	})
}

// Key returns the Key parameter from the KeyValue struct
func (myKeyValue *KeyValue) Key() (param string) {
	return myKeyValue.Keyvar
}

// Value returns the Value parameter from the KeyValue struct
func (myKeyValue *KeyValue) Value() (param string) {
	return myKeyValue.Valuevar
}

// SetKey sets the Key parameter from the KeyValue struct
func (myKeyValue *KeyValue) SetKey(val string) {
	myKeyValue.Keyvar = val
}

// SetValue sets the Value parameter from the KeyValue struct
func (myKeyValue *KeyValue) SetValue(val string) {
	myKeyValue.Valuevar = val
}
