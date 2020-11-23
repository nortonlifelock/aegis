package nexpose

// Asset creates the structure of a json request to the API endpoint that can then be marshalled from
// the nexpose return for use with asset endpoints
type Asset struct {

	// All addresses discovered on the asset.
	Addresses []struct {

		// The IPv4 or IPv6 address.
		IP string `json:"ip,omitempty"`

		// The Media Access Control (MAC) address. The format is six groups of two hexadecimal digits separated by colons.
		Mac string `json:"mac,omitempty"`
	} `json:"addresses,omitempty"`

	// The history of changes to the asset over time.
	History []struct {

		// The date the asset information was collected or changed.
		Date string `json:"date,omitempty"`

		// Additional information describing the change.
		Description string `json:"description,omitempty"`

		// If a scan-oriented change, the identifier of the corresponding scan the asset was scanned in.
		ScanID int64 `json:"scanId,omitempty"`

		// Type identifies the type of the change
		// ASSET-IMPORT
		// EXTERNAL-IMPORT - External source such as the API
		// EXTERNAL-IMPORT-APPSPIDER - Rapid7 InsightAppSec (previously known as AppSpider)
		// SCAN - Scan engine scan
		// ACTIVE-SYNC - ActiveSync
		// SCAN-LOG-IMPORT - Manual import of a scan log
		// VULNERABILITY_EXCEPTION_APPLIED - Vulnerability exception applied
		// VULNERABILITY_EXCEPTION_UNAPPLIED - Vulnerability exception unapplied
		Type string `json:"type,omitempty"`

		// If a vulnerability exception change, the login name of the user that performed the operation.
		User string `json:"user,omitempty"`

		// The version number of the change (a chronological incrementing number starting from 1).
		Version int32 `json:"version,omitempty"`

		// If a vulnerability exception change, the identifier of the vulnerability exception that caused the change.
		VulnerabilityExceptionID int32 `json:"vulnerabilityExceptionId,omitempty"`
	} `json:"history,omitempty"`

	// The primary host name (local or FQDN) of the asset.
	HostName string `json:"hostName,omitempty"`

	// All host names or aliases discovered on the asset.
	HostNames []struct {

		// The host name (local or FQDN).
		Name string `json:"name"`

		// The source used to detect the host name. `user` indicates the host name source is user-supplied
		// (e.g. in a site target definition).
		Source string `json:"source,omitempty"`
	} `json:"hostNames,omitempty"`

	// The identifier of the asset.
	ID int `json:"id,omitempty"`

	// Unique identifiers found on the asset, such as hardware or operating system identifiers.
	Ids []struct {

		// The unique identifier.
		ID string `json:"id"`

		// The source of the unique identifier.
		Source string `json:"source,omitempty"`
	} `json:"ids,omitempty"`

	// The primary IPv4 or IPv6 address of the asset.
	IP string `json:"ip,omitempty"`

	// Hypermedia links to corresponding or related resources.
	Links []Link `json:"links,omitempty"`

	// The primary Media Access Control (MAC) address of the asset. The format is six groups of two hexadecimal
	// digits separated by colons.
	MAC string `json:"mac,omitempty"`

	// The full description of the operating system of the asset.
	OS string `json:"os,omitempty"`

	Vulnerabilities []Finding `json:"-"`
}
