package domain

// User defines the interface
type User interface {
	Email() (param string)
	FirstName() (param string)
	ID() (param string)
	IsDisabled() (param bool)
	LastName() (param string)
	Username() (param *string)
}
