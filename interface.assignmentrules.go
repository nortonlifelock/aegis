package domain

type AssignmentRules interface {
	AssignmentGroup() *string
	Assignee() *string
	OrganizationID() string
	VulnTitleRegex() *string
	HostnameRegex() *string
	TagKeyID() *int
	TagKeyRegex() *string
	Priority() int
}
