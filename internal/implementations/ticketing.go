package implementations

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nortonlifelock/aegis/internal/database/dal"
	"github.com/nortonlifelock/aegis/internal/integrations"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/log"
	"github.com/nortonlifelock/scaffold"
	"github.com/pkg/errors"
)

// TicketingPayload decides which asset groups to ticket on, as well as defining the min date which is used to calculate the SLA if the calculated
// due date is in the past
type TicketingPayload struct {
	MinDate *time.Time `json:"mindate,omitempty"`
	Groups  []string   `json:"groups,omitempty"`
}

// OrgPayload contains the SLA information for how long a vulnerability has to be remediated given the severity
// it is located from the Payload field of the organization table
type OrgPayload struct {
	LowestCVSS        float32       `json:"lowest_ticketed_cvss"`
	CVSSVersion       int           `json:"cvss_version"`
	Severities        []OrgSeverity `json:"severities"`
	DescriptionFooter string        `json:"description_footer"`
}

const (
	cvssVersion2 = 2
	cvssVersion3 = 3
)

// OrgSeverity holds the information pertaining to the severity and it's relation to CVSS. The severities are organized based on their CVSS minimum score
// CVSSMin dictates the lowest score required for a vulnerability to be associated with this severity. If another severity has a higher CVSS min that
// the vulnerability is also above, the vulnerability is associated with that CVSS min. The duration is the amount of time in days that a remediator would
// have to fix the vulnerability after discovery
type OrgSeverity struct {
	Name     string  `json:"name"`
	Duration int     `json:"duration"`
	CVSSMin  float32 `json:"cvss_min"`
}

// Len implements the sort interface so the severities may be organized
func (payload *OrgPayload) Len() int {
	return len(payload.Severities)
}

// Less identifies which severity entry has a lower CVSS minimum
func (payload *OrgPayload) Less(i, j int) bool {
	return payload.Severities[i].CVSSMin < payload.Severities[j].CVSSMin
}

// Swap swaps two severity entries
func (payload *OrgPayload) Swap(i, j int) {
	payload.Severities[i], payload.Severities[j] = payload.Severities[j], payload.Severities[i]
}

// Validate ensures there is a severity description for an organization, sorts them, and ensures all the numerical values
// held are valid
// additionally, it checks that the cvss version is set within the organization payload
func (payload *OrgPayload) Validate() (valid bool) {
	if len(payload.Severities) > 0 {
		sort.Sort(payload)

		var allNonZero = true
		for _, entry := range payload.Severities {
			if entry.CVSSMin < 0 || entry.Duration < 0 {
				allNonZero = false
				break
			}
		}

		if allNonZero {

			var noOverlap = true
			for index := range payload.Severities {
				if index > 0 {
					if payload.Severities[index].CVSSMin <= payload.Severities[index-1].CVSSMin {
						noOverlap = false
					}
				}
			}

			if noOverlap {
				valid = payload.CVSSVersion == cvssVersion2 || payload.CVSSVersion == cvssVersion3
			}
		}
	}

	return valid
}

// TicketingJob implements the IJob interface required to run the job
type TicketingJob struct {
	Payload *TicketingPayload

	ticketMutex     *sync.Mutex
	ticketingEngine integrations.TicketingEngine
	duplicatesMap   sync.Map

	// TODO: remove the port flag from the code, these should always create multiple tickets
	OrgPayload *OrgPayload

	id          string
	payloadJSON string
	ctx         context.Context
	db          domain.DatabaseConnection
	lstream     log.Logger
	appconfig   domain.Config
	config      domain.JobConfig
	insource    domain.SourceConfig
	outsource   domain.SourceConfig

	cachedReportedBy string
	assignmentRules  []assignmentRule
	tagMaps          []domain.TagMap
}

// vulnerabilityPayload is passed through the pipeline of the ticketing job
type vulnerabilityPayload struct {
	// ticketing engine is cached in order for multiple threads to share a connection
	tickets integrations.TicketingEngine

	// the organization code is used in the ticket and must be pulled from the database, so it is cached
	orgCode string

	combo domain.Detection
	// device, vuln, and detectedDate are pulled off combo using Accessor methods, but are cached to prevent repeated error checking
	device       domain.Device
	vuln         domain.Vulnerability
	detectedDate *time.Time

	// holds the statuses that are used to query existing tickets when checking for duplicates
	statuses map[string]bool

	// ticket is populated at the end of the process for creation in the ticketing engine
	ticket domain.Ticket
}

// Tag mapping options
const (
	// Append states that the tag mapping information should be included in addition to the information from the scanner
	Append = "Append"

	// Overwrite states that the tag mapping information should replace the information from the scanner
	Overwrite = "Overwrite"
)

// buildPayload loads the Payload from the job history into the Payload object
func (job *TicketingJob) buildPayload(pjson string) (err error) {

	if len(pjson) > 0 {

		job.Payload = &TicketingPayload{}
		job.ticketMutex = &sync.Mutex{}

		err = json.Unmarshal([]byte(pjson), job.Payload)
	} else {
		err = errors.New("Payload length is 0")
	}

	return err
}

func (job *TicketingJob) buildOrgPayload(org domain.Organization) (err error) {
	if len(org.Payload()) > 0 {
		job.OrgPayload = &OrgPayload{}

		err = json.Unmarshal([]byte(org.Payload()), job.OrgPayload)
		if err == nil {
			if !job.OrgPayload.Validate() {
				err = fmt.Errorf("organization payload validation failed")
			}
		}
	} else {
		err = errors.New("Payload length is 0")
	}

	return err
}

// Process the ticketing job loads device information from a scanner, and creates a ticket for each device/vulnerability combination where one does not
// already exist. First, it checks for an entry in the ignore table to see if that device/vulnerability combination is a known exception or false positive
func (job *TicketingJob) Process(ctx context.Context, id string, appconfig domain.Config, db domain.DatabaseConnection, lstream log.Logger, payload string, jobConfig domain.JobConfig, inSource []domain.SourceConfig, outSource []domain.SourceConfig) (err error) {

	var ok bool
	if job.ctx, job.id, job.appconfig, job.db, job.lstream, job.payloadJSON, job.config, job.insource, job.outsource, ok = validInputs(ctx, id, appconfig, db, lstream, payload, jobConfig, inSource, outSource); ok {

		if err = job.buildPayload(job.payloadJSON); err == nil {

			var org domain.Organization
			if org, err = job.db.GetOrganizationByID(job.config.OrganizationID()); err == nil {
				var vscanner integrations.Vscanner
				if vscanner, err = integrations.NewVulnScanner(job.ctx, job.insource.Source(), job.db, job.lstream, job.appconfig, job.insource); vscanner != nil && err == nil {
					if org != nil {

						// the organization Payload holds the SLA configuration
						if err = job.buildOrgPayload(org); err == nil {
							if job.assignmentRules, err = job.loadAssignmentRules(); err == nil {
								if job.tagMaps, err = job.db.GetTagMapsByOrg(job.config.OrganizationID()); err == nil {
									job.lstream.Send(log.Debug("Scanner connection initialized."))

									if len(job.Payload.Groups) == 0 {

										job.lstream.Send(log.Infof("Pulling ALL detections since %s", tord1970(job.config.LastJobStart()).String()))

										startTime := time.Now() // must be before we load the detections from the db
										var detections []domain.Detection
										if detections, err = job.db.GetDetectionsAfter(tord1970(job.config.LastJobStart()), job.config.OrganizationID()); err == nil {
											job.processVulnerabilities(vscanner, pushDetectionsToChannel(job.ctx, detections))

											// passing an empty group ID updates the last ticketing date of ALL groups
											_, _, err = job.db.UpdateAssetGroupLastTicket("", job.config.OrganizationID(), startTime)
											if err != nil {
												job.lstream.Send(log.Criticalf(err, "Error while updating the last ticketed date to %s", startTime.String()))
											}
										} else {
											job.lstream.Send(log.Error("Error occurred while loading device vulnerability information", err))
										}
									} else {
										for _, groupID := range job.Payload.Groups {

											if len(groupID) > 0 {
												var assetGroup domain.AssetGroup
												if assetGroup, err = job.db.GetAssetGroupForOrgNoScanner(job.config.OrganizationID(), groupID); err == nil {

													var after time.Time
													if assetGroup.LastTicketing() == nil || assetGroup.LastTicketing().IsZero() {
														after = tord1970(nil)
													} else {
														after = *assetGroup.LastTicketing()
													}

													job.lstream.Send(log.Infof("Pulling all detections for group [%s] since %s", groupID, after.String()))
													startTime := time.Now() // must be before we load the detections from the db

													var detections []domain.Detection
													if detections, err = job.db.GetDetectionForGroupAfter(after, job.config.OrganizationID(), groupID); err == nil {

														job.processVulnerabilities(vscanner, pushDetectionsToChannel(job.ctx, detections))

														// passing an empty group ID updates the last ticketing date of ALL groups
														_, _, err = job.db.UpdateAssetGroupLastTicket(groupID, job.config.OrganizationID(), startTime)
														if err != nil {
															job.lstream.Send(log.Criticalf(err, "Error while updating the last ticketed date to %s", startTime.String()))
														}
													} else {
														job.lstream.Send(log.Error("Error occurred while loading device vulnerability information", err))
													}
												} else {
													err = fmt.Errorf("error while loading asset group for [%s] - %s", groupID, err.Error())
												}
											} else {
												err = fmt.Errorf("empty group id in payload")
											}
										}

									}

								} else {
									job.lstream.Send(log.Error("error while loading tag maps", err))
								}
							} else {
								job.lstream.Send(log.Errorf(err, "error while loading assignment rules"))
							}
						} else {
							job.lstream.Send(log.Error("error while processing the organization Payload", err))
						}

					} else {
						job.lstream.Send(log.Errorf(nil, "Null org object returned."))
					}
				} else {
					err = fmt.Errorf("error while creating the vuln scanner: [%v]", err)
				}
			} else {
				err = fmt.Errorf("could not find organization by this ID: [%s] - %s", job.config.OrganizationID(), err.Error())
				job.lstream.Send(log.Error("Error while getting organization.", err))
			}
		} else {
			err = fmt.Errorf("error while building payload - %s", err.Error())
		}
	} else {
		err = fmt.Errorf("input validation failed")
	}

	return err
}

// processVulnerabilities creates a pipeline of channels. each method takes a channel as an input, and creates a channel as an output
// the first pipe in the pipeline is process vulnerability, and the final pipe is create ticket. each method takes an input
// from a channel, performs some transformation on the input, and pushes the result on the output channel for the next method
// to handle
func (job *TicketingJob) processVulnerabilities(vscanner integrations.Vscanner, in <-chan domain.Detection) {
	job.createTicket(
		job.prepareTicketCreation(
			job.checkForExistingTicket(
				job.processVulnerability(in),
			),
		),
	)
}

func (job *TicketingJob) processVulnerability(in <-chan domain.Detection) <-chan *vulnerabilityPayload {
	defer handleRoutinePanic(job.lstream)
	var out = make(chan *vulnerabilityPayload)

	go func() {
		defer handleRoutinePanic(job.lstream)
		defer close(out)
		wg := &sync.WaitGroup{}

		var orgcode = job.getOrgCode()
		var err error

		job.lstream.Send(log.Debugf("Opening connection to job engine for Job ID [%v].", job.id))
		var tickets integrations.TicketingEngine
		if tickets, err = integrations.GetEngine(job.ctx, job.outsource.Source(), job.db, job.lstream, job.appconfig, job.outsource); err == nil {
			job.lstream.Send(log.Debugf("Connection opened to job engine for Job ID [%v].", job.id))
			job.ticketingEngine = tickets

			const numConcurrentThreads = 100
			var permit = make(chan bool, numConcurrentThreads)
			for i := 0; i < numConcurrentThreads; i++ {
				permit <- true
			}

			func() {
				for {
					select {
					case <-job.ctx.Done():
						return
					case item, ok := <-in:
						if ok {

							select {
							case <-job.ctx.Done():
								return
							case <-permit:
							}

							wg.Add(1)
							go func(dvCombo domain.Detection) {
								defer wg.Done()
								defer handleRoutinePanic(job.lstream)
								defer func() {
									permit <- true
								}()

								if strings.ToLower(dvCombo.Status()) != strings.ToLower(domain.Fixed) {
									var err error

									var device domain.Device
									var vuln domain.Vulnerability
									var detectedDate *time.Time

									device, err = dvCombo.Device()
									if err == nil {
										vuln, err = dvCombo.Vulnerability()
										if err == nil {
											detectedDate, err = dvCombo.Detected()
										}
									}

									if err == nil {
										if device != nil && vuln != nil && detectedDate != nil {
											if job.getCVSSScore(vuln) >= job.OrgPayload.LowestCVSS {

												statuses := make(map[string]bool)
												loadStatuses(tickets, statuses)

												job.lstream.Send(log.Infof("Processing vulnerability [%s] on device [%v]", vuln.SourceID(), sord(device.SourceID())))

												var payload = &vulnerabilityPayload{
													tickets,
													orgcode,
													dvCombo,
													device,
													vuln,
													detectedDate,
													statuses,
													nil,
												}

												select {
												case <-job.ctx.Done():
													return
												case out <- payload:
												}
											} else {
												job.lstream.Send(log.Debugf("Skipping vulnerability [%s] on device [%v] with CVSS [%v].", vuln.SourceID(), sord(device.SourceID()), job.getCVSSScore(vuln)))
											}
										} else {
											job.lstream.Send(log.Errorf(err, "failed to load vulnerability information for [%v|%v|%v]", device, vuln, detectedDate))
										}
									} else {
										job.lstream.Send(log.Errorf(err, "error while processing vulnerability %v", dvCombo.VulnerabilityID()))
									}
								} else {
									// vulnerability fixed - don't create ticket
								}
							}(item)
						} else {
							return
						}
					}
				}
			}()

			wg.Wait()
		} else {
			job.lstream.Send(log.Error("Error while getting job object.", err))
		}
	}()

	return out
}

func (job *TicketingJob) getOrgCode() (orgCode string) {
	if len(job.config.OrganizationID()) > 0 {

		// Get the organization from the database using the id in the ticket object
		if org, err := job.db.GetOrganizationByID(job.config.OrganizationID()); err == nil {
			// Ensure there is only one return
			if org != nil {
				orgCode = org.Code()
			} else {
				job.lstream.Send(log.Criticalf(err, "failed to load the organization for ID [%v]", job.config.OrganizationID()))
			}
		}
	}
	return orgCode
}

func loadStatuses(tickets integrations.TicketingEngine, statuses map[string]bool) {
	// Statuses to Query when looking for existing tickets for the vulnerabilities

	// TODO TODO do we want these hardcoded or configurable?
	statuses[tickets.GetStatusMap(domain.StatusOpen)] = true
	statuses[tickets.GetStatusMap(domain.StatusReopened)] = true
	statuses[tickets.GetStatusMap(domain.StatusResolvedRemediated)] = true
	statuses[tickets.GetStatusMap(domain.StatusResolvedDecom)] = true
	statuses[tickets.GetStatusMap(domain.StatusResolvedException)] = true
	statuses[tickets.GetStatusMap(domain.StatusResolvedFalsePositive)] = true
	statuses[tickets.GetStatusMap(domain.StatusClosedCerf)] = true

	// TODO: Remove this once the closed-error status is part of exceptions
	statuses[tickets.GetStatusMap(domain.StatusClosedError)] = true
}

func (job *TicketingJob) checkForExistingTicket(in <-chan *vulnerabilityPayload) <-chan *vulnerabilityPayload {
	defer handleRoutinePanic(job.lstream)

	var out = make(chan *vulnerabilityPayload)
	go func() {
		defer handleRoutinePanic(job.lstream)
		defer close(out)
		wg := &sync.WaitGroup{}

		for {

			var payload *vulnerabilityPayload
			var ok bool

			select {
			case <-job.ctx.Done():
				return
			case payload, ok = <-in:
				// do nothing
			}

			if ok {

				var port int
				var protocol string

				port = payload.combo.Port()
				protocol = payload.combo.Protocol()

				var keyToPreventDuplicates = fmt.Sprintf("%v-%v-%v-%v", sord(payload.device.SourceID()), payload.vuln.SourceID(), port, protocol)

				var exists bool
				if _, exists = job.duplicatesMap.LoadOrStore(keyToPreventDuplicates, true); !exists { // doesn't exist in sync map

					wg.Add(1)
					go func(payload *vulnerabilityPayload, exists bool, port int, protocol string) {
						defer handleRoutinePanic(job.lstream)
						defer wg.Done()

						var err error

						var existingTicket domain.TicketSummary
						if existingTicket, err = job.db.GetTicketByDeviceIDVulnID(sord(payload.device.SourceID()), payload.vuln.ID(), port, protocol, job.config.OrganizationID()); err == nil { // TODO is this vuln ID correct? I would be happiest if the device lookup didn't use the source id
							if existingTicket == nil {

								var existingTicketChan <-chan domain.Ticket
								var statuses = make(map[string]bool)
								statuses[job.ticketingEngine.GetStatusMap(domain.StatusOpen)] = true
								statuses[job.ticketingEngine.GetStatusMap(domain.StatusInProgress)] = true
								statuses[job.ticketingEngine.GetStatusMap(domain.StatusReopened)] = true
								statuses[job.ticketingEngine.GetStatusMap(domain.StatusResolvedRemediated)] = true
								statuses[job.ticketingEngine.GetStatusMap(domain.StatusResolvedFalsePositive)] = true
								statuses[job.ticketingEngine.GetStatusMap(domain.StatusResolvedDecom)] = true
								statuses[job.ticketingEngine.GetStatusMap(domain.StatusResolvedException)] = true
								existingTicketChan, err = job.ticketingEngine.GetTicketsByDeviceIDVulnID(job.insource.Source(), payload.orgCode, sord(payload.device.SourceID()), payload.vuln.SourceID(), statuses, payload.combo.Port(), payload.combo.Protocol())
								if err == nil {

									if emptyChannel(existingTicketChan) {
										job.lstream.Send(log.Infof("No ticket found for vulnerability [%s] on device [%v]. Creating new ticket...", payload.vuln.SourceID(), sord(payload.device.SourceID())))
										select {
										case <-job.ctx.Done():
											return
										case out <- payload:
										}
									}
								} else {
									job.lstream.Send(log.Error(
										fmt.Sprintf(
											"Error issues from JIRA with vuln title [%v] and ID [%v].",
											payload.vuln.Name(),
											payload.vuln.SourceID(),
										),
										err,
									))
								}

							} else {
								job.lstream.Send(log.Info(
									fmt.Sprintf(
										"EXISTING TICKET: [%v] for vulnerability [%v] with Vuln ID [%v] on device [%v]. Skipping...",
										existingTicket.Title(),
										payload.vuln.Name(),
										payload.vuln.SourceID(),
										sord(payload.device.SourceID()),
									)))
							}
						} else {
							job.lstream.Send(log.Warning(
								fmt.Sprintf(
									"Error getting issues from database with vuln title [%v] and ID [%v].",
									payload.vuln.Name(),
									payload.vuln.SourceID(),
								),
								err,
							))
						}
					}(payload, exists, port, protocol)
				} else {
					job.lstream.Send(log.Info(
						fmt.Sprintf(
							"ALREADY PROCESSED: A ticket was already created for vulnerability [%v] with Vuln ID [%v] on device [%v] during this run. Skipping...",
							payload.vuln.Name(),
							payload.vuln.SourceID(),
							sord(payload.device.SourceID()),
						)))
				}

			} else {
				break
			}

		}

		wg.Wait()
	}()

	return out
}

func emptyChannel(in <-chan domain.Ticket) bool {
	for {
		select {
		case _, ok := <-in:
			if ok {
				go func() {
					for {
						if _, ok := <-in; !ok {
							return
						}
					}
				}()
				return false
			} else {
				return true
			}
		}
	}
}

// takes the Payload and transforms it to a ticket. overwrites/appends information in the ticket fields from cloud service tags if a tag mapping & tags
// for the device are found
func (job *TicketingJob) prepareTicketCreation(in <-chan *vulnerabilityPayload) <-chan *vulnerabilityPayload {
	defer handleRoutinePanic(job.lstream)
	var out = make(chan *vulnerabilityPayload)

	go func() {
		defer handleRoutinePanic(job.lstream)
		defer close(out)
		wg := &sync.WaitGroup{}

		for {

			var payload *vulnerabilityPayload
			var ok bool

			select {
			case <-job.ctx.Done():
				return
			case payload, ok = <-in:
				// do nothing
			}

			if ok {

				wg.Add(1)
				go func(payload *vulnerabilityPayload) {
					defer handleRoutinePanic(job.lstream)
					defer wg.Done()

					var err error
					payload.ticket = &dal.Ticket{}
					var create bool
					payload.ticket, create = job.payloadToTicket(payload)
					if create {

						var tagsForDevice []domain.Tag
						// map cloud service fields to ticket if necessary
						tagsForDevice, err = job.handleCloudTagMappings(payload.ticket)
						if err == nil {
							job.getAssignmentInformation(tagsForDevice, payload)
						} else {
							// we still want to create the ticket, but log the error
							job.lstream.Send(log.Errorf(err, "error while managing job mappings for [%s]", payload.ticket.Title()))
						}

						select {
						case <-job.ctx.Done():
							return
						case out <- payload:
						}
					} else {
						job.lstream.Send(log.Infof("Skipping vulnerability with CVSS [%f]", payload.vuln.CVSS2()))
					}
				}(payload)
			} else {
				break
			}
		}

		wg.Wait()
	}()

	return out
}

func (job *TicketingJob) createTicket(in <-chan *vulnerabilityPayload) {
	defer handleRoutinePanic(job.lstream)

	var wg = &sync.WaitGroup{}
	for {

		payload, ok := <-in
		if ok {

			if payload != nil {

				if len(payload.ticket.VulnerabilityID()) > 0 {
					wg.Add(1)
					go func(payload *vulnerabilityPayload) {
						defer handleRoutinePanic(job.lstream)
						defer wg.Done()
						job.createIndividualTicket(payload)
					}(payload)
				} else {
					var err = errors.Errorf("%s had an invalid vulnerability id in createTicket", payload.ticket.VulnerabilityID())
					job.lstream.Send(log.Error(err.Error(), err))
				}
			} else {
				var err = errors.Errorf("Ticket received NIL from channel in createTicket | %v", payload)
				job.lstream.Send(log.Error(err.Error(), err))
			}

		} else {
			break
		}
	}
	wg.Wait()
}

func (job *TicketingJob) calculateSLA(vuln domain.Vulnerability, alertDate time.Time) (priority string, dueDate time.Time, create bool) {
	severity := job.getSLAForVuln(vuln)
	if severity != nil {
		create = true
		priority = severity.Name
		dueDate = job.calculateDueDate(alertDate, severity.Duration)
	}

	return priority, dueDate, create
}

func (job *TicketingJob) getSLAForVuln(vuln domain.Vulnerability) (highestApplicableSeverity *OrgSeverity) {
	var cvssScore = job.getCVSSScore(vuln)

	// we iterate over the sorted list of custom severity ranges and find the highest applicable severity
	for index := range job.OrgPayload.Severities {
		if cvssScore >= job.OrgPayload.Severities[index].CVSSMin {
			highestApplicableSeverity = &job.OrgPayload.Severities[index]
		}
	}

	return highestApplicableSeverity
}

func (job *TicketingJob) calculateDueDate(alertDate time.Time, durationInDays int) (dueDate time.Time) {
	dueDate = alertDate.AddDate(0, 0, durationInDays)

	if job.Payload.MinDate != nil {
		var minDate = job.Payload.MinDate.AddDate(0, 0, durationInDays)
		if dueDate.Before(minDate) {
			dueDate = minDate
		}
	}

	return dueDate
}

func (job *TicketingJob) createIndividualTicket(payload *vulnerabilityPayload) {
	if _, ticketTitle, err := job.ticketingEngine.CreateTicket(payload.ticket); err == nil {

		if len(ticketTitle) > 0 {
			job.lstream.Send(log.Info(
				fmt.Sprintf(
					"Ticket created for vulnerability [%s] on device [%v]. [Title: %s]",
					payload.ticket.VulnerabilityID(),
					payload.ticket.DeviceID(),
					ticketTitle,
				)))

			// track the created ticket in our database
			_, _, err = job.db.CreateTicket(
				ticketTitle,
				domain.StatusOpen,
				payload.combo.ID(),
				job.config.OrganizationID(),
				tord1970(payload.ticket.CreatedDate()),
				tord1970(payload.ticket.DueDate()),
				time.Now(),
				tord1970(payload.ticket.ResolutionDate()),
				tord1970(nil), // used to set the resolution date to nil in the DB if the ticket doesn't have one
			)

			if err != nil {
				job.lstream.Send(log.Errorf(err, "error while creating database entry for ticket [%v]", ticketTitle))
			}
		} else {
			job.lstream.Send(log.Error(
				fmt.Sprintf(
					"Could not retrieve ticket title created for vulnerability [%s] with vuln ID [%v] on device [%v]",
					*payload.ticket.VulnerabilityTitle(),
					payload.ticket.VulnerabilityID(),
					payload.ticket.DeviceID(),
				),
				err,
			))
		}
	} else {
		job.lstream.Send(log.Error(
			fmt.Sprintf(
				"Error while creating ticket for vulnerability [%s] with Vuln ID [%v] on device [%v]",
				*payload.ticket.VulnerabilityTitle(),
				payload.ticket.VulnerabilityID(),
				payload.ticket.DeviceID(),
			),
			err,
		))
	}
}

// takes a Payload for a ticket and transforms it to a dal ticket for creation
func (job *TicketingJob) payloadToTicket(payload *vulnerabilityPayload) (newtix *dal.Ticket, create bool) {

	// Handle address fields
	var macs string
	var hosts string
	var ips string
	macs, ips, hosts = job.gatherHostInfoFromDevice(payload)

	// Determine Due Date and Priority
	var duedate time.Time
	var alertdate = time.Now()
	if payload.detectedDate != nil {
		alertdate = *payload.detectedDate
	}
	var priority string
	priority, duedate, create = job.calculateSLA(payload.vuln, alertdate)
	if create {

		cves, vendorRefs := job.gatherReferences(payload)
		var configs string
		if len(vendorRefs) == 0 {
			// Anything other than CVE should be as a config vuln
			configs = "True"
		}

		// TODO: This needs to be updated to a better method in the next releases
		var servicePorts string
		if payload.combo.Port() >= 0 && payload.combo.Port() <= 65535 && len(payload.combo.Protocol()) > 0 {
			servicePorts = fmt.Sprintf("%d %s", payload.combo.Port(), payload.combo.Protocol())
		}

		var ticketType = "Request"
		var operatingSystem = job.gatherOSDropdown(payload.device.OS())

		// TODO make configurable
		var summary = fmt.Sprintf("VRR (%s - %s): %s", ips, hosts, payload.vuln.Name())

		var template *scaffold.Template
		template = scaffold.NewTemplateEmpty()
		template.UpdateBase(descriptionTemplate)
		template.Repl("%description", payload.vuln.Description()).
			Repl("%proof", payload.combo.Proof())

		if len(sord(payload.vuln.Threat())) > 0 {
			template.Repl("%threat", fmt.Sprintf("\n*Threat:*\n%s\n", sord(payload.vuln.Threat())))
		} else {
			template.Repl("%threat", "")
		}

		var description = template.Get()

		if len(job.OrgPayload.DescriptionFooter) > 0 {
			description = fmt.Sprintf("%s\n\n%s", description, job.OrgPayload.DescriptionFooter)
		}

		var solution = removeHTMLTags(job.gatherSolution(payload))
		var methodOfDiscovery = job.insource.Source()
		var vulnerabilityTitle = payload.vuln.Name()
		var cvss = job.getCVSSScore(payload.vuln)
		var fullOSName = payload.device.OS()
		var reportedBy = job.getCachedReportedBy()
		var created = time.Now()

		newtix = &dal.Ticket{
			DeviceIDvar:          sord(payload.device.SourceID()),
			GroupIDvar:           sord(payload.device.GroupID()),
			VulnerabilityIDvar:   payload.vuln.SourceID(),
			MethodOfDiscoveryvar: &methodOfDiscovery,

			Descriptionvar:        &description,
			Summaryvar:            &summary,
			Solutionvar:           &solution,
			VulnerabilityTitlevar: &vulnerabilityTitle,
			CVSSvar:               &cvss,

			OSDetailedvar:      &fullOSName,
			OperatingSystemvar: &operatingSystem,
			MacAddressvar:      &macs,
			IPAddressvar:       &ips,
			HostNamevar:        &hosts,

			ReportedByvar:     &reportedBy,
			TicketTypevar:     &ticketType,
			OrganizationIDvar: job.config.OrganizationID(),
			Priorityvar:       &priority,

			Configsvar:          configs,
			ServicePortsvar:     &servicePorts,
			VendorReferencesvar: &vendorRefs,
			CVEReferencesvar:    &cves,

			CreatedDatevar: &created,
			AlertDatevar:   &alertdate,
			DueDatevar:     &duedate,
			OrgCodevar:     &payload.orgCode,
		}
	}

	return newtix, create
}

func (job *TicketingJob) gatherSolution(payload *vulnerabilityPayload) (solution string) {

	ctx, cancel := context.WithCancel(job.ctx)
	defer cancel()

	sols, err := payload.vuln.Solutions(ctx)
	if err == nil {
		for {
			select {
			case <-job.ctx.Done():
				return
			case sol, ok := <-sols:
				if ok {
					solution = sol.String()
				}

				return
			}
		}
	} else {
		job.lstream.Send(log.Errorf(err, "error while gathering solution for vulnerability %s", payload.vuln.SourceID()))
	}

	return solution
}

func (job *TicketingJob) gatherReferences(payload *vulnerabilityPayload) (cves string, vendorRefs string) {
	refs, err := payload.vuln.References(job.ctx)
	if err == nil {
		func() {
			for {
				select {
				case <-job.ctx.Done():
					return
				case ref, ok := <-refs:
					if ok {
						if strings.Contains(ref.Reference(), "CVE") {
							cves += ref.Reference() + ","
						} else {
							vendorRefs += ref.Reference() + ","
						}
					} else {
						return
					}
				}
			}
		}()

		cves = strings.TrimRight(cves, ",")
		vendorRefs = strings.TrimRight(vendorRefs, ",")
	} else {
		job.lstream.Send(log.Errorf(err, "error while gathering references for vulnerability %v", payload.vuln.SourceID()))
	}

	return cves, vendorRefs
}

func (job *TicketingJob) gatherHostInfoFromDevice(payload *vulnerabilityPayload) (string, string, string) {
	var macs = payload.device.MAC()
	var hosts = payload.device.HostName()
	var ips = payload.device.IP()

	return macs, ips, hosts
}

// the cloud sync job pulls tag information from cloud service providers. we can use that tag information to overwrite JIRA fields or append
// the information to a JIRA field
func (job *TicketingJob) handleCloudTagMappings(tic domain.Ticket) (tagsForDevice []domain.Tag, err error) {
	tagsForDevice = make([]domain.Tag, 0)

	if len(sord(tic.IPAddress())) > 0 {
		var ips = strings.Split(sord(tic.IPAddress()), ",")

		var device domain.Device
		device, err = job.getDeviceByIPList(ips)

		if err == nil {
			if device != nil { // device with ip found in database, check for it's tags
				// grab all the cloud tags for a device
				tagsForDevice, err = job.db.GetTagsForDevice(device.ID())
				if err == nil {
					if len(job.tagMaps) > 0 {
						err = job.mapAllTagsForDevice(tic, tagsForDevice, job.tagMaps)
					}
				}
			} else {
				// TODO no device found in db - email warning
				job.lstream.Send(log.Warningf(nil, "could not find device with any of ips [%s]", sord(tic.IPAddress())))
			}
		} else {
			err = fmt.Errorf("error while grabbing device - %s", err.Error())
		}

	} else {
		err = fmt.Errorf("ticket [%s] did not have an associated IP", tic.Title())
	}

	return tagsForDevice, err
}

// this ticket takes all tags found for a particular device, and maps them to fields within the domain.Ticket if necessary
func (job *TicketingJob) mapAllTagsForDevice(tic domain.Ticket, tagsForDevice []domain.Tag, tagMaps []domain.TagMap) (err error) {
	for index := range tagsForDevice {
		tagForDevice := tagsForDevice[index]

		var tagForDeviceKey domain.TagKey
		tagForDeviceKey, err = job.db.GetTagKeyByID(strconv.Itoa(tagForDevice.TagKeyID()))
		if err == nil {
			if tagForDeviceKey != nil {
				err = job.mapTagForDevice(tic, tagForDeviceKey, tagForDevice, tagMaps)
				if err != nil {
					break
				}
			} else {
				err = fmt.Errorf("could not find tag key [%d] in the database", tagForDevice.TagKeyID())
				break
			}
		} else {
			err = fmt.Errorf("error while grabbing tag key from database - %s", err.Error())
			break
		}
	}

	return err
}

// check to see if the tags found for a ticket match any of the fields in the tag map
// a tag map associates a JIRA field to a cloud service tag
func (job *TicketingJob) mapTagForDevice(tic domain.Ticket, tagForDeviceKey domain.TagKey, tagForDevice domain.Tag, tagMaps []domain.TagMap) (err error) {
	for mapIndex := range tagMaps {
		tagMap := tagMaps[mapIndex]

		// see if the cloud tag is mapped to a job field
		if strings.ToLower(strings.ToLower(tagMap.CloudTag())) == strings.ToLower(tagForDeviceKey.KeyValue()) {
			var ticketKey = tagMap.TicketingTag()

			var option = tagMap.Options()
			if tagMap.Options() == Append || tagMap.Options() == Overwrite {

				tic = tagMappedTicket{
					tic,
					strings.ToLower(ticketKey),
					option,
					tagForDevice,
					tagMap.CloudTag(),
					sord(tic.HostName()),
					sord(tic.AssignmentGroup()),
					sord(tic.Labels()),
				}
			} else {
				err = fmt.Errorf("unrecognized tag mapping option: %s", tagMap.Options())
				job.lstream.Send(log.Error("mapping error", err))
			}
		}
	}

	return err
}

type tagMappedTicket struct {
	domain.Ticket
	ticketKeyLower string
	option         string
	tagForDevice   domain.Tag
	cloudTag       string

	hostname        string
	assignmentGroup string
	labels          string
}

func (tmt tagMappedTicket) HostName() *string {
	val := tmt.hostname
	if tmt.option == Append && len(tmt.hostname) > 0 {
		val = fmt.Sprintf("%s,%s", tmt.hostname, tmt.tagForDevice.Value())
	} else { //overwrite
		val = tmt.tagForDevice.Value()
	}
	return &val
}

func (tmt tagMappedTicket) AssignmentGroup() *string {
	val := tmt.assignmentGroup
	if tmt.option == Append && len(tmt.assignmentGroup) > 0 {
		val = fmt.Sprintf("%s,%s", tmt.assignmentGroup, tmt.tagForDevice.Value())
	} else { //overwrite
		val = tmt.tagForDevice.Value()
	}
	return &val
}

func (tmt tagMappedTicket) Labels() *string {
	val := fmt.Sprintf("%s-%s", strings.ToLower(tmt.cloudTag), tmt.tagForDevice.Value())

	if tmt.option == Append && len(tmt.labels) > 0 {
		val = fmt.Sprintf("%s,%s", tmt.labels, val)
	}
	return &val
}

func (job *TicketingJob) getDeviceByIPList(ips []string) (device domain.Device, err error) {
	for index := range ips {
		ip := ips[index]

		device, err = job.db.GetDeviceByIP(ip, job.config.OrganizationID())
		if err == nil {
			if device != nil {
				break
			}
		} else {
			break
		}
	}

	return device, err
}

// transforms the specific OS from the scanner and transforms it to a generic OS that can be chosen in a dropdown field
func (job *TicketingJob) gatherOSDropdown(input string) (output string) {
	var ost domain.OperatingSystemType
	var err error
	if ost, err = job.db.GetOperatingSystemType(input); err == nil {
		output = ost.Type()
	} else {
		output = unknown
		job.lstream.Send(log.Errorf(err, "error while loading operating system type for [%s]", input))
	}

	return output
}

const (
	descriptionTemplate = `
	*Description:*
	%description
	%threat
	*Proof:*
	%proof
	`
)

var reportedByMutex sync.Mutex

func (job *TicketingJob) getCachedReportedBy() (reportedBy string) {

	if len(job.cachedReportedBy) > 0 {
		reportedBy = job.cachedReportedBy
	} else {
		reportedByMutex.Lock()
		defer reportedByMutex.Unlock()

		var parseReporter domain.BasicAuth
		var err error
		if err = json.Unmarshal([]byte(job.outsource.AuthInfo()), &parseReporter); err == nil {
			if len(parseReporter.Username) > 0 {
				reportedBy = parseReporter.Username
				job.cachedReportedBy = reportedBy
			} else {
				err = fmt.Errorf("could not parse the reported from the source config")
			}
		}

		if err != nil {
			job.lstream.Send(log.Error("could not find the reporter from the out source config", err))
		}
	}

	return reportedBy
}

type assignedTicket struct {
	domain.Ticket
	assignee        string
	assignmentGroup string
}

func (t *assignedTicket) AssignedTo() *string {
	if len(t.assignee) > 0 {
		return &t.assignee
	} else {
		return nil
	}
}

func (t *assignedTicket) AssignmentGroup() (param *string) {
	if len(t.assignmentGroup) > 0 {
		return &t.assignmentGroup
	} else {
		return nil
	}
}

func (job *TicketingJob) getAssignmentInformation(tagsForDevice []domain.Tag, payload *vulnerabilityPayload) {
	var assignmentGroup, assignee string

	for _, rule := range job.assignmentRules {
		var match = true

		if rule.vulnTitleRegex != nil {
			match = rule.vulnTitleRegex.MatchString(payload.vuln.Name())
		}

		if match && rule.tagKey != nil {
			var found bool
			for _, deviceTag := range tagsForDevice {
				if deviceTag.TagKeyID() == iord(rule.TagKeyID()) {
					found = true
					match = rule.tagKeyRegex.MatchString(deviceTag.Value())
					break
				}
			}

			if !found {
				match = false
			}
		}

		if match {
			assignmentGroup = sord(rule.AssignmentGroup())
			assignee = sord(rule.Assignee())
			break // the rules are pulled highest-priority first, so the first match found should be the match taken
		}
	}

	if len(assignmentGroup) == 0 && len(sord(payload.ticket.IPAddress())) > 0 {
		// Handle the assignment using the data in config which is the scanner assignment for the IPs
		if ag, err := job.db.GetAssignmentGroupByIP(job.insource.SourceID(), job.config.OrganizationID(), sord(payload.ticket.IPAddress())); err == nil {
			if ag != nil && len(ag) > 0 {
				assignmentGroup = ag[0].GroupName()
			}
		} else {
			job.lstream.Send(log.Errorf(err, "error while loading assignment group for device [%s]", sord(payload.ticket.IPAddress())))
		}
	}

	payload.ticket = &assignedTicket{
		Ticket:          payload.ticket,
		assignee:        assignee,
		assignmentGroup: assignmentGroup,
	}
}

type assignmentRule struct {
	domain.AssignmentRules
	vulnTitleRegex *regexp.Regexp
	tagKeyRegex    *regexp.Regexp
	tagKey         domain.TagKey
}

func (job *TicketingJob) loadAssignmentRules() (assignmentRules []assignmentRule, err error) {
	assignmentRules = make([]assignmentRule, 0)

	var rules []domain.AssignmentRules
	if rules, err = job.db.GetAssignmentRulesByOrg(job.config.OrganizationID()); err == nil {

		for _, rule := range rules {
			var currentRule = assignmentRule{
				AssignmentRules: rule,
			}

			if rule.VulnTitleRegex() != nil {
				var regex *regexp.Regexp
				if regex, err = regexp.Compile(sord(rule.VulnTitleRegex())); err == nil {
					currentRule.vulnTitleRegex = regex
				} else {
					err = fmt.Errorf("error while compiling vuln title regex [%s] - %s", sord(rule.VulnTitleRegex()), err.Error())
					break
				}
			}

			if rule.TagKeyID() != nil {
				var tagKey domain.TagKey
				if tagKey, err = job.db.GetTagKeyByID(strconv.Itoa(iord(rule.TagKeyID()))); err == nil {
					if tagKey != nil {
						currentRule.tagKey = tagKey
					} else {
						err = fmt.Errorf("could not find a tag key for %d", iord(rule.TagKeyID()))
					}
				} else {
					err = fmt.Errorf("error while loading tag key - %s", err.Error())
					break
				}
			}

			if rule.TagKeyRegex() != nil {
				var regex *regexp.Regexp
				if regex, err = regexp.Compile(sord(rule.TagKeyRegex())); err == nil {
					currentRule.tagKeyRegex = regex
				} else {
					err = fmt.Errorf("error while compiling tag key regex [%s] - %s", sord(rule.TagKeyRegex()), err.Error())
					break
				}
			}

			if (rule.TagKeyID() != nil) != (rule.TagKeyRegex() != nil) { // != is equivalent to an xor operation, meaning if only one is set
				err = fmt.Errorf("entry exists where both TagKeyID and TagKeyRegex are not nil (xor)")
				break
			}

			assignmentRules = append(assignmentRules, currentRule)
		}

	} else {
		err = fmt.Errorf("error while loading assignment rules - %s", err.Error())
	}

	return assignmentRules, err
}

func (job *TicketingJob) getCVSSScore(vuln domain.Vulnerability) (score float32) {
	if job.OrgPayload.CVSSVersion == cvssVersion3 && vuln.CVSS3() != nil {
		score = *vuln.CVSS3()
	} else {
		score = vuln.CVSS2()
	}

	return score
}

func pushDetectionsToChannel(ctx context.Context, detections []domain.Detection) <-chan domain.Detection {
	out := make(chan domain.Detection)
	go func() {
		defer close(out)

		for _, detection := range detections {
			select {
			case <-ctx.Done():
				return
			case out <- detection:
			}
		}
	}()

	return out
}
