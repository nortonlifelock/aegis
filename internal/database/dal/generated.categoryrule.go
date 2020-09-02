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

// CategoryRule defines the struct that implements the CategoryRule interface
type CategoryRule struct {
	Categoryvar              string
	IDvar                    string
	OrganizationIDvar        string
	SourceIDvar              string
	VulnerabilityCategoryvar *string
	VulnerabilityTitlevar    *string
	VulnerabilityTypevar     *string
}

//**********************************************************
// Struct Methods
//**********************************************************

// MarshalJSON marshals the struct by converting it to a map
func (myCategoryRule CategoryRule) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Category":              myCategoryRule.Categoryvar,
		"ID":                    myCategoryRule.IDvar,
		"OrganizationID":        myCategoryRule.OrganizationIDvar,
		"SourceID":              myCategoryRule.SourceIDvar,
		"VulnerabilityCategory": myCategoryRule.VulnerabilityCategoryvar,
		"VulnerabilityTitle":    myCategoryRule.VulnerabilityTitlevar,
		"VulnerabilityType":     myCategoryRule.VulnerabilityTypevar,
	})
}

// Category returns the Category parameter from the CategoryRule struct
func (myCategoryRule *CategoryRule) Category() (param string) {
	return myCategoryRule.Categoryvar
}

// ID returns the ID parameter from the CategoryRule struct
func (myCategoryRule *CategoryRule) ID() (param string) {
	return myCategoryRule.IDvar
}

// OrganizationID returns the OrganizationID parameter from the CategoryRule struct
func (myCategoryRule *CategoryRule) OrganizationID() (param string) {
	return myCategoryRule.OrganizationIDvar
}

// SourceID returns the SourceID parameter from the CategoryRule struct
func (myCategoryRule *CategoryRule) SourceID() (param string) {
	return myCategoryRule.SourceIDvar
}

// VulnerabilityCategory returns the VulnerabilityCategory parameter from the CategoryRule struct
func (myCategoryRule *CategoryRule) VulnerabilityCategory() (param *string) {
	return myCategoryRule.VulnerabilityCategoryvar
}

// VulnerabilityTitle returns the VulnerabilityTitle parameter from the CategoryRule struct
func (myCategoryRule *CategoryRule) VulnerabilityTitle() (param *string) {
	return myCategoryRule.VulnerabilityTitlevar
}

// VulnerabilityType returns the VulnerabilityType parameter from the CategoryRule struct
func (myCategoryRule *CategoryRule) VulnerabilityType() (param *string) {
	return myCategoryRule.VulnerabilityTypevar
}

// SetCategory sets the Category parameter from the CategoryRule struct
func (myCategoryRule *CategoryRule) SetCategory(val string) {
	myCategoryRule.Categoryvar = val
}

// SetID sets the ID parameter from the CategoryRule struct
func (myCategoryRule *CategoryRule) SetID(val string) {
	myCategoryRule.IDvar = val
}

// SetOrganizationID sets the OrganizationID parameter from the CategoryRule struct
func (myCategoryRule *CategoryRule) SetOrganizationID(val string) {
	myCategoryRule.OrganizationIDvar = val
}

// SetSourceID sets the SourceID parameter from the CategoryRule struct
func (myCategoryRule *CategoryRule) SetSourceID(val string) {
	myCategoryRule.SourceIDvar = val
}

// SetVulnerabilityCategory sets the VulnerabilityCategory parameter from the CategoryRule struct
func (myCategoryRule *CategoryRule) SetVulnerabilityCategory(val string) {
	myCategoryRule.VulnerabilityCategoryvar = &val
}

// SetVulnerabilityTitle sets the VulnerabilityTitle parameter from the CategoryRule struct
func (myCategoryRule *CategoryRule) SetVulnerabilityTitle(val string) {
	myCategoryRule.VulnerabilityTitlevar = &val
}

// SetVulnerabilityType sets the VulnerabilityType parameter from the CategoryRule struct
func (myCategoryRule *CategoryRule) SetVulnerabilityType(val string) {
	myCategoryRule.VulnerabilityTypevar = &val
}
