package connector

import (
	"context"
	"fmt"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/log"
	"github.com/nortonlifelock/qualys"
	"strconv"
	"strings"
	"sync"
	"time"
)

type scanBundle struct {
	groupID    int
	networkID  int
	appliances []int
	external   bool
	ips        []string
	vulns      []string
	seenIP     map[string]bool
	seenVuln   map[string]bool
}

func intArrayToStringArray(intIn []int) (stringOut []string) {
	stringOut = make([]string, 0)

	if intIn != nil && len(intIn) > 0 {

		for _, value := range intIn {
			stringOut = append(stringOut, strconv.Itoa(value))
		}
	}

	return stringOut
}

func cleanIPList(ipList string) (ips []string) {
	ips = make([]string, 0)

	ipList = strings.Replace(ipList, " ", "", -1)
	split := strings.Split(ipList, ",")

	for _, ip := range split {
		if len(ip) > 0 {
			ips = append(ips, ip)
		}
	}

	return ips
}

func getMatchesCoveredInScanBundle(bundle *scanBundle, matches []domain.Match) (matchesCoveredByBundle []domain.Match) {
	matchesCoveredByBundle = make([]domain.Match, 0)
	if bundle != nil {
		for _, match := range matches {
			if bundle.seenIP[match.IP()] {
				matchesCoveredByBundle = append(matchesCoveredByBundle, match)
			}
		}
	}

	return matchesCoveredByBundle
}

func (session *QsSession) createVulnerabilityScanForGroup(ctx context.Context, out chan<- domain.Scan, bundle *scanBundle, matches []domain.Match) (err error) {
	var scanRef string
	var optionProfileID, searchListID string

	if len(bundle.ips) > 0 && len(bundle.vulns) > 0 {
		if optionProfileID, searchListID, err = session.createOptionProfileWithSearchList(bundle.vulns, session.payload.OptionProfileID); err == nil {
			var scanTitle = fmt.Sprintf(session.payload.ScanNameFormatString, time.Now().Format(time.RFC3339))
			if _, scanRef, err = session.apiSession.CreateScan(scanTitle, optionProfileID, intArrayToStringArray(bundle.appliances), bundle.networkID, bundle.ips, bundle.external); err == nil {

				scan := &scan{
					Name:       scanTitle,
					ScanID:     scanRef,
					TemplateID: fmt.Sprintf("%s%s%s", optionProfileID, templateDelimiter, searchListID),

					AssetGroupID: strconv.Itoa(bundle.groupID),
					EngineIDs:    intArrayToStringArray(bundle.appliances),

					Created: time.Now(),
					matches: getMatchesCoveredInScanBundle(bundle, matches),
				}

				session.lstream.Send(log.Infof("scan %v created for group %v", scan.ScanID, bundle.groupID))

				select {
				case <-ctx.Done():
					return
				case out <- scan:
				}
			}
		} else {
			session.lstream.Send(log.Errorf(err, "error while creating option profile and search list"))
		}
	} else {
		// do nothing
	}

	return err
}

func (session *QsSession) createScanForWebApplication(ctx context.Context, detections []domain.Match, out chan<- domain.Scan) {
	var seen = make(map[string]bool)

	for _, detection := range detections {
		var findingUID = detection.Device()

		if !seen[findingUID] {
			seen[findingUID] = true

			_, err := session.apiSession.CreateRetestForWebAppVulnerabilityFinding(findingUID)
			if err == nil {
				scan := &scan{
					Name:       fmt.Sprintf("was_aegis_retest_%s_%s", findingUID, time.Now().Format(time.RFC3339)),
					ScanID:     fmt.Sprintf("%s_%s_%s", webPrefix, findingUID, time.Now().Format(time.RFC3339)),
					TemplateID: findingUID,

					AssetGroupID: detection.GroupID(),
					EngineIDs:    []string{},

					Created: time.Now(),
					matches: detections,
				}

				select {
				case <-ctx.Done():
					return
				case out <- scan:
				}
			} else {
				session.lstream.Send(log.Errorf(err, "error while retesting finding UID [%s]", findingUID))
			}
		}
	}
}

func (session *QsSession) createScanForDetections(ctx context.Context, detections []domain.Match, out chan<- domain.Scan) {
	var err error

	var groupIDToScanBundle map[string]*scanBundle
	if groupIDToScanBundle, err = session.prepareIPsAndAGMapping(detections); err == nil {
		if err = session.populateGroupVulnerabilityChecks(detections, groupIDToScanBundle); err == nil {
			// wg to ensure we don't close the out channel before the writing threads finish
			wg := &sync.WaitGroup{}

			for groupID := range groupIDToScanBundle {
				wg.Add(1)
				go func(bundle *scanBundle) {
					defer handleRoutinePanic(session.lstream)
					defer wg.Done()

					if len(bundle.ips) > 0 && len(bundle.vulns) > 0 {
						session.lstream.Send(log.Infof("Creating vulnerability scan for group %v", bundle.groupID))
						// error intentionally scoped out
						err := session.createVulnerabilityScanForGroup(ctx, out, bundle, detections)
						if err != nil {
							session.lstream.Send(log.Errorf(err, "error while creating scan for group %v", bundle.groupID))
						}
					}
				}(groupIDToScanBundle[groupID])
			}

			wg.Wait()
		} else {
			session.lstream.Send(log.Errorf(err, "error while populating group vulnerability checks"))
		}
	} else {
		session.lstream.Send(log.Errorf(err, "error while creating assignment group mapping for vulnerability scan"))
	}
}

/*

 */

func (session *QsSession) createDiscoveryScanForGroup(ctx context.Context, out chan<- domain.Scan, bundle *scanBundle, matches []domain.Match) (err error) {
	if len(bundle.ips) > 0 {

		if session.payload.DiscoveryOptionProfileID > 0 {
			var scanRef string
			var optionProfileID string

			if optionProfileID, err = session.createCopyOfOptionProfile(session.payload.DiscoveryOptionProfileID); err == nil {
				var scanTitle = fmt.Sprintf(session.payload.ScanNameFormatString, time.Now().Format(time.RFC3339))

				if _, scanRef, err = session.apiSession.CreateScan(scanTitle, optionProfileID, intArrayToStringArray(bundle.appliances), bundle.networkID, bundle.ips, bundle.external); err == nil {

					scan := &scan{
						Name:       scanTitle,
						ScanID:     scanRef,
						TemplateID: optionProfileID,

						AssetGroupID: strconv.Itoa(bundle.groupID),
						EngineIDs:    intArrayToStringArray(bundle.appliances),

						Created: time.Now(),
						matches: getMatchesCoveredInScanBundle(bundle, matches),
					}

					select {
					case <-ctx.Done():
						return
					case out <- scan:
						session.lstream.Send(log.Infof("created discovery scan for group %v", bundle.groupID))
					}
				}
			} else {
				err = fmt.Errorf("error while creating option profile for discovery scan - %v", err.Error())
			}
		} else {
			err = fmt.Errorf("empty discovery option profile ID in Qualys payload")
		}
	} else {
		// do nothing
	}

	return err
}

func (session *QsSession) populateGroupVulnerabilityChecks(detections []domain.Match, groupIDToScanBundle map[string]*scanBundle) (err error) {
	for _, match := range detections {
		var matchFound bool

		// check every group for each match to see which group it belongs to
		for _, group := range groupIDToScanBundle {
			if group.seenIP[match.IP()] {
				matchFound = true
				group.vulns = append(group.vulns, match.Vulnerability())
				break
			}
		}

		if !matchFound {
			if len(match.GroupID()) > 0 {
				err = fmt.Errorf("empty group ID for IP [%s]", match.IP())
			} else {
				err = fmt.Errorf("could not find a group for %v - check to see if there is an online appliance for its expected group", match.IP())
			}

			break
		}
	}

	return err
}

func (session *QsSession) prepareIPsAndAGMapping(matches []domain.Match) (groupIDToScanBundle map[string]*scanBundle, err error) {
	var groups []*qualys.QSAssetGroup
	if groups, err = session.getAssetGroups(append(session.payload.AssetGroups, session.payload.ExternalGroups...)); err == nil {
		groupIDToScanBundle = make(map[string]*scanBundle)

		// even though slices in golang are pass-by-value, groups is a slice of pointers, so modifying the elements on the slice
		// within a method call will effect the elements of the slice of the caller
		if err = session.populateOnlineAppliances(groups); err == nil {

			// initialize the group map for each group that has at least one online appliance (scanning engine)
			for _, group := range groups {
				if groupIDToScanBundle[strconv.Itoa(group.ID)] == nil {

					externalGroup := session.isExternalGroup(group.ID)

					if len(group.OnlineAppliances) > 0 || externalGroup {
						groupIDToScanBundle[strconv.Itoa(group.ID)] = &scanBundle{
							groupID:    group.ID,
							networkID:  group.NetworkID,
							appliances: group.OnlineAppliances,
							external:   externalGroup,
							ips:        make([]string, 0),
							vulns:      make([]string, 0),
							seenIP:     make(map[string]bool),
							seenVuln:   make(map[string]bool),
						}
					}
				}
			}

			// map a device IP to an assignment groups that contain it
			var ipToAGs = session.mapIPToAssetGroup(matches)
			for ip, ags := range ipToAGs {
				var found bool
				for _, ag := range ags {
					if groupIDToScanBundle[ag] != nil {

						if !groupIDToScanBundle[ag].seenIP[ip] {
							groupIDToScanBundle[ag].seenIP[ip] = true
							groupIDToScanBundle[ag].ips = append(groupIDToScanBundle[ag].ips, ip)
						}
						found = true
						break
					} else {
						// potentialGroupID does not have any online appliances
						// do nothing
					}
				}

				if !found {
					session.lstream.Send(log.Errorf(err, "could not find asset group with online engine for IP [%s]", ip))
				}
			}

		} else {
			err = fmt.Errorf("error while gathering online appliances - %s", err.Error())
		}
	} else {
		err = fmt.Errorf("error while processing asset groups - %s", err.Error())
	}

	return groupIDToScanBundle, err
}

func (session *QsSession) isExternalGroup(groupID int) (val bool) {
	for _, externalGroup := range session.payload.ExternalGroups {
		if groupID == externalGroup {
			val = true
			break
		}
	}

	return val
}
