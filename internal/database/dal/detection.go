package dal

import (
	"github.com/nortonlifelock/domain"
	"sync"
	"time"
)

// Detection holds information regarding the instance of a vulnerability on a particular device
type Detection struct {
	Conn domain.DatabaseConnection
	Info domain.DetectionInfo

	cacheDev  domain.Device
	cacheVuln domain.Vulnerability
	cacheLock sync.Mutex
}

// ID returns the ID of the detection as tracked by the Aegis database
func (detection *Detection) ID() string {
	return detection.Info.ID()
}

// VulnerabilityID returns the ID of the vulnerability as tracked by the vulnerability scanner that found it
func (detection *Detection) VulnerabilityID() string {
	var id string
	detection.cacheVuln, _ = detection.Vulnerability()
	if detection.cacheVuln != nil {
		id = detection.cacheVuln.SourceID()
	}

	return id
}

// Status returns the state of the vulnerability on the device (i.e. whether it is still active)
func (detection *Detection) Status() string {
	// default to vulnerable in case the detection loading fails
	var status = domain.Vulnerable
	detectionStatus, err := detection.Conn.GetDetectionStatusByID(detection.Info.DetectionStatusID())
	if err == nil {
		status = detectionStatus.Name()
	}
	return status
}

// ActiveKernel returns a nullable integer regarding to which kernel the vulnerability applies to (non-nil when there are multiple kernels on the device)
func (detection *Detection) ActiveKernel() *int {
	return detection.Info.ActiveKernel()
}

// Detected returns the date that the detection was identified
func (detection *Detection) Detected() (*time.Time, error) {
	date := detection.Info.AlertDate()
	return &date, nil
}

// TimesSeen returns the amount of times the detection has been identified according to the vulnerability scanner
func (detection *Detection) TimesSeen() int {
	return detection.Info.TimesSeen()
}

// Proof returns evidence that the vulnerability exists on the device
func (detection *Detection) Proof() string {
	return detection.Info.Proof()
}

// Port returns the port that the vulnerability applies to (when applicable)
func (detection *Detection) Port() int {
	return detection.Info.Port()
}

// Protocol returns the protocol that the vulnerability applies to (when applicable)
func (detection *Detection) Protocol() string {
	return detection.Info.Protocol()
}

// IgnoreID returns the ID of the ignore entry that pertains to the detection (if one exists(
func (detection *Detection) IgnoreID() (*string, error) {
	return detection.Info.IgnoreID(), nil
}

func (detection *Detection) Updated() time.Time {
	return detection.Info.Updated()
}

// Device returns on object implementing a corresponding Device interface that the detection exists on
func (detection *Detection) Device() (device domain.Device, err error) {
	if detection.cacheDev != nil {
		device = detection.cacheDev
	} else {
		detection.cacheLock.Lock()
		defer detection.cacheLock.Unlock()

		device, err = detection.Conn.GetDeviceByAssetOrgID(detection.Info.DeviceID(), detection.Info.OrganizationID())
	}

	return device, err
}

// Vulnerability returns on object implementing a corresponding Vulnerability interface that the detection applies to
func (detection *Detection) Vulnerability() (vulnerability domain.Vulnerability, err error) {
	if detection.cacheVuln != nil {
		vulnerability = detection.cacheVuln
	} else {
		detection.cacheLock.Lock()
		defer detection.cacheLock.Unlock()

		var vulnerabilityInfo domain.VulnerabilityInfo
		vulnerabilityInfo, err = detection.Conn.GetVulnInfoByID(detection.Info.VulnerabilityID())
		if err == nil && vulnerabilityInfo != nil {
			vulnerability = &Vulnerability{
				Conn: detection.Conn,
				Info: vulnerabilityInfo,
			}

			detection.cacheVuln = vulnerability
		}
	}

	return vulnerability, err
}
