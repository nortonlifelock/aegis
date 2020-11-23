package domain

import (
	"time"
)

// Ticket defines the interface
type Ticket interface {
	AlertDate() (param *time.Time)
	AssignedTo() (param *string)
	AssignmentGroup() (param *string)
	ApplicationName() (param *string)
	Category() (param *string)
	CERF() (param string)
	ExceptionExpiration() (param time.Time)
	CVEReferences() (param *string)
	CVSS() (param *float32)
	CloudID() (param string)
	Configs() (param string)
	CreatedDate() (param *time.Time)
	DBCreatedDate() (param time.Time)
	DBUpdatedDate() (param *time.Time)
	Description() (param *string)
	DeviceID() (param string)
	DueDate() (param *time.Time)
	ExceptionDate() (param *time.Time)
	GroupID() string
	HostName() (param *string)
	ID() (param int)
	IPAddress() (param *string)
	Labels() (param *string)
	LastChecked() (param *time.Time)
	MacAddress() (param *string)
	MethodOfDiscovery() (param *string)
	OSDetailed() (param *string)
	OperatingSystem() (param *string)
	OrgCode() (param *string)
	OrganizationID() (param string)
	OWASP() (param *string)
	Patchable() (param *string)
	Priority() (param *string)
	Project() (param *string)
	ReportedBy() (param *string)
	ResolutionDate() (param *time.Time)
	ResolutionStatus() (param *string)
	ScanID() (param int)
	ServicePorts() (param *string)
	Solution() (param *string)
	Status() (param *string)
	Summary() (param *string)
	SystemName() (param *string)
	TicketType() (param *string)
	Title() (param string)
	UpdatedDate() (param *time.Time)
	VendorReferences() (param *string)
	VulnerabilityID() (param string)
	VulnerabilityTitle() (param *string)

	// fields relevant to BlackDuck
	HubProjectName() *string
	HubProjectVersion() *string
	HubSeverity() *string
	ComponentName() *string
	ComponentVersion() *string
	PolicyRule() *string
	PolicySeverity() *string
}
