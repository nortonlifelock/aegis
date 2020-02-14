package dal

//**********************************************************
// GENERATED CODE - DO NOT CHANGE
// This file is generated using scaffolding. Any changes to
// this file will be overwritten on the next build
//**********************************************************

import (
	"encoding/json"
	"github.com/nortonlifelock/aegis/internal/domain"
)

//**********************************************************
// Struct Declaration
//**********************************************************

// Permission defines the struct that implements the Permission interface
type Permission struct {
	Adminvar               bool
	Managervar             bool
	OrgIDvar               string
	ParentOrgPermissionvar domain.Permission
	Readervar              bool
	Reportervar            bool
	UserIDvar              string
}

//**********************************************************
// Struct Methods
//**********************************************************

// MarshalJSON marshals the struct by converting it to a map
func (myPermission Permission) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Admin":               myPermission.Adminvar,
		"Manager":             myPermission.Managervar,
		"OrgID":               myPermission.OrgIDvar,
		"ParentOrgPermission": myPermission.ParentOrgPermissionvar,
		"Reader":              myPermission.Readervar,
		"Reporter":            myPermission.Reportervar,
		"UserID":              myPermission.UserIDvar,
	})
}

// Admin returns the Admin parameter from the Permission struct
func (myPermission *Permission) Admin() (param bool) {
	return myPermission.Adminvar
}

// Manager returns the Manager parameter from the Permission struct
func (myPermission *Permission) Manager() (param bool) {
	return myPermission.Managervar
}

// OrgID returns the OrgID parameter from the Permission struct
func (myPermission *Permission) OrgID() (param string) {
	return myPermission.OrgIDvar
}

// ParentOrgPermission returns the ParentOrgPermission parameter from the Permission struct
func (myPermission *Permission) ParentOrgPermission() (param domain.Permission) {
	return myPermission.ParentOrgPermissionvar
}

// Reader returns the Reader parameter from the Permission struct
func (myPermission *Permission) Reader() (param bool) {
	return myPermission.Readervar
}

// Reporter returns the Reporter parameter from the Permission struct
func (myPermission *Permission) Reporter() (param bool) {
	return myPermission.Reportervar
}

// UserID returns the UserID parameter from the Permission struct
func (myPermission *Permission) UserID() (param string) {
	return myPermission.UserIDvar
}

// SetOrgID sets the OrgID parameter from the Permission struct
func (myPermission *Permission) SetOrgID(val string) {
	myPermission.OrgIDvar = val
}

// SetUserID sets the UserID parameter from the Permission struct
func (myPermission *Permission) SetUserID(val string) {
	myPermission.UserIDvar = val
}
