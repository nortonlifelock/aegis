package dome9

const (
	getAccountBundles = "/v2/CompliancePolicy/view"
	getBundleByID     = "/v2/CompliancePolicy/{bundleId}"
	getBundleFindings = "/v2/Compliance/Finding/bundle/{bundleId}"

	getCloudAccounts      = "/v2/CloudAccounts"
	getCloudAccountByID   = "/v2/CloudAccounts/{id}"
	getAzureCloudAccounts = "/v2/AzureCloudAccount"
	getAzureCloudAccount  = "/v2/AzureCloudAccount/{id}"

	// POST
	runAssessmentOnBundle = "/v2/assessment/bundleV2"
)
