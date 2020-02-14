package tool

import (
	"fmt"
	"github.com/nortonlifelock/domain"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type csvTicket struct {
	domain.Ticket
	command     string
	updateLine  []string
	descToIndex map[string]int
}

func (t csvTicket) AlertDate() (param *time.Time) {
	key := "alert date"

	if keyIsInMap(key, t.descToIndex) && len(t.updateLine[t.descToIndex[key]]) > 0 {
		if t.command == Create {
			alertDate, err := time.Parse(TimeLayout, t.updateLine[t.descToIndex[key]])
			if err == nil {
				param = &alertDate
			} else {
				fmt.Println(fmt.Sprintf("Error while parsing alert date: %s", err.Error()))
			}
		}
	} else {
		param = t.Ticket.AlertDate()
	}

	return param
}

func (t csvTicket) AssignedTo() (param *string) {
	key := "assignee"

	if keyIsInMap(key, t.descToIndex) && len(t.updateLine[t.descToIndex[key]]) > 0 {
		param = &t.updateLine[t.descToIndex[key]]
	} else {
		param = t.Ticket.AssignedTo()
	}

	return param
}

func (t csvTicket) AssignmentGroup() (param *string) {
	key := "assignment group"

	if keyIsInMap(key, t.descToIndex) && len(t.updateLine[t.descToIndex[key]]) > 0 {
		param = &t.updateLine[t.descToIndex[key]]
	} else {
		param = t.Ticket.AssignmentGroup()
	}

	return param
}

func (t csvTicket) CERF() (param string) {
	key := "cerf"

	if keyIsInMap(key, t.descToIndex) && len(t.updateLine[t.descToIndex[key]]) > 0 {
		var cerf = t.updateLine[t.descToIndex[key]]
		var desiredString = "CERF-"
		var prefixIndex = strings.Index(cerf, desiredString)
		if prefixIndex > -1 {
			prefixIndex += len(desiredString)
			for prefixIndex < len(cerf) && unicode.IsNumber(rune(cerf[prefixIndex])) {
				prefixIndex++
			}
			if prefixIndex <= len(cerf) {
				var cerfLink = cerf[0:prefixIndex]
				param = cerfLink
			}
		}
	} else {
		param = t.Ticket.CERF()
	}

	return param
}

func (t csvTicket) CERFExpirationDate() (param time.Time) {
	return t.Ticket.CERFExpirationDate()
}

func (t csvTicket) CVEReferences() (param *string) {
	key := "cve references"

	if keyIsInMap(key, t.descToIndex) && len(t.updateLine[t.descToIndex[key]]) > 0 {
		param = &t.updateLine[t.descToIndex[key]]
	} else {
		param = t.Ticket.CVEReferences()
	}

	return param
}

func (t csvTicket) CVSS() (param *float32) {
	key := "cvss"

	if keyIsInMap(key, t.descToIndex) && len(t.updateLine[t.descToIndex[key]]) > 0 {
		floatValue, err := strconv.ParseFloat(t.updateLine[t.descToIndex[key]], 32)
		if err == nil {
			val := float32(floatValue)
			param = &val
		} else {
			fmt.Println(fmt.Sprintf("Error while parsing float for cvss: %s", err.Error()))
		}
	} else {
		param = t.Ticket.CVSS()
	}

	return param
}

func (t csvTicket) CloudID() (param string) {
	return t.Ticket.CloudID()
}

func (t csvTicket) Configs() (param string) {
	return t.Ticket.Configs()
}

func (t csvTicket) CreatedDate() (param *time.Time) {
	return t.Ticket.CreatedDate()
}

func (t csvTicket) DBCreatedDate() (param time.Time) {
	return t.Ticket.DBCreatedDate()
}

func (t csvTicket) Description() (param *string) {
	key := "description"

	if keyIsInMap(key, t.descToIndex) && len(t.updateLine[t.descToIndex[key]]) > 0 {
		param = &t.updateLine[t.descToIndex[key]]
	} else {
		param = t.Ticket.Description()
	}

	return param
}

func (t csvTicket) DeviceID() (param string) {
	return t.Ticket.DeviceID()
}

func (t csvTicket) DueDate() (param *time.Time) {
	key := "due date"

	if keyIsInMap(key, t.descToIndex) && len(t.updateLine[t.descToIndex[key]]) > 0 {
		dueDate, err := time.Parse(TimeLayout, t.updateLine[t.descToIndex[key]])
		if err == nil {
			param = &dueDate
		} else {
			fmt.Println(fmt.Sprintf("Error while parsing due date: %s", err.Error()))
		}
	} else {
		param = t.Ticket.DueDate()
	}

	return param
}

func (t csvTicket) HostName() (param *string) {
	key := "hostname"

	if keyIsInMap(key, t.descToIndex) && len(t.updateLine[t.descToIndex[key]]) > 0 {
		val := t.updateLine[t.descToIndex[key]]
		param = &val
	} else {
		param = t.Ticket.HostName()
	}

	return param
}

func (t csvTicket) ID() (param int) {
	return t.Ticket.ID()
}

func (t csvTicket) IPAddress() (param *string) {
	key := "ip address"

	if keyIsInMap(key, t.descToIndex) && len(t.updateLine[t.descToIndex[key]]) > 0 {
		val := t.updateLine[t.descToIndex[key]]
		param = &val
	} else {
		param = t.Ticket.IPAddress()
	}

	return param
}

func (t csvTicket) Labels() (param *string) {
	key := "labels"

	if keyIsInMap(key, t.descToIndex) && len(t.updateLine[t.descToIndex[key]]) > 0 {
		val := t.updateLine[t.descToIndex[key]]
		param = &val
	} else {
		param = t.Ticket.Labels()
	}

	return param
}

func (t csvTicket) LastChecked() (param *time.Time) {
	return t.Ticket.LastChecked()
}

func (t csvTicket) MacAddress() (param *string) {
	return t.Ticket.MacAddress()
}

func (t csvTicket) MethodOfDiscovery() (param *string) {
	key := "method of discovery"

	if keyIsInMap(key, t.descToIndex) && len(t.updateLine[t.descToIndex[key]]) > 0 {
		val := t.updateLine[t.descToIndex[key]]
		param = &val
	} else {
		param = t.Ticket.MethodOfDiscovery()
	}

	return param
}

func (t csvTicket) OSDetailed() (param *string) {
	return t.Ticket.OSDetailed()
}

func (t csvTicket) OperatingSystem() (param *string) {
	return t.Ticket.OperatingSystem()
}

func (t csvTicket) OrgCode() (param *string) {
	return t.Ticket.OrgCode()
}

func (t csvTicket) OrganizationID() (param string) {
	return t.Ticket.OrganizationID()
}

func (t csvTicket) Priority() (param *string) {
	key := "priority"

	if keyIsInMap(key, t.descToIndex) && len(t.updateLine[t.descToIndex[key]]) > 0 {
		val := t.updateLine[t.descToIndex[key]]
		param = &val
	} else {
		param = t.Ticket.Priority()
	}

	return param
}

func (t csvTicket) Project() (param *string) {
	return t.Ticket.Project()
}

func (t csvTicket) ReportedBy() (param *string) {
	return t.Ticket.ReportedBy()
}

func (t csvTicket) ResolutionDate() (param *time.Time) {
	key := "resolution date"

	if keyIsInMap(key, t.descToIndex) && len(t.updateLine[t.descToIndex[key]]) > 0 {
		resolutionDate, err := time.Parse(TimeLayout, t.updateLine[t.descToIndex[key]])
		if err == nil {
			param = &resolutionDate
		} else {
			fmt.Println(fmt.Sprintf("Error while parsing resolution date: %s", err.Error()))
		}
	} else {
		param = t.Ticket.ResolutionDate()
	}

	return param
}

func (t csvTicket) ResolutionStatus() (param *string) {
	return t.Ticket.ResolutionStatus()
}

func (t csvTicket) ScanID() (param int) {
	return t.Ticket.ScanID()
}

func (t csvTicket) ServicePorts() (param *string) {
	key := "service ports"

	if keyIsInMap(key, t.descToIndex) && len(t.updateLine[t.descToIndex[key]]) > 0 {
		val := t.updateLine[t.descToIndex[key]]
		param = &val
	} else {
		param = t.Ticket.ServicePorts()
	}

	return param
}

func (t csvTicket) Solution() (param *string) {
	key := "solution"

	if keyIsInMap(key, t.descToIndex) && len(t.updateLine[t.descToIndex[key]]) > 0 {
		val := t.updateLine[t.descToIndex[key]]
		param = &val
	} else {
		param = t.Ticket.Solution()
	}

	return param
}

func (t csvTicket) Status() (param *string) {
	key := "status"

	if keyIsInMap(key, t.descToIndex) && len(t.updateLine[t.descToIndex[key]]) > 0 {
		val := t.updateLine[t.descToIndex[key]]
		param = &val
	} else {
		param = t.Ticket.Status()
	}

	return param
}

func (t csvTicket) Summary() (param *string) {
	key := "summary"

	if keyIsInMap(key, t.descToIndex) && len(t.updateLine[t.descToIndex[key]]) > 0 {
		val := t.updateLine[t.descToIndex[key]]
		param = &val
	} else {
		param = t.Ticket.Summary()
	}

	return param
}

func (t csvTicket) TicketType() (param *string) {
	val := "Request"
	return &val
}

func (t csvTicket) Title() (param string) {
	key := "title"

	if keyIsInMap(key, t.descToIndex) && len(t.updateLine[t.descToIndex[key]]) > 0 {
		param = t.updateLine[t.descToIndex[key]]
	} else {
		param = t.Ticket.Title()
	}

	return param
}

func (t csvTicket) UpdatedDate() (param *time.Time) {
	return t.Ticket.UpdatedDate()
}

func (t csvTicket) VendorReferences() (param *string) {
	key := "vendor references"

	if keyIsInMap(key, t.descToIndex) && len(t.updateLine[t.descToIndex[key]]) > 0 {
		val := t.updateLine[t.descToIndex[key]]
		param = &val
	} else {
		param = t.Ticket.VendorReferences()
	}

	return param
}

func (t csvTicket) VulnerabilityID() (param string) {
	key := "vulnerabilityid"

	if keyIsInMap(key, t.descToIndex) && len(t.updateLine[t.descToIndex[key]]) > 0 {
		param = t.updateLine[t.descToIndex[key]]
	} else {
		param = t.Ticket.VulnerabilityID()
	}

	return param
}

func (t csvTicket) VulnerabilityTitle() (param *string) {
	key := "vulnerability title"

	if keyIsInMap(key, t.descToIndex) && len(t.updateLine[t.descToIndex[key]]) > 0 {
		val := t.updateLine[t.descToIndex[key]]
		param = &val
	} else {
		param = t.Ticket.VulnerabilityTitle()
	}

	return param
}

func keyIsInMap(key string, vals map[string]int) bool {
	var seen bool
	for possKey := range vals {
		if key == possKey {
			seen = true
			break
		}
	}

	return seen
}
