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
	backendOperatingSystem,
	backendVulnerability,
	backendCVSS,
	backendVulnerabilityID,
	backendGroupID,
	backendDeviceID,
	backendScanID,
	backendOrg,
	backendCERF,
	backendCERFExpiration,
	backendCVEReferences,
	backendVendorReferences,
	backendOSDetailed,
	backendConfig,
	backendLastChecked,
	backendCloudID,
}

const (
	// DateFormatJira is used to parse the string dates that JIRA returns from API requests
	DateFormatJira = "2006-01-02T15:04:05.999-0700"

	// DateOnlyJira is used to set the date of JIRA tickets where more fine-grained dates is not needed
	DateOnlyJira = "2006-01-02"

	// QueryDateTimeFormatJira is the format of dates in a JQL
	QueryDateTimeFormatJira = "2006/01/02 15:04"
)

// Const for status maps
// These are used as KEYS to grab the equivalent mapped status from the JIRA payload
const (
	// MaxThreads defines the amount of threads create API calls for JIRA
	MaxThreads = 25

	// StatusReopened is the status of a ticket which a scanner confirmed its vulnerability still exists after it was marked resolved
	StatusReopened = "Reopened"

	// StatusClosedRemediated is the status of a ticket which was marked resolved, and had the vulnerability resolution confirmed by a scanner
	StatusClosedRemediated = "Closed-Remediated"

	// StatusClosedFalsePositive is the status of a ticket which has the vulnerability confirmed to be a false positive
	StatusClosedFalsePositive = "Closed-False-Positive"

	// StatusClosedDecommissioned is the status of a ticket for a vulnerability on a device that is no longer active
	StatusClosedDecommissioned = "Closed-Decommission"

	// StatusOpen is the status of a ticket which has not had any remediation steps taken
	StatusOpen = "Open"

	// StatusInProgress is the status of a ticket which is in the process of remediation
	StatusInProgress = "In-Progress"

	// StatusResolvedException is the status of a ticket which does not need to have the vulnerability remediated but not confirmed by a scanner
	StatusResolvedException = "Resolved-Exception"

	// StatusClosedException is the status of a ticket which does not need to have the vulnerability remediated
	StatusClosedException = "Closed-Exception"

	// StatusResolvedDecom is the status of a ticket for a vulnerability on a device that is no longer active but not confirmed by a scanner
	StatusResolvedDecom = "Resolved-Decommissioned"

	// StatusResolvedRemediated is the status of a ticket which was marked resolved but not confirmed by a scanner
	StatusResolvedRemediated = "Resolved-Remediated"

	// StatusResolvedFalsePositive is the status of a ticket which has the vulnerability confirmed to be a false positive but not verified by a scanner
	StatusResolvedFalsePositive = "Resolved-FalsePositive"

	// StatusClosedCerf is the status of a ticket that was closed due to an associated CERF (Candidate Exception Request Form)
	StatusClosedCerf = "Closed-CERF"

	// StatusClosedError is
	StatusClosedError = "Closed-Error"

	// StatusScanError is used for denoting when a scan failed to reach a device
	StatusScanError = "Scan-Error"
)

// API Endpoints
const (
	jsearch                = "/rest/api/latest/search"
	jfield                 = "/rest/api/latest/field"
	jresolution            = "/rest/api/2/resolution"
	jstatus                = "/rest/api/2/status"
	jissuetype             = "/rest/api/2/issuetype"
	jticket                = "/rest/api/2/issue/{issueIdOrKey}"
	jeditablefields        = "/rest/api/2/issue/{issueIdOrKey}/editmeta"
	jcreateissue           = "rest/api/2/issue/"
	jeditableprojectfields = "/rest/api/2/issue/createmeta?projectKeys={project}&expand=projects.issuetypes.fields"
)

// The following consts correspond to ticket values that Aegis tracks. Because any given project can have varying fields, these
// consts are used in collaboration with the field map to find th equivalent custom field in any given JIRA project
const (
	backendKey              = "Key"
	backendProject          = "Project"
	backendMOD              = "Method of Discovery"
	backendStatus           = "Status"
	backendCreated          = "created"
	backendSummary          = "Summary"
	backendHostname         = "Hostname"
	backendIPAddress        = "IP Address"
	backendMACAddress       = "MAC Address"
	backendServicePort      = "Service Port"
	backendDescription      = "Description"
	backendSolution         = "Solution"
	backendVRRPriority      = "VRR Priority"
	backendScanDate         = "Scan/Alert Date"
	backendAssignmentGroup  = "Assignment Group"
	backendResolutionDate   = "Resolution Date"
	backendOperatingSystem  = "Operating System"
	backendVulnerability    = "Vulnerablility" // typo is present in JIRA and intentional
	backendCVSS             = "cvss_score"
	backendVulnerabilityID  = "VulnerabilityID"
	backendGroupID          = "GroupID"
	backendDeviceID         = "DeviceID"
	backendScanID           = "ScanID"
	backendOrg              = "Org"
	backendAutomationID     = "AutomationID"
	backendCERF             = "CERF"
	backendCERFExpiration   = "Actual Expiration Date"
	backendCVEReferences    = "CVE References"
	backendVendorReferences = "VendorRef"
	backendUpdated          = "Updated"
	backendOSDetailed       = "OS_Detailed"
	backendConfig           = "Config"
	backendLastChecked      = "LastChecked"
	backendCloudID          = "CloudID"
)
