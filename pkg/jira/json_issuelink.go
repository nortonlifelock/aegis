package jira

type IssueLinkRequest struct {
	//Fields Fields `json:"fields"`
	Update *UpdateIssueLinks `json:"update,omitempty"`
}
//type Project struct {
//	Key string `json:"key"`
//}
//type Issuetype struct {
//	Name string `json:"name"`
//}
//type Priority struct {
//	Name string `json:"name"`
//}
//type Fields struct {
//	Project     Project   `json:"project"`
//	Summary     string    `json:"summary"`
//	Description string    `json:"description"`
//	Issuetype   Issuetype `json:"issuetype"`
//	Priority    Priority  `json:"priority"`
//}

type Type struct {
	Name    string `json:"name"`
	Inward  string `json:"inward"`
	Outward string `json:"outward"`
}
type OutwardIssue struct {
	Key string `json:"key"`
}
type Add struct {
	Type         Type         `json:"type"`
	OutwardIssue OutwardIssue `json:"outwardIssue"`
}
type Issuelinks struct {
	Add Add `json:"add"`
}
type UpdateIssueLinks struct {
	Issuelinks []Issuelinks `json:"issuelinks"`
}