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

// Source defines the struct that implements the Source interface
type Source struct {
	DBCreatedDatevar time.Time
	DBUpdatedDatevar *time.Time
	IDvar            string
	Sourcevar        string
	SourceTypeIDvar  int
}

//**********************************************************
// Struct Methods
//**********************************************************

// MarshalJSON marshals the struct by converting it to a map
func (mySource Source) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"DBCreatedDate": mySource.DBCreatedDatevar,
		"DBUpdatedDate": mySource.DBUpdatedDatevar,
		"ID":            mySource.IDvar,
		"Source":        mySource.Sourcevar,
		"SourceTypeID":  mySource.SourceTypeIDvar,
	})
}

// DBCreatedDate returns the DBCreatedDate parameter from the Source struct
func (mySource *Source) DBCreatedDate() (param time.Time) {
	return mySource.DBCreatedDatevar
}

// DBUpdatedDate returns the DBUpdatedDate parameter from the Source struct
func (mySource *Source) DBUpdatedDate() (param *time.Time) {
	return mySource.DBUpdatedDatevar
}

// ID returns the ID parameter from the Source struct
func (mySource *Source) ID() (param string) {
	return mySource.IDvar
}

// Source returns the Source parameter from the Source struct
func (mySource *Source) Source() (param string) {
	return mySource.Sourcevar
}

// SourceTypeID returns the SourceTypeID parameter from the Source struct
func (mySource *Source) SourceTypeID() (param int) {
	return mySource.SourceTypeIDvar
}

// SetID sets the ID parameter from the Source struct
func (mySource *Source) SetID(val string) {
	mySource.IDvar = val
}

// SetSource sets the Source parameter from the Source struct
func (mySource *Source) SetSource(val string) {
	mySource.Sourcevar = val
}
