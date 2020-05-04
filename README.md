# Funnel - Request Throttler for http.Client

[![CI](https://github.com/nortonlifelock/funnel/workflows/CI/badge.svg)](https://github.com/nortonlifelock/funnel/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/nortonlifelock/funnel)](https://goreportcard.com/report/github.com/nortonlifelock/funnel)
[![GoDoc](https://godoc.org/github.com/nortonlifelock/funnel?status.svg)](https://pkg.go.dev/github.com/nortonlifelock/funnel)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](http://makeapullrequest.com)

Funnel implements the interface shown below

```go
type Client interface {
    Do(request *http.Request) (*http.Response, error)
}
```

This interface is also implemented by the default `http.Client`.

`funnel` replaces the hard implementation of `http.Client` with an
implementation of a shared interface such that anything accepting the
above interface can use `funnel` to throttle their API requests through
the configuration of funnel.

An example of this is handling API rate limiting from an API you do not
control. Funnel can be configured through the `funnel.New` method.

