package connector

import (
	"bytes"
	"fmt"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/log"
	"github.com/nortonlifelock/qualys"
	"net"
	"sort"
	"strconv"
	"strings"
)

// GetAGsForIPs returns a map that attaches an ip to a list of assignment groups that it belongs to
func (session *QsSession) GetAGsForIPs(ips []string) (ipToAGs map[string][]int, err error) {
	if ipToAGs, err = session.getAGMapping(ips); err == nil {
		// verify that there were AGs found for every IP
		for _, ip := range ips {
			if ipToAGs[ip] != nil {
				sort.Ints(ipToAGs[ip])
			} else {
				err = fmt.Errorf("could not find the asset groups from Qualys API for [%s]", ip)
				break
			}
		}
	} else {
		err = fmt.Errorf("error while loading AG information - %v", err)
	}

	return ipToAGs, err
}

func (session *QsSession) mapIPToAssetGroup(matches []domain.Match) (ipToAGs map[string][]string) {
	var seenIPAndGroup = make(map[string]bool)

	ipToAGs = make(map[string][]string)
	for _, match := range matches {
		if len(match.GroupID()) > 0 {

			if session.payload.EC2ScanSettings[match.GroupID()] == nil { // the ec2 scan settings are in the Qualys payload and we don't need to load them from the API
				if ipToAGs[match.IP()] == nil {
					ipToAGs[match.IP()] = make([]string, 0)
				}

				var key = fmt.Sprintf("%s;%s", match.IP(), match.GroupID())
				if !seenIPAndGroup[key] {
					seenIPAndGroup[key] = true

					ipToAGs[match.IP()] = append(ipToAGs[match.IP()], match.GroupID())
				}
			}
		} else {
			session.lstream.Send(log.Errorf(nil, "Device with IP [%s] did not provide an associated GroupID", match.IP()))
		}
	}

	return ipToAGs
}

func (session *QsSession) getAGMapping(ips []string) (ipToAGs map[string][]int, err error) {
	ipToAGs = make(map[string][]int)
	chunkedIPs := breakIPsIntoSmallerGroups(ips)

	func() {
		for _, ipList := range chunkedIPs {
			var output *qualys.HostListOutput
			output, err = session.apiSession.GetHostAGInfo(ipList)
			if err == nil {
				// TODO is it worth checking the Network ID here?
				for _, host := range output.Response.HostList.Host {

					if host.TrackingMethod == "IP" {
						if len(host.AssetGroupIDs) > 0 {
							if ipToAGs[host.IP] == nil {
								ipToAGs[host.IP] = make([]int, 0)
							}

							var agsAsString = strings.Split(host.AssetGroupIDs, ",")
							for _, ag := range agsAsString {
								var agAsInt int
								if agAsInt, err = strconv.Atoi(ag); err == nil {
									ipToAGs[host.IP] = append(ipToAGs[host.IP], agAsInt)
								} else {
									return
								}
							}
						}
					}
				}
			} else {
				return
			}
		}
	}()

	return ipToAGs, err
}

// Qualys can only handle URIs of length 7000 or so, so we need to break the ips into smaller lists before gathering their AG information
func breakIPsIntoSmallerGroups(ips []string) (ipChunks [][]string) {
	const maxLenOfIPsAllowedForURI = 7000
	ipChunks = make([][]string, 0)

	if len(ips) > 0 {
		var ipList = []string{ips[0]}
		for index := 1; index < len(ips); index++ {
			if len(strings.Join(ipList, ",")) < maxLenOfIPsAllowedForURI {
				ipList = append(ipList, ips[index])
			} else {
				ipChunks = append(ipChunks, ipList)
				ipList = []string{ips[index]}
			}
		}
		ipChunks = append(ipChunks, ipList)
	}

	return ipChunks
}

func (session *QsSession) doesIPExistForThisGroup(ip string, sGroups *qualys.QSAssetGroup) bool {
	var found bool
	var groupIps = sGroups.IPs

	for ipIndex := range groupIps {
		gIP := strings.TrimSpace(groupIps[ipIndex])
		if ip == gIP {
			if sGroups.ID > 0 {
				found = true
				break
			}
		}
	}

	if !found {
		for rangeIndex := range sGroups.Ranges {
			rangeIps := strings.Split(strings.TrimSpace(sGroups.Ranges[rangeIndex]), "-")
			if len(rangeIps) == 2 {
				found = session.isIPInRange(ip, rangeIps[0], rangeIps[1])
			}
			if found {
				break
			}

		}

	}

	return found
}

func (session *QsSession) isIPInRange(ip string, ipFrom string, ipTo string) bool {
	var found bool
	var err error
	ipF := net.ParseIP(ipFrom)
	ipT := net.ParseIP(ipTo)
	ipToCompare := net.ParseIP(ip)
	if ipToCompare.To4() == nil {
		session.lstream.Send(log.Errorf(err, "%v for IP is not an IPv4 address", ipToCompare))

	}
	if ipF.To4() == nil {
		session.lstream.Send(log.Errorf(err, "%v for IP From(Range) is not an IPv4 address", ipF))

	}
	if ipT.To4() == nil {
		session.lstream.Send(log.Errorf(err, "%v for IP To(Range) is not an IPv4 address", ipT))

	}

	if bytes.Compare(ipToCompare, ipF) >= 0 && bytes.Compare(ipToCompare, ipT) <= 0 {
		session.lstream.Send(log.Infof("%v is between %v and %v", ipToCompare, ipF, ipT))
		found = true
	}

	return found
}

// getAppliancesForAssetGroups loads asset groups from the Qualys API and using that information gets all scan
// appliances from them after filtering them against the systems configured asset groups
func (session *QsSession) getAssetGroups(assetGroups []int) (groups []*qualys.QSAssetGroup, err error) {
	if session.assetGroupCache == nil {
		session.lstream.Send(log.Info("Loading asset groups with Engines from Qualys"))

		// Load the asset groups from Qualys
		var ags *qualys.QSAGListOutput
		if ags, err = session.apiSession.LoadAssetGroups(assetGroups); err == nil {

			if ags != nil {
				session.lstream.Send(log.Info("Completed Loading asset groups with Engines from Qualys"))
				groups = ags.Groups
				session.assetGroupCache = ags.Groups
			} else {
				err = fmt.Errorf("failed to load Qualys asset groups")
			}
		}
	} else {
		groups = session.assetGroupCache
	}

	return groups, err
}

func (session *QsSession) getAssetGroupWithOnlineAppliancesForIP(ip string, groups []*qualys.QSAssetGroup) (applicable []*qualys.QSAssetGroup, err error) {
	applicable = make([]*qualys.QSAssetGroup, 0)

	for index := range groups {
		if len(groups[index].OnlineAppliances) > 0 {
			if session.doesIPExistForThisGroup(ip, groups[index]) {
				applicable = append(applicable, groups[index])
			}
		}
	}

	if len(applicable) == 0 {
		err = fmt.Errorf("could not find asset group with online appliances for [%s]", ip)
	}

	return applicable, err
}

func (session *QsSession) populateOnlineAppliances(groups []*qualys.QSAssetGroup) (err error) {

	var appliances = make([]string, 0)
	var seenAppliance = make(map[string]bool)
	for _, group := range groups {
		applianceString := strings.Replace(group.Appliances, " ", "", -1)
		var groupAppliances = strings.Split(applianceString, ",")

		for _, groupAppliance := range groupAppliances {
			if !seenAppliance[groupAppliance] {
				seenAppliance[groupAppliance] = true

				if len(groupAppliance) > 0 {
					if intVal, err := strconv.Atoi(groupAppliance); err == nil && intVal > 0 {
						appliances = append(appliances, groupAppliance)
					} else if err != nil {
						session.lstream.Send(log.Errorf(err, "could not parse [%s] as an integer for gathering appliances", groupAppliance))
					}
				}
			}
		}
	}

	if len(appliances) > 0 {
		var output *qualys.QAppliances
		output, err = session.apiSession.GetApplianceInformation(appliances)

		if err == nil {
			for _, group := range groups {
				groupAppliances := strings.Split(group.Appliances, ",")
				group.OnlineAppliances = make([]int, 0)

				for _, appliance := range output.Appliances {
					if appliance.Status == "Online" {
						if elementExistsInSlice(groupAppliances, strconv.Itoa(appliance.ID)) {
							group.OnlineAppliances = append(group.OnlineAppliances, appliance.ID)
						}
					}
				}

			}
		}
	}

	return err
}

func elementExistsInSlice(slice []string, element string) (exists bool) {
	for _, val := range slice {
		if element == val {
			exists = true
			break
		}
	}

	return exists
}
