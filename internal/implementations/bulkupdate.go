package implementations

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/nortonlifelock/aegis/interal/files"
	"github.com/nortonlifelock/aegis/internal/domain"
	"github.com/nortonlifelock/jira-tool/tool"
	"github.com/nortonlifelock/log"
	"io"
	"os"
	"sort"
	"strings"
	"time"
)

// BulkUpdateJob implements the job structure and holds the connection in order to contact the API with log information that should be relayed
// to the user who started the job
type BulkUpdateJob struct {
	Payload *BulkUpdatePayload

	websocketConnection *websocket.Conn
	messages            chan []byte

	id          string
	payloadJSON string
	ctx         context.Context
	db          domain.DatabaseConnection
	lstream     log.Logger
	appconfig   domain.Config
	config      domain.JobConfig
	insource    domain.SourceConfig
	outsource   domain.SourceConfig
}

// BulkUpdatePayload parses the information from the job history. It holds the files that are going to be ran by the bulk update job as well
// as the org that's running it/who started the job/the JIRA URL to use
type BulkUpdatePayload struct {
	Filenames           []string `json:"file"`
	UsernameOfRequester string   `json:"user"`
	ServiceURL          string   `json:"serviceURL"`
	OrgID               string   `json:"orgId"`
}

// buildPayload parses the information from the job history and parses it into the Payload object
func (job *BulkUpdateJob) buildPayload(pjson string) (err error) {
	job.Payload = &BulkUpdatePayload{}
	err = json.Unmarshal([]byte(pjson), job.Payload)
	if err == nil {
		if len(job.Payload.Filenames) == 0 {
			err = fmt.Errorf("at least one file must be included in Payload")
		} else if len(job.Payload.ServiceURL) == 0 {
			err = fmt.Errorf("service URL cannot be empty")
		}
	}
	return err
}

// BulkUpdateMessage holds information from the job execution that should be relayed to the user
type BulkUpdateMessage struct {
	User       string `json:"user,omitempty"`
	Success    string `json:"success,omitempty"`
	Error      string `json:"failure,omitempty"`
	Repeat     string `json:"repeat,omitempty"`
	ErrorCount int    `json:"repeat_count"`
	Progress   string `json:"progress"`
}

// Process pulls the files for the bulk update, establishes a connection with the user, and executes the changes against JIRA and relays relevant information back
// to the user
// TODO should we have a max number of files allowed per history?
func (job *BulkUpdateJob) Process(ctx context.Context, id string, appconfig domain.Config, db domain.DatabaseConnection, lstream log.Logger, payload string, jobConfig domain.JobConfig, inSource []domain.SourceConfig, outSource []domain.SourceConfig) (err error) {

	var ok bool
	if job.ctx, job.id, job.appconfig, job.db, job.lstream, job.payloadJSON, job.config, job.insource, job.outsource, ok = validInputs(ctx, id, appconfig, db, lstream, payload, jobConfig, inSource, outSource); ok {

		if err = job.buildPayload(job.payloadJSON); err == nil {

			job.createSocketChannel()
			messages := job.sendMessageToAPI(job.ctx)
			job.messages = messages

			var jiraPayload *tool.JiraToolPayload
			var filePath string
			jiraPayload, filePath, err = job.createJiraToolPayload()
			if err == nil {

				err = tool.ProcessCSVContents(jiraPayload)
				if err != nil {
					if err == io.EOF {
						err = nil
					} else {
						job.lstream.Send(log.Error("error while processing csv", err))
						job.sendMessageToUser(err.Error(), false)
					}
				}

				job.processFileAfterCompletion(jiraPayload.CommandFailure, *jiraPayload.LineNumOfFailed, jiraPayload.LineNumToDescLine, jiraPayload.FileAsSlice, jiraPayload.Separator, filePath, messages)
			} else {
				job.sendMessageToUser(fmt.Sprintf("error during job initialization"), false)
			}

			if err == nil {
				job.sendMessageToUser("Job completed successfully!", true)
			} else {
				job.sendMessageToUser(fmt.Sprintf("error during this processing - %s", err.Error()), false)
			}
		} else {
			err = fmt.Errorf("error while building payload - %s", err.Error())
		}
	} else {
		err = fmt.Errorf("input validation failed")
	}

	return err
}

func (job *BulkUpdateJob) createJiraToolPayload() (jiraPayload *tool.JiraToolPayload, filePath string, err error) {
	var appConfig = job.appconfig
	var combinedFileContents string
	combinedFileContents, filePath, err = job.loadBulkUpdateFiles()
	if err == nil {
		var configTemplate domain.SourceConfig
		configTemplate, err = job.db.GetSourceOauthByOrgURL(job.Payload.ServiceURL, job.Payload.OrgID)
		if err == nil {
			if configTemplate != nil {
				var separator string

				combinedFileContents, separator = job.identifySeparator(combinedFileContents)
				var fileAsSlice []string
				fileAsSlice, err = job.addRequesterFieldToFile(combinedFileContents, separator)

				if err == nil {
					combinedFileContents = strings.Join(fileAsSlice, separator)

					jiraPayload = tool.MakePayload(
						job.db,
						combinedFileContents,
						configTemplate,
						appConfig,
						job.uiProgressPrint,
					)

				}

			} else {
				job.lstream.Send(log.Errorf(nil, "could not find source config with address [%s] in org [%s]", job.Payload.ServiceURL, job.Payload.OrgID))
			}
		} else {
			var originalErr = err
			err = fmt.Errorf("could not find source config for URL [%s] in org [%s]", job.Payload.ServiceURL, job.Payload.OrgID)
			job.lstream.Send(log.Error(err.Error(), originalErr))
		}
	} else {
		job.lstream.Send(log.Error("error while loading bulk update files", err))
	}

	return jiraPayload, filePath, err
}

func (job *BulkUpdateJob) addRequesterFieldToFile(combinedFileContents string, separator string) ([]string, error) {
	var err error

	var fileAsSlice = strings.Split(combinedFileContents, separator)
	for index := range fileAsSlice {
		fileAsSlice[index] = strings.Replace(fileAsSlice[index], "\n", "", -1) //last line has trailing newline
		if appearsToBeDescriptionLine(fileAsSlice[index]) {
			if strings.Index(fileAsSlice[index], "requester") < 0 {
				fileAsSlice[index] += ",requester"
			} else {
				err = fmt.Errorf("CSV cannot provide a requester field")
				break
			}

		} else {
			if len(fileAsSlice[index]) > 1 {
				fileAsSlice[index] += fmt.Sprintf(",requested by: %s", job.Payload.UsernameOfRequester)
			}
		}
	}
	return fileAsSlice, err
}

func (job *BulkUpdateJob) loadBulkUpdateFiles() (combinedFileContents string, filePath string, err error) {
	filePath = fmt.Sprintf("%s/bulk_updates", job.appconfig.AegisPath())
	var finishedFilePath = fmt.Sprintf("%s/bulk_updates/finished", job.appconfig.AegisPath())

	for index := range job.Payload.Filenames {
		var currentFile = job.Payload.Filenames[index]
		var updatePath = fmt.Sprintf("%s/%s", filePath, currentFile)

		var fileContents string
		if fileContents, err = files.GetStringFromFile(updatePath); err == nil {
			combinedFileContents += fileContents

			err = os.Rename(updatePath, fmt.Sprintf("%s/%s", finishedFilePath, currentFile))
			if err != nil {
				job.lstream.Send(log.Errorf(err, "could not move update file [%s] to finished directory", currentFile))
			}
		} else {
			var originalErr = err
			err = fmt.Errorf("could not find update file %s", currentFile)
			job.lstream.Send(log.Error(err.Error(), originalErr))
		}
	}

	return combinedFileContents, filePath, err
}

func (job *BulkUpdateJob) identifySeparator(combinedFileContents string) (string, string) {
	var separator = "\r\n"

	if strings.Index(combinedFileContents, "\r\n") < 0 { //\r\n is not the separator
		if strings.Index(combinedFileContents, "\r") >= 0 {
			combinedFileContents = strings.Replace(combinedFileContents, "\r", "\n", -1)
			separator = "\n"
		} else if strings.Index(combinedFileContents, "\n") >= 0 {
			separator = "\n"
		} else {
			fmt.Printf("WARNING could not identify separator\n")
		}
	}
	return combinedFileContents, separator
}

func (job *BulkUpdateJob) processFileAfterCompletion(commandFailure *int, lineNumOfFailed []int, lineNumToDescLine map[int]string, fileAsSlice []string, separator string, filePath string, messages chan []byte) {
	var err error

	var fileName = time.Now().Format(time.RFC3339)
	if *commandFailure > 0 {
		sort.Ints(lineNumOfFailed)
		var failedFileContents = ""
		var recentDescLine = ""

		for _, failedLine := range lineNumOfFailed {

			//don't reprint the description line if the failed lines use the same description line
			if recentDescLine == lineNumToDescLine[failedLine] {
				var failedLine = fileAsSlice[failedLine-1]

				//remove the requester field (the trailing field)
				failedLine = failedLine[:strings.LastIndex(failedLine, ",")]

				failedFileContents = fmt.Sprintf("%s\n%s", failedFileContents, failedLine)
			} else {
				recentDescLine = lineNumToDescLine[failedLine]

				var failedLine = fileAsSlice[failedLine-1]
				failedLine = strings.Replace(failedLine, separator, "", -1)

				//remove the requester field (the trailing field)
				failedLine = failedLine[:strings.LastIndex(failedLine, ",")]

				// the LastIndex on the recentDescLine removed the requested by, which should not be in the failed file
				failedFileContents = fmt.Sprintf("%s\n%s\n%s", failedFileContents, recentDescLine[:strings.LastIndex(recentDescLine, ",")], failedLine)
			}

		}
		if failedFileContents[0] == '\n' {
			failedFileContents = failedFileContents[1:]
		}
		//var failedFileName = fmt.Sprintf("failed_%s", fileName)
		//job.lstream.Send(log.Infof("Writing %s...", failedFileName))
		//err = files.WriteFile(fmt.Sprintf("%s/%s.csv", filePath, failedFileName), failedFileContents)
		//if err == nil {
		//	message := BulkUpdateMessage{
		//		User:       job.Payload.UsernameOfRequester,
		//		Repeat:     failedFileContents,
		//		ErrorCount: *commandFailure,
		//	}
		//
		//	byteVal, byteErr := json.Marshal(message)
		//	if byteErr == nil {
		//		messages <- byteVal
		//	}
		//} else {
		//	job.lstream.Send(log.Infof("Error while writing failed csv [%s]", err.Error()))
		//}

		message := BulkUpdateMessage{
			User:       job.Payload.UsernameOfRequester,
			Repeat:     failedFileContents,
			ErrorCount: *commandFailure,
		}

		byteVal, byteErr := json.Marshal(message)
		if byteErr == nil {
			messages <- byteVal
		}
	} else {
		// Let the UI know that the file completed successfully without errors
		message := BulkUpdateMessage{
			User:       job.Payload.UsernameOfRequester,
			Repeat:     "NA",
			ErrorCount: 0,
		}

		byteVal, byteErr := json.Marshal(message)
		if byteErr == nil {
			messages <- byteVal
		}
	}
}

func (job *BulkUpdateJob) sendMessageToUser(messageToSend string, success bool) {
	var message BulkUpdateMessage

	messageToSend = job.formatJobMessageToUI(messageToSend)

	if success {
		message = BulkUpdateMessage{
			User:    job.Payload.UsernameOfRequester,
			Success: messageToSend,
		}
	} else {
		message = BulkUpdateMessage{
			User:  job.Payload.UsernameOfRequester,
			Error: messageToSend,
		}
	}

	byteVal, byteErr := json.Marshal(message)
	if byteErr == nil {
		job.messages <- byteVal
	}

}

func appearsToBeDescriptionLine(line string) bool {
	line = strings.ToLower(line)
	return strings.Index(line, tool.Update) == 0 || strings.Index(line, tool.Create) == 0 || strings.Index(line, tool.Delete) == 0
}

func (job *BulkUpdateJob) sendMessageToAPI(ctx context.Context) (messages chan []byte) {
	messages = make(chan []byte)
	alreadyClosed := false

	go func() {
		defer handleRoutinePanic(job.lstream)

		for {
			select {
			case <-ctx.Done():
				if !alreadyClosed {
					job.lstream.Send(log.Info("Closing API websocket"))
					close(messages)
					alreadyClosed = true
				}
			case message, ok := <-messages:
				if ok {
					if job.websocketConnection != nil {
						err := job.websocketConnection.WriteMessage(websocket.TextMessage, message)
						if err != nil {
							job.lstream.Send(log.Error("Error while writing websocket message", err))
						}
					}
				} else {
					return
				}
			}
		}
	}()

	return messages
}

func (job *BulkUpdateJob) formatJobMessageToUI(input string) string {
	return fmt.Sprintf("%s] %s", job.id, input)
}

func (job *BulkUpdateJob) uiProgressPrint(input string, lineCount int, commandSuccess int, commandFailure int, startTime time.Time) {
	message := BulkUpdateMessage{}
	message.User = job.Payload.UsernameOfRequester
	message.Success = job.formatJobMessageToUI(input)
	message.Progress = tool.CalculateProgress(commandSuccess, commandFailure, lineCount, startTime)

	byteVal, err := json.Marshal(message)
	if err == nil {
		job.messages <- byteVal
	}

	job.lstream.Send(log.Info(strings.Replace(input, "\n", "", -1)))
}

func (job *BulkUpdateJob) createSocketChannel() {
	var err error
	var url = fmt.Sprintf("ws://localhost:%d", job.appconfig.APIPort())
	var endpoint = "connect/InternalWebSocket"
	var combinedURL = fmt.Sprintf("%s/%s", url, endpoint)

	job.lstream.Send(log.Infof("connecting to %s", combinedURL))

	fin := make(chan bool)
	tic := time.Tick(10 * time.Second)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer handleRoutinePanic(job.lstream)

		var connection *websocket.Conn
		connection, _, err = websocket.DefaultDialer.DialContext(ctx, combinedURL, nil)
		if err == nil {
			job.websocketConnection = connection
		} else {
			job.lstream.Send(log.Error("dial err", err))
		}
		fin <- true
	}()

	select {
	case <-fin:
	case <-tic:
		cancel()
		job.lstream.Send(log.Errorf(err, "timeout while creating websocket connection"))
	}
}
