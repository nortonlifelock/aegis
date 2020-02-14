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

// TagMap defines the struct that implements the TagMap interface
type TagMap struct {
	CloudSourceIDvar     string
	CloudTagvar          string
	IDvar                string
	Optionsvar           string
	TicketingSourceIDvar string
	TicketingTagvar      string
}

//**********************************************************
// Struct Methods
//**********************************************************

// MarshalJSON marshals the struct by converting it to a map
func (myTagMap TagMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"CloudSourceID":     myTagMap.CloudSourceIDvar,
		"CloudTag":          myTagMap.CloudTagvar,
		"ID":                myTagMap.IDvar,
		"Options":           myTagMap.Optionsvar,
		"TicketingSourceID": myTagMap.TicketingSourceIDvar,
		"TicketingTag":      myTagMap.TicketingTagvar,
	})
}

// CloudSourceID returns the CloudSourceID parameter from the TagMap struct
func (myTagMap *TagMap) CloudSourceID() (param string) {
	return myTagMap.CloudSourceIDvar
}

// CloudTag returns the CloudTag parameter from the TagMap struct
func (myTagMap *TagMap) CloudTag() (param string) {
	return myTagMap.CloudTagvar
}

// ID returns the ID parameter from the TagMap struct
func (myTagMap *TagMap) ID() (param string) {
	return myTagMap.IDvar
}

// Options returns the Options parameter from the TagMap struct
func (myTagMap *TagMap) Options() (param string) {
	return myTagMap.Optionsvar
}

// TicketingSourceID returns the TicketingSourceID parameter from the TagMap struct
func (myTagMap *TagMap) TicketingSourceID() (param string) {
	return myTagMap.TicketingSourceIDvar
}

// TicketingTag returns the TicketingTag parameter from the TagMap struct
func (myTagMap *TagMap) TicketingTag() (param string) {
	return myTagMap.TicketingTagvar
}

// SetCloudSourceID sets the CloudSourceID parameter from the TagMap struct
func (myTagMap *TagMap) SetCloudSourceID(val string) {
	myTagMap.CloudSourceIDvar = val
}

// SetCloudTag sets the CloudTag parameter from the TagMap struct
func (myTagMap *TagMap) SetCloudTag(val string) {
	myTagMap.CloudTagvar = val
}

// SetID sets the ID parameter from the TagMap struct
func (myTagMap *TagMap) SetID(val string) {
	myTagMap.IDvar = val
}

// SetOptions sets the Options parameter from the TagMap struct
func (myTagMap *TagMap) SetOptions(val string) {
	myTagMap.Optionsvar = val
}

// SetTicketingSourceID sets the TicketingSourceID parameter from the TagMap struct
func (myTagMap *TagMap) SetTicketingSourceID(val string) {
	myTagMap.TicketingSourceIDvar = val
}

// SetTicketingTag sets the TicketingTag parameter from the TagMap struct
func (myTagMap *TagMap) SetTicketingTag(val string) {
	myTagMap.TicketingTagvar = val
}
