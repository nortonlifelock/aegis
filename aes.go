package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// AESClient performs encryption and decryption using AES256
type AESClient struct {
	eKey []byte
}

// createAESClient first verifies that the provided key is the correct length (must be 32 bytes for AES256), and then returns a client that
// performs encryption and decryption using AES
func createAESClient(eKey string) (client *AESClient, err error) {
	var eKeyByte = []byte(eKey)
	if len(eKeyByte) == 32 {
		client = &AESClient{
			eKey: []byte(eKey),
		}
	} else {
		err = fmt.Errorf("encryption key must be 32 bytes, but was %d", len(eKeyByte))
	}

	return client, err
}

// Encrypt encrypts the input using AES256 and then encodes it in base64
func (client *AESClient) Encrypt(in string) (out string, err error) {
	var block cipher.Block
	block, err = aes.NewCipher(client.eKey)
	if err == nil {
		plaintext := []byte(in)
		cipherText := make([]byte, aes.BlockSize+len(plaintext))
		iv := cipherText[:aes.BlockSize]
		if _, err = io.ReadFull(rand.Reader, iv); err == nil {
			stream := cipher.NewCFBEncrypter(block, iv)
			stream.XORKeyStream(cipherText[aes.BlockSize:], plaintext)
			out = base64.StdEncoding.EncodeToString(cipherText)
		}
	}

	return out, err
}

// Decrypt takes a base64 encoded encrypted string, first decodes the string, and then decrypts it
func (client *AESClient) Decrypt(in string) (out string, err error) {
	var cipherText []byte
	cipherText, err = base64.StdEncoding.DecodeString(in)
	if err == nil {
		var block cipher.Block
		block, err = aes.NewCipher(client.eKey)
		if err == nil {
			if len(cipherText) >= aes.BlockSize {
				iv := cipherText[:aes.BlockSize]
				cipherText = cipherText[aes.BlockSize:]
				stream := cipher.NewCFBDecrypter(block, iv)
				// XORKeyStream can work in-place if the two arguments are the same.
				stream.XORKeyStream(cipherText, cipherText)
				out = fmt.Sprintf("%s", cipherText)
			} else {
				err = fmt.Errorf("cipher text is of length %d but must be at least %d", len(cipherText), aes.BlockSize)
			}
		}
	}

	return out, err
}
