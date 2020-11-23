package funnel

import "net/http"

// Client defines an interface with a Do method for use with http.Client or other
// mocked instances that implement a do method that accept a request and return a response error combo
type Client interface {
	Do(req *http.Request) (*http.Response, error)
}
