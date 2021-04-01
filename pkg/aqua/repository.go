package aqua

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (cli *APIClient) GetRepositories() (repos []RepositoryResult, err error) {
	repos = make([]RepositoryResult, 0)
	page := 1

	for {
		var request *http.Request
		if request, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s&page=%d&pagesize=50", cli.baseURL, getRepositories, page), nil); err == nil {
			var body []byte
			if body, err = cli.executeRequest(request); err == nil {
				repoPage := &Repositories{}
				if err = json.Unmarshal(body, repoPage); err == nil {
					if len(repoPage.Result) == 0 {
						break
					} else {
						repos = append(repos, repoPage.Result...)
						page++
					}
				} else {
					err = fmt.Errorf("error while parsing repositories from response - %s", err.Error())
					break
				}
			} else {
				err = fmt.Errorf("error while gathering repositories - %s", err.Error())
			}
		} else {
			err = fmt.Errorf("error while making request - %s", err.Error())
		}
	}

	return repos, err
}
