package implementations

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nortonlifelock/aegis/internal/integrations"
	"github.com/nortonlifelock/aegis/pkg/domain"
	"github.com/nortonlifelock/aegis/pkg/log"
	"github.com/pkg/errors"
	"strings"
)

// RescanPayload is used to parse the Payload from the job history table. The Payload is generated automatically from the rescan queue job which creates the
// job history for the rescan job
type RescanPayload struct {
	Group   string   `json:"group"`
	Tickets []string `json:"tickets"`
	Type    string   `json:"type"`
}

// RescanJob implements the Job interface required to run the job
type RescanJob struct {
	Payload *RescanPayload
	State   string

	id          string
	payloadJSON string
	ctx         context.Context
	db          domain.DatabaseConnection
	lstream     log.Logger
	appConfig   domain.Config
	config      domain.JobConfig
	inSource    domain.SourceConfig
	outSource   domain.SourceConfig
}

const (
	reopenComment         = "Remediation Failed"
	closeComment          = "Remediation successful; Vulnerability not found"
	inactiveKernelComment = "Remediation successful"
	potentialComment      = "Closing ticket; Potential vulnerability"
)

// buildPayload loads the Payload from the job history into the RescanPayload
func (job *RescanJob) buildPayload(pjson string) (err error) {
	// Parse json to RescanPayload
	// Verify pJson length > 0
	if len(pjson) > 0 {

		job.Payload = &RescanPayload{}

		err = json.Unmarshal([]byte(pjson), job.Payload)
	} else {
		err = errors.New("Payload length is 0")
	}

	return err
}

// Process loads tickets that are in a status that requires rescanning. The job kicks off a rescan for the tickets using the scanning engine
func (job *RescanJob) Process(ctx context.Context, id string, appconfig domain.Config, db domain.DatabaseConnection, lstream log.Logger, payload string, jobConfig domain.JobConfig, inSource []domain.SourceConfig, outSource []domain.SourceConfig) (err error) {

	var ok bool
	if job.ctx, job.id, job.appConfig, job.db, job.lstream, job.payloadJSON, job.config, job.inSource, job.outSource, ok = validInputs(ctx, id, appconfig, db, lstream, payload, jobConfig, inSource, outSource); ok {

		if err = job.buildPayload(job.payloadJSON); err == nil {

			// Initialize ticketing object
			var ticketing integrations.TicketingEngine
			if ticketing, err = integrations.GetEngine(job.ctx, job.outSource.Source(), job.db, job.lstream, job.appConfig, job.outSource); err == nil {
				job.lstream.Send(log.Debug("Ticketing Engine Connection Initialized"))

				var tickets []domain.Ticket
				if tickets, err = loadTickets(job.lstream, ticketing, job.Payload.Tickets); err == nil {
					job.lstream.Send(log.Infof("Tickets Loaded [%s]", strings.Join(job.Payload.Tickets, ",")))

					const batchSize = 400
					for i := 0; i < len(tickets); i += batchSize {

						if i+batchSize <= len(tickets) {
							err = job.createAndMonitorScan(job.ticketToMatch(tickets[i:i+batchSize], job.Payload.Group), tickets[i:i+batchSize])
						} else {
							err = job.createAndMonitorScan(job.ticketToMatch(tickets[i:], job.Payload.Group), tickets[i:])
						}

						if err != nil {
							job.lstream.Send(log.Errorf(err, "error while creating scan"))
						}
					}

				} else {
					job.lstream.Send(log.Errorf(err, "Error while loading tickets from JIRA"))
				}
			} else {
				job.lstream.Send(log.Error("Ticketing Engine Connection Cannot Be Initialized", err))
			}
		} else {
			err = fmt.Errorf("error while building payload - %s", err.Error())
		}
	} else {
		err = fmt.Errorf("input validation failed")
	}

	return err
}

// creates a connection to the scanning engine in order to create a scan for the ips and vulnerabilities
// after the scan is created, the method monitors the status of the scan in the database and does not end
// until the scan leaves the queued status
func (job *RescanJob) createAndMonitorScan(matches []domain.Match, tickets []domain.Ticket) (err error) {
	// Initialize scanner object
	var vscanner integrations.Vscanner
	if vscanner, err = integrations.NewVulnScanner(job.ctx, job.inSource.Source(), job.db, job.lstream, job.appConfig, job.inSource); err == nil {
		job.lstream.Send(log.Debugf("scanning engine connection initialized"))

		var scans <-chan domain.Scan
		if job.Payload.Type == domain.RescanDecommission {
			if _, _, canCreate := canCreateCloudDecommJob(job.db, job.lstream, job.config.OrganizationID(), job.Payload.Group); canCreate {
				createCloudDecommissionJob(job.id, job.db, job.lstream, job.config.OrganizationID(), job.Payload.Group, getIPsFromMatches(matches), getTitlesFromTickets(tickets))
				out := make(chan domain.Scan)
				close(out)
				scans = out
			} else {
				scans = vscanner.Discovery(job.ctx, matches)
			}
		} else {
			scans, err = vscanner.Scan(job.ctx, matches)
		}

		if err == nil {

			func() {
				for {
					select {
					case <-job.ctx.Done():
						return
					case scan, ok := <-scans:
						if ok {
							ticketsToRescan := getTicketsBelongingToMatches(scan.Matches(), tickets)
							job.lstream.Send(log.Infof("New scan created [ID: %v] for [%v] tickets", scan.ID(), len(ticketsToRescan)))

							if len(ticketsToRescan) != len(scan.Matches()) {
								job.lstream.Send(log.Warningf(err, "mismatch between match length and ticket length [%d != %d]", len(ticketsToRescan), len(scan.Matches())))
							}

							// to be used in the Payload for the scan summary
							scanClosePayload := job.createScanClosePayload(scan, scan.Matches(), ticketsToRescan)

							var bytePayload []byte
							bytePayload, err = json.Marshal(scanClosePayload)

							if err == nil {

								// Log the scan Id in the database for monitoring by the scan sync this
								if _, _, err = job.db.CreateScanSummary(
									job.inSource.SourceID(),
									job.inSource.ID(),
									job.config.OrganizationID(),
									scan.ID(),
									domain.ScanQUEUED,
									string(bytePayload),
									job.id,
								); err != nil {
									job.lstream.Send(log.Errorf(err, "Error Saving Scan Summary Details for Scan [%v]", scan.ID()))
								}
							} else {
								job.lstream.Send(log.Errorf(err, "Error while marshalling Payload for ScanCloseJob"))
							}
						} else {
							return
						}
					}
				}
			}()
		} else {
			job.lstream.Send(log.Errorf(err, "Error while creating scan"))
		}
	} else {
		job.lstream.Send(log.Errorf(err, "Error when connecting to scan engine"))
	}

	return err
}

func (job *RescanJob) createScanClosePayload(scan domain.Scan, matches []domain.Match, tickets []string) *ScanClosePayload {
	var devices = make([]string, 0)
	for _, match := range matches {
		devices = append(devices, match.IP())
	}

	scanClosePayload := &ScanClosePayload{
		RescanPayload: RescanPayload{
			Group:   job.Payload.Group,
			Tickets: tickets,
			Type:    job.Payload.Type,
		},
		Scan:    scan,
		Devices: devices,
		ScanID:  scan.ID(),
	}

	return scanClosePayload
}

func (job *RescanJob) getIPAddressFromTicket(ticket domain.Ticket) (ip string, err error) {
	if ticket != nil {
		if ticket.IPAddress() != nil {

			var ticketIps = strings.Split(*ticket.IPAddress(), ",")
			if ticketIps != nil && len(ticketIps) > 0 {
				ip = ticketIps[0]
			} else {
				err = fmt.Errorf("%s did not have an ip address listed", ticket.Title())
			}
		} else {
			err = fmt.Errorf("%s did not have an ip address listed", ticket.Title())
		}
	} else {
		err = errors.New("null ticket passed to getIPAddressFromTicket")
	}

	return ip, err
}

func (job *RescanJob) ticketToMatch(t []domain.Ticket, groupID string) (m []domain.Match) {
	m = make([]domain.Match, 0)
	deviceIDToInstanceID := make(map[string]*string)
	deviceIDToRegion := make(map[string]*string)
	for _, tic := range t {

		// if this is the first time you've seen the device, try and load the relevant cloud information (if it exists)
		// if it does not exist or it is not a cloud device, this will still only execute a single time
		if deviceIDToInstanceID[tic.DeviceID()] == nil {
			device, err := job.db.GetDeviceInfoByAssetOrgID(tic.DeviceID(), job.config.OrganizationID())
			if err == nil {
				instanceID := sord(device.InstanceID())
				region := sord(device.Region())
				deviceIDToInstanceID[tic.DeviceID()] = &instanceID
				deviceIDToRegion[tic.DeviceID()] = &region
			} else {
				val := ""
				deviceIDToInstanceID[tic.DeviceID()] = &val
				deviceIDToRegion[tic.DeviceID()] = &val
				job.lstream.Send(log.Errorf(err, "error while loading device match information for [%s]", tic.DeviceID()))
			}
		}

		m = append(m, matchTicket{
			t:          tic,
			groupID:    groupID,
			instanceID: sord(deviceIDToInstanceID[tic.DeviceID()]),
			region:     sord(deviceIDToRegion[tic.DeviceID()]),
		})
	}
	return m
}

type matchTicket struct {
	t          domain.Ticket
	groupID    string
	instanceID string
	region     string
}

// IP returns the IP contained within the ticket
func (m matchTicket) IP() string {
	return sord(m.t.IPAddress())
}

// Device returns the device ID of the ticket
func (m matchTicket) Device() string {
	return m.t.DeviceID()
}

// Vulnerability returns the vulnerability ID contained within the ticket
func (m matchTicket) Vulnerability() string {
	return strings.Split(m.t.VulnerabilityID(), domain.VulnPathConcatenator)[0]
}

// GroupID returns the group that the ticket belongs to. This is used to create the scan within the scanning engine
func (m matchTicket) GroupID() string {
	if len(m.t.GroupID()) > 0 {
		return m.t.GroupID()
	} else {
		return m.groupID
	}
}

func (m matchTicket) Region() string {
	return m.region
}

func (m matchTicket) InstanceID() string {
	return m.instanceID
}

func getTitlesFromTickets(tickets []domain.Ticket) (titles []string) {
	titles = make([]string, 0)

	for index := range tickets {
		titles = append(titles, tickets[index].Title())
	}

	return titles
}

func getIPsFromMatches(matches []domain.Match) (ips []string) {
	ips = make([]string, 0)
	seen := make(map[string]bool)
	for _, match := range matches {
		if !seen[match.IP()] {
			seen[match.IP()] = true
			ips = append(ips, match.IP())
		}
	}
	return ips
}

func getIPsFromTickets(matches []domain.Ticket) (ips []string) {
	ips = make([]string, 0)
	seen := make(map[string]bool)
	for _, match := range matches {
		if len(sord(match.IPAddress())) > 0 {
			if !seen[sord(match.IPAddress())] {
				seen[sord(match.IPAddress())] = true
				ips = append(ips, sord(match.IPAddress()))
			}
		}
	}
	return ips
}

func getTicketsBelongingToMatches(matches []domain.Match, tickets []domain.Ticket) (ticketsBelongingToMatches []string) {
	ticketsBelongingToMatches = make([]string, 0)

	for _, match := range matches {
		for _, ticket := range tickets {
			if ticket.DeviceID() == match.Device() &&
				strings.Contains(ticket.VulnerabilityID(), match.Vulnerability()) && // strings.Contains because there might be a version at the end of the ticket vuln ID
				(ticket.GroupID() == match.GroupID() || len(ticket.GroupID()) == 0) {
				ticketsBelongingToMatches = append(ticketsBelongingToMatches, ticket.Title())
				break
			}
		}
	}

	return ticketsBelongingToMatches
}
