package nexpose

// AdhocScan is the json struct representation for the scan creation body for creating a new adhoc scan in nexpose
type AdhocScan struct {

	// The identifier of the scan engine.
	EngineID string `json:"engineId,omitempty"`

	// The hosts that should be included as a part of the scan. This should be a mixture of IP Addresses and Hostnames as a String array.
	Hosts []string `json:"hosts,omitempty"`

	// The user-driven scan name for the scan.
	Name string `json:"name,omitempty"`

	// The identifier of the scan template
	TemplateID string `json:"templateId,omitempty"`
}

// Scan is the json struct representation of a scan record in nexpose
type Scan struct {

	// The number of assets found in the scan.
	Assets int32 `json:"assets,omitempty"`

	// The duration of the scan in ISO8601 format.
	Duration string `json:"duration,omitempty"`

	// The end time of the scan in ISO8601 format.
	EndTime string `json:"endTime,omitempty"`

	// The identifier of the scan engine.
	EngineID int `json:"engineId,omitempty"`

	// The name of the scan engine.
	EngineName string `json:"engineName,omitempty"`

	// The identifier of the scan.
	ID int `json:"id,omitempty"`

	// Hypermedia links to corresponding or related resources.
	Links []Link `json:"links,omitempty"`

	// The reason for the scan status.
	Message string `json:"message,omitempty"`

	// The user-driven scan name for the scan.
	ScanName string `json:"scanName,omitempty"`

	// The scan type (automated, manual, scheduled).
	ScanType string `json:"scanType,omitempty"`

	// The start time of the scan in ISO8601 format.
	StartTime string `json:"startTime,omitempty"`

	// The name of the user that started the scan.
	StartedBy string `json:"startedBy,omitempty"`

	// The scan status.
	Status string `json:"status,omitempty"`
}
