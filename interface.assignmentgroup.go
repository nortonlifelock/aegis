package domain

import (
	"time"
)

// AssignmentGroup defines the interface
type AssignmentGroup interface {
	DBCreatedDate() (param time.Time)
	DBUpdatedDate() (param *time.Time)
	GroupName() (param string)
	IPAddress() (param string)
	OrganizationID() (param string)
	SourceID() (param int)
	VulnTitleRegex() (param *string)
}
