package connector

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/benjivesterby/validator"
	"github.com/pkg/errors"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/nexpose"
	"github.com/nortonlifelock/funnel"
	"github.com/nortonlifelock/log"
	"github.com/nortonlifelock/ttl"
)

// TODO: Setup default sorting for the api so that it correctly returns the data

type client interface {
	Do(r *http.Request) (*http.Response, error)
}

// Connection is the object that is used to interact with the Nexpose API
type Connection struct {
	ctx             context.Context
	api             *nexpose.Session
	settings        *Payload
	authentication  *domain.Host
	logger          log.Logger
	vulnerabilities ttl.Cache
	engines         ttl.Cache
}

// Discovery is used to create scans in Nexpose that are intended to find dead/decommissioned hosts
func (conn *Connection) Discovery(ctx context.Context, matches []domain.Match) <-chan domain.Scan {
	out := make(chan domain.Scan)

	go func(out chan<- domain.Scan) {
		defer handleRoutinePanic(conn.logger)
		defer close(out)
		ctx = ctxtest(ctx)

		if matches != nil {

			if strings.Count(conn.settings.DiscoveryNameFormat, "%s") == 1 {
				var groupToMatch = make(map[string][]domain.Match)

				for _, detection := range matches {
					if groupToMatch[detection.GroupID()] == nil {
						groupToMatch[detection.GroupID()] = make([]domain.Match, 0)
					}

					groupToMatch[detection.GroupID()] = append(groupToMatch[detection.GroupID()], detection)
				}

				wg := &sync.WaitGroup{}
				for groupToRescan, detectionsToRescan := range groupToMatch {
					wg.Add(1)
					go func(groupToRescan string, detectionsToRescan []domain.Match) {
						defer handleRoutinePanic(conn.logger)
						defer wg.Done()

						if templateID, err := conn.api.DuplicateScanTemplate(conn.settings.DiscoveryTemplate); err == nil {
							scanName := fmt.Sprintf(conn.settings.ScanNameFormat)
							var scan *Scan
							scan, err = conn.createScanForDetections(matches, groupToRescan, scanName, templateID, strconv.Itoa(conn.settings.RescanSite))

							if err == nil {
								select {
								case <-ctx.Done():
									return
								case out <- scan:
								}
							} else {
								conn.logger.Send(log.Errorf(err, "failed to create scan"))
							}
						} else {
							conn.logger.Send(log.Errorf(err, "error while creating copy of discovery scan template"))
						}
					}(groupToRescan, detectionsToRescan)
				}
				wg.Wait()
			} else {
				conn.logger.Send(log.Error("the discovery name template must contain a single %s placeholder", nil))
			}
		} else {
			conn.logger.Send(log.Errorf(nil, "nil match slice passed for discovery scan"))
		}
	}(out)

	return out
}

// Detections returns the vulnerability detections for a set of sites or asset groups
func (conn *Connection) Detections(ctx context.Context, groupsIDs []string) (<-chan domain.Detection, error) {
	var detections = make(chan domain.Detection)

	if len(groupsIDs) > 0 {
		ctx = ctxtest(ctx)

		go func(ctx context.Context, detections chan<- domain.Detection) {
			defer handleRoutinePanic(conn.logger)
			defer close(detections)

			var wg = &sync.WaitGroup{}

			for _, id := range groupsIDs {
				select {
				case <-ctx.Done():
					return
				default:
					wg.Add(1)
					go func(id string) {
						defer wg.Done()
						defer handleRoutinePanic(conn.logger)
						conn.loadSiteDetections(ctx, detections, id)
					}(id)
				}
			}

			wg.Wait()
		}(ctx, detections)
	} else {
		var err = errors.New("empty list of scan groups")
		conn.logger.Send(log.Error(err.Error(), err))
	}

	return detections, nil
}

// ScanResults returns a channel of scan results for the scan
func (conn *Connection) ScanResults(ctx context.Context, payload []byte) (<-chan domain.Detection, error) {
	var detections chan domain.Detection
	var err error

	if len(payload) > 0 {
		ctx = ctxtest(ctx)

		// Create the detections channel
		detections = make(chan domain.Detection)

		go func(ctx context.Context, detections chan<- domain.Detection) {
			defer handleRoutinePanic(conn.logger)
			defer close(detections)
			var err error

			// Unmarshall the payload here to the scan object
			scan := &Scan{api: conn.api}
			if err = json.Unmarshal(payload, scan); err == nil {

				// Attempt to delete the template used in the scan from nexpose to ensure that we're cleaning up after ourselves
				var temp *nexpose.ScanTemplate
				if temp, err = conn.api.GetScanTemplate(scan.TemplateID); err == nil && temp != nil {
					if err = conn.api.DeleteScanTemplate(scan.TemplateID); err == nil {
						// determine the template was deleted
						if temp, err = conn.api.GetScanTemplate(scan.TemplateID); err == nil {
							if temp != nil {
								// template still exists, deletion failed
								conn.logger.Send(log.Errorf(err, "deletion failed: template [%s] still exists and was not properly deleted", scan.TemplateID))
							}
						} else {
							if strings.Contains(err.Error(), "404") {
								err = nil
							} else {
								// unable to determine deletion status of template
								conn.logger.Send(log.Errorf(err, "unable to determine deletion status of template [%s]", scan.TemplateID))
							}
						}
					} else {
						// error deleting scan template from nexpose
						conn.logger.Send(log.Errorf(err, "error while deleting template [%s]", scan.TemplateID))
					}
				} else if err != nil {
					//error that the template was not able to be loaded
					conn.logger.Send(log.Errorf(err, "error occurred while retrieving scan template [%s]", scan.TemplateID))
				}

				var status string
				if status, err = scan.Status(); err == nil {
					if status == scanStatuses["finished"] {

						// Load the asset detections into the detections channel for processing
						conn.loadAssetDetections(conn.ctx, detections, scan.Assets)
					} else {
						// error here because we shouldn't pull results from non-finished scans
						conn.logger.Send(log.Errorf(nil, "scan [%s] not in the finished status, cannot pull detections", scan.ID()))
					}
				} else {
					conn.logger.Send(log.Errorf(err, "error occurred while pulling scan status for scan [%s]", scan.ID()))
				}
			} else {
				conn.logger.Send(log.Error("error while unmarshalling scan payload in nexpose driver", err))
			}
		}(ctx, detections)

	} else {
		err = errors.New("empty payload for loading scan results")
		conn.logger.Send(log.Error("error while gathering scan results", err))
	}

	return detections, err
}

// Scan creates a scan for a set of detections in Nexpose
func (conn *Connection) Scan(ctx context.Context, detections []domain.Match) (<-chan domain.Scan, error) {
	var scans = make(chan domain.Scan)
	var err error

	if detections != nil {
		if strings.Count(conn.settings.ScanNameFormat, "%s") == 1 {
			ctx = ctxtest(ctx)

			go func(ctx context.Context, scans chan<- domain.Scan) {
				defer handleRoutinePanic(conn.logger)
				defer close(scans)

				var groupToMatch = make(map[string][]domain.Match)

				for _, detection := range detections {
					if groupToMatch[detection.GroupID()] == nil {
						groupToMatch[detection.GroupID()] = make([]domain.Match, 0)
					}

					groupToMatch[detection.GroupID()] = append(groupToMatch[detection.GroupID()], detection)
				}

				wg := &sync.WaitGroup{}
				for groupToRescan, detectionsToRescan := range groupToMatch {
					wg.Add(1)
					go func(groupToRescan string, detectionsToRescan []domain.Match) {
						defer handleRoutinePanic(conn.logger)
						defer wg.Done()

						var vulns = make([]string, 0)
						var seenVuln = make(map[string]bool)

						for _, detection := range detectionsToRescan {
							if !seenVuln[detection.Vulnerability()] {
								seenVuln[detection.Vulnerability()] = true
								vulns = append(vulns, detection.Vulnerability())
							}
						}

						if templateID, err := conn.api.DuplicateScanTemplateWVulns(ctx, conn.settings.ScanTemplate, vulns); err == nil {
							scanName := fmt.Sprintf(conn.settings.ScanNameFormat, time.Now().Format(time.RFC3339))

							var scan *Scan
							scan, err = conn.createScanForDetections(detectionsToRescan, groupToRescan, scanName, templateID, strconv.Itoa(conn.settings.RescanSite))

							if err == nil {
								select {
								case <-ctx.Done():
									return
								case scans <- scan:
								}
							} else {
								conn.logger.Send(log.Errorf(err, "failed to create scan"))
							}
						} else {
							conn.logger.Send(log.Errorf(err, "error while creating scan template"))
						}
					}(groupToRescan, detectionsToRescan)
				}
				wg.Wait()
			}(ctx, scans)
		} else {
			err = errors.New("scan_name_format must contain a single '%s' symbol")
		}
	} else {
		err = fmt.Errorf("nil detection slice")
	}

	return scans, err
}

// Scans returns a list of scans from Nexpose
func (conn *Connection) Scans(ctx context.Context, scanPayloads <-chan []byte) <-chan domain.Scan {
	var scans = make(chan domain.Scan)

	if scanPayloads != nil {
		ctx = ctxtest(ctx)

		go func(ctx context.Context, scans chan<- domain.Scan) {
			defer handleRoutinePanic(conn.logger)
			defer close(scans)
			var err error

			for {
				select {
				case <-ctx.Done():
					return
				case payload, ok := <-scanPayloads:
					if ok {

						if len(payload) > 0 {

							// Unmarshall the payload here to the scan object
							scan := &Scan{api: conn.api}
							if err = json.Unmarshal(payload, scan); err == nil {

								// Push the unmarshalled scan onto the channel back out
								select {
								case <-ctx.Done():
									return
								case scans <- scan:
								}
							} else {
								conn.logger.Send(log.Error("error while unmarshalling scan payload in nexpose driver", err))
							}
						} else {
							err := errors.New("empty scan payload received by nexpose driver")
							conn.logger.Send(log.Error(err.Error(), err))
						}
					} else {
						return
					}
				}
			}

		}(ctx, scans)
	}

	return scans
}

// KnowledgeBase loads all vulnerabilities from the nexpose knowledgebase and returns them on a concurrency channel
func (conn *Connection) KnowledgeBase(ctx context.Context, since *time.Time) <-chan domain.Vulnerability {
	var vulnerabilities = make(chan domain.Vulnerability)
	ctx = ctxtest(ctx)

	vulns := conn.api.GetVulnerabilities(ctx, "")

	// Push off the loading of the vulnerabilities to a go routine
	go func(ctx context.Context, vulnerabilities chan<- domain.Vulnerability) {
		defer handleRoutinePanic(conn.logger)
		defer close(vulnerabilities)

		for {
			select {
			case <-ctx.Done():
				return
			case vuln, ok := <-vulns:
				if ok {

					if vuln != nil {

						wrappedVuln := &vulnerability{
							vuln:   vuln,
							api:    conn.api,
							logger: conn.logger,
						}

						// Setup a TTL cache for the vulnerability record
						conn.vulnerabilities.Store(ctx, vuln.ID, wrappedVuln, conn.ttl())

						// Wrap the vulnerability and send on the channel
						select {
						case <-ctx.Done():
							return
						case vulnerabilities <- wrappedVuln:
						}
					}
				} else {
					return
				}
			}
		}
	}(ctx, vulnerabilities)

	return vulnerabilities
}

// Validate determines whether the connection object is properly
// formed by checking the important struct elements for expected
// values
func (conn *Connection) Validate() (valid bool) {

	if conn.ctx != nil {
		if conn.api != nil {
			if conn.logger != nil {
				if conn.authentication != nil {
					if conn.settings != nil {
						valid = true
					}
				}
			}
		}
	}

	return valid
}

func (conn *Connection) ttl() (duration *time.Duration) {

	// Setup a TTL cache for the vulnerability record
	if conn.authentication.CacheTTLSeconds != nil {
		var temp = time.Duration(*conn.authentication.CacheTTLSeconds)
		duration = &temp
	}

	return duration
}

func unmarshalAuthCreateClient(ctx context.Context, authPayload string, lstream log.Logger) (authentication *domain.Host, client client, err error) {

	if len(authPayload) > 0 {

		authentication = &domain.Host{}
		if err = json.Unmarshal([]byte(authPayload), authentication); err == nil {
			if validator.IsValid(authentication) {

				// Create the HTTP client
				client, err = funnel.New(ctx,
					&http.Client{
						Transport: &http.Transport{
							TLSClientConfig: &tls.Config{
								InsecureSkipVerify: authentication.VerifyTLS(),
							},
						},
					},
					lstream,
					authentication.Delay(),
					authentication.Retries(),
					authentication.Concurrency())
			} else {
				err = errors.New("unable to validate authentication for nexpose api")
			}
		}
	} else {
		err = errors.New("authentication payload is empty")
	}

	return authentication, client, err
}

func unmarshalPayload(payload string) (settings *Payload, err error) {

	if len(payload) > 0 {

		settings = &Payload{}
		if err = json.Unmarshal([]byte(payload), settings); err == nil {
			if !validator.IsValid(settings) {
				err = errors.New("unable to validate settings for nexpose api")
			}
		}
	} else {
		err = errors.New("settings payload is empty")
	}

	return settings, err
}

func createSession(ctx context.Context, client client, host *domain.Host, lstream log.Logger) (session *nexpose.Session, err error) {
	if client != nil {
		if host != nil {
			if lstream != nil {
				if session, err = nexpose.Connect(ctx, client, host, lstream); err == nil {
					if !validator.IsValid(session) {
						err = errors.New("unable to validate the api connection to Nexpose")
					}
				}
			} else {
				err = errors.New("logger is nil for nexpose api connection")
			}
		} else {
			err = errors.New("authentication is nil for nexpose api connection")
		}
	} else {
		err = errors.New("client is nil for nexpose api connection")
	}

	return session, err
}

func (conn *Connection) loadSiteDetections(ctx context.Context, detections chan<- domain.Detection, id string) {
	var assets <-chan *nexpose.Asset
	var err error
	if assets, err = conn.api.GetAssetsForSite(ctx, id, ""); err == nil {

		wg := &sync.WaitGroup{}

		func() {
			for {
				select {
				case <-ctx.Done():
					return
				case a, ok := <-assets:
					if ok {
						wg.Add(1)
						go func(a asset) {
							defer handleRoutinePanic(a.conn.logger)
							defer wg.Done()

							var vulns <-chan domain.Detection
							var err error
							if vulns, err = a.Vulnerabilities(ctx); err == nil {
								for {
									select {
									case <-ctx.Done():
										return
									case v, ok := <-vulns:
										if ok {
											detections <- v
										} else {
											return
										}
									}
								}
							}
						}(asset{conn: conn, asset: a})
					} else {
						return
					}
				}
			}
		}()

		wg.Wait()
	} else {
		conn.logger.Send(log.Errorf(err, "error while gathering assets for site [%v]", id))
	}
}

func (conn *Connection) loadAssetDetections(ctx context.Context, detections chan<- domain.Detection, assetIDs []int) {
	var wg = sync.WaitGroup{}

	for _, assetID := range assetIDs {
		var a *nexpose.Asset
		var err error

		if a, err = conn.api.GetAsset(assetID); err == nil {
			wg.Add(1)

			go func(a asset) {
				defer handleRoutinePanic(a.conn.logger)
				defer wg.Done()
				var vulns <-chan domain.Detection
				var err error
				if vulns, err = a.Vulnerabilities(ctx); err == nil {
					for {
						select {
						case <-ctx.Done():
							return
						case v, ok := <-vulns:
							if ok {
								detections <- v
							} else {
								return
							}
						}
					}
				}
			}(asset{conn: conn, asset: a})
		} else {
			conn.logger.Send(log.Errorf(err, "error while loading asset [%v] from nexpose", assetID))
		}
	}

	wg.Wait()
}
