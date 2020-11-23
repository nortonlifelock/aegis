package connector

import (
	"github.com/nortonlifelock/aegis/pkg/domain"
	"github.com/nortonlifelock/aegis/pkg/qualys"
	"strings"
	"time"
)

type scan struct {
	session *QsSession

	Name string ` json:"name,omitempty"`

	ScanID string `json:"scanId,omitempty"`

	// The identifier of the scan template
	TemplateID string `json:"templateId,omitempty"`

	// AssetGroupID holds the ID of the asset group that the scan is being executed against
	AssetGroupID string `json:"groupId,omitempty"`

	// The identifier of the scan engine.
	EngineIDs []string `json:"engineId,omitempty"`

	Created time.Time `json:"created"`

	Scheduled bool `json:"scheduled"`

	// matches holds the device/vuln combos that are covered in the scan
	// they are not included in the json intentionally as it is not required
	matches []domain.Match
}

func (s *scan) ID() string {
	return s.ScanID
}

func (s *scan) Title() string {
	return s.Name
}

func (s *scan) GroupID() string {
	return s.AssetGroupID
}

func (s *scan) Matches() []domain.Match {
	return s.matches
}

func (s *scan) Status() (status string, err error) {
	if !strings.Contains(s.ScanID, webPrefix) {
		var scan qualys.ScanQualys
		scan, err = s.session.apiSession.GetScanByReference(s.ScanID)
		if err == nil {
			status = scan.Status.State

			if strings.ToLower(status) == strings.ToLower(domain.ScanFINISHED) {
				if scan.Processed == 0 {
					status = domain.ScanPROCESSING
				}
			}
		}
	} else {
		// WAS Retest - the finding UID is stored in the TemplateID field
		// the only way I've found that we can tell that a previously kicked off retest is finished
		// is by attempting to kickoff a new restest. if the old retest is finished, a new one will be kicked off
		var count string
		count, err = s.session.apiSession.CreateRetestForWebAppVulnerabilityFinding(s.TemplateID)
		if count == "0" {
			status = domain.ScanPROCESSING
		} else {
			status = domain.ScanFINISHED
		}
	}

	status = strings.ToLower(status)

	return status, err
}
