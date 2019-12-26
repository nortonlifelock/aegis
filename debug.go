package jira

import (
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"github.com/nortonlifelock/domain"
	"strings"
)

// GetEditableFields is not used by the project, but is useful for debugging as it prints the editable JIRA fields
func (connector *ConnectorJira) GetEditableFields(id string) (err error) {
	var replace = "{issueIdOrKey}"
	var endPoint = html.EscapeString(jeditablefields)
	endPoint = strings.Replace(endPoint, replace, id, 1)

	var request *http.Request

	if request, err = connector.client.NewRequest(http.MethodGet, endPoint, nil); err == nil {

		var response *http.Response
		if response, err = connector.funnelClient.Do(request); err == nil {
			if response != nil {
				defer response.Body.Close()

				var body []byte
				if body, err = ioutil.ReadAll(response.Body); err == nil {
					fmt.Println(string(body))
				}
			}
		}
	}

	return err
}

// GetByCustomJQL returns the tickets that JIRA returns for a JQL statement. Not used by the application but is useful for testing
func (connector *ConnectorJira) GetByCustomJQL(JQL string) (tickets []domain.Ticket, err error) {
	var issues <-chan domain.Ticket
	if issues, err = connector.getByCustomJQL(JQL); err == nil {

		tickets = make([]domain.Ticket, 0)
		for {
			if ticket, ok := <-issues; ok {
				tickets = append(tickets, ticket)
			} else {
				break
			}
		}

	}
	return tickets, err
}
