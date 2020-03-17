package domain

import "time"

// AssetGroup defines the interface
type AssetGroup interface {
	GroupID() int
	ScannerSourceID() string
	CloudSourceID() *string
	LastTicketing() *time.Time
}
