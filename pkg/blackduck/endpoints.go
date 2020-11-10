package blackduck

const (
	getProject           = "/api/projects/{PROJECT_ID}"
	getProjectVersions   = "/api/projects/{PROJECT_ID}/versions"
	getComponents        = "/api/projects/{PROJECT_ID}/versions/{VERSION_ID}/components"
	getVulnerabilityBOM  = "/api/projects/{PROJECT_ID}/versions/{VERSION_ID}/vulnerability-bom"
	getPolicyStatus      = "/api/projects/{PROJECT_ID}/versions/{VERSION_ID}/components/{COMPONENT_ID}/versions/{COMPONENT_VERSION}/policy-rules"
	getVulnerabilityInfo = "/api/vulnerabilities/{VULNERABILITY_ID}"

	getComponentVulnerabilityNoProject = "/api/components/{COMPONENT_ID}/versions/{COMPONENT_VERSION}/vulnerabilities"
	getComponentVulnerability          = "/api/projects/{PROJECT_ID}/versions/{VERSION_ID}/components/{COMPONENT_ID}/versions/{COMPONENT_VERSION}/vulnerabilities"
)
