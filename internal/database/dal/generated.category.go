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

// Category defines the struct that implements the Category interface
type Category struct {
	Categoryvar         string
	IDvar               string
	ParentCategoryIDvar *string
}

//**********************************************************
// Struct Methods
//**********************************************************

// MarshalJSON marshals the struct by converting it to a map
func (myCategory Category) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Category":         myCategory.Categoryvar,
		"ID":               myCategory.IDvar,
		"ParentCategoryID": myCategory.ParentCategoryIDvar,
	})
}

// Category returns the Category parameter from the Category struct
func (myCategory *Category) Category() (param string) {
	return myCategory.Categoryvar
}

// ID returns the ID parameter from the Category struct
func (myCategory *Category) ID() (param string) {
	return myCategory.IDvar
}

// ParentCategoryID returns the ParentCategoryID parameter from the Category struct
func (myCategory *Category) ParentCategoryID() (param *string) {
	return myCategory.ParentCategoryIDvar
}

// SetCategory sets the Category parameter from the Category struct
func (myCategory *Category) SetCategory(val string) {
	myCategory.Categoryvar = val
}

// SetID sets the ID parameter from the Category struct
func (myCategory *Category) SetID(val string) {
	myCategory.IDvar = val
}

// SetParentCategoryID sets the ParentCategoryID parameter from the Category struct
func (myCategory *Category) SetParentCategoryID(val string) {
	myCategory.ParentCategoryIDvar = &val
}
