package servicenow

import (
	"fmt"
	"github.com/nortonlifelock/aegis/internal/database/dal"
	"github.com/nortonlifelock/domain"
	"strconv"
	"strings"
	"time"
)

const (
	svcNowDate = "2006-01-02"
)

func mapToDalTickets(connector *SvcNowConnector, svcNowTickets SvcNowTickets) <-chan domain.Ticket {
	var tickets = make(chan domain.Ticket)

	go func() {
		defer close(tickets)

		for index := range svcNowTickets.Results {
			var dalTicket domain.Ticket
			var err error

			ticket := svcNowTickets.Results[index]
			dalTicket, err = mapSvcNowIssueToDalTicket(connector, ticket)
			if err == nil {
				tickets <- dalTicket
			} else {
				break
			}
		}
	}()

	return tickets
}

func mapDalTicketToSvcNowIssue(connector *SvcNowConnector, ticket domain.Ticket) (svcNowTicket *SvcNowRequest, err error) {
	//Not usable fields: approval
	if ticket != nil {
		// Build out the JiraIssue object with the data from the dal object
		svcNowTicket = &SvcNowRequest{}
		var cmdbSysID string
		var vulnSysID string

		if len(*ticket.IPAddress()) > 0 {
			if cmdbSysID, err = connector.getSysIDFor("ip_address", *ticket.IPAddress(), "cmdb_ci"); err == nil {
				svcNowTicket.CmdbCi = cmdbSysID
			}
		}
		if ticket.VulnerabilityID() != "" {
			vulnID := "QID-" + ticket.VulnerabilityID()
			if vulnSysID, err = connector.getSysIDFor("id", vulnID, "sn_vul_entry"); err == nil {
				svcNowTicket.Vulnerability = vulnSysID
			}
		}
		//business service is the imapcted service (splunk as example)
		//svcNowTicket.BusinessService = utilities.Truncate(utilities.RemoveHTMLTags(*ticket.Summary()), 255)

		svcNowTicket.Description = *ticket.Description() // TODO remove html tags

		svcNowTicket.CloseNotes = *ticket.Solution()
		svcNowTicket.AssignedTo = *ticket.AssignedTo()
		svcNowTicket.Requestor = *ticket.ReportedBy()

		//automatic open status upon creation
		//map states and substates in here
		svcNowTicket.State = *ticket.Status()
		//example:
		//if *ticket.Status() == "closedfalsePositive" { svcNowTicket.Substate=2}

		mapVulnInformation(ticket, svcNowTicket)
		err = mapDeviceInformation(ticket, svcNowTicket, connector)
	}
	return svcNowTicket, err
}

func mapVulnInformation(ticket domain.Ticket, svcNowTicket *SvcNowRequest) {
	//we need to find out hwo to populate the vuln in service now, give the vulnID below
	if len(ticket.VulnerabilityID()) > 0 {
		svcNowTicket.AdditionalAssigneeList = ticket.VulnerabilityID()
	}
	if ticket.AlertDate() != nil {
		// Strip the date off of the alert date when storing it
		svcNowTicket.ActivityDue = ticket.AlertDate().Format(svcNowDate)
	}
	if len(*ticket.Priority()) > 0 {
		svcNowTicket.Priority = *ticket.Priority()
	}
	if len(*ticket.CVEReferences()) > 0 {
		svcNowTicket.CorrelationID = *ticket.CVEReferences()
	}
	//it should be mapped to source filed but we cant change it
	if len(*ticket.MethodOfDiscovery()) > 0 {
		svcNowTicket.UserInput = *ticket.MethodOfDiscovery()
	}
	if len(*ticket.VulnerabilityTitle()) > 0 {
		svcNowTicket.WatchList = Truncate(*ticket.VulnerabilityTitle(), 255)
	}
	if ticket.DueDate() != nil {
		// Strip the date off of the due date when storing it
		//Format(JIRADATEONLY) if we need to format the date
		svcNowTicket.DueDate = *ticket.DueDate().Format(svcNowDateFormat)
	}
	if ticket.ResolutionDate() != nil {
		// Strip the date off of the resolution date when storing it
		svcNowTicket.ClosedAt = *ticket.ResolutionDate().Format(svcNowDate)
	}
	if ticket.CVSS() != nil {
		// Reduce the CVSS to one decimal place
		if *ticket.CVSS() >= 0 && *ticket.CVSS() <= 10 {

			svcNowTicket.QualysSeverity = fmt.Sprintf("%.1f", *ticket.CVSS())
		} else {

			svcNowTicket.QualysSeverity = "None"
		}
	}
}

func mapDeviceInformation(ticket domain.Ticket, svcNowTicket *SvcNowRequest, connector *SvcNowConnector) (err error) {
	if len(*ticket.MacAddress()) > 0 {
		svcNowTicket.UAdditionalInformation = Truncate(*ticket.MacAddress(), 250)
	}
	if len(*ticket.ServicePorts()) > 0 {
		svcNowTicket.Port = Truncate(*ticket.ServicePorts(), 250)
	}
	if len(*ticket.IPAddress()) > 0 {
		svcNowTicket.IPAddress = Truncate(*ticket.IPAddress(), 255)

	}
	if len(*ticket.HostName()) > 0 {
		svcNowTicket.DNS = Truncate(*ticket.HostName(), 255)
	}
	if len(ticket.GroupID()) > 0 {
		svcNowTicket.GroupList = ticket.GroupID()
	}
	// lets revist
	if len(ticket.DeviceID()) > 0 {
		svcNowTicket.QualysTicket = ticket.DeviceID()
	}
	// Ensure that the id is a proper database id integer
	if len(ticket.OrganizationID()) > 0 {

		// Get the organization from the database using the id in the ticket object
		var torg domain.Organization
		if torg, err = connector.db.GetOrganizationByID(ticket.OrganizationID()); err == nil {

			// Ensure there is only one return
			if torg != nil {
				svcNowTicket.Skills = torg.Code()
			} else {
				//TODO:
			}
		}
	}
	//we need the assignment group from svcnow like this: SYMC-GSO-Svcs-Splunk
	if len(*ticket.AssignmentGroup()) > 0 {
		svcNowTicket.AssignmentGroup = *ticket.AssignmentGroup()
	}
	//this value either will be setup on the servicenow end or we need to figure
	//how to get the cmdb from svnow , given IP or device id
	if len(*ticket.OperatingSystem()) > 0 {
		//need a mapper
		svcNowTicket.CorrelationDisplay = *ticket.OperatingSystem()

	}
	// we need to use different field , this should be reserved for SVCNow change control
	if ticket.ScanID() > 0 {
		svcNowTicket.QualysAssigneeEmail = strconv.Itoa(int(ticket.ScanID()))
	}
	return err
}

func mapSvcNowIssueToDalTicket(connector *SvcNowConnector, svcNowTicket *Result) (ticket domain.Ticket, err error) {
	ticket = &dal.Ticket{}

	//Not usable fields: approval
	if svcNowTicket != nil {
		reportedBy := svcNowTicket.Requestor.DisplayValue
		title := svcNowTicket.SysID

		description := svcNowTicket.Description

		reportedBy := svcNowTicket.Requestor.DisplayValue
		var status string
		if svcNowTicket.State == "2" || svcNowTicket.State == "Analysis" {
			status = "Resolved-Remediated"
		} else {
			status = svcNowTicket.State
		}
		var score float64
		var cvss float32
		if score, err = strconv.ParseFloat(svcNowTicket.QualysSeverity, 32); err == nil {
			cvss = float32(score)
		}

		var dueDate, resolutionDate, alertDate time.Time
		if svcNowTicket.DueDate != "" {
			dueDate = stringToDalDate(svcNowTicket.DueDate)
		}
		if svcNowTicket.ClosedAt != "" {
			resolutionDate = stringToDalDate(svcNowTicket.ClosedAt)
		}
		if svcNowTicket.ActivityDue != "" {
			alertDate = stringToDalDate(svcNowTicket.ActivityDue)
		}

		err = setTicketsOrg(svcNowTicket, ticket, connector)

		macAddress := svcNowTicket.UAdditionalInformation
		servicePorts := svcNowTicket.Port
		ipAddress := svcNowTicket.IPAddress
		hostName := svcNowTicket.DNS

		var scanID = svcNowTicket.QualysAssigneeEmail // is this right?
		var deviceID = svcNowTicket.QualysTicket      // is this right?
		var groupID = svcNowTicket.GroupList

		operatingSystem := svcNowTicket.CorrelationDisplay
		priority := svcNowTicket.Priority
		cveReferences := svcNowTicket.CorrelationID
		methodOfDiscovery := svcNowTicket.UserInput
		vulnerabilityTitle := svcNowTicket.WatchList
		vulnerabilityID := svcNowTicket.AdditionalAssigneeList
		solution := svcNowTicket.CloseNotes
	}
	return ticket, err
}

func setTicketsOrg(svcNowTicket *Result, ticket domain.Ticket, connector *SvcNowConnector) (err error) {
	// Organization Id
	var org interface{}
	if len(svcNowTicket.Skills) > 0 {

		OrgCode := svcNowTicket.Skills

		// Get the organization from the database using the org code
		var torg domain.Organization
		if len(org.(string)) >= 20 {
			if torg, err = connector.db.GetOrganizationByCode(svcNowTicket.Skills[:20]); err == nil { // TODO this slice falls out of range

				if torg != nil {
					// Set the organization id
					OrganizationID := torg.ID()
				} else {
					// TODO:
				}
			}
		} else {
			//TODO
		}
	} else {
		// TODO:
	}
	return err
}

func setTicketDates(svcNowTicket *Result, ticket domain.Ticket) {
	if svcNowTicket.DueDate != "" {
		DueDate := stringToDalDate(svcNowTicket.DueDate)
	}
	if svcNowTicket.ClosedAt != "" {
		ResolutionDate := stringToDalDate(svcNowTicket.ClosedAt)
	}
	if svcNowTicket.ActivityDue != "" {
		AlertDate := stringToDalDate(svcNowTicket.ActivityDue)
	}
}

func stringToDalDate(dateString string) time.Time {
	returnedDate, _ := time.Parse("2006-01-02T15:04:05.999-0700", strings.Replace(dateString+".000-0700", " ", "T", -1))

	return returnedDate.UTC()
}

// Truncate truncates a string if the length of the first parameter is longer than the value of the second parameter
func Truncate(s string, size int) (truncatedString string) {

	truncatedString = s

	if len(s) > size {
		truncatedString = s[:size]
	}

	return truncatedString
}
