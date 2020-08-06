package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/nortonlifelock/aegis/tools/jira-tool/tool"
	"os"
	"os/signal"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/nortonlifelock/aegis/internal/config"
	"github.com/nortonlifelock/aegis/internal/database"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/files"
)

var (
	separator = "\r\n"
)

//password hashing in go form benji
//there are tickets that relatively relate, can close as duplicates and create more relevant

func main() {
	/*
		want to talk to them about what method of discovery should be used for imports
			specific to tool or something that they set
				for inserts

		have a mode where a JQL query is made, and each ticket returned is updated
			make this an extra ticket ^
			ask scott if this would actually be useful along with the CSV import
	*/

	var appConfigPath = flag.String("config", "", "Path to the app json")
	var filePath = flag.String("path", "", "Path to the CSV file for bulk updates")
	var fileName = flag.String("file", "", "CSV file for bulk updates")
	//var username = flag.String("u", "", "Username for JIRA connector")
	//var password = flag.String("p", "", "Password for JIRA connector")
	var address = flag.String("url", "", "Address for JIRA")
	var help = flag.Bool("h", false, "Print out help dialogue")
	var timeFormat = flag.String("time", "", "Sets time format for input")
	var msec = flag.Int("msec", 300, "Sets the time between JIRA requests")

	flag.Parse()

	var err error
	appconfig, filePath, err := processFlags(msec, timeFormat, help, filePath, appConfigPath, fileName, address)
	if err == nil {
		var ms domain.DatabaseConnection

		if ms, err = database.NewConnection(appconfig); err == nil {
			var templateConfig domain.SourceConfig

			templateConfig, err = ms.GetSourceOauthByURL(*address)
			if err == nil {

				var fileContentsAsBytes *bytes.Buffer
				var err error
				fileContentsAsBytes, err = files.GetFileContents(fmt.Sprintf("%s/%s", *filePath, *fileName))

				if err == nil {
					var fileContents = fileContentsAsBytes.String()
					if strings.Index(fileContents, "\r\n") < 0 { //\r\n is not the separator
						if strings.Index(fileContents, "\r") >= 0 {
							fileContents = strings.Replace(fileContents, "\r", "\n", -1)
							separator = "\n"
						} else if strings.Index(fileContents, "\n") >= 0 {
							separator = "\n"
						} else {
							fmt.Printf("WARNING could not identify separator\n")
						}
					}

					fileContents = removeNonASCIICharacters(fileContents)

					var fileAsSlice []string
					fileAsSlice = strings.Split(fileContents, separator)

					payload := tool.MakePayload(ms, fileContents, templateConfig, appconfig, progressPrint)

					//This starts a thread that listens for a sigint. Once heard, it creates files for the failed and unfinished lines
					signalChan := make(chan os.Signal, 1)
					signal.Notify(signalChan, os.Interrupt)
					go func() {
						for range signalChan {
							fmt.Println() //ensures ^C isn't on the same line as the next print statement
							_ = cleanUp(payload.BlockWG, fileAsSlice, *fileName, *filePath, nil, payload.LineNumOfSucceed, payload.LineNumOfFailed, payload.LineNumToDescLine, payload.CommandSuccess, payload.CommandFailure)
							os.Exit(0)
						}
					}()

					err = tool.ProcessCSVContents(payload)

					if err != nil && !strings.Contains(err.Error(), "EOF") {
						fmt.Println(fmt.Sprintf("\nError during processing of CSV - %s", err.Error()))
					}

					elapsed := time.Since(payload.StartTime).Round(time.Second)
					fmt.Println(fmt.Sprintf("\nExecution took %s", elapsed))
					err = cleanUp(payload.BlockWG, fileAsSlice, *fileName, *filePath, err, payload.LineNumOfSucceed, payload.LineNumOfFailed, payload.LineNumToDescLine, payload.CommandSuccess, payload.CommandFailure)

					if err != nil {
						fmt.Println(fmt.Sprintf("Error while creating termination files - %s", err.Error()))
					}
				} else {
					fmt.Println(err.Error())
				}

			} else {
				fmt.Println(err.Error())
			}

		} else {
			fmt.Println(fmt.Sprintf("Error while opening file %s - %s", *fileName, err.Error()))
		}
	} else {
		fmt.Println(err.Error())
	}
}

func processFlags(msec *int, timeFormat *string, help *bool, filePath *string, appConfigPath *string, fileName *string, address *string) (appconfig config.AppConfig, path *string, err error) {
	if msec != nil {
		tool.TimeToWaitBetweenLines = time.Duration(*msec) * time.Millisecond
	}
	if len(*timeFormat) > 0 {
		tool.TimeLayout = *timeFormat
	}
	if *help {
		printHelp()
	}
	if len(*filePath) == 0 {
		workingDir, err := os.Getwd()
		if err == nil {
			*filePath = workingDir
		} else {
			panic(fmt.Sprint("a path to your file must be supplied with the -path flag, as the working directory could not be inferred at runtime"))
		}
	}

	if len(*appConfigPath) > 0 {
		*appConfigPath = strings.TrimRight(strings.TrimLeft(*appConfigPath, "“"), "”")
		appconfig, err = config.LoadConfigByPath(*appConfigPath)
		if err == nil {
			if len(*fileName) > 0 {
				if len(*address) > 0 {

				} else {
					err = fmt.Errorf("a JIRA address must be specified using the -url flag")
				}
			} else {
				err = fmt.Errorf("must pass the filename of your CSV file using the -file flag")
			}
		} else {
			err = fmt.Errorf("error while loading app json - " + err.Error())
		}
	} else {
		err = fmt.Errorf("must provide a path to your app.json using the -config flag")
	}

	return appconfig, filePath, err
}

func cleanUp(blockWG *sync.WaitGroup, fileAsSlice []string, fileName string, filePath string, err error, lineNumOfSucceed *[]int, lineNumOfFailed *[]int, lineNumToDescLine map[int]string, commandSuccess *int, commandFailure *int) error {
	// give time for the threads to finish up
	blockWG.Wait()

	var totalCommandsHandled = append(*lineNumOfSucceed, *lineNumOfFailed...)
	if err != nil { //should be an EOF error
		fmt.Printf("\nProgram exiting: %s\n\n", err.Error())
		err = nil
	} else {
		fmt.Printf("\nProgram exiting\n\n")
	}

	//fmt.Printf("Finished with [%d] successful connections [%d] failed connections\n",
	//	connectSuccess, connectFailure)
	fmt.Printf("There were [%d] successful commands and [%d] failed commands\n", *commandSuccess, *commandFailure)
	if len(*lineNumOfFailed) > 0 {
		sort.Ints(*lineNumOfFailed)
		var failedFileContents = ""
		var recentDescLine = ""

		for _, failedLine := range *lineNumOfFailed {

			//don't reprint the description line if the failed lines use the same description line
			if recentDescLine == lineNumToDescLine[failedLine] {
				failedFileContents = fmt.Sprintf("%s\n%s", failedFileContents, fileAsSlice[failedLine-1])
			} else {
				recentDescLine = lineNumToDescLine[failedLine]
				var failedLine = fileAsSlice[failedLine-1]
				failedLine = strings.Replace(failedLine, separator, "", -1)
				failedFileContents = fmt.Sprintf("%s\n%s\n%s", failedFileContents, recentDescLine, failedLine)
			}

		}
		if failedFileContents[0] == '\n' {
			failedFileContents = failedFileContents[1:]
		}
		var failedFileName = fmt.Sprintf("failed_%s", fileName)
		fmt.Printf("Writing %s...\n", failedFileName)
		err = files.WriteFile(fmt.Sprintf("%s/%s.csv", filePath, failedFileName), failedFileContents)
		if err != nil {
			fmt.Printf("Error while writing failed csv [%s]\n", err.Error())
		}
	}

	if len(totalCommandsHandled) > 0 {
		//because the slice is sorted, I can resume my search at the position where I saw the last element
		sort.Ints(totalCommandsHandled)

		var unfinishedCommandFileContents = ""

		var lastSeenIndex = 0
		//starting at the number 1 up to and equaling the length of fileAsASlice, check if the number exists in the combined slice
		for i := 1; i < len(fileAsSlice); i++ {
			indexOfLineInSlice := indexExistsInSlice(totalCommandsHandled, lastSeenIndex, i)

			if indexOfLineInSlice > -1 {
				lastSeenIndex = indexOfLineInSlice
			} else {
				//i did not exist in either the succeeded slice or the failed slice
				unfinishedCommandFileContents += fileAsSlice[i-1]
				if i < len(fileAsSlice)-1 {
					unfinishedCommandFileContents += "\n"
				}
			}
		}

		if len(unfinishedCommandFileContents) > len(fileAsSlice[0])+1 { //+1 for the newline
			var unfinishedFileName = fmt.Sprintf("unfinished_%s", fileName)
			fmt.Println(fmt.Sprintf("Writing %s...", unfinishedFileName))
			err = files.WriteFile(fmt.Sprintf("%s/%s", filePath, unfinishedFileName), unfinishedCommandFileContents)
			if err != nil {
				fmt.Printf("Error while writing unfinished csv [%s]\n", err.Error())
			}
		}
	}

	return err
}

func progressPrint(input string, lineCount int, commandSuccess int, commandFailure int, startTime time.Time) {
	bar := tool.CalculateProgress(commandSuccess, commandFailure, lineCount, startTime)
	fmt.Println("\r" + strings.Replace(input, "\n", "", -1) + strings.Repeat(" ", len(bar)))
	fmt.Print(bar)
}

func indexExistsInSlice(slice []int, startingIndex int, desiredElement int) int {
	for i := startingIndex; i < len(slice); i++ {
		if slice[i] == desiredElement {
			return i
		} else if slice[i] > desiredElement {
			return -1
		}
	}
	return -1
}

func printHelp() {
	fmt.Println("\nWelcome to the JIRA tool!")
	fmt.Println("Specify the path to your CSV file using the -path flag")
	fmt.Println("The username and password of the JIRA account must be specified with the -u and -p flags")
	fmt.Println("A JIRA address must be specified using the -url flag")
	fmt.Println("If you want to use a custom time format, specify it with a time template with the -time flag (default m/d/yy)")
	fmt.Println("The CSV file must have descriptor lines preceding the details of the tickets")
	fmt.Println("The first column of each descriptor line must have a command [update, delete, create]")
	fmt.Println("The second column of each descriptor line must have the field 'title' when updating/deleting a ticket")
	fmt.Println("In the first column of each data line, the project should precede the title")
	fmt.Println("In the update command, the fields you want to update should follow the title field")
	fmt.Println("The updatable fields are as follows:")
	fmt.Println("assets affected, assigned to, assignment group, cve references, cvss" +
		", description, due date, ip address, mac address, method of discovery, operating system, " +
		"priority, resolution date, scan errata, service ports, summary, ticket type, vendor references, " +
		"vulnerability title\n")
	fmt.Println("The fields may be specified in any order, and are all optional")
	fmt.Println("If updating a status, however, an assigned to must be specified")
	fmt.Println("In the delete command, no additional fields are required")
	fmt.Println("At any point, a new descriptor line can be written, modifying the command/fields of the tool")
	fmt.Println("Example.csv:")
	fmt.Println("update, title, hostname\nAegisDEV2, AegisDEV2-00000, newhost\nAegisDEV2, AegisDEV2-00001, changehost\nupdate, title, hostname, mac address\nAegisDEV2, AegisDEV2-00002, anotherhost, AB:BA:AB:BA\ndelete, title\nAegisDEV2, AegisDEV2-00003")
	fmt.Printf("Please note that in each data line uses AegisDEV2 as the project\n\n")
}

func removeNonASCIICharacters(s string) string {
	re := regexp.MustCompile("[[:^ascii:]]")
	return strings.TrimSpace(re.ReplaceAllLiteralString(s, ""))
}
