package qualys

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nortonlifelock/aegis/pkg/domain"
	"github.com/nortonlifelock/aegis/pkg/funnel"
	"github.com/nortonlifelock/aegis/pkg/log"
	"net/http"
	"sync"
	"time"
)

// each endpoint has a separate rate limit, so we create a funnel for each endpoint
var (
	funnelMap  = make(map[string]funnel.Client)
	funnelLock = &sync.Mutex{}
)

type logger interface {
	Send(log log.Log)
}

// Session is the struct that is responsible for making Qualys API calls
type Session struct {
	lstream logger

	// Source configuration which holds the authentication information for the Qualys API
	Config domain.SourceConfig

	ctx              context.Context
	concurrencyLimit int
	rateLimit        time.Duration

	// Web Application Scanning API uses a different base URL from the Qualys VM API
	// we can gather the second URL from the AuthInfo column for the Qualys connection in the SourceConfig table
	// the url JSON key is "host_web_app"
	webAppBaseURL string
}

// NewQualysAPISession initializes the Qualys session object but currently all authentication is done with basic auth rather
// than session management which will be changing in a future release
func NewQualysAPISession(ctx context.Context, lstream logger, config domain.SourceConfig) (session *Session, err error) {
	session = &Session{
		lstream:       lstream,
		Config:        config,
		webAppBaseURL: getWebAppURLIfPresent(config),
		ctx:           ctx,
	}

	var rates Rates
	if rates, err = session.pullRateInfoForInitialization(); err == nil {

		// RLLimit is requests allowed within a time period
		// RLWindowSec is the time period which the requests are limited
		// The ratio between of the two is the frequency one is able to make API calls for any particular endpoint
		var reqsPerSecond float32
		if rates.RLLimit > 0 {
			reqsPerSecond = float32(rates.RLWindowSec) / float32(rates.RLLimit)
		}

		if reqsPerSecond < 1 {
			reqsPerSecond = 1
		}

		session.rateLimit = time.Duration(reqsPerSecond)

		session.concurrencyLimit = rates.CLimit
		if session.concurrencyLimit < 1 {
			session.concurrencyLimit = 1
		}
	}

	return session, err
}

func (session *Session) getFunnelForEndpoint(endpoint string) (client funnel.Client) {
	funnelLock.Lock()
	defer funnelLock.Unlock()

	mapKey := fmt.Sprintf("%s;%s", session.Config.Address(), endpoint)
	if funnelMap[mapKey] != nil {
		client = funnelMap[mapKey]
	} else {
		if session.concurrencyLimit < 1 {
			session.concurrencyLimit = 1
		}

		if session.rateLimit < 1 {
			session.rateLimit = 1
		}

		var err error
		if client, err = funnel.New(context.Background(), &http.Client{}, session.lstream, session.rateLimit*time.Second, 1, session.concurrencyLimit); err == nil && client != nil {
			funnelMap[mapKey] = client
		} else {
			session.lstream.Send(log.Error("error while creating a Qualys funnel, defaulting to http client", err))
			client = &http.Client{}
		}
	}

	return client
}

func (session *Session) pullRateInfoForInitialization() (rates Rates, err error) {
	var req *http.Request
	// we make a request to a random authenticated endpoint to get the rate limit information in the response headers
	req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%s?action=list", session.Config.Address()+qsAppliance), nil)
	if err == nil {
		err = session.makeRequest(true, req, func(response *http.Response) (err error) {
			defer response.Body.Close()

			var status int
			status, rates, err = session.pullResponseHeaders(response)
			if err == nil {
				if status != http.StatusOK {
					err = fmt.Errorf("response [%d] from server", status)
				}
			}

			return err
		})
	}

	return rates, err
}

func getWebAppURLIfPresent(sourceConfig domain.SourceConfig) (url string) {
	type parseWebApp struct {
		URL string `json:"host_web_app"`
	}
	parse := &parseWebApp{}
	_ = json.Unmarshal([]byte(sourceConfig.AuthInfo()), parse)
	url = parse.URL
	return url
}
