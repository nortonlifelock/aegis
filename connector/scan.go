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
	devices    []string
	vulns      []string
	seenDevice map[string]bool
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

func (session *QsSession) createVulnerabilityScanForGroup(ctx context.Context, out chan<- domain.Scan, bundle *scanBundle) (err error) {
	var scanRef string
	var optionProfileID, searchListID string

	if len(bundle.devices) > 0 && len(bundle.vulns) > 0 {
		if optionProfileID, searchListID, err = session.createOptionProfileWithSearchList(bundle.vulns, session.payload.OptionProfileID); err == nil {
			var scanTitle = fmt.Sprintf(session.payload.ScanNameFormatString, time.Now().Format(time.RFC3339))
			if _, scanRef, err = session.apiSession.CreateScan(scanTitle, optionProfileID, intArrayToStringArray(bundle.appliances), bundle.networkID, bundle.devices, bundle.external); err == nil {

				scan := &scan{
					Name:       scanTitle,
					ScanID:     scanRef,
					TemplateID: fmt.Sprintf("%s%s%s", optionProfileID, templateDelimiter, searchListID),

					AssetGroupID: strconv.Itoa(bundle.groupID),
					EngineIDs:    intArrayToStringArray(bundle.appliances),

					Created: time.Now(),
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

func (session *QsSession) createScanForDetections(ctx context.Context, detections []domain.Match, out chan<- domain.Scan) {
	var err error
	// make a list of unique IPs so we can assign them to a group
	var seen = make(map[string]bool)
	var ips = make([]string, 0)
	for _, detection := range detections {
		if !seen[detection.IP()] {
			seen[detection.IP()] = true
			ips = append(ips, detection.IP())
		}
	}

	var groupIDToScanBundle map[int]*scanBundle
	if groupIDToScanBundle, err = session.prepareIPsAndAGMapping(ips); err == nil {
		if err = session.populateGroupVulnerabilityChecks(detections, groupIDToScanBundle); err == nil {
			// wg to ensure we don't close the out channel before the writing threads finish
			wg := &sync.WaitGroup{}

			for groupID := range groupIDToScanBundle {
				wg.Add(1)
				go func(bundle *scanBundle) {
					defer handleRoutinePanic(session.lstream)
					defer wg.Done()

					if len(bundle.devices) > 0 && len(bundle.vulns) > 0 {
						session.lstream.Send(log.Infof("Creating vulnerability scan for group %v", bundle.groupID))
						// error intentionally scoped out
						err := session.createVulnerabilityScanForGroup(ctx, out, bundle)
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

func (session *QsSession) createDiscoveryScanForGroup(ctx context.Context, out chan<- domain.Scan, bundle *scanBundle) (err error) {
	if len(bundle.devices) > 0 {

		if session.payload.DiscoveryOptionProfileID > 0 {
			var scanRef string
			var optionProfileID string

			if optionProfileID, err = session.createCopyOfOptionProfile(session.payload.DiscoveryOptionProfileID); err == nil {
				var scanTitle = fmt.Sprintf(session.payload.ScanNameFormatString, time.Now().Format(time.RFC3339))

				if _, scanRef, err = session.apiSession.CreateScan(scanTitle, optionProfileID, intArrayToStringArray(bundle.appliances), bundle.networkID, bundle.devices, bundle.external); err == nil {

					scan := &scan{
						Name:       scanTitle,
						ScanID:     scanRef,
						TemplateID: optionProfileID,

						AssetGroupID: strconv.Itoa(bundle.groupID),
						EngineIDs:    intArrayToStringArray(bundle.appliances),

						Created: time.Now(),
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

func (session *QsSession) populateGroupVulnerabilityChecks(detections []domain.Match, groupIDToScanBundle map[int]*scanBundle) (err error) {
	for _, match := range detections {
		var matchFound bool

		// check every group for each match to see which group it belongs to
		for _, group := range groupIDToScanBundle {
			if group.seenDevice[match.IP()] {
				matchFound = true
				group.vulns = append(group.vulns, match.Vulnerability())
				break
			}
		}

		if !matchFound {
			err = fmt.Errorf("could not find a group for %v - check to see if there is an online appliance for its expected group", match.IP())
			break
		}
	}

	return err
}

func (session *QsSession) prepareIPsAndAGMapping(ips []string) (groupIDToScanBundle map[int]*scanBundle, err error) {
	var groups []*qualys.QSAssetGroup
	if groups, err = session.getAssetGroups(append(session.payload.AssetGroups, session.payload.ExternalGroups...)); err == nil {
		groupIDToScanBundle = make(map[int]*scanBundle)

		// even though slices in golang are pass-by-value, groups is a slice of pointers, so modifying the elements on the slice
		// within a method call will effect the elements of the slice of the caller
		if err = session.populateOnlineAppliances(groups); err == nil {

			// initialize the group map for each group that has at least one online appliance (scanning engine)
			for _, group := range groups {
				if groupIDToScanBundle[group.ID] == nil {

					externalGroup := session.isExternalGroup(group.ID)

					if len(group.OnlineAppliances) > 0 || externalGroup {
						groupIDToScanBundle[group.ID] = &scanBundle{
							groupID:    group.ID,
							networkID:  group.NetworkID,
							appliances: group.OnlineAppliances,
							external:   externalGroup,
							devices:    make([]string, 0),
							vulns:      make([]string, 0),
							seenDevice: make(map[string]bool),
							seenVuln:   make(map[string]bool),
						}
					}
				}
			}

			// map a device IP to an assignment groups that contain it
			var ipToAGs map[string][]int
			if ipToAGs, err = session.GetAGsForIPs(ips); err == nil {
				for _, device := range ips {

					// there's at least one assignment group for the device
					if len(ipToAGs[device]) > 0 {

						var found bool

						// find an assignment group to assign the device to for the scan
						for _, potentialGroupID := range ipToAGs[device] {

							// find a group that has at least one online appliance
							if groupIDToScanBundle[potentialGroupID] != nil {
								found = true
								if !groupIDToScanBundle[potentialGroupID].seenDevice[device] {
									groupIDToScanBundle[potentialGroupID].seenDevice[device] = true
									groupIDToScanBundle[potentialGroupID].devices = append(groupIDToScanBundle[potentialGroupID].devices, device)
								}
								break
							} else {
								// potentialGroupID does not have any online appliances
								// do nothing
							}
						}

						if !found {
							session.lstream.Send(log.Errorf(fmt.Errorf("could not find asset group with online engine for %s, check to see if it's asset group is in the Qualys source config payload", device), "no online engines found"))
						}

					} else {
						session.lstream.Send(log.Errorf(fmt.Errorf("could not find any assignment groups for %s", device), "empty AG list for device"))
					}
				}

			} else {
				err = fmt.Errorf("error while gathering Qualys asset groups - %s", err.Error())
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
