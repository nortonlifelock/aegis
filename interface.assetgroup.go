package domain

import "time"

// AssetGroup defines the interface
type AssetGroup interface {
	GroupID() string
	ScannerSourceID() string
	CloudSourceID() *string
	LastTicketing() *time.Time
}
