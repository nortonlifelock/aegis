package nexpose

// Link is a struct representation of several returns from the nexpose api
type Link struct {

	// A hypertext reference, which is either a URI (see <a target=\"_blank\" href=\"https://tools.ietf.org/html/rfc3986\">RFC 3986</a>) or URI template (see <a target=\"_blank\" href=\"https://tools.ietf.org/html/rfc6570\">RFC 6570</a>).
	Href string `json:"href,omitempty"`

	// The link relation type. This value is one from the <a target=\"_blank\" href=\"https://tools.ietf.org/html/rfc5988#section-6.2\">Link Relation Type Registry</a> or is the type of resource being linked to.
	Rel string `json:"rel,omitempty"`
}

// Links is a struct representation of several returns from the nexpose api
type Links struct {

	// Hypermedia links to corresponding or related resources.
	Links []Link `json:"links,omitempty"`
}

// Reference is a struct representation of several returns from the nexpose api
// the ID field is an integer
type Reference struct {

	// The identifier of the resource.
	ID interface{} `json:"id,omitempty"`

	// Hypermedia links to corresponding or related resources.
	Links []Link `json:"links,omitempty"`

	// The identifiers of the associated resources.
	Resources []string `json:"resources,omitempty"`
}
