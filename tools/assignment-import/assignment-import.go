package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/benjivesterby/validator"
	"github.com/nortonlifelock/aegis/internal/config"
	"github.com/nortonlifelock/aegis/internal/database"
	"github.com/nortonlifelock/aegis/internal/integrations"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/files"
	"github.com/nortonlifelock/log"
)

type logger struct{}

// Send prints a log to the console
func (logger logger) Send(log log.Log) {
	val := log.Text
	if log.Error != nil {
		val += fmt.Sprintf(" %s", log.Error.Error())
	}
	fmt.Println(val)
}

func main() {
	var err error

	// Setting up config arguments for starting the job runner
	configFile := flag.String("config", "app.json", "The filename of the config to load.")
	configPath := flag.String("cpath", "", "The directory path of the config to load.")
	agFilePath := flag.String("agfile", "", "The path to the input assignment group csv.")
	orgCode := flag.String("org", "", "The organization that the assignment group belongs to.")
	source := flag.String("source", "", "The method of discovery that found the asset for the assignment group.")
	ticketingService := flag.String("ticketing", "", "The name of the ticketing service you're using to verify the existence of the AG")

	flag.Parse()

	logger := logger{}
	ctx := context.Background()

	ms, appconfig, ticketingSrc, sourceID, orgID, err := loadConfigs(configFile, configPath, orgCode, source, ticketingService, agFilePath)

	if err == nil {
		var engine integrations.TicketingEngine
		engine, err = integrations.GetEngine(ctx, *ticketingService, ms, logger, appconfig, ticketingSrc)
		if err == nil {
			fmt.Printf("Loading File [%s]\n", *agFilePath)

			var file *os.File
			if file, err = os.Open(*agFilePath); err == nil {

				fmt.Printf("[%s] Loaded\n", *agFilePath)

				var reader *bufio.Reader
				reader = bufio.NewReader(file)
				// read header
				_, _ = reader.ReadString('\r')

				var InvalidGroups = make(map[string]bool)

				for {

					var line string
					if line, err = reader.ReadString('\r'); err != io.EOF {
						err = nil

						if len(line) > 0 {

							err = processLine(line, engine, ms, sourceID, orgID, InvalidGroups)
							if err != nil {
								fmt.Println(err.Error())
							}
						} else {
							// TODO:
						}
					} else {
						break
					}
				}

				for k, v := range InvalidGroups {
					if v {
						fmt.Printf("JIRA could not find assignment group [%s]\n", k)
					}
				}

			} else {
				// TODO:
			}
		} else {
			// TODO
		}
	}

	if err != nil {
		fmt.Println(err)
	}
}

func processLine(line string, engine integrations.TicketingEngine, ms domain.DatabaseConnection, sourceID string, orgID string, InvalidGroups map[string]bool) (err error) {
	line = strings.Replace(line, "\n", "", -1)
	line = strings.Replace(line, "\r", "", -1)
	fmt.Printf("Processing Line [%s] \n", line)
	columns := strings.Split(line, ",")
	if len(columns) >= 2 {
		columns[0] = strings.TrimSpace(columns[0])
		columns[1] = strings.Replace(strings.TrimSpace(columns[1]), "\"", "", -1)

		if len(columns[0]) > 0 && len(columns[1]) > 0 {

			var exists bool
			if exists, err = engine.AssignmentGroupExists(columns[1]); err == nil && exists {
				fmt.Printf("UPDATED - IP: %s | GROUP: %s \n", columns[0], columns[1])
				if _, _, err = ms.SaveAssignmentGroup(sourceID, orgID, columns[0], columns[1]); err != nil {
					fmt.Printf("ERROR: %s \n", err.Error())
				}
			} else {
				InvalidGroups[columns[1]] = true
			}
		} else {
			// TODO:
		}
	} else {
		// TODO:
	}
	return err
}

func loadConfigs(configFile *string, configPath *string, orgCode *string, source *string, ticketingService *string, agFilePath *string) (ms domain.DatabaseConnection, appconfig config.AppConfig, ticketingSC domain.SourceConfig, sourceID, orgID string, err error) {
	if configFile != nil && configPath != nil {

		if orgCode != nil && len(*orgCode) > 0 {

			if source != nil && len(*source) > 0 {

				if ticketingService != nil && len(*ticketingService) > 0 {

					if files.ValidFile(*agFilePath) {

						fmt.Printf("Loading Config\n")
						if appconfig, err = config.LoadConfig(*configPath, *configFile); err == nil {

							if validator.IsValid(appconfig) {

								fmt.Printf("Opening Database Connection\n")
								if ms, err = database.NewConnection(appconfig); err == nil {
									fmt.Printf("Database Connection Established\n")
									orgID, sourceID, ticketingSC, err = loadOrgAndSource(orgCode, ms, orgID, source, sourceID, ticketingService, ticketingSC)
								} else {
									// TODO:
								}
							} else {
								err = fmt.Errorf("config is not valid, ensure that all proper fields are set in the config")
							}
						} else {
							// TODO:
						}
					} else {
						err = fmt.Errorf("assignment group file path invalid")
					}
				} else {
					err = fmt.Errorf("must include a ticketing service with the -ticketing flag")
				}
			} else {
				err = fmt.Errorf("must include a source with the -source flag")
			}
		} else {
			err = fmt.Errorf("must include an org with the -org flag")
		}
	} else {
		err = fmt.Errorf("must include a config file and path using -config and -cpath flags")
	}

	return ms, appconfig, ticketingSC, sourceID, orgID, err
}

func loadOrgAndSource(orgCode *string, ms domain.DatabaseConnection, orgID string, source *string, sourceID string, ticketingService *string, ticketingSC domain.SourceConfig) (string, string, domain.SourceConfig, error) {
	var err error

	fmt.Printf("Loading Org [%s]\n", *orgCode)
	// Get the organization from the database using the id in the ticket object
	var torg domain.Organization
	if torg, err = ms.GetOrganizationByCode(*orgCode); err == nil {

		if torg != nil {

			fmt.Printf("Org [%s] Loaded\n", *orgCode)

			orgID = torg.ID()

			fmt.Printf("Loading Source [%s]\n", *source)

			var src domain.Source
			if src, err = ms.GetSourceByName(*source); err == nil {
				fmt.Printf("Source [%s] Loaded\n", *source)

				// Ensure there is only one return
				sourceID = src.ID()

				fmt.Printf("Loading %s source config\n", *ticketingService)
				var ticketingSrc []domain.SourceConfig
				if ticketingSrc, err = ms.GetSourceConfigByNameOrg(*ticketingService, torg.ID()); err == nil && ticketingSrc != nil && len(ticketingSrc) > 0 {
					fmt.Printf("%s source config loaded", *ticketingService)
					ticketingSC = ticketingSrc[0]
				} else {
					err = fmt.Errorf("could not find ticketing source config - [%v]", err)
				}
			} else {
				// TODO:
			}
		} else {
			err = fmt.Errorf("organization Not Found")
		}
	} else {
		// TODO:
	}

	return orgID, sourceID, ticketingSC, err
}
