package aqua

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/nortonlifelock/log"
	"io/ioutil"
	"net/http"
)

type APIClient struct {
	client   *http.Client
	lstream  log.Logger
	username string
	password string
	jwt      string
	baseURL  string
}

func CreateClient(URL, username, password string, lstream log.Logger) (client *APIClient, err error) {
	client = &APIClient{
		client:   &http.Client{},
		username: username,
		password: password,
		baseURL:  URL,
		lstream:  lstream,
	}

	client.jwt, err = client.getJWT(client.username, client.password)

	return client, err
}

func (cli *APIClient) executeRequest(req *http.Request) (body []byte, err error) {
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	if len(cli.jwt) > 0 {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", cli.jwt))
	} else {
		req.SetBasicAuth(cli.username, cli.password)
	}

	var response *http.Response
	if response, err = cli.client.Do(req); err == nil {
		if response != nil {
			var success = response.StatusCode < 300

			if response.Body != nil {
				defer response.Body.Close()

				if body, err = ioutil.ReadAll(response.Body); err != nil {
					err = fmt.Errorf("error while reading response body - %s", err.Error())
				}
			}

			if !success {
				err = fmt.Errorf("CODE [%d] - response [%s]", response.StatusCode, string(body))
			}
		} else {
			err = fmt.Errorf("nil response")
		}
	} else {
		err = fmt.Errorf("error while making request - %s", err.Error())
	}

	return body, err
}

type login struct {
	ID   string `json:"id"`
	Pass string `json:"password"`
}

type bearer struct {
	Token string `json:"token"`
}

func (cli *APIClient) getJWT(username string, password string) (jwt string, err error) {
	log := &login{ID: username, Pass: password}

	var body []byte
	if body, err = json.Marshal(log); err == nil {
		var req *http.Request
		if req, err = http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v1/login", cli.baseURL), bytes.NewReader(body)); err == nil {
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Accept", "application/json")

			var resp *http.Response
			if resp, err = cli.client.Do(req); err == nil {
				if resp != nil {
					defer resp.Body.Close()

					if resp.StatusCode < 300 {
						if body, err = ioutil.ReadAll(resp.Body); err == nil {
							bearer := &bearer{}
							if err = json.Unmarshal(body, bearer); err == nil {
								jwt = bearer.Token
							}
						}
					} else {
						err = fmt.Errorf("error while logging in - code [%d]", resp.StatusCode)
					}
				} else {
					err = fmt.Errorf("empty response while attempting to login")
				}
			}
		}
	}

	return jwt, err
}
