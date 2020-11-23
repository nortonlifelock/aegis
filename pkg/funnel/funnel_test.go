package funnel

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/devnw/validator"
)

func TestFunnel_Do(t *testing.T) {

	tests := []struct {
		name    string
		client  *httpclient
		request *http.Request
		success tstruct
	}{
		{
			"ValidWValidClient",
			&httpclient{
				0,
				1,
				http.StatusOK,
				0,
				0,
				0,
				1,
				&testlogger{},
				false,
			},
			newReq(http.MethodGet, "ValidWValidClient"),
			tstruct{false},
		},
		{
			"ValidWValidClientConcurrency",
			&httpclient{
				0,
				1,
				http.StatusOK,
				0,
				0,
				0,
				10,
				&testlogger{},
				false,
			},
			newReq(http.MethodGet, "ValidWValidClientConcurrency"),
			tstruct{false},
		},
		{
			"ValidWValidClientDelay",
			&httpclient{
				time.Millisecond * 25,
				1,
				http.StatusOK,
				0,
				0,
				0,
				1,
				&testlogger{},
				false,
			},
			newReq(http.MethodGet, "ValidWValidClientDelay"),
			tstruct{false},
		},
		{
			"ValidWValidClientDelayAndConcurrency",
			&httpclient{
				time.Millisecond * 25,
				1,
				http.StatusOK,
				0,
				0,
				1,
				10,
				&testlogger{},
				false,
			},
			newReq(http.MethodGet, "ValidWValidClientDelayAndConcurrency"),
			tstruct{false},
		},
		{
			"ValidWValidClientW5Retries",
			&httpclient{
				0,
				1,
				http.StatusOK,
				5,
				0,
				0,
				1,
				&testlogger{},
				false,
			},
			newReq(http.MethodGet, "ValidWValidClientW5Retries"),
			tstruct{false},
		},
		{
			"FailWValidClientW5Retries4Attempts",
			&httpclient{
				0,
				1,
				http.StatusOK,
				5,
				-1,
				0,
				1,
				&testlogger{},
				false,
			},
			newReq(http.MethodGet, "FailWValidClientW5Retries4Attempts"),
			tstruct{true},
		},
		{
			"FailWBadStatus",
			&httpclient{
				0,
				1,
				http.StatusBadRequest,
				0,
				0,
				0,
				1,
				&testlogger{},
				false,
			},
			newReq(http.MethodGet, "FailWBadStatus"),
			tstruct{true},
		},
		{
			"FailWBadConcurrency",
			&httpclient{
				0,
				1,
				http.StatusOK,
				0,
				0,
				0,
				0,
				&testlogger{},
				false,
			},
			newReq(http.MethodGet, "FailWBadConcurrency"),
			tstruct{true},
		},
		{
			"FailWBadRetries",
			&httpclient{
				0,
				1,
				http.StatusOK,
				-1,
				0,
				0,
				1,
				&testlogger{},
				false,
			},
			newReq(http.MethodGet, "FailWBadRetries"),
			tstruct{true},
		},
		{
			"FailWBadLogger",
			&httpclient{
				0,
				1,
				http.StatusOK,
				0,
				0,
				0,
				1,
				nil,
				false,
			},
			newReq(http.MethodGet, "FailWBadLogger"),
			tstruct{true},
		},
		{
			"FailWBadDelay",
			&httpclient{
				-1,
				1,
				http.StatusOK,
				0,
				0,
				0,
				1,
				&testlogger{},
				false,
			},
			newReq(http.MethodGet, "FailWBadDelay"),
			tstruct{true},
		},
		{
			"FailByCancellation",
			&httpclient{
				0,
				1,
				http.StatusOK,
				0,
				0,
				0,
				1,
				&testlogger{},
				true,
			},
			newReq(http.MethodGet, "FailByCancellation"),
			tstruct{true},
		},
		{
			"FailByCancellation",
			&httpclient{
				time.Millisecond,
				1,
				http.StatusOK,
				0,
				0,
				0,
				1,
				&testlogger{},
				true,
			},
			newReq(http.MethodGet, "FailByCancellation"),
			tstruct{true},
		},
		{
			"FailByNilRequest",
			&httpclient{
				time.Millisecond,
				1,
				http.StatusOK,
				0,
				0,
				0,
				1,
				&testlogger{},
				true,
			},
			nil,
			tstruct{true},
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			defer func() {
				if panicMessage := recover(); panicMessage != nil {
					t.Errorf("test [%s] had a panic", test.name)
				}
			}()

			var client Client
			var err error

			var ctx, cancel = context.WithCancel(context.Background())
			defer cancel()

			if client, err = New(ctx, test.client, test.client.logger, test.client.delay, test.client.retries, test.client.concurrency); err == nil {

				// Cancellation test
				if test.client.cancel {
					cancel()
				}

				var resp *http.Response
				if resp, err = client.Do(test.request); err == nil {
					if resp == nil {
						if err = test.success.correct(true, false); err != nil {
							t.Errorf("[%s] failed; %s", test.name, err.Error())
						}
					}
				} else {

					if test.client.retries > 0 && test.client.retries != test.client.attempts && !test.success.error {
						t.Errorf("[%s] failed; number of attempts doesn't match the expected retries [%v:%v]", test.name, test.client.attempts, test.client.retries)
					} else if err = test.success.correct(true, false); err != nil {
						t.Errorf("[%s] failed; %s", test.name, err.Error())
					}
				}

			} else {
				if err = test.success.correct(true, false); err != nil {
					t.Errorf("[%s] failed; %s", test.name, err.Error())
				}
			}
		})
	}
}

func TestFunnel_DoBadClient(t *testing.T) {
	tests := []struct {
		name    string
		client  *badclient
		request *http.Request
		success tstruct
	}{
		{
			"PanicyClient",
			&badclient{
				true,
				0,
				1,
				http.StatusOK,
				0,
				0,
				0,
				1,
				&testlogger{},
			},
			newReq(http.MethodGet, "fakeurl"),
			tstruct{true},
		},
		{
			"ErroringClietn",
			&badclient{
				false,
				0,
				1,
				http.StatusOK,
				0,
				0,
				0,
				1,
				&testlogger{},
			},
			newReq(http.MethodGet, "fakeurl"),
			tstruct{true},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				if panicMessage := recover(); panicMessage != nil {
					t.Errorf("test [%s] had a panic", test.name)
				}
			}()

			var client Client
			var err error

			if client, err = New(context.Background(), test.client, test.client.logger, test.client.delay, test.client.retries, test.client.concurrency); err == nil {

				var resp *http.Response
				if resp, err = client.Do(test.request); err == nil {
					if resp != nil {
					} else {
						if err = test.success.correct(true, false); err != nil {
							t.Errorf("[%s] failed; %s", test.name, err.Error())
						}
					}
				} else {

					if test.client.retries > 0 && test.client.retries != test.client.attempts && !test.success.error {
						t.Errorf("[%s] failed; number of attempts doesn't match the expected retries [%v:%v]", test.name, test.client.attempts, test.client.retries)
					} else if err = test.success.correct(true, false); err != nil {
						t.Errorf("[%s] failed; %s", test.name, err.Error())
					}
				}
			} else {
				if err = test.success.correct(true, false); err != nil {
					t.Errorf("[%s] failed; %s", test.name, err.Error())
				}
			}
		})
	}
}

func TestFunnel_Do_Direct(t *testing.T) {

	tests := []struct {
		name    string
		fun     *funnel
		request *http.Request
		success tstruct
	}{
		{
			"InvalidFunnel_Concurency",
			&funnel{
				context.Background(),
				&httpclient{},
				0,
				nil,
				0,
				0,
				make(chan bool),
				make(chan requestWrapper),
				&testlogger{},
			},
			newReq(http.MethodGet, "fakeurl"),
			tstruct{true},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				if panicMessage := recover(); panicMessage != nil {
					t.Errorf("test [%s] had a panic", test.name)
				}
			}()

			var err error

			var resp *http.Response
			if resp, err = test.fun.Do(test.request); err == nil {
				if resp != nil {
				} else {
					if err = test.success.correct(true, false); err != nil {
						t.Errorf("[%s] failed; %s", test.name, err.Error())
					}
				}
			} else if err = test.success.correct(true, false); err != nil {
				t.Errorf("[%s] failed; %s", test.name, err.Error())
			}
		})
	}
}

func TestFunnel_Validate(t *testing.T) {
	tests := []struct {
		name    string
		fun     *funnel
		success tstruct
	}{
		{
			"ValidFunnel",
			&funnel{
				context.Background(),
				&httpclient{},
				0,
				nil,
				0,
				1,
				make(chan bool),
				make(chan requestWrapper),
				&testlogger{},
			},
			tstruct{false},
		},
		{
			"InvalidFunnel_Concurency",
			&funnel{
				context.Background(),
				&httpclient{},
				0,
				nil,
				0,
				0,
				make(chan bool),
				make(chan requestWrapper),
				&testlogger{},
			},
			tstruct{true},
		},
		{
			"InvalidFunnel_Retries",
			&funnel{
				context.Background(),
				&httpclient{},
				-1,
				nil,
				0,
				1,
				make(chan bool),
				make(chan requestWrapper),
				&testlogger{},
			},
			tstruct{true},
		},
		{
			"InvalidFunnel_Logger",
			&funnel{
				context.Background(),
				&httpclient{},
				0,
				nil,
				0,
				1,
				make(chan bool),
				make(chan requestWrapper),
				nil,
			},
			tstruct{true},
		},
		{
			"InvalidFunnel_Delay",
			&funnel{
				context.Background(),
				&httpclient{},
				0,
				nil,
				-1,
				1,
				make(chan bool),
				make(chan requestWrapper),
				&testlogger{},
			},
			tstruct{true},
		},
		{
			"InvalidFunnel_Context",
			&funnel{
				nil,
				&httpclient{},
				0,
				nil,
				0,
				1,
				make(chan bool),
				make(chan requestWrapper),
				&testlogger{},
			},
			tstruct{true},
		},
		{
			"InvalidFunnel_ConcurrencyTicker",
			&funnel{
				context.Background(),
				&httpclient{},
				0,
				nil,
				0,
				1,
				nil,
				make(chan requestWrapper),
				&testlogger{},
			},
			tstruct{true},
		},
		{
			"InvalidFunnel_RequestChannel",
			&funnel{
				context.Background(),
				&httpclient{},
				0,
				nil,
				0,
				1,
				make(chan bool),
				nil,
				&testlogger{},
			},
			tstruct{true},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				if panicMessage := recover(); panicMessage != nil {
					t.Errorf("test [%s] had a panic", test.name)
				}
			}()

			var err error

			if err = test.success.correct(!validator.Valid(test.fun), false); err != nil {
				t.Errorf("[%s] failed; %s", test.name, err.Error())
			}
		})
	}
}
