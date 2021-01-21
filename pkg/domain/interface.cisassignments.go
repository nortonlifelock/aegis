package domain

// CISAssignments defines the interface
type CISAssignments interface {
	AssignmentGroup() (param string)
	BundleID() (param *string)
	CloudAccountID() (param *string)
	OrganizationID() (param string)
	RuleID() (param *string)
	DeviceIDRegex() (param *string)
	RuleRegex() (param *string)
	Priority() (param int)
}
