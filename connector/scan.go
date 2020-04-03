package connector

import (
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/log"
	"github.com/nortonlifelock/nexpose"
	"github.com/pkg/errors"
	"strconv"
)

// Scan defines a struct for mapping scan information from nexpose so
// that the scan can be effectively tracked and cleaned up after the scan finishes
type Scan struct {
	api *nexpose.Session

	// The identifier of the scan engine.
	EngineID string `json:"engineId,omitempty"`

	// The hosts that should be included as a part of the scan. This should be a mixture of IP Addresses and Hostnames as a String array.
	IPs []string `json:"ips,omitempty"`

	// The asset ids that are included in this scan
	Assets []int `json:"assets,omitempty"`

	// The user-driven scan name for the scan.
	Name string `json:"name,omitempty"`

	// The identifier of the scan template
	TemplateID string `json:"templateId,omitempty"`

	// ScanID holds the scan's identifier as returned from nexpose
	ScanID string `json:"scanId,omitempty"`

	// AssetGroupID holds the ID of the asset group that the scan is being executed against
	AssetGroupID string `json:"groupId,omitempty"`

	// VulnerabilityIDs holds the list of vulnerability identifiers scanned for in the scan
	VulnerabilityIDs []string `json:"vulnerabilities,omitempty"`
}

// ID returns the scan id from nexpose
func (s *Scan) ID() string {
	return s.ScanID
}

func (s *Scan) Title() string {
	return s.Name
}

func (s *Scan) GroupID() string {
	return s.AssetGroupID
}

// Status returns the status of the scan in Nexpose
func (s *Scan) Status() (status string, err error) {

	// get the scan from the api
	var scan *nexpose.Scan
	if scan, err = s.api.GetScan(s.ScanID); err == nil {
		if scan != nil {
			// use the normalized scan status
			status = scanStatuses[scan.Status]
		} else {
			err = errors.Errorf("scan with id [%s] returned nil from nexpose", s.ScanID)
		}
	}

	return status, err
}

// Devices returns the list of ip addresses for the hosts in the scan
func (s *Scan) Devices() []string {
	return s.IPs
}

// Vulnerabilities returns the list of vulnerability ids for the scan
func (s *Scan) Vulnerabilities() []string {
	return s.VulnerabilityIDs
}

func (conn *Connection) createScanForDetections(detectionsToRescan []domain.Match, groupToRescan string, scanName string, scanTemplate string, rescanSite string) (scan *Scan, err error) {
	var vulns = make([]string, 0)
	var seenVuln = make(map[string]bool)

	var ips = make([]string, 0)
	var seenIP = make(map[string]bool)

	var devices = make([]int, 0)
	var seenDevice = make(map[string]bool)

	for _, detection := range detectionsToRescan {
		if !seenVuln[detection.Vulnerability()] {
			seenVuln[detection.Vulnerability()] = true
			vulns = append(vulns, detection.Vulnerability())
		}

		if !seenIP[detection.IP()] {
			seenIP[detection.IP()] = true
			ips = append(ips, detection.IP())
		}

		if !seenDevice[detection.Device()] {
			seenDevice[detection.Device()] = true

			if deviceID, err := strconv.Atoi(detection.Device()); err == nil {
				devices = append(devices, deviceID)
			} else {
				conn.logger.Send(log.Errorf(err, "error while parsing device id [%v]", detection.Device()))
			}
		}
	}

	var engineID string
	if engineID, err = conn.GetSiteEngine(groupToRescan); err == nil {

		var scanID string
		if scanID, err = conn.api.CreateScan(rescanSite, engineID, scanTemplate, scanName, ips); err == nil {
			scan = &Scan{
				api:              conn.api,
				EngineID:         engineID,
				IPs:              ips,
				Name:             scanName,
				TemplateID:       scanTemplate,
				ScanID:           scanID,
				AssetGroupID:     groupToRescan,
				Assets:           devices,
				VulnerabilityIDs: vulns,
			}
		} else {
			conn.logger.Send(log.Errorf(err, "error while creating Nexpose scan"))
		}
	} else {
		conn.logger.Send(log.Errorf(err, "error while gathering engines for site %v", groupToRescan))
	}

	return scan, err
}
