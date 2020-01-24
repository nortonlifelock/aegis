package jira

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/nortonlifelock/funnel"
	"net/http"
	"net/url"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/dghubble/oauth1"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/log"
)

// Create a jira client instance
// Pulled from https://gist.github.com/Lupus/edafe9a7c5c6b13407293d795442fe67
func (connector *ConnectorJira) getOauthClient(JiraURL, JiraPk, JiraCk, JiraToken string) (client funnel.Client, token string, err error) {
	if len(JiraURL) > 0 {
		if len(JiraPk) > 0 {
			if len(JiraCk) > 0 {

				ctx := context.Background()
				if keyDERBlock, _ := pem.Decode([]byte(JiraPk)); keyDERBlock != nil &&
					(keyDERBlock.Type == "PRIVATE KEY" ||
						strings.HasSuffix(keyDERBlock.Type, " PRIVATE KEY")) {

					client, token, err = connector.loadOauthToken(ctx, keyDERBlock, JiraCk, JiraURL, JiraToken)

				} else {
					err = fmt.Errorf("unable to decode key PEM block")
				}

			} else {
				err = fmt.Errorf("JIRA Consumer Key is empty")
			}
		} else {
			err = fmt.Errorf("JIRA Private Key is empty")
		}
	} else {
		err = fmt.Errorf("missing JIRA API path")
	}

	return client, token, err
}

func (connector *ConnectorJira) loadOauthToken(ctx context.Context, keyDERBlock *pem.Block, JiraCk string, JiraURL string, JiraToken string) (client *http.Client, token string, err error) {
	var privateKey *rsa.PrivateKey
	privateKey, err = x509.ParsePKCS1PrivateKey(keyDERBlock.Bytes)
	if err == nil {

		config := oauth1.Config{
			ConsumerKey: JiraCk,
			CallbackURL: "oob", /* for command line usage */
			Endpoint: oauth1.Endpoint{
				RequestTokenURL: JiraURL + "/plugins/servlet/oauth/request-token",
				AuthorizeURL:    JiraURL + "/plugins/servlet/oauth/authorize",
				AccessTokenURL:  JiraURL + "/plugins/servlet/oauth/access-token",
			},
			Signer: &oauth1.RSASigner{
				PrivateKey: privateKey,
			},
		}

		tok := &oauth1.Token{}

		// Create the client using the cached token from the database
		if err = json.Unmarshal([]byte(JiraToken), tok); err != nil { // Get a new token from JIRA

			if tok, err = connector.getTokenFromJIRA(&config); err == nil {
				token, err = connector.extractAndStoreToken(tok)
			}
		}

		client = config.Client(ctx, tok)
	} else {
		err = fmt.Errorf("unable to parse PKCS1 private key. %v", err)
	}

	return client, token, err
}

func (connector *ConnectorJira) extractAndStoreToken(tok *oauth1.Token) (token string, err error) {
	var tokenJSON []byte
	if tokenJSON, err = json.Marshal(tok); err == nil {

		// Update in memory value in case a new connection is initialized
		err = connector.updateToken(string(tokenJSON))
		if err == nil {
			token = string(tokenJSON)
		} else {
			connector.lstream.Send(log.Warning("error while updating JIRA Oauth token in memory", err))
		}

	} else {
		err = fmt.Errorf("error while marshalling token - %s", err.Error())
	}

	return token, err
}

func (connector *ConnectorJira) updateToken(token string) (err error) {
	var authInfo domain.OauthAuth
	if err = json.Unmarshal([]byte(connector.config.AuthInfo()), &authInfo); err == nil {
		authInfo.Token = token
		var authInfoUpdated []byte
		if authInfoUpdated, err = json.Marshal(authInfo); err == nil {
			//connector.config.SetAuthInfo(string(authInfoUpdated))
			_ = authInfoUpdated
		}
	}

	return err
}

func (connector *ConnectorJira) initJIRALayerClient(httpClient funnel.Client, baseURL string) (err error) {
	connector.client, err = jira.NewClient(httpClient, baseURL)
	return err
}

// pulled from https://gist.github.com/Lupus/edafe9a7c5c6b13407293d795442fe67
func (connector *ConnectorJira) getTokenFromJIRA(config *oauth1.Config) (token *oauth1.Token, err error) {
	var requestToken string
	var requestSecret string
	if requestToken, requestSecret, err = config.RequestToken(); err == nil {

		var authorizationURL *url.URL
		if authorizationURL, err = config.AuthorizationURL(requestToken); err == nil {
			fmt.Printf("Go to the following link in your browser then type the "+
				"authorization code: \n%v\n", authorizationURL.String())

			var code string
			if _, err = fmt.Scan(&code); err == nil {
				var accessToken string
				var accessSecret string

				if accessToken, accessSecret, err = config.AccessToken(requestToken, requestSecret, code); err == nil {

					token = oauth1.NewToken(accessToken, accessSecret)

				} else {
					err = fmt.Errorf("unable to get access token. %v", err)
				}
			} else {
				err = fmt.Errorf("unable to read authorization code. %v", err)
			}
		} else {
			err = fmt.Errorf("unable to get authorization url. %v", err)
		}
	} else {
		err = fmt.Errorf("unable to get request token. %v", err)
	}

	return token, err
}

func (connector *ConnectorJira) initBasicClient(JiraURL, JiraUser, JiraPass string) (client funnel.Client, err error) {
	tp := jira.BasicAuthTransport{
		Username: JiraUser,
		Password: JiraPass,
	}

	return tp.Client(), err
}
