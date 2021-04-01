package domain

import "time"

// Finding holds information pertaining to a CIS scanner
type Finding interface {
	// ID corresponds to a vulnerability ID
	ID() string

	// DeviceID corresponds to the entity violating the rule
	DeviceID() string

	// AccountID corresponds to the cloud account that the entity lies within
	AccountID() string

	// ScanID corresponds to the assessment that found the finding
	ScanID() int

	Summary() string
	VulnerabilityTitle() string
	Priority() string

	// String extracts relevant information from the finding
	String() string

	BundleID() string

	LastFound() time.Time
}
