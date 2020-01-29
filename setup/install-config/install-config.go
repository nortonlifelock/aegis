package install_config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/nortonlifelock/config"
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

	fmt.Print("Enter a url or IP pointing to your database (url recommended if IP is dynamic): ")
	conf.DatabasePath = getInput(reader)

	fmt.Print("Enter the port the database is listening on: ")
	conf.DatabasePort = getInput(reader)

	fmt.Print("Enter the username for database authentication: ")
	conf.DatabaseUser = getInput(reader)

	fmt.Print("Enter the password for database authentication: ")
	passwordByte, err := terminal.ReadPassword(0)
	check(err)
	conf.DatabasePassword = string(passwordByte)
	fmt.Println()

	fmt.Print("Enter the schema name that Aegis will be operating in: ")
	conf.DatabaseSchema = getInput(reader)

	fmt.Print("Enter AWS SNS topic ID that will be used to alert on critical logs (optional): ")
	conf.TopicKey = getInput(reader)

	fmt.Print("Enter AWS KMS key ID that will be used for encryption (not optional): ")
	conf.EKey = getInput(reader)

	fmt.Print("Enter the region that the AWS KMS key exists in (e.g. us-west-1)")
	conf.RegionKMS = getInput(reader)

	conf.ProfileKMS = "default"

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

	kmsClient, err := crypto.CreateKMSClientWithProfile(conf.EncryptionKey(), conf.ProfileKMS, conf.RegionKMS)
	check(err)

	conf.DatabasePassword, err = kmsClient.Encrypt(conf.DatabasePassword)
	check(err)

	body, err := json.MarshalIndent(conf, "", "\t")
	check(err)

	// the reader leaves trailing escaped newlines from their input
	body = []byte(strings.Replace(string(body), "\\n", "", -1))

	err = ioutil.WriteFile(fmt.Sprintf("%s/app.json", path), body, os.ModePerm)
	check(err)

	fmt.Println("Created app config, stored at", fmt.Sprintf("%s/app.json", path))
}

func getInput(reader *bufio.Reader) (userInput string) {
	var err error
	userInput, err = reader.ReadString('\n')
	check(err)

	userInput = strings.TrimSuffix(userInput, "\n")
	return userInput
}

func check(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
