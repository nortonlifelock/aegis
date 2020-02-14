package domain

// ExceptionType defines the interface
type ExceptionType interface {
	ID() (param int)
	Name() (param string)
	Type() (param string)
}
