package blackduck

import "time"

type PolicyStatusResp struct {
	ApprovalStatus           string                     `json:"approvalStatus"`
	PolicyRuleViolationViews []PolicyRuleViolationViews `json:"policyRuleViolationViews"`
	Meta                     Meta                       `json:"_meta"`
}

type Data struct {
	Data interface{} `json:"data"`
}

type Parameters struct {
	Values interface{} `json:"values"`
	Data   []Data      `json:"data"`
}

type Expressions struct {
	Name        string     `json:"name"`
	Operation   string     `json:"operation"`
	Parameters  Parameters `json:"parameters"`
	DisplayName string     `json:"displayName"`
}

type PolicyRuleView struct {
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	Enabled       bool       `json:"enabled"`
	Overridable   bool       `json:"overridable"`
	Severity      string     `json:"severity"`
	Expression    Expression `json:"expression"`
	Category      string     `json:"category"`
	CreatedAt     time.Time  `json:"createdAt"`
	CreatedBy     string     `json:"createdBy"`
	CreatedByUser string     `json:"createdByUser"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	UpdatedBy     string     `json:"updatedBy"`
	UpdatedByUser string     `json:"updatedByUser"`
	Meta          Meta       `json:"_meta"`
}

type PolicyRuleViolationViews struct {
	PolicyRuleView    PolicyRuleView `json:"policyRuleView"`
	ApprovalStatus    string         `json:"approvalStatus"`
	OverriddenBy      []interface{}  `json:"overriddenBy"`
	UpdatedBy         []interface{}  `json:"updatedBy"`
	ApprovalStatusURL string         `json:"approvalStatusUrl,omitempty"`
}
