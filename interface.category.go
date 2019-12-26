package domain

import ()

// Category defines the interface
type Category interface {
	Category() (param string)
	ID() (param string)
	ParentCategoryID() (param *string)
}
