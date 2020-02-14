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

// Session defines the struct that implements the Session interface
type Session struct {
	IsDisabledvar bool
	OrgIDvar      string
	SessionKeyvar string
	UserIDvar     string
}

//**********************************************************
// Struct Methods
//**********************************************************

// MarshalJSON marshals the struct by converting it to a map
func (mySession Session) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"IsDisabled": mySession.IsDisabledvar,
		"OrgID":      mySession.OrgIDvar,
		"SessionKey": mySession.SessionKeyvar,
		"UserID":     mySession.UserIDvar,
	})
}

// IsDisabled returns the IsDisabled parameter from the Session struct
func (mySession *Session) IsDisabled() (param bool) {
	return mySession.IsDisabledvar
}

// OrgID returns the OrgID parameter from the Session struct
func (mySession *Session) OrgID() (param string) {
	return mySession.OrgIDvar
}

// SessionKey returns the SessionKey parameter from the Session struct
func (mySession *Session) SessionKey() (param string) {
	return mySession.SessionKeyvar
}

// UserID returns the UserID parameter from the Session struct
func (mySession *Session) UserID() (param string) {
	return mySession.UserIDvar
}

// SetOrgID sets the OrgID parameter from the Session struct
func (mySession *Session) SetOrgID(val string) {
	mySession.OrgIDvar = val
}

// SetUserID sets the UserID parameter from the Session struct
func (mySession *Session) SetUserID(val string) {
	mySession.UserIDvar = val
}
