package dome9

import "time"

type createAssessmentRequestBody struct {
	BundleID          int    `json:"id"`
	BundleName        string `json:"name"`
	BundleDescription string `json:"description"`

	// indicates request target is a CFT (template file)
	IsCft bool `json:"isCft"`

	// The Dome9 account id
	Dome9CloudAccountID string `json:"dome9CloudAccountId"`

	// account id on cloud provider (AWS, Azure, GCP)
	ExternalCloudAccountID string `json:"externalCloudAccountId"`

	// account id on cloud provider (AWS, Azure, GCP)
	CloudAccountID string `json:"cloudAccountId"`

	// cloud provider (AWS/Azure/GCP)
	CloudNetwork string `json:"cloudNetwork"`

	// the cloud provider (AWS/Azure/GCP)
	CloudAccountType string `json:"cloudAccountType"`

	// the assessment id (returned in the response)
	RequestID string `json:"requestId"`
}

// AssessmentResult contains the information about a specific assessment
type AssessmentResult struct {
	Tests            []Test    `json:"tests"`
	CreatedTime      time.Time `json:"createdTime"`
	ID               int       `json:"id"`
	AssessmentPassed bool      `json:"assessmentPassed"`
	HasErrors        bool      `json:"hasErrors"`
}

// ExclusionResult is the ... TODO:
type ExclusionResult struct {
	ValidationStatus string     `json:"validationStatus"`
	IsRelevant       bool       `json:"isRelevant"`
	IsValid          bool       `json:"isValid"`
	IsExcluded       bool       `json:"isExcluded"`
	ExclusionID      string     `json:"exclusionId"`
	Error            string     `json:"error"`
	Obj              TestObject `json:"testObj"`
}

// TestObject ... TODO:
type TestObject struct {
	ID          string `json:"id"`
	Dome9ID     string `json:"dome9Id"`
	EntityType  string `json:"entityType"`
	EntityIndex int    `json:"entityIndex"`
}

// Test ... TODO:
type Test struct {
	Error             string `json:"error"`
	TestedCount       int    `json:"testedCount"`
	RelevantCount     int    `json:"relevantCount"`
	NonComplyingCount int    `json:"nonComplyingCount"`
	ExclusionStats    struct {
		TestedCount       int `json:"testedCount"`
		RelevantCount     int `json:"relevantCount"`
		NonComplyingCount int `json:"nonComplyingCount"`
	} `json:"exclusionStats"`
	EntityResults []ExclusionResult `json:"entityResults"`
	Rule          Rule              `json:"rule"`
	TestPassed    bool              `json:"testPassed"`
}
