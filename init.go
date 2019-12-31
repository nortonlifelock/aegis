package funnel

import (
	"context"
	"github.com/benjivesterby/validator"
	"github.com/nortonlifelock/log"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

// New creates a new instance of the funnel for use with an api. New accepts a context, client interface which implements the Do method
// from the http.Client struct, as well as a logger, delay time duration for how often requests should be allwed as well as the number
// of retries each request should be allowed and the number of concurrent requests that should be allowed as part of the funnel.
// New returns an interface implemntation of Client which replaces the implementation of an http.Client interface so that it looks like
// an http.Client and can perform the same functions but it limits the requests using the parameters defined when created
func New(ctx context.Context, client Client, lstream log.Logger, delay time.Duration, retries int, concurrency int) (Client, error) {
	var err error
	var f *funnel

	if retries >= 0 {
		if lstream != nil {
			var ticker *time.Ticker
			if delay > 0 {
				ticker = time.NewTicker(delay)
			}

			// ensure the concurrency is setup above zero
			if concurrency > 0 {

				// Setup a background context if no context is passed
				if ctx == nil {
					ctx = context.Background()
				}

				// If a nil client is passed to the funnel then initialize a new http client
				if client == nil {
					client = &http.Client{}
				}

				f = &funnel{
					ctx:               ctx,
					client:            client,
					lstream:           lstream,
					retries:           retries,
					delay:             delay,
					ticker:            ticker,
					concurrency:       concurrency,
					concurrencyticker: make(chan bool, concurrency),
				}

				// Initialize the concurrency channel for managing concurrent calls
				for i := 0; i < f.concurrency; i++ {
					select {
					case <-f.ctx.Done():
					case f.concurrencyticker <- true:
					}
				}

				// Setup requests channel
				f.requests = f.receive()

				if !validator.IsValid(f) {
					err = errors.New("funnel is not valid after initialization")
				}
			} else {
				err = errors.Errorf("concurrency limit of [%v] is not valid; concurrency limit must be > 0", concurrency)
			}
		} else {
			err = errors.New("nil logger passed to funnel initialization")
		}
	} else {
		err = errors.Errorf("retries [%v] is invalid; retries must be greater than or equal to 0", retries)
	}

	return f, err
}
