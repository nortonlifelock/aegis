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
	Assigneevar           *string
	AssignmentGroupvar    *string
	OrganizationIDvar     string
	Priorityvar           int
	TagKeyIDvar           *int
	TagKeyValuevar        *string
	VulnTitleRegexvar     *string
	VulnTitleSubstringvar *string
}

//**********************************************************
// Struct Methods
//**********************************************************

// MarshalJSON marshals the struct by converting it to a map
func (myAssignmentRules AssignmentRules) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Assignee":           myAssignmentRules.Assigneevar,
		"AssignmentGroup":    myAssignmentRules.AssignmentGroupvar,
		"OrganizationID":     myAssignmentRules.OrganizationIDvar,
		"Priority":           myAssignmentRules.Priorityvar,
		"TagKeyID":           myAssignmentRules.TagKeyIDvar,
		"TagKeyValue":        myAssignmentRules.TagKeyValuevar,
		"VulnTitleRegex":     myAssignmentRules.VulnTitleRegexvar,
		"VulnTitleSubstring": myAssignmentRules.VulnTitleSubstringvar,
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

// TagKeyValue returns the TagKeyValue parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) TagKeyValue() (param *string) {
	return myAssignmentRules.TagKeyValuevar
}

// VulnTitleRegex returns the VulnTitleRegex parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) VulnTitleRegex() (param *string) {
	return myAssignmentRules.VulnTitleRegexvar
}

// VulnTitleSubstring returns the VulnTitleSubstring parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) VulnTitleSubstring() (param *string) {
	return myAssignmentRules.VulnTitleSubstringvar
}

// SetAssignee sets the Assignee parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) SetAssignee(val string) {
	myAssignmentRules.Assigneevar = &val
}

// SetAssignmentGroup sets the AssignmentGroup parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) SetAssignmentGroup(val string) {
	myAssignmentRules.AssignmentGroupvar = &val
}

// SetOrganizationID sets the OrganizationID parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) SetOrganizationID(val string) {
	myAssignmentRules.OrganizationIDvar = val
}

// SetTagKeyValue sets the TagKeyValue parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) SetTagKeyValue(val string) {
	myAssignmentRules.TagKeyValuevar = &val
}

// SetVulnTitleRegex sets the VulnTitleRegex parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) SetVulnTitleRegex(val string) {
	myAssignmentRules.VulnTitleRegexvar = &val
}

// SetVulnTitleSubstring sets the VulnTitleSubstring parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) SetVulnTitleSubstring(val string) {
	myAssignmentRules.VulnTitleSubstringvar = &val
}
