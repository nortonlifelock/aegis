package aqua

type Repositories struct {
	Count    int                `json:"count"`
	Page     int                `json:"page"`
	Pagesize int                `json:"pagesize"`
	Result   []RepositoryResult `json:"result"`
	Query    Query              `json:"query"`
	MoreData int                `json:"more_data"`
}
type RepositoryResult struct {
	Name                   string `json:"name"`
	Registry               string `json:"registry"`
	Author                 string `json:"author"`
	Policy                 string `json:"policy"`
	DynamicProfiling       bool   `json:"dynamic_profiling"`
	NumImages              int    `json:"num_images"`
	NumDisallowed          int    `json:"num_disallowed"`
	NumFailed              int    `json:"num_failed"`
	CritVulns              int    `json:"crit_vulns"`
	HighVulns              int    `json:"high_vulns"`
	MedVulns               int    `json:"med_vulns"`
	LowVulns               int    `json:"low_vulns"`
	NegVulns               int    `json:"neg_vulns"`
	SensitiveData          int    `json:"sensitive_data"`
	Malware                int    `json:"malware"`
	TrustedBaseCount       int    `json:"trusted_base_count"`
	WhitelistedImagesCount int    `json:"whitelisted_images_count"`
	IsDefaultPolicy        bool   `json:"is_default_policy"`
}

type Variables struct {
	Attribute string `json:"attribute"`
	Value     string `json:"value"`
}

type Scopes struct {
	Expression string      `json:"expression"`
	Variables  []Variables `json:"variables"`
}

type Query struct {
	Time    string   `json:"time"`
	Orderby string   `json:"orderby"`
	Scopes  []Scopes `json:"Scopes"`
}
