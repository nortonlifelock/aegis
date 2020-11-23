package domain

import (
	"time"
)

// SourceConfig defines the interface
type SourceConfig interface {
	Address() (param string)
	AuthInfo() (param string)
	DBCreatedDate() (param time.Time)
	DBUpdatedDate() (param *time.Time)
	ID() (param string)
	OrganizationID() (param string)
	Payload() (param *string)
	Port() (param string)
	Source() (param string)
	SourceID() (param string)
}
