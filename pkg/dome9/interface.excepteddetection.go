package domain

import "time"

type ExceptedDetection interface {
	Title() *string
	IP() *string
	Hostname() *string
	VulnerabilityID() *string
	VulnerabilityTitle() *string
	Approval() *string
	DueDate() *time.Time
	AssignmentGroup() *string
	OS() *string
	OSRegex() *string
	IgnoreID() string
	IgnoreType() int
}
