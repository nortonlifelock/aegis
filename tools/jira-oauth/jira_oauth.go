package main

import (
	"flag"
	"fmt"

	"context"

	"github.com/benjivesterby/validator"
	"github.com/nortonlifelock/aegis/internal/config"
	"github.com/nortonlifelock/aegis/internal/database"
	"github.com/nortonlifelock/aegis/pkg/crypto"
	"github.com/nortonlifelock/aegis/pkg/domain"
	"github.com/nortonlifelock/aegis/pkg/jira"
	"github.com/nortonlifelock/aegis/pkg/log"
)

type logger struct{}

// Send prints a log to console
func (logger logger) Send(log log.Log) {
	val := log.Text
	if log.Error != nil {
		val += fmt.Sprintf(" %s", log.Error.Error())
	}
	fmt.Println(val)
}

func main() {
	// Setting up config arguments
	configFile := flag.String("config", "app.json", "The filename of the config to load.")
	configPath := flag.String("cpath", "", "The directory path of the config to load.")
	orgCode := flag.String("org", "", "The organization code for which you're editing. Example: LOCK, NORTON, IDA")

	flag.Parse()
	ctx := context.Background()

	if configFile != nil && configPath != nil {

		if appconfig, err := config.LoadConfig(*configPath, *configFile); err == nil {

			if validator.IsValid(appconfig) {

				var ms domain.DatabaseConnection
				if ms, err = database.NewConnection(appconfig); err == nil {

					var orgID string
					var sourceConfig domain.SourceConfig
					sourceConfig, orgID, err = getSourceConfig(orgCode, ms)
					if err == nil {
						logger := logger{}

						var decryptedConfig domain.SourceConfig
						if decryptedConfig, err = crypto.DecryptSourceConfig(ms, sourceConfig, appconfig); err == nil {

							var token string
							if _, token, err = jira.NewJiraConnector(ctx, logger, decryptedConfig); err == nil {
								encryptTokenInDatabase(token, ms, appconfig, orgID, sourceConfig)
							} else {
								fmt.Println(err.Error())
							}
						} else {
							fmt.Println(err.Error())
						}
					} else {
						fmt.Println(err.Error())
					}
				} else {
					fmt.Println(err.Error())
				}
			} else {
				fmt.Println("Invalid app config")
			}
		} else {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println("flags did not load properly")
	}
}

func encryptTokenInDatabase(token string, ms domain.DatabaseConnection, appconfig config.AppConfig, orgID string, sourceConfig domain.SourceConfig) {
	var err error

	if len(token) > 0 {
		var client crypto.Client
		if client, err = crypto.NewEncryptionClient(crypto.AES256, ms, appconfig.EncryptionKey(), orgID); err == nil {
			var encryptedToken string
			if encryptedToken, err = client.Encrypt(token); err == nil {
				if _, _, err = ms.UpdateSourceConfigToken(sourceConfig.ID(), encryptedToken); err == nil {
					fmt.Println("Token updated successfully")
				} else {
					fmt.Println(err.Error())
				}
			} else {
				fmt.Println(err.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println("Token did not extract. Either already exists in DB, or error occurred")
	}
}

func getSourceConfig(orgCode *string, ms domain.DatabaseConnection) (config domain.SourceConfig, orgID string, err error) {
	var sourceID string
	if orgCode != nil && len(*orgCode) > 0 {
		// Get the organization from the database using the id in the ticket object
		var torg domain.Organization
		if torg, err = ms.GetOrganizationByCode(*orgCode); err == nil {

			if torg != nil {
				fmt.Printf("Org [%s] Loaded\n", *orgCode)

				orgID = torg.ID()

				var src domain.Source
				if src, err = ms.GetSourceByName("JIRA"); err == nil {

					sourceID = src.ID()

					var sc []domain.SourceConfig
					if sc, err = ms.GetSourceConfigBySourceID(orgID, sourceID); err == nil {
						if len(sc) > 0 {
							config = sc[0]
						}
					}
				}
			}
		}
	}

	return config, orgID, err
}
