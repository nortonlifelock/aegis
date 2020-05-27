package domain

// DeviceInfo defines the interface
type DeviceInfo interface {
	// ID is the ID of the device as reported by the backend database of Aegis
	ID() string

	// SourceID is the ID of the device as reported by the scanner
	SourceID() *string

	// ScannerSourceID is the id of the source (vulnerability scanner) that found the device
	ScannerSourceID() *string

	OS() string
	HostName() string
	MAC() string
	IP() string

	// Region is the area that the device is stored in (if the device is a cloud device)
	Region() *string

	GroupID() *string

	// InstanceID identifies which instance a device is (the the device is a cloud device)
	InstanceID() *string

	State() *string

	//TrackingMethod() *string
}
