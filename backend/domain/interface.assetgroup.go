package domain

// AssetGroup defines the interface
type AssetGroup interface {
	GroupID() int
	ScannerSourceID() string
	CloudSourceID() *string
}
