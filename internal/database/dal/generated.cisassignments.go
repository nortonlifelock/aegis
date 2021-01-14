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

// CISAssignments defines the struct that implements the CISAssignments interface
type CISAssignments struct {
	AssignmentGroupvar string
	BundleIDvar        *string
	CloudAccountIDvar  *string
	OrganizationIDvar  string
	RuleIDvar          *string
	RuleRegexvar       *string
}

//**********************************************************
// Struct Methods
//**********************************************************

// MarshalJSON marshals the struct by converting it to a map
func (myCISAssignments CISAssignments) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"AssignmentGroup": myCISAssignments.AssignmentGroupvar,
		"BundleID":        myCISAssignments.BundleIDvar,
		"CloudAccountID":  myCISAssignments.CloudAccountIDvar,
		"OrganizationID":  myCISAssignments.OrganizationIDvar,
		"RuleID":          myCISAssignments.RuleIDvar,
		"RuleRegex":       myCISAssignments.RuleRegexvar,
	})
}

// AssignmentGroup returns the AssignmentGroup parameter from the CISAssignments struct
func (myCISAssignments *CISAssignments) AssignmentGroup() (param string) {
	return myCISAssignments.AssignmentGroupvar
}

// BundleID returns the BundleID parameter from the CISAssignments struct
func (myCISAssignments *CISAssignments) BundleID() (param *string) {
	return myCISAssignments.BundleIDvar
}

// CloudAccountID returns the CloudAccountID parameter from the CISAssignments struct
func (myCISAssignments *CISAssignments) CloudAccountID() (param *string) {
	return myCISAssignments.CloudAccountIDvar
}

// OrganizationID returns the OrganizationID parameter from the CISAssignments struct
func (myCISAssignments *CISAssignments) OrganizationID() (param string) {
	return myCISAssignments.OrganizationIDvar
}

// RuleID returns the RuleID parameter from the CISAssignments struct
func (myCISAssignments *CISAssignments) RuleID() (param *string) {
	return myCISAssignments.RuleIDvar
}

// RuleRegex returns the RuleRegex parameter from the CISAssignments struct
func (myCISAssignments *CISAssignments) RuleRegex() (param *string) {
	return myCISAssignments.RuleRegexvar
}

// SetAssignmentGroup sets the AssignmentGroup parameter from the CISAssignments struct
func (myCISAssignments *CISAssignments) SetAssignmentGroup(val string) {
	myCISAssignments.AssignmentGroupvar = val
}

// SetBundleID sets the BundleID parameter from the CISAssignments struct
func (myCISAssignments *CISAssignments) SetBundleID(val string) {
	myCISAssignments.BundleIDvar = &val
}

// SetCloudAccountID sets the CloudAccountID parameter from the CISAssignments struct
func (myCISAssignments *CISAssignments) SetCloudAccountID(val string) {
	myCISAssignments.CloudAccountIDvar = &val
}

// SetOrganizationID sets the OrganizationID parameter from the CISAssignments struct
func (myCISAssignments *CISAssignments) SetOrganizationID(val string) {
	myCISAssignments.OrganizationIDvar = val
}

// SetRuleID sets the RuleID parameter from the CISAssignments struct
func (myCISAssignments *CISAssignments) SetRuleID(val string) {
	myCISAssignments.RuleIDvar = &val
}

// SetRuleRegex sets the RuleRegex parameter from the CISAssignments struct
func (myCISAssignments *CISAssignments) SetRuleRegex(val string) {
	myCISAssignments.RuleRegexvar = &val
}
