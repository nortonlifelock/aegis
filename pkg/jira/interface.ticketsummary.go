package domain

import (
	"time"
)

// TicketSummary defines the interface
type TicketSummary interface {
	DetectionID() (param string)
	DueDate() (param time.Time)
	OrganizationID() (param string)
	ResolutionDate() (param *time.Time)
	Status() (param string)
	Title() (param string)
	CreatedDate() (param *time.Time)
	UpdatedDate() (param *time.Time)
}
