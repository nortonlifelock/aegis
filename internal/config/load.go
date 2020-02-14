package config

import (
	"encoding/json"
	"fmt"
	"github.com/nortonlifelock/crypto"
	"github.com/nortonlifelock/files"
	"os"
)

// LoadConfigByPath loads the application config using a single path, decrypts the encrypted fields, and returns a struct that implements the
// config interface
func LoadConfigByPath(configPath string) (AppConfig, error) {
	var config = &AppConfig{}
	var err error

	if files.ValidFile(configPath) {

		var configJSON string
		if configJSON, err = files.GetStringFromFile(configPath); err == nil {
			if len(configJSON) > 0 {
				if err = json.Unmarshal([]byte(configJSON), config); err == nil {

					// the database information is encrypted using KMS
					var client crypto.Client
					client, err = crypto.NewEncryptionClientWithDirectKey(crypto.KMS, config.EncryptionKey(), config.KMSRegion())
					if err == nil {
						var decrypted string
						decrypted, err = client.Decrypt(config.DBPassword())

						if err == nil {
							config.DatabasePassword = decrypted
						} else {
							err = fmt.Errorf("could not decrypt database password: %s", err.Error())
						}
					} else {
						err = fmt.Errorf("could not load client for key decryption - %s", err.Error())
					}
				}
			} else {
				err = fmt.Errorf("empty config file at path [%s]. Unable to start application", configPath)
			}
		}
	} else {
		err = fmt.Errorf("invalid config path passed to Load AppConfig [%s]. Unable to start application", configPath)
	}

	return *config, err
}

// LoadConfig loads the application config using a path to the config, and the name of the config,
// and then decrypts the encrypted fields, and returns a struct that implements the
// config interface
func LoadConfig(path string, file string) (config AppConfig, err error) {

	if len(file) > 0 {

		var workingDir string
		if workingDir, err = os.Getwd(); err == nil {

			var configPath string

			// Determine if the path passed in is a valid path and if the file exists if not use the default
			if len(path) > 0 {
				configPath = fmt.Sprintf("%s%s%s", path, string(os.PathSeparator), file)
			} else {
				configPath = fmt.Sprintf("%s%s%s", workingDir, string(os.PathSeparator), file)
			}

			config, err = LoadConfigByPath(configPath)
		}
	} else {
		err = fmt.Errorf("invalid filename passed to Load AppConfig. Unable to start application")
	}

	return config, err
}
