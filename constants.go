package domain

// Job statuses
//noinspection GoUnusedConst,GoUnusedConst,GoUnusedConst,GoUnusedConst
const (
	JobStatusPending    = 1
	JobStatusInProgress = 2
	JobStatusCompleted  = 3
	JobStatusError      = 4
	JobStatusCancelled  = 5
)

// Normalized Scan Statuses
const (
	ScanQUEUED     = "queued"
	ScanPROCESSING = "processing"
	ScanPAUSED     = "paused"
	ScanFINISHED   = "finished"
	ScanERRORED    = "error"
	ScanSTOPPED    = "stopped"
	ScanCANCELED   = "canceled"
)

// Statuses that Aegis utilizes for vulnerability tickets
// These consts are used in status mapping as KEYS to grab the equivalent mapped status from the JIRA payload
const (
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

const (
	// Vulnerable denotes that the vulnerability is not fixed
	Vulnerable = "vulnerable"

	// Exceptioned denotes that the vulnerability can be ignored
	Exceptioned = "exceptioned"

	// Fixed denotes that the vulnerability is no longer present
	Fixed = "fixed"

	// DeadHost denotes that the vulnerability exists on a host that is no longer online
	DeadHost = "dead host"
)

const (
	// DeviceRunning denotes that the device is online
	DeviceRunning = "running"

	// DeviceStopped denotes that the device is offline but still exists
	DeviceStopped = "stopped"

	// DeviceDecommed denotes that the device no longer exists
	DeviceDecommed = "decommissioned"
)

// Ignore Types
const (
	// Exception delineates an entry in an ignore table that is an exception
	Exception = iota

	// FalsePositive delineates an entry in an ignore table that is a false positive
	FalsePositive

	// DecommAsset delineates an entry in an ignore table that is a decommissioned asset
	DecommAsset
)

const (
	// RescanExceptions is a constant that dictates the type of rescan job is currently running. This controls, for example, the types of tickets collected
	RescanExceptions = "EXCEPTIONS"

	// RescanPassive is a constant that dictates the type of rescan job is currently running. This controls, for example, the types of tickets collected
	RescanPassive = "PASSIVE"

	// RescanNormal is a constant that dictates the type of rescan job is currently running. This controls, for example, the types of tickets collected
	RescanNormal = "NORMAL"

	// RescanDecommission is a constant that dictates the type of rescan job is currently running. This controls, for example, the types of tickets collected
	RescanDecommission = "DECOMMISSIONED"

	// RescanScheduled is a constant that dictates the type of rescan job is currently running. Scheduled rescans are treated like normal rescans
	RescanScheduled = "SCHEDULED"
)

// Reference Types string to filter references coming from Nexpose/Qualys api
const (
	// MSType corresponds to a vendor reference that is a Microsoft Security bulletin
	MSType = "ms"

	// CVEType corresponds to a vendor reference that is a Common Vulnerability Exposure
	CVEType = "cve"

	// CVEPrefix is used to find the prefix as a substring in the title of a reference
	CVEPrefix = "cve-"
)

// Reference Types Enum
const (
	// CVE is the ID that corresponds to a CVE entry in the VulnerabilityReference table
	CVE = 0

	// MS is the ID that corresponds to a MS entry in the VulnerabilityReference table
	MS = 1

	// Vendor is the ID that corresponds to a generic entry in the VulnerabilityReference table
	Vendor = 2
)

const (
	// VulnPathConcatenator is used for when the same vulnerability is present on multiple paths. The concatenator joins the vulnerability ID and the path to ensure no instance of the vulnerability is skipped
	VulnPathConcatenator = ";"
)
