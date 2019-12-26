package dome9

import (
	"encoding/base64"
	"net/http"

	"github.com/nortonlifelock/log"
)

// logger is the interface that defines the required method for processing a log
type logger interface {
	Send(log log.Log)
}

// Client holds the authentication information as well as the Dome9 URL for future API calls
type Client struct {
	authString string
	client     *http.Client
	baseURL    string
	lstream    logger
}

// CreateClient decrypts the Dome9 source configs and builds the auth string for future API calls
func CreateClient(user string, password string, address string, lstream logger) (client *Client, err error) {
	client = &Client{}
	client.client = &http.Client{}
	client.baseURL = address
	client.lstream = lstream
	client.authString = basicAuth(user, password)

	return client, err
}

func basicAuth(keyID string, keySecret string) string {
	auth := keyID + ":" + keySecret
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
