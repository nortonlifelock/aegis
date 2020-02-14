package main

import (
	"flag"
	"fmt"
	"github.com/nortonlifelock/aegis/internal/config"
	"github.com/nortonlifelock/crypto"
)

func main() {
	inputMessage := flag.String("m", "", "The value to encrypt/decrypt")
	key := flag.String("k", "", "The key to use for encryption")
	decrypt := flag.Bool("d", false, "Decrypt the message")
	profile := flag.String("p", "", "The KMS profile to use in ~/.aws")
	region := flag.String("r", "default", "The region your KMS key exists in")

	client := flag.String("client", "kms", "[kms|aes] delineates the encryption scheme")

	configFile := flag.String("config", "app.json", "The filename of the config to load.")
	configPath := flag.String("cpath", "", "The directory path of the config to load.")
	flag.Parse()

	var EDFlag = crypto.EncryptMode
	if *decrypt {
		EDFlag = crypto.DecryptMode
	}

	if len(*inputMessage) > 0 {

		if key != nil && len(*key) > 0 {
			encryptOrDecrypt(*client, *key, EDFlag, *inputMessage, *profile, *region)
		} else {
			appConfig, err := config.LoadConfig(*configPath, *configFile)
			if err == nil {
				encryptOrDecrypt(*client, appConfig.EncryptionKey(), EDFlag, *inputMessage, appConfig.KMSProfile(), appConfig.KMSRegion())
			} else {
				fmt.Println(fmt.Sprintf("Error while loading the application config  [%s]", err.Error()))
			}
		}

	} else {
		fmt.Printf("The message given for cryptography was empty")
	}
}

func encryptOrDecrypt(client string, key string, EDFlag int, inputMessage string, profile string, region string) {
	var encryptor crypto.Client
	var err error
	if len(profile) == 0 {
		encryptor, err = crypto.NewEncryptionClientWithDirectKey(client, key, region)
	} else {
		if client == crypto.KMS {
			encryptor, err = crypto.CreateKMSClientWithProfile(key, profile, region)
		} else {
			err = fmt.Errorf("can only specify profiles with KMS")
		}
	}

	if err == nil {
		if EDFlag == crypto.EncryptMode {
			result, err := encryptor.Encrypt(inputMessage)
			if err == nil {
				fmt.Printf("Encrypted message: %s\n", result)
			} else {
				fmt.Println(err.Error())
			}
		} else {
			result, err := encryptor.Decrypt(inputMessage)
			if err == nil {
				fmt.Printf("Decrypted message: %s\n", result)
			} else {
				fmt.Println(err.Error())
			}
		}
	} else {
		fmt.Println(err.Error())
	}
}
