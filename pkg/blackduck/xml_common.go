package blackduck

type Links struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

type Meta struct {
	Allow []string `json:"allow"`
	Href  string   `json:"href"`
	Links []Links  `json:"links"`
}
