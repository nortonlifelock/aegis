package jira

import (
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

// UpdatedDate gets the UpdatedDate parameter from the Ticket struct
func (ji *Issue) UpdatedDate() (param *time.Time) {
	param = new(time.Time)
	*param = time.Time(ji.Issue.Fields.Updated)
	return param
}

// Solution gets the Solution parameter from the Ticket struct
func (ji *Issue) Solution() (param *string) {
	return ji.getStringPointer(backendSolution)
}

// Labels gets the Labels parameter from the Ticket struct
func (ji *Issue) Labels() (param *string) {
	if ji.Issue.Fields.Labels != nil && len(ji.Issue.Fields.Labels) > 0 {
		var ret = strings.Join(ji.Issue.Fields.Labels, ",")
		param = &ret
	}

	return param
}

// AlertDate gets the AlertDate parameter from the Ticket struct
func (ji *Issue) AlertDate() (param *time.Time) {
	return ji.getTimePointer(backendScanDate)
}

// DeviceID gets the DeviceID parameter from the Ticket struct
func (ji *Issue) DeviceID() (param string) {
	return ji.getString(backendDeviceID)
}

func (ji *Issue) SystemName() (param *string) {
	return ji.getStringPointer(backendSystemName)
}

func (ji *Issue) Patchable() (param *string) {
	return ji.getStringPointer(backendPatchable)
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

// CreatedDate gets the CreatedDate parameter from the Ticket struct
func (ji *Issue) CreatedDate() (param *time.Time) {
	param = new(time.Time)
	*param = time.Time(ji.Issue.Fields.Created)
	return param
}

// VendorReferences gets the VendorReferences parameter from the Ticket struct
func (ji *Issue) VendorReferences() (param *string) {
	return ji.getStringPointer(backendVendorReferences)
}

// AssignedTo gets the AssignedTo parameter from the Ticket struct
func (ji *Issue) AssignedTo() (param *string) {
	if ji.Issue.Fields.Assignee != nil {
		param = &ji.Issue.Fields.Assignee.Name
	}
	return param
}

// OrganizationID gets the OrganizationID parameter from the Ticket struct
// The OrganizationID is a database field, and is not relevant to JIRA tickets
func (ji *Issue) OrganizationID() (param string) {
	return param
}

// IPAddress gets the IPAddress parameter from the Ticket struct
func (ji *Issue) IPAddress() (param *string) {
	return ji.getStringPointer(backendIPAddress)
}

// CERF gets the CERF parameter from the Ticket struct
func (ji *Issue) CERF() (param string) {
	return ji.getString(backendCERF)
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

// OrgCode gets the OrgCode parameter from the Ticket struct
func (ji *Issue) OrgCode() (param *string) {
	return ji.getStringPointer(backendOrg)
}

// DBUpdatedDate gets the DBUpdatedDate parameter from the Ticket struct
func (ji *Issue) DBUpdatedDate() (param *time.Time) {
	return // jira will not have this information
}

// Description gets the Description parameter from the Ticket struct
func (ji *Issue) Description() (param *string) {
	return &ji.Issue.Fields.Description
}

// Configs gets the Configs parameter from the Ticket struct
func (ji *Issue) Configs() (param string) {
	return ji.getString(backendConfig)
}

// LastChecked gets the LastChecked parameter from the Ticket struct
func (ji *Issue) LastChecked() (param *time.Time) {
	return ji.getTimePointer(backendLastChecked)
}

// Project gets the Project parameter from the Ticket struct
func (ji *Issue) Project() (param *string) {
	return &ji.Issue.Fields.Project.Name
}

// AssignmentGroup gets the AssignmentGroup parameter from the Ticket struct
func (ji *Issue) AssignmentGroup() (param *string) {
	return ji.getStringPointer(backendAssignmentGroup)
}

// GroupID gets the GroupID parameter from the Ticket struct
func (ji *Issue) GroupID() (param string) {
	return ji.getString(backendGroupID)
}

// HostName gets the HostName parameter from the Ticket struct
func (ji *Issue) HostName() (param *string) {
	return ji.getStringPointer(backendHostname)
}

// ResolutionDate gets the ResolutionDate parameter from the Ticket struct
func (ji *Issue) ResolutionDate() (param *time.Time) {
	param = new(time.Time)
	*param = time.Time(ji.Issue.Fields.Resolutiondate)
	if param.IsZero() {
		param = ji.getTimePointer("Resolution Date") // TODO
	}
	return param
}

// VulnerabilityTitle gets the VulnerabilityTitle parameter from the Ticket struct
func (ji *Issue) VulnerabilityTitle() (param *string) {
	return ji.getStringPointer(backendVulnerability)
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

// Priority gets the Priority parameter from the Ticket struct
func (ji *Issue) Priority() (param *string) {
	return ji.getStringPointer(backendVRRPriority)
}

// Status gets the Status parameter from the Ticket struct
func (ji *Issue) Status() (param *string) {
	if ji.Issue.Fields.Status != nil {
		param = &ji.Issue.Fields.Status.Name
	}

	return
}

// OSDetailed gets the OSDetailed parameter from the Ticket struct
func (ji *Issue) OSDetailed() (param *string) {
	return ji.getStringPointer(backendOSDetailed)
}

func (ji *Issue) Category() (param *string) {
	return ji.getStringPointer(backendCategory)
}

// MethodOfDiscovery gets the MethodOfDiscovery parameter from the Ticket struct
func (ji *Issue) MethodOfDiscovery() (param *string) {
	return ji.getStringPointer(backendMOD)
}

// CVEReferences gets the CVEReferences parameter from the Ticket struct
func (ji *Issue) CVEReferences() (param *string) {
	return ji.getStringPointer(backendCVEReferences)
}

// Title gets the Title parameter from the Ticket struct
func (ji *Issue) Title() (param string) {
	return ji.Issue.Key
}

// ServicePorts gets the ServicePorts parameter from the Ticket struct
func (ji *Issue) ServicePorts() (param *string) {
	return ji.getStringPointer(backendServicePort)
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

// MacAddress gets the MacAddress parameter from the Ticket struct
func (ji *Issue) MacAddress() (param *string) {
	return ji.getStringPointer(backendMACAddress)
}

// Summary gets the Summary parameter from the Ticket struct
func (ji *Issue) Summary() (param *string) {
	return &ji.Issue.Fields.Summary
}

// DueDate gets the DueDate parameter from the Ticket struct
func (ji *Issue) DueDate() (param *time.Time) {
	var retVal = time.Time(ji.Issue.Fields.Duedate)
	return &retVal
}

// CVSS gets the CVSS parameter from the Ticket struct
func (ji *Issue) CVSS() (param *float32) {
	return ji.getFloatPointer(backendCVSS)
}

// TicketType gets the TicketType parameter from the Ticket struct
func (ji *Issue) TicketType() (param *string) {
	return &ji.Issue.Fields.Type.Name
}

// OperatingSystem gets the OperatingSystem parameter from the Ticket struct
func (ji *Issue) OperatingSystem() (param *string) {
	return ji.getStringPointer(backendOperatingSystem)
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
