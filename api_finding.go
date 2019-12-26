package dome9

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func (client *Client) getFindingsForRule(ruleHash string, bundleID int) (findings []*Finding, err error) {
	var body []byte

	// TODO implement paging like qualys
	var queryParams = fmt.Sprintf("?ruleLogicHash=%s&pageNumber=1&pageSize=100", ruleHash)
	body, err = client.executeRequest(http.MethodGet, strings.Replace(getBundleFindings+queryParams, "{bundleId}", strconv.Itoa(bundleID), 1), nil)

	if err == nil {
		findings = make([]*Finding, 0)
		err = json.Unmarshal(body, &findings)

		if err == nil {
			for index := range findings {
				findings[index].RuleHash = ruleHash
			}
		}
	}

	return findings, err
}
