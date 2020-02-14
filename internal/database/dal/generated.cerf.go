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

// CERF defines the struct that implements the CERF interface
type CERF struct {
	CERFormvar string
}

//**********************************************************
// Struct Methods
//**********************************************************

// MarshalJSON marshals the struct by converting it to a map
func (myCERF CERF) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"CERForm": myCERF.CERFormvar,
	})
}

// CERForm returns the CERForm parameter from the CERF struct
func (myCERF *CERF) CERForm() (param string) {
	return myCERF.CERFormvar
}

// SetCERForm sets the CERForm parameter from the CERF struct
func (myCERF *CERF) SetCERForm(val string) {
	myCERF.CERFormvar = val
}
