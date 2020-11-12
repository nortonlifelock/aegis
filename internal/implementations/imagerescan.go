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

type ImageRescanJob struct {
	Payload *ImageRescanPayload

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

type ImageRescanPayload struct {
	RegistryImage     []string `json:"registry_image"`
	ExceptionAssignee string   `json:"exception_assignee"`
}

// buildPayload parses the information from the Payload of the job history entry
func (job *ImageRescanJob) buildPayload(pjson string) (err error) {
	// Parse json to RescanPayload
	// Verify pJson length > 0
	job.Payload = &ImageRescanPayload{}
	if len(pjson) > 0 {
		if err = json.Unmarshal([]byte(pjson), job.Payload); err == nil && len(job.Payload.RegistryImage) == 0 {
			err = fmt.Errorf("no registry_image combinations included in the payload")
		}
	}

	return err
}

func (job *ImageRescanJob) Process(ctx context.Context, id string, appconfig domain.Config, db domain.DatabaseConnection, lstream log.Logger, payload string, jobConfig domain.JobConfig, inSource []domain.SourceConfig, outSource []domain.SourceConfig) (err error) {
	var ok bool
	if job.ctx, job.id, job.appconfig, job.db, job.lstream, job.payloadJSON, job.config, job.insource, job.outsource, ok = validInputs(ctx, id, appconfig, db, lstream, payload, jobConfig, inSource, outSource); ok {
		if err = job.buildPayload(job.payloadJSON); err == nil {

			if job.orgCode, job.orgPayload, err = getOrgInfo(job.db, job.config.OrganizationID()); err == nil {
				sort.Sort(job.orgPayload)

				var engine integrations.TicketingEngine
				if engine, err = integrations.GetEngine(job.ctx, job.outsource.Source(), job.db, job.lstream, job.appconfig, job.outsource); err == nil {

					var scanner integrations.IScanner
					if scanner, err = integrations.GetImageScanner(job.insource.Source(), job.db, job.insource, job.appconfig, job.lstream); err == nil {
						for _, registryImage := range job.Payload.RegistryImage {
							if len(registryImage) > 0 {
								err = job.processImageFindings(engine, scanner, registryImage)
								if err != nil {
									job.lstream.Send(log.Errorf(err, "error while processing [%s]", registryImage))
								}
							} else {
								job.lstream.Send(log.Errorf(err, "empty registry;image found in the payload"))
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

func (job *ImageRescanJob) processImageFindings(engine integrations.TicketingEngine, scanner integrations.IScanner, registryImage string) (err error) {
	if len(strings.Split(registryImage, ";")) == 2 {
		var registry = strings.Split(registryImage, ";")[0]
		var image = strings.Split(registryImage, ";")[1]

		var tickets <-chan domain.Ticket
		tickets, err = engine.GetOpenTicketsByGroupID(job.insource.Source(), job.orgCode, registry)
		if err == nil {
			var findings []domain.ImageFinding
			if findings, err = scanner.RescanImage(job.ctx, image, registry); err == nil {
				var findingsAsTickets = make([]domain.Ticket, 0)
				for index := range findings {
					finding := findings[index]

					var priority string
					var dueDate time.Time
					if priority, dueDate, err = calculateSLAForImageFinding(finding, job.orgPayload); err == nil {
						findingTic := &ImageFinding{
							finding,
							job.insource.Source(),
							job.orgCode,
							dueDate,
							priority,
						}

						if len(findingTic.DeviceID()) > 0 && len(findingTic.VulnerabilityID()) > 0 {
							findingsAsTickets = append(findingsAsTickets, findingTic)
						}
					} else {
						job.lstream.Send(log.Errorf(err, "error while calculating priority for ticket"))
					}
				}

				findingMap := mapImageFindingsByDeviceIDVulnID(findings)
				ticketSlice := fanInChannel(tickets)
				ticketsForImage := make([]domain.Ticket, 0)
				for index := range ticketSlice {
					if ticketSlice[index].DeviceID() == image {
						ticketsForImage = append(ticketsForImage, ticketSlice[index])
					}
				}

				_, _, ticketsWithFindings := processFindingsAndTickets(
					job.lstream,
					job.db,
					job.config.OrganizationID(),
					job.insource.SourceID(),
					engine,
					ticketsForImage,
					findingsAsTickets,
					fmt.Sprintf("finding was NOT found by %s", job.insource.Source()),
					fmt.Sprintf("finding still detected by %s", job.insource.Source()),
					func(ticket domain.Ticket) string {
						return fmt.Sprintf("%s;%s;%s", ticket.DeviceID(), ticket.VulnerabilityID(), sord(ticket.ServicePorts()))
					},
					func(ticket domain.Ticket) bool {
						finding := findingMap[fmt.Sprintf("%s;%s;%s", ticket.DeviceID(), ticket.VulnerabilityID(), sord(ticket.ServicePorts()))]
						if finding != nil {
							if finding.Exception() {
								return false
							}
						}
						return ford(ticket.CVSS()) >= job.orgPayload.LowestCVSS
					},
				)

				for _, pair := range ticketsWithFindings {
					key := fmt.Sprintf("%s;%s;%s", pair.ticket.DeviceID(), pair.ticket.VulnerabilityID(), sord(pair.ticket.ServicePorts()))
					finding := findingMap[key]
					if finding != nil {
						if sord(pair.ticket.Status()) == engine.GetStatusMap(domain.StatusClosedException) {
							if !finding.Exception() { // if the finding isn't already marked as an exception in Aqua, set it as an exception in Aqua
								err = scanner.CreateException(finding, fmt.Sprintf("%s marked as Closed-Exception on %s", pair.ticket.Title(), time.Now().Format(time.RFC3339)))
								if err != nil {
									job.lstream.Send(log.Errorf(err, "error marking ticket as an exception in Aqua [%s]", pair.ticket.Title()))
								}
							}
						} else if finding.Exception() {
							// close the ticket
							err = engine.Transition(pair.ticket, engine.GetStatusMap(domain.StatusClosedException), fmt.Sprintf("Moving ticket to closed exception as it was acknowledged in %s", job.insource.Source()), job.Payload.ExceptionAssignee)
							if err != nil {
								job.lstream.Send(log.Errorf(err, "error while setting [%s] to Closed-Exception", pair.ticket.Title()))
							}
						}
					} else {
						job.lstream.Send(log.Errorf(err, "could not find finding [%s] while checking for exceptions", key))
					}
				}
			}
		}
	} else {
		err = fmt.Errorf("registry_image appears to be malformed, expected it to be in the form [REGISTRY;IMAGE]")
	}

	return err
}

func calculateSLAForImageFinding(finding domain.ImageFinding, orgPayload *OrgPayload) (highestApplicableSeverity string, dueDate time.Time, err error) {
	var cvssScore float32
	if orgPayload.CVSSVersion == cvssVersion2 && finding.CVSS2() != nil {
		cvssScore = ford(finding.CVSS2())
	} else if orgPayload.CVSSVersion == cvssVersion3 && finding.CVSS3() != nil {
		cvssScore = ford(finding.CVSS3())
	} else {
		err = fmt.Errorf("either CVSS version not present [%d] or CVSS score not present [%v|%v]", orgPayload.CVSSVersion, finding.CVSS2(), finding.CVSS3())
	}

	// we iterate over the sorted list of custom severity ranges and find the highest applicable severity
	for index := range orgPayload.Severities {
		if cvssScore >= orgPayload.Severities[index].CVSSMin {
			highestApplicableSeverity = orgPayload.Severities[index].Name
			dueDate = time.Now().AddDate(0, 0, orgPayload.Severities[index].Duration)
		}
	}

	return highestApplicableSeverity, dueDate, err
}

type ImageFinding struct {
	finding  domain.ImageFinding
	mod      string
	orgCode  string
	dueDate  time.Time
	priority string
}

func (i *ImageFinding) AlertDate() (param *time.Time) {
	return
}

func (i *ImageFinding) Category() (param *string) {
	return
}

func (i *ImageFinding) AssignedTo() (param *string) {
	return
}

func (i *ImageFinding) AssignmentGroup() (param *string) {
	return
}

func (i *ImageFinding) CERF() (param string) {
	return
}

func (i *ImageFinding) ExceptionExpiration() (param time.Time) {
	return
}

func (i *ImageFinding) CVEReferences() (param *string) {
	return
}

func (i *ImageFinding) CVSS() (param *float32) {
	return i.finding.CVSS2()
}

func (i *ImageFinding) CloudID() (param string) {
	return
}

func (i *ImageFinding) Configs() (param string) {
	return
}

func (i *ImageFinding) CreatedDate() (param *time.Time) {
	return
}

func (i *ImageFinding) DBCreatedDate() (param time.Time) {
	return
}

func (i *ImageFinding) DBUpdatedDate() (param *time.Time) {
	return
}

func (i *ImageFinding) Description() (param *string) {
	return i.finding.Summary()
}

func (i *ImageFinding) DeviceID() (param string) {
	return i.finding.ImageName()
}

func (i *ImageFinding) DueDate() (param *time.Time) {
	return &i.dueDate
}

func (i *ImageFinding) GroupID() string {
	return i.finding.Registry()
}

func (i *ImageFinding) HostName() (param *string) {
	val := i.finding.ImageVersion()
	return &val
}

func (i *ImageFinding) ID() (param int) {
	return
}

func (i *ImageFinding) IPAddress() (param *string) {
	return
}

func (i *ImageFinding) Labels() (param *string) {
	return
}

func (i *ImageFinding) LastChecked() (param *time.Time) {
	return i.finding.LastFound()
}

func (i *ImageFinding) MacAddress() (param *string) {
	val := i.finding.ImageTag()
	return &val
}

func (i *ImageFinding) MethodOfDiscovery() (param *string) {
	return &i.mod
}

func (i *ImageFinding) OSDetailed() (param *string) {
	return
}

func (i *ImageFinding) OperatingSystem() (param *string) {
	return
}

func (i *ImageFinding) OrgCode() (param *string) {
	return &i.orgCode
}

func (i *ImageFinding) OrganizationID() (param string) {
	return
}

func (i *ImageFinding) Patchable() (param *string) {
	return i.finding.Patchable()
}

func (i *ImageFinding) Priority() (param *string) {
	return &i.priority
}

func (i *ImageFinding) Project() (param *string) {
	return
}

func (i *ImageFinding) ReportedBy() (param *string) {
	return
}

func (i *ImageFinding) ResolutionDate() (param *time.Time) {
	return
}

func (i *ImageFinding) ResolutionStatus() (param *string) {
	return
}

func (i *ImageFinding) ScanID() (param int) {
	return
}

func (i *ImageFinding) ServicePorts() (param *string) {
	val := fmt.Sprintf("0 %s", i.finding.VulnerabilityLocation())
	return &val
}

func (i *ImageFinding) Solution() (param *string) {
	return i.finding.Solution()
}

func (i *ImageFinding) Status() (param *string) {
	return
}

func (i *ImageFinding) Summary() (param *string) {
	val := fmt.Sprintf("%s on %s", i.finding.VulnerabilityID(), i.finding.ImageName())
	return &val
}

func (i *ImageFinding) TicketType() (param *string) {
	val := "Request"
	return &val
}

func (i *ImageFinding) Title() (param string) {
	return
}

func (i *ImageFinding) SystemName() (param *string) {
	return nil
}

func (i *ImageFinding) UpdatedDate() (param *time.Time) {
	return i.finding.LastUpdated()
}

func (i *ImageFinding) VendorReferences() (param *string) {
	val := i.finding.VendorReference()
	return &val
}

func (i *ImageFinding) VulnerabilityID() (param string) {
	return i.finding.VulnerabilityID()
}

func (i *ImageFinding) VulnerabilityTitle() (param *string) {
	return
}

func (i *ImageFinding) OWASP() (param *string) {
	return nil
}

func (i *ImageFinding) ExceptionDate() (param *time.Time) {
	return nil
}

func (i *ImageFinding) ApplicationName() (param *string) {
	return nil
}

func (i *ImageFinding) TrackingMethod() (param *string) {
	return nil
}

func mapImageFindingsByDeviceIDVulnID(findings []domain.ImageFinding) (entityIDToRuleHashToFinding map[string]domain.ImageFinding) {
	// DeviceID = Image name
	// VulnID = CVE

	entityIDToRuleHashToFinding = make(map[string]domain.ImageFinding)
	for _, finding := range findings {
		if len(finding.ImageName()) > 0 {
			// leading 0 in the last element is to match up with the [port protocol] format in the ServicePorts
			// since no value for port is relevant to an image finding, a 0 is left here
			key := fmt.Sprintf("%s;%s;0 %s", finding.ImageName(), finding.VulnerabilityID(), finding.VulnerabilityLocation())
			entityIDToRuleHashToFinding[key] = finding
		}
	}

	return entityIDToRuleHashToFinding
}
