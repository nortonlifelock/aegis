package domain

import (
	"time"
)

// Detection defines the interface
type Detection interface {
	ID() string
	VulnerabilityID() string
	Status() string
	ActiveKernel() *int
	Detected() (*time.Time, error)
	TimesSeen() int
	Proof() string
	Port() int
	Protocol() string
	IgnoreID() (*string, error)
	Updated() time.Time
	Device() (Device, error)
	Vulnerability() (Vulnerability, error)
}
