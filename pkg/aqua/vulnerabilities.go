package aqua

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (cli *APIClient) GetVulnerabilitiesForImage(ctx context.Context, image string, registry string) (vulns []*VulnerabilityResult, err error) {
	vulns = make([]*VulnerabilityResult, 0)
	page := 1

	endpoint := strings.Replace(getVulnerabilities, "$IMAGENAME", image, 1)

	if len(registry) > 0 {
		endpoint = fmt.Sprintf("%s&registry_name=%s", endpoint, registry)
	}


	endpoint = strings.Replace(endpoint, " ", "%20", -1)

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
