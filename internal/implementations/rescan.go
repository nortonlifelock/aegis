package implementations

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/integrations"
	"github.com/nortonlifelock/log"
	"github.com/pkg/errors"
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
	inactiveKernelComment = "Remediation successful; Vulnerable Kernel is no longer active"
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
							err = job.createAndMonitorScan(ticketToMatch(tickets[i:i+batchSize], job.Payload.Group), getTicketTitles(tickets[i:i+batchSize]))
						} else {
							err = job.createAndMonitorScan(ticketToMatch(tickets[i:], job.Payload.Group), getTicketTitles(tickets[i:]))
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

func getTicketTitles(tickets []domain.Ticket) (titles []string) {
	titles = make([]string, 0)

	for _, ticket := range tickets {
		titles = append(titles, ticket.Title())
	}

	return titles
}

// creates a connection to the scanning engine in order to create a scan for the ips and vulnerabilities
// after the scan is created, the method monitors the status of the scan in the database and does not end
// until the scan leaves the queued status
func (job *RescanJob) createAndMonitorScan(matches []domain.Match, tickets []string) (err error) {
	// Initialize scanner object
	var vscanner integrations.Vscanner
	if vscanner, err = integrations.NewVulnScanner(job.ctx, job.inSource.Source(), job.db, job.lstream, job.appConfig, job.inSource); err == nil {

		job.lstream.Send(log.Debugf("Scanning Engine Connection Initialized. Loading Tickets [%s]", strings.Join(tickets, ",")))

		var scans <-chan domain.Scan

		if job.Payload.Type == domain.RescanDecommission {
			scans = vscanner.Discovery(job.ctx, matches)
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
							job.lstream.Send(log.Infof("New scan created [ID: %v] for [%v] tickets", scan.ID(), len(matches)))

							// to be used in the Payload for the scan summary
							scanClosePayload := job.createScanClosePayload(scan, matches, tickets)

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

// takes an input slice of tickets, and returns a list of ip addresses and vulnerabilities associated with those tickets
func (job *RescanJob) loadIPAddressesAndVulnerabilitiesFromTickets(tickets []domain.Ticket) (ipAddresses []string, vulnerabilities []string, err error) {
	var ipsLoaded = make(map[string]bool)
	ipAddresses = make([]string, 0)
	vulnerabilities = make([]string, 0)

	for index := range tickets {
		select {
		case <-job.ctx.Done():
			return
		default:

			var ipAddress string
			if ipAddress, err = job.getIPAddressFromTicket(tickets[index]); err == nil {
				// TODO: We need to eventually verify the formatting of this IP to ensure that it's actually an IP address
				if len(ipAddress) > 0 {

					// Add the device to the list to be scanned if not already there
					if !ipsLoaded[ipAddress] {

						ipAddresses = append(ipAddresses, ipAddress)
						ipsLoaded[ipAddress] = true
					}

					vulnerabilities = append(vulnerabilities, tickets[index].VulnerabilityID())

				} else {
					job.lstream.Send(log.Errorf(err, "Invalid Ip Address for ticket [%s]", tickets[index].Title()))
				}
			} else {
				job.lstream.Send(log.Errorf(err, "Error while loading Ip Address for ticket [%s]", tickets[index].Title()))
			}
		}
	}

	return ipAddresses, vulnerabilities, err
}

func (job *RescanJob) createScanClosePayload(scan domain.Scan, matches []domain.Match, tickets []string) *ScanClosePayload {
	var devices = make([]string, 0)
	for _, match := range matches {
		devices = append(devices, match.IP())
	}

	scanClosePayload := &ScanClosePayload{}
	scanClosePayload.Scan = scan
	scanClosePayload.ScanID = scan.ID()
	scanClosePayload.Tickets = tickets
	scanClosePayload.Devices = devices
	scanClosePayload.Group = job.Payload.Group
	scanClosePayload.Type = job.Payload.Type

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

func ticketToMatch(t []domain.Ticket, groupID string) (m []domain.Match) {
	m = make([]domain.Match, 0)
	for _, tic := range t {
		m = append(m, matchTicket{
			t:       tic,
			groupID: groupID,
		})
	}
	return m
}

type matchTicket struct {
	t       domain.Ticket
	groupID string
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
	return m.groupID
}
