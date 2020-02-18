package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/benjivesterby/validator"
	"github.com/nortonlifelock/aegis/internal/config"
	"github.com/nortonlifelock/aegis/internal/database"
	"github.com/nortonlifelock/aegis/internal/integrations"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/files"
	"github.com/nortonlifelock/log"
)

const csvDate = "01-02-2006"
const filenameDate = "01-02-2006_1504MST"

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
	outputPath := flag.String("o", "", "The output file for the resulting file")
	configFile := flag.String("config", "app.json", "The filename of the config to load.")
	configPath := flag.String("cpath", "", "The directory path of the config to load.")
	orgCode := flag.String("org", "", "The organization code for which you're editing. Example: LOCK, NORTON, IDA")

	flag.Parse()

	var appConfig config.AppConfig

	if configFile != nil && configPath != nil {

		var outputFile = *outputPath
		if len(outputFile) == 0 {
			if outputFile, err = os.Getwd(); err == nil {
				outputFile = fmt.Sprintf("%s%sQualysExport_%s.csv", outputFile, string(os.PathSeparator), time.Now().Format(filenameDate))
			} else {
				panic(0)
			}
		}

		if appConfig, err = config.LoadConfig(*configPath, *configFile); err == nil {

			if validator.IsValid(appConfig) {

				var ms domain.DatabaseConnection
				if ms, err = database.NewConnection(appConfig); err == nil {
					err = getConfigsAndPerformExport(orgCode, ms, appConfig, outputFile)
				} else {
					fmt.Printf("ERROR Opening DB Connection : %s", err.Error())
				}
			}
		}
	}

	fmt.Println("Export Complete")

}

func getConfigsAndPerformExport(orgCode *string, ms domain.DatabaseConnection, appConfig config.AppConfig, outputFile string) (err error) {
	var orgID string
	var sourceID string
	if orgCode != nil && len(*orgCode) > 0 {
		// Get the organization from the database using the id in the ticket object
		var torg domain.Organization
		if torg, err = ms.GetOrganizationByCode(*orgCode); err == nil {

			if torg != nil {
				fmt.Printf("Org [%s] Loaded\n", *orgCode)

				orgID = torg.ID()

				var src domain.Source
				if src, err = ms.GetSourceByName("QUALYS"); err == nil {

					sourceID = src.ID()

					var sc []domain.SourceConfig
					if sc, err = ms.GetSourceConfigBySourceID(orgID, sourceID); err == nil {
						if len(sc) > 0 {
							var source = sc[0]

							var entry string

							reader := bufio.NewReader(os.Stdin)

							for len(entry) <= 0 && err == nil {
								fmt.Print("Enter Asset Groups to Export (Comma Separated): ")
								if entry, err = reader.ReadString('\n'); err == nil {
									entry = strings.TrimSpace(entry)

									err = exportAssetGroups(entry, src, ms, appConfig, source, outputFile)
								} else {
									fmt.Printf("Error reading entry from command line : %s", err.Error())
								}
							}
						} else {
							fmt.Printf("No Source Configs")
						}
					} else {
						fmt.Printf("ERROR Getting Source Config : %s", err.Error())
					}
				} else {
					fmt.Printf("Unable to load Qualys Source Config : %s", err.Error())
				}
			} else {
				fmt.Println("Organization Not Found")
			}
		} else {
			fmt.Printf("Unable to load organization : %s", err.Error())
		}
	} else {
		fmt.Println("Organization (-org) is required. Example: LOCK, NORTON, IDA")
	}
	return err
}

func exportAssetGroups(entry string, src domain.Source, ms domain.DatabaseConnection, appConfig config.AppConfig, source domain.SourceConfig, outputFile string) (err error) {
	var wg = &sync.WaitGroup{}

	if len(entry) > 0 {
		var groups = strings.Split(entry, ",")

		if len(groups) > 0 {

			var dgs = make([]string, 0)

			for group := range groups {
				var groupID int

				if groupID, err = strconv.Atoi(groups[group]); err == nil {
					dgs = append(dgs, strconv.Itoa(groupID))
				} else {
					fmt.Printf("Unable to parse [%s] Asset Group\n", groups[group])
				}
			}

			if err == nil {

				fmt.Println("Connecting to Qualys")
				ctx := context.Background()

				var scanner integrations.Vscanner
				if scanner, err = integrations.NewVulnScanner(ctx, src.Source(), ms, logger{}, appConfig, source); err == nil {

					if devVulns, err := scanner.Detections(context.Background(), dgs); err == nil {

						dump(outputFile, fmt.Sprintf("%v,%s,%s,%s,%s,%s,%s,%s,%s,%s\n",
							"QID",
							"IP",
							"HOSTNAME",
							"OS",
							"VULNERABILITY",
							"ACTIVEKERNEL",
							"PORTS",
							"CVSS",
							"CVE",
							"LASTDETECTED"))

						for {
							if dv, ok := <-devVulns; ok {
								wg.Add(1)

								go func(devvuln domain.Detection) {
									defer wg.Done()
									extractInformationAndDumpToFile(devvuln, outputFile)
								}(dv)

							} else {
								break
							}
						}
					}

					wg.Wait()
				} else {
					fmt.Printf("Error initializing Qualys : %s", err.Error())
				}
			} else {
				fmt.Printf("Error reading groups : %s", err.Error())
			}
		} else {
			fmt.Println("No groups loaded")
		}
	} else {
		fmt.Println("No groups loaded")
	}
	return err
}

func extractInformationAndDumpToFile(devvuln domain.Detection, outputFile string) {
	if device, err := devvuln.Device(); err == nil {
		if vuln, err := devvuln.Vulnerability(); err == nil {
			if detected, err := devvuln.Detected(); err == nil {
				var ips = device.IP()
				var hosts = device.HostName()

				var port string
				if devvuln.Port() > 0 {
					port = fmt.Sprintf("%v %s", devvuln.Port(), devvuln.Protocol())
				}

				var activeKernel string
				if devvuln.ActiveKernel() != nil {
					activeKernel = fmt.Sprintf("%v", *devvuln.ActiveKernel())
				}

				dump(outputFile, fmt.Sprintf("%v,%s,\"%s\",\"%s\",\"%s\",%s,%s,%v,%s,%s\n",
					vuln.SourceID(),
					ips,
					hosts,
					device.OS(),
					vuln.Name(),
					activeKernel,
					port,
					vuln.CVSS2(),
					getReferences(vuln),
					detected.Format(csvDate)))

			} else {
				fmt.Println("error while gathering detection date", err.Error())
			}
		} else {
			fmt.Println("error while gathering vulnerability", err.Error())
		}
	} else {
		fmt.Println("error while gathering device", err.Error())
	}
}

func getReferences(vuln domain.Vulnerability) string {
	var result string
	if out, err := vuln.References(context.Background()); err == nil {
		func() {
			for {
				if reference, ok := <-out; ok {
					result += reference.Reference() + ","
				} else {
					return
				}
			}
		}()
	} else {
		fmt.Println("error while gathering references", err.Error())
	}

	result = strings.TrimRight(result, ",")
	return result
}

var mut = sync.Mutex{}

func dump(path string, lineOut string) {
	mut.Lock()

	_, err := os.Stat(path)
	//longest filesize is a megabyte
	if err == nil {
		var logFile *os.File
		logFile, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err == nil {
			defer func() {
				_ = logFile.Close()
			}()
			_, err = logFile.WriteString(lineOut)
		} else {
			fmt.Println(fmt.Sprintf("FAILED APPENDING LOG CONTENTS - %s", err.Error()))
		}
	} else { //The file didn't exist yet
		err = files.WriteFile(path, lineOut)
		if err != nil {
			fmt.Println("FAILED WRITING INITIAL LOG - ", err.Error())
		}
	}

	mut.Unlock()
}
