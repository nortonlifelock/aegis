package domain

// TagMap defines the interface
type TagMap interface {
	CloudSourceID() (param string)
	CloudTag() (param string)
	ID() (param string)
	Options() (param string)
	TicketingSourceID() (param string)
	TicketingTag() (param string)
}
