package domain

import ()

// DetectionStatus defines the interface
type DetectionStatus interface {
	ID() (param int)
	Name() (param string)
	Status() (param string)
}
