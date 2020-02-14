package domain

// Tag defines the interface
type Tag interface {
	DeviceID() (param string)
	ID() (param string)
	TagKeyID() (param int)
	Value() (param string)
}
