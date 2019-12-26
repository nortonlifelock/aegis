package domain

import "context"

// Device defines the interface
type Device interface {
	// ID is the ID of the device as reported by the backend database of PDE
	ID() string

	// SourceID is the ID of the device as reported by the scanner
	SourceID() *string

	OS() string
	HostName() string
	MAC() string
	IP() string
	Vulnerabilities(ctx context.Context) (param <-chan Detection, err error)

	// Region is the area that the device is stored in (if the device is a cloud device)
	Region() *string
	// InstanceID identifies which instance a device is (the the device is a cloud device)
	InstanceID() *string
}
