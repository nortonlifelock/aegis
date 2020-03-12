package nexpose

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/benjivesterby/validator"
	"github.com/nortonlifelock/log"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

type resource interface {
	Resources(ctx context.Context, out chan<- interface{})
	TotalPages() int
}

type client interface {
	Do(r *http.Request) (*http.Response, error)
}

type auth interface {
	User() string
	Pass() string
}

type host interface {
	auth
	Host() string
	Port() int
}

// Session is the base api struct for managing the api calls to nexpose
type Session struct {
	ctx     context.Context
	c       client
	host    host
	lstream log.Logger
}

// Connect creates a new Session connection to the nexpose API that implements the interface for scanner abstraction
func Connect(ctx context.Context, client client, host host, lstream log.Logger) (nexpose *Session, err error) {

	if ctx == nil {
		ctx = context.Background()
	}

	if client != nil {

		if host != nil {

			nexpose = &Session{
				c:       client,
				ctx:     ctx,
				host:    host,
				lstream: lstream,
			}
		} else {
			err = fmt.Errorf("nil host passed to Nexpose client creation")
		}
	} else {
		err = fmt.Errorf("nil client passed to Nexpose client creation")
	}

	return nexpose, err
}

// Validate ensures that the Session connection is valid so that it can properly reach out to the
// nexpose Session and that the connection is able to log, etc...
func (a *Session) Validate() (valid bool) {

	if a.ctx != nil &&
		len(a.host.Host()) > 0 &&
		a.host.Port() > 0 &&
		len(a.host.User()) > 0 &&
		len(a.host.Pass()) > 0 &&
		a.lstream != nil {
		valid = true
	}

	return valid
}

// getPagedData handles the paging of nexpose api calls and returns the data back over an interface channel as each page
// load is completed
func (a *Session) getPagedData(ctx context.Context, binding resource, url string, qfields map[string]string) (<-chan interface{}, error) {

	var data = make(chan interface{})
	var err error

	if validator.IsValid(a) {

		go func(data chan<- interface{}) {
			defer handleRoutinePanic(a.lstream)
			defer close(data)
			var wg = sync.WaitGroup{}

			a.page(func(fields map[string]string) (totalPages int) {

				// Add the additional non-paging fields
				for k, v := range qfields {
					if k != "size" && k != "page" {
						fields[k] = v
					}
				}

				select {
				case <-ctx.Done():
					return
				default:

					// Create a new pointer instance of the object to be mapped
					page := reflect.New(reflect.TypeOf(binding).Elem()).Interface().(resource)

					if err = a.execute(http.MethodGet, url, fields, nil, page); err == nil {
						totalPages = page.TotalPages()

						wg.Add(1)
						// Process the page
						go func() {
							defer handleRoutinePanic(a.lstream)
							defer wg.Done()

							// Return the resources over the channel
							page.Resources(ctx, data)
						}()
					} else {
						a.lstream.Send(log.Errorf(err, "error occurred while making request to endpoint [%s]", url))
					}
				}

				return totalPages
			})

			// Wait for the rest of the data to be loaded
			wg.Wait()
		}(data)
	} else {
		err = errors.New("invalid nexpose Session instance")
	}

	return data, err
}

// page manages the page counts and executions of the endpoint calls
func (a *Session) page(execute func(fields map[string]string) (totalPages int)) {
	var err error
	var wg = sync.WaitGroup{}

	// Setup paging variables
	currentPage := 0
	totalPages := 0

	// Execute the first page for the total number of pages
	totalPages = execute(map[string]string{"page": strconv.Itoa(currentPage), "size": strconv.Itoa(pageSize)})
	currentPage++

	for ((currentPage == 0 && totalPages == 0) || currentPage < totalPages) && err == nil {

		wg.Add(1)
		go func(page int, size int) {
			defer handleRoutinePanic(a.lstream)
			defer wg.Done()
			execute(map[string]string{"page": strconv.Itoa(page), "size": strconv.Itoa(size)})
		}(currentPage, pageSize)

		// Increment the page
		currentPage++
		time.Sleep(time.Second * 3)
	}

	wg.Wait()

	// Log any errors encountered as part of executing a paged request
	if err != nil {
		a.lstream.Send(log.Errorf(err, "error encountered when paging api request [page %v, size %v]", currentPage, pageSize))
	}
}

// execute the request passed to the execute method against the endpoint specificed and
// unmarshal the response body to the interface passed in
func (a *Session) execute(method string, url string, fields map[string]string, body io.Reader, v interface{}) (err error) {
	url = a.mapQString(url, fields)

	// Create a new http request
	var request *http.Request
	if request, err = a.newRequest(method, url, body); err == nil {

		// Execute the request against the API
		var response *http.Response
		if response, err = a.c.Do(request); err == nil {
			if response != nil {
				defer response.Body.Close()

				// Read the response body from the response
				var data []byte
				if data, err = ioutil.ReadAll(response.Body); err == nil {

					// only unmarshal if the status code is in the success range
					if response.StatusCode >= http.StatusOK && response.StatusCode < 300 {

						// Ensure the value to unmarshal to is not nil
						if v != nil {
							// unmarshal the response to the asset variable
							err = json.Unmarshal(data, v)
						}
					} else {
						err = errors.Errorf("http response code [%v] returned from nexpose [%s]", response.StatusCode, data)
					}
				}
			}
		}
	}

	return err
}

// mapQString takes in a map of strings and creates a query string from them
func (a *Session) mapQString(path string, fields map[string]string) (qstring string) {

	if fields != nil && len(fields) > 0 {

		var first bool
		if !strings.Contains(path, "?") {
			first = true
			qstring = "?"
		}

		for k, v := range fields {
			k = encode(k)
			v = encode(v)

			var separator = ""

			if first {
				first = false
			} else {
				separator = "&"
			}

			qstring = fmt.Sprintf("%s%s%s=%s", qstring, separator, k, v)
		}

		qstring = fmt.Sprintf("%s%s", path, qstring)
	} else {
		qstring = path
	}

	return qstring
}

// newRequest creates an Session request and establishes basic auth on the request using the internal credentials
func (a *Session) newRequest(method string, url string, body io.Reader) (request *http.Request, err error) {

	// Update the url with the full path and configured endpoint
	url = fmt.Sprintf(apiEndpoint, a.host.Host(), a.host.Port(), url)

	if request, err = http.NewRequest(method, url, body); err == nil {
		request.SetBasicAuth(a.host.User(), a.host.Pass())

		// Set the content type of the request
		request.Header.Set("Content-Type", "application/json")
	}

	return request, err
}
