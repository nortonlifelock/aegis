package blackduck

import "time"

type ProjectResponse struct {
	Name                     string    `json:"name"`
	ProjectLevelAdjustments  bool      `json:"projectLevelAdjustments"`
	CloneCategories          []string  `json:"cloneCategories"`
	CustomSignatureEnabled   bool      `json:"customSignatureEnabled"`
	CustomSignatureDepth     int       `json:"customSignatureDepth"`
	DeepLicenseDataEnabled   bool      `json:"deepLicenseDataEnabled"`
	SnippetAdjustmentApplied bool      `json:"snippetAdjustmentApplied"`
	CreatedAt                time.Time `json:"createdAt"`
	CreatedBy                string    `json:"createdBy"`
	CreatedByUser            string    `json:"createdByUser"`
	UpdatedAt                time.Time `json:"updatedAt"`
	UpdatedBy                string    `json:"updatedBy"`
	UpdatedByUser            string    `json:"updatedByUser"`
	Source                   string    `json:"source"`
	Meta                     Meta      `json:"_meta"`
}

type ProjectVersionResponse struct {
	TotalCount     int           `json:"totalCount"`
	Items          []ProjectItem `json:"items"`
	AppliedFilters []interface{} `json:"appliedFilters"`
	Meta           Meta          `json:"_meta"`
}

type Licenses struct {
	License              string               `json:"license"`
	Licenses             []interface{}        `json:"licenses"`
	Name                 string               `json:"name"`
	Ownership            string               `json:"ownership"`
	LicenseDisplay       string               `json:"licenseDisplay"`
	LicenseFamilySummary LicenseFamilySummary `json:"licenseFamilySummary"`
}

type License struct {
	Type           string     `json:"type"`
	Licenses       []Licenses `json:"licenses"`
	LicenseDisplay string     `json:"licenseDisplay"`
}

type ProjectItem struct {
	VersionName          string    `json:"versionName"`
	Phase                string    `json:"phase"`
	Distribution         string    `json:"distribution"`
	License              License   `json:"license"`
	CreatedAt            time.Time `json:"createdAt"`
	CreatedBy            string    `json:"createdBy"`
	CreatedByUser        string    `json:"createdByUser"`
	SettingUpdatedAt     time.Time `json:"settingUpdatedAt"`
	SettingUpdatedBy     string    `json:"settingUpdatedBy"`
	SettingUpdatedByUser string    `json:"settingUpdatedByUser"`
	LastBOMUpdate        time.Time `json:"lastBomUpdateDate"`
	LastScanDate         time.Time `json:"lastScanDate"`
	Source               string    `json:"source"`
	Meta                 Meta      `json:"_meta"`
}

type LicenseFamilySummary struct {
	Name string `json:"name"`
	Href string `json:"href"`
}
