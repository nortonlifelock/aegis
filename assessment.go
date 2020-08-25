package dome9

import (
	"fmt"

	"github.com/nortonlifelock/domain"
)

// RescanBundle tests a rule set against a cloud account and returns a slice of findings for each instance of a violation of a
// rule within that cloud account
func (client *Client) RescanBundle(bundleID int, cloudSubscriptionID string) (findings []domain.Finding, err error) {
	findings = make([]domain.Finding, 0)

	var vendor, externalAccountNumber, cloudAccountName string
	if vendor, externalAccountNumber, cloudAccountName, err = client.determineCloudAccountTypeFromSubscriptionID(cloudSubscriptionID); err == nil {
		var assessmentResult *AssessmentResult
		assessmentResult, err = client.RunAssessmentOnBundle(bundleID, cloudSubscriptionID, vendor)
		if err == nil {

			for _, test := range assessmentResult.Tests {

				for _, vulnerableEntity := range test.EntityResults {
					finding := &Finding{
						assessmentID:             assessmentResult.ID,
						bundleID:                 bundleID,
						CloudAccountID:           cloudSubscriptionID,
						externalCloudAccountID:   externalAccountNumber,
						externalCloudAccountName: cloudAccountName,

						EntityType:       vulnerableEntity.Obj.EntityType,
						EntityName:       vulnerableEntity.Obj.Dome9ID,
						EntityExternalID: vulnerableEntity.Obj.ID,

						FindingDescription: test.Rule.Description,
						Remediation:        test.Rule.Remediation,
						RuleName:           test.Rule.Name,
						RuleLogic:          test.Rule.Logic,
						RuleHash:           test.Rule.LogicHash,
						Severity:           test.Rule.Severity,
					}

					findings = append(findings, finding)
				}
			}
		} else {
			err = fmt.Errorf("error while running assessment on bundle - %s", err.Error())
		}
	} else {
		err = fmt.Errorf("error while determining cloud account vendor for [%v] - %v", cloudSubscriptionID, err.Error())
	}

	return findings, err
}
