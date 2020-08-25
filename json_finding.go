package dome9

import "time"

// Finding holds information pertaining to rule violations
type Finding struct {
	//FindingID string  `json:"id"`
	//FindingKey string `json:"findingKey"`

	assessmentID int // not populated by endpoint

	EntityType       string `json:"entityType"`
	EntityName       string `json:"entityName"`
	EntityExternalID string `json:"entityExternalId"`

	FindingDescription string `json:"description"`
	Remediation        string `json:"remediation"`

	RuleName  string `json:"ruleName"`
	RuleLogic string `json:"ruleLogic"`
	RuleHash  string // Not populated by endpoint

	Severity string `json:"severity"`

	//CreatedTime  time.Time `json:"createdTime"`
	//UpdatedTime  time.Time `json:"updatedTime"`
	//LastSeenTime time.Time `json:"lastSeenTime"`

	//CloudAccountType string `json:"cloudAccountType"`
	CloudAccountID           string `json:"cloudAccountId"`
	externalCloudAccountID   string
	externalCloudAccountName string

	// bundleID is not returned by the API, but is populated by the api wrapper
	bundleID int

	// dueDate is not returned by the API, and must be set manually
	dueDate time.Time
}

// Bundle holds the rules associated to a particular binding, which are then used to get findings that violate the rule
type Bundle struct {
	Rules     []Rule `json:"rules"`
	AccountID int    `json:"accountId"`
	//CreatedTime      time.Time `json:"createdTime"`
	//UpdatedTime      time.Time `json:"updatedTime"`
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	IsTemplate       bool   `json:"isTemplate"`
	HideInCompliance bool   `json:"hideInCompliance"`
	MinFeatureTier   string `json:"minFeatureTier"`
	Section          int    `json:"section"`
	TooltipText      string `json:"tooltipText"`
	ShowBundle       bool   `json:"showBundle"`
	SystemBundle     bool   `json:"systemBundle"`
	CloudVendor      string `json:"cloudVendor"`
	Version          int    `json:"version"`
	Language         string `json:"language"`
	RulesCount       int    `json:"rulesCount"`
}

// Rule holds information pertaining to compliance in a CIS engine
type Rule struct {
	Name          string `json:"name"`
	Severity      string `json:"severity"`
	Logic         string `json:"logic"`
	Description   string `json:"description"`
	Remediation   string `json:"remediation"`
	ComplianceTag string `json:"complianceTag"`
	Domain        string `json:"domain"`
	Priority      string `json:"priority"`
	ControlTitle  string `json:"controlTitle"`
	RuleID        string `json:"ruleId"`
	LogicHash     string `json:"logicHash"`
	IsDefault     bool   `json:"isDefault"`
}

// CloudAccount ... TODO:
type CloudAccount struct {
	ID                    string `json:"id"`
	Vendor                string `json:"vendor"`
	Name                  string `json:"name"`
	ExternalAccountNumber string `json:"externalAccountNumber"`
	Error                 string `json:"error"`
	//CreationDate          time.Time `json:"creationDate"`
	//Credentials           struct {
	//	Apikey     string `json:"apikey"`
	//	Arn        string `json:"arn"`
	//	Secret     string `json:"secret"`
	//	IamUser    string `json:"iamUser"`
	//	Type       string `json:"type"`
	//	IsReadOnly bool   `json:"isReadOnly"`
	//} `json:"credentials"`
	//IamSafe struct {
	//	AwsGroupArn         string `json:"awsGroupArn"`
	//	AwsPolicyArn        string `json:"awsPolicyArn"`
	//	Mode                string `json:"mode"`
	//	State               string `json:"state"`
	//	ExcludedIamEntities struct {
	//		RolesArns []string `json:"rolesArns"`
	//		UsersArns []string `json:"usersArns"`
	//	} `json:"excludedIamEntities"`
	//	RestrictedIamEntities struct {
	//		RolesArns []string `json:"rolesArns"`
	//		UsersArns []string `json:"usersArns"`
	//	} `json:"restrictedIamEntities"`
	//} `json:"iamSafe"`
	//NetSec struct {
	//	Regions []struct {
	//		Region           string `json:"region"`
	//		Name             string `json:"name"`
	//		Hidden           bool   `json:"hidden"`
	//		NewGroupBehavior string `json:"newGroupBehavior"`
	//	} `json:"regions"`
	//} `json:"netSec"`
	//Magellan               bool   `json:"magellan"`
	//FullProtection         bool   `json:"fullProtection"`
	//AllowReadOnly          bool   `json:"allowReadOnly"`
	//OrganizationalUnitID   string `json:"organizationalUnitId"`
	//OrganizationalUnitPath string `json:"organizationalUnitPath"`
	//OrganizationalUnitName string `json:"organizationalUnitName"`
	//LambdaScanner          bool   `json:"lambdaScanner"`
}
