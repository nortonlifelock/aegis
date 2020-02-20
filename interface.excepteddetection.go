package domain

import "time"

type ExceptedDetection interface {
	Title() *string
	IP() *string
	Hostname() *string
	VulnerabilityID() *string
	VulnerabilityTitle() *string
	Approval() *string
	Expires() *time.Time
	AssignmentGroup() *string
	OS() *string
	IgnoreID() string
}
