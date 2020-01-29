package crypto

import (
	"fmt"
	"github.com/nortonlifelock/domain"
)

// Client manages
type Client interface {
	Encrypt(string) (string, error)
	Decrypt(string) (string, error)
}

const (
	// EncryptMode tells KMS to perform an encryption operation
	EncryptMode = iota

	// DecryptMode tells KMS to perform an decryption operation
	DecryptMode
)

const (
	// KMS is a const that delineates the type of encryption used (AWS Key Management System)
	KMS = "kms"

	// AES256 is a const that delineates the type of encryption used (AES256)
	AES256 = "aes"
)

// NewEncryptionClient takes in an application level encryption key (a KMS key)
// The fields in the database are encrypted with an organization specific encryption key, which is not the same as the application level encryption key
// This ensures that one organization cannot read the fields of another.
// The organization encryption key in the database itself is encrypted using the KMS application level encryption key
// The organization encryption key must be pulled from the database and decrypted before the client is created
// The application level encryption key should only exist in the root organization of an organization hierarchy
func NewEncryptionClient(clientType string, db domain.DatabaseConnection, applicationEncryptionKey string, orgID string, profile string, region string) (client Client, err error) {
	var orgKey string
	orgKey, err = decryptOrganizationKey(db, applicationEncryptionKey, orgID, profile, region)
	if err == nil {
		client, err = NewEncryptionClientWithDirectKey(clientType, orgKey, region)
	}

	return client, err
}

// NewEncryptionClientWithDirectKey takes the key used for encryption as a direct argument, and does not grab an encrypted, organization specific key from
// the database like NewEncryptionClient does
// region is only needed for kms encryption and can be empty when and AES client is being created
func NewEncryptionClientWithDirectKey(clientType string, key string, region string) (client Client, err error) {
	switch clientType {
	case AES256:
		client, err = createAESClient(key)
	case KMS:
		client, err = createKMSClient(key, region)
	default:
		err = fmt.Errorf("unrecognized encryption type [%s]", clientType)
	}

	return client, err
}

// The organization key is encrypted using KMS
func decryptOrganizationKey(db domain.DatabaseConnection, applicationEncryptionKey string, orgID string, profile string, region string) (orgKey string, err error) {
	if len(orgID) > 0 {
		var rootOrg domain.Organization
		rootOrg, err = getRootOrganization(db, orgID)
		if err == nil {
			if rootOrg != nil {
				if len(sord(rootOrg.EncryptionKey())) > 0 {
					orgKey, err = kmsDoEncryption(applicationEncryptionKey, DecryptMode, sord(rootOrg.EncryptionKey()), profile, region)
				} else {
					err = fmt.Errorf("root organization [%s] did not have an encryption key in the database", orgID)
				}
			} else {
				err = fmt.Errorf("could not find root organization for [%s]", orgID)
			}
		}
	} else {
		err = fmt.Errorf("no organization ID was provided to decryptOrganizationKey")
	}

	return orgKey, err
}

func getRootOrganization(db domain.DatabaseConnection, orgID string) (rootOrg domain.Organization, err error) {
	var traverse domain.Organization
	traverse, err = db.GetOrganizationByID(orgID)
	if err == nil {
		if traverse != nil {
			for traverse != nil && len(sord(traverse.ParentOrgID())) > 0 && err == nil {
				traverse, err = db.GetOrganizationByID(sord(traverse.ParentOrgID()))
				if err == nil {
					if traverse == nil {
						err = fmt.Errorf("could not find root organization of the org [%s]", orgID)
						break
					}
				}
			}
			rootOrg = traverse
		} else {
			err = fmt.Errorf("could not find an organization with ID [%s]", orgID)
		}
	}

	return rootOrg, err
}
