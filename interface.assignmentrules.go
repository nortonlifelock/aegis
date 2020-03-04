package domain

type AssignmentRules interface {
	AssignmentGroup() *string
	Assignee() *string
	OrganizationID() string
	VulnTitleSubstring() *string
	VulnTitleRegex() *string
	TagKeyID() *int
	TagKeyValue() *string
	Priority() int
}
