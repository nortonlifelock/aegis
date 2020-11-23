package domain

type CategoryRule interface {
	ID() string
	OrganizationID() string
	SourceID() string
	VulnerabilityTitle() *string
	VulnerabilityCategory() *string
	VulnerabilityType() *string
	Category() string
}
