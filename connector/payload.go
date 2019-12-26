package connector

import "time"

// Payload is used to parse the information from the payload in the Nexpose source config entry
type Payload struct {

	// ScanTemplate is the name of the Nexpose template that should
	// be used for rescans.
	ScanTemplate string `json:"template,omitempty"`

	// ScanNameFormat holds the format string for the name of the scans that are created in Nexpose
	// should have a single '%s' symbol
	ScanNameFormat string `json:"scan_name_format"`

	// RescanSite specifies the site that will be used for rescans
	// A rescan site is required for proper rescans in nexpose as
	// this bypasses an issue with aws discovery sites and automated scans
	RescanSite int `json:"rescansite"`

	// DiscoveryTemplate is the name of the Nexpose template that should
	// be used for decommission rescans
	DiscoveryTemplate string `json:"discoverytemplate,omitempty"`

	// DiscoveryNameFormat holds the format string for the name of the discscans that are created in Nexpose
	//	// should have a single '%s' symbol
	DiscoveryNameFormat string `json:"discovery_name_format"`

	// EngineCacheTTL describes how long engine information should be cached in minutes
	EngineCacheTTL *time.Duration `json:"engine_cache_ttle,omitempty"`
}
