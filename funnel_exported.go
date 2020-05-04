package funnel

import (
	"github.com/devnw/validator"
	"github.com/pkg/errors"
	"net/http"
)

// Do sends the http request through the funnel to be executed against the endpoint
// when there are available threads to do so. This returns an http response which
// is returned from the executation of the http request as well as an error
func (f *funnel) Do(request *http.Request) (response *http.Response, err error) {
	if request != nil {

		// Validate the funnel object and make sure it's properly functioning
		if validator.Valid(f) {

			var responsechan = make(chan responseWrapper)

			// Create the request wrapper to send to receive
			req := requestWrapper{
				request:  request,
				response: responsechan,
			}

			// setup a defer to close the channel
			defer close(req.response)

			// Send the request to the processing channel of the funnel
			go func() {
				defer f.handle()
				select {
				case <-f.ctx.Done():
					req.response <- responseWrapper{nil, errors.New("context closed for funnel")}
				case f.requests <- req:
				}
			}()

			// Wait for the response from the request
			if resp, ok := <-responsechan; ok {

				// Update the returns with the response from the server
				response = resp.response
				err = resp.err
			} else {
				err = errors.Errorf("response channel for request closed prematurely")
			}
		} else {
			err = errors.Errorf("invalid funnel instance, re-initialize properly and retry")
		}
	} else {
		err = errors.New("request cannot be nil")
	}

	return response, err
}

// Validate determines if the funnel is valid or not based on the
// internal state of the struct
func (f *funnel) Validate() (valid bool) {

	if f.concurrency > 0 {
		if f.retries >= 0 {
			if f.lstream != nil {
				if f.delay >= 0 {
					if f.ctx != nil {
						if f.concurrencyticker != nil {
							if f.requests != nil {
								valid = true
							}
						}
					}
				}
			}
		}
	}

	return valid
}
