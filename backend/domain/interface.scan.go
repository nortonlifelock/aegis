package domain

// Scan is an interface that captures important information regarding the scan that is created inside of a scanner
type Scan interface {
	ID() string
	Status() (string, error)
}
