package domain

import "time"

// JobConfigAudit holds information for the JobConfigAudit table in the database. It holds all the information that a job config has with some additional metadata
type JobConfigAudit interface {
	JobConfig
	EventType() string
	EventDate() time.Time
}
