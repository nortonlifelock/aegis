package jira

// MappableFields contains a list of all custom fields used by Aegis that may need to be mapped to a different name for a field that
// uses custom fields with different names
var MappableFields = []string{
	backendMOD,
	backendSummary,
	backendHostname,
	backendIPAddress,
	backendMACAddress,
	backendServicePort,
	backendDescription,
	backendSolution,
	backendVRRPriority,
	backendScanDate,
	backendAssignmentGroup,
	backendResolutionDate,
	backendReopenReason,
	backendOperatingSystem,
	backendVulnerability,
	backendCVSS,
	backendVulnerabilityID,
	backendOWASP,
	backendExceptionDate,
	backendGroupID,
	backendSystemName,
	backendPatchable,
	backendDeviceID,
	backendScanID,
	backendOrg,
	backendCERF,
	backendExceptionExpiration,
	backendCVEReferences,
	backendVendorReferences,
	backendOSDetailed,
	backendCategory,
	backendApplicationName,
	backendConfig,
	backendLastChecked,
	backendCloudID,
	backendHubProjectName,
	backendHubProjectVersion,
	backendHubSeverity,
	backendComponentName,
	backendComponentVersion,
	backendPolicyRule,
	backendPolicySeverity,
}

const (
	// DateFormatJira is used to parse the string dates that JIRA returns from API requests
	DateFormatJira = "2006-01-02T15:04:05.999-0700"

	// DateOnlyJira is used to set the date of JIRA tickets where more fine-grained dates is not needed
	DateOnlyJira = "2006-01-02"

	// QueryDateTimeFormatJira is the format of dates in a JQL
	QueryDateTimeFormatJira = "2006/01/02 15:04"
)

// API Endpoints
const (
	jsearch                = "/rest/api/latest/search"
	jfield                 = "/rest/api/latest/field"
	jresolution            = "/rest/api/2/resolution"
	jstatus                = "/rest/api/2/status"
	jissuetype             = "/rest/api/2/issuetype"
	jticket                = "/rest/api/2/issue/{issueIdOrKey}"
	jproject               = "/rest/api/2/project/{projectIdOrKey}"
	jeditablefields        = "/rest/api/2/issue/{issueIdOrKey}/editmeta"
	jcreateissue           = "rest/api/2/issue/"
	jeditableprojectfields = "/rest/api/2/issue/createmeta?projectKeys={project}&expand=projects.issuetypes.fields"
)

// The following consts correspond to ticket values that Aegis tracks. Because any given project can have varying fields, these
// consts are used in collaboration with the field map to find th equivalent custom field in any given JIRA project
const (
	backendKey                 = "Key"
	backendProject             = "Project"
	backendMOD                 = "Method of Discovery"
	backendStatus              = "Status"
	backendCreated             = "created"
	backendSummary             = "Summary"
	backendHostname            = "Hostname"
	backendIPAddress           = "IP Address"
	backendMACAddress          = "MAC Address"
	backendServicePort         = "Service Port"
	backendDescription         = "Description"
	backendSolution            = "Solution"
	backendVRRPriority         = "VRR Priority"
	backendScanDate            = "Scan/Alert Date"
	backendAssignmentGroup     = "Assignment Group"
	backendResolutionDate      = "Resolution Date"
	backendReopenReason        = "Reopen Reason"
	backendOperatingSystem     = "Operating System"
	backendVulnerability       = "Vulnerablility" // typo is present in JIRA and intentional
	backendCVSS                = "cvss_score"
	backendVulnerabilityID     = "VulnerabilityID"
	backendOWASP               = "OWASP"
	backendExceptionDate       = "Exception Date"
	backendGroupID             = "GroupID"
	backendDeviceID            = "DeviceID"
	backendPatchable           = "Patchable"
	backendScanID              = "ScanID"
	backendOrg                 = "Org"
	backendSystemName          = "SystemName"
	backendCERF                = "CERF"
	backendExceptionExpiration = "Exception Expiration"
	backendCVEReferences       = "CVE References"
	backendVendorReferences    = "VendorRef"
	backendUpdated             = "Updated"
	backendOSDetailed          = "OS_Detailed"
	backendCategory            = "Category"
	backendConfig              = "Config"
	backendLastChecked         = "LastChecked"
	backendCloudID             = "CloudID"
	backendApplicationName     = "Application Name"

	backendHubProjectName    = "Hub Project Name"
	backendHubProjectVersion = "Hub Project Version"
	backendHubSeverity       = "Hub Severity"
	backendComponentName     = "Component Name"
	backendComponentVersion  = "Component Version"
	backendPolicyRule        = "Policy Rule"
	backendPolicySeverity    = "Policy Severity"
)
