package dome9

import (
	"fmt"
	"strings"
	"time"
)

// Description returns the Description parameter of the ticket
func (finding *Finding) String() (param string) {
	return finding.getDescription()
}

// ID returns the VulnerabilityID parameter of the ticket
func (finding *Finding) ID() (param string) {
	return finding.RuleHash
}

// DeviceID returns the DeviceID parameter of the ticket
func (finding *Finding) DeviceID() (param string) {
	return finding.EntityExternalID
}

// AccountID returns the AccountID parameter of the ticket
func (finding *Finding) AccountID() (param string) {
	return finding.CloudAccountID
}

// Priority returns the Priority parameter of the ticket
func (finding *Finding) Priority() (param string) {
	return finding.Severity
}

// ScanID returns the ScanID parameter of the ticket
func (finding *Finding) ScanID() (param int) {
	return finding.assessmentID
}

func (finding *Finding) LastFound() (param time.Time) {
	return time.Time{}
}

// Summary returns the Summary parameter of the ticket
func (finding *Finding) Summary() (param string) {
	return fmt.Sprintf("Aegis (%s) - %s", finding.EntityName, finding.VulnerabilityTitle())
}

// BundleID returns the bundle id of the finding
func (finding *Finding) BundleID() string {
	return finding.bundleID
}

// VulnerabilityTitle returns the VulnerabilityTitle parameter of the ticket
func (finding *Finding) VulnerabilityTitle() (param string) {
	return finding.RuleName
}

func (finding *Finding) getDescription() string {
	descriptionTemp := `
		*Scan Data:*

		Rule Name: %rulename
		Rule Hash: %rulehash
		Bundle ID: %bundle
	
		Entity Name: %entityname
		Entity Type: %entitytype
	
		Solution: %solution

		Cloud Account:               %cloudAccountName
		Dome9 Cloud Account Id:      %dome9ID
		External Cloud Account Id:   %externalCloudID`

	//descriptionTemp = strings.Replace(descriptionTemp, "%accountType", finding.CloudAccountType, 1)
	descriptionTemp = strings.Replace(descriptionTemp, "%dome9ID", finding.CloudAccountID, 1)
	descriptionTemp = strings.Replace(descriptionTemp, "%externalCloudID", finding.externalCloudAccountID, 1)
	descriptionTemp = strings.Replace(descriptionTemp, "%cloudAccountName", finding.externalCloudAccountName, 1)
	descriptionTemp = strings.Replace(descriptionTemp, "%rulename", finding.RuleName, 1)
	descriptionTemp = strings.Replace(descriptionTemp, "%rulehash", finding.RuleHash, 1)
	descriptionTemp = strings.Replace(descriptionTemp, "%entityname", finding.EntityName, 1)
	descriptionTemp = strings.Replace(descriptionTemp, "%entitytype", finding.EntityType, 1)
	descriptionTemp = strings.Replace(descriptionTemp, "%solution", finding.Remediation, 1)
	descriptionTemp = strings.Replace(descriptionTemp, "%bundle", finding.bundleID, 1)

	return descriptionTemp
}
