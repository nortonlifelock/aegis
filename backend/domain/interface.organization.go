package domain

import (
	"time"
)

// Organization defines the interface
type Organization interface {
	Code() (param string)
	Created() (param time.Time)
	Description() (param *string)
	EncryptionKey() (param *string)
	ID() (param string)
	ParentOrgID() (param *string)
	Payload() (param string)
	TimeZoneOffset() (param float32)
	Updated() (param *time.Time)
}
