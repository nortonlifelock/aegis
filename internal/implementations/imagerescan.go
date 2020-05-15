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
	"sync"
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
	RegistryImage []string `json:"registry_image"`
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

type imageFindingTicketPair struct {
	finding domain.ImageFinding
	ticket  domain.Ticket
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
				// findings without tickets need tickets created for  them
				var findingsWithoutTickets = make([]domain.ImageFinding, 0)

				// tickets without findings may be closed
				var ticketsWithoutFindings = make([]domain.Ticket, 0)

				// tickets with findings can have their last seen date updated and should be reopened
				var ticketsWithFindings = make([]imageFindingTicketPair, 0)

				ticketSlice := fanInChannel(tickets)

				deviceIDToVulnIDToTicket := mapTicketsByDeviceIDVulnID(ticketSlice)
				deviceIDToVulnIDToFinding := mapImageFindingsByDeviceIDVulnID(findings)

				for _, finding := range findings {

					if deviceIDToVulnIDToTicket[finding.ImageName()] != nil {
						if deviceIDToVulnIDToTicket[finding.ImageName()][finding.VulnerabilityID()] != nil {
							for _, tieTicketToFinding := range deviceIDToVulnIDToTicket[finding.ImageName()][finding.VulnerabilityID()] {
								ticketsWithFindings = append(ticketsWithFindings, imageFindingTicketPair{
									finding: finding,
									ticket:  tieTicketToFinding,
								})
							}
						} else {
							findingsWithoutTickets = append(findingsWithoutTickets, finding)
						}
					} else {
						findingsWithoutTickets = append(findingsWithoutTickets, finding)
					}
				}

				for _, ticket := range ticketSlice {
					if deviceIDToVulnIDToFinding[ticket.DeviceID()] != nil {
						if deviceIDToVulnIDToFinding[ticket.DeviceID()][ticket.VulnerabilityID()] == nil {
							ticketsWithoutFindings = append(ticketsWithoutFindings, ticket)
						}
					} else {
						ticketsWithoutFindings = append(ticketsWithoutFindings, ticket)
					}
				}

				job.updateTicketsAccordingToFindings(engine, findingsWithoutTickets, ticketsWithFindings, ticketsWithoutFindings)
			}
		}
	} else {
		err = fmt.Errorf("registry_image appears to be malformed, expected it to be in the form [REGISTRY;IMAGE]")
	}

	return err
}

func (job *ImageRescanJob) updateTicketsAccordingToFindings(engine integrations.TicketingEngine, findingsWithoutTickets []domain.ImageFinding, ticketsWithFindings []imageFindingTicketPair, ticketsWithoutFindings []domain.Ticket) {
	wg := &sync.WaitGroup{}
	wg.Add(3)
	go func() {
		defer handleRoutinePanic(job.lstream)
		defer wg.Done()
		job.createTicketsForUnticketedFindings(engine, findingsWithoutTickets)
	}()
	go func() {
		defer handleRoutinePanic(job.lstream)
		defer wg.Done()
		job.updateTicketsWithStaleFindings(engine, ticketsWithFindings)
	}()
	go func() {
		defer handleRoutinePanic(job.lstream)
		defer wg.Done()
		job.closeTicketsWithMissingFindings(engine, ticketsWithoutFindings)
	}()
	wg.Wait()
}

func (job *ImageRescanJob) closeTicketsWithMissingFindings(engine integrations.TicketingEngine, tickets []domain.Ticket) {
	wg := &sync.WaitGroup{}
	for index := range tickets {
		wg.Add(1)
		go func(ticket domain.Ticket) {
			defer handleRoutinePanic(job.lstream)
			defer wg.Done()

			_, _, err := engine.UpdateTicket(
				closedTicket{
					ticket,
					engine,
				},
				fmt.Sprintf("finding was NOT found by %s", job.insource.Source()),
			)

			if err == nil {
				job.lstream.Send(log.Infof("finding for %s was NOT by %s, closing ticket...", ticket.Title(), job.insource.Source()))
			} else {
				job.lstream.Send(log.Errorf(err, "error while closing ticket [%s]", ticket.Title()))
			}
		}(tickets[index])
	}
	wg.Wait()
}

func (job *ImageRescanJob) updateTicketsWithStaleFindings(engine integrations.TicketingEngine, ticketsWithFindings []imageFindingTicketPair) {
	wg := &sync.WaitGroup{}
	for index := range ticketsWithFindings {
		wg.Add(1)
		go func(pair imageFindingTicketPair) {
			defer handleRoutinePanic(job.lstream)
			defer wg.Done()

			_, _, err := engine.UpdateTicket(
				&staleTicket{
					pair.ticket,
					engine,
				},
				fmt.Sprintf("finding still detected by %s on [%s]", job.insource.Source(), time.Now().Format(time.RFC822)),
			)
			if err == nil {
				job.lstream.Send(log.Infof("finding for %s still detected by %s", pair.ticket.Title(), job.insource.Source()))
			} else {
				job.lstream.Send(log.Errorf(err, "error while updating ticket [%s]", pair.ticket.Title()))
			}
		}(ticketsWithFindings[index])
	}
	wg.Wait()
}

func (job *ImageRescanJob) createTicketsForUnticketedFindings(engine integrations.TicketingEngine, findings []domain.ImageFinding) {
	wg := &sync.WaitGroup{}
	for index := range findings {

		if ford(findings[index].CVSS2()) >= job.orgPayload.LowestCVSS {
			wg.Add(1)
			go func(finding domain.ImageFinding) {
				defer handleRoutinePanic(job.lstream)
				defer wg.Done()

				ignore, err := job.db.HasIgnore(job.insource.SourceID(), finding.VulnerabilityID(), finding.ImageName(), job.config.OrganizationID(), "", time.Now())
				if err != nil {
					job.lstream.Send(log.Errorf(err, "error while loading ignore for Aqua entry [%s|%s]", finding.ImageName(), finding.VulnerabilityID()))
				}

				if ignore == nil {
					var priority string
					var dueDate time.Time
					if priority, dueDate, err = calculateSLAForImageFinding(finding, job.orgPayload); err == nil {
						var ticket = &ImageFinding{
							finding,
							job.insource.Source(),
							job.orgCode,
							dueDate,
							priority,
						}

						_, sourceKey, err := engine.CreateTicket(ticket)
						if err == nil {
							job.lstream.Send(log.Infof("Created ticket [%s] for image [%s] on vulnerability [%s]", sourceKey, finding.ImageName(), finding.VulnerabilityID()))
						} else {
							job.lstream.Send(log.Errorf(err, "error while creating ticket for image [%s] on vulnerability [%s]", finding.ImageName(), finding.VulnerabilityID()))
						}
					} else {
						job.lstream.Send(log.Errorf(err, "error while calculating SLAs for Image/Vuln [%s|%s]", finding.ImageName(), finding.VulnerabilityID()))
					}

				} else {
					job.lstream.Send(log.Infof("SKIPPING ticket for [%s|%s] as it has an ignore entry", finding.ImageName(), finding.VulnerabilityID()))
				}

			}(findings[index])
		}

	}
	wg.Wait()
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

func (i *ImageFinding) AssignedTo() (param *string) {
	return
}

func (i *ImageFinding) AssignmentGroup() (param *string) {
	return
}

func (i *ImageFinding) CERF() (param string) {
	return
}

func (i *ImageFinding) CERFExpirationDate() (param time.Time) {
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
	return
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
	return
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

func mapImageFindingsByDeviceIDVulnID(findings []domain.ImageFinding) (entityIDToRuleHashToFinding map[string]map[string]domain.ImageFinding) {
	// DeviceID = Image name
	// VulnID = CVE

	entityIDToRuleHashToFinding = make(map[string]map[string]domain.ImageFinding)
	for _, finding := range findings {
		if len(finding.ImageName()) > 0 {
			if entityIDToRuleHashToFinding[finding.ImageName()] == nil {
				entityIDToRuleHashToFinding[finding.ImageName()] = make(map[string]domain.ImageFinding)
			}

			entityIDToRuleHashToFinding[finding.ImageName()][finding.VulnerabilityID()] = finding
		}
	}

	return entityIDToRuleHashToFinding
}
