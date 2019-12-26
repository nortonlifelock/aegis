package domain

import "fmt"

// Solution defines the interface
type Solution interface {
	fmt.Stringer
	Summary() string
	Steps() string
}
