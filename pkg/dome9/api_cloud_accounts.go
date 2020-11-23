package dome9

import (
	"encoding/json"
	"net/http"
	"strings"
)

// GetCloudAccounts helps gather cloud account IDs which is useful for running assessments on bundles
func (client *Client) GetCloudAccounts() (cloudAccounts []*CloudAccount, err error) {
	var body []byte
	if body, err = client.executeRequest(http.MethodGet, getCloudAccounts, nil); err == nil {
		cloudAccounts = make([]*CloudAccount, 0)
		err = json.Unmarshal(body, &cloudAccounts)
	}

	return cloudAccounts, err
}

// GetAzureCloudAccounts gathers the cloud accounts from azure subscriptions
func (client *Client) GetAzureCloudAccounts() (cloudAccounts []*CloudAccount, err error) {
	var body []byte
	if body, err = client.executeRequest(http.MethodGet, getAzureCloudAccounts, nil); err == nil {
		cloudAccounts = make([]*CloudAccount, 0)
		err = json.Unmarshal(body, &cloudAccounts)
	}

	return cloudAccounts, err
}

// GetCloudAccountByID helps gather cloud account IDs which is useful for running assessments on bundles
func (client *Client) GetCloudAccountByID(id string) (cloudAccount *CloudAccount, err error) {
	var body []byte
	if body, err = client.executeRequest(http.MethodGet, strings.Replace(getCloudAccountByID, "{id}", id, 1), nil); err == nil {
		cloudAccount = &CloudAccount{}
		err = json.Unmarshal(body, &cloudAccount)
	}

	return cloudAccount, err
}

// GetAzureCloudAccountByID is necessary as the normal cloud account endpoint does not return Azure cloud accounts
func (client *Client) GetAzureCloudAccountByID(id string) (cloudAccount *CloudAccount, err error) {
	var body []byte
	if body, err = client.executeRequest(http.MethodGet, strings.Replace(getAzureCloudAccount, "{id}", id, 1), nil); err == nil {
		cloudAccount = &CloudAccount{}
		err = json.Unmarshal(body, &cloudAccount)
	}

	return cloudAccount, err
}
