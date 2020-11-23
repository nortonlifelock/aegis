package tool

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/nortonlifelock/aegis/internal/database/dal"
	"github.com/nortonlifelock/aegis/pkg/crypto"
	"github.com/nortonlifelock/aegis/pkg/domain"
	"github.com/nortonlifelock/aegis/pkg/jira"
	"github.com/nortonlifelock/aegis/pkg/log"
	"github.com/pkg/errors"
)

const (
	// Update is a const for the update command
	Update = "update"

	// Delete is a const for the delete command
	Delete = "delete"

	// Create is a const for the create command
	Create = "create"

	// Unknown is a const for an unrecognized command
	Unknown = "unknown"
)

// TODO there should not be any global variables
var (
	// TimeLayout is the layout for time parsing in the case of due date, can be overwritten VIA flag
	TimeLayout = "1/2/06" //"2006-01-02T15:04:05.000Z"

	// TimeToWaitBetweenLines defines the duration the tool waits between lines as to not overload the JIRA api
	TimeToWaitBetweenLines = 300 * time.Millisecond
)

// JiraConnection is a single JIRAConnection that is shared between threads that use the same project
type JiraConnection struct {
	Connector      *jira.ConnectorJira
	Config         domain.SourceConfig
	connectionLock *sync.Mutex
	project        string
}

func calculateLineCount(fileAsSlice []string) (lineCount int) {
	for index := range fileAsSlice {
		var currentLine = fileAsSlice[index]
		if len(currentLine) > 1 {
			if strings.Index(currentLine, Update) != 0 && strings.Index(currentLine, Create) != 0 && strings.Index(currentLine, Delete) != 0 {
				lineCount++
			}
		}
	}

	return lineCount
}

// JiraToolPayload contains all the information required throughout the processing of a file
type JiraToolPayload struct {
	BlockWG             *sync.WaitGroup
	db                  domain.DatabaseConnection
	fileContents        string
	templateConfig      domain.SourceConfig
	LineNumOfSucceed    *[]int
	LineNumOfFailed     *[]int
	LineNumToDescLine   map[int]string
	CommandSuccess      *int
	CommandFailure      *int
	appconfig           jira.ConfigJira
	Separator           string
	progressPrint       func(string, int, int, int, time.Time)
	StartTime           time.Time
	FileAsSlice         []string
	mapLock             *sync.Mutex //for locking projectToConnection
	numLock             *sync.Mutex //locks the numsuccess and numfailed variables as multiple threads utilize them
	lineCount           int
	projectToConnection map[string]*JiraConnection //projects must have separate connections, this map returns a connection to the project, one is created if it does not exist yet
}

// MakePayload creates a JIRAPayload, and calculates minor helpful information (such as the line count)
func MakePayload(db domain.DatabaseConnection, fileContents string, templateConfig domain.SourceConfig, appconfig jira.ConfigJira, progressPrint func(string, int, int, int, time.Time)) *JiraToolPayload {
	payload := &JiraToolPayload{
		BlockWG:             &sync.WaitGroup{},
		db:                  db,
		fileContents:        fileContents,
		templateConfig:      templateConfig,
		appconfig:           appconfig,
		progressPrint:       progressPrint,
		mapLock:             &sync.Mutex{},
		numLock:             &sync.Mutex{},
		projectToConnection: make(map[string]*JiraConnection),
	}

	var lineNumOfFailed = make([]int, 0)
	var lineNumOfSucceed = make([]int, 0)
	var lineNumToDescLine = make(map[int]string)
	var commandSuccess = new(int)
	var commandFailure = new(int)
	*commandSuccess = 0
	*commandFailure = 0

	payload.LineNumOfFailed = &lineNumOfFailed
	payload.LineNumOfSucceed = &lineNumOfSucceed
	payload.LineNumToDescLine = lineNumToDescLine
	payload.CommandSuccess = commandSuccess
	payload.CommandFailure = commandFailure

	var separator = "\r\n"

	if strings.Index(fileContents, "\r\n") < 0 { //\r\n is not the separator
		if strings.Index(fileContents, "\r") >= 0 {
			payload.fileContents = strings.Replace(payload.fileContents, "\r", "\n", -1)
			separator = "\n"
		} else if strings.Index(fileContents, "\n") >= 0 {
			separator = "\n"
		} else {
			fmt.Printf("WARNING could not identify separator\n")
		}
	}

	payload.Separator = separator
	payload.FileAsSlice = strings.Split(payload.fileContents, payload.Separator)
	payload.lineCount = calculateLineCount(payload.FileAsSlice)

	return payload
}

// ProcessCSVContents parses a file and creates the API calls to update/create/delete a JIRA issue
func ProcessCSVContents(payload *JiraToolPayload) (err error) {
	payload.StartTime = time.Now()

	var lineNum = 0
	var status = Unknown

	//each index holds a value separated by value
	var currentLine []string

	//helps map which ticket field is stored in which position for the description lines
	var descToIndex map[string]int

	reader := csv.NewReader(strings.NewReader(payload.fileContents))
	reader.Comma = ','

	var count = 0

	for err == nil {
		count++
		if count%100 == 0 {
			payload.BlockWG.Wait()
			count = 0
		}

		currentLine, err = reader.Read()
		lineNum++

		if err != nil { //allows different groups of fields to be specified between commands
			if strings.Contains(err.Error(), "wrong number of fields in line") {
				err = nil
			}
		}

		if err == nil {
			if len(currentLine) > 0 {
				command := strings.ToLower(strings.TrimSpace(currentLine[0]))

				if command == Delete || command == Update || command == Create {
					descToIndex, err = descriptionLineToMap(currentLine, command)
					if err == nil {
						status = command
					} else {
						err = errors.Errorf("Error while creating map for %s command on line %d: %s\n", command, lineNum, err.Error())
					}
				} else if status == Unknown {
					err = errors.Errorf("unrecognized command on line %d: %s, must be either update, delete, or create", lineNum, command)
					//i want this to be fatal, don't overwrite
				} else {
					ProcessLine(payload, currentLine, lineNum, status, descToIndex)
				}
			}
		}
	}
	//wait for threads to finish being created
	payload.BlockWG.Wait()

	return err
}

// ProcessLine kicks off a thread to handle a single line from the JIRA change file
func ProcessLine(payload *JiraToolPayload, currentLine []string, lineNum int, status string, descToIndex map[string]int) {
	time.Sleep(TimeToWaitBetweenLines)
	payload.BlockWG.Add(1)
	go func(currentLine []string, lineNum int, status string, descToIndex map[string]int, lineCount int) {
		defer payload.BlockWG.Done()

		var project = strings.TrimSpace(currentLine[0])
		var connection *JiraConnection
		connection = gatherConnection(connection, payload, project, lineNum)

		var err error
		var comment string
		if status == Update {
			if len(descToIndex) == len(currentLine) && len(currentLine) > 0 {

				var updateTicket domain.Ticket
				updateTicket, comment, err = updateLineToTicket(status, currentLine, descToIndex, connection, false, lineCount, payload)
				if err == nil {
					payload.progressPrint(fmt.Sprintf("Updating %s...\n", updateTicket.Title()), lineCount, *payload.CommandSuccess, *payload.CommandFailure, payload.StartTime)
					threadedUpdate(payload, connection, updateTicket, lineNum, status, comment, descToIndex)
				} else {
					failedLine(fmt.Sprintf("Error while parsing line to ticket object on line %d: %s\n", lineNum, err.Error()), payload, descToIndex, lineNum, status)
					err = nil
				}

			} else {
				failedLine(fmt.Sprintf("was expecting %d arguments on line %d but got %d\n", len(descToIndex), lineNum, len(currentLine)), payload, descToIndex, lineNum, status)
			}
		} else if status == Delete {
			if len(descToIndex) == len(currentLine) && len(currentLine) > 1 {
				ticketTitle := strings.TrimSpace(currentLine[1])
				threadedDelete(payload, connection, ticketTitle, lineNum, status, descToIndex)
			} else {
				failedLine(fmt.Sprintf("delete rows should have 2 columns [project, title] one line %d\n", lineNum), payload, descToIndex, lineNum, status)
			}
		} else if status == Create {
			if len(descToIndex) == len(currentLine) && len(currentLine) > 0 {

				var updateTicket domain.Ticket
				updateTicket, _, err = updateLineToTicket(status, currentLine, descToIndex, connection, true, lineCount, payload)
				if err == nil {
					payload.progressPrint(fmt.Sprintf("Creating ticket...\n"), lineCount, *payload.CommandSuccess, *payload.CommandFailure, payload.StartTime) //should i show information here
					threadedCreate(connection, updateTicket, payload, lineNum, descToIndex)
				} else {
					failedLine(fmt.Sprintf("Error while parsing line to ticket object on line %d: %s\n", lineNum, err.Error()), payload, descToIndex, lineNum, status)
					err = nil
				}

			} else {
				failedLine(fmt.Sprintf("was expecting %d arguments on line %d but got %d\n", len(descToIndex), lineNum, len(currentLine)), payload, descToIndex, lineNum, status)
			}
		}
	}(currentLine, lineNum, status, descToIndex, payload.lineCount)
}

func threadedDelete(payload *JiraToolPayload, connection *JiraConnection, ticketTitle string, lineNum int, status string, descToIndex map[string]int) {
	var err error
	err = connection.Connector.DeleteTicket(ticketTitle)
	if err == nil {
		succeedingLine(fmt.Sprintf("Successfully deleted %s\n", ticketTitle), payload, lineNum)
	} else {
		failedLine(fmt.Sprintf("Failed to delete %s on line %d: %s\n", ticketTitle, lineNum, err.Error()), payload, descToIndex, lineNum, status)
		err = nil
	}
}

func threadedCreate(connection *JiraConnection, updateTicket domain.Ticket, payload *JiraToolPayload, lineNum int, descToIndex map[string]int) {
	var newTicketName string
	var err error

	_, newTicketName, err = connection.Connector.CreateTicket(updateTicket)
	if err == nil {
		succeedingLine(fmt.Sprintf("Successfully created %s\n", newTicketName), payload, lineNum)
	} else {
		var message = fmt.Sprintf("Error while creating ticket [%s] on line %d: %s\n", updateTicket.Title(), lineNum, err.Error())

		if strings.Contains(err.Error(), "issue type") {
			message += "\nIssue types: "
			index := 0
			for key := range connection.Connector.IssueTypes {

				message += key

				if index < len(connection.Connector.IssueTypes)-1 {
					message += ", "
				} else {
					message += "\n"
				}

				index++

			}
		}

		failedLine(message, payload, descToIndex, lineNum, Create)
	}
}

func threadedUpdate(payload *JiraToolPayload, connection *JiraConnection, updateTicket domain.Ticket, lineNum int, status string, comment string, descToIndex map[string]int) {
	var err error
	_, _, err = connection.Connector.UpdateTicket(updateTicket, comment)
	if err == nil {
		succeedingLine(fmt.Sprintf("Successfully updated %s\n", updateTicket.Title()), payload, lineNum)
	} else {
		failedLine(fmt.Sprintf("Error while updating ticket [%s] on line %d: %s\n", updateTicket.Title(), lineNum, err.Error()), payload, descToIndex, lineNum, status)
	}
}

type payloadConfig struct {
	domain.SourceConfig
	payload string
}

// Payload overrides the underlying Payload of the SourceConfig interface contained within the struct
func (p *payloadConfig) Payload() *string {
	return &p.payload
}

func gatherConnection(connection *JiraConnection, payload *JiraToolPayload, project string, lineNum int) *JiraConnection {
	payload.mapLock.Lock()
	connection = payload.projectToConnection[project]
	if payload.projectToConnection[project] == nil && len(project) > 0 {
		payload.projectToConnection[project] = &JiraConnection{}
		connection = payload.projectToConnection[project]
		connection.connectionLock = &sync.Mutex{}
		connection.project = project

		decryptedConfig, err := crypto.DecryptSourceConfig(payload.db, payload.templateConfig, payload.appconfig)
		if err == nil {
			connection.Config = decryptedConfig
		} else {
			panic(err)
		}

		payloadByte, _ := json.Marshal(sord(payload.templateConfig.Payload()))
		payloadString := strings.Replace(string(payloadByte), "\"{", "{", 1)
		payloadString = strings.Replace(payloadString, "}\"", "}", 1)
		payloadString = strings.Replace(payloadString, "\\\"", "\"", -1)
		payload.templateConfig = &payloadConfig{payload.templateConfig, payloadString}
		checkConnection(connection, payload.db, lineNum, payload.appconfig)
	}
	payload.mapLock.Unlock()
	return connection
}

func updateLineToTicket(command string, updateLine []string, descToIndex map[string]int, connection *JiraConnection, createTicket bool, lineCount int, payload *JiraToolPayload) (updateTicket domain.Ticket, comment string, err error) {
	if connection != nil {
		if connection.Connector != nil {
			for index := range updateLine {
				updateLine[index] = strings.TrimSpace(updateLine[index])
				updateLine[index] = strings.Replace(updateLine[index], "\\n", "\n", -1)
				updateLine[index] = strings.Replace(updateLine[index], "/n", "\n", -1)
			}

			var ticketTitle string
			if descToIndex["title"] > 0 { //title should never be 0, that's the project
				ticketTitle = updateLine[descToIndex["title"]]
			}

			if createTicket {
				// do nothing
			} else {
				updateTicket, err = grabTicketFromJIRA(connection, ticketTitle, payload, lineCount)
			}

			updateTicket, err = createTicketObjectIfUpdateCommand(updateTicket, command)

			if err == nil { //will i need to use this method on non-updates? if so the ticket might not exist...
				updateTicket = csvTicket{
					updateTicket,
					command,
					updateLine,
					descToIndex,
				}

				for key := range descToIndex {
					if key == "comment" {
						comment = updateLine[descToIndex["comment"]]
						break
					}
				}

			} else {
				err = fmt.Errorf("error while grabbing [%s]: %s", ticketTitle, err.Error())
			}
		} else {
			err = errors.New("the JIRA connector was passed nil to updateLineToTicket")
		}
	} else {
		err = errors.New("the JIRA connection was passed nil to updateLineToTicket")
	}

	return updateTicket, comment, err
}

func grabTicketFromJIRA(connection *JiraConnection, ticketTitle string, payload *JiraToolPayload, lineCount int) (updateTicket domain.Ticket, err error) {
	var retryCount = 15
	var attemptCount = 0
	updateTicket = nil
	//how many times should i retry finding a ticket?
	for attemptCount < retryCount && updateTicket == nil {

		updateTicket, err = connection.Connector.GetTicket(ticketTitle)
		attemptCount++
		//sometimes jira returns a ticket nil even though it exists, let's try again
		if err != nil {
			if strings.Contains(err.Error(), "Issue returned from JIRA is null") && attemptCount < retryCount {
				payload.progressPrint(fmt.Sprintf("JIRA couldn't find %s, reattempting...\n", ticketTitle), lineCount, *payload.CommandSuccess, *payload.CommandFailure, payload.StartTime)
				updateTicket = nil
				err = nil //set the err to nil so the next block can enter
				time.Sleep(time.Millisecond * 500)
			} else {
				//ticket grab failed for some other reason, don't reattempt
				break
			}
		} else if attemptCount > 1 {
			payload.progressPrint(fmt.Sprintf("Found %s on a reattempt\n", ticketTitle), lineCount, *payload.CommandSuccess, *payload.CommandFailure, payload.StartTime)
		}
	}
	return updateTicket, err
}

func createTicketObjectIfUpdateCommand(updateTicket domain.Ticket, command string) (domain.Ticket, error) {
	var err error

	if updateTicket == nil && command != Update {
		updateTicket = &dal.Ticket{}
	} else if updateTicket == nil && command == Update {
		err = fmt.Errorf("could not find ticket")
	}
	return updateTicket, err
}

func descriptionLineToMap(descriptionLine []string, command string) (descToIndex map[string]int, err error) {
	var seenStatus = false
	var seenAssignedTo = false
	var seenTitle = false
	descToIndex = make(map[string]int)

	if len(descriptionLine) > 0 {

		for index, description := range descriptionLine {

			if strings.TrimSpace(strings.ToLower(description)) == "status" {
				seenStatus = true
			} else if strings.TrimSpace(strings.ToLower(description)) == "assigned to" {
				seenAssignedTo = true
			} else if strings.TrimSpace(strings.ToLower(description)) == "assignee" {
				seenAssignedTo = true
			} else if strings.TrimSpace(strings.ToLower(description)) == "title" {
				seenTitle = true
			}

			if index == 0 {
				descToIndex["project"] = index
			} else {
				descToIndex[strings.ToLower(strings.TrimSpace(description))] = index
			}
		}

	} else {
		err = errors.New("slice given to descriptionLineToMap was empty")
	}

	// TODO: Re-evaluate this code
	if seenStatus && !seenAssignedTo {
		// err = errors.New("if including a status column for updates, there must also be an assigned to column")
		err = nil // <-- certain status transitions don't need an assigned to (e.g. Closed-Decomm)
	}

	if !seenTitle && command != Create {
		err = errors.New("update commands must include a title field")
	}

	return descToIndex, err
}

func mapToDescriptionLine(command string, descToIndex map[string]int, lineNum int, lineNumToDescLine map[int]string) {
	var returnVal = command
	for index := 0; index < len(descToIndex); index++ {
		for mapKey, mapValue := range descToIndex {
			if mapValue == index && mapValue > 0 { //mapValue > 0 ignores the project, which we don't want to include in a description line
				returnVal += ","
				returnVal += mapKey
			}
		}
	}

	lineNumToDescLine[lineNum] = returnVal
}

type logger struct{}

func (logger logger) Send(log log.Log) {
	text := log.Text
	if log.Error != nil {
		text += " " + log.Error.Error()
	}
	fmt.Println(text)
}

func checkConnection(connection *JiraConnection, db domain.DatabaseConnection, lineNum int, appconfig jira.ConfigJira) {
	var err error

	if connection.Connector == nil {
		var newConnector *jira.ConnectorJira
		newConnector, _, err = jira.NewJiraConnector(context.Background(), logger{}, connection.Config)
		if err == nil && newConnector != nil {

			//connectSuccess++
			connection.Connector = newConnector
		} else {

			if err != nil {
				fmt.Println(fmt.Sprintf("Error while creating ticketing connection: %s", err.Error()))
			} else {
				fmt.Println("Null ticketing engine returned, exiting...")
				err = errors.New("null ticketing engine returned")
			}
			//connectFailure++
			checkConnection(connection, db, lineNum, appconfig)
		}
	}
}

func succeedingLine(message string, payload *JiraToolPayload, lineNum int) {
	payload.progressPrint(message, payload.lineCount, *payload.CommandSuccess, *payload.CommandFailure, payload.StartTime)
	payload.numLock.Lock()
	*payload.CommandSuccess++
	*payload.LineNumOfSucceed = append(*payload.LineNumOfSucceed, lineNum)
	payload.numLock.Unlock()
}

func failedLine(message string, payload *JiraToolPayload, descToIndex map[string]int, lineNum int, status string) {
	payload.progressPrint(message, payload.lineCount, *payload.CommandSuccess, *payload.CommandFailure, payload.StartTime)
	payload.numLock.Lock()
	*payload.CommandFailure++
	mapToDescriptionLine(status, descToIndex, lineNum, payload.LineNumToDescLine)
	*payload.LineNumOfFailed = append(*payload.LineNumOfFailed, lineNum)
	payload.numLock.Unlock()
}

// CalculateProgress creates a progress bar that includes an estimate of the amount of time until termination
func CalculateProgress(commandSuccess int, commandFailure int, lineCount int, startTime time.Time) string {
	percentNumRatio := float64(commandSuccess+commandFailure) / float64(lineCount)
	percentNum := percentNumRatio * 100
	elapsed := float64(time.Now().Sub(startTime))
	total := time.Duration(elapsed / percentNumRatio)
	estimated := startTime.Add(total)
	left := estimated.Sub(time.Now()).Round(time.Second)
	percent := fmt.Sprintf("%.2f", percentNum) + "%"
	bar := fmt.Sprintf("%d/%d %s %s [S:%d|F:%d]",
		commandSuccess+commandFailure, lineCount, percent, left, commandSuccess, commandFailure)
	return bar
}
