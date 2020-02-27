package qualys

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/log"
	"github.com/pkg/errors"
)

//----------------------------------------------------------------
// Rate Limiting
//----------------------------------------------------------------

// CData allows text to be parsed within the XML
type CData struct {
	Text string `xml:",cdata"`
}

// Rates contains the information for Qualys rate limiting which is pulled from the headers after
// each qualys request
type Rates struct {

	// RLLimit is the total Rate Limit
	RLLimit int

	// RLWindowSec is the total Time In Seconds which you have to wait between rate limit reset
	RLWindowSec int

	// RLRemaining is the total number of requests that are allowed given the number of requests against this account
	RLRemaining int

	// RLToWaitSec is the total wait period in seconds before you can make the next call without being blocked
	RLToWaitSec int

	// CLimit is the total Concurrency limit
	CLimit int

	// CRunning is the current concurrent threads
	CRunning int
}

// makeRequest creates a http request to the Qualys API while also taking into account the rate limiting implementation
// by Qualys so that the requesting methods don't crash immediately
func (session *Session) makeRequest(request *http.Request, action func(response *http.Response) (err error)) (err error) {

	var status int
	var timeout int

	// Loop over the request until the status is 200 or there is an error returned from the loop
	for status != 200 && err == nil {

		// Ensure the timeout hasn't been met before making the request, otherwise break the loop
		if timeout < 500 {

			var authInfo domain.BasicAuth
			if err = json.Unmarshal([]byte(session.Config.AuthInfo()), &authInfo); err == nil {
				// Set the basic auth information and required Qualys Headers
				request.SetBasicAuth(authInfo.Username, authInfo.Password)
				request.Header.Add("X-Requested-With", "Aegis")
				request.Header.Add("Content-Type", "application/xml")

				// Execute the HTTP request against the Qualys API
				var response *http.Response
				if response, err = session.getFunnelForEndpoint(request.URL.Path).Do(request); err == nil {

					// Ensure that the response from the API was not NIL
					if response != nil {

						// Pull the response headers and validate that the rate limits have not been exceeded
						var rates Rates
						if status, rates, err = session.pullResponseHeaders(response); err == nil {
							err = session.processResponse(status, action, response, request, rates, timeout)
						} else {
							err = errors.Errorf("Error occurred when pulling rate limit information from Qualys response header [Endpoint: %s | Error: %s]", request.URL, err.Error())
						}
					} else {
						err = errors.Errorf("No Response from Qualys API [Endpoint: %s]", request.URL)
					}
				} else {
					err = errors.Errorf("Error when making API Request to Qualys [Endpoint: %s | Error: %s]", request.URL, err.Error())
				}
			} else {
				err = fmt.Errorf("error while parsing authentication information - %s", err.Error())
			}

		} else {
			session.lstream.Send(log.Warningf(err, "Qualys Timeout Reached for API [%s]", request.URL))
			break
		}
	}

	return err
}

func (session *Session) processResponse(status int, action func(response *http.Response) (err error), response *http.Response, request *http.Request, rates Rates, timeout int) (err error) {
	// If the status is 200 then the API request went through correctly
	switch status {
	case 200:
		// If the action function literal was passed to the makeRequest method then execute the action method
		if action != nil {
			err = action(response)
			_ = response.Body.Close()
		}
	default: // A non-rate limiting error was returned from the API, handle it appropriately
		var data []byte
		defer response.Body.Close()
		data, err = ioutil.ReadAll(response.Body)

		// Error code that we're not monitoring
		err = errors.Errorf("Qualys API [%s] Error Returned [Status: %v | Error: %s]", request.URL, status, string(data))
	}
	return err
}

// pullResponseHeaders reads the response headers from the response object and parses them into the rate limit
// properties specified in the Qualys documentation
func (session *Session) pullResponseHeaders(response *http.Response) (status int, rates Rates, err error) {

	if response != nil {

		// Get the response code from the error
		status = response.StatusCode

		var rllimit = response.Header.Get("x-ratelimit-limit")
		if len(rllimit) > 0 {
			if rates.RLLimit, err = strconv.Atoi(rllimit); err == nil {

				var rlwindowsec = response.Header.Get("x-ratelimit-window-sec")

				if len(rlwindowsec) > 0 {
					if rates.RLWindowSec, err = strconv.Atoi(rlwindowsec); err == nil {

						var rlremaining = response.Header.Get("x-ratelimit-remaining")

						if len(rlremaining) > 0 {
							if rates.RLRemaining, err = strconv.Atoi(rlremaining); err == nil {

								var rltowaitsec = response.Header.Get("x-ratelimit-towait-sec")

								if len(rltowaitsec) > 0 {
									if rates.RLToWaitSec, err = strconv.Atoi(rltowaitsec); err == nil {

										var climit = response.Header.Get("x-concurrency-limit-limit")

										if len(climit) > 0 {
											if rates.CLimit, err = strconv.Atoi(climit); err == nil {

												var crunning = response.Header.Get("x-concurrency-limit-running")

												if len(crunning) > 0 {
													if rates.CRunning, err = strconv.Atoi(crunning); err == nil {

													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return status, rates, err
}

//----------------------------------------------------------------
// HTTP Methods
//----------------------------------------------------------------

// Execute a POST call against the Qualys API
func (session *Session) post(path string, fields map[string]string, obj interface{}) (err error) {
	err = session.httpCall(http.MethodPost, path, fields, nil, obj)
	return err
}

// Execute a POST call against the Qualys API that contains a binary
func (session *Session) httpCall(method string, path string, fields map[string]string, in *string, obj interface{}) (err error) {

	var qstring string
	if qstring, err = mapToQueryString(path, fields); err == nil {

		var reader *strings.Reader
		if in == nil {
			reader = strings.NewReader("")
		} else {
			reader = strings.NewReader(*in)
		}

		var request *http.Request
		if request, err = http.NewRequest(method, fmt.Sprintf("%s%s", path, qstring), reader); err == nil {

			err = session.makeRequest(request, func(response *http.Response) (err error) {
				if response != nil {
					defer response.Body.Close()

					var data []byte
					if data, err = ioutil.ReadAll(response.Body); err == nil {
						//fmt.Println(string(data))
						ioutil.WriteFile("hostlist.xml", data, 0444)
						if retResponse, ok := obj.(*simpleReturn); ok {

							if err = xml.Unmarshal(data, &retResponse); err == nil {

								if retResponse.Response.Code > 0 {
									err = fmt.Errorf("error While Accessing Qualys. Error [%v]: %s | URL [%s] | DATA [%s]", retResponse.Response.Code, retResponse.Response.Message, path, string(data))
								}
							}
						} else {

							if err = xml.Unmarshal(data, obj); err != nil {
								var e1 = err

								var retResponse simpleReturn
								if err = xml.Unmarshal(data, &retResponse); err == nil {

									if retResponse.Response.Code > 0 {
										err = fmt.Errorf("error While Accessing Qualys: [%s] Error [%v]: %s |  URL [%s] | DATA [%s]", e1, retResponse.Response.Code, retResponse.Response.Message, path, string(data))
									}
								} else {
									err = fmt.Errorf("error While Accessing Qualys: [%s] Error attempting to unmarshal simple return. URL [%s] | DATA [%s]", e1, path, string(data))
								}
							}
						}
					}
				}

				return err
			})
		}
	}

	return err
}

func mapToQueryString(path string, fields map[string]string) (qstring string, err error) {

	if fields != nil {
		if len(fields) > 0 {

			var first bool
			if !strings.Contains(path, "?") {
				first = true
				qstring = "?"
			}

			for k, v := range fields {
				var separator = ""

				if first {
					first = false
				} else {
					separator = "&"
				}

				qstring = fmt.Sprintf("%s%s%s=%s", qstring, separator, k, url.QueryEscape(v))
			}
		}
	} else {
		err = errors.New("Invalid list of fields passed to map to query string")
	}

	return qstring, err
}
