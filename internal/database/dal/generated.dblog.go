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

// DBLog defines the struct that implements the DBLog interface
type DBLog struct {
	CreateDatevar   time.Time
	Errorvar        string
	IDvar           int
	JobHistoryIDvar string
	Logvar          string
	TypeIDvar       int
}

//**********************************************************
// Struct Methods
//**********************************************************

// MarshalJSON marshals the struct by converting it to a map
func (myDBLog DBLog) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"CreateDate":   myDBLog.CreateDatevar,
		"Error":        myDBLog.Errorvar,
		"ID":           myDBLog.IDvar,
		"JobHistoryID": myDBLog.JobHistoryIDvar,
		"Log":          myDBLog.Logvar,
		"TypeID":       myDBLog.TypeIDvar,
	})
}

// CreateDate returns the CreateDate parameter from the DBLog struct
func (myDBLog *DBLog) CreateDate() (param time.Time) {
	return myDBLog.CreateDatevar
}

// Error returns the Error parameter from the DBLog struct
func (myDBLog *DBLog) Error() (param string) {
	return myDBLog.Errorvar
}

// ID returns the ID parameter from the DBLog struct
func (myDBLog *DBLog) ID() (param int) {
	return myDBLog.IDvar
}

// JobHistoryID returns the JobHistoryID parameter from the DBLog struct
func (myDBLog *DBLog) JobHistoryID() (param string) {
	return myDBLog.JobHistoryIDvar
}

// Log returns the Log parameter from the DBLog struct
func (myDBLog *DBLog) Log() (param string) {
	return myDBLog.Logvar
}

// TypeID returns the TypeID parameter from the DBLog struct
func (myDBLog *DBLog) TypeID() (param int) {
	return myDBLog.TypeIDvar
}

// SetError sets the Error parameter from the DBLog struct
func (myDBLog *DBLog) SetError(val string) {
	myDBLog.Errorvar = val
}

// SetJobHistoryID sets the JobHistoryID parameter from the DBLog struct
func (myDBLog *DBLog) SetJobHistoryID(val string) {
	myDBLog.JobHistoryIDvar = val
}

// SetLog sets the Log parameter from the DBLog struct
func (myDBLog *DBLog) SetLog(val string) {
	myDBLog.Logvar = val
}
