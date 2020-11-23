package nexpose

// Engine is the json struct representation of the nexpose scan engine
type Engine struct {

	// The address the scan engine is hosted.
	Address string `json:"address"`

	// The content version of the scan engine.
	ContentVersion string `json:"contentVersion,omitempty"`

	// A list of identifiers of engine pools this engine is included in.
	EnginePools []int32 `json:"enginePools,omitempty"`

	// The identifier of the scan engine.
	ID int `json:"id"`

	// The date the engine was last refreshed. Date format is in ISO 8601.
	LastRefreshedDate string `json:"lastRefreshedDate,omitempty"`

	// The date the engine was last updated. Date format is in ISO 8601.
	LastUpdatedDate string `json:"lastUpdatedDate,omitempty"`

	// Hypermedia links to corresponding or related resources.
	Links []Link `json:"links,omitempty"`

	// The name of the scan engine.
	Name string `json:"name"`

	// The port used by the scan engine to communicate with the Security Console.
	Port int32 `json:"port"`

	// The product version of the scan engine.
	ProductVersion string `json:"productVersion,omitempty"`

	// A list of identifiers of each site the scan engine is assigned to.
	Sites []int32 `json:"sites,omitempty"`
}
