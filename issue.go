package jira

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/log"
	"strconv"
	"strings"
	"time"
)

// Issue holds the issue returned from JIRA and implements the domain.Ticket interface
type Issue struct {
	Issue     *jira.Issue
	connector *ConnectorJira
}

// SetHostName sets the HostName parameter from the Ticket struct
func (ji *Issue) SetHostName(val string) {
	_ = ji.setCF(ji.connector, backendHostname, val)
}

// SetDueDate sets the DueDate parameter from the Ticket struct
func (ji *Issue) SetDueDate(val time.Time) {
	ji.Issue.Fields.Duedate = jira.Date(val)
}

// SetAlertDate sets the AlertDate parameter from the Ticket struct
func (ji *Issue) SetAlertDate(val time.Time) {
	_ = ji.setCF(ji.connector, backendScanDate, val.Format(DateFormatJira))
}

// UpdatedDate gets the UpdatedDate parameter from the Ticket struct
func (ji *Issue) UpdatedDate() (param *time.Time) {
	param = new(time.Time)
	*param = time.Time(ji.Issue.Fields.Updated)
	return param
}

// UpdatedDateOrDefault gets the UpdatedDate parameter from the Ticket struct
func (ji *Issue) UpdatedDateOrDefault() (param time.Time) {
	// TODO configurable time zone
	param = time.Time(ji.Issue.Fields.Updated).Add(-7)
	return param
}

// Solution gets the Solution parameter from the Ticket struct
func (ji *Issue) Solution() (param *string) {
	return ji.getStringPointer(backendSolution)
}

// SolutionOrDefault gets the Solution parameter from the Ticket struct
func (ji *Issue) SolutionOrDefault() (param string) {
	return ji.getString(backendSolution)
}

// Labels gets the Labels parameter from the Ticket struct
func (ji *Issue) Labels() (param *string) {
	if ji.Issue.Fields.Labels != nil && len(ji.Issue.Fields.Labels) > 0 {
		var ret = strings.Join(ji.Issue.Fields.Labels, ",")
		param = &ret
	}

	return param
}

// LabelsOrDefault gets the Labels parameter from the Ticket struct
func (ji *Issue) LabelsOrDefault() (param string) {
	if ji.Issue.Fields.Labels != nil && len(ji.Issue.Fields.Labels) > 0 {
		var ret = strings.Join(ji.Issue.Fields.Labels, ",")
		param = ret
	}

	return param
}

// SetCERF sets the CERF parameter from the Ticket struct
func (ji *Issue) SetCERF(val string) {
	_ = ji.setCF(ji.connector, backendCERF, val)
}

// AlertDate gets the AlertDate parameter from the Ticket struct
func (ji *Issue) AlertDate() (param *time.Time) {
	return ji.getTimePointer(backendScanDate)
}

// AlertDateOrDefault gets the AlertDate parameter from the Ticket struct
func (ji *Issue) AlertDateOrDefault() (param time.Time) {
	return ji.getTime(backendScanDate)
}

// DeviceID gets the DeviceID parameter from the Ticket struct
func (ji *Issue) DeviceID() (param string) {
	return ji.getString(backendDeviceID)
}

// ResolutionStatus gets the ResolutionStatus parameter from the Ticket struct
func (ji *Issue) ResolutionStatus() (param *string) {
	if ji.Issue.Fields.Resolution != nil {
		param = &ji.Issue.Fields.Resolution.Name
	} else {
		param = new(string)
		*param = "Unresolved"
	}
	return param
}

// ResolutionStatusOrDefault gets the ResolutionStatus parameter from the Ticket struct
func (ji *Issue) ResolutionStatusOrDefault() (param string) {
	if ji.Issue.Fields.Resolution != nil {
		param = ji.Issue.Fields.Resolution.Name
	} else {
		param = "Unresolved"
	}
	return param
}

// CreatedDate gets the CreatedDate parameter from the Ticket struct
func (ji *Issue) CreatedDate() (param *time.Time) {
	param = new(time.Time)
	*param = time.Time(ji.Issue.Fields.Created)
	return param
}

// CreatedDateOrDefault gets the CreatedDate parameter from the Ticket struct
func (ji *Issue) CreatedDateOrDefault() (param time.Time) {
	param = time.Time(ji.Issue.Fields.Created)
	return param
}

// VendorReferences gets the VendorReferences parameter from the Ticket struct
func (ji *Issue) VendorReferences() (param *string) {
	return ji.getStringPointer(backendVendorReferences)
}

// VendorReferencesOrDefault gets the VendorReferences parameter from the Ticket struct
func (ji *Issue) VendorReferencesOrDefault() (param string) {
	return ji.getString(backendVendorReferences)
}

// AssignedTo gets the AssignedTo parameter from the Ticket struct
func (ji *Issue) AssignedTo() (param *string) {
	if ji.Issue.Fields.Assignee != nil {
		param = &ji.Issue.Fields.Assignee.Name
	}
	return param
}

// AssignedToOrDefault gets the AssignedTo parameter from the Ticket struct
func (ji *Issue) AssignedToOrDefault() (param string) {
	if ji.Issue.Fields.Assignee != nil {
		param = ji.Issue.Fields.Assignee.Name
	}
	return param
}

// OrganizationID gets the OrganizationID parameter from the Ticket struct
// The OrganizationID is a database field, and is not relevant to JIRA tickets
func (ji *Issue) OrganizationID() (param string) {
	return param
}

// SetScanID sets the ScanID parameter from the Ticket struct
func (ji *Issue) SetScanID(val int) {
	_ = ji.setCF(ji.connector, backendScanID, strconv.Itoa(val))
}

// SetUpdatedDate sets the UpdatedDate parameter from the Ticket struct
func (ji *Issue) SetUpdatedDate(val time.Time) {
	ji.Issue.Fields.Updated = jira.Time(val)
}

// IPAddress gets the IPAddress parameter from the Ticket struct
func (ji *Issue) IPAddress() (param *string) {
	return ji.getStringPointer(backendIPAddress)
}

// IPAddressOrDefault gets the IPAddress parameter from the Ticket struct
func (ji *Issue) IPAddressOrDefault() (param string) {
	return ji.getString(backendIPAddress)
}

// SetDBCreatedDate sets the DBCreatedDate parameter from the Ticket struct
func (ji *Issue) SetDBCreatedDate(val time.Time) {
	return // we won't have this for the jira issue
}

// SetOperatingSystem sets the OperatingSystem parameter from the Ticket struct
func (ji *Issue) SetOperatingSystem(val string) {
	_ = ji.setCF(ji.connector, backendOperatingSystem, val)
}

// CERF gets the CERF parameter from the Ticket struct
func (ji *Issue) CERF() (param string) {
	return ji.getString(backendCERF)
}

// SetVulnerabilityID sets the VulnerabilityID parameter from the Ticket struct
func (ji *Issue) SetVulnerabilityID(val string) {
	_ = ji.setCF(ji.connector, backendVulnerabilityID, val)
}

// CERFExpirationDate gets the CERFExpirationDate parameter from the Ticket struct
func (ji *Issue) CERFExpirationDate() (param time.Time) {
	var cerfExpiration time.Time

	// the ticket itself is a CERF ticket
	if !ji.getTime(backendCERFExpiration).IsZero() {
		cerfExpiration = ji.getTime(backendCERFExpiration)
	} else {
		// the ticket has a referenced CERF, let's grab it's expiration
		if len(ji.CERF()) > 0 {
			if val, exists := ji.connector.CERFs.Load(ji.CERF()); exists {
				if cerf, ok := val.(domain.Ticket); ok {
					cerfExpiration = cerf.CERFExpirationDate()
				} else {
					ji.connector.lstream.Send(log.Errorf(nil, "cerf [%v] failed to load from cache", ji.CERF()))
				}
			} else {
				if cerf, err := ji.connector.GetTicket(ji.CERF()); err == nil {
					if cerf != nil {
						cerfExpiration = cerf.CERFExpirationDate()
						ji.connector.CERFs.Store(ji.CERF(), cerf)
					} else {
						ji.connector.lstream.Send(log.Errorf(err, "cerf [%v] returned nil from JIRA", ji.CERF()))
					}
				} else {
					ji.connector.lstream.Send(log.Errorf(err, "cerf [%v] failed to load from JIRA", ji.CERF()))
				}
			}
		}
	}

	return cerfExpiration
}

// SetCERFExpirationDate sets the CERFExpirationDate parameter from the Ticket struct
func (ji *Issue) SetCERFExpirationDate(val time.Time) {
	_ = ji.setCF(ji.connector, backendCERFExpiration, val.Format(DateFormatJira))
}

// SetVulnerabilityTitle sets the VulnerabilityTitle parameter from the Ticket struct
func (ji *Issue) SetVulnerabilityTitle(val string) {
	_ = ji.setCF(ji.connector, backendVulnerability, val)
}

// SetOSDetailed sets the OSDetailed parameter from the Ticket struct
func (ji *Issue) SetOSDetailed(val string) {
	_ = ji.setCF(ji.connector, backendOSDetailed, val)
}

// OrgCode gets the OrgCode parameter from the Ticket struct
func (ji *Issue) OrgCode() (param *string) {
	return ji.getStringPointer(backendOrg)
}

// OrgCodeOrDefault gets the OrgCode parameter from the Ticket struct
func (ji *Issue) OrgCodeOrDefault() (param string) {
	return ji.getString(backendOrg)
}

// SetOrgCode sets the OrgCode parameter from the Ticket struct
func (ji *Issue) SetOrgCode(val string) {
	_ = ji.setCF(ji.connector, backendOrg, val)
}

// SetVendorReferences sets the VendorReferences parameter from the Ticket struct
func (ji *Issue) SetVendorReferences(val string) {
	_ = ji.setCF(ji.connector, backendVendorReferences, val)
}

// SetServicePorts sets the ServicePorts parameter from the Ticket struct
func (ji *Issue) SetServicePorts(val string) {
	_ = ji.setCF(ji.connector, backendServicePort, val)
}

// DBUpdatedDate gets the DBUpdatedDate parameter from the Ticket struct
func (ji *Issue) DBUpdatedDate() (param *time.Time) {
	return // jira will not have this information
}

// DBUpdatedDateOrDefault gets the DBUpdatedDate parameter from the Ticket struct
func (ji *Issue) DBUpdatedDateOrDefault() (param time.Time) {
	return // jira will not have this information
}

// SetCVEReferences sets the CVEReferences parameter from the Ticket struct
func (ji *Issue) SetCVEReferences(val string) {
	_ = ji.setCF(ji.connector, backendCVEReferences, val)
}

// SetStatus sets the Status parameter from the Ticket struct
func (ji *Issue) SetStatus(val string) {
	if ji.Issue.Fields.Status != nil {
		ji.Issue.Fields.Status.Name = val
	}
}

// SetGroupID sets the GroupID parameter from the Ticket struct
func (ji *Issue) SetGroupID(val string) {
	_ = ji.setCF(ji.connector, backendGroupID, val)
}

// Description gets the Description parameter from the Ticket struct
func (ji *Issue) Description() (param *string) {
	return &ji.Issue.Fields.Description
}

// DescriptionOrDefault gets the Description parameter from the Ticket struct
func (ji *Issue) DescriptionOrDefault() (param string) {
	return ji.Issue.Fields.Description
}

// Configs gets the Configs parameter from the Ticket struct
func (ji *Issue) Configs() (param string) {
	return ji.getString(backendConfig)
}

// LastChecked gets the LastChecked parameter from the Ticket struct
func (ji *Issue) LastChecked() (param *time.Time) {
	return ji.getTimePointer(backendLastChecked)
}

// LastCheckedOrDefault gets the LastChecked parameter from the Ticket struct
func (ji *Issue) LastCheckedOrDefault() (param time.Time) {
	return ji.getTime(backendLastChecked)
}

// Project gets the Project parameter from the Ticket struct
func (ji *Issue) Project() (param *string) {
	return &ji.Issue.Fields.Project.Name
}

// ProjectOrDefault gets the Project parameter from the Ticket struct
func (ji *Issue) ProjectOrDefault() (param string) {
	return ji.Issue.Fields.Project.Name
}

// Packages gets the Packages parameter from the Ticket struct
func (ji *Issue) Packages() (param *string) {
	// pretty sure this field was removed?
	return
}

// PackagesOrDefault gets the Packages parameter from the Ticket struct
func (ji *Issue) PackagesOrDefault() (param string) {
	// pretty sure this field was removed?
	return
}

// SetCreatedDate sets the CreatedDate parameter from the Ticket struct
func (ji *Issue) SetCreatedDate(val time.Time) {
	ji.Issue.Fields.Created = jira.Time(val.Add(-7))
}

// SetReportedBy sets the ReportedBy parameter from the Ticket struct
func (ji *Issue) SetReportedBy(val string) {
	ji.Issue.Fields.Creator.Name = val
}

// SetTicketType sets the TicketType parameter from the Ticket struct
func (ji *Issue) SetTicketType(val string) {
	ji.Issue.Fields.Type.Name = val
}

// AssignmentGroup gets the AssignmentGroup parameter from the Ticket struct
func (ji *Issue) AssignmentGroup() (param *string) {
	return ji.getStringPointer(backendAssignmentGroup)
}

// AssignmentGroupOrDefault gets the AssignmentGroup parameter from the Ticket struct
func (ji *Issue) AssignmentGroupOrDefault() (param string) {
	return ji.getString(backendAssignmentGroup)
}

// SetMethodOfDiscovery sets the MethodOfDiscovery parameter from the Ticket struct
func (ji *Issue) SetMethodOfDiscovery(val string) {
	_ = ji.setCF(ji.connector, backendMOD, val)
}

// SetPackages sets the Packages parameter from the Ticket struct
func (ji *Issue) SetPackages(val string) {
	// pretty sure this field was removed?
	return
}

// GroupID gets the GroupID parameter from the Ticket struct
func (ji *Issue) GroupID() (param string) {
	return ji.getString(backendGroupID)
}

// SetID sets the ID parameter from the Ticket struct
func (ji *Issue) SetID(val int) {
	_ = ji.setCF(ji.connector, backendAutomationID, val)
}

// SetLabels sets the Labels parameter from the Ticket struct
func (ji *Issue) SetLabels(val string) {
	ji.Issue.Fields.Labels = strings.Split(val, ",")
}

// HostName gets the HostName parameter from the Ticket struct
func (ji *Issue) HostName() (param *string) {
	return ji.getStringPointer(backendHostname)
}

// HostNameOrDefault gets the HostName parameter from the Ticket struct
func (ji *Issue) HostNameOrDefault() (param string) {
	return ji.getString(backendHostname)
}

// ResolutionDate gets the ResolutionDate parameter from the Ticket struct
func (ji *Issue) ResolutionDate() (param *time.Time) {
	param = new(time.Time)
	*param = time.Time(ji.Issue.Fields.Resolutiondate)
	return param
}

// ResolutionDateOrDefault gets the ResolutionDate parameter from the Ticket struct
func (ji *Issue) ResolutionDateOrDefault() (param time.Time) {
	param = time.Time(ji.Issue.Fields.Resolutiondate)
	return param
}

// SetCVSS sets the CVSS parameter from the Ticket struct
func (ji *Issue) SetCVSS(val float32) {
	_ = ji.setCF(ji.connector, backendCVSS, fmt.Sprintf("%f", val))
}

// VulnerabilityTitle gets the VulnerabilityTitle parameter from the Ticket struct
func (ji *Issue) VulnerabilityTitle() (param *string) {
	return ji.getStringPointer(backendVulnerability)
}

// VulnerabilityTitleOrDefault gets the VulnerabilityTitle parameter from the Ticket struct
func (ji *Issue) VulnerabilityTitleOrDefault() (param string) {
	return ji.getString(backendVulnerability)
}

// SetResolutionStatus sets the ResolutionStatus parameter from the Ticket struct
func (ji *Issue) SetResolutionStatus(val string) {
	if ji.Issue.Fields.Resolution != nil {
		ji.Issue.Fields.Resolution.Name = val
	}
}

// DBCreatedDate gets the DBCreatedDate parameter from the Ticket struct
func (ji *Issue) DBCreatedDate() (param time.Time) {
	return
}

// ReportedBy gets the ReportedBy parameter from the Ticket struct
func (ji *Issue) ReportedBy() (param *string) {
	if ji.Issue.Fields.Creator != nil {
		param = &ji.Issue.Fields.Creator.Name
	}
	return
}

// ReportedByOrDefault gets the ReportedBy parameter from the Ticket struct
func (ji *Issue) ReportedByOrDefault() (param string) {
	if ji.Issue.Fields.Creator != nil {
		param = ji.Issue.Fields.Creator.Name
	}
	return
}

// Priority gets the Priority parameter from the Ticket struct
func (ji *Issue) Priority() (param *string) {
	return ji.getStringPointer(backendVRRPriority)
}

// PriorityOrDefault gets the Priority parameter from the Ticket struct
func (ji *Issue) PriorityOrDefault() (param string) {
	return ji.getString(backendVRRPriority)
}

// SetAssignedTo sets the AssignedTo parameter from the Ticket struct
func (ji *Issue) SetAssignedTo(val string) {
	if ji.Issue.Fields.Assignee == nil {
		ji.Issue.Fields.Assignee = &jira.User{}
	}
	ji.Issue.Fields.Assignee.Name = val
}

// Status gets the Status parameter from the Ticket struct
func (ji *Issue) Status() (param *string) {
	if ji.Issue.Fields.Status != nil {
		param = &ji.Issue.Fields.Status.Name
	}

	return
}

// StatusOrDefault gets the Status parameter from the Ticket struct
func (ji *Issue) StatusOrDefault() (param string) {
	if ji.Issue.Fields.Status != nil {
		param = ji.Issue.Fields.Status.Name
	}
	return param
}

// OSDetailed gets the OSDetailed parameter from the Ticket struct
func (ji *Issue) OSDetailed() (param *string) {
	return ji.getStringPointer(backendOSDetailed)
}

// OSDetailedOrDefault gets the OSDetailed parameter from the Ticket struct
func (ji *Issue) OSDetailedOrDefault() (param string) {
	return ji.getString(backendOSDetailed)
}

// SetSummary sets the Summary parameter from the Ticket struct
func (ji *Issue) SetSummary(val string) {
	ji.Issue.Fields.Summary = val
}

// SetAssignmentGroup sets the AssignmentGroup parameter from the Ticket struct
func (ji *Issue) SetAssignmentGroup(val string) {
	_ = ji.setCF(ji.connector, backendAssignmentGroup, val)
}

// MethodOfDiscovery gets the MethodOfDiscovery parameter from the Ticket struct
func (ji *Issue) MethodOfDiscovery() (param *string) {
	return ji.getStringPointer(backendMOD)
}

// MethodOfDiscoveryOrDefault gets the MethodOfDiscovery parameter from the Ticket struct
func (ji *Issue) MethodOfDiscoveryOrDefault() (param string) {
	return ji.getString(backendMOD)
}

// CVEReferences gets the CVEReferences parameter from the Ticket struct
func (ji *Issue) CVEReferences() (param *string) {
	return ji.getStringPointer(backendCVEReferences)
}

// CVEReferencesOrDefault gets the CVEReferences parameter from the Ticket struct
func (ji *Issue) CVEReferencesOrDefault() (param string) {
	return ji.getString(backendCVEReferences)
}

// Title gets the Title parameter from the Ticket struct
func (ji *Issue) Title() (param string) {
	return ji.Issue.Key
}

// SetLastChecked sets the LastChecked parameter from the Ticket struct
func (ji *Issue) SetLastChecked(val time.Time) {
	_ = ji.setCF(ji.connector, backendLastChecked, val.Format(DateFormatJira))
}

// ServicePorts gets the ServicePorts parameter from the Ticket struct
func (ji *Issue) ServicePorts() (param *string) {
	return ji.getStringPointer(backendServicePort)
}

// ServicePortsOrDefault gets the ServicePorts parameter from the Ticket struct
func (ji *Issue) ServicePortsOrDefault() (param string) {
	return ji.getString(backendServicePort)
}

// SetDBUpdatedDate sets the DBUpdatedDate parameter from the Ticket struct
func (ji *Issue) SetDBUpdatedDate(val time.Time) {
	return
}

// SetConfigs sets the Configs parameter from the Ticket struct
func (ji *Issue) SetConfigs(val string) {
	_ = ji.setCF(ji.connector, backendConfig, val)
}

// SetOrganizationID sets the OrganizationID parameter from the Ticket struct
func (ji *Issue) SetOrganizationID(val string) {
	return // not mutable
}

// SetResolutionDate sets the ResolutionDate parameter from the Ticket struct
func (ji *Issue) SetResolutionDate(val time.Time) {
	ji.Issue.Fields.Resolutiondate = jira.Time(val)
}

// ScanID gets the ScanID parameter from the Ticket struct
func (ji *Issue) ScanID() (param int) {
	return ji.getInt(backendScanID)
}

// VulnerabilityID gets the VulnerabilityID parameter from the Ticket struct
func (ji *Issue) VulnerabilityID() (param string) {
	return ji.getString(backendVulnerabilityID)
}

// ID gets the ID parameter from the Ticket struct
func (ji *Issue) ID() (param int) {
	return // DB field
}

// CloudID gets the CloudID parameter from the Ticket struct
func (ji *Issue) CloudID() (param string) {
	return ji.getString(backendCloudID)
}

// SetCloudID sets the CloudID parameter from the Ticket struct
func (ji *Issue) SetCloudID(val string) {
	_ = ji.setCF(ji.connector, backendCloudID, val)
}

// SetPriority sets the Priority parameter from the Ticket struct
func (ji *Issue) SetPriority(val string) {
	_ = ji.setCF(ji.connector, backendVRRPriority, val)
}

// SetSolution sets the Solution parameter from the Ticket struct
func (ji *Issue) SetSolution(val string) {
	_ = ji.setCF(ji.connector, backendSolution, val)
}

// MacAddress gets the MacAddress parameter from the Ticket struct
func (ji *Issue) MacAddress() (param *string) {
	return ji.getStringPointer(backendMACAddress)
}

// MacAddressOrDefault gets the MacAddress parameter from the Ticket struct
func (ji *Issue) MacAddressOrDefault() (param string) {
	return ji.getString(backendMACAddress)
}

// SetMacAddress sets the MacAddress parameter from the Ticket struct
func (ji *Issue) SetMacAddress(val string) {
	_ = ji.setCF(ji.connector, backendMACAddress, val)
}

// Summary gets the Summary parameter from the Ticket struct
func (ji *Issue) Summary() (param *string) {
	return &ji.Issue.Fields.Summary
}

// SummaryOrDefault gets the Summary parameter from the Ticket struct
func (ji *Issue) SummaryOrDefault() (param string) {
	return ji.Issue.Fields.Summary
}

// DueDate gets the DueDate parameter from the Ticket struct
func (ji *Issue) DueDate() (param *time.Time) {
	var retVal = time.Time(ji.Issue.Fields.Duedate)
	return &retVal
}

// DueDateOrDefault gets the DueDate parameter from the Ticket struct
func (ji *Issue) DueDateOrDefault() (param time.Time) {
	return time.Time(ji.Issue.Fields.Duedate)
}

// CVSS gets the CVSS parameter from the Ticket struct
func (ji *Issue) CVSS() (param *float32) {
	return ji.getFloatPointer(backendCVSS)
}

// CVSSOrDefault gets the CVSS parameter from the Ticket struct
func (ji *Issue) CVSSOrDefault() (param float32) {
	return ji.getFloat(backendCVSS)
}

// TicketType gets the TicketType parameter from the Ticket struct
func (ji *Issue) TicketType() (param *string) {
	return &ji.Issue.Fields.Type.Name
}

// TicketTypeOrDefault gets the TicketType parameter from the Ticket struct
func (ji *Issue) TicketTypeOrDefault() (param string) {
	return ji.Issue.Fields.Type.Name
}

// SetProject sets the Project parameter from the Ticket struct
func (ji *Issue) SetProject(val string) {
	ji.Issue.Fields.Project.Name = val
}

// OperatingSystem gets the OperatingSystem parameter from the Ticket struct
func (ji *Issue) OperatingSystem() (param *string) {
	return ji.getStringPointer(backendOperatingSystem)
}

// OperatingSystemOrDefault gets the OperatingSystem parameter from the Ticket struct
func (ji *Issue) OperatingSystemOrDefault() (param string) {
	return ji.getString(backendOperatingSystem)
}

// SetTitle sets the Title parameter from the Ticket struct
func (ji *Issue) SetTitle(val string) {
	ji.Issue.Key = val
}

// SetDeviceID sets the DeviceID parameter from the Ticket struct
func (ji *Issue) SetDeviceID(val string) {
	_ = ji.setCF(ji.connector, backendDeviceID, val)
}

// SetDescription sets the Description parameter from the Ticket struct
func (ji *Issue) SetDescription(val string) {
	ji.Issue.Fields.Description = val
}

// SetIPAddress sets the IPAddress parameter from the Ticket struct
func (ji *Issue) SetIPAddress(val string) {
	_ = ji.setCF(ji.connector, backendIPAddress, val)
}

func (ji *Issue) getStringPointer(key string) (param *string) {
	val := ji.getString(key)
	param = &val
	return param
}

func (ji *Issue) getString(key string) (param string) {
	ret, err := ji.getCF(ji.connector, key)
	if err == nil {

		switch v := ret.(type) {
		case string:
			param = v
		case map[string]interface{}:
			value := v["value"]
			if valueString, ok := value.(string); ok {
				param = valueString
			}

			if len(param) == 0 {
				value = v["name"]
				if valueString, ok := value.(string); ok {
					param = valueString
				}
			}
		case CF:
			if val, ok := v.Value.(string); ok {
				param = val
			}

			if len(param) == 0 {
				if val, ok := v.Name.(string); ok {
					param = val
				}
			}
		}

	}
	return param
}

func (ji *Issue) getTimePointer(key string) (param *time.Time) {
	val := ji.getTime(key)
	param = &val
	return param
}

func (ji *Issue) getTime(key string) (param time.Time) {
	retVal, err := ji.getCF(ji.connector, key)
	if err == nil {
		retValString, ok := retVal.(string)
		if ok {
			paramVal, err := time.Parse(DateFormatJira, retValString+"T00:00:00.000-0700")
			if err == nil {
				param = paramVal
			} else {
				paramVal, err := time.Parse(DateFormatJira, retValString)
				if err == nil {
					param = paramVal
				}
			}
		}
	}
	return param
}

func (ji *Issue) getInt(key string) (param int) {
	ret, err := ji.getCF(ji.connector, key)
	if err == nil {
		retValString, ok := ret.(string)
		if ok {
			param, _ = strconv.Atoi(retValString)
		}

	}
	return param
}

func (ji *Issue) getFloat(key string) (param float32) {
	ret, err := ji.getCF(ji.connector, key)
	if err == nil {

		switch v := ret.(type) {
		case float32:
			param = v
		case float64:
			param = float32(v)
		case string:
			param32, err := strconv.ParseFloat(v, 32)
			if err == nil {
				param = float32(param32)
			}
		case map[string]interface{}:
			value := v["value"]
			if valueString, ok := value.(string); ok {
				param32, err := strconv.ParseFloat(valueString, 32)
				if err == nil {
					param = float32(param32)
				}
			}
		case CF:
			if val, ok := v.Value.(string); ok {
				param32, err := strconv.ParseFloat(val, 32)
				if err == nil {
					param = float32(param32)
				}
			}
		}
	}
	return param
}

func (ji *Issue) getFloatPointer(key string) (param *float32) {
	val := ji.getFloat(key)
	param = &val
	return param
}
