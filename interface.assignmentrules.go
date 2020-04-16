package domain

type AssignmentRules interface {
	AssignmentGroup() *string
	Assignee() *string
	OrganizationID() string
	GroupID() *string
	VulnTitleRegex() *string
	HostnameRegex() *string
	TagKeyID() *int
	TagKeyRegex() *string
	Priority() int
}
