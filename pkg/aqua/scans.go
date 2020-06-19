package aqua

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (cli *apiClient) GetScanSummaries(ctx context.Context, registry string, imageAndTag string) (scan *Scan, err error) {
	endpoint := "/api/v1/scanqueue?order_by=-created" // order by created desc, the earlier the scan is in the response, the more recently it was created

	page := 1

	for scan == nil {
		select {
		case <-ctx.Done():
			return
		default:
		}

		var request *http.Request
		if request, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s&page=%d&pagesize=50", cli.baseURL, endpoint, page), nil); err == nil {
			var body []byte
			if body, err = cli.executeRequest(request); err == nil {
				scanPage := &ScanPage{}
				if err = json.Unmarshal(body, scanPage); err == nil {
					if len(scanPage.Result) == 0 {
						break
					} else {
						for _, findScan := range scanPage.Result {
							if findScan.Registry == registry && imageAndTag == findScan.Image && findScan.InitiatingUser == cli.username {
								scan = &findScan // we order by created desc, so the first scan we find is the most recent one
								break
							}
						}

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

	if err == nil && scan == nil {
		err = fmt.Errorf("could not find scan for registry/image [%s|%s]", registry, imageAndTag)
	}

	return scan, err
}
