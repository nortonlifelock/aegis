package domain

import (
	"time"
)

// JobConfig defines the interface
type JobConfig interface {
	Active() (param bool)
	AutoStart() (param bool)
	Continuous() (param bool)
	CreatedBy() (param string)
	CreatedDate() (param time.Time)
	DataInSourceConfigID() (param *string)
	DataOutSourceConfigID() (param *string)
	ID() (param string)
	JobID() (param int)
	LastJobStart() (param *time.Time)
	MaxInstances() (param int)
	//Organization() (param Organization)
	OrganizationID() (param string)
	Payload() (param *string)
	PriorityOverride() (param *int)
	UpdatedBy() (param *string)
	UpdatedDate() (param *time.Time)
	WaitInSeconds() (param int)
}
