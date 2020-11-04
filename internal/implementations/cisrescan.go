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

	"github.com/nortonlifelock/aegis/internal/integrations"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/log"
)

// CISRescanJob implements the Job interface and pulls findings from Dome9 and creates tickets when applicable
type CISRescanJob struct {
	Payload *CISRescanPayload

	id         string
	orgCode    string
	orgPayload *OrgPayload
	catRules   []domain.CategoryRule

	payloadJSON string
	ctx         context.Context
	db          domain.DatabaseConnection
	lstream     log.Logger
	appconfig   domain.Config
	config      domain.JobConfig
	insource    domain.SourceConfig
	outsource   domain.SourceConfig
}

// CISRescanPayload holds information that dictates how the rescan is run, and on what account
// The BundleID points towards a bundle, which holds a series of rules
// The cloud account IDs points to the cloud account (e.g. AWS/Azure) that we which to test the rules against
type CISRescanPayload struct {
	BundleID        int      `json:"bundle_id"`
	CloudAccountIDs []string `json:"cloud_accounts"`
}

// buildPayload parses the BundleID from the job history Payload
func (job *CISRescanJob) buildPayload(pjson string) (err error) {
	job.Payload = &CISRescanPayload{}
	if err = json.Unmarshal([]byte(pjson), job.Payload); err == nil {
		if len(job.Payload.CloudAccountIDs) == 0 {
			err = fmt.Errorf("job payload did not include cloud account IDs")
		}
	}

	return err
}

// Process pulls findings from a particular bundle, and creates a ticket in the ticketing engine if one did not exist
func (job *CISRescanJob) Process(ctx context.Context, id string, appconfig domain.Config, db domain.DatabaseConnection, lstream log.Logger, payload string, jobConfig domain.JobConfig, inSource []domain.SourceConfig, outSource []domain.SourceConfig) (err error) {

	var ok bool
	if job.ctx, job.id, job.appconfig, job.db, job.lstream, job.payloadJSON, job.config, job.insource, job.outsource, ok = validInputs(ctx, id, appconfig, db, lstream, payload, jobConfig, inSource, outSource); ok {

		if err = job.buildPayload(job.payloadJSON); err == nil {

			if job.orgCode, job.orgPayload, err = getOrgInfo(job.db, job.config.OrganizationID()); err == nil {
				var engine integrations.TicketingEngine
				if engine, err = integrations.GetEngine(job.ctx, job.outsource.Source(), job.db, job.lstream, job.appconfig, job.outsource); err == nil {

					job.lstream.Send(log.Debug("Establishing connection to CIS scanner..."))

					var scanner integrations.CISScanner
					if scanner, err = integrations.GetCISScanner(job.ctx, job.insource.Source(), job.db, job.insource, job.appconfig, job.lstream); err == nil {

						if job.catRules, err = job.db.GetCategoryRules(job.config.OrganizationID(), job.insource.SourceID()); err == nil {
							wg := &sync.WaitGroup{}
							for _, cloudID := range job.Payload.CloudAccountIDs {
								wg.Add(1)
								go func(cloudID string) {
									defer handleRoutinePanic(job.lstream)
									defer wg.Done()

									var err error // error is scoped intentionally
									err = job.processBundleOnCloud(scanner, engine, job.Payload.BundleID, cloudID)
									if err != nil {
										job.lstream.Send(log.Errorf(err, "error while processing bundle ID [%d] for cloud account [%s]", job.Payload.BundleID, cloudID))
									}
								}(cloudID)
							}
							wg.Wait()
						} else {
							err = fmt.Errorf("error while loading category rules [%s]", err.Error())
						}
					} else {
						err = fmt.Errorf("error while building scanning connection - %v", err)
					}
				} else {
					err = fmt.Errorf("error while builing ticketing connection - %v", err)
				}
			} else {
				err = fmt.Errorf("error while loading organization information - %s", err.Error())
			}
		} else {
			err = fmt.Errorf("error while building payload - %s", err.Error())
		}
	} else {
		err = fmt.Errorf("input validation failed")
	}

	return err
}

func getOrgInfo(db domain.DatabaseConnection, orgID string) (orgCode string, orgPayload *OrgPayload, err error) {
	var torg domain.Organization
	if torg, err = db.GetOrganizationByID(orgID); err == nil {
		// Ensure there is only one return
		if torg != nil {
			orgCode = torg.Code()
			orgPayload = &OrgPayload{}
			err = json.Unmarshal([]byte(torg.Payload()), orgPayload)
			if err == nil {
				if len(orgPayload.Severities) > 0 {
					sort.Sort(orgPayload)
				} else {
					err = fmt.Errorf("no severities found in the organization payload")
				}
			}
		} else {
			err = fmt.Errorf("no organization found for [%s]", orgID)
		}
	}

	return orgCode, orgPayload, err
}

type findingTicketPair struct {
	finding domain.Ticket
	ticket  domain.Ticket
}

func (job *CISRescanJob) processBundleOnCloud(scanner integrations.CISScanner, engine integrations.TicketingEngine, bundleID int, cloudAccountID string) (err error) {

	var findings []domain.Finding
	findings, err = scanner.RescanBundle(bundleID, cloudAccountID)
	if err == nil {

		var tickets <-chan domain.Ticket
		tickets, err = engine.GetOpenTicketsByGroupID(job.insource.Source(), job.orgCode, cloudAccountID)
		if err == nil {

			assignmentInformation, err := job.db.GetCISAssignments(job.config.OrganizationID())
			if err != nil {
				job.lstream.Send(log.Errorf(err, "error while loading assignment group information"))
			}

			var findingsAsTickets = make([]domain.Ticket, 0)
			for index := range findings {
				finding := findings[index]

				findingTic := &FindingWrapper{
					finding,
					job,
					job.getAssignmentGroupForFinding(assignmentInformation, finding),
					getCategoryBasedOnRule(job.catRules, finding.VulnerabilityTitle(), "", ""),
				}

				if len(findingTic.DeviceID()) > 0 && len(findingTic.VulnerabilityID()) > 0 {
					findingsAsTickets = append(findingsAsTickets, findingTic)
				}
			}

			var assessmentID int
			if len(findings) > 0 {
				assessmentID = findings[0].ScanID()
			}

			processFindingsAndTickets(
				job.lstream,
				job.db,
				job.config.OrganizationID(),
				job.insource.SourceID(),
				engine,
				fanInChannel(tickets),
				findingsAsTickets,
				fmt.Sprintf("finding was NOT by %s in assessment [%d]", job.insource.Source(), assessmentID),
				fmt.Sprintf("finding still detected by %s in assessment [%d]", job.insource.Source(), assessmentID),
				func(ticket domain.Ticket) string {
					return fmt.Sprintf("%s;%s", ticket.DeviceID(), ticket.VulnerabilityID())
				},
				func(ticket domain.Ticket) bool {
					return strings.ToLower(sord(ticket.Priority())) != "low"
				},
			)
		}
	}

	return err
}

func processFindingsAndTickets(lstream log.Logger, db domain.DatabaseConnection, orgID string, sourceID string, engine integrations.TicketingEngine, tickets []domain.Ticket, findings []domain.Ticket, closingComment, updatingComment string, getKey func(ticket domain.Ticket) string, shouldCreateTicket func(ticket domain.Ticket) bool) (findingsWithoutTickets []domain.Ticket, ticketsWithoutFindings []domain.Ticket, ticketsWithFindings []findingTicketPair) {
	// findings without tickets need tickets created for  them
	findingsWithoutTickets = make([]domain.Ticket, 0)

	// tickets without findings may be closed
	ticketsWithoutFindings = make([]domain.Ticket, 0)

	// tickets with findings can have their last seen date updated and should be reopened
	ticketsWithFindings = make([]findingTicketPair, 0)

	var deviceIDToVulnIDToTicket = mapTicketsByDeviceIDVulnID(tickets, getKey)
	var deviceIDToVulnIDToFinding = mapTicketsByDeviceIDVulnID(findings, getKey)

	for _, finding := range findings {
		key := getKey(finding)

		if deviceIDToVulnIDToTicket[key] != nil {
			for _, tieTicketToFinding := range deviceIDToVulnIDToTicket[key] {
				ticketsWithFindings = append(ticketsWithFindings, findingTicketPair{
					finding: finding,
					ticket:  tieTicketToFinding,
				})
			}
		} else {
			findingsWithoutTickets = append(findingsWithoutTickets, finding)
		}
	}

	for _, ticket := range tickets {
		key := getKey(ticket)
		if deviceIDToVulnIDToFinding[key] == nil {
			ticketsWithoutFindings = append(ticketsWithoutFindings, ticket)
		}
	}

	updateTicketsAccordingToFindings(
		lstream,
		db,
		orgID,
		sourceID,
		engine,
		findingsWithoutTickets,
		ticketsWithFindings,
		ticketsWithoutFindings,
		closingComment,
		updatingComment,
		shouldCreateTicket,
	)
	findingDetectionCreation(lstream, db, orgID, sourceID, findings)

	return findingsWithoutTickets, ticketsWithoutFindings, ticketsWithFindings
}

func findingDetectionCreation(lstream log.Logger, db domain.DatabaseConnection, orgID string, sourceID string, findings []domain.Ticket) {
	var createdOrFoundDeviceID = make(map[string]bool)
	var vulnIDToVulnInfo = make(map[string]domain.VulnerabilityInfo)
	if vulnerableStatus, err := db.GetDetectionStatusByName(domain.Vulnerable); err == nil {
		if unknownOST, err := grabAndCreateOsType(db, unknown); err == nil {
			for _, finding := range findings {
				if len(finding.DeviceID()) > 0 && len(finding.VulnerabilityID()) > 0 {
					if !createdOrFoundDeviceID[finding.DeviceID()] {
						deviceInfo, err := db.GetDeviceInfoByAssetOrgID(finding.DeviceID(), orgID)
						if err == nil {

							if deviceInfo == nil {
								_, _, err = db.CreateDevice(
									finding.DeviceID(),
									sourceID,
									sord(finding.IPAddress()),
									sord(finding.HostName()),
									finding.CloudID(),
									sord(finding.MacAddress()),
									finding.GroupID(),
									orgID,
									sord(finding.OperatingSystem()),
									unknownOST.ID(),
									"",
								)

								if err == nil {
									createdOrFoundDeviceID[finding.DeviceID()] = true
								} else {
									lstream.Send(log.Errorf(err, "error while creating device [%s]", finding.DeviceID()))
								}
							} else {
								createdOrFoundDeviceID[finding.DeviceID()] = true
							}
						} else {
							lstream.Send(log.Errorf(err, "error while pulling device [%s]", finding.DeviceID()))
						}
					}

					var vulnInfo domain.VulnerabilityInfo
					var err error
					if vulnIDToVulnInfo[finding.VulnerabilityID()] == nil {
						vulnInfo, err = createAndGetVulnInfoForFinding(lstream, db, sourceID, finding)
					} else {
						vulnInfo = vulnIDToVulnInfo[finding.VulnerabilityID()]
					}

					if err == nil {
						var port, protocol string
						var portInt int
						if len(strings.Split(sord(finding.ServicePorts()), " ")) == 2 {
							port = strings.Split(sord(finding.ServicePorts()), " ")[0]
							protocol = strings.Split(sord(finding.ServicePorts()), " ")[1]
							portInt, _ = strconv.Atoi(port)
						}

						if detectionInfo, err := db.GetDetectionInfo(finding.DeviceID(), finding.VulnerabilityID(), portInt, protocol); err == nil {
							var ignoreID string
							ignore, err := db.HasIgnore(sourceID, finding.VulnerabilityID(), finding.DeviceID(), orgID, sord(finding.ServicePorts()), time.Now())
							if err != nil {
								lstream.Send(log.Errorf(err, "error while loading ignore entry [%s|%s]", finding.DeviceID(), finding.VulnerabilityID()))
							}
							if ignore != nil {
								ignoreID = ignore.ID()
							}

							if detectionInfo == nil {
								_, _, err = db.CreateDetection(
									orgID,
									sourceID,
									finding.DeviceID(),
									vulnInfo.ID(),
									ignoreID,
									time.Now(), // alert date
									tord1970(nil),
									tord1970(finding.LastChecked()),
									"",
									portInt,
									protocol,
									0,
									vulnerableStatus.ID(),
									0,
									tord1970(nil),
								)

								if err != nil {
									lstream.Send(log.Errorf(err, "error while creating detection for [%s|%s]", finding.DeviceID(), finding.VulnerabilityID()))
								}
							} else {
								if sord(detectionInfo.IgnoreID()) != ignoreID {
									_, _, err = db.UpdateDetectionIgnore(
										finding.DeviceID(),
										finding.VulnerabilityID(),
										portInt,
										protocol,
										ignoreID,
									)

									if err != nil {
										lstream.Send(log.Errorf(err, "error while updating detection for [%s|%s]", finding.DeviceID(), finding.VulnerabilityID()))
									}
								}

							}
						} else {
							lstream.Send(log.Errorf(err, "error while gathering detection info for [%s|%s]", finding.DeviceID(), finding.VulnerabilityID()))
						}
					}
				}
			}
		} else {
			lstream.Send(log.Errorf(err, ""))
		}
	} else {
		lstream.Send(log.Errorf(err, "error while loading vulnerable detection status"))
	}
}

func createAndGetVulnInfoForFinding(lstream log.Logger, db domain.DatabaseConnection, sourceID string, finding domain.Ticket) (vulnInfo domain.VulnerabilityInfo, err error) {
	vulnInfo, err = db.GetVulnInfoBySourceVulnID(finding.VulnerabilityID())
	if err == nil {
		if vulnInfo == nil {
			_, _, err = db.CreateVulnInfo(
				finding.VulnerabilityID(),
				sord(finding.VulnerabilityTitle()),
				sourceID,
				ford(finding.CVSS()),
				0.0,
				sord(finding.Description()),
				"",
				sord(finding.Solution()),
				"",
				sord(finding.Patchable()),
				sord(finding.Category()),
				"",
			)

			if err == nil {
				vulnInfo, err = db.GetVulnInfoBySourceVulnID(finding.VulnerabilityID())
				if err == nil && vulnInfo == nil {
					err = fmt.Errorf("attempted and failed to create vuln info for [%s]", finding.VulnerabilityID())
				}
			} else {
				lstream.Send(log.Errorf(err, "error while creating vuln info for [%s]", finding.VulnerabilityID()))
			}
		}
	} else {
		lstream.Send(log.Errorf(err, "error while pulling vuln info for [%s]", finding.VulnerabilityID()))
	}

	return vulnInfo, err
}

func updateTicketsAccordingToFindings(lstream log.Logger, db domain.DatabaseConnection, orgID string, sourceID string, engine integrations.TicketingEngine, findingsWithoutTickets []domain.Ticket, ticketsWithFindings []findingTicketPair, ticketsWithoutFindings []domain.Ticket, closingComment, updatingComment string, shouldCreateTicket func(ticket domain.Ticket) bool) {
	wg := &sync.WaitGroup{}
	wg.Add(3)
	go func() {
		defer handleRoutinePanic(lstream)
		defer wg.Done()
		createTicketsForUnticketedFindings(db, lstream, orgID, sourceID, engine, findingsWithoutTickets, shouldCreateTicket)
	}()
	go func() {
		defer handleRoutinePanic(lstream)
		defer wg.Done()
		updateTicketsWithStaleFindings(db, lstream, engine, ticketsWithFindings, updatingComment, orgID, sourceID)
	}()
	go func() {
		defer handleRoutinePanic(lstream)
		defer wg.Done()
		closeTicketsWithMissingFindings(lstream, engine, ticketsWithoutFindings, closingComment)
	}()
	wg.Wait()
}

func createTicketsForUnticketedFindings(db domain.DatabaseConnection, lstream log.Logger, orgID string, sourceID string, engine integrations.TicketingEngine, findings []domain.Ticket, shouldCreateTicket func(ticket domain.Ticket) bool) {
	wg := &sync.WaitGroup{}
	for index := range findings {
		if shouldCreateTicket(findings[index]) {
			wg.Add(1)
			go func(finding domain.Ticket) {
				defer handleRoutinePanic(lstream)
				defer wg.Done()

				ignore, err := db.HasIgnore(sourceID, finding.VulnerabilityID(), finding.DeviceID(), orgID, sord(finding.ServicePorts()), time.Now())
				if err != nil {
					lstream.Send(log.Errorf(err, "error while loading ignore entry [%s|%s]", finding.DeviceID(), finding.VulnerabilityID()))
				}

				if ignore == nil {

					_, sourceKey, err := engine.CreateTicket(finding)
					if err == nil {
						lstream.Send(log.Infof("Created ticket [%s] for device [%s] with vulnerability [%s]", sourceKey, finding.DeviceID(), finding.VulnerabilityID()))
					} else {
						lstream.Send(log.Errorf(err, "error while creating ticket for device [%s] with vulnerability [%s]", finding.DeviceID(), finding.VulnerabilityID()))
					}
				} else {
					lstream.Send(log.Infof("SKIPPING ticket for [%s|%s] as it has an ignore entry", finding.DeviceID(), finding.VulnerabilityID()))
				}

			}(findings[index])
		}

	}
	wg.Wait()
}

// the assignment groups mapping information is stored in the database. The following fields show the hierarchy of the prioritization for finding an assignment group
// CloudAccountID->BundleID->RuleRegex->RuleHash
// The only required field is the cloud account ID. The rest of the fields may be nil. If the other fields are non-nil, and their values don't match that of the finding, the match is not considered
func (job *CISRescanJob) getAssignmentGroupForFinding(assignmentInformation []domain.CISAssignments, finding domain.Finding) (assignmentGroup string) {
	if assignmentInformation != nil {

		// currentDepth tracks how much of a match the current assignment group match is. a higher depth means a greater match
		var currentDepth = -1

		for _, info := range assignmentInformation {

			// match contains a value of true as long as none of the specifications are violated
			// if a specification is violated, the assignmentInformation is not taken into account
			var match = true
			var depthOfMatch = 0

			if len(sord(info.CloudAccountID())) > 0 {
				if sord(info.CloudAccountID()) == finding.AccountID() {
					// matching on the cloud account id alone has a match value of 1
					depthOfMatch = 1
				} else {
					match = false
				}
			}

			// a bundle id match implies a greater match than the cloud account
			if len(sord(info.BundleID())) > 0 {
				if sord(info.BundleID()) == finding.BundleID() {
					depthOfMatch = 2
				} else {
					match = false
				}
			}

			// the rule name matching a regex implies a greater match than a bundle id
			if len(sord(info.RuleRegex())) > 0 {
				if valid, err := regexp.Match(sord(info.RuleRegex()), []byte(finding.VulnerabilityTitle())); err == nil {
					if valid {
						depthOfMatch = 3
					} else {
						match = false
					}
				} else {
					job.lstream.Send(log.Errorf(err, "error while compiling regex [%v]", *info.RuleRegex()))
				}
			}

			// a specific rule hash implies a greater match than a regex in the rule name
			if len(sord(info.RuleHash())) > 0 {
				if sord(info.RuleHash()) == finding.ID() {
					depthOfMatch = 4
				} else {
					match = false
				}
			}

			if match {

				// the current iteration contained the most closely specified assignment group information
				if depthOfMatch > currentDepth {
					currentDepth = depthOfMatch
					assignmentGroup = info.AssignmentGroup()
				}
			}
		}
	}

	return assignmentGroup
}

func (job *CISRescanJob) calculateSLAForCISTicket(severity string) (due time.Time) {
	// the organizations are sorted by severity, so the highest severity should be the last one in the list
	var highestSeverity = job.orgPayload.Severities[len(job.orgPayload.Severities)-1]
	for _, orgSeverity := range job.orgPayload.Severities {
		if strings.ToLower(orgSeverity.Name) == strings.ToLower(severity) {
			due = time.Now().AddDate(0, 0, orgSeverity.Duration)
			break
		}
	}

	if due.IsZero() {
		due = time.Now().AddDate(0, 0, highestSeverity.Duration)
	}

	return due
}

type staleTicket struct {
	domain.Ticket
	engine integrations.TicketingEngine
}

// LastChecked overrides the domain.Ticket method
func (t *staleTicket) LastChecked() *time.Time {
	val := time.Now()
	return &val
}

// Status opens the stale ticket if it's in resolved-remediated
func (t *staleTicket) Status() (val *string) {
	val = t.Ticket.Status()

	if strings.ToLower(sord(t.Ticket.Status())) == strings.ToLower(t.engine.GetStatusMap(domain.StatusResolvedRemediated)) {
		s := t.engine.GetStatusMap(domain.StatusReopened)
		val = &s
	}

	return val
}

func updateTicketsWithStaleFindings(db domain.DatabaseConnection, lstream log.Logger, engine integrations.TicketingEngine, pairs []findingTicketPair, updatingComment string, orgID, sourceID string) {
	wg := &sync.WaitGroup{}
	for index := range pairs {
		wg.Add(1)
		go func(pair findingTicketPair) {
			defer handleRoutinePanic(lstream)
			defer wg.Done()

			_, _, err := engine.UpdateTicket(
				&staleTicket{
					pair.ticket,
					engine,
				},
				updatingComment,
			)
			if err == nil {
				lstream.Send(log.Infof("finding for %s still detected", pair.ticket.Title()))
			} else {
				lstream.Send(log.Errorf(err, "error while updating ticket [%s]", pair.ticket.Title()))
			}
		}(pairs[index])
	}
	wg.Wait()
}

type closedTicket struct {
	domain.Ticket
	engine integrations.TicketingEngine
}

// Status overrides the domain.Ticket method
func (c closedTicket) Status() *string {
	val := c.engine.GetStatusMap(domain.StatusClosedRemediated)
	return &val
}

func closeTicketsWithMissingFindings(lstream log.Logger, engine integrations.TicketingEngine, tickets []domain.Ticket, closingComment string) {
	wg := &sync.WaitGroup{}
	for index := range tickets {
		wg.Add(1)
		go func(ticket domain.Ticket) {
			defer handleRoutinePanic(lstream)
			defer wg.Done()

			_, _, err := engine.UpdateTicket(
				closedTicket{
					ticket,
					engine,
				},
				closingComment,
			)

			if err == nil {
				lstream.Send(log.Infof("finding for %s was NOT in assessment, closing ticket...", ticket.Title()))
			} else {
				lstream.Send(log.Errorf(err, "error while updating ticket [%s]", ticket.Title()))
			}
		}(tickets[index])
	}
	wg.Wait()
}

func mapTicketsByDeviceIDVulnID(tickets []domain.Ticket, getKey func(ticket domain.Ticket) string) (entityIDToRuleHashToTicket map[string][]domain.Ticket) {
	entityIDToRuleHashToTicket = make(map[string][]domain.Ticket)

	for _, ticket := range tickets {
		key := getKey(ticket)
		if entityIDToRuleHashToTicket[key] == nil {
			entityIDToRuleHashToTicket[key] = make([]domain.Ticket, 0)
		}

		entityIDToRuleHashToTicket[key] = append(entityIDToRuleHashToTicket[key], ticket)
	}

	return entityIDToRuleHashToTicket
}

// fanInChannel is useful because we want to reuse the ticket information, so we store it in a slice
func fanInChannel(in <-chan domain.Ticket) (out []domain.Ticket) {
	out = make([]domain.Ticket, 0)
	for {
		if ticket, ok := <-in; ok {
			out = append(out, ticket)
		} else {
			break
		}
	}

	return out
}

// FindingWrapper implements the domain.Ticket interface so the finding may be converted into a ticket
type FindingWrapper struct {
	domain.Finding
	job *CISRescanJob
	ag  string
	cat string
}

// AlertDate returns the AlertDate of the ticket
func (wrapper *FindingWrapper) AlertDate() (param *time.Time) {
	return
}

// AssignedTo returns the AssignedTo of the ticket
func (wrapper *FindingWrapper) AssignedTo() (param *string) {
	return
}

func (wrapper *FindingWrapper) Category() (param *string) {
	return &wrapper.cat
}

// AssignmentGroup returns the AssignmentGroup of the ticket
func (wrapper *FindingWrapper) AssignmentGroup() (param *string) {
	return &wrapper.ag
}

// CERF returns the CERF of the ticket
func (wrapper *FindingWrapper) CERF() (param string) {
	return
}

// ExceptionExpiration returns the ExceptionExpiration of the ticket
func (wrapper *FindingWrapper) ExceptionExpiration() (param time.Time) {
	return
}

// CVEReferences returns the CVEReferences of the ticket
func (wrapper *FindingWrapper) CVEReferences() (param *string) {
	return
}

// CVSS returns the CVSS of the ticket
func (wrapper *FindingWrapper) CVSS() (param *float32) {
	return
}

// CloudID returns the CloudID of the ticket
func (wrapper *FindingWrapper) CloudID() (param string) {
	return
}

// Configs returns the Configs of the ticket
func (wrapper *FindingWrapper) Configs() (param string) {
	return
}

// CreatedDate returns the CreatedDate of the ticket
func (wrapper *FindingWrapper) CreatedDate() (param *time.Time) {
	return
}

// DBCreatedDate returns the DBCreatedDate of the ticket
func (wrapper *FindingWrapper) DBCreatedDate() (param time.Time) {
	return
}

// DBUpdatedDate returns the DBUpdatedDate of the ticket
func (wrapper *FindingWrapper) DBUpdatedDate() (param *time.Time) {
	return
}

// Description returns the Description of the ticket
func (wrapper *FindingWrapper) Description() (param *string) {
	val := wrapper.Finding.String()
	return &val
}

// DeviceID returns the DeviceID of the ticket
func (wrapper *FindingWrapper) DeviceID() (param string) {
	return wrapper.Finding.DeviceID()
}

// DueDate returns the DueDate of the ticket
func (wrapper *FindingWrapper) DueDate() (param *time.Time) {
	val := wrapper.job.calculateSLAForCISTicket(wrapper.Finding.Priority())
	return &val
}

// GroupID returns the GroupID of the ticket
func (wrapper *FindingWrapper) GroupID() (param string) {
	return wrapper.Finding.AccountID()
}

// HostName returns the HostName of the ticket
func (wrapper *FindingWrapper) HostName() (param *string) {
	return
}

// ID returns the ID of the ticket
func (wrapper *FindingWrapper) ID() (param int) {
	return
}

// IPAddress returns the IPAddress of the ticket
func (wrapper *FindingWrapper) IPAddress() (param *string) {
	return
}

// Labels returns the Labels of the ticket
func (wrapper *FindingWrapper) Labels() (param *string) {
	return
}

// LastChecked returns the LastChecked of the ticket
func (wrapper *FindingWrapper) LastChecked() (param *time.Time) {
	return
}

// MacAddress returns the MacAddress of the ticket
func (wrapper *FindingWrapper) MacAddress() (param *string) {
	return
}

// MethodOfDiscovery returns the MethodOfDiscovery of the ticket
func (wrapper *FindingWrapper) MethodOfDiscovery() (param *string) {
	val := wrapper.job.insource.Source()
	return &val
}

// OSDetailed returns the OSDetailed of the ticket
func (wrapper *FindingWrapper) OSDetailed() (param *string) {
	return
}

// OperatingSystem returns the OperatingSystem of the ticket
func (wrapper *FindingWrapper) OperatingSystem() (param *string) {
	return
}

// OrgCode returns the OrgCode of the ticket
func (wrapper *FindingWrapper) OrgCode() (param *string) {
	return &wrapper.job.orgCode
}

// OrganizationID returns the OrganizationID of the ticket
func (wrapper *FindingWrapper) OrganizationID() (param string) {
	return
}

func (wrapper *FindingWrapper) Patchable() (param *string) {
	return nil
}

// Priority returns the Priority of the ticket
func (wrapper *FindingWrapper) Priority() (param *string) {
	val := wrapper.Finding.Priority()
	return &val
}

// Project returns the Project of the ticket
func (wrapper *FindingWrapper) Project() (param *string) {
	return
}

// ReportedBy returns the ReportedBy of the ticket
func (wrapper *FindingWrapper) ReportedBy() (param *string) {
	return
}

// ResolutionDate returns the ResolutionDate of the ticket
func (wrapper *FindingWrapper) ResolutionDate() (param *time.Time) {
	return
}

// ResolutionStatus returns the ResolutionStatus of the ticket
func (wrapper *FindingWrapper) ResolutionStatus() (param *string) {
	return
}

// ScanID returns the ScanID of the ticket
func (wrapper *FindingWrapper) ScanID() (param int) {
	return wrapper.Finding.ScanID()
}

// ServicePorts returns the ServicePorts of the ticket
func (wrapper *FindingWrapper) ServicePorts() (param *string) {
	return
}

// Solution returns the Solution of the ticket
func (wrapper *FindingWrapper) Solution() (param *string) {
	return
}

// Status returns the Status of the ticket
func (wrapper *FindingWrapper) Status() (param *string) {
	return
}

// Summary returns the Summary of the ticket
func (wrapper *FindingWrapper) Summary() (param *string) {
	val := wrapper.Finding.Summary()
	return &val
}

// TicketType returns the TicketType of the ticket
func (wrapper *FindingWrapper) TicketType() (param *string) {
	val := "Request"
	return &val
}

// Title returns the Title of the ticket
func (wrapper *FindingWrapper) Title() (param string) {
	return
}

// UpdatedDate returns the UpdatedDate of the ticket
func (wrapper *FindingWrapper) UpdatedDate() (param *time.Time) {
	return
}

// VendorReferences returns the VendorReferences of the ticket
func (wrapper *FindingWrapper) VendorReferences() (param *string) {
	return
}

// VulnerabilityID returns the VulnerabilityID of the ticket
func (wrapper *FindingWrapper) VulnerabilityID() (param string) {
	return wrapper.Finding.ID()
}

// VulnerabilityTitle returns the VulnerabilityTitle of the ticket
func (wrapper *FindingWrapper) VulnerabilityTitle() (param *string) {
	val := wrapper.Finding.VulnerabilityTitle()
	return &val
}

func (wrapper *FindingWrapper) SystemName() (param *string) {
	return nil
}

func (wrapper *FindingWrapper) OWASP() (param *string) {
	return nil
}

func (wrapper *FindingWrapper) ExceptionDate() (param *time.Time) {
	return nil
}

func (wrapper *FindingWrapper) ApplicationName() (param *string) {
	return nil
}

func (wrapper *FindingWrapper) HubProjectName() (param *string) {
	return
}

func (wrapper *FindingWrapper) HubProjectVersion() (param *string) {
	return
}

func (wrapper *FindingWrapper) HubSeverity() (param *string) {
	return
}

func (wrapper *FindingWrapper) ComponentName() (param *string) {
	return
}

func (wrapper *FindingWrapper) ComponentVersion() (param *string) {
	return
}

func (wrapper *FindingWrapper) PolicyRule() (param *string) {
	return
}

func (wrapper *FindingWrapper) PolicySeverity() (param *string) {
	return
}
