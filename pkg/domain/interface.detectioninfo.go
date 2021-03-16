package domain

import "time"

// DetectionInfo defines the interface
type DetectionInfo interface {
	ID() string
	OrganizationID() string
	SourceID() string
	DeviceID() string
	VulnerabilityID() string
	AlertDate() time.Time
	LastFound() *time.Time
	LastUpdated() *time.Time
	Proof() string
	DetectionStatusID() int
	TimesSeen() int
	Port() int
	IgnoreID() *string
	Updated() time.Time
	Protocol() string
	ActiveKernel() *int
	ParentDetectionID() *string
}
