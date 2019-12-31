package funnel

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/nortonlifelock/log"
	"github.com/pkg/errors"
)

type funnel struct {
	ctx               context.Context
	client            Client
	retries           int
	ticker            *time.Ticker
	delay             time.Duration
	concurrency       int
	concurrencyticker chan bool
	requests          chan<- requestWrapper
	lstream           log.Logger
}

// Receive creates a new listening channel for http request wrappers. After creating the request channel it then
// monitors the delay timer (aka ticker) for each tick then checks for an available concurrency entry on the concurrency
// channel to process work. Once it's cleared the ticker and concurrency channel it then pulls an available request
// from the request channel and executes the http request against the endpoint and returns the response across the
// response channel of the request along with any errors that occurred when making the request
func (f *funnel) receive() chan<- requestWrapper {

	var reqs = make(chan requestWrapper)

	go func(reqs chan requestWrapper) {
		defer f.handle()
		defer func() {

			// Re-initialize unless the context is done
			select {
			case <-f.ctx.Done():
			default:
				// TODO: make sure that this doesn't cause a send to be lost in the event of a panic
				f.requests = f.receive()
			}
		}()
		defer func() {
			if f.ticker != nil {
				f.ticker.Stop()
			}
		}()

		f.processRequestWhenAppropriate(reqs)
	}(reqs)

	return reqs
}

func (f *funnel) processRequestWhenAppropriate(reqs chan requestWrapper) {
	if f.ticker != nil {
		for {
			select {
			case <-f.ctx.Done():
				return
			case _, ok := <-f.ticker.C:
				if ok {

					select {
					case <-f.ctx.Done():
						return
					case _, ok = <-f.concurrencyticker:
						if ok {
							f.process(reqs)
						} else {
							err := errors.New("concurrency ticker is closed, exiting")
							f.lstream.Send(log.Error(err.Error(), err))
							return
						}
					}
				} else {
					err := errors.New("delay ticker is closed, exiting")
					f.lstream.Send(log.Error(err.Error(), err))
					return
				}
			}
		}
	} else {
		for {

			select {
			case <-f.ctx.Done():
				return
			case _, ok := <-f.concurrencyticker:
				if ok {
					f.process(reqs)
				} else {
					err := errors.New("concurrency ticker is closed, exiting")
					f.lstream.Send(log.Error(err.Error(), err))
					return
				}
			}
		}
	}
}

func (f *funnel) process(requests chan requestWrapper) {
	defer f.handle()

	// Must have the context done in the select here otherwise if the ctx is closed
	// then this will cause a panic because of sending on a closed channel
	defer func() {
		select {
		case <-f.ctx.Done():
		case f.concurrencyticker <- true:
		}
	}()

	select {
	case <-f.ctx.Done():
	case req, ok := <-requests:
		if ok {
			if req.request != nil {
				go f.handleRequest(req, requests)
			} else {
				// Handle the nil request path
				select {
				case <-f.ctx.Done():
				case req.response <- responseWrapper{nil, errors.New("request is nil, cannot process nil request")}:
				}
			}
		} else {
			err := errors.New("request channel in the funnel closed")
			f.lstream.Send(log.Error(err.Error(), err))
		}
	}
}

func (f *funnel) handleRequest(req requestWrapper, requests chan requestWrapper) {
	// increment the attempt counter
	req.attempts++
	// Execute a call against the endpoint handling any potential panics from the http client
	resp, err := func() (resp *http.Response, err error) {
		defer func() {
			if panicMessage := recover(); panicMessage != nil {
				err = errors.New("panic occurred while executing http request")
			}
		}()

		// Execute the http request and return the response to the requester
		resp, err = f.client.Do(req.request)

		return resp, err
	}()
	// If the request is successful or the retries have run out then return the response, otherwise
	// push the request back onto the request channel until the retries counter has been exceeded or
	// the request has been successful
	if (err == nil && resp != nil && resp.StatusCode < 300) || f.retries == 0 || (f.retries > 0 && req.attempts >= f.retries) {

		// Add an error for the number of retries being exce
		if err == nil {
			err = f.handleFailedRequest(resp)
		}
		req.response <- responseWrapper{resp, err}
	} else {
		if resp != nil {
			if resp.StatusCode != http.StatusNotFound {
				f.lstream.Send(log.Warningf(nil, "code: %v - retrying: %s", resp.StatusCode, req.request.URL))

				// Send the request back on the channel
				go func() {
					select {
					case <-f.ctx.Done():
						req.response <- responseWrapper{nil, errors.New("context closed for funnel")}
					case requests <- req:
					}
				}()
			} else {
				req.response <- responseWrapper{
					response: resp,
					err:      fmt.Errorf("request returned a 404 response"),
				}
			}
		} else {
			req.response <- responseWrapper{
				response: resp,
				err:      fmt.Errorf("request returned a null response"),
			}
		}
	}
}

func (f *funnel) handleFailedRequest(resp *http.Response) (err error) {
	if resp != nil && resp.StatusCode >= 300 {
		err = errors.Errorf("retries exceeded for request | code: %v", resp.StatusCode)

		if resp.Body != nil {
			defer resp.Body.Close()
			if body, readErr := ioutil.ReadAll(resp.Body); readErr == nil {
				err = fmt.Errorf("%v - %v", err.Error(), string(body))
			}
		}
	}
	return err
}

// handle panics that occur, this method is called in a defer
func (f *funnel) handle() {
	if panicMessage := recover(); panicMessage != nil {
		if f.lstream != nil {
			var err = errors.Errorf("panic occurred in funnel instance [%v]", panicMessage)
			f.lstream.Send(log.Criticalf(err, err.Error()))
		}
	}
}

// wrapper for transporting requests along a channel along with a response channel for returning
// the response from the endpoint as well as an attempt counter for tracking the number of times
// a request has been attempted in the event that it continues to fail
type requestWrapper struct {
	request  *http.Request
	response chan<- responseWrapper
	attempts int
}

// wrapper for tracking the response of execuing a client.Do against an http request. This returns any errors
// from the funnel attempting to execute the request as well as the http response in the event of a response
type responseWrapper struct {
	response *http.Response
	err      error
}
