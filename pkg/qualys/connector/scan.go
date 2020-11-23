package connector

import (
	"context"
	"fmt"
	"github.com/nortonlifelock/aegis/pkg/domain"
	"github.com/nortonlifelock/aegis/pkg/log"
	"github.com/nortonlifelock/aegis/pkg/qualys"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type scanBundle struct {
	groupID    string
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

func getMatchesCoveredInScanBundle(bundle *scanBundle, matches []domain.Match) (matchesCoveredByBundle []domain.Match) {
	matchesCoveredByBundle = make([]domain.Match, 0)
	if bundle != nil {
		for _, match := range matches {
			if bundle.seenDevice[match.IP()] {
				matchesCoveredByBundle = append(matchesCoveredByBundle, match)
			} else if bundle.seenDevice[match.InstanceID()] {
				matchesCoveredByBundle = append(matchesCoveredByBundle, match)
			}
		}
	}

	return matchesCoveredByBundle
}

func (session *QsSession) createVulnerabilityScanForGroup(ctx context.Context, out chan<- domain.Scan, bundle *scanBundle, matches []domain.Match) (err error) {
	var scanRef string
	var optionProfileID, searchListID string

	if len(bundle.devices) > 0 && len(bundle.vulns) > 0 {
		if optionProfileID, searchListID, err = session.createOptionProfileWithSearchList(bundle.vulns, session.payload.OptionProfileID); err == nil {

			var scanCreationFunctions = make([]func() (string, string, error), 0)

			if session.payload.EC2ScanSettings[bundle.groupID] == nil {
				scanCreationFunctions = append(scanCreationFunctions, func() (string, string, error) {
					var scanTitle = fmt.Sprintf(session.payload.ScanNameFormatString, time.Now().Format(time.RFC3339))
					_, scanRef, err = session.apiSession.CreateScan(scanTitle, optionProfileID, intArrayToStringArray(bundle.appliances), bundle.networkID, bundle.devices, bundle.external)
					return scanTitle, scanRef, err
				})
			} else {
				var instanceIDs []string
				var region string
				instanceIDs, region, err = session.getEC2ScanData(matches)
				if err == nil {
					const batchSize = 10 // Qualys allows only rescan of 10 instances at a time
					for i := 0; i < len(instanceIDs); i += batchSize {
						i := i // scope the iterating variable so the outer loop doesn't overwrite it when the function is called at a later time
						settings := session.payload.EC2ScanSettings[bundle.groupID]
						var instancesCoveredInThisScan []string
						if i+batchSize <= len(instanceIDs) {
							instancesCoveredInThisScan = instanceIDs[i : i+batchSize]
						} else {
							instancesCoveredInThisScan = instanceIDs[i:]
						}

						// TODO have this method return the matches that the scan covers?
						scanCreationFunctions = append(scanCreationFunctions, func() (string, string, error) {
							bundle.seenDevice = make(map[string]bool) // have each scan overwrite the devices that it covers so we can separate out the matches that each scan covers
							for _, instanceID := range instancesCoveredInThisScan {
								bundle.seenDevice[instanceID] = true
							}
							var scanTitle = fmt.Sprintf(session.payload.ScanNameFormatString, time.Now().Format(time.RFC3339))
							_, scanRef, err = session.apiSession.CreateEC2Scan(scanTitle, optionProfileID, instancesCoveredInThisScan, region, settings.ConnectorName, settings.ScannerName)
							return scanTitle, scanRef, err
						})
					}
				} else {
					session.lstream.Send(log.Errorf(err, "error while pulling ec2 data while creating scan for group [%s]", bundle.groupID))
				}
			}

			if err == nil {
				for _, createScanFunction := range scanCreationFunctions {
					scanTitle, scanRef, err := createScanFunction()
					if err != nil {
						var retries = 0
						for errShowsThatScanLimitHit(err) {
							retries++
							select {
							case <-ctx.Done():
								return fmt.Errorf("context closed")
							default:
							}

							session.lstream.Send(log.Warningf(err, "scan limit hit while trying to create the scan, waiting 15 minutes before trying again - times tried [%d]", retries))
							time.Sleep(time.Minute * 15)
							scanTitle, scanRef, err = createScanFunction()
						}
					}

					if err == nil {
						scan := &scan{
							Name:       scanTitle,
							ScanID:     scanRef,
							TemplateID: fmt.Sprintf("%s%s%s", optionProfileID, templateDelimiter, searchListID),

							AssetGroupID: bundle.groupID,
							EngineIDs:    intArrayToStringArray(bundle.appliances),

							Created: time.Now(),
							matches: getMatchesCoveredInScanBundle(bundle, matches),
						}

						session.lstream.Send(log.Infof("scan %v created for group %v", scan.ScanID, bundle.groupID))

						select {
						case <-ctx.Done():
							return fmt.Errorf("context closed")
						case out <- scan:
						}
					} else {
						session.lstream.Send(log.Errorf(err, "error while creating scan for group [%s]", bundle.groupID))
					}
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

func (session *QsSession) getEC2ScanData(matches []domain.Match) (instanceIDs []string, region string, err error) {
	instanceIDs = make([]string, 0)
	var seen = make(map[string]bool) // TODO max out at 10 at a time
	for _, match := range matches {
		if len(match.InstanceID()) > 0 {
			if !seen[match.InstanceID()] {
				seen[match.InstanceID()] = true
				instanceIDs = append(instanceIDs, match.InstanceID())
			}
		} else {
			err = fmt.Errorf("empty instance ID found for device %s", match.Device())
			break
		}

		if len(match.Region()) > 0 {
			if len(region) == 0 {
				region = match.Region()
			} else if region != match.Region() {
				err = fmt.Errorf("found multiple regions within same ec2 group [%s|%s]", region, match.Region())
				break
			}
		}
	}

	if err == nil && len(region) == 0 {
		err = fmt.Errorf("could not determine region for any of the instance IDs [%s]", strings.Join(instanceIDs, ","))
	}

	return instanceIDs, region, err
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
					ScanID:     fmt.Sprintf("%s%s_%s", webPrefix, findingUID, time.Now().Format(time.RFC3339)),
					TemplateID: findingUID,

					AssetGroupID: detection.GroupID(),
					EngineIDs:    []string{},

					Created: time.Now(),
					matches: []domain.Match{detection},
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

					if len(bundle.devices) > 0 && len(bundle.vulns) > 0 {
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

func (session *QsSession) createDiscoveryScanForGroup(ctx context.Context, out chan<- domain.Scan, bundle *scanBundle, matches []domain.Match) (err error) {
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

						AssetGroupID: bundle.groupID,
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
		var matchIsEc2DeviceThatHasSettingsInPayload = len(match.InstanceID()) > 0 && session.payload.EC2ScanSettings[match.GroupID()] != nil

		if !matchIsEc2DeviceThatHasSettingsInPayload {
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
				if len(match.GroupID()) > 0 {
					err = fmt.Errorf("could not find a group for %v - check to see if there is an online appliance for its expected group - also check if the group ID is in the Qualys SourceConfig payload", match.IP())
				} else {
					err = fmt.Errorf("empty group ID for IP [%s]", match.IP())
				}

				break
			}
		} else {
			// here we've found a match for an ec2 scan, and need to populate the information using the instanceID instead of the IP
			if groupIDToScanBundle[match.GroupID()] == nil {
				groupIDToScanBundle[match.GroupID()] = &scanBundle{
					groupID:    match.GroupID(),
					networkID:  0,
					appliances: nil,
					external:   false,
					devices:    make([]string, 0),
					vulns:      make([]string, 0),
					seenDevice: make(map[string]bool),
					seenVuln:   make(map[string]bool),
				}
			}

			if !groupIDToScanBundle[match.GroupID()].seenDevice[match.InstanceID()] {
				groupIDToScanBundle[match.GroupID()].seenDevice[match.InstanceID()] = true
				groupIDToScanBundle[match.GroupID()].devices = append(groupIDToScanBundle[match.GroupID()].devices, match.InstanceID())
			}

			if !groupIDToScanBundle[match.GroupID()].seenVuln[match.Vulnerability()] {
				groupIDToScanBundle[match.GroupID()].seenVuln[match.Vulnerability()] = true
				groupIDToScanBundle[match.GroupID()].vulns = append(groupIDToScanBundle[match.GroupID()].vulns, match.Vulnerability())
			}
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
							groupID:    strconv.Itoa(group.ID),
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
			var ipToAGs = session.mapIPToAssetGroup(matches)
			for ip, ags := range ipToAGs {
				var found bool
				for _, ag := range ags {
					if groupIDToScanBundle[ag] != nil {

						if !groupIDToScanBundle[ag].seenDevice[ip] {
							groupIDToScanBundle[ag].seenDevice[ip] = true
							groupIDToScanBundle[ag].devices = append(groupIDToScanBundle[ag].devices, ip)
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

func errShowsThatScanLimitHit(err error) (scanLimitHit bool) {
	if err != nil {
		errContents := err.Error()
		if len(errContents) > 0 {
			scanLimitHit = regexp.MustCompile("You are allowed to run \\d+ concurrent scans").MatchString(errContents) &&
				regexp.MustCompile("This limit has already been reached").MatchString(errContents)
		}
	}

	return scanLimitHit
}
