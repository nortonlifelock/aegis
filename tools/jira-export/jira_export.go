package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/files"
	"github.com/nortonlifelock/jira"
	"github.com/nortonlifelock/log"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
	"sync"
	"time"
)

const csvDate = "01-02-2006"
const filenameDate = "01-02-2006_1504MST"

type exportLogger struct{}

// Send prints a log to console
func (logger exportLogger) Send(log log.Log) {
	if log.Error != nil {
		log.Text = fmt.Sprintf("%s - %s", log.Text, log.Error.Error())
	}

	fmt.Println(log.Text)
}

func main() {

	var err error

	// Setting up config arguments for starting the job runner
	outputPath := flag.String("o", "", "The output file for the resulting file")
	apiPath := flag.String("url", "", "The output file for the resulting file")
	user := flag.String("user", "", "The output file for the resulting file")

	// GetByCustomJQLChan

	flag.Parse()

	if outputPath != nil {
		var outputFile = *outputPath

		if len(outputFile) == 0 {
			outputFile, err = os.Getwd()

			outputFile = fmt.Sprintf("%s%sJiraExport_%s.csv", outputFile, string(os.PathSeparator), time.Now().Format(filenameDate))
		}

		if apiPath != nil && len(*apiPath) > 0 {
			if user != nil && len(*user) > 0 {

				var password []byte
				fmt.Println("Enter JIRA Password: ")

				if password, err = terminal.ReadPassword(0); err == nil {

					var tickets *jira.ConnectorJira
					if tickets, err = jira.ConnectJira(*apiPath, *user, string(password), exportLogger{}); err == nil {

						var JQL string

						reader := bufio.NewReader(os.Stdin)
						fmt.Print("Enter JQL: ")
						if JQL, err = reader.ReadString('\n'); err == nil {
							JQL = strings.TrimSpace(JQL)

							// TODO: Remove this once the closed-error status is part of exceptions
							//statuses["Closed-Error"] = true
							var wg = sync.WaitGroup{}
							var sm = sync.Map{}

							var tix = tickets.GetByCustomJQLChan(JQL)

							dump(outputFile, fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%v,%s,%s,%s,%s,%s,%s,%s,%s\n",
								"KEY",
								"SUMMARY",
								"ASSIGNMENTGROUP",
								"STATUS",
								"METHODOFDISCOVERY",
								"ORG",
								"CREATED",
								"ALERTDATE",
								"DUE",
								"RESOLUTIONDATE",
								"PRIORITY",
								"CVSS",
								"CERF",
								"PORTS",
								"OS",
								"HOSTNAME",
								"IPADDRESS",
								"DEVICEID",
								"VULNERABILITY",
								"VULNERABILITYID"))

							var tickets = make([]domain.Ticket, 0)
							var ticketLock = &sync.Mutex{}

							for {
								if issue, ok := <-tix; ok {
									wg.Add(1)

									go func(ticket domain.Ticket) {
										defer wg.Done()
										saveUniqueTicket(ticket, &sm, ticketLock, &tickets)
									}(issue)

								} else {
									break
								}
							}

							wg.Wait()

							fmt.Printf("Writing %d tickets\n", len(tickets))

							for index := range tickets {
								var ticket = tickets[index]

								var resolved string
								if ticket.ResolutionDate() != nil {
									resolved = tord(ticket.ResolutionDate()).Format(csvDate)
								}

								dump(outputFile, fmt.Sprintf("%s,\"%s\",%s,%s,%s,%s,%s,%s,%s,%s,%s,%v,%s,%s,%s,\"%s\",%s,%s,\"%s\",%v\n",
									ticket.Title(),
									sord(ticket.Summary()),
									sord(ticket.AssignmentGroup()),
									sord(ticket.Status()),
									sord(ticket.MethodOfDiscovery()),
									sord(ticket.OrgCode()),
									tord(ticket.CreatedDate()).Format(csvDate),
									tord(ticket.AlertDate()).Format(csvDate),
									tord(ticket.DueDate()).Format(csvDate),
									resolved,
									sord(ticket.Priority()),
									ford(ticket.CVSS()),
									ticket.CERF(),
									sord(ticket.ServicePorts()),
									sord(ticket.OperatingSystem()),
									sord(ticket.HostName()),
									sord(ticket.IPAddress()),
									ticket.DeviceID(),
									sord(ticket.VulnerabilityTitle()),
									ticket.VulnerabilityID()))
							}
						}
					} else {
						fmt.Printf("Error opening JIRA connection - check your username and password [%s]\n", err.Error())
					}
				} else {
					fmt.Println("Error while reading password from terminal")
				}
			} else {
				fmt.Println("Username (-user) is required")
			}
		} else {
			fmt.Println("API Path (-url) is required")
		}
	}

	fmt.Println("Export Complete")
}

func saveUniqueTicket(ticket domain.Ticket, sm *sync.Map, ticketLock *sync.Mutex, tickets *[]domain.Ticket) {
	if ticket != nil {
		if key, ok := sm.Load(ticket.Title()); ok {
			fmt.Printf("TICKET %s Already Exists\n", key)
		} else {

			sm.Store(ticket.Title(), ticket.Title())

			ticketLock.Lock()
			*tickets = append(*tickets, ticket)
			ticketLock.Unlock()
		}
	} else {
		fmt.Println("nil issue pulled off channel")
	}
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
			if err != nil {
				fmt.Println("ERROR WRITING - ", err.Error())
			}
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
