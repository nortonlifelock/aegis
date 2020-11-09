package implementations

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nortonlifelock/aegis/internal/integrations"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/log"
	"sort"
	"strings"
	"time"
)

type CodeRescanJob struct {
	Payload *CodeRescanPayload

	id         string
	orgCode    string
	orgPayload *OrgPayload

	payloadJSON string
	ctx         context.Context
	db          domain.DatabaseConnection
	lstream     log.Logger
	appconfig   domain.Config
	config      domain.JobConfig
	insource    domain.SourceConfig
	outsource   domain.SourceConfig
}

// TODO specify project versions?
type CodeRescanPayload struct {
	ProjectIDs []string `json:"project_ids"`
}

// buildPayload parses the information from the Payload of the job history entry
func (job *CodeRescanJob) buildPayload(pjson string) (err error) {
	// Parse json to RescanPayload
	// Verify pJson length > 0
	job.Payload = &CodeRescanPayload{}
	if len(pjson) > 0 {
		if err = json.Unmarshal([]byte(pjson), job.Payload); err == nil && len(job.Payload.ProjectIDs) == 0 {
			err = fmt.Errorf("no projects_ids included in the payload")
		}
	}

	return err
}

func (job *CodeRescanJob) Process(ctx context.Context, id string, appconfig domain.Config, db domain.DatabaseConnection, lstream log.Logger, payload string, jobConfig domain.JobConfig, inSource []domain.SourceConfig, outSource []domain.SourceConfig) (err error) {
	var ok bool
	if job.ctx, job.id, job.appconfig, job.db, job.lstream, job.payloadJSON, job.config, job.insource, job.outsource, ok = validInputs(ctx, id, appconfig, db, lstream, payload, jobConfig, inSource, outSource); ok {
		if err = job.buildPayload(job.payloadJSON); err == nil {

			if job.orgCode, job.orgPayload, err = getOrgInfo(job.db, job.config.OrganizationID()); err == nil {
				sort.Sort(job.orgPayload)

				var engine integrations.TicketingEngine
				if engine, err = integrations.GetEngine(job.ctx, job.outsource.Source(), job.db, job.lstream, job.appconfig, job.outsource); err == nil {

					var scanner integrations.CodeScanner
					if scanner, err = integrations.GetCodeScanner(job.insource.Source(), job.db, job.insource, job.appconfig, job.lstream); err == nil {
						for _, projectID := range job.Payload.ProjectIDs {
							if len(projectID) > 0 {
								err = job.processCodeImageFindings(engine, scanner, projectID)
								if err != nil {
									job.lstream.Send(log.Errorf(err, "error while processing [%s]", projectID))
								}
							} else {
								job.lstream.Send(log.Errorf(err, "empty project ID found in the payload"))
							}
						}
					}
				}
			}
		}
	} else {
		err = fmt.Errorf("input validation failed")
	}

	return err
}

func (job *CodeRescanJob) processCodeImageFindings(engine integrations.TicketingEngine, scanner integrations.CodeScanner, projectID string) (err error) {
	var tickets <-chan domain.Ticket
	tickets, err = engine.GetOpenTicketsByGroupID(job.insource.Source(), job.orgCode, projectID)
	if err == nil {
		var findings []domain.CodeFinding
		if findings, err = scanner.GetProjectVulnerabilities(job.ctx, projectID); err == nil {
			var findingsAsTickets = make([]domain.Ticket, 0)
			for index := range findings {
				finding := findings[index]

				priority, dueDate := calculateSeverityAccordingToOrgPayload(job.orgPayload, finding.CVSS())

				findingTic := &CodeFinding{
					finding,
					job.insource.Source(),
					job.orgCode,
					dueDate,
					priority,
				}

				if len(findingTic.DeviceID()) > 0 && len(findingTic.VulnerabilityID()) > 0 {
					findingsAsTickets = append(findingsAsTickets, findingTic)
				}
			}

			ticketSlice := fanInChannel(tickets)
			processFindingsAndTickets(
				job.lstream,
				job.db,
				job.config.OrganizationID(),
				job.insource.SourceID(),
				engine,
				ticketSlice,
				findingsAsTickets,
				fmt.Sprintf("finding was NOT found by %s", job.insource.Source()),
				fmt.Sprintf("finding still detected by %s", job.insource.Source()),
				func(ticket domain.Ticket) string {
					return fmt.Sprintf("%s;%s;%s", ticket.DeviceID(), ticket.VulnerabilityID(), sord(ticket.ServicePorts()))
				},
				func(ticket domain.Ticket) bool {
					return ford(ticket.CVSS()) >= job.orgPayload.LowestCVSS
				},
			)
		}
	}

	return err
}

type CodeFinding struct {
	finding  domain.CodeFinding
	mod      string
	orgCode  string
	dueDate  time.Time
	priority string
}

func (c *CodeFinding) HubProjectName() (param *string) {
	val := c.finding.ProjectName()
	return &val
}

func (c *CodeFinding) HubProjectVersion() (param *string) {
	val := c.finding.ProjectVersion()
	return &val
}

func (c *CodeFinding) HubSeverity() (param *string) {
	// Medium: If at least one vulnerability is of medium risk OR policy severity is 'Major'
	// High : If at least one vulnerability is of high risk OR Policy Severity is 'Critical'
	var val string
	if strings.ToLower(sord(c.Priority())) == "high" || strings.ToLower(sord(c.Priority())) == "critical" || strings.Contains(strings.ToLower(sord(c.PolicySeverity())), "critical") {
		val = "High"
	} else if strings.ToLower(sord(c.Priority())) == "medium" || strings.Contains(strings.ToLower(sord(c.PolicySeverity())), "major") {
		val = "Medium"
	}
	return &val
}

func (c *CodeFinding) ComponentName() (param *string) {
	val := c.finding.ComponentName()
	return &val
}

func (c *CodeFinding) ComponentVersion() (param *string) {
	val := c.finding.ComponentVersion()
	return &val
}

func (c *CodeFinding) PolicyRule() (param *string) {
	val := c.finding.ViolatedPolicyName()
	return &val
}

func (c *CodeFinding) PolicySeverity() (param *string) {
	val := c.finding.ViolatedPolicySeverity()
	return &val
}

func (c *CodeFinding) AlertDate() (param *time.Time) {
	return
}

func (c *CodeFinding) AssignedTo() (param *string) {
	return
}

func (c *CodeFinding) AssignmentGroup() (param *string) {
	return
}

func (c *CodeFinding) Category() (param *string) {
	return
}

func (c *CodeFinding) CERF() (param string) {
	return
}

func (c *CodeFinding) ExceptionExpiration() (param time.Time) {
	return
}

func (c *CodeFinding) CVEReferences() (param *string) {
	return
}

func (c *CodeFinding) CVSS() (param *float32) {
	val := c.finding.CVSS()
	return &val
}

func (c *CodeFinding) CloudID() (param string) {
	return
}

func (c *CodeFinding) Configs() (param string) {
	return
}

func (c *CodeFinding) CreatedDate() (param *time.Time) {
	return
}

func (c *CodeFinding) DBCreatedDate() (param time.Time) {
	return
}

func (c *CodeFinding) DBUpdatedDate() (param *time.Time) {
	return
}

func (c *CodeFinding) Description() (param *string) {
	val := c.finding.Description()
	return &val
}

func (c *CodeFinding) DeviceID() (param string) {
	return fmt.Sprintf("%s %s", c.finding.ComponentName(), c.finding.ComponentVersion())
}

func (c *CodeFinding) DueDate() (param *time.Time) {
	return
}

func (c *CodeFinding) ExceptionDate() (param *time.Time) {
	return
}

func (c *CodeFinding) GroupID() string {
	return c.finding.ProjectID()
}

func (c *CodeFinding) HostName() (param *string) {
	return
}

func (c *CodeFinding) ID() (param int) {
	return
}

func (c *CodeFinding) IPAddress() (param *string) {
	return
}

func (c *CodeFinding) Labels() (param *string) {
	return
}

func (c *CodeFinding) LastChecked() (param *time.Time) {
	val := c.finding.Updated()
	return &val
}

func (c *CodeFinding) MacAddress() (param *string) {
	return
}

func (c *CodeFinding) MethodOfDiscovery() (param *string) {
	val := c.mod
	return &val
}

func (c *CodeFinding) OSDetailed() (param *string) {
	return
}

func (c *CodeFinding) OperatingSystem() (param *string) {
	return
}

func (c *CodeFinding) OrgCode() (param *string) {
	val := c.orgCode
	return &val
}

func (c *CodeFinding) OrganizationID() (param string) {
	return
}

func (c *CodeFinding) OWASP() (param *string) {
	return
}

func (c *CodeFinding) Patchable() (param *string) {
	return
}

func (c *CodeFinding) Priority() (param *string) {
	val := c.priority
	return &val
}

func (c *CodeFinding) Project() (param *string) {
	return
}

func (c *CodeFinding) ReportedBy() (param *string) {
	projectOwner := c.finding.ProjectOwner() // reporter?
	projectOwner = ""                        // TODO remove?
	return &projectOwner
}

func (c *CodeFinding) ResolutionDate() (param *time.Time) {
	return
}

func (c *CodeFinding) ResolutionStatus() (param *string) {
	return
}

func (c *CodeFinding) ScanID() (param int) {
	return
}

func (c *CodeFinding) ServicePorts() (param *string) {
	return
}

func (c *CodeFinding) Solution() (param *string) {
	return
}

func (c *CodeFinding) Status() (param *string) {
	return
}

func (c *CodeFinding) Summary() (param *string) {
	val := c.finding.Summary()
	return &val
}

func (c *CodeFinding) SystemName() (param *string) {
	return
}

func (c *CodeFinding) TicketType() (param *string) {
	val := "Request"
	return &val
}

func (c *CodeFinding) Title() (param string) {
	return
}

func (c *CodeFinding) UpdatedDate() (param *time.Time) {
	return
}

func (c *CodeFinding) VendorReferences() (param *string) {
	return
}

func (c *CodeFinding) VulnerabilityID() (param string) {
	return c.finding.VulnerabilityID()
}

func (c *CodeFinding) VulnerabilityTitle() (param *string) {
	return
}

func (c *CodeFinding) ApplicationName() (param *string) {
	return
}
