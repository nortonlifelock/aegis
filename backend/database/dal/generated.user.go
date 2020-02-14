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

// User defines the struct that implements the User interface
type User struct {
	Emailvar      string
	FirstNamevar  string
	IDvar         string
	IsDisabledvar bool
	LastNamevar   string
	Usernamevar   *string
}

//**********************************************************
// Struct Methods
//**********************************************************

// MarshalJSON marshals the struct by converting it to a map
func (myUser User) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Email":      myUser.Emailvar,
		"FirstName":  myUser.FirstNamevar,
		"ID":         myUser.IDvar,
		"IsDisabled": myUser.IsDisabledvar,
		"LastName":   myUser.LastNamevar,
		"Username":   myUser.Usernamevar,
	})
}

// Email returns the Email parameter from the User struct
func (myUser *User) Email() (param string) {
	return myUser.Emailvar
}

// FirstName returns the FirstName parameter from the User struct
func (myUser *User) FirstName() (param string) {
	return myUser.FirstNamevar
}

// ID returns the ID parameter from the User struct
func (myUser *User) ID() (param string) {
	return myUser.IDvar
}

// IsDisabled returns the IsDisabled parameter from the User struct
func (myUser *User) IsDisabled() (param bool) {
	return myUser.IsDisabledvar
}

// LastName returns the LastName parameter from the User struct
func (myUser *User) LastName() (param string) {
	return myUser.LastNamevar
}

// Username returns the Username parameter from the User struct
func (myUser *User) Username() (param *string) {
	return myUser.Usernamevar
}

// SetEmail sets the Email parameter from the User struct
func (myUser *User) SetEmail(val string) {
	myUser.Emailvar = val
}

// SetID sets the ID parameter from the User struct
func (myUser *User) SetID(val string) {
	myUser.IDvar = val
}
