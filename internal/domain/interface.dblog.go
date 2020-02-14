package domain

import (
	"time"
)

// DBLog defines the interface
type DBLog interface {
	CreateDate() (param time.Time)
	Error() (param string)
	ID() (param int)
	JobHistoryID() (param string)
	Log() (param string)
	TypeID() (param int)
}
