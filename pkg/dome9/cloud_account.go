package dome9

import "fmt"

func (client *Client) determineCloudAccountTypeFromSubscriptionID(subscriptionID string) (vendor, externalAccountNumber, cloudAccountName string, err error) {
	var cloudAccount *CloudAccount
	if cloudAccount, err = client.GetCloudAccountByID(subscriptionID); err == nil {
		if cloudAccount != nil {
			vendor = cloudAccount.Vendor
			externalAccountNumber = cloudAccount.ExternalAccountNumber
			cloudAccountName = cloudAccount.Name
		}
	}

	if len(vendor) == 0 {
		if cloudAccount, err = client.GetAzureCloudAccountByID(subscriptionID); err == nil {
			if cloudAccount != nil {
				vendor = cloudAccount.Vendor
				externalAccountNumber = cloudAccount.ExternalAccountNumber
				cloudAccountName = cloudAccount.Name
			}
		}
	}

	if err == nil && len(vendor) == 0 {
		err = fmt.Errorf("could not determine vendor")
	}
	return vendor, externalAccountNumber, cloudAccountName, err
}
