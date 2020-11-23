package nexpose

const (
	apiEndpoint = "%s:%v/api/3/%s" // Takes url, port, path
	//apiMediaType = "application/json"

	// Assets
	// GET: /api/3/assets
	// Returns all assets for which you have access.
	apiGetAssets = "assets"

	// Asset
	// GET: api/3/assets/{id}
	// Returns the specified asset.
	apiGetAsset = apiGetAssets + "/%v" // Takes id of asset

	// Site Assets
	// GET: /api/3/sites/{id}/assets
	// Retrieves a paged resource of assets linked with the specified site.
	apiGetSiteAssets = "sites/%v/" + apiGetAssets

	// Asset Vulnerabilities
	// GET: /api/3/assets/{id}/vulnerabilities
	// Retrieves all vulnerability findings on an asset. A finding may be invulnerable if all instances have exceptions applied.
	apiGetAssetVulnerabilities = apiGetAssets + "/%v/" + apiGetVulnerabilities // takes asset id

	// Asset Vulnerability
	// GET: /api/3/assets/{id}/vulnerabilities/{vulnerabilityId}
	// Retrieves the details for a vulnerability finding on an asset.
	apiGetAssetVulnerabilityDetail = apiGetAssetVulnerabilities + "/%v" // takes asset id then vulnerability id

	// Scans
	// GET: /api/3/scans
	// Returns all scans.
	apiGetScans = "scans"

	// Scan
	// GET: /api/3/scans/{id}
	// Returns the specified scan.
	apiGetScan = "scans/%v" // Takes id of scan

	// Create Site Scan
	// POST: /api/3/sites/{id}/scans
	// Starts a scan for the specified site.
	apiPostSiteScan = "sites/%v/scans"

	// Get Sites
	// GET: /api/3/sites
	// Retrieves all sites in the Nexpose instance
	apiGetSites = "sites"

	// Site Scan Engine
	// GET: /api/3/sites/{id}/scan_engine
	// Retrieves the resource of the scan engine assigned to the site.
	apiGetSiteScanEngine = "sites/%v/scan_engine" // takes id of site

	// Get Scan Templates
	// GET: /api/3/scan_templates
	// Returns all scan templates.
	apiGetScanTemplates = "scan_templates"

	// Get Scan Template
	// GET: /api/3/scan_templates/{id}
	// Returns a scan template.
	apiGetScanTemplate = apiGetScanTemplates + "/%s" // Takes id of scan template

	// Create Scan Templates
	// POST: /api/3/scan_templates
	// Creates a new scan template.
	apiPostScanTemplates = apiGetScanTemplates

	// Delete Scan Template
	// DELETE: /api/3/scan_templates/{id}
	// Deletes a scan template.
	apiDeleteScanTemplate = apiGetScanTemplates + "/%s" // Takes id of scan template

	// Vulnerabilities
	// GET: /api/3/vulnerabilities
	// Returns all vulnerabilities that can be assessed during a scan.
	apiGetVulnerabilities = "vulnerabilities"

	// Vulnerability
	// GET: /api/3/vulnerabilities/{id}
	// Returns the details for a vulnerability.
	apiGetVulnerability = apiGetVulnerabilities + "/%v" // Takes id of vulnerability

	// Vulnerability References
	// GET: /api/3/vulnerabilities/{id}/references
	// Returns the external references that may be associated to a vulnerability.
	apiGetVulnerabilityReferences = apiGetVulnerabilities + "/%v/references" // Takes id of vulnerability

	// Get All Solutions
	// GET: /api/3/solutions
	// Returns the details for all solutions.
	apiGetSolutions = "solutions"

	// Get Solution
	// GET: /api/3/solutions/{id}
	// Returns the details for a solution that can remediate one or more vulnerabilities.
	apiGetSolution = apiGetSolutions + "/%s"

	// Vulnerability Solutions
	// GET: /api/3/vulnerabilities/{id}/solutions
	// Returns all solutions (across all platforms) that may be used to remediate this vulnerability.
	apiGetVulnerabilitySolutions = apiGetVulnerabilities + "/%v/solutions"

	// Vulnerability Checks
	// GET: /api/3/vulnerabilities/{id}/checks
	// Returns the vulnerability checks that assess for a specific vulnerability during a scan.
	apiGetVulnerabilityChecks = apiGetVulnerabilities + "/%v/checks" // Takes id of vulnerability

	// Vulnerability Check
	// GET: /api/3/vulnerability_checks/{id}
	// Returns the vulnerability check.
	apiGetVulnerabilityCheck = "vulnerability_checks/%s" // Takes id of vulnerability check

	// pageSize is used for paging and indicates the max return from an endpoint
	pageSize = 500
)
