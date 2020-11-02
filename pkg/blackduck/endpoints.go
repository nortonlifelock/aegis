package blackduck

const (
	GetProject                = "/api/projects/{PROJECT_ID}"
	GetProjectVersions        = "/api/projects/{PROJECT_ID}/versions"
	GetComponents             = "/api/projects/{PROJECT_ID}/versions/{VERSION_ID}/components"
	GetVulnerabilityBOM       = "/api/projects/{PROJECT_ID}/versions/{VERSION_ID}/vulnerability-bom"
	GetComponentVulnerability = "/api/projects/{PROJECT_ID}/versions/{VERSION_ID}/components/{COMPONENT_ID}/versions/{COMPONENT_VERSION}/vulnerabilities"
	GetPolicyStatus           = "/api/projects/{PROJECT_ID}/versions/{VERSION_ID}/components/{COMPONENT_ID}/versions/{COMPONENT_VERSION}/policy-rules"
)
