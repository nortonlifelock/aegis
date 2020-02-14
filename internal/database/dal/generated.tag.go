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

// Tag defines the struct that implements the Tag interface
type Tag struct {
	DeviceIDvar string
	IDvar       string
	TagKeyIDvar int
	Valuevar    string
}

//**********************************************************
// Struct Methods
//**********************************************************

// MarshalJSON marshals the struct by converting it to a map
func (myTag Tag) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"DeviceID": myTag.DeviceIDvar,
		"ID":       myTag.IDvar,
		"TagKeyID": myTag.TagKeyIDvar,
		"Value":    myTag.Valuevar,
	})
}

// DeviceID returns the DeviceID parameter from the Tag struct
func (myTag *Tag) DeviceID() (param string) {
	return myTag.DeviceIDvar
}

// ID returns the ID parameter from the Tag struct
func (myTag *Tag) ID() (param string) {
	return myTag.IDvar
}

// TagKeyID returns the TagKeyID parameter from the Tag struct
func (myTag *Tag) TagKeyID() (param int) {
	return myTag.TagKeyIDvar
}

// Value returns the Value parameter from the Tag struct
func (myTag *Tag) Value() (param string) {
	return myTag.Valuevar
}

// SetDeviceID sets the DeviceID parameter from the Tag struct
func (myTag *Tag) SetDeviceID(val string) {
	myTag.DeviceIDvar = val
}

// SetID sets the ID parameter from the Tag struct
func (myTag *Tag) SetID(val string) {
	myTag.IDvar = val
}

// SetValue sets the Value parameter from the Tag struct
func (myTag *Tag) SetValue(val string) {
	myTag.Valuevar = val
}
