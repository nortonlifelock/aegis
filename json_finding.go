package nexpose

// Finding is the json struct representation of a vulnerabilities existence
// on an asset in nexpose
type Finding struct {

	// The identifier of the vulnerability.
	ID string `json:"id"`

	// The number of vulnerable occurrences of the vulnerability. This does not include `invulnerable` instances.
	Instances int32 `json:"instances"`

	// Hypermedia links to corresponding or related resources.
	Links []Link `json:"links,omitempty"`

	// The vulnerability check results for the finding. Multiple instances may be present if one or more checks fired, or a check has multiple independent results.
	Results []struct {

		// The identifier of the vulnerability check.
		CheckID string `json:"checkId,omitempty"`

		// If the result is vulnerable with exceptions applied, the identifier(s) of the exceptions actively applied to the result.
		Exceptions []int32 `json:"exceptions,omitempty"`

		// An additional discriminating key used to uniquely identify between multiple instances of results on the same finding.
		ID string `json:"key,omitempty"`

		// Hypermedia links to corresponding or related resources.
		Links []Link `json:"links,omitempty"`

		// The port of the service the result was discovered on.
		Port int32 `json:"port,omitempty"`

		// The proof explaining why the result was found vulnerable. The proof may container embedded HTML formatting markup.
		Proof string `json:"proof,omitempty"`

		// The protocol of the service the result was discovered on.
		Protocol string `json:"protocol,omitempty"`

		// The date and time the result was first recorded, in the ISO8601 format. If the result changes status this value is the date and time of the status change.
		Since string `json:"since,omitempty"`

		// The status of the vulnerability check result.
		Status string `json:"status"`
	} `json:"results,omitempty"`

	// The date and time the finding was was first recorded, in the ISO8601 format.
	// If the result changes status this value is the date and time of the status change.
	Since string `json:"since,omitempty"`

	// The status of the finding.
	Status string `json:"status"`
}
