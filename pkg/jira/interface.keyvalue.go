package domain

// KeyValue defines the interface
type KeyValue interface {
	Key() (param string)
	Value() (param string)
}
