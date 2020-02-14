package domain

// CloudIP defines the interface
type CloudIP interface {
	IP() string
	Region() string
	State() string
	MAC() string
	InstanceID() string
}
