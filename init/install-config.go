package init

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/nortonlifelock/aegis/internal/config"
	"github.com/nortonlifelock/crypto"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"os"
	"strings"
)

func InstallConfig(path string) {
	conf := config.AppConfig{}
	var err error

	reader := bufio.NewReader(os.Stdin)

	conf.PathToAegis = path

	fmt.Print("Enter the path to the directory where logs will be stored: ")
	conf.LogFilePath = getInput(reader)

	if strings.Index(conf.LogFilePath, "/") < 0 && strings.Index(conf.LogFilePath, "\\") < 0 {
		conf.LogFilePath = fmt.Sprintf("%s/%s", path, conf.LogFilePath)
	}

	if _, err := os.Stat(conf.LogFilePath); os.IsNotExist(err) {
		_ = os.Mkdir(conf.LogFilePath, os.ModePerm)
	}

	// Controls whether logs are printed to the console
	conf.LogToConsole = true

	// Controls whether logs are stored in the database
	conf.LogToDb = true

	// Controls whether logs are written to a file
	conf.LogToFile = true

	// Controls whether logs are deleted after a day
	conf.LogDontDelete = false

	// Controls whether debug logs are processed
	conf.Debug = false

	// The port on which the Aegis API listens
	conf.APIServicePort = 4040

	// Options [ws|wss]
	conf.SocketProtocol = "wss"

	// Options [http|https]
	conf.Protocol = "https"

	// The url that the UI is being served at
	conf.UI = "localhost:4200"

	for conf.EncryptionMethod != crypto.KMS && conf.EncryptionMethod != crypto.AES256 && conf.EncryptionMethod != crypto.VAULT {
		fmt.Printf("What form of encryption would you like to use? Supported: AWS KMS/AES256/Hashicorp Vault secret storage. Ente one of [%s|%s|%s]\n", crypto.KMS, crypto.AES256, crypto.VAULT)
		conf.EncryptionMethod = getInput(reader)
	}

	var encryptionClient crypto.Client
	if conf.EncryptionMethod == crypto.KMS {
		fmt.Print("Enter AWS SNS topic ID that will be used to alert on critical logs (optional): ")
		conf.TopicKey = getInput(reader)

		fmt.Print("Enter AWS KMS symmetric encryption key ID that will be used for encryption (not optional): ")
		conf.EKey = getInput(reader)

		fmt.Print("Enter the region that the AWS KMS key exists in (e.g. us-west-1): ")
		conf.RegionKMS = getInput(reader)

		conf.ProfileKMS = "default"

		// The AWS SDK has the environment variables take precedence over the credentials in the shared file. CreateKMSClientWithProfile forces the loading of credentials from
		// the shared file, so if we have environment variables present we have to make sure that method isn't called
		var environmentVarsPresent bool
		if len(os.Getenv("AWS_ACCESS_KEY_ID")) > 0 && len(os.Getenv("AWS_SECRET_ACCESS_KEY")) > 0 {
			environmentVarsPresent = true
		}

		if environmentVarsPresent {
			encryptionClient, err = crypto.CreateKMSClient(conf.EncryptionKey(), conf.RegionKMS)
		} else {
			encryptionClient, err = crypto.CreateKMSClientWithProfile(conf.EncryptionKey(), conf.ProfileKMS, conf.RegionKMS)
		}

		check(err)
	} else if conf.EncryptionMethod == crypto.VAULT {
		fmt.Println("Because you are using Vault secret storage, two environment variables are required: VAULT_ADDR & VAULT_TOKEN")
		if len(os.Getenv("VAULT_ADDR")) == 0 {
			fmt.Println("Your VAULT_ADDR environment variable isn't set, we can set it for right now, but you will need to add it to a bash profile to make it persistent")
			vaultAddr := getInput(reader)
			_ = os.Setenv("VAULT_ADDR", vaultAddr)
		}

		if len(os.Getenv("VAULT_TOKEN")) == 0 {
			fmt.Println("Your VAULT_TOKEN environment variable isn't set, we can set it for right now, but you will need to add it to a bash profile to make it persistent")
			vaultAddr := getInput(reader)
			_ = os.Setenv("VAULT_TOKEN", vaultAddr)
		}

		encryptionClient, err = crypto.NewEncryptionClientWithDirectKey(crypto.VAULT, "", "")
	} else if conf.EncryptionMethod == crypto.AES256 {
		fmt.Println("Enter your AES encryption key (must be 32 bytes)")
		conf.EKey = getInput(reader)
		encryptionClient, err = crypto.NewEncryptionClientWithDirectKey(crypto.AES256, conf.EKey, "")
	} else {
		fmt.Printf("Unrecognized encryption scheme: [%s]", conf.EncryptionMethod)
		os.Exit(1)
	}

	fmt.Print("Enter a url or IP pointing to your database (url recommended if IP is dynamic): ")
	conf.DatabasePath = getInput(reader)

	fmt.Print("Enter the port the database is listening on: ")
	conf.DatabasePort = getInput(reader)

	fmt.Print("Enter the username for database authentication: ")
	conf.DatabaseUser = getInput(reader)

	if conf.EncryptionMethod == crypto.VAULT {
		fmt.Println("Enter the path to where your database password is stored in Vault, followed by a semicolon, followed by the key attached to the database password [e.g. secret/data/db;password]")
		conf.DatabasePassword = getInput(reader)
	} else {
		fmt.Print("Enter the password for database authentication: ")
		passwordByte, err := terminal.ReadPassword(0)
		check(err)
		conf.DatabasePassword = string(passwordByte)
		fmt.Println()

		conf.DatabasePassword, err = encryptionClient.Encrypt(conf.DatabasePassword)
		check(err)
	}

	fmt.Print("Enter the schema name that Aegis will be operating in: ")
	conf.DatabaseSchema = getInput(reader)

	body, err := json.MarshalIndent(conf, "", "\t")
	check(err)

	// the reader leaves trailing escaped newlines from their input
	body = []byte(strings.Replace(string(body), "\\n", "", -1))

	err = ioutil.WriteFile(fmt.Sprintf("%s/app.json", path), body, 0400)
	check(err)

	fmt.Println("Created app config, stored at", fmt.Sprintf("%s/app.json", path))
}
