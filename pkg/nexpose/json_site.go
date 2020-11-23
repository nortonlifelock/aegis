package nexpose

// Site is the json representation of the nexpose scan site
type Site struct {

	// The number of assets that belong to the site.
	Assets int32 `json:"assets,omitempty"`

	// The type of discovery connection configured for the site. This property only applies to dynamic sites.
	ConnectionType string `json:"connectionType,omitempty"`

	// The site description.
	Description string `json:"description,omitempty"`

	// The identifier of the site.
	ID int `json:"id,omitempty"`

	// The site importance.
	Importance string `json:"importance,omitempty"`

	// The date and time of the site's last scan.
	LastScanTime string `json:"lastScanTime,omitempty"`

	// Hypermedia links to corresponding or related resources.
	Links []Link `json:"links,omitempty"`

	// The site name.
	Name string `json:"name,omitempty"`

	// The risk score (with criticality adjustments) of the site.
	RiskScore float32 `json:"riskScore,omitempty"`

	// The identifier of the scan engine configured in the site.
	ScanEngineID int `json:"scanEngine,omitempty"`

	// The identifier of the scan template configured in the site.
	ScanTemplateID string `json:"scanTemplate,omitempty"`

	// The type of the site.
	Type string `json:"type,omitempty"`
}
