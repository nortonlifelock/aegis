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

// TagKey defines the struct that implements the TagKey interface
type TagKey struct {
	IDvar       string
	KeyValuevar string
}

//**********************************************************
// Struct Methods
//**********************************************************

// MarshalJSON marshals the struct by converting it to a map
func (myTagKey TagKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"ID":       myTagKey.IDvar,
		"KeyValue": myTagKey.KeyValuevar,
	})
}

// ID returns the ID parameter from the TagKey struct
func (myTagKey *TagKey) ID() (param string) {
	return myTagKey.IDvar
}

// KeyValue returns the KeyValue parameter from the TagKey struct
func (myTagKey *TagKey) KeyValue() (param string) {
	return myTagKey.KeyValuevar
}

// SetID sets the ID parameter from the TagKey struct
func (myTagKey *TagKey) SetID(val string) {
	myTagKey.IDvar = val
}

// SetKeyValue sets the KeyValue parameter from the TagKey struct
func (myTagKey *TagKey) SetKeyValue(val string) {
	myTagKey.KeyValuevar = val
}
