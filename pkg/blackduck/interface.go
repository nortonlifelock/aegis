package blackduck

import (
	"fmt"
	"strings"
	"time"
)

const (
	policyInViolation = "IN_VIOLATION"
)

func (cli *BlackDuckClient) GetVulnerabilitiesForProject(projectID string) (findings []*BlackDuckFinding) {
	findings = make([]*BlackDuckFinding, 0)

	projectResponse, err := cli.GetProject(projectID)
	check(err)

	projectVersionsResponse, err := cli.GetProjectVersions(projectID)
	check(err)

	for _, version := range projectVersionsResponse.Items {
		linkContainingVersionID := version.Meta.Href
		var lookingFor = "/versions/"
		projectVersionID := linkContainingVersionID[strings.Index(linkContainingVersionID, lookingFor)+len(lookingFor):]

		findingsForVersion := cli.getVulnerabilityFindings(projectID, projectVersionID, projectResponse, version)

		findings = append(findings, findingsForVersion...)
	}

	return findings
}

func (cli *BlackDuckClient) getVulnerabilityFindings(projectID string, projectVersionID string, projectResponse *ProjectResponse, version ProjectItem) (findings []*BlackDuckFinding) {
	findings = make([]*BlackDuckFinding, 0)
	//vulnerabilityBOMResponse, err := cli.GetVulnerabilityBOM(projectID, projectVersionID)
	//check(err)

	componentInfoResp, err := cli.GetComponentInformation(projectID, projectVersionID)
	check(err)

	for index := range componentInfoResp.Items {
		component := componentInfoResp.Items[index]
		lowerRange := "/components/"
		upperRange := "/versions/"

		componentID := component.ComponentVersion[strings.Index(component.ComponentVersion, lowerRange)+len(lowerRange) : strings.Index(component.ComponentVersion, upperRange)]
		componentVersionID := component.ComponentVersion[strings.Index(component.ComponentVersion, upperRange)+len(upperRange):]

		componentVulnerabilityResponse, err := cli.GetComponentVulnerabilities(projectID, projectVersionID, componentID, componentVersionID)
		check(err)

		policyResp, err := cli.GetPolicyStatus(projectID, projectVersionID, componentID, componentVersionID)
		check(err)

		for index := range componentVulnerabilityResponse.Items {
			vuln := componentVulnerabilityResponse.Items[index]
			finding := &BlackDuckFinding{
				ProjectInfo:    projectResponse,
				ProjectVersion: &version,
				Component:      &component,
				ComponentVuln:  &vuln,
				PolicyRules:    policyResp,
			}

			projectName := finding.ProjectInfo.Name
			projectVersion := finding.ProjectVersion.VersionName
			projectOwner := finding.ProjectVersion.CreatedBy // reporter?

			componentName := finding.Component.ComponentName
			componentVersion := finding.Component.ComponentVersionName
			componentVulnerabilityStatus := finding.ComponentVuln.RemediationStatus
			vulnCreatedDate, vulnUpdatedDate := finding.ComponentVuln.CreatedAt, finding.ComponentVuln.UpdatedAt

			vulnerabilityID := finding.ComponentVuln.ID
			vulnSummary := finding.ComponentVuln.Summary
			cwes := strings.Join(finding.ComponentVuln.CweIds, ",")
			cvss := finding.ComponentVuln.Cvss2.OverallScore

			policyRule, policySeverity := "", ""           // TODO
			resolvedDate, issuePriority := time.Time{}, "" // issue priority = issue severity? resolved date set in jira?

			_, _, _, _, _, _, _, _, _, _, _, _ = projectOwner, projectName, projectVersion, componentName, componentVersion, policyRule, policySeverity, vulnCreatedDate, vulnUpdatedDate, componentVulnerabilityStatus, resolvedDate, issuePriority
			_, _, _, _ = vulnerabilityID, vulnSummary, cwes, cvss

			findings = append(findings, finding)
		}
	}

	return findings
}

type BlackDuckFinding struct {
	ProjectInfo    *ProjectResponse
	ProjectVersion *ProjectItem
	Component      *ComponentItems
	ComponentVuln  *ComponentVulnerabilityItem
	PolicyRules    *PolicyRulesResp
}

// TODO remove
func check(err error) {
	if err != nil {
		panic(err)
	}
}

func (finding *BlackDuckFinding) HubProjectName() (param *string) {
	return &finding.ProjectInfo.Name
}

func (finding *BlackDuckFinding) HubProjectVersion() (param *string) {
	return &finding.ProjectVersion.VersionName
}

func (finding *BlackDuckFinding) ComponentName() (param *string) {
	return &finding.Component.ComponentName
}

func (finding *BlackDuckFinding) ComponentVersion() (param *string) {
	return &finding.Component.ComponentVersionName
}

func (finding *BlackDuckFinding) PolicyRule() (param *string) {
	var rules = make([]string, 0)
	for _, rule := range finding.PolicyRules.Items {
		if rule.PolicyApprovalStatus == policyInViolation {
			rules = append(rules, rule.Name)
		}
	}

	val := strings.Join(rules, ",")
	return &val
}

func (finding *BlackDuckFinding) PolicySeverity() (param *string) {
	var severities = make([]string, 0)
	for _, rule := range finding.PolicyRules.Items {
		if rule.PolicyApprovalStatus == policyInViolation {
			severities = append(severities, rule.Severity)
		}
	}

	val := strings.Join(severities, ",")
	return &val
}

func (finding *BlackDuckFinding) AlertDate() (param *time.Time) {
	return
}

func (finding *BlackDuckFinding) AssignedTo() (param *string) {
	return
}

func (finding *BlackDuckFinding) AssignmentGroup() (param *string) {
	return
}

func (finding *BlackDuckFinding) Category() (param *string) {
	return
}

func (finding *BlackDuckFinding) CERF() (param string) {
	return
}

func (finding *BlackDuckFinding) ExceptionExpiration() (param time.Time) {
	return
}

func (finding *BlackDuckFinding) CVEReferences() (param *string) {
	return
}

func (finding *BlackDuckFinding) CVSS() (param *float32) {
	val := float32(finding.ComponentVuln.Cvss2.OverallScore)
	param = &val
	return param
}

func (finding *BlackDuckFinding) CloudID() (param string) {
	return
}

func (finding *BlackDuckFinding) Configs() (param string) {
	return
}

func (finding *BlackDuckFinding) CreatedDate() (param *time.Time) {
	return
}

func (finding *BlackDuckFinding) DBCreatedDate() (param time.Time) {
	return
}

func (finding *BlackDuckFinding) DBUpdatedDate() (param *time.Time) {
	return
}

func (finding *BlackDuckFinding) Description() (param *string) {
	projectOwner := finding.ProjectVersion.CreatedBy
	vulnSummary := finding.ComponentVuln.Summary
	vals := fmt.Sprintf("Project owner: %s\n\n%s", projectOwner, vulnSummary)
	return &vals
}

func (finding *BlackDuckFinding) DeviceID() (param string) {
	componentName := finding.Component.ComponentName
	componentVersion := finding.Component.ComponentVersionName
	return fmt.Sprintf("%s %s", componentName, componentVersion)
}

func (finding *BlackDuckFinding) DueDate() (param *time.Time) {
	return
}

func (finding *BlackDuckFinding) ExceptionDate() (param *time.Time) {
	return
}

func (finding *BlackDuckFinding) GroupID() string {
	projectName := finding.ProjectInfo.Name
	projectVersion := finding.ProjectVersion.VersionName
	return fmt.Sprintf("%s [%s]", projectName, projectVersion)
}

func (finding *BlackDuckFinding) HostName() (param *string) {
	return
}

func (finding *BlackDuckFinding) ID() (param int) {
	return
}

func (finding *BlackDuckFinding) IPAddress() (param *string) {
	return
}

func (finding *BlackDuckFinding) Labels() (param *string) {
	return
}

func (finding *BlackDuckFinding) LastChecked() (param *time.Time) {
	param = &finding.ComponentVuln.UpdatedAt
	return param
}

func (finding *BlackDuckFinding) MacAddress() (param *string) {
	return
}

func (finding *BlackDuckFinding) MethodOfDiscovery() (param *string) {
	val := "Black Duck"
	return &val
}

func (finding *BlackDuckFinding) OSDetailed() (param *string) {
	return
}

func (finding *BlackDuckFinding) OperatingSystem() (param *string) {
	return
}

func (finding *BlackDuckFinding) OrgCode() (param *string) {
	val := "LOCK"
	return &val
}

func (finding *BlackDuckFinding) OrganizationID() (param string) {
	return
}

func (finding *BlackDuckFinding) OWASP() (param *string) {
	return
}

func (finding *BlackDuckFinding) Patchable() (param *string) {
	return
}

func (finding *BlackDuckFinding) Priority() (param *string) {
	val := "Low"
	if *finding.CVSS() <= 3 {
		val = "Low"
	} else if *finding.CVSS() <= 6 {
		val = "Medium"
	} else if *finding.CVSS() <= 9 {
		val = "High"
	} else {
		val = "Critical"
	}
	return &val
}

func (finding *BlackDuckFinding) Project() (param *string) {
	val := "RYAN120"
	return &val
}

func (finding *BlackDuckFinding) ReportedBy() (param *string) {
	projectOwner := finding.ProjectVersion.CreatedBy // reporter?
	projectOwner = ""                                // TODO remove
	return &projectOwner
}

func (finding *BlackDuckFinding) ResolutionDate() (param *time.Time) {
	return
}

func (finding *BlackDuckFinding) ResolutionStatus() (param *string) {
	return
}

func (finding *BlackDuckFinding) ScanID() (param int) {
	return
}

func (finding *BlackDuckFinding) ServicePorts() (param *string) {
	return
}

func (finding *BlackDuckFinding) Solution() (param *string) {
	return
}

func (finding *BlackDuckFinding) Status() (param *string) {
	return
}

func (finding *BlackDuckFinding) Summary() (param *string) {
	var header string
	if finding.PolicyRule() != nil && len(*finding.PolicyRule()) > 0 {
		header = fmt.Sprintf("%s / %s", *finding.PolicyRule(), finding.VulnerabilityID())
	} else {
		header = finding.VulnerabilityID()
	}

	val := fmt.Sprintf("(%s) %s, %s", header, finding.GroupID(), finding.DeviceID())
	return &val
}

func (finding *BlackDuckFinding) SystemName() (param *string) {
	return
}

func (finding *BlackDuckFinding) TicketType() (param *string) {
	val := "Request"
	return &val
}

func (finding *BlackDuckFinding) Title() (param string) {
	return
}

func (finding *BlackDuckFinding) UpdatedDate() (param *time.Time) {
	return
}

func (finding *BlackDuckFinding) VendorReferences() (param *string) {
	return
}

func (finding *BlackDuckFinding) VulnerabilityID() (param string) {
	vulnerabilityID := finding.ComponentVuln.ID
	return vulnerabilityID
}

func (finding *BlackDuckFinding) VulnerabilityTitle() (param *string) {
	return
}

func (finding *BlackDuckFinding) ApplicationName() (param *string) {
	return
}
