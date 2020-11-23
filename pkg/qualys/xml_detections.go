package qualys

import (
	"encoding/xml"
	"time"
)

type simpleReturn struct {
	XMLName  xml.Name `xml:"SIMPLE_RETURN"`
	Response SimpleReturnResponse
}

// QBaseInfo is a member of many structs that need to be marshaled and must be exported
type QBaseInfo struct {
	XMLName xml.Name `xml:"CUSTOM"`
	ID      int      `xml:"ID"`
	Title   *CData   `xml:"TITLE"`
	Type    *string  `xml:"TYPE,omitempty"`
}

// SimpleReturnItem is a member of SimpleReturnResponse and must be exported in order to be marshaled
type SimpleReturnItem struct {
	Key   string `xml:"KEY"`
	Value string `xml:"VALUE"`
}

// SimpleReturnResponse is a member of simpleReturn and must be exported in order to be marshaled
type SimpleReturnResponse struct {
	XMLName xml.Name           `xml:"RESPONSE"`
	Date    string             `xml:"DATETIME"`
	Code    int                `xml:"CODE"`
	Message string             `xml:"TEXT"`
	Items   []SimpleReturnItem `xml:"ITEM_LIST>ITEM"`
}

// QScanListOutput holds a list of scans in qualys
type QScanListOutput struct {
	XMLName  xml.Name `xml:"SCAN_LIST_OUTPUT"`
	Response ScanListOutputResponse
}

// ScanListOutputResponse is a member of QScanListOutput and must be exported in order to be marshaled
type ScanListOutputResponse struct {
	XMLName xml.Name     `xml:"RESPONSE"`
	Date    time.Time    `xml:"DATETIME"` // Will need to parse to date
	Scans   []ScanQualys `xml:"SCAN_LIST>SCAN"`
}

// ScanQualys is a member of ScanListOutputResponse and must be exported in order to be marshaled
type ScanQualys struct {
	XMLName    xml.Name  `xml:"SCAN"`
	Reference  string    `xml:"REF"`
	Type       string    `xml:"TYPE"`
	Title      string    `xml:"TITLE"`
	User       string    `xml:"USER_LOGIN"`
	LaunchDate time.Time `xml:"LAUNCH_DATETIME"`
	Duration   string    `xml:"DURATION"`
	Priority   string    `xml:"PROCESSING_PRIORITY"` // Will need to parse to in enum priority
	Processed  int       `xml:"PROCESSED"`
	Status     ScanStatusQualys
	Target     string `xml:"TARGET"`
}

// ScanStatusQualys is a member of ScanQualys and must be exported in order to be marshaled
type ScanStatusQualys struct {
	XMLName xml.Name `xml:"STATUS"`
	State   string   `xml:"STATE"`
}

// QHostListDetectionOutput holds vulnerability information pertaining to the hosts
type QHostListDetectionOutput struct {
	XMLName xml.Name  `xml:"HOST_LIST_VM_DETECTION_OUTPUT"`
	Date    time.Time `xml:"RESPONSE>DATETIME"`
	Hosts   []QHost   `xml:"RESPONSE>HOST_LIST>HOST"`
	Warning *QWarning `xml:"RESPONSE>WARNING,omitempty"`
}

// QHost is a member of qHostListOutput and must be exported in order to be marshaled
type QHost struct {
	XMLName                xml.Name     `xml:"HOST"`
	HostID                 int          `xml:"ID"`
	IPAddress              string       `xml:"IP"`
	TrackingMethod         string       `xml:"TRACKING_METHOD"`
	OperatingSystem        CData        `xml:"OS"`
	DNS                    CData        `xml:"DNS"`
	Netbios                string       `xml:"NETBIOS"`
	LastScan               time.Time    `xml:"LAST_SCAN_DATETIME"`
	LastVMScan             time.Time    `xml:"LAST_VM_SCANNED_DATE"`
	LastVMScanDuration     int          `xml:"LAST_VM_SCANNED_DURATION"`
	LastVMAuthScan         *time.Time   `xml:"LAST_VM_AUTH_SCANNED_DATE,omitempty"`
	LastVMAuthScanDuration *int         `xml:"LAST_VM_AUTH_SCANNED_DURATION,omitempty"`
	LastPCScan             time.Time    `xml:"LAST_PC_SCANNED_DATE"`
	Detections             []QDetection `xml:"DETECTION_LIST>DETECTION"`
	NetworkID              *string      `xml:"NETWORK_ID,omitempty"`
	AssetGroupIDs          string       `xml:"ASSET_GROUP_IDS"`

	// Host List Additions

	EC2Id string `xml:"EC2_INSTANCE_ID"`
	//QGHostI					string					`xml:"QG_HOSTID, omitempty"`
	//Tags					string					`xml:"TAGS>TAG, omitempty"` // TODO
	Metadata []interface{} `xml:"METADATA>EC2, omitempty"`
}

// QDetection is a member of QHost and must be exported in order to be marshaled
type QDetection struct {
	XMLName  xml.Name `xml:"DETECTION"`
	QualysID int      `xml:"QID"`
	Type     string   `xml:"TYPE"`
	Severity int      `xml:"SEVERITY"`
	Port     *int     `xml:"PORT,omitempty"`
	Protocol *string  `xml:"PROTOCOL,omitempty"`
	SSL      bool     `xml:"SSL"`
	Proof    string   `xml:"RESULTS"`

	// Status can hold values [New, Active, Re-Opened, Fixed]
	// when holding multiple values, it is displayed in CSV
	Status string `xml:"STATUS"`

	FirstFound           time.Time  `xml:"FIRST_FOUND_DATETIME"`
	LastFound            time.Time  `xml:"LAST_FOUND_DATETIME"`
	LastCheck            time.Time  `xml:"LAST_TEST_DATETIME"`
	LastUpdate           time.Time  `xml:"LAST_UPDATE_DATETIME"`
	LastFixed            *time.Time `xml:"LAST_FIXED_DATETIME,omitempty"`
	Ignored              bool       `xml:"IS_IGNORED"`
	Disabled             bool       `xml:"IS_DISABLED"`
	TimeFound            int        `xml:"TIMES_FOUND"`
	AffectsRunningKernel *int       `xml:"AFFECT_RUNNING_KERNEL,omitempty"`
}

// QWarning is a member of several other structs and must be exported in order to be marshaled
type QWarning struct {
	XMLName xml.Name `xml:"WARNING"`
	Code    int      `xml:"CODE"`
	Text    string   `xml:"TEXT"`
	URL     string   `xml:"URL"`
}

// QKnowledgeBaseVulnOutput holds a list of all vulnerabilities in the Qualys vulnerability base
type QKnowledgeBaseVulnOutput struct {
	XMLName         xml.Name         `xml:"KNOWLEDGE_BASE_VULN_LIST_OUTPUT"`
	Date            time.Time        `xml:"RESPONSE>DATETIME"`
	Vulnerabilities []QVulnerability `xml:"RESPONSE>VULN_LIST>VULN"`
	Warning         *QWarning        `xml:"RESPONSE>WARNING"`
}

// QVulnerability is a member of QKnowledgeBaseVulnOutput and must be exported in order to be marshaled
type QVulnerability struct {
	XMLName                 xml.Name           `xml:"VULN"`
	QualysID                int                `xml:"QID"`
	Type                    string             `xml:"VULN_TYPE"`
	Severity                int                `xml:"SEVERITY_LEVEL"`
	Title                   string             `xml:"TITLE"`
	Category                string             `xml:"CATEGORY,omitempty"`
	DetectionInformation    string             `xml:"DETECTION_INFO,omitempty"`
	LastCustomization       *QChange           `xml:"LAST_CUSTOMIZATION,omitempty"`
	LastServiceModification time.Time          `xml:"LAST_SERVICE_MODIFICATION_DATETIME,omitempty"`
	Published               time.Time          `xml:"PUBLISHED_DATETIME"`
	BugtraqList             []QBugTrack        `xml:"BUGTRAQ_LIST>BUGTRAQ,omitempty"`
	Patchable               bool               `xml:"PATCHABLE"`
	Software                []QSoftware        `xml:"SOFTWARE_LIST>SOFTWARE,omitempty"`
	VendorReferences        []QVendorReference `xml:"VENDOR_REFERENCE_LIST>VENDOR_REFERENCE,omitempty"`
	CVEs                    []QCVE             `xml:"CVE_LIST>CVE,omitempty"`
	Diagnosis               string             `xml:"DIAGNOSIS,omitempty"`
	DiagnosisComment        string             `xml:"DIAGNOSIS_COMMENT,omitempty"`
	Consequence             string             `xml:"CONSEQUENCE,omitempty"`
	ConsequenceComment      string             `xml:"CONSEQUENCE_COMMENT,omitempty"`
	Solution                string             `xml:"SOLUTION,omitempty"`
	SolutionComment         string             `xml:"SOLUTION_COMMENT,omitempty"`
	ComplianceList          []QCompliance      `xml:"COMPLIANCE_LIST>COMPLIANCE,omitempty"`
	Correlations            []QCorrelation     `xml:"CORRELATION>CORRELATION,omitempty"`
	CVSS                    *QCVSS             `xml:"CVSS,omitempty"`
	CVSS3                   *QCVSS3            `xml:"CVSS_V3,omitempty"`
	PCI                     bool               `xml:"PCI_FLAG"`
	AutoPCIFail             bool               `xml:"AUTOMATIC_PCI_FAIL,omitempty"`
	PCIReasons              []string           `xml:"PCI_REASONS>PCI_REASON,omitempty"`
	SupportedModules        string             `xml:"SUPPORTED_MODULES,omitempty"`
	Discovery               QDiscovery         `xml:"DISCOVERY"`
	Disabled                bool               `xml:"IS_DISABLED,omitempty"`
	ThreatIntel             []QThreatIntel     `xml:"THREAT_INTELLIGENCE>THREAT_INTEL,omitempty"`
}

// QChange is a member of QVulnerability and must be exported in order to be marshaled
type QChange struct {
	Date time.Time `xml:"DATETIME"`
	User time.Time `xml:"USER_LOGIN"`
}

// QBugTrack is a member of QVulnerability and must be exported in order to be marshaled
type QBugTrack struct {
	XMLName xml.Name `xml:"BUGTRAQ"`
	ID      int      `xml:"ID"`
	URL     string   `xml:"URL"`
}

// QThreatIntel is a member of QVulnerability and must be exported in order to be marshaled
type QThreatIntel struct {
	XMLName xml.Name `xml:"THREAT_INTEL"`
	ID      int      `xml:"id,attr"`
	Intel   string   `xml:",innerxml"`
}

// QSoftware is a member of QVulnerability and must be exported in order to be marshaled
type QSoftware struct {
	XMLName xml.Name `xml:"SOFTWARE"`
	Product string   `xml:"PRODUCT"`
	Vendor  string   `xml:"VENDOR"`
}

// QVendorReference is a member of QVulnerability and must be exported in order to be marshaled
type QVendorReference struct {
	XMLName xml.Name `xml:"VENDOR_REFERENCE"`
	ID      string   `xml:"ID"`
	URL     string   `xml:"URL"`
}

// QCVE is a member of QVulnerability and must be exported in order to be marshaled
type QCVE struct {
	XMLName xml.Name `xml:"CVE"`
	ID      string   `xml:"ID"`
	URL     string   `xml:"URL"`
}

// QCompliance is a member of QVulnerability and must be exported in order to be marshaled
type QCompliance struct {
	XMLName     xml.Name `xml:"COMPLIANCE"`
	Type        string   `xml:"TYPE"`
	Section     string   `xml:"SECTION"`
	Description string   `xml:"DESCRIPTION"`
}

// QCorrelation is a member of QVulnerability and must be exported in order to be marshaled
type QCorrelation struct {
	XMLName     xml.Name       `xml:"CORRELATION"`
	ExploitList []QExploitList `xml:"EXPLOITS>EXPLT_SRC"`
	MalwareList []QMalwareList `xml:"MALWARE>MW_SRC"`
}

// QExploitList is a member of QCorrelation and must be exported in order to be marshaled
type QExploitList struct {
	XMLName  xml.Name   `xml:"EXPLT_SRC"`
	Source   string     `xml:"SRC_NAME"`
	Exploits []QExploit `xml:"EXPLT_LIST>EXPLT"`
}

// QMalwareList is a member of QCorrelation and must be exported in order to be marshaled
type QMalwareList struct {
	XMLName     xml.Name   `xml:"MW_SRC"`
	Source      string     `xml:"SRC_NAME"`
	MalwareList []QMalware `xml:"MW_LIST>MW_INFO"`
}

// QExploit is a member of QExploitList and must be exported in order to be marshaled
type QExploit struct {
	XMLName     xml.Name `xml:"EXPLT"`
	Reference   string   `xml:"REF"`
	Description string   `xml:"DESC"`
	Link        string   `xml:"LINK"`
}

// QMalware is a member of QMalwareList and must be exported in order to be marshaled
type QMalware struct {
	XMLName  xml.Name `xml:"MW_INFO"`
	ID       string   `xml:"MW_ID"`
	Type     string   `xml:"MW_TYPE,omitempty"`
	Platform string   `xml:"MW_PLATFORM,omitempty"`
	Alias    string   `xml:"MW_ALIAS,omitempty"`
	Rating   string   `xml:"MW_RATING,omitempty"`
	Link     string   `xml:"MW_LINK,omitempty"`
}

// QCVSS is a member of QVulnerability and must be exported in order to be marshaled
type QCVSS struct {
	XMLName          xml.Name     `xml:"CVSS"`
	Base             float32      `xml:"BASE"`
	Temporal         float32      `xml:"TEMPORAL,omitempty"`
	Access           *QCVSSAccess `xml:"ACCESS,omitempty"`
	Impact           *QCVSSImpact `xml:"IMPACT,omitempty"`
	Authentication   int          `xml:"AUTHENTICATION,omitempty"`
	Exploitability   int          `xml:"EXPLOITABILITY,omitempty"`
	RemediationLevel int          `xml:"REMEDIATION_LEVEL,omitempty"`
	ReportConfidence int          `xml:"REPORT_CONFIDENCE,omitempty"`
}

// QCVSS3 is a member of QVulnerability and must be exported in order to be marshaled
type QCVSS3 struct {
	QCVSS
	XMLName xml.Name `xml:"CVSS_V3"`
}

// QCVSSAccess is a member of QCVSS and must be exported in order to be marshaled
type QCVSSAccess struct {
	XMLName    xml.Name `xml:"ACCESS"`
	Vector     int      `xml:"VECTOR,omitempty"`
	Complexity int      `xml:"COMPLEXITY,omitempty"`
}

// QCVSSImpact is a member of QCVSS and must be exported in order to be marshaled
type QCVSSImpact struct {
	XMLName         xml.Name `xml:"IMPACT"`
	Confidentiality int      `xml:"CONFIDENTIALITY,omitempty"`
	Integrity       int      `xml:"INTEGRITY,omitempty"`
	Availability    int      `xml:"AVAILABILITY,omitempty"`
}

// QDiscovery is a member of QVulnerability and must be exported in order to be marshaled
type QDiscovery struct {
	XMLName   xml.Name `xml:"DISCOVERY"`
	Remote    bool     `xml:"REMOTE"`
	AuthTypes []string `xml:"AUTH_TYPE_LIST>AUTH_TYPE,omitempty"`
	Info      string   `xml:"ADDITIONAL_INFO,omitempty"`
}
