package domain

import (
	"time"
)

// ScanSummary defines the interface
type ScanSummary interface {
	CreatedDate() (param time.Time)
	OrgID() (param string)
	ParentJobID() (param string)
	ScanClosePayload() (param string)
	ScanStatus() (param string)
	Source() (param string)
	SourceID() (param string)
	SourceKey() (param *string)
	TemplateID() (param *string)
	UpdatedDate() (param *time.Time)
}
