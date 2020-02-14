package domain

import (
	"time"
)

// Source defines the interface
type Source interface {
	DBCreatedDate() (param time.Time)
	DBUpdatedDate() (param *time.Time)
	ID() (param string)
	Source() (param string)
	SourceTypeID() (param int)
}
