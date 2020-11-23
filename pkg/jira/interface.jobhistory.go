package domain

import (
	"time"
)

// JobHistory defines the interface
type JobHistory interface {
	ConfigID() (param string)
	CreatedDate() (param time.Time)
	CurrentIteration() (param *int)
	ID() (param string)
	Identifier() (param *string)
	JobID() (param int)
	ParentJobID() (param *string)
	Payload() (param string)
	Priority() (param int)
	PulseDate() (param *time.Time)
	StatusID() (param int)
	ThreadID() (param *string)
	UpdatedDate() (param *time.Time)

	MaxInstances() int
}
