package dal

//**********************************************************
// GENERATED CODE - DO NOT CHANGE
// This file is generated using scaffolding. Any changes to
// this file will be overwritten on the next build
//**********************************************************

import (
	"encoding/json"
	"time"
)

//**********************************************************
// Struct Declaration
//**********************************************************

// Ticket defines the struct that implements the Ticket interface
type Ticket struct {
	AlertDatevar           *time.Time
	ApplicationNamevar     *string
	AssignedTovar          *string
	AssignmentGroupvar     *string
	CERFvar                string
	CVEReferencesvar       *string
	CVSSvar                *float32
	Categoryvar            *string
	CloudIDvar             string
	ComponentNamevar       *string
	ComponentVersionvar    *string
	Configsvar             string
	CreatedDatevar         *time.Time
	DBCreatedDatevar       time.Time
	DBUpdatedDatevar       *time.Time
	Descriptionvar         *string
	DeviceIDvar            string
	DueDatevar             *time.Time
	ExceptionDatevar       *time.Time
	ExceptionExpirationvar time.Time
	GroupIDvar             string
	HostNamevar            *string
	HubProjectNamevar      *string
	HubProjectVersionvar   *string
	IDvar                  int
	IPAddressvar           *string
	Labelsvar              *string
	LastCheckedvar         *time.Time
	MacAddressvar          *string
	MethodOfDiscoveryvar   *string
	OSDetailedvar          *string
	OWASPvar               *string
	OperatingSystemvar     *string
	OrgCodevar             *string
	OrganizationIDvar      string
	Patchablevar           *string
	PolicyRulevar          *string
	PolicySeverityvar      *string
	Priorityvar            *string
	Projectvar             *string
	ReportedByvar          *string
	ResolutionDatevar      *time.Time
	ResolutionStatusvar    *string
	ScanIDvar              int
	ServicePortsvar        *string
	Solutionvar            *string
	Statusvar              *string
	Summaryvar             *string
	SystemNamevar          *string
	TicketTypevar          *string
	Titlevar               string
	UpdatedDatevar         *time.Time
	VendorReferencesvar    *string
	VulnerabilityIDvar     string
	VulnerabilityTitlevar  *string
}

//**********************************************************
// Struct Methods
//**********************************************************

// MarshalJSON marshals the struct by converting it to a map
func (myTicket Ticket) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"AlertDate":           myTicket.AlertDatevar,
		"ApplicationName":     myTicket.ApplicationNamevar,
		"AssignedTo":          myTicket.AssignedTovar,
		"AssignmentGroup":     myTicket.AssignmentGroupvar,
		"CERF":                myTicket.CERFvar,
		"CVEReferences":       myTicket.CVEReferencesvar,
		"CVSS":                myTicket.CVSSvar,
		"Category":            myTicket.Categoryvar,
		"CloudID":             myTicket.CloudIDvar,
		"ComponentName":       myTicket.ComponentNamevar,
		"ComponentVersion":    myTicket.ComponentVersionvar,
		"Configs":             myTicket.Configsvar,
		"CreatedDate":         myTicket.CreatedDatevar,
		"DBCreatedDate":       myTicket.DBCreatedDatevar,
		"DBUpdatedDate":       myTicket.DBUpdatedDatevar,
		"Description":         myTicket.Descriptionvar,
		"DeviceID":            myTicket.DeviceIDvar,
		"DueDate":             myTicket.DueDatevar,
		"ExceptionDate":       myTicket.ExceptionDatevar,
		"ExceptionExpiration": myTicket.ExceptionExpirationvar,
		"GroupID":             myTicket.GroupIDvar,
		"HostName":            myTicket.HostNamevar,
		"HubProjectName":      myTicket.HubProjectNamevar,
		"HubProjectVersion":   myTicket.HubProjectVersionvar,
		"ID":                  myTicket.IDvar,
		"IPAddress":           myTicket.IPAddressvar,
		"Labels":              myTicket.Labelsvar,
		"LastChecked":         myTicket.LastCheckedvar,
		"MacAddress":          myTicket.MacAddressvar,
		"MethodOfDiscovery":   myTicket.MethodOfDiscoveryvar,
		"OSDetailed":          myTicket.OSDetailedvar,
		"OWASP":               myTicket.OWASPvar,
		"OperatingSystem":     myTicket.OperatingSystemvar,
		"OrgCode":             myTicket.OrgCodevar,
		"OrganizationID":      myTicket.OrganizationIDvar,
		"Patchable":           myTicket.Patchablevar,
		"PolicyRule":          myTicket.PolicyRulevar,
		"PolicySeverity":      myTicket.PolicySeverityvar,
		"Priority":            myTicket.Priorityvar,
		"Project":             myTicket.Projectvar,
		"ReportedBy":          myTicket.ReportedByvar,
		"ResolutionDate":      myTicket.ResolutionDatevar,
		"ResolutionStatus":    myTicket.ResolutionStatusvar,
		"ScanID":              myTicket.ScanIDvar,
		"ServicePorts":        myTicket.ServicePortsvar,
		"Solution":            myTicket.Solutionvar,
		"Status":              myTicket.Statusvar,
		"Summary":             myTicket.Summaryvar,
		"SystemName":          myTicket.SystemNamevar,
		"TicketType":          myTicket.TicketTypevar,
		"Title":               myTicket.Titlevar,
		"UpdatedDate":         myTicket.UpdatedDatevar,
		"VendorReferences":    myTicket.VendorReferencesvar,
		"VulnerabilityID":     myTicket.VulnerabilityIDvar,
		"VulnerabilityTitle":  myTicket.VulnerabilityTitlevar,
	})
}

// AlertDate returns the AlertDate parameter from the Ticket struct
func (myTicket *Ticket) AlertDate() (param *time.Time) {
	return myTicket.AlertDatevar
}

// ApplicationName returns the ApplicationName parameter from the Ticket struct
func (myTicket *Ticket) ApplicationName() (param *string) {
	return myTicket.ApplicationNamevar
}

// AssignedTo returns the AssignedTo parameter from the Ticket struct
func (myTicket *Ticket) AssignedTo() (param *string) {
	return myTicket.AssignedTovar
}

// AssignmentGroup returns the AssignmentGroup parameter from the Ticket struct
func (myTicket *Ticket) AssignmentGroup() (param *string) {
	return myTicket.AssignmentGroupvar
}

// CERF returns the CERF parameter from the Ticket struct
func (myTicket *Ticket) CERF() (param string) {
	return myTicket.CERFvar
}

// CVEReferences returns the CVEReferences parameter from the Ticket struct
func (myTicket *Ticket) CVEReferences() (param *string) {
	return myTicket.CVEReferencesvar
}

// CVSS returns the CVSS parameter from the Ticket struct
func (myTicket *Ticket) CVSS() (param *float32) {
	return myTicket.CVSSvar
}

// Category returns the Category parameter from the Ticket struct
func (myTicket *Ticket) Category() (param *string) {
	return myTicket.Categoryvar
}

// CloudID returns the CloudID parameter from the Ticket struct
func (myTicket *Ticket) CloudID() (param string) {
	return myTicket.CloudIDvar
}

// ComponentName returns the ComponentName parameter from the Ticket struct
func (myTicket *Ticket) ComponentName() (param *string) {
	return myTicket.ComponentNamevar
}

// ComponentVersion returns the ComponentVersion parameter from the Ticket struct
func (myTicket *Ticket) ComponentVersion() (param *string) {
	return myTicket.ComponentVersionvar
}

// Configs returns the Configs parameter from the Ticket struct
func (myTicket *Ticket) Configs() (param string) {
	return myTicket.Configsvar
}

// CreatedDate returns the CreatedDate parameter from the Ticket struct
func (myTicket *Ticket) CreatedDate() (param *time.Time) {
	return myTicket.CreatedDatevar
}

// DBCreatedDate returns the DBCreatedDate parameter from the Ticket struct
func (myTicket *Ticket) DBCreatedDate() (param time.Time) {
	return myTicket.DBCreatedDatevar
}

// DBUpdatedDate returns the DBUpdatedDate parameter from the Ticket struct
func (myTicket *Ticket) DBUpdatedDate() (param *time.Time) {
	return myTicket.DBUpdatedDatevar
}

// Description returns the Description parameter from the Ticket struct
func (myTicket *Ticket) Description() (param *string) {
	return myTicket.Descriptionvar
}

// DeviceID returns the DeviceID parameter from the Ticket struct
func (myTicket *Ticket) DeviceID() (param string) {
	return myTicket.DeviceIDvar
}

// DueDate returns the DueDate parameter from the Ticket struct
func (myTicket *Ticket) DueDate() (param *time.Time) {
	return myTicket.DueDatevar
}

// ExceptionDate returns the ExceptionDate parameter from the Ticket struct
func (myTicket *Ticket) ExceptionDate() (param *time.Time) {
	return myTicket.ExceptionDatevar
}

// ExceptionExpiration returns the ExceptionExpiration parameter from the Ticket struct
func (myTicket *Ticket) ExceptionExpiration() (param time.Time) {
	return myTicket.ExceptionExpirationvar
}

// GroupID returns the GroupID parameter from the Ticket struct
func (myTicket *Ticket) GroupID() (param string) {
	return myTicket.GroupIDvar
}

// HostName returns the HostName parameter from the Ticket struct
func (myTicket *Ticket) HostName() (param *string) {
	return myTicket.HostNamevar
}

// HubProjectName returns the HubProjectName parameter from the Ticket struct
func (myTicket *Ticket) HubProjectName() (param *string) {
	return myTicket.HubProjectNamevar
}

// HubProjectVersion returns the HubProjectVersion parameter from the Ticket struct
func (myTicket *Ticket) HubProjectVersion() (param *string) {
	return myTicket.HubProjectVersionvar
}

// ID returns the ID parameter from the Ticket struct
func (myTicket *Ticket) ID() (param int) {
	return myTicket.IDvar
}

// IPAddress returns the IPAddress parameter from the Ticket struct
func (myTicket *Ticket) IPAddress() (param *string) {
	return myTicket.IPAddressvar
}

// Labels returns the Labels parameter from the Ticket struct
func (myTicket *Ticket) Labels() (param *string) {
	return myTicket.Labelsvar
}

// LastChecked returns the LastChecked parameter from the Ticket struct
func (myTicket *Ticket) LastChecked() (param *time.Time) {
	return myTicket.LastCheckedvar
}

// MacAddress returns the MacAddress parameter from the Ticket struct
func (myTicket *Ticket) MacAddress() (param *string) {
	return myTicket.MacAddressvar
}

// MethodOfDiscovery returns the MethodOfDiscovery parameter from the Ticket struct
func (myTicket *Ticket) MethodOfDiscovery() (param *string) {
	return myTicket.MethodOfDiscoveryvar
}

// OSDetailed returns the OSDetailed parameter from the Ticket struct
func (myTicket *Ticket) OSDetailed() (param *string) {
	return myTicket.OSDetailedvar
}

// OWASP returns the OWASP parameter from the Ticket struct
func (myTicket *Ticket) OWASP() (param *string) {
	return myTicket.OWASPvar
}

// OperatingSystem returns the OperatingSystem parameter from the Ticket struct
func (myTicket *Ticket) OperatingSystem() (param *string) {
	return myTicket.OperatingSystemvar
}

// OrgCode returns the OrgCode parameter from the Ticket struct
func (myTicket *Ticket) OrgCode() (param *string) {
	return myTicket.OrgCodevar
}

// OrganizationID returns the OrganizationID parameter from the Ticket struct
func (myTicket *Ticket) OrganizationID() (param string) {
	return myTicket.OrganizationIDvar
}

// Patchable returns the Patchable parameter from the Ticket struct
func (myTicket *Ticket) Patchable() (param *string) {
	return myTicket.Patchablevar
}

// PolicyRule returns the PolicyRule parameter from the Ticket struct
func (myTicket *Ticket) PolicyRule() (param *string) {
	return myTicket.PolicyRulevar
}

// PolicySeverity returns the PolicySeverity parameter from the Ticket struct
func (myTicket *Ticket) PolicySeverity() (param *string) {
	return myTicket.PolicySeverityvar
}

// Priority returns the Priority parameter from the Ticket struct
func (myTicket *Ticket) Priority() (param *string) {
	return myTicket.Priorityvar
}

// Project returns the Project parameter from the Ticket struct
func (myTicket *Ticket) Project() (param *string) {
	return myTicket.Projectvar
}

// ReportedBy returns the ReportedBy parameter from the Ticket struct
func (myTicket *Ticket) ReportedBy() (param *string) {
	return myTicket.ReportedByvar
}

// ResolutionDate returns the ResolutionDate parameter from the Ticket struct
func (myTicket *Ticket) ResolutionDate() (param *time.Time) {
	return myTicket.ResolutionDatevar
}

// ResolutionStatus returns the ResolutionStatus parameter from the Ticket struct
func (myTicket *Ticket) ResolutionStatus() (param *string) {
	return myTicket.ResolutionStatusvar
}

// ScanID returns the ScanID parameter from the Ticket struct
func (myTicket *Ticket) ScanID() (param int) {
	return myTicket.ScanIDvar
}

// ServicePorts returns the ServicePorts parameter from the Ticket struct
func (myTicket *Ticket) ServicePorts() (param *string) {
	return myTicket.ServicePortsvar
}

// Solution returns the Solution parameter from the Ticket struct
func (myTicket *Ticket) Solution() (param *string) {
	return myTicket.Solutionvar
}

// Status returns the Status parameter from the Ticket struct
func (myTicket *Ticket) Status() (param *string) {
	return myTicket.Statusvar
}

// Summary returns the Summary parameter from the Ticket struct
func (myTicket *Ticket) Summary() (param *string) {
	return myTicket.Summaryvar
}

// SystemName returns the SystemName parameter from the Ticket struct
func (myTicket *Ticket) SystemName() (param *string) {
	return myTicket.SystemNamevar
}

// TicketType returns the TicketType parameter from the Ticket struct
func (myTicket *Ticket) TicketType() (param *string) {
	return myTicket.TicketTypevar
}

// Title returns the Title parameter from the Ticket struct
func (myTicket *Ticket) Title() (param string) {
	return myTicket.Titlevar
}

// UpdatedDate returns the UpdatedDate parameter from the Ticket struct
func (myTicket *Ticket) UpdatedDate() (param *time.Time) {
	return myTicket.UpdatedDatevar
}

// VendorReferences returns the VendorReferences parameter from the Ticket struct
func (myTicket *Ticket) VendorReferences() (param *string) {
	return myTicket.VendorReferencesvar
}

// VulnerabilityID returns the VulnerabilityID parameter from the Ticket struct
func (myTicket *Ticket) VulnerabilityID() (param string) {
	return myTicket.VulnerabilityIDvar
}

// VulnerabilityTitle returns the VulnerabilityTitle parameter from the Ticket struct
func (myTicket *Ticket) VulnerabilityTitle() (param *string) {
	return myTicket.VulnerabilityTitlevar
}

// SetApplicationName sets the ApplicationName parameter from the Ticket struct
func (myTicket *Ticket) SetApplicationName(val string) {
	myTicket.ApplicationNamevar = &val
}

// SetAssignedTo sets the AssignedTo parameter from the Ticket struct
func (myTicket *Ticket) SetAssignedTo(val string) {
	myTicket.AssignedTovar = &val
}

// SetAssignmentGroup sets the AssignmentGroup parameter from the Ticket struct
func (myTicket *Ticket) SetAssignmentGroup(val string) {
	myTicket.AssignmentGroupvar = &val
}

// SetCERF sets the CERF parameter from the Ticket struct
func (myTicket *Ticket) SetCERF(val string) {
	myTicket.CERFvar = val
}

// SetCVEReferences sets the CVEReferences parameter from the Ticket struct
func (myTicket *Ticket) SetCVEReferences(val string) {
	myTicket.CVEReferencesvar = &val
}

// SetCategory sets the Category parameter from the Ticket struct
func (myTicket *Ticket) SetCategory(val string) {
	myTicket.Categoryvar = &val
}

// SetCloudID sets the CloudID parameter from the Ticket struct
func (myTicket *Ticket) SetCloudID(val string) {
	myTicket.CloudIDvar = val
}

// SetComponentName sets the ComponentName parameter from the Ticket struct
func (myTicket *Ticket) SetComponentName(val string) {
	myTicket.ComponentNamevar = &val
}

// SetComponentVersion sets the ComponentVersion parameter from the Ticket struct
func (myTicket *Ticket) SetComponentVersion(val string) {
	myTicket.ComponentVersionvar = &val
}

// SetConfigs sets the Configs parameter from the Ticket struct
func (myTicket *Ticket) SetConfigs(val string) {
	myTicket.Configsvar = val
}

// SetDeviceID sets the DeviceID parameter from the Ticket struct
func (myTicket *Ticket) SetDeviceID(val string) {
	myTicket.DeviceIDvar = val
}

// SetGroupID sets the GroupID parameter from the Ticket struct
func (myTicket *Ticket) SetGroupID(val string) {
	myTicket.GroupIDvar = val
}

// SetHostName sets the HostName parameter from the Ticket struct
func (myTicket *Ticket) SetHostName(val string) {
	myTicket.HostNamevar = &val
}

// SetHubProjectName sets the HubProjectName parameter from the Ticket struct
func (myTicket *Ticket) SetHubProjectName(val string) {
	myTicket.HubProjectNamevar = &val
}

// SetHubProjectVersion sets the HubProjectVersion parameter from the Ticket struct
func (myTicket *Ticket) SetHubProjectVersion(val string) {
	myTicket.HubProjectVersionvar = &val
}

// SetIPAddress sets the IPAddress parameter from the Ticket struct
func (myTicket *Ticket) SetIPAddress(val string) {
	myTicket.IPAddressvar = &val
}

// SetLabels sets the Labels parameter from the Ticket struct
func (myTicket *Ticket) SetLabels(val string) {
	myTicket.Labelsvar = &val
}

// SetMacAddress sets the MacAddress parameter from the Ticket struct
func (myTicket *Ticket) SetMacAddress(val string) {
	myTicket.MacAddressvar = &val
}

// SetMethodOfDiscovery sets the MethodOfDiscovery parameter from the Ticket struct
func (myTicket *Ticket) SetMethodOfDiscovery(val string) {
	myTicket.MethodOfDiscoveryvar = &val
}

// SetOSDetailed sets the OSDetailed parameter from the Ticket struct
func (myTicket *Ticket) SetOSDetailed(val string) {
	myTicket.OSDetailedvar = &val
}

// SetOWASP sets the OWASP parameter from the Ticket struct
func (myTicket *Ticket) SetOWASP(val string) {
	myTicket.OWASPvar = &val
}

// SetOperatingSystem sets the OperatingSystem parameter from the Ticket struct
func (myTicket *Ticket) SetOperatingSystem(val string) {
	myTicket.OperatingSystemvar = &val
}

// SetOrgCode sets the OrgCode parameter from the Ticket struct
func (myTicket *Ticket) SetOrgCode(val string) {
	myTicket.OrgCodevar = &val
}

// SetOrganizationID sets the OrganizationID parameter from the Ticket struct
func (myTicket *Ticket) SetOrganizationID(val string) {
	myTicket.OrganizationIDvar = val
}

// SetPatchable sets the Patchable parameter from the Ticket struct
func (myTicket *Ticket) SetPatchable(val string) {
	myTicket.Patchablevar = &val
}

// SetPolicyRule sets the PolicyRule parameter from the Ticket struct
func (myTicket *Ticket) SetPolicyRule(val string) {
	myTicket.PolicyRulevar = &val
}

// SetPolicySeverity sets the PolicySeverity parameter from the Ticket struct
func (myTicket *Ticket) SetPolicySeverity(val string) {
	myTicket.PolicySeverityvar = &val
}

// SetPriority sets the Priority parameter from the Ticket struct
func (myTicket *Ticket) SetPriority(val string) {
	myTicket.Priorityvar = &val
}

// SetProject sets the Project parameter from the Ticket struct
func (myTicket *Ticket) SetProject(val string) {
	myTicket.Projectvar = &val
}

// SetReportedBy sets the ReportedBy parameter from the Ticket struct
func (myTicket *Ticket) SetReportedBy(val string) {
	myTicket.ReportedByvar = &val
}

// SetResolutionStatus sets the ResolutionStatus parameter from the Ticket struct
func (myTicket *Ticket) SetResolutionStatus(val string) {
	myTicket.ResolutionStatusvar = &val
}

// SetServicePorts sets the ServicePorts parameter from the Ticket struct
func (myTicket *Ticket) SetServicePorts(val string) {
	myTicket.ServicePortsvar = &val
}

// SetStatus sets the Status parameter from the Ticket struct
func (myTicket *Ticket) SetStatus(val string) {
	myTicket.Statusvar = &val
}

// SetSummary sets the Summary parameter from the Ticket struct
func (myTicket *Ticket) SetSummary(val string) {
	myTicket.Summaryvar = &val
}

// SetSystemName sets the SystemName parameter from the Ticket struct
func (myTicket *Ticket) SetSystemName(val string) {
	myTicket.SystemNamevar = &val
}

// SetTicketType sets the TicketType parameter from the Ticket struct
func (myTicket *Ticket) SetTicketType(val string) {
	myTicket.TicketTypevar = &val
}

// SetTitle sets the Title parameter from the Ticket struct
func (myTicket *Ticket) SetTitle(val string) {
	myTicket.Titlevar = val
}

// SetVulnerabilityID sets the VulnerabilityID parameter from the Ticket struct
func (myTicket *Ticket) SetVulnerabilityID(val string) {
	myTicket.VulnerabilityIDvar = val
}

// SetVulnerabilityTitle sets the VulnerabilityTitle parameter from the Ticket struct
func (myTicket *Ticket) SetVulnerabilityTitle(val string) {
	myTicket.VulnerabilityTitlevar = &val
}
