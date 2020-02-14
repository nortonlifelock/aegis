package domain

import (
	"time"
)

// JobRegistration defines the interface
type JobRegistration interface {
	CreatedBy() (param string)
	CreatedDate() (param time.Time)
	GoStruct() (param string)
	ID() (param int)
	Priority() (param int)
	UpdatedBy() (param *string)
	UpdatedDate() (param *time.Time)
}
