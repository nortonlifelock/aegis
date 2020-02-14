package domain

// JobSchedule defines the interface
type JobSchedule interface {
	ConfigID() (param string)
	ID() (param string)
	Payload() (param *string)
}
