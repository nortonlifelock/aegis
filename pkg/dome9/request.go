package dome9

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func (client *Client) executeRequest(method string, endpoint string, data io.Reader) (body []byte, err error) {
	headers := map[string][]string{
		"Content-Type":  {"application/json"},
		"Accept":        {"application/json"},
		"Authorization": {"Basic " + client.authString},
	}

	var base = client.baseURL
	url := fmt.Sprintf("%s%s", base, endpoint)

	var req *http.Request
	if req, err = http.NewRequest(method, url, data); err == nil {
		req.Header = headers

		var resp *http.Response
		if resp, err = client.client.Do(req); err == nil {

			defer resp.Body.Close()
			if body, err = ioutil.ReadAll(resp.Body); err == nil {
				if !strings.Contains(resp.Status, strconv.Itoa(http.StatusOK)) {
					err = fmt.Errorf("response code [%s] returned - %s", resp.Status, string(body))
				}
			}
		}
	}

	return body, err
}
