package dome9

import (
	"encoding/json"
	"net/http"
	"strings"
)

// The standard getAllBundles() method doesn't include the rule hash, this method does
func (client *Client) getBundleDetailedInfo(bundleID string) (bundleDetailed *Bundle, err error) {
	var body []byte
	body, err = client.executeRequest(http.MethodGet, strings.Replace(getBundleByID, "{bundleId}", bundleID, 1), nil)

	if err == nil {
		bundleDetailed = &Bundle{}
		err = json.Unmarshal(body, bundleDetailed)
	}

	return bundleDetailed, err
}

// GetAllBundles is useful for testing - do not delete
func (client *Client) GetAllBundles() (bundles []*Bundle, err error) {
	var body []byte
	body, err = client.executeRequest(http.MethodGet, getAccountBundles, nil)
	if err == nil {
		bundles = make([]*Bundle, 0)
		err = json.Unmarshal(body, &bundles)
	}

	return bundles, err
}
