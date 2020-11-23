package aqua

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/nortonlifelock/aegis/pkg/domain"
	"net/http"
)

type createExceptionReq struct {
	Issues  []ImageIssue `json:"issues"`
	Comment string       `json:"comment"`
}

type ImageIssue struct {
	IssueType string `json:"issue_type"`
	IssueName string `json:"issue_name"`
	ImageName string `json:"image_name"`

	ResourceType    string `json:"resource_type"`
	ResourceCpe     string `json:"resource_cpe"`
	ResourceName    string `json:"resource_name"`
	ResourceVersion string `json:"resource_version"`
	ResourcePath    string `json:"resource_path"`

	RegistryName string `json:"registry_name"`
}

func (cli *APIClient) CreateException(finding domain.ImageFinding, comment string) (err error) {
	if vulnerabilityResult, ok := finding.(*VulnerabilityResult); ok {
		req := &createExceptionReq{}
		req.Comment = comment
		req.Issues = []ImageIssue{{
			IssueType: "vulnerability",
			IssueName: vulnerabilityResult.Name,
			ImageName: vulnerabilityResult.ImageNameVar,

			ResourceType:    vulnerabilityResult.Resource.Type,
			ResourceCpe:     vulnerabilityResult.Resource.Cpe,
			ResourceName:    vulnerabilityResult.Resource.Name,
			ResourceVersion: vulnerabilityResult.Resource.Version,
			ResourcePath:    vulnerabilityResult.Resource.Path,

			RegistryName: vulnerabilityResult.RegistryVar,
		}}

		var body []byte
		if body, err = json.Marshal(req); err == nil {
			endpoint := postCreateException

			var request *http.Request
			if request, err = http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", cli.baseURL, endpoint), bytes.NewReader(body)); err == nil {
				if body, err = cli.executeRequest(request); err == nil {
				} else {
					err = fmt.Errorf("error while creating image scan - %s", err.Error())
				}
			} else {
				err = fmt.Errorf("error while making request - %s", err.Error())
			}
		} else {
			err = fmt.Errorf("error while creating request body")
		}
	} else {
		err = fmt.Errorf("error - did not get vulnerability result")
	}

	return err
}

func (cli *APIClient) GetExceptions(ctx context.Context) (vulns []domain.ImageFinding, err error) {
	vulns = make([]domain.ImageFinding, 0)
	page := 1

	endpoint := getExceptions

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		var request *http.Request
		if request, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s&page=%d&pagesize=50&order_by=-vulnerability", cli.baseURL, endpoint, page), nil); err == nil {
			var body []byte
			if body, err = cli.executeRequest(request); err == nil {
				vulnPage := &VulnerabilityPage{}
				if err = json.Unmarshal(body, vulnPage); err == nil {
					if len(vulnPage.Result) == 0 {
						break
					} else {

						for index := range vulnPage.Result {
							vulns = append(vulns, &vulnPage.Result[index])
						}
						page++
					}
				} else {
					err = fmt.Errorf("error while parsing vulnerabilities from response - %s", err.Error())
					break
				}
			} else {
				err = fmt.Errorf("error while gathering vulnerabilities - %s", err.Error())
				break
			}
		} else {
			err = fmt.Errorf("error while making request - %s", err.Error())
			break
		}
	}

	return vulns, err
}
