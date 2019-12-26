package crypto

import (
	"bytes"
	"crypto/rand"
	b64 "encoding/base64"
	"encoding/gob"
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/pkg/errors"
	"golang.org/x/crypto/nacl/secretbox"
)

const (
	// UsEast defines the AWS region in which KMS is operating
	usEast = "us-east-1"

	// EncryptionType256 is a flag that tells the KMS client to encrypt using AES 256
	encryptionType256 = "AES_256"

	// EncryptionType128 is a flag that tells the KMS client to encrypt using AES 128
	encryptionType128 = "AES_128"

	keyLength   = 32
	nonceLength = 24
)

type payload struct {
	Key     []byte
	Nonce   *[nonceLength]byte
	Message []byte
}

// KMSClient holds all information required to perform encryption and decryption. Once the object is created, once can simply call
// encrypt or decrypt on it
type KMSClient struct {
	Client  *kms.KMS
	keyID   string
	KeySpec string
}

// CreateKMSClientWithProfile creates a KMSClient object. The keyID is the AWS KMS key ID. The profile is optional and may be passed as
// an empty string
func CreateKMSClientWithProfile(keyID string, profile string) (client *KMSClient, err error) {
	var kmsClient *kms.KMS
	var sess *session.Session
	if sess, err = session.NewSessionWithOptions(session.Options{
		Profile: profile,
	}); err == nil {

		var verifyCred credentials.Value
		if verifyCred, err = sess.Config.Credentials.Get(); err == nil {
			if len(verifyCred.AccessKeyID) > 0 && len(verifyCred.SecretAccessKey) > 0 {
				kmsClient = kms.New(sess, aws.NewConfig().WithRegion(usEast))
				client = &KMSClient{kmsClient, keyID, encryptionType256}
			} else {
				err = errors.Errorf("kms could not find profile [%s] in ~/.aws/credentials", profile)
			}
		}
	}

	return client, err
}

// createKMSClient performs the same operation as CreateKMSClientWithProfile, but does not use an associated profile
func createKMSClient(keyID string) (client *KMSClient, err error) {
	var kmsClient *kms.KMS
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(usEast),
	})

	if err == nil {
		kmsClient = kms.New(sess, aws.NewConfig().WithRegion(usEast))
		client = &KMSClient{kmsClient, keyID, encryptionType256}
	}
	return client, err
}

// Encrypt encrypts the argument using AWS KMS and the key ID in the KMS client
func (kmsClient *KMSClient) Encrypt(message string) (encryptedString string, err error) {
	var plaintext []byte
	var stringReader = strings.NewReader(message)
	plaintext, err = ioutil.ReadAll(stringReader)

	var buf = &bytes.Buffer{}

	var dataKeyInput = kms.GenerateDataKeyInput{KeyId: &kmsClient.keyID, KeySpec: &kmsClient.KeySpec}
	var dataKeyOutput *kms.GenerateDataKeyOutput
	dataKeyOutput, err = kmsClient.Client.GenerateDataKey(&dataKeyInput)

	if err == nil {

		payload := &payload{
			Key:   dataKeyOutput.CiphertextBlob,
			Nonce: &[nonceLength]byte{},
		}

		_, err = rand.Read(payload.Nonce[:])

		if err == nil {
			key := &[keyLength]byte{}
			copy(key[:], dataKeyOutput.Plaintext)

			payload.Message = secretbox.Seal(payload.Message, plaintext, payload.Nonce, key)
			err = gob.NewEncoder(buf).Encode(payload)
			encryptedString = b64.StdEncoding.EncodeToString([]byte(buf.Bytes()))
		}
	}
	return encryptedString, err
}

// Decrypt decrypts the argument using AWS KMS and the key ID in the KMS client
func (kmsClient *KMSClient) Decrypt(encryptedText string) (message string, err error) {
	var dataKeyOutput *kms.DecryptOutput
	var payload payload
	var plaintext []byte
	var encryptedTextInSlice []uint8

	encryptedTextInSlice, err = b64.StdEncoding.DecodeString(encryptedText)

	if err == nil {

		err = gob.NewDecoder(bytes.NewReader(encryptedTextInSlice)).Decode(&payload)
		if err == nil {
			dataKeyOutput, err = kmsClient.Client.Decrypt(&kms.DecryptInput{
				CiphertextBlob: payload.Key,
			})

			if err == nil {

				var success bool
				key := &[keyLength]byte{}
				copy(key[:], dataKeyOutput.Plaintext)
				plaintext, success = secretbox.Open(plaintext, payload.Message, payload.Nonce, key)

				if !success {
					err = errors.New("Unable to decrypt encrypted message")
				}

			}
		}
	}

	return string(plaintext), err
}

// kmsDoEncryption performs either an encryption or decryption operation depending on EDFlag using the AWS encryption key
func kmsDoEncryption(encryptionKey string, EDFlag int, message string, profile string) (enOrDe string, err error) {
	var client Client
	if len(profile) > 0 {
		client, err = CreateKMSClientWithProfile(encryptionKey, profile)
	} else {
		client, err = createKMSClient(encryptionKey)
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
