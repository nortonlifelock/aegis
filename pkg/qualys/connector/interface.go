package connector

import (
	"context"
	"encoding/json"
	"github.com/nortonlifelock/aegis/pkg/domain"
	"github.com/nortonlifelock/aegis/pkg/log"
	"github.com/nortonlifelock/aegis/pkg/qualys"
	"strings"
	"sync"
	"time"
)

const tagPrefix = "tag-"
const webPrefix = "web-"

// KnowledgeBase grabs all vulnerabilities from the Qualys knowledge base and pushes them onto a channel
func (session *QsSession) KnowledgeBase(ctx context.Context, since *time.Time) <-chan domain.Vulnerability {
	var out = make(chan domain.Vulnerability, 50)

	go func(out chan<- domain.Vulnerability) {
		defer handleRoutinePanic(session.lstream)
		defer close(out)
		var err error

		start := time.Now()
		if err = session.loadAndCacheQualysKB(since); err == nil {
			session.lstream.Send(log.Infof("%d vulnerabilities loaded, took %s - beginning processing", len(session.vulnerabilities), time.Since(start).Round(time.Second)))

			var wg = &sync.WaitGroup{}
			var count = 0
			for index := range session.vulnerabilities {

				select {
				case <-ctx.Done():
					return
				default:
					// Create 50 vulnerabilities at a time so we don't have tons of goroutines sitting around forever waiting to finish
					count++
					if count%50 == 0 {
						wg.Wait()
					}

					wg.Add(1)
					go func(v *qualys.QVulnerability) {
						defer handleRoutinePanic(session.lstream)
						defer wg.Done()

						select {
						case <-ctx.Done():
							return
						case out <- &vulnerabilityInfo{v: v}:
						}
					}(session.vulnerabilities[index])
				}
			}

			wg.Wait()
		} else {
			session.lstream.Send(log.Error("Error while loading vulnerabilities", err))
		}
	}(out)

	return out
}

// Detections returns a channel which contains combinations of devices/vulnerabilities within a group where the vulnerability was found on the device
// and the device exists in the asset group
func (session *QsSession) Detections(ctx context.Context, ids []string) (detections <-chan domain.Detection, err error) {
	var out = make(chan domain.Detection)

	go func(out chan<- domain.Detection) {
		defer handleRoutinePanic(session.lstream)
		defer close(out)

		var tags = make([]string, 0)
		var webAppIDs = make([]string, 0)
		var groupIDs = make([]string, 0)
		for _, id := range ids {

			if strings.Index(id, tagPrefix) >= 0 {
				// ids with the prefix "tag-" before their integer ID are considered a tag tracked asset in Qualys
				tags = append(tags, id[strings.Index(id, tagPrefix)+len(tagPrefix):])
			} else if strings.Index(id, webPrefix) >= 0 {
				// ids with the prefix "web-" before their integer ID are considered a Qualys Web Application Security WebApp ID
				webAppIDs = append(webAppIDs, id[strings.Index(id, webPrefix)+len(webPrefix):])
			} else {
				// plain integer IDs are considered a Qualys asset group
				groupIDs = append(groupIDs, id)
			}
		}

		if len(groupIDs) > 0 {
			err = session.pushDetectionsForAssetGroup(ctx, out, groupIDs)
		}

		if len(webAppIDs) > 0 {
			var contextClosed bool
			err, contextClosed = session.pushDetectionsForWebApplications(ctx, out, webAppIDs)
			if contextClosed {
				return
			}
		}

		if len(tags) > 0 {
			err = session.pushDetectionsForTags(ctx, out, tags)
		}

	}(out)

	return out, err
}



// ScanResults takes a scanID and returns a series of detections that were found by the corresponding scan
func (session *QsSession) ScanResults(ctx context.Context, payload []byte) (<-chan domain.Detection, <-chan domain.KeyValue, error) {
	var out = make(chan domain.Detection)
	var deadIPToProof = make(chan domain.KeyValue)

	go func(out chan<- domain.Detection, deadIPToProof chan<- domain.KeyValue) {
		defer handleRoutinePanic(session.lstream)
		defer close(out)

		var needToCloseDeadIPChannel = true
		defer func() {
			if needToCloseDeadIPChannel {
				close(deadIPToProof)
			}
		}()

		scanInfo := &scan{session: session}
		if err := json.Unmarshal(payload, scanInfo); err == nil {

			if !scanInfo.Scheduled && session.payload.EC2ScanSettings[scanInfo.AssetGroupID] == nil {
				if len(scanInfo.ScanID) > 0 {
					if strings.Contains(scanInfo.ScanID, "scan") {
						// the scan was a VM scan
						needToCloseDeadIPChannel = session.pushDetectionsByScanTarget(ctx, scanInfo, out, deadIPToProof)
					} else if strings.Contains(scanInfo.ScanID, webPrefix) {
						// the scan was a web application scan
						session.pushDetectionsByAssetGroup(ctx, out, scanInfo)
					} else {
						session.lstream.Send(log.Infof("malformed scan ID [%s]", scanInfo.ScanID))
					}
				} else {
					session.lstream.Send(log.Errorf(err, "zero length scan ID received in payload"))
				}
			} else if len(scanInfo.AssetGroupID) > 0 {
				// this block hits for scheduled scans or EC2 scans
				session.pushDetectionsByAssetGroup(ctx, out, scanInfo)
			} else {
				session.lstream.Send(log.Errorf(err, "Scheduled scan [%s] did not specify the group IDs or cloud tags that it executed against", scanInfo.Name))
			}

		} else {
			session.lstream.Send(log.Errorf(err, "error while unmarshalling scan"))
		}

	}(out, deadIPToProof)

	return out, deadIPToProof, nil
}

// Discovery kicks of a Qualys scan to identify which devices corresponding to the IPs are online
func (session *QsSession) Discovery(ctx context.Context, matches []domain.Match) (scanID <-chan domain.Scan) {
	var out = make(chan domain.Scan)

	go func(out chan<- domain.Scan) {
		defer handleRoutinePanic(session.lstream)
		defer close(out)

		var err error
		var groupIDToScanBundle map[string]*scanBundle
		if groupIDToScanBundle, err = session.prepareIPsAndAGMapping(matches); err == nil {
			// wg to ensure we don't close the out channel before the threads finish
			wg := &sync.WaitGroup{}

			for groupID := range groupIDToScanBundle {
				wg.Add(1)
				go func(bundle *scanBundle) {
					defer handleRoutinePanic(session.lstream)
					defer wg.Done()

					// error intentionally scoped out
					err := session.createDiscoveryScanForGroup(ctx, out, bundle, matches)
					if err != nil {
						session.lstream.Send(log.Errorf(err, "error while creating scan for group %v", bundle.groupID))
					}
				}(groupIDToScanBundle[groupID])
			}

			wg.Wait()
		} else {
			session.lstream.Send(log.Errorf(err, "error while creating assignment group mapping for discovery scan"))
		}
	}(out)

	return out
}

// Scan creates a Qualys scan for all the devices/vulnerabilities passed onto the channel. All vulnerabilities passed will be searched for on all devices
// passed. The scanID is passed onto a channel before closing the channel immediately.
func (session *QsSession) Scan(ctx context.Context, detections []domain.Match) (scanID <-chan domain.Scan, err error) {
	var out = make(chan domain.Scan)

	go func(out chan<- domain.Scan) {
		defer handleRoutinePanic(session.lstream)
		defer close(out)

		if len(detections) > 0 {
			if !strings.Contains(detections[0].GroupID(), webPrefix) {
				session.createScanForDetections(ctx, detections, out)
			} else {
				session.createScanForWebApplication(ctx, detections, out)
			}
		} else {
			session.lstream.Send(log.Error("empty list of detections sent for scan", nil))
		}
	}(out)

	return out, err
}

// Scans takes a channel of scan payloads, and unmarshals into a scan struct that implements the Scan interface
// these interfaces are then used to gather the statuses of the scans by the caller
func (session *QsSession) Scans(ctx context.Context, payloads <-chan []byte) (scans <-chan domain.Scan) {
	var out = make(chan domain.Scan)

	go func(out chan<- domain.Scan) {
		defer handleRoutinePanic(session.lstream)
		defer close(out)

		var seen = make(map[string]bool)
		for {
			select {
			case <-ctx.Done():
				return
			case payload, ok := <-payloads:
				if ok {
					scan := &scan{session: session}
					if err := json.Unmarshal(payload, scan); err == nil {

						if len(scan.ScanID) > 0 {
							seen[scan.Name] = true
							select {
							case <-ctx.Done():
								return
							case out <- scan:
							}
						} else if !seen[scan.Name] && len(scan.Name) > 0 {
							// this block hits when the title for an expected scheduled scan was passed instead of a scan reference
							// here we must check to see if one of those scan schedules actually have a scan running, if it does - we push it on the channel

							// the empty scan ID means that a scheduled scan isn't running with that name currently, or the recently created scan with that name hasn't had it's scan ID
							// loaded yet. We check for a running scan and populate the scan ID of the scheduled
							scheduledScan, err := session.apiSession.GetScheduledScan(scan.Name)
							if err == nil {
								if scheduledScan != nil {
									seen[scan.Name] = true
									scan.ScanID = scheduledScan.Reference
									scan.Created = scheduledScan.LaunchDate
									scan.Scheduled = true

									select {
									case <-ctx.Done():
										return
									case out <- scan:
									}
								}
							} else {
								session.lstream.Send(log.Errorf(err, "error while finding scheduled scan for [%s]", scan.Name))
							}
						}
					} else {
						session.lstream.Send(log.Errorf(err, "error while marshaling scan"))
					}
				} else {
					return
				}
			}
		}
	}(out)

	return out
}
