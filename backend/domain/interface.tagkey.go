package domain

// TagKey defines the interface
type TagKey interface {
	ID() (param string)
	KeyValue() (param string)
}
