package domain

// DeviceGroup defines the interface
type DeviceGroup interface {
	Description() (param *string)
	SourceIdentifier() (param int)
}
