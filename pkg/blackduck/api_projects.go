package blackduck

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (cli *BlackDuckClient) GetProject(projectID string) (resp *ProjectResponse, err error) {
	endpoint := strings.Replace(GetProject, "{PROJECT_ID}", projectID, 1)
	resp = &ProjectResponse{}

	var body []byte
	if body, err = cli.executeRequest(http.MethodGet, endpoint, nil); err == nil {
		err = json.Unmarshal(body, resp)
		if err != nil {
			err = fmt.Errorf("error while parsing response body - %s", err.Error())
		}
	} else {
		err = fmt.Errorf("error while getting project [%s] - %s", projectID, err.Error())
	}

	return resp, err
}

func (cli *BlackDuckClient) GetProjectVersions(projectID string) (resp *ProjectVersionResponse, err error) {
	endpoint := strings.Replace(GetProjectVersions, "{PROJECT_ID}", projectID, 1)
	resp = &ProjectVersionResponse{}

	var body []byte
	if body, err = cli.executeRequest(http.MethodGet, endpoint, nil); err == nil {
		err = json.Unmarshal(body, resp)
		if err != nil {
			err = fmt.Errorf("error while parsing response body - %s", err.Error())
		}
	} else {
		err = fmt.Errorf("error while getting project versions [%s] - %s", projectID, err.Error())
	}

	return resp, err
}
