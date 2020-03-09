package asdf

import (
	"context"
	"flag"
	"fmt"
	"github.com/nortonlifelock/aegis/internal/config"
	"github.com/nortonlifelock/aegis/internal/database"
	"github.com/nortonlifelock/aegis/internal/integrations"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/jira"
	"github.com/nortonlifelock/log"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"
)

var data domain.DatabaseConnection
var appConfig config.AppConfig

type copyTicket struct {
	domain.Ticket
}

func (e *copyTicket) ReportedBy() *string {
	return nil
}

func (e *copyTicket) ResolutionDate() *time.Time {
	return nil
}

func (e *copyTicket) ResolutionStatus() *string {
	return nil
}

func (e *copyTicket) AssignmentGroup() (param *string) {
	return nil
}

func (e *copyTicket) AssignedTo() (param *string) {
	return nil
}

func main() {
	if len(filePath) > 0 {
		copyTicketsByFile()
	} else if len(ticketJQL) > 0 {
		copyTicketsByJQL()
	} else {
		fmt.Println("Neither provided file nor JQL")
	}
}

func copyTicketsByJQL() {
	prod := getProdJiraSession()
	stage := getStageJiraSession()

	tics, err := prod.GetByCustomJQL(ticketJQL)
	check(err)

	var wg sync.WaitGroup
	for _, tic := range tics {

		wg.Add(1)
		go func(line domain.Ticket) {
			defer wg.Done()

			_, newKey, err := stage.CreateTicket(&copyTicket{tic})
			check(err)

			fmt.Printf("%s copied to %s\n", line, newKey)
		}(tic)
	}
	wg.Wait()
}

func copyTicketsByFile() {
	prod := getProdJiraSession()
	stage := getStageJiraSession()

	body, err := ioutil.ReadFile(filePath)
	check(err)

	var separator = "\r\n"
	fileContents := string(body)
	if strings.Index(fileContents, "\r\n") < 0 { //\r\n is not the separator
		if strings.Index(fileContents, "\r") >= 0 {
			fileContents = strings.Replace(fileContents, "\r", "\n", -1)
			separator = "\n"
		} else if strings.Index(fileContents, "\n") >= 0 {
			separator = "\n"
		} else {
			fmt.Printf("WARNING could not identify separator\n")
			os.Exit(0)
		}
	}

	lines := strings.Split(fileContents, separator)

	var wg sync.WaitGroup
	for _, line := range lines {

		wg.Add(1)
		go func(line string) {
			defer wg.Done()

			tic, err := prod.GetTicket(line)
			check(err)

			_, newKey, err := stage.CreateTicket(&copyTicket{tic})
			check(err)

			fmt.Printf("%s copied to %s\n", line, newKey)
		}(line)
	}
	wg.Wait()
}

func getProdJiraSession() *jira.ConnectorJira {
	sc, err := data.GetSourceConfigByID(prodJiraSourceID)

	if err == nil {
		sess, err := integrations.GetEngine(context.Background(), integrations.JIRA, data, Logger{}, appConfig, sc)
		if err == nil {
			if isJira, ok := sess.(*jira.ConnectorJira); ok {
				return isJira
			}
		} else {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println(err.Error())
	}

	panic("failed to create")
}

func getStageJiraSession() *jira.ConnectorJira {
	sc, err := data.GetSourceConfigByID(stageJiraSourceID)

	if err == nil {
		sess, err := integrations.GetEngine(context.Background(), integrations.JIRA, data, Logger{}, appConfig, sc)
		if err == nil {
			if isJira, ok := sess.(*jira.ConnectorJira); ok {
				return isJira
			}
		} else {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println(err.Error())
	}

	panic("failed to create")
}

func check(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

type Logger struct {
}

func (logger Logger) Send(log log.Log) {
	_ = logFunc("", log.Text, log.Error)
}

func logFunc(logType string, log string, logError error) (err error) {
	_ = logType
	var message = log
	if logError != nil {
		message = fmt.Sprintf("%s - %s", message, logError.Error())
	}
	fmt.Println(message)
	return err
}

var (
	stageJiraSourceID string
	prodJiraSourceID  string
	filePath          string
	ticketJQL         string
)

// TODO could also support JQLs
func init() {
	path := flag.String("path", "", "The path to (and including) your app.json")
	prodID := flag.String("prod", "", "The source config ID of your prod instance")
	stageID := flag.String("stage", "", "The source config ID of your stage instance")
	filePathP := flag.String("file", "", "The path to your newline delimited file containing the titles of tickets you'd like to copy")
	jql := flag.String("jql", "", "The JQL that will be used to copy tickets over")
	flag.Parse()

	if len(*prodID) == 0 || len(*stageID) == 0 {
		fmt.Println("Please provide the source config ID of your stage and prod instance with the -prod and -stage flag")
		os.Exit(1)
	}

	if len(*path) == 0 {
		fmt.Println("Please provide the path to (and including) your app.json with -path")
	}

	if len(*filePathP) == 0 && len(*jql) == 0 {
		fmt.Println("Please provide a either path to your newline delimited file containing the titles of tickets you'd like to copy with -file or a JQL to grab tickets to copy with -jql")
	}

	stageJiraSourceID, prodJiraSourceID, filePath, ticketJQL = *stageID, *prodID, *filePathP, *jql

	var err error
	appConfig, err = config.LoadConfigByPath(*path)
	if err == nil {
		data, err = database.NewConnection(appConfig)
		if err == nil {

		} else {
			panic(err)
		}
	} else {
		panic(err)
	}
}
