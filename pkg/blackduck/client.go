package blackduck

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type BlackDuckClient struct {
	bearerToken string
	baseUrl     string
	client      *http.Client
}

func NewBlackDuckClient(baseURL string, apiToken string, insecureSkipVerify bool) (client *BlackDuckClient, err error) {
	client = &BlackDuckClient{
		baseUrl: baseURL,
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: insecureSkipVerify,
				},
			},
		},
	}

	client.bearerToken, err = client.authenticateToken(apiToken)

	return client, err
}

func (cli *BlackDuckClient) executeRequest(method string, endpoint string, requestBody io.Reader) (respBody []byte, err error) {
	var request *http.Request
	if request, err = http.NewRequest(method, fmt.Sprintf("%s/%s", cli.baseUrl, endpoint), requestBody); err == nil {
		request.AddCookie(&http.Cookie{Name: "AUTHORIZATION_BEARER", Value: cli.bearerToken})

		var response *http.Response
		if response, err = cli.client.Do(request); err == nil {
			if respBody, err = ioutil.ReadAll(response.Body); err == nil {
				_ = response.Body.Close()
				if response.StatusCode >= 300 {
					err = fmt.Errorf("STATUS [%d] - %s", response.StatusCode, string(respBody))
				}
			} else {
				err = fmt.Errorf("error while reading resopnse body - %s", err.Error())
			}
		} else {
			err = fmt.Errorf("error while executing request - %s", err.Error())
		}
	} else {
		err = fmt.Errorf("error while making request - %s", err.Error())
	}

	return respBody, err
}

type BearerToken struct {
	BearerToken string `json:"bearerToken"`
}

// TODO this endpoint returns an expiration in milliseconds
// if we pass a context, we could kickoff a thread that waits that length in milliseconds, and if the context isn't closed
// we reauthorize
func (cli *BlackDuckClient) authenticateToken(apiToken string) (bearerToken string, err error) {
	var req *http.Request
	req, err = http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/tokens/authenticate", cli.baseUrl), nil)
	if err == nil {
		req.Header.Add("Authorization", fmt.Sprintf("token %s", apiToken))

		var resp *http.Response
		resp, err = cli.client.Do(req)
		if err == nil {
			var body []byte
			body, err = ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err == nil {
				parseBearer := &BearerToken{}

				err = json.Unmarshal(body, parseBearer)
				if err == nil {
					if len(parseBearer.BearerToken) > 0 {
						bearerToken = parseBearer.BearerToken
					} else {
						err = fmt.Errorf("empty bearer token")
					}
				} else {
					err = fmt.Errorf("error while parsing login response - %s", err.Error())
				}
			} else {
				err = fmt.Errorf("error while reading login response - %s", err.Error())
			}
		} else {
			err = fmt.Errorf("error while logging in - %s", err.Error())
		}
	} else {
		err = fmt.Errorf("error while building login request - %s", err.Error())
	}

	return bearerToken, err
}
