package endpoints

import "github.com/nortonlifelock/aegis/pkg/domain"

// ADConfig holds the fields that are required to authenticate against AD as well as grab
// the groups that a user is a member of
type ADConfig struct {
	Servers             []string `json:"ad_servers"`
	ADLdapTLSPort       int      `json:"ad_ldap_tls_port"`
	ADBaseDN            string   `json:"ad_base_dn"`
	ADSkipTLSVerify     bool     `json:"ad_skip_verify"`
	ADMemberOfAttribute string   `json:"ad_member_of_attribute"`
	ADSearchString      string   `json:"ad_search_string"`
}

// OrgConfigWrapper ties together an organization along with it's AD configuration
type OrgConfigWrapper struct {
	Org domain.Organization
	Con *ADConfig
}
