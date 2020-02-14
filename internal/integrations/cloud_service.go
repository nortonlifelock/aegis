package integrations

import (
	"fmt"
	"github.com/nortonlifelock/aws"
	"github.com/nortonlifelock/azure"
	"github.com/nortonlifelock/crypto"
	"github.com/nortonlifelock/domain"
)

const (
	// AWS is a string that identifies the connection is an AWS connection
	AWS = "AWS"

	// Azure is a string that identifies the connection is an Azure connection
	Azure = "Azure"
)

// CloudServiceConnection defines the methods that are required to grab tag information from a cloud service provider
type CloudServiceConnection interface {
	GetAllTagNames() (tagNames []string, err error)
	GetIPTagMapping() (ipToKeyToValue map[domain.CloudIP]map[string]string, err error)
	IPAddresses() (ips []domain.CloudIP, err error)
}

type config interface {
	EncryptionKey() string
	KMSRegion() string
	KMSProfile() string
}

// GetCloudServiceConnection returns a struct that implements the CloudServiceConnection interface
func GetCloudServiceConnection(ms domain.DatabaseConnection, cloudServiceID string, config domain.SourceConfig, appconfig config, lstream logger) (connection CloudServiceConnection, err error) {
	var decryptedConfig domain.SourceConfig
	decryptedConfig, err = crypto.DecryptSourceConfig(ms, config, appconfig)

	if err == nil {
		if len(cloudServiceID) > 0 {

			switch cloudServiceID {
			case AWS:
				connection, err = awsclient.CreateConnection(decryptedConfig.AuthInfo(), sord(config.Payload()))
				break
			case Azure:
				connection, err = azure.CreateConnection(decryptedConfig.AuthInfo(), decryptedConfig.Address(), lstream)
				break
			default:
				err = fmt.Errorf("unrecognized cloud service [%s]", cloudServiceID)
			}

		} else {
			err = fmt.Errorf("cloudServiceId passed empty to GetCloudServiceConnection")
		}
	}

	return connection, err
}
