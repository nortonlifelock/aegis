package domain

type AssignmentRules interface {
	AssignmentGroup() *string
	Assignee() *string
	ApplicationName() *string
	OrganizationID() string
	GroupID() *string
	VulnTitleRegex() *string
	ExcludeVulnTitleRegex() *string
	HostnameRegex() *string
	OSRegex() *string
	CategoryRegex() *string
	TagKeyID() *int
	TagKeyRegex() *string
	PortCSV() *string
	ExcludePortCSV() *string
	Priority() int
}
