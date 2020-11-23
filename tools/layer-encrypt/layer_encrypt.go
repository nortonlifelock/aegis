package main

import (
	"flag"
	"fmt"
	"github.com/nortonlifelock/aegis/internal/config"
	"github.com/nortonlifelock/aegis/internal/database"
	"github.com/nortonlifelock/aegis/pkg/crypto"
	"github.com/nortonlifelock/aegis/pkg/domain"
)

var appConfig config.AppConfig
var data domain.DatabaseConnection

func main() {
	path := flag.String("p", "", "Path to the app config")
	inputMessage := flag.String("m", "", "The value to encrypt/decrypt")
	decrypt := flag.Bool("d", false, "Decrypt the message")
	org := flag.String("org", "", "")
	profile := flag.String("profile", "", "Profile to use for kms")
	flag.Parse()

	start(*path)

	if org == nil || len(*org) == 0 {
		panic("must supply and organization code")
	}

	client, err := crypto.NewEncryptionClient(crypto.AES256, data, appConfig.EncryptionKey(), *org, *profile, appConfig.KMSRegion())
	if err != nil {
		panic(err)
	}

	if *decrypt {
		message, err := client.Decrypt(*inputMessage)
		if err == nil {
			fmt.Printf("Decrypted message:\n%s\n", message)
		} else {
			panic(err)
		}
	} else {
		message, err := client.Encrypt(*inputMessage)
		if err == nil {
			fmt.Printf("Encrypted message:\n%s\n", message)
		} else {
			panic(err)
		}
	}
}

func start(path string) {
	var err error
	appConfig, err = config.LoadConfigByPath(path)
	if err == nil {
		data, err = database.NewConnection(appConfig)
		if err == nil {

		} else {
			panic(err)
		}
	} else {
		panic(err)
	}
}
