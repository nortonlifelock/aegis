package jira

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"

	"github.com/andygrunwald/go-jira"
	"github.com/nortonlifelock/aegis/pkg/domain"
	"github.com/nortonlifelock/aegis/pkg/log"
)

// Request contains the fields required to make a JIRA API call and process the result
type Request struct {
	req      *http.Request
	v        interface{}
	callback func(*jira.Response, error)
}

func (connector *ConnectorJira) buildPOSTSearchRequest(query *Query, startIndex int, max int) (req *postSearchRequest, err error) {
	if query != nil {

		if query.Size == 0 {
			query.Size = 1000
		}

		// Assign the Fields map to a string slice for passing to the jsearch API
		var fields []string

		req = &postSearchRequest{
			JQL:        query.JQL,
			StartAt:    startIndex,
			MaxResults: max,
			Fields:     fields,
		}
	} else {
		err = errors.New("query object cannot be nil")
	}

	return req, err
}

func (connector *ConnectorJira) getSearchResults(query *Query) <-chan domain.Ticket {
	results := make(chan domain.Ticket)

	go func(results chan<- domain.Ticket) {
		defer close(results)

		var err error
		var startIndex int
		var total int
		var ok bool

		if total, err = connector.getQueryTotal(query); err == nil {

			connector.lstream.Send(log.Infof("Ticket Total: %v", total))

			for startIndex < total {

				endIndex := startIndex + query.Size
				if endIndex > total {
					endIndex = total
				}

				connector.lstream.Send(log.Infof("Loading Tickets [%v - %v] out of %v", startIndex, endIndex, total))

				var issues interface{} = newSearchResult()
				if issues, err = connector.getSearchResponse(query, issues, startIndex); err == nil {
					var res *searchResult
					res, ok = issues.(*searchResult)
					if ok {
						for id := range res.Issues {
							results <- &Issue{Issue: res.Issues[id], connector: connector}
						}
					} else {
						connector.lstream.Send(log.Errorf(nil, "[%v] was NOT a JIRA issue", res))
					}
				} else {
					connector.lstream.Send(log.Error(fmt.Sprintf("error while gathering search results [%s]", query.JQL), err))
				}

				startIndex = startIndex + 1000
			}
		} else {
			connector.lstream.Send(log.Error(fmt.Sprintf("Error while gathering query total for request [%s]", query.JQL), err))
		}
	}(results)

	return results
}

func (connector *ConnectorJira) getQueryTotal(query *Query) (total int, err error) {
	var searchResult = newSearchResult()

	// Verify the proper inputs for executing a jira jsearch query
	if connector.client != nil {
		if query != nil {
			if len(query.JQL) > 0 {

				if connector.client != nil {
					var request *http.Request

					var req *postSearchRequest

					if req, err = connector.buildPOSTSearchRequest(query, 0, 1); err == nil {
						// Create a jira api request for the custom Fields in JIRA
						if request, err = connector.client.NewRequest(http.MethodPost, html.EscapeString(jsearch), req); err == nil {

							var response *http.Response
							if response, err = connector.funnelClient.Do(request); err == nil {
								if response != nil {
									defer response.Body.Close()

									var body []byte
									if body, err = ioutil.ReadAll(response.Body); err == nil {
										if err = json.Unmarshal(body, searchResult); err == nil {
											total = searchResult.Total
										}
									}
								}
							}
						}
					}

				} else {
					err = errors.New("JIRA Client is not authenticated")
				}

			} else {
				err = errors.New("empty JQL String")
			}
		} else {
			err = errors.New("query cannot be nil")
		}
	} else {
		err = errors.New("JIRA Client cannot be nil")
	}

	return total, err
}

// Accepts a client connection, a Query object and the object that will be returned when filled with data
func (connector *ConnectorJira) getSearchResponse(query *Query, ret interface{}, startIndex int) (interface{}, error) {

	var err error

	// Verify the proper inputs for executing a jira jsearch query
	if connector.client != nil {
		if query != nil {
			if len(query.JQL) > 0 {

				if connector.client != nil {
					var request *http.Request

					var req *postSearchRequest

					if req, err = connector.buildPOSTSearchRequest(query, startIndex, query.Size); err == nil {
						// Create a jira api request for the custom Fields in JIRA
						if request, err = connector.client.NewRequest(http.MethodPost, html.EscapeString(jsearch), req); err == nil {

							var response *http.Response
							if response, err = connector.funnelClient.Do(request); err == nil {
								if response != nil {
									defer response.Body.Close()

									var body []byte
									if body, err = ioutil.ReadAll(response.Body); err == nil {
										err = json.Unmarshal(body, ret)
									}
								}
							}
						}
					}

				} else {
					err = errors.New("JIRA Client is not authenticated")
				}

			} else {
				err = errors.New("empty JQL String")
			}
		} else {
			err = errors.New("query cannot be nil")
		}
	} else {
		err = errors.New("JIRA Client cannot be nil")
	}

	return ret, err
}

func newSearchResult() (result *searchResult) {
	result = &searchResult{
		StartAt:    0,
		MaxResults: 0,
		Issues:     make([]*jira.Issue, 0),
	}

	return result
}
