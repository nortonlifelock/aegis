package dome9

import "fmt"

func (client *Client) determineCloudAccountTypeFromSubscriptionID(subscriptionID string) (vendor string, err error) {
	var cloudAccount *CloudAccount
	if cloudAccount, err = client.GetCloudAccountByID(subscriptionID); err == nil {
		if cloudAccount != nil {
			vendor = cloudAccount.Vendor
		}
	}

	if len(vendor) == 0 {
		if cloudAccount, err = client.GetAzureCloudAccountByID(subscriptionID); err == nil {
			if cloudAccount != nil {
				vendor = cloudAccount.Vendor
			}
		}
	}

	if err == nil && len(vendor) == 0 {
		err = fmt.Errorf("could not determine vendor")
	}
	return vendor, err
}
