package domain

import "time"

// Ignore defines the interface
type Ignore interface {
	DeviceID() (param string)
	DueDate() (param *time.Time)
	ID() (param string)
	OrganizationID() (param string)
	VulnerabilityID() (param string)

	Approval() (param string)
	Active() (param bool)
	Port() (param string)
	TypeID() (param int)
	SourceID() string
	CreatedBy() (param *string)
	UpdatedBy() (param *string)
	DBCreatedDate() (param time.Time)
	DBUpdatedDate() (param *time.Time)
}
