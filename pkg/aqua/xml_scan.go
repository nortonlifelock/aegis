package aqua

import "time"

type ScanPage struct {
	Count    int    `json:"count"`
	Page     int    `json:"page"`
	Pagesize int    `json:"pagesize"`
	Result   []Scan `json:"result"`
}
type Scan struct {
	Registry       string      `json:"registry"`
	Image          string      `json:"image"`
	Created        time.Time   `json:"created"`
	LastUpdated    time.Time   `json:"last_updated"`
	StatusVar      string      `json:"status"`
	InitiatingUser string      `json:"initiating_user"`
	OsType         string      `json:"os_type"`
	OsVersions     interface{} `json:"os_versions"`
	Priority       int         `json:"priority"`
	ForcePull      bool        `json:"force_pull"`
	WebhookURL     string      `json:"webhook_url"`
	ScanError      string      `json:"scan_error"`
	HostID         string      `json:"host_id"`
	DockerID       string      `json:"docker_id"`
	EntityType     int         `json:"entity_type"`
	ImportedScan   interface{} `json:"imported_scan"`
}
