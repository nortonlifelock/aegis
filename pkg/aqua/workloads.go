package aqua

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (cli *APIClient) GetContainersForNamespace(namespace string) (containers []ContainerInfo, err error) {
	containers = make([]ContainerInfo, 0)
	page := 1

	endpoint := strings.Replace(getContainersFromNameSpace, "$NAMESPACE", namespace, 1)

	endpoint = strings.Replace(endpoint, " ", "%20", -1)

	for {
		var request *http.Request
		if request, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s&page=%d&pagesize=50", cli.baseURL, endpoint, page), nil); err == nil {
			var body []byte
			if body, err = cli.executeRequest(request); err == nil {
				containerResp := &containersForNamespaceResp{}
				if err = json.Unmarshal(body, containerResp); err == nil {
					if len(containerResp.Result) == 0 {
						break
					} else {
						containers = append(containers, containerResp.Result...)
						page++
					}
				} else {
					err = fmt.Errorf("error while parsing vulnerabilities from response - %s", err.Error())
					break
				}
			} else {
				err = fmt.Errorf("error while gathering vulnerabilities - %s", err.Error())
			}
		} else {
			err = fmt.Errorf("error while making request - %s", err.Error())
		}
	}

	return containers, err
}

type containersForNamespaceResp struct {
	Count    int             `json:"count"`
	Page     int             `json:"page"`
	Pagesize int             `json:"pagesize"`
	Result   []ContainerInfo `json:"result"`
	Query    struct {
		IdentifiersOnly bool        `json:"identifiers_only"`
		Status          string      `json:"status"`
		Groupby         string      `json:"groupby"`
		Namespace       string      `json:"namespace"`
		ShowScanStatus  bool        `json:"show_scan_status"`
		ImageNames      interface{} `json:"image_names"`
	} `json:"query"`
}

type ContainerInfo struct {
	ID              string `json:"id"`
	HostID          string `json:"host_id"`
	Name            string `json:"name"`
	HaveSecrets     bool   `json:"have_secrets"`
	HostName        string `json:"host_name"`
	HostLname       string `json:"host_lname"`
	ImageDigest     string `json:"image_digest"`
	ServerDigest    string `json:"server_digest"`
	ImageName       string `json:"image_name"`
	OriginImageName string `json:"origin_image_name"`
	OwnerName       string `json:"owner_name"`
	ImageID         string `json:"image_id"`
	Status          string `json:"status"`
	RiskLevel       string `json:"risk_level"`
	IsEvaluated     bool   `json:"is_evaluated"`
	ValidDigest     bool   `json:"valid_digest"`
	IsProfiling     bool   `json:"is_profiling"`
	IsRegistered    bool   `json:"is_registered"`
	IsDisallowed    bool   `json:"is_disallowed"`
	IsRoot          bool   `json:"is_root"`
	IsPrivileged    bool   `json:"is_privileged"`
	ScanStatus      string `json:"scan_status"`
	CreateTime      int    `json:"create_time"`
	Vulnerabilities struct {
		Total      int `json:"total"`
		Critical   int `json:"critical"`
		High       int `json:"high"`
		Medium     int `json:"medium"`
		Low        int `json:"low"`
		Sensitive  int `json:"sensitive"`
		Malware    int `json:"malware"`
		Negligible int `json:"negligible"`
	} `json:"vulnerabilities"`
	SystemContainer    bool        `json:"system_container"`
	StartTime          int         `json:"start_time"`
	ModifyTime         int         `json:"modify_time"`
	AquaService        string      `json:"aqua_service"`
	Secrets            interface{} `json:"secrets"`
	Risk               int         `json:"risk"`
	NetworkMode        string      `json:"network_mode"`
	RegistryImageName  string      `json:"registry_image_name"`
	RuntimeProfileName string      `json:"runtime_profile_name"`
	ContainerType      string      `json:"container_type"`
	Total              int         `json:"total"`
	Critical           int         `json:"critical"`
	High               int         `json:"high"`
	Medium             int         `json:"medium"`
	Low                int         `json:"low"`
	Sensitive          int         `json:"sensitive"`
	Malware            int         `json:"malware"`
	Negligible         int         `json:"negligible"`
}
