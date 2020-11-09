package blackduck

import (
	"context"
	"fmt"
	"github.com/nortonlifelock/domain"
	"strings"
	"time"
)

const (
	policyInViolation = "IN_VIOLATION"
)

val := "\"Method of Discovery\" = Qualys AND project in (VRR, \"VRR SurfEasy\") AND status not in (Closed-Decommission, Closed-Remediated, CLOSED-NA, Closed-Error, \"Closed-Moved to IDA\") AND LastFound <= 2020-10-15 AND GroupID ~ \"tag*\""

func (cli *BlackDuckClient) GetProjectVulnerabilities(ctx context.Context, projectID string) (findings []domain.CodeFinding, err error) {
	findings = make([]domain.CodeFinding, 0)

	projectResponse, err := cli.GetProject(projectID)
	if err == nil {
		projectVersionsResponse, err := cli.GetProjectVersions(projectID)
		// TODO only use most recent project version
		if err == nil {
			for _, version := range projectVersionsResponse.Items {
				linkContainingVersionID := version.Meta.Href
				var lookingFor = "/versions/"
				projectVersionID := linkContainingVersionID[strings.Index(linkContainingVersionID, lookingFor)+len(lookingFor):]

				findingsForVersion, err := cli.getVulnerabilityFindings(ctx, projectID, projectVersionID, projectResponse, version)

				select {
				case <-ctx.Done():
					return nil, fmt.Errorf("context closed")
				default:
				}

				if err == nil {
					for _, ffv := range findingsForVersion {
						findings = append(findings, ffv)
					}
				} else {
					break
				}
			}
		} else {
			err = fmt.Errorf("error while getting project versions for [%s] - %s", projectID, err.Error())
		}
	} else {
		err = fmt.Errorf("error while getting project information for [%s] - %s", projectID, err.Error())
	}

	return findings, err
}

func (cli *BlackDuckClient) getVulnerabilityFindings(ctx context.Context, projectID string, projectVersionID string, projectResponse *ProjectResponse, version ProjectItem) (findings []*BlackDuckFinding, err error) {
	findings = make([]*BlackDuckFinding, 0)

	componentInfoResp, err := cli.GetComponentInformation(projectID, projectVersionID)
	if err == nil {
		for index := range componentInfoResp.Items {
			select {
			case <-ctx.Done():
				return nil, fmt.Errorf("context closed")
			default:
			}

			component := componentInfoResp.Items[index]
			lowerRange := "/components/"
			upperRange := "/versions/"

			componentID := component.ComponentVersion[strings.Index(component.ComponentVersion, lowerRange)+len(lowerRange) : strings.Index(component.ComponentVersion, upperRange)]
			componentVersionID := component.ComponentVersion[strings.Index(component.ComponentVersion, upperRange)+len(upperRange):]

			componentVulnerabilityResponse, err := cli.GetComponentVulnerabilities(projectID, projectVersionID, componentID, componentVersionID)
			if err == nil {
				policyResp, err := cli.GetPolicyStatus(projectID, projectVersionID, componentID, componentVersionID)
				if err == nil {
					for index := range componentVulnerabilityResponse.Items {
						vuln := componentVulnerabilityResponse.Items[index]
						finding := &BlackDuckFinding{
							ProjectUUID:   projectID,
							ProjectInfo:   projectResponse,
							ProjectItem:   &version,
							Component:     &component,
							ComponentVuln: &vuln,
							PolicyRules:   policyResp,
						}

						findings = append(findings, finding)
					}
				} else {
					err = fmt.Errorf("error while getting policies for [%s/%s/%s/%s] - %s", projectID, projectVersionID, componentID, componentVersionID, err.Error())
					break
				}
			} else {
				err = fmt.Errorf("error while getting component vulnerabilities for [%s/%s/%s/%s] - %s", projectID, projectVersionID, componentID, componentVersionID, err.Error())
				break
			}
		}
	} else {
		err = fmt.Errorf("error while getting component information for [%s/%s] - %s", projectID, projectVersionID, err.Error())
	}

	return findings, err
}

type BlackDuckFinding struct {
	ProjectUUID   string
	ProjectInfo   *ProjectResponse
	ProjectItem   *ProjectItem
	Component     *ComponentItems
	ComponentVuln *ComponentVulnerabilityItem
	PolicyRules   *PolicyRulesResp
}

func (finding *BlackDuckFinding) ProjectID() string {
	return finding.ProjectUUID
}
func (finding *BlackDuckFinding) ProjectName() string {
	return finding.ProjectInfo.Name
}
func (finding *BlackDuckFinding) ProjectVersion() string {
	return finding.ProjectItem.VersionName
}
func (finding *BlackDuckFinding) ProjectOwner() string {
	return finding.ProjectItem.CreatedBy
}
func (finding *BlackDuckFinding) ComponentName() string {
	return finding.Component.ComponentName
}
func (finding *BlackDuckFinding) ComponentVersion() string {
	return finding.Component.ComponentVersionName
}
func (finding *BlackDuckFinding) ViolatedPolicyName() string {
	var rules = make([]string, 0)
	for _, rule := range finding.PolicyRules.Items {
		if rule.PolicyApprovalStatus == policyInViolation {
			rules = append(rules, rule.Name)
		}
	}

	val := strings.Join(rules, ",")
	return val
}
func (finding *BlackDuckFinding) ViolatedPolicySeverity() string {
	_, severity := finding.getHighestSeverityAndPolicy()
	return severity
}

func (finding *BlackDuckFinding) getHighestSeverityAndPolicy() (policy, severity string) {
	severityNameToSeverityLevel := map[string]int{
		"TRIVIAL":  0,
		"MINOR":    1,
		"MAJOR":    2,
		"CRITICAL": 3,
		"BLOCKER":  4,
	}

	var highestSeverityLevel = -1
	for _, rule := range finding.PolicyRules.Items {
		if rule.PolicyApprovalStatus == policyInViolation {
			if severityNameToSeverityLevel[rule.Severity] > highestSeverityLevel {
				highestSeverityLevel = severityNameToSeverityLevel[rule.Severity]
				severity = rule.Severity
				policy = rule.Name
			}
		}
	}
	return policy, severity
}
func (finding *BlackDuckFinding) CVSS() float32 {
	return float32(finding.ComponentVuln.Cvss2.OverallScore)
}

func (finding *BlackDuckFinding) Description() string {
	projectOwner := finding.ProjectItem.CreatedBy
	vulnSummary := finding.ComponentVuln.Summary
	vals := fmt.Sprintf("Project owner: %s\n\n%s", projectOwner, vulnSummary)
	return vals
}

func (finding *BlackDuckFinding) Summary() string {
	var header string
	if len(finding.ViolatedPolicyName()) > 0 {
		highestSeverityPolicy, _ := finding.getHighestSeverityAndPolicy()
		header = fmt.Sprintf("(%s) ", highestSeverityPolicy)
	}

	val := fmt.Sprintf("%s%s %s, %s %s", header, finding.ProjectName(), finding.ProjectVersion(), finding.ComponentName(), finding.ComponentVersion())
	return val
}
func (finding *BlackDuckFinding) Updated() time.Time {
	return finding.ComponentVuln.UpdatedAt
}
func (finding *BlackDuckFinding) VulnerabilityID() string {
	return finding.ComponentVuln.ID
}
