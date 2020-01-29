package crypto

import (
	"encoding/json"
	"fmt"

	"github.com/nortonlifelock/domain"
	"github.com/pkg/errors"
)

func sord(in *string) (out string) {
	if in != nil {
		out = *in
	}

	return out
}

type config interface {
	EncryptionKey() string
	KMSRegion() string
	KMSProfile() string
}

const (
	password    = "Password"
	privateKey  = "PrivateKey"
	consumerKey = "ConsumerKey"
	token       = "Token"
)

type overrideAuthInfo struct {
	domain.SourceConfig
	authInfo string
}

func (sc *overrideAuthInfo) AuthInfo() string {
	return sc.authInfo
}

// DecryptSourceConfig takes in a source config as an argument, and decrypts the fields that are expected to be encrypted
// Should not store encrypted pass inside the sourceConfig because when a client reconnects, it will try to decrypt the
// already decrypted password
func DecryptSourceConfig(ms domain.DatabaseConnection, sourceConfig domain.SourceConfig, config config) (domain.SourceConfig, error) {
	var err error

	if sourceConfig != nil {
		if len(config.EncryptionKey()) > 0 {

			var client Client
			if len(config.EncryptionKey()) > 0 {
				client, err = NewEncryptionClient(AES256, ms, config.EncryptionKey(), sourceConfig.OrganizationID(), config.KMSProfile(), config.KMSRegion())
				if err == nil {
					var authInfo map[string]interface{}
					if err = json.Unmarshal([]byte(sourceConfig.AuthInfo()), &authInfo); err == nil {
						encryptedPassword, _ := authInfo[password].(string)
						if len(encryptedPassword) > 0 {
							authInfo[password], err = client.Decrypt(encryptedPassword)
						}
						if err == nil {
							if err = decryptOauthFieldsForConfig(authInfo, client); err == nil {
								var authInfoByte []byte
								if authInfoByte, err = json.Marshal(authInfo); err == nil {
									sourceConfig = &overrideAuthInfo{
										sourceConfig,
										string(authInfoByte),
									}
								}
							}
						} else {
							err = fmt.Errorf("error while decrypting password - %s", err.Error())
						}
					}
				} else {
					err = errors.Errorf("could not create KMS Client - %s", err.Error())
				}
			} else {
				err = errors.New("no encryption key was found from the config file")
			}
		}
	} else {
		err = fmt.Errorf("cannot decrypt a nil source config")
	}

	return sourceConfig, err
}

func decryptOauthFieldsForConfig(authInfo map[string]interface{}, client Client) (err error) {
	encryptedPrivateKey, _ := authInfo[privateKey].(string)
	encryptedConsumerKey, _ := authInfo[consumerKey].(string)
	encryptedToken, _ := authInfo[token].(string)

	if len(encryptedPrivateKey) > 0 && len(encryptedConsumerKey) > 0 {

		authInfo[privateKey], err = client.Decrypt(encryptedPrivateKey)

		// Make sure we don't overwrite previous error
		if err == nil {
			authInfo[consumerKey], err = client.Decrypt(encryptedConsumerKey)
		}
	}
	if len(encryptedToken) > 0 && err == nil {
		authInfo[token], err = client.Decrypt(encryptedToken)
	}

	return err
}
