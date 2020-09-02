package aqua

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// TODO don't export
func (cli *APIClient) GetContainersForWorkflow(ctx context.Context) (err error) {
	//vulns = make([]domain.ImageFinding, 0)
	page := 1

	//endpoint := strings.Replace(getVulnerabilities, "$IMAGENAME", image, 1)
	//endpoint = strings.Replace(endpoint, "$REGISTRYNAME", registry, 1)
	//
	//endpoint = strings.Replace(endpoint, " ", "%20", -1)
	endpoint := "/api/v1/containers?groupby=containers&status=running&page=1&pagesize=50"

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
				//vulnPage := &VulnerabilityPage{}
				//if err = json.Unmarshal(body, vulnPage); err == nil {
				//	if len(vulnPage.Result) == 0 {
				//		break
				//	} else {
				//
				//		for index := range vulnPage.Result {
				//			vulns = append(vulns, &vulnPage.Result[index])
				//		}
				//		page++
				//	}
				//} else {
				//	err = fmt.Errorf("error while parsing vulnerabilities from response - %s", err.Error())
				//	break
				//}

				var buf = &bytes.Buffer{}
				json.Indent(buf, body, "", "\t")
				fmt.Println(string(buf.Bytes()))
				break
			} else {
				err = fmt.Errorf("error while gathering vulnerabilities - %s", err.Error())
				break
			}
		} else {
			err = fmt.Errorf("error while making request - %s", err.Error())
			break
		}
	}

	return err
}
