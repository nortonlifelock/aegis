package domain

// Result defines the interface
type Result interface {
	Success() bool
	Err() error
}
