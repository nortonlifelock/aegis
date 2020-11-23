package jira

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/nortonlifelock/aegis/pkg/domain"
	"github.com/pkg/errors"
	"github.com/trivago/tgo/tcontainer"
	"strconv"
	"strings"
)

func (connector *ConnectorJira) mapIssues(issues []Issue) <-chan domain.Ticket {
	var tickets = make(chan domain.Ticket)

	go func() {
		defer close(tickets)

		if issues != nil {
			for index := range issues {
				tickets <- &issues[index]
			}
		} else {
			// this error block is intentionally left blank
		}
	}()

	return tickets
}

func (connector *ConnectorJira) isMappableField(mappableField string) (isMappable bool) {
	for _, field := range connector.payload.MappableFields {
		if strings.ToLower(mappableField) == strings.ToLower(field) {
			isMappable = true
			break
		}
	}

	return isMappable
}

func (connector *ConnectorJira) mapDalTicketToJiraIssue(ticket domain.Ticket) (ji *Issue, err error) {

	if ticket != nil {
		// Build out the Issue object with the data from the dal object
		ji = &Issue{
			Issue:     new(jira.Issue),
			connector: connector,
		}

		err = connector.setJINonCustomFields(ji, ticket)

		if ticket.CVSS() != nil {
			// Reduce the CVSS to one decimal place
			if ford(ticket.CVSS()) >= 0 && ford(ticket.CVSS()) <= 10 {
				_ = ji.setCF(connector, backendCVSS, fmt.Sprintf("%.1f", ford(ticket.CVSS())))
			} else {
				_ = ji.setCF(connector, backendCVSS, "None")
			}
		}

		if len(sord(ticket.ApplicationName())) > 0 {
			setJIField(connector, ji, backendApplicationName, ticket.ApplicationName())
		}

		if len(sord(ticket.TrackingMethod())) > 0 {
			setJIField(connector, ji, backendTrackingMethod, ticket.TrackingMethod())
		}

		if len(ticket.DeviceID()) > 0 {
			_ = ji.setCF(connector, backendDeviceID, ticket.DeviceID())
		}

		if ticket.AlertDate() != nil {
			// Strip the date off of the alert date when storing it
			_ = ji.setCF(connector, backendScanDate, tord(ticket.AlertDate()).Format(DateOnlyJira))
		}

		connector.setJIFieldsRequiringTruncation(ticket, ji)
		setJIField(connector, ji, backendMACAddress, ticket.MacAddress())

		if ticket.ScanID() > 0 {
			setJIField(connector, ji, backendScanID, strconv.Itoa(ticket.ScanID()))
		}

		if len(sord(ticket.Patchable())) > 0 {
			setJIField(connector, ji, backendPatchable, sord(ticket.Patchable()))
		}

		// group ID may not be necessary for vulnerability remediators, but it is useful for tracking CIS tickets by the cloud account
		if len(ticket.GroupID()) > 0 {
			setJIField(connector, ji, backendGroupID, ticket.GroupID())
		}

		if len(sord(ticket.SystemName())) > 0 {
			setJIField(connector, ji, backendSystemName, sord(ticket.SystemName()))
		}

		if ticket.ExceptionDate() != nil {
			setJIField(connector, ji, backendExceptionDate, ticket.ExceptionDate())
		}
		setJIField(connector, ji, backendOWASP, ticket.OWASP())
		setJIField(connector, ji, backendVulnerabilityID, ticket.VulnerabilityID())
		setJIField(connector, ji, backendOperatingSystem, ticket.OperatingSystem())
		setJIField(connector, ji, backendVRRPriority, ticket.Priority()) // TODO: This should default to critical if it's not properly set?
		setJIField(connector, ji, backendCVEReferences, ticket.CVEReferences())
		setJIField(connector, ji, backendVendorReferences, ticket.VendorReferences())
		setJIField(connector, ji, backendMOD, ticket.MethodOfDiscovery())
		setJIField(connector, ji, backendAssignmentGroup, ticket.AssignmentGroup())
		setJIField(connector, ji, backendOSDetailed, ticket.OSDetailed())
		setJIField(connector, ji, backendCategory, ticket.Category())
		setJIField(connector, ji, backendConfig, ticket.Configs())
		setJIField(connector, ji, backendOrg, ticket.OrgCode())
		setJIField(connector, ji, backendSolution, ticket.Solution())

		// TODO need to ensure the org code is in the ticket
		setJIField(connector, ji, backendOrg, sord(ticket.OrgCode()))
	} else {
		err = errors.Errorf("ticket was passed nil to mapDalTicketToJiraIssue")
	}

	return ji, err
}

func (connector *ConnectorJira) setJINonCustomFields(ji *Issue, ticket domain.Ticket) (err error) {
	ji.Issue.Key = ticket.Title()
	ji.Issue.Fields = &jira.IssueFields{}
	ji.Issue.Fields.Unknowns = tcontainer.NewMarshalMap()
	ji.Issue.Fields.Summary = truncate(removeHTMLTags(sord(ticket.Summary())), 255)
	ji.Issue.Fields.Description = removeHTMLTags(sord(ticket.Description()))

	// TODO: Get the project from the JIRA system for this project name
	//issue.Fields.Project = sord(ticket.GetProject())

	// TODO: Get ticket type from JIRA by name
	if len(sord(ticket.TicketType())) > 0 {

		// Get the issue type from the jira system to properly set the pointer
		if len(connector.IssueTypes) > 0 {
			issueType := connector.IssueTypes[sord(ticket.TicketType())]

			// Set the resolution to the proper resolution for this ticket
			ji.Issue.Fields.Type = issueType
		}
	}
	if len(sord(ticket.ResolutionStatus())) > 0 {

		// Get the resolution status from the jira system to properly set the pointer
		if len(connector.Resolutions) > 0 && connector.Resolutions[sord(ticket.ResolutionStatus())] != nil {
			resolution := connector.Resolutions[sord(ticket.ResolutionStatus())]

			// Set the resolution to the proper resolution for this ticket
			ji.Issue.Fields.Resolution = resolution
		}
	}

	if len(sord(ticket.Status())) > 0 {

		// Get the status from the jira system to properly set the pointer
		if len(connector.Statuses) > 0 && connector.Statuses[sord(ticket.Status())] != nil {
			status := connector.Statuses[sord(ticket.Status())]

			// Set the resolution to the proper resolution for this ticket
			ji.Issue.Fields.Status = status
		}
	}

	if len(sord(ticket.Labels())) > 0 {
		labels := strings.Split(sord(ticket.Labels()), ",")
		for index := range labels {
			labels[index] = strings.Replace(labels[index], " ", "-", -1)
		}
		ji.Issue.Fields.Labels = labels
	}

	connector.setJINonCustomDateFields(ticket, ji)
	err = connector.setJINonCustomUserFields(ticket, ji)

	return err
}

func (connector *ConnectorJira) setJINonCustomUserFields(ticket domain.Ticket, ji *Issue) (err error) {
	if len(sord(ticket.AssignedTo())) > 0 {

		var user *jira.User
		if user, _, err = connector.client.User.Get(sord(ticket.AssignedTo())); err == nil {
			ji.Issue.Fields.Assignee = user
		}
	}
	if len(sord(ticket.ReportedBy())) > 0 {

		var user *jira.User
		if user, _, err = connector.client.User.Get(sord(ticket.ReportedBy())); err == nil {
			ji.Issue.Fields.Reporter = user
		}
	}
	return err
}

func (connector *ConnectorJira) setJINonCustomDateFields(ticket domain.Ticket, ji *Issue) {
	if ticket.DueDate() != nil {
		// Strip the date off of the due date when storing it
		ji.Issue.Fields.Duedate = jira.Date(tord(ticket.DueDate()))
	}
	if ticket.CreatedDate() != nil {
		ji.Issue.Fields.Created = jira.Time(tord(ticket.CreatedDate()).Add(-7))
	}
	if ticket.UpdatedDate() != nil {
		ji.Issue.Fields.Updated = jira.Time(tord(ticket.UpdatedDate()).Add(-7))
	}
	if ticket.ResolutionDate() != nil {
		// Strip the date off of the resolution date when storing it
		ji.Issue.Fields.Resolutiondate = jira.Time(tord(ticket.ResolutionDate()))
	}
}

func (connector *ConnectorJira) setJIFieldsRequiringTruncation(ticket domain.Ticket, ji *Issue) {
	if ticket.IPAddress() != nil && len(sord(ticket.IPAddress())) > 0 {
		_ = ji.setCF(connector, backendIPAddress, truncate(sord(ticket.IPAddress()), 255))
	}
	if ticket.HostName() != nil && len(sord(ticket.HostName())) > 0 {
		_ = ji.setCF(connector, backendHostname, truncate(sord(ticket.HostName()), 255))
	}
	if ticket.VulnerabilityTitle() != nil && len(sord(ticket.VulnerabilityTitle())) > 0 {
		_ = ji.setCF(connector, backendVulnerability, truncate(sord(ticket.VulnerabilityTitle()), 255))
	}
	if ticket.ServicePorts() != nil && len(sord(ticket.ServicePorts())) > 0 {
		if sord(ticket.ServicePorts()) != "0" {
			_ = ji.setCF(connector, backendServicePort, truncate(sord(ticket.ServicePorts()), 255))
		}
	}
}

func setJIField(connector *ConnectorJira, ji *Issue, backendField string, value interface{}) {
	if value != nil {
		var set bool
		switch val := value.(type) {

		case *string:
			if val != nil {
				set = len(*val) > 0
			}
		case string:
			set = len(val) > 0
		case *int:
			if val != nil {
				set = *val > 0
			}
		case int:
			set = val > 0
		default:
			set = true
		}

		if set {
			_ = ji.setCF(connector, backendField, value)
		}
	}
}

// truncate truncates a string if the length of the first parameter is longer than the value of the second parameter
func truncate(s string, size int) (truncatedString string) {

	truncatedString = s

	if len(s) > size {
		truncatedString = s[:size]
	}

	return truncatedString
}
