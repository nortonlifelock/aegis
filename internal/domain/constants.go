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
