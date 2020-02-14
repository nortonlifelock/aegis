package domain

// OperatingSystemType defines the interface
type OperatingSystemType interface {
	ID() (param int)
	Type() (param string)
	Match() (param string)
	Priority() (param int)
}
