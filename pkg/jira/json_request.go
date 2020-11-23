package jira

import (
	"github.com/andygrunwald/go-jira"
	"time"
)

type postSearchRequest struct {
	JQL        string   `json:"jql,omitempty"`
	StartAt    int      `json:"startAt,omitempty"`
	MaxResults int      `json:"maxResults,omitempty"`
	Fields     []string `json:"Fields,omitempty"`
}

type assignmentGroupResponseWrapper struct {
	TotalResults int                       `json:"total"`
	Groups       []AssignmentGroupResponse `json:"groups"`
}

// AssignmentGroupResponse is a member of assignmentGroupResponseWrapper and must be exported in order to be marshaled
type AssignmentGroupResponse struct {
	Name string `json:"name"`
}

// CF is used for parsing custom field information from the JIRA API back to the JIRA driver
type CF struct {
	Value interface{} `json:"value,omitempty"`
	ID    int         `json:"id,omitempty"`
	Self  string      `json:"self,omitempty"`
	Name  interface{} `json:"name,omitempty"`
}

// Field is used for parsing field information from the JIRA API
type Field struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Custom      bool     `json:"custom,omitempty"`
	Orderable   bool     `json:"orderable,omitempty"`
	Navigable   bool     `json:"navigable,omitempty"`
	Searchable  bool     `json:"searchable,omitempty"`
	ClauseNames []string `json:"clauseNames,omitempty"`
	Schema      Schema   `json:"schema,omitempty"`
}

// Schema is used in field and must be exported in order to be marshaled
type Schema struct {
	Type     string `json:"type,omitempty"`
	Custom   string `json:"custom,omitempty"`
	CustomID int    `json:"customId,omitempty"`
}

type searchResult struct {
	Expand     string        `json:"expand,omitempty"`
	StartAt    int           `json:"startAt,omitempty"`
	MaxResults int           `json:"maxResults,omitempty"`
	Total      int           `json:"total,omitempty"`
	Issues     []*jira.Issue `json:"issues,omitempty"`
}

// ValueField is a member of FieldList and must be exported in order to be marshaled
type ValueField struct {
	Value string `json:"value,omitempty"`
	Name  string `json:"name,omitempty"`
}

// FieldList is a member of updateBlock and must be exported in order to be marshaled
type FieldList struct {
	UpdatedDate         *time.Time  `json:"updateddate,omitempty"`
	ResolutionDate      string      `json:"resolutiondate,omitempty"`
	Project             *string     `json:"project,omitempty"`
	Summary             *string     `json:"summary,omitempty"`
	ProposedExpiration  *time.Time  `json:"proposedexpiration,omitempty"`
	ReportedBy          *ValueField `json:"reportedby,omitempty"`
	Priority            *ValueField `json:"priority,omitempty"`
	MACAddress          *string     `json:"macaddress,omitempty"`
	Hostname            *string     `json:"hostname,omitempty"`
	GroupID             string      `json:"groupid,omitempty"`
	SystemName          *string     `json:"systemname,omitempty"`
	Category            *ValueField `json:"category,omitempty"`
	ApplicationName     *string     `json:"applicationname,omitempty"`
	ScanID              string      `json:"scanid,omitempty"`
	VulnerabilityID     string      `json:"vulnerabilityid,omitempty"`
	CreatedDate         *time.Time  `json:"created,omitempty"`
	Labels              *[]string   `json:"labels,omitempty"`
	CerfLink            string      `json:"cerflink,omitempty"`
	DeviceID            string      `json:"deviceid,omitempty"`
	AssignedTo          *Assignee   `json:"assignee,omitempty"`
	AssetsAffected      *string     `json:"assetsaffected,omitempty"`
	ScanErrata          *string     `json:"scanerrata,omitempty"`
	ResolutionStatus    *string     `json:"resolutionstatus,omitempty"`
	DueDate             *time.Time  `json:"duedate,omitempty"`
	AlertDate           *time.Time  `json:"alertdate,omitempty"`
	ServicePorts        *string     `json:"serviceports,omitempty"`
	TicketType          *string     `json:"issuetype,omitempty"`
	AssignmentGroup     *Assignee   `json:"assignmentgroup,omitempty"`
	MethodOfDiscovery   *ValueField `json:"methodofdiscovery,omitempty"`
	VulnerabilityTitle  *string     `json:"vulnerabilitytitle,omitempty"`
	CveReferences       *string     `json:"cve_references,omitempty"`
	IPAddress           *string     `json:"ipaddress,omitempty"`
	ExceptionDate       *time.Time  `json:"exceptiondate,omitempty"`
	ExceptionExpiration *time.Time  `json:"exceptionexpiration"`
	OWASP               *ValueField `json:"owasp,omitempty"`
	ID                  int32       `json:"id,omitempty"`
	Title               string      `json:"title,omitempty"`
	Status              *string     `json:"status,omitempty"`
	OperatingSystem     *ValueField `json:"operatingsystem,omitempty"`
	CVSS                *ValueField `json:"cvss,omitempty"`
	OrganizationID      int32       `json:"organizationid,omitempty"`
	OrgCode             *string     `json:"org,omitempty"`
	Description         *string     `json:"description,omitempty"`
	VendorReferences    *string     `json:"vendor_references,omitempty"`
	Solution            *string     `json:"solution,omitempty"`
	LastChecked         *time.Time  `json:"lastchecked,omitempty"`

	HubProjectName    *string `json:"hubprojectname"`
	HubProjectVersion *string `json:"hubprojectversion"`
	HubSeverity       *string `json:"hubseverity"`
	ComponentName     *string `json:"componentname"`
	ComponentVersion  *string `json:"componentversion"`
	PolicyRule        *string `json:"policyrule"`
	PolicySeverity    *string `json:"policyseverity"`
}

type updateBlock struct {
	Fields FieldList `json:"fields,omitempty"`
}
