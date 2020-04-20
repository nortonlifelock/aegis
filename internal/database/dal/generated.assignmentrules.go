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

// AssignmentRules defines the struct that implements the AssignmentRules interface
type AssignmentRules struct {
	Assigneevar        *string
	AssignmentGroupvar *string
	GroupIDvar         *string
	HostnameRegexvar   *string
	OrganizationIDvar  string
	Priorityvar        int
	TagKeyIDvar        *int
	TagKeyRegexvar     *string
	VulnTitleRegexvar  *string
}

//**********************************************************
// Struct Methods
//**********************************************************

// MarshalJSON marshals the struct by converting it to a map
func (myAssignmentRules AssignmentRules) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Assignee":        myAssignmentRules.Assigneevar,
		"AssignmentGroup": myAssignmentRules.AssignmentGroupvar,
		"GroupID":         myAssignmentRules.GroupIDvar,
		"HostnameRegex":   myAssignmentRules.HostnameRegexvar,
		"OrganizationID":  myAssignmentRules.OrganizationIDvar,
		"Priority":        myAssignmentRules.Priorityvar,
		"TagKeyID":        myAssignmentRules.TagKeyIDvar,
		"TagKeyRegex":     myAssignmentRules.TagKeyRegexvar,
		"VulnTitleRegex":  myAssignmentRules.VulnTitleRegexvar,
	})
}

// Assignee returns the Assignee parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) Assignee() (param *string) {
	return myAssignmentRules.Assigneevar
}

// AssignmentGroup returns the AssignmentGroup parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) AssignmentGroup() (param *string) {
	return myAssignmentRules.AssignmentGroupvar
}

// GroupID returns the GroupID parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) GroupID() (param *string) {
	return myAssignmentRules.GroupIDvar
}

// HostnameRegex returns the HostnameRegex parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) HostnameRegex() (param *string) {
	return myAssignmentRules.HostnameRegexvar
}

// OrganizationID returns the OrganizationID parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) OrganizationID() (param string) {
	return myAssignmentRules.OrganizationIDvar
}

// Priority returns the Priority parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) Priority() (param int) {
	return myAssignmentRules.Priorityvar
}

// TagKeyID returns the TagKeyID parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) TagKeyID() (param *int) {
	return myAssignmentRules.TagKeyIDvar
}

// TagKeyRegex returns the TagKeyRegex parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) TagKeyRegex() (param *string) {
	return myAssignmentRules.TagKeyRegexvar
}

// VulnTitleRegex returns the VulnTitleRegex parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) VulnTitleRegex() (param *string) {
	return myAssignmentRules.VulnTitleRegexvar
}

// SetAssignee sets the Assignee parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) SetAssignee(val string) {
	myAssignmentRules.Assigneevar = &val
}

// SetAssignmentGroup sets the AssignmentGroup parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) SetAssignmentGroup(val string) {
	myAssignmentRules.AssignmentGroupvar = &val
}

// SetGroupID sets the GroupID parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) SetGroupID(val string) {
	myAssignmentRules.GroupIDvar = &val
}

// SetHostnameRegex sets the HostnameRegex parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) SetHostnameRegex(val string) {
	myAssignmentRules.HostnameRegexvar = &val
}

// SetOrganizationID sets the OrganizationID parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) SetOrganizationID(val string) {
	myAssignmentRules.OrganizationIDvar = val
}

// SetTagKeyRegex sets the TagKeyRegex parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) SetTagKeyRegex(val string) {
	myAssignmentRules.TagKeyRegexvar = &val
}

// SetVulnTitleRegex sets the VulnTitleRegex parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) SetVulnTitleRegex(val string) {
	myAssignmentRules.VulnTitleRegexvar = &val
}
