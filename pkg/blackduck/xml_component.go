package blackduck

import "time"

type ComponentResponse struct {
	TotalCount     int              `json:"totalCount"`
	Items          []ComponentItems `json:"items"`
	AppliedFilters []AppliedFilters `json:"appliedFilters"`
	Meta           Meta             `json:"_meta"`
}

type ComponentLicenses struct {
	LicenseDisplay string        `json:"licenseDisplay"`
	License        string        `json:"license"`
	SpdxID         string        `json:"spdxId"`
	Licenses       []interface{} `json:"licenses"`
}

type Counts struct {
	CountType string `json:"countType"`
	Count     int    `json:"count"`
}
type LicenseRiskProfile struct {
	Counts []Counts `json:"counts"`
}
type SecurityRiskProfile struct {
	Counts []Counts `json:"counts"`
}
type VersionRiskProfile struct {
	Counts []Counts `json:"counts"`
}
type ActivityRiskProfile struct {
	Counts []Counts `json:"counts"`
}
type OperationalRiskProfile struct {
	Counts []Counts `json:"counts"`
}
type ActivityData1 struct {
	LastCommitDate time.Time `json:"lastCommitDate"`
	NewerReleases  int       `json:"newerReleases"`
}
type ActivePolicies struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}
type ActivityData2 struct {
	ContributorCount12Month int       `json:"contributorCount12Month"`
	CommitCount12Month      int       `json:"commitCount12Month"`
	LastCommitDate          time.Time `json:"lastCommitDate"`
	Trending                string    `json:"trending"`
	NewerReleases           int       `json:"newerReleases"`
}
type ComponentItems struct {
	ComponentName          string                 `json:"componentName"`
	ComponentVersionName   string                 `json:"componentVersionName"`
	Component              string                 `json:"component"`
	ComponentVersion       string                 `json:"componentVersion"`
	TotalFileMatchCount    int                    `json:"totalFileMatchCount"`
	Licenses               []ComponentLicenses    `json:"licenses"`
	Origins                []interface{}          `json:"origins"`
	Usages                 []string               `json:"usages"`
	Ignored                bool                   `json:"ignored"`
	MatchTypes             []string               `json:"matchTypes"`
	ReleasedOn             time.Time              `json:"releasedOn"`
	LicenseRiskProfile     LicenseRiskProfile     `json:"licenseRiskProfile"`
	SecurityRiskProfile    SecurityRiskProfile    `json:"securityRiskProfile"`
	VersionRiskProfile     VersionRiskProfile     `json:"versionRiskProfile"`
	ActivityRiskProfile    ActivityRiskProfile    `json:"activityRiskProfile"`
	OperationalRiskProfile OperationalRiskProfile `json:"operationalRiskProfile"`
	ReviewStatus           string                 `json:"reviewStatus"`
	ApprovalStatus         string                 `json:"approvalStatus"`
	PolicyStatus           string                 `json:"policyStatus"`
	ManuallyAdjusted       bool                   `json:"manuallyAdjusted"`
	ComponentModified      bool                   `json:"componentModified"`
	InAttributionReport    bool                   `json:"inAttributionReport"`
	ComponentType          string                 `json:"componentType"`
	Meta                   Meta                   `json:"_meta"`
	CommentCount           int                    `json:"commentCount"`
	ActivePolicies         []ActivePolicies       `json:"activePolicies"`
	CustomFields           []interface{}          `json:"customFields"`

	ActivityData1 interface{} `json:"activityData,omitempty"`
	ActivityData2 interface{} `json:"activityData,omitempty"`
}
type Selected struct {
	Key   string `json:"key"`
	Label string `json:"label"`
}
type AppliedFilters struct {
	Name     string     `json:"name"`
	Label    string     `json:"label"`
	Selected []Selected `json:"selected"`
}
