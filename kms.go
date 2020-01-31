package crypto

import (
	b64 "encoding/base64"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/pkg/errors"
)

const (
	// EncryptionType256 is a flag that tells the KMS client to encrypt using AES 256
	encryptionType256 = "AES_256"

	// EncryptionType128 is a flag that tells the KMS client to encrypt using AES 128
	encryptionType128 = "AES_128"
)

// KMSClient holds all information required to perform encryption and decryption. Once the object is created, once can simply call
// encrypt or decrypt on it
type KMSClient struct {
	Client  *kms.KMS
	keyID   string
	KeySpec string
}

// CreateKMSClientWithProfile creates a KMSClient object. The keyID is the AWS KMS key ID. The profile is optional and may be passed as
// an empty string
func CreateKMSClientWithProfile(keyID string, profile string, region string) (client *KMSClient, err error) {
	var kmsClient *kms.KMS
	var sess *session.Session
	if sess, err = session.NewSessionWithOptions(session.Options{
		Profile: profile,
	}); err == nil {

		var verifyCred credentials.Value
		if verifyCred, err = sess.Config.Credentials.Get(); err == nil {
			if len(verifyCred.AccessKeyID) > 0 && len(verifyCred.SecretAccessKey) > 0 {
				kmsClient = kms.New(sess, aws.NewConfig().WithRegion(region))
				client = &KMSClient{kmsClient, keyID, encryptionType256}
			} else {
				err = errors.Errorf("kms could not find profile [%s] in ~/.aws/credentials", profile)
			}
		}
	}

	return client, err
}

// CreateKMSClient performs the same operation as CreateKMSClientWithProfile, but does not use an associated profile
func CreateKMSClient(keyID string, region string) (client *KMSClient, err error) {
	var kmsClient *kms.KMS
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	if err == nil {
		kmsClient = kms.New(sess, aws.NewConfig().WithRegion(region))
		client = &KMSClient{kmsClient, keyID, encryptionType256}
	}
	return client, err
}

// Encrypt encrypts the argument using AWS KMS and the key ID in the KMS client
func (kmsClient *KMSClient) Encrypt(message string) (encryptedString string, err error) {
	var output *kms.EncryptOutput
	output, err = kmsClient.Client.Encrypt(&kms.EncryptInput{
		KeyId:     &kmsClient.keyID,
		Plaintext: []byte(message),
	})

	if err == nil {
		encryptedString = b64.StdEncoding.EncodeToString(output.CiphertextBlob)
	}

	return encryptedString, err
}

// Decrypt decrypts the argument using AWS KMS and the key ID in the KMS client
func (kmsClient *KMSClient) Decrypt(encryptedText string) (message string, err error) {
	var encryptedTextInSlice []byte
	encryptedTextInSlice, err = b64.StdEncoding.DecodeString(encryptedText)

	var output *kms.DecryptOutput
	output, err = kmsClient.Client.Decrypt(&kms.DecryptInput{
		KeyId:          &kmsClient.keyID,
		CiphertextBlob: encryptedTextInSlice,
	})

	if err == nil {
		message = string(output.Plaintext)
	}

	return message, err
}

// kmsDoEncryption performs either an encryption or decryption operation depending on EDFlag using the AWS encryption key
func kmsDoEncryption(encryptionKey string, EDFlag int, message string, profile string, region string) (enOrDe string, err error) {
	var client Client
	if len(profile) > 0 {
		client, err = CreateKMSClientWithProfile(encryptionKey, profile, region)
	} else {
		client, err = CreateKMSClient(encryptionKey, region)
	}

	if err == nil {

		if EDFlag == EncryptMode {
			var encryptedMessage string
			encryptedMessage, err = client.Encrypt(message)
			if err == nil {
				enOrDe = encryptedMessage
			} else {
				err = errors.Errorf("error while encrypting message [%s]", err.Error())
			}

		} else if EDFlag == DecryptMode {
			var decryptedMessage string
			decryptedMessage, err = client.Decrypt(message)
			if err == nil {
				enOrDe = decryptedMessage
			} else {
				err = errors.Errorf("error while decrypting message [%s]", err.Error())
			}
		} else {
			err = errors.New("neither the encryption nor the decryption flag were set")
		}
	} else {
		err = errors.Errorf("error while creating kms client [%s]", err.Error())
	}

	return enOrDe, err
}
