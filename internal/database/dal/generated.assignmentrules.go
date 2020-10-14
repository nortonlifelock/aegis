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
	ApplicationNamevar       *string
	Assigneevar              *string
	AssignmentGroupvar       *string
	CategoryRegexvar         *string
	ExcludePortCSVvar        *string
	ExcludeVulnTitleRegexvar *string
	GroupIDvar               *string
	HostnameRegexvar         *string
	OSRegexvar               *string
	OrganizationIDvar        string
	PortCSVvar               *string
	Priorityvar              int
	TagKeyIDvar              *int
	TagKeyRegexvar           *string
	VulnTitleRegexvar        *string
}

//**********************************************************
// Struct Methods
//**********************************************************

// MarshalJSON marshals the struct by converting it to a map
func (myAssignmentRules AssignmentRules) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"ApplicationName":       myAssignmentRules.ApplicationNamevar,
		"Assignee":              myAssignmentRules.Assigneevar,
		"AssignmentGroup":       myAssignmentRules.AssignmentGroupvar,
		"CategoryRegex":         myAssignmentRules.CategoryRegexvar,
		"ExcludePortCSV":        myAssignmentRules.ExcludePortCSVvar,
		"ExcludeVulnTitleRegex": myAssignmentRules.ExcludeVulnTitleRegexvar,
		"GroupID":               myAssignmentRules.GroupIDvar,
		"HostnameRegex":         myAssignmentRules.HostnameRegexvar,
		"OSRegex":               myAssignmentRules.OSRegexvar,
		"OrganizationID":        myAssignmentRules.OrganizationIDvar,
		"PortCSV":               myAssignmentRules.PortCSVvar,
		"Priority":              myAssignmentRules.Priorityvar,
		"TagKeyID":              myAssignmentRules.TagKeyIDvar,
		"TagKeyRegex":           myAssignmentRules.TagKeyRegexvar,
		"VulnTitleRegex":        myAssignmentRules.VulnTitleRegexvar,
	})
}

// ApplicationName returns the ApplicationName parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) ApplicationName() (param *string) {
	return myAssignmentRules.ApplicationNamevar
}

// Assignee returns the Assignee parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) Assignee() (param *string) {
	return myAssignmentRules.Assigneevar
}

// AssignmentGroup returns the AssignmentGroup parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) AssignmentGroup() (param *string) {
	return myAssignmentRules.AssignmentGroupvar
}

// CategoryRegex returns the CategoryRegex parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) CategoryRegex() (param *string) {
	return myAssignmentRules.CategoryRegexvar
}

// ExcludePortCSV returns the ExcludePortCSV parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) ExcludePortCSV() (param *string) {
	return myAssignmentRules.ExcludePortCSVvar
}

// ExcludeVulnTitleRegex returns the ExcludeVulnTitleRegex parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) ExcludeVulnTitleRegex() (param *string) {
	return myAssignmentRules.ExcludeVulnTitleRegexvar
}

// GroupID returns the GroupID parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) GroupID() (param *string) {
	return myAssignmentRules.GroupIDvar
}

// HostnameRegex returns the HostnameRegex parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) HostnameRegex() (param *string) {
	return myAssignmentRules.HostnameRegexvar
}

// OSRegex returns the OSRegex parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) OSRegex() (param *string) {
	return myAssignmentRules.OSRegexvar
}

// OrganizationID returns the OrganizationID parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) OrganizationID() (param string) {
	return myAssignmentRules.OrganizationIDvar
}

// PortCSV returns the PortCSV parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) PortCSV() (param *string) {
	return myAssignmentRules.PortCSVvar
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

// SetApplicationName sets the ApplicationName parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) SetApplicationName(val string) {
	myAssignmentRules.ApplicationNamevar = &val
}

// SetAssignee sets the Assignee parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) SetAssignee(val string) {
	myAssignmentRules.Assigneevar = &val
}

// SetAssignmentGroup sets the AssignmentGroup parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) SetAssignmentGroup(val string) {
	myAssignmentRules.AssignmentGroupvar = &val
}

// SetCategoryRegex sets the CategoryRegex parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) SetCategoryRegex(val string) {
	myAssignmentRules.CategoryRegexvar = &val
}

// SetExcludePortCSV sets the ExcludePortCSV parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) SetExcludePortCSV(val string) {
	myAssignmentRules.ExcludePortCSVvar = &val
}

// SetExcludeVulnTitleRegex sets the ExcludeVulnTitleRegex parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) SetExcludeVulnTitleRegex(val string) {
	myAssignmentRules.ExcludeVulnTitleRegexvar = &val
}

// SetGroupID sets the GroupID parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) SetGroupID(val string) {
	myAssignmentRules.GroupIDvar = &val
}

// SetHostnameRegex sets the HostnameRegex parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) SetHostnameRegex(val string) {
	myAssignmentRules.HostnameRegexvar = &val
}

// SetOSRegex sets the OSRegex parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) SetOSRegex(val string) {
	myAssignmentRules.OSRegexvar = &val
}

// SetOrganizationID sets the OrganizationID parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) SetOrganizationID(val string) {
	myAssignmentRules.OrganizationIDvar = val
}

// SetPortCSV sets the PortCSV parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) SetPortCSV(val string) {
	myAssignmentRules.PortCSVvar = &val
}

// SetTagKeyRegex sets the TagKeyRegex parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) SetTagKeyRegex(val string) {
	myAssignmentRules.TagKeyRegexvar = &val
}

// SetVulnTitleRegex sets the VulnTitleRegex parameter from the AssignmentRules struct
func (myAssignmentRules *AssignmentRules) SetVulnTitleRegex(val string) {
	myAssignmentRules.VulnTitleRegexvar = &val
}
