package funnel

import (
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"github.com/nortonlifelock/log"
	"strings"
	"time"
)

type httpclient struct {
	delay             time.Duration
	requests          int
	status            int
	retries           int
	attempts          int
	attemptedRequests int
	concurrency       int
	logger            log.Logger
	cancel            bool
}

func (client *httpclient) Do(r *http.Request) (*http.Response, error) {
	if client.delay > 0 {
		time.Sleep(client.delay)
	}

	var status = client.status
	if client.retries > 0 {
		client.attempts++

		if client.attempts < client.retries {
			status = http.StatusBadRequest
		}
	}

	return &http.Response{StatusCode: status}, nil
}

type badclient struct {
	panic             bool
	delay             time.Duration
	requests          int
	status            int
	retries           int
	attempts          int
	attemptedRequests int
	concurrency       int
	logger            log.Logger
}

func (client *badclient) Do(r *http.Request) (*http.Response, error) {
	if client.panic {
		panic("panic")
	}

	return nil, errors.New("error")
}

// Test logger for printing internal logs to the screen as part of the testing
type testlogger struct{}

func (t *testlogger) Send(l log.Log) {
	fmt.Println(l.Text)
}

type tstruct struct {
	error bool
}

func (t *tstruct) correct(errored, paniced bool) (err error) {
	if !paniced {

		if t.error {
			if !errored {
				err = errors.New("expected error but success instead")
			}
		} else if !t.error && errored {
			err = errors.New("expected success but errored instead")
		}
	} else {
		err = errors.New("unexpected panic")
	}

	return err
}

func newReq(method string, url string) *http.Request {
	r, _ := http.NewRequest(method, url, strings.NewReader(""))

	return r
}
