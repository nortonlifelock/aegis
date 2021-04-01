package connector

import (
	"context"
	"fmt"
	"github.com/nortonlifelock/aegis/pkg/domain"
	"github.com/nortonlifelock/aegis/pkg/log"
	"github.com/nortonlifelock/aegis/pkg/qualys"
	"strconv"
	"strings"
	"sync"
)

func (session *QsSession) pushCombosForHost(ctx context.Context, h qualys.QHost, devVulnMutex *sync.Mutex, processedDevVulns map[string]bool, out chan<- domain.Detection) {
	for index := range h.Detections {
		v := h.Detections[index]

		// Only load detections that are NOT fixed, and are CONFIRMED as vulnerabilities that affect the specific host
		// Confirmed is important because in Qualys there are potential vulnerabilities in the KB and those vulnerabilities
		// can be ACTUAL vulnerabilities when scanned for but the KB always shows them as "Potential" so the actual
		// status on the host itself is what determines if it is actually a vulnerability or not
		if v.Type == "Confirmed" {

			// Read the port information from the detection if a port is specified
			var port = -1
			var protocol = ""
			if v.Port != nil {
				port = *(v.Port)

				if v.Protocol != nil {
					protocol = *(v.Protocol)
				}
			}

			// Lock threads accessing this code and update the processed vulns map so that the records aren't duplicated
			// along the channel for processing
			devVulnMutex.Lock()
			if processedDevVulns[fmt.Sprintf("%v-%v-%v%s", h.HostID, v.QualysID, port, protocol)] == false {
				processedDevVulns[fmt.Sprintf("%v-%v-%v%s", h.HostID, v.QualysID, port, protocol)] = true
				devVulnMutex.Unlock()

				select {
				case <-ctx.Done():
					return
				case out <- &hostDetectionCombo{
					host: &host{
						h: h,
					},
					detection: &detection{
						d:       v,
						session: session,
					},
				}:
				}
			} else {
				devVulnMutex.Unlock()
			}
		}
	}
}

func (session *QsSession) pushDetectionsForTags(ctx context.Context, out chan<- domain.Detection, tags []string) (err error) {
	session.lstream.Send(log.Infof("Loading Detections from Qualys using tags [%s]", strings.Join(tags, ",")))

	var hosts <-chan qualys.QHost
	if hosts, err = session.apiSession.GetTagDetections(tags, session.payload.KernelFilter); err == nil {

		var processedDevVulns = make(map[string]bool)
		var devVulnMutex = &sync.Mutex{}

		wg := &sync.WaitGroup{}
		func() {
			for {
				select {
				case <-ctx.Done():
					return
				case h, ok := <-hosts:
					if ok {
						wg.Add(1)
						go func(h qualys.QHost) {
							defer handleRoutinePanic(session.lstream)
							defer wg.Done()
							session.pushCombosForHost(ctx, h, devVulnMutex, processedDevVulns, out)
						}(h)
					} else {
						return
					}
				}
			}
		}()
		wg.Wait()

	} else {
		session.lstream.Send(log.Error("Error while loading host detections from Qualys", err))
	}
	return err
}

func (session *QsSession) pushDetectionsForWebApplications(ctx context.Context, out chan<- domain.Detection, webAppIDs []string) (err error, contextClosed bool) {
	for _, webAppID := range webAppIDs {
		session.lstream.Send(log.Infof("Loading web application detection from Qualys using ID [%s]", webAppID))

		var findings []*qualys.WebAppFinding
		if findings, err = session.apiSession.GetVulnerabilitiesForSite(webAppID); err == nil {

			filteredFindings := session.getParentFindingsAndAttachChildren(findings)
			for _, detectionWrapper := range filteredFindings {
				if detectionWrapper.Status() != domain.Informational {
					select {
					case <-ctx.Done():
						return nil, true
					case out <- detectionWrapper:
					}
				}
			}
		} else {
			session.lstream.Send(log.Errorf(err, "error while gathering vulnerabilities for site [%s]", webAppID))
		}
	}

	return err, false
}

func (session *QsSession) pushDetectionsForAssetGroup(ctx context.Context, out chan<- domain.Detection, groupIDs []string) (err error) {
	session.lstream.Send(log.Infof("Loading Detections from Qualys using group IDs [%s]", strings.Join(groupIDs, ",")))

	var hosts <-chan qualys.QHost
	if hosts, err = session.apiSession.GetHostDetections(groupIDs, session.payload.KernelFilter); err == nil {

		var processedDevVulns = make(map[string]bool)
		var devVulnMutex = &sync.Mutex{}

		wg := &sync.WaitGroup{}
		func() {
			for {
				select {
				case <-ctx.Done():
					return
				case h, ok := <-hosts:
					if ok {
						wg.Add(1)
						go func(h qualys.QHost) {
							defer handleRoutinePanic(session.lstream)
							defer wg.Done()
							session.pushCombosForHost(ctx, h, devVulnMutex, processedDevVulns, out)
						}(h)
					} else {
						return
					}
				}
			}
		}()
		wg.Wait()

	} else {
		session.lstream.Send(log.Error("Error while loading host detections from Qualys", err))
	}
	return err
}

func (session *QsSession) pushDetectionsByAssetGroup(ctx context.Context, out chan<- domain.Detection, scanInfo *scan) {
	// scheduled scans should be provided the asset group ID they are covering
	detections, err := session.Detections(ctx, strings.Split(scanInfo.AssetGroupID, ","))
	if err == nil {
		func() {
			for {
				select {
				case <-ctx.Done():
					return
				case val, ok := <-detections:
					if ok {
						select {
						case <-ctx.Done():
							return
						case out <- val:
						}
					} else {
						return
					}
				}
			}
		}()
	} else {
		session.lstream.Send(log.Errorf(err, "error while loading detections for [%s]", scanInfo.AssetGroupID))
	}
}

func (session *QsSession) pushDetectionsOnChannel(ctx context.Context, output *qualys.QHostListDetectionOutput, deadHostIPToProof map[string]string, out chan<- domain.Detection) bool {
	for _, h := range output.Hosts {
		for _, d := range h.Detections {

			var unconfirmedDetection bool

			var deadHostProof = deadHostIPToProof[h.IPAddress]
			if len(deadHostProof) > 0 {
				d.Status = domain.DeadHost
				d.Proof = deadHostProof
			} else if d.Type != "Confirmed" {
				unconfirmedDetection = true
			} else if d.Status == "Fixed" {
				d.Status = domain.Fixed
			}

			if !unconfirmedDetection {
				select {
				case <-ctx.Done():
					return true
				case out <- &hostDetectionCombo{
					host: &host{
						h: h,
					},
					detection: &detection{
						d:       d,
						session: session,
					},
				}:
				}
			}
		}
	}

	return false
}

func (session *QsSession) pushDetectionsByScanTarget(ctx context.Context, scanInfo *scan, out chan<- domain.Detection, deadIPToProof chan<- domain.KeyValue) (needToClose bool) {
	var err error

	// Ask Qualys for the scan so we can find what IPs it scanned
	var scan qualys.ScanQualys
	scan, err = session.apiSession.GetScanByReference(scanInfo.ScanID)
	if err == nil {
		ipList := cleanIPList(scan.Target)

		// Use the IPs to grab the host detections
		var output *qualys.QHostListDetectionOutput
		output, err = session.apiSession.GetHostSpecificDetections(ipList, []string{scanInfo.AssetGroupID}, session.payload.KernelFilter)
		if err == nil {

			var deadHostIPToProof map[string]string
			if deadHostIPToProof, err = session.getDeadHostsForScan(scanInfo.ScanID, scanInfo.Created); err == nil {

				needToClose = false
				go func() {
					defer close(deadIPToProof)

					for deadIP, proof := range deadHostIPToProof {
						select {
						case <-ctx.Done():
							return
						case deadIPToProof <- deadIPProofCombo{
							ip:    deadIP,
							proof: proof,
						}:
						}
					}
				}()

				if session.pushDetectionsOnChannel(ctx, output, deadHostIPToProof, out) {
					return
				}
			} else {
				session.lstream.Send(log.Errorf(err, "error while loading dead hosts for scan %v", scanInfo.ScanID))
			}

			// TODO refactor this to own method
			if scanInfo.TemplateID != strconv.Itoa(session.payload.DiscoveryOptionProfileID) && scanInfo.TemplateID != strconv.Itoa(session.payload.OptionProfileID) {
				if len(scanInfo.TemplateID) > 0 {
					templateFields := strings.Split(scanInfo.TemplateID, templateDelimiter)
					if err = session.apiSession.DeleteOptionProfile(templateFields[0]); err != nil {
						session.lstream.Send(log.Errorf(err, "error while deleting option profile for scan %v", scanInfo.ScanID))
					}

					// search list only exists in vulnerability scans
					if len(templateFields) > 1 {
						if err = session.apiSession.DeleteSearchList(templateFields[1]); err != nil {
							session.lstream.Send(log.Errorf(err, "error while deleting search list for scan %v", scanInfo.ScanID))
						}
					}
				} else {
					session.lstream.Send(log.Warningf(err, "no template found in payload of scan %v", scanInfo.ScanID))
				}
			} else {
				// do nothing - we don't want to delete the option profile which we make copies of
				// this block should never hit, but we keep it just in case
			}
		} else {
			session.lstream.Send(log.Errorf(err, "error while getting host detections for scan %v", scanInfo.ScanID))
		}
	} else {
		session.lstream.Send(log.Errorf(err, "error while gathering the scan %v", scanInfo.ScanID))
	}

	return
}

type deadIPProofCombo struct {
	ip    string
	proof string
}

func (d deadIPProofCombo) Key() string {
	return d.ip
}

func (d deadIPProofCombo) Value() string {
	return d.proof
}