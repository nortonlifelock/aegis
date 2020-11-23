package domain

import (
	"time"
)

// AllAuth is used as a container for all authentication methods. It is useful for the API which often has to deal with source configs generically
type AllAuth struct {
	BasicAuth
	OauthAuth
	Host
}

// Host is a struct that contains all of the important information for a host that is
// going to be authenticated to including funneling information and tuning
type Host struct {
	Path             string `json:"Host,omitempty"`
	Porticus         int    `json:"Port,omitempty"`
	Verify           bool   `json:"VerifyTLS,omitempty"`
	TimeDelay        int    `json:"Delay,omitempty"`
	ConcurrencyLimit int    `json:"Concurrency,omitempty"`
	RetryLimit       int    `json:"Retries,omitempty"`
	CacheTTLSeconds  *int   `json:"CacheTTLSeconds,omitempty"`
	*BasicAuth
}

// User returns the username for the basic auth
func (h *Host) User() string {
	return h.Username
}

// Pass returns the password for the basic auth
func (h *Host) Pass() string {
	return h.Password
}

// Host returns the host name of the endpoint
func (h *Host) Host() string {
	return h.Path
}

// Port returns the port for the API connection
func (h *Host) Port() int {
	return h.Porticus
}

// VerifyTLS indicates whether to attempt a tls certificate verification
// on the endpoint
func (h *Host) VerifyTLS() bool {
	return h.Verify
}

// Delay returns the time delay to use when accessing the api endpoints for this asset
func (h *Host) Delay() time.Duration {
	return time.Duration(h.TimeDelay)
}

// Concurrency indicates the maximum number of concurrent requests
// to make at a time against this asset
func (h *Host) Concurrency() int {
	return h.ConcurrencyLimit
}

// Retries indicates the number of retries to attempt of a call to an endpoint
// fails against the API
func (h *Host) Retries() int {
	return h.RetryLimit
}

// BasicAuth is used to parse the authentication information from the AuthInfo field in the SourceConfig database
type BasicAuth struct {
	Username string `json:"Username,omitempty"`
	Password string `json:"Password,omitempty"`
}

// User returns the username for the basic auth
func (ba *BasicAuth) User() string {
	return ba.Username
}

// Pass returns the password for the basic auth
func (ba *BasicAuth) Pass() string {
	return ba.Password
}

// OauthAuth is used to parse the authentication information from the AuthInfo field in the SourceConfig database
type OauthAuth struct {
	PrivateKey  string `json:"PrivateKey,omitempty"`
	ConsumerKey string `json:"ConsumerKey,omitempty"`
	Token       string `json:"Token,omitempty"`
}

// PK returns the private key to use for creating an OAuth connection
func (oa *OauthAuth) PK() string {
	return oa.PrivateKey
}

// CK returns the consumer key to use for creating an OAuth connection
func (oa *OauthAuth) CK() string {
	return oa.ConsumerKey
}

// TK returns the token to use for creating an OAuth connection
func (oa *OauthAuth) TK() string {
	return oa.Token
}

// TODO: implement validators
