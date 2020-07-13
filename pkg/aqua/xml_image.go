package aqua

import "time"

type ImagePage struct {
	Count    int           `json:"count"`
	Page     int           `json:"page"`
	Pagesize int           `json:"pagesize"`
	Result   []ImageResult `json:"result"`
}
type PolicyFailures struct {
	Blocking   bool     `json:"blocking"`
	Controls   []string `json:"controls"`
	PolicyID   int      `json:"policy_id"`
	PolicyName string   `json:"policy_name"`
}
type AssuranceResults struct {
	Disallowed      bool        `json:"disallowed"`
	ChecksPerformed interface{} `json:"checks_performed"`
}
type ImageResult struct {
	Registry             string      `json:"registry"`
	Name                 string      `json:"name"`
	VulnsFound           int         `json:"vulns_found"`
	CritVulns            int         `json:"crit_vulns"`
	HighVulns            int         `json:"high_vulns"`
	MedVulns             int         `json:"med_vulns"`
	LowVulns             int         `json:"low_vulns"`
	NegVulns             int         `json:"neg_vulns"`
	RegistryType         string      `json:"registry_type"`
	Repository           string      `json:"repository"`
	Tag                  string      `json:"tag"`
	Created              time.Time   `json:"created"`
	Author               string      `json:"author"`
	Digest               string      `json:"digest"`
	Size                 int         `json:"size"`
	Labels               interface{} `json:"labels"`
	Os                   string      `json:"os"`
	OsVersion            string      `json:"os_version"`
	ScanStatus           string      `json:"scan_status"`
	ScanDate             time.Time   `json:"scan_date"`
	ScanError            string      `json:"scan_error"`
	SensitiveData        int         `json:"sensitive_data"`
	Malware              int         `json:"malware"`
	Disallowed           bool        `json:"disallowed"`
	Whitelisted          bool        `json:"whitelisted"`
	Blacklisted          bool        `json:"blacklisted"`
	PermissionLastupdate int         `json:"permission_lastupdate"`
	PermissionAuthor     string      `json:"permission_author"`
	//PolicyFailures        []PolicyFailures `json:"policy_failures"`
	PartialResults        bool             `json:"partial_results"`
	NewerImageExists      bool             `json:"newer_image_exists"`
	AssuranceResults      AssuranceResults `json:"assurance_results"`
	PendingDisallowed     bool             `json:"pending_disallowed"`
	MicroenforcerDetected bool             `json:"microenforcer_detected"`
}
