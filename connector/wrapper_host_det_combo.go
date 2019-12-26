package connector

import (
	"strings"
	"time"

	"github.com/nortonlifelock/domain"
)

type hostDetectionCombo struct {
	host      *host
	detection *detection
}

// ID returns the PDE database ID, which is not present in the Qualys object
func (combo *hostDetectionCombo) ID() string {
	return ""
}

// VulnerabilityID returns the vulnerability ID of the host/detection combo
func (combo *hostDetectionCombo) VulnerabilityID() string {
	return combo.detection.SourceID()
}

// Status dictates whether a detection is considered active or not
func (combo *hostDetectionCombo) Status() string {
	var status = combo.detection.d.Status

	switch strings.ToLower(combo.detection.d.Status) {
	case "new":
		status = domain.Vulnerable
	case "active":
		status = domain.Vulnerable
	case "re-opened":
		status = domain.Vulnerable
	case "fixed":
		status = domain.Fixed
	default:
		// do nothing
	}

	return status
}

// ActiveKernel refers whether or not the detection applies to the kernel that is currently running on the host. There are three values this can take
// nil     - the detection
// 0       - the detection
// nonzero - the detection
func (combo *hostDetectionCombo) ActiveKernel() *int {
	return combo.detection.d.AffectsRunningKernel
}

// Detected returns the date that the detection was last found
func (combo *hostDetectionCombo) Detected() (*time.Time, error) {
	val := combo.detection.Updated()
	return &val, nil
}

func (combo *hostDetectionCombo) TimesSeen() int {
	return combo.detection.d.TimeFound
}

func (combo *hostDetectionCombo) Proof() string {
	return combo.detection.d.Proof
}

func (combo *hostDetectionCombo) Port() int {
	var port int
	if combo.detection.d.Port != nil {
		port = *combo.detection.d.Port
	}
	return port
}

func (combo *hostDetectionCombo) Protocol() string {
	var protocol = ""
	if combo.detection.d.Protocol != nil {
		protocol = *combo.detection.d.Protocol
	}
	return protocol
}

func (combo *hostDetectionCombo) Device() (domain.Device, error) {
	return combo.host, nil
}

func (combo *hostDetectionCombo) Vulnerability() (domain.Vulnerability, error) {
	combo.detection.lazyLoadVulnerabilityInfo()
	return combo.detection, nil
}
