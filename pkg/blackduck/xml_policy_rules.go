package blackduck

import "time"

type PolicyRulesResp struct {
	TotalCount     int           `json:"totalCount"`
	Items          []PolicyItem  `json:"items"`
	AppliedFilters []interface{} `json:"appliedFilters"`
	Meta           Meta          `json:"_meta"`
}

type Expression struct {
	Operator    string        `json:"operator"`
	Expressions []Expressions `json:"expressions"`
}

type PolicyItem struct {
	Name                 string     `json:"name"`
	Description          string     `json:"description"`
	Enabled              bool       `json:"enabled"`
	Overridable          bool       `json:"overridable"`
	Severity             string     `json:"severity"`
	Expression           Expression `json:"expression"`
	CreatedAt            time.Time  `json:"createdAt"`
	CreatedBy            string     `json:"createdBy"`
	CreatedByUser        string     `json:"createdByUser"`
	UpdatedAt            time.Time  `json:"updatedAt"`
	UpdatedBy            string     `json:"updatedBy"`
	UpdatedByUser        string     `json:"updatedByUser"`
	PolicyApprovalStatus string     `json:"policyApprovalStatus"`
	Meta                 Meta       `json:"_meta"`
}
