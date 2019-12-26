package connector

import (
	"fmt"
	"github.com/benjivesterby/validator"
	"github.com/pkg/errors"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/nexpose"
	"time"
)

type detection struct {
	asset           *asset
	conn            *Connection
	vulnerabilityID string
	resultID        string
	status          string
	detected        string
	proof           string
	port            int
	protocol        string
}

// ID returns the PDE database ID, which is not present in the Nexpose object
func (d *detection) ID() string {
	return ""
}

// VulnerabilityID concatenates the result id if it is present
func (d *detection) VulnerabilityID() (ID string) {

	ID = d.vulnerabilityID
	if len(d.resultID) > 0 {
		// Hash the result id
		var resultID string

		// TODO: Add hashing here
		// concat the two together
		ID = fmt.Sprintf("%s;%s", ID, resultID)
	}

	return ID
}

// Status returns the current vulnerability status on an asset
func (d *detection) Status() string {
	return detectionStatus(d.status)
}

// Detected returns the date the vulnerability was first detected on the asset
func (d *detection) Detected() (*time.Time, error) {
	val, err := time.Parse(time.RFC3339, d.detected)
	return &val, err
}

// TimesSeen is not implemented in nexpose
func (d *detection) TimesSeen() int { return 1 }

// Proof returns the proof that the vulnerability exists on this asset
func (d *detection) Proof() (param string) {
	return d.proof
}

func (d *detection) Port() int {
	return d.port
}

func (d *detection) Protocol() string {
	return d.protocol
}

// ActiveKernel is not implemented in nexpose
func (d *detection) ActiveKernel() *int {
	return nil
}

func (d *detection) Device() (device domain.Device, err error) {
	return d.asset, err
}

func (d *detection) Vulnerability() (v domain.Vulnerability, err error) {

	if validator.IsValid(d) {

		if obj, ok := d.conn.vulnerabilities.Load(d.vulnerabilityID); ok {
			if validator.IsValid(obj) {
				var vuln domain.Vulnerability
				if vuln, ok = obj.(domain.Vulnerability); ok {
					v = vuln
				}
			}
		}

		// Load it from the API because it's not in the cache
		if v == nil {

			var vuln *nexpose.Vulnerability
			if vuln, err = d.conn.api.GetVulnerability(d.vulnerabilityID); err == nil {

				if vuln != nil {

					v = &vulnerability{
						vuln:   vuln,
						api:    d.conn.api,
						logger: d.conn.logger,
					}

					d.conn.vulnerabilities.Store(d.conn.ctx, d.vulnerabilityID, v, d.conn.ttl())
				} else {
					err = errors.Errorf("vulnerability [%s] returned nil from nexpose api", d.vulnerabilityID)
				}
			}
		}
	} else {
		err = fmt.Errorf("detection did not pass validation")
	}

	return v, err
}

func (d *detection) Validate() (valid bool) {
	if d != nil {
		if validator.IsValid(d.asset) &&
			validator.IsValid(d.conn.api) &&
			len(d.vulnerabilityID) > 0 {

			valid = true
		}
	}

	return valid
}
