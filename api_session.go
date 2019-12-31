package qualys

import (
	"context"
	"fmt"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/funnel"
	"github.com/nortonlifelock/log"
	"net/http"
	"sync"
	"time"
)

type logger interface {
	Send(log log.Log)
}

// Session is the struct that is responsible for making Qualys API calls
type Session struct {
	lstream logger

	// Source configuration which holds the authentication information for the Qualys API
	Config domain.SourceConfig

	// each endpoint has a separate rate limit, so we create a funnel for each endpoint
	funnelMap  map[string]funnel.Client
	funnelLock *sync.Mutex

	ctx              context.Context
	concurrencyLimit int
	rateLimit        time.Duration
}

// NewQualysAPISession initializes the Qualys session object but currently all authentication is done with basic auth rather
// than session management which will be changing in a future release
func NewQualysAPISession(ctx context.Context, lstream logger, config domain.SourceConfig) (session *Session, err error) {
	session = &Session{
		lstream: lstream,
		Config:  config,

		funnelMap:  make(map[string]funnel.Client),
		funnelLock: &sync.Mutex{},

		ctx: ctx,
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
	session.funnelLock.Lock()
	defer session.funnelLock.Unlock()

	if session.funnelMap[endpoint] != nil {
		client = session.funnelMap[endpoint]
	} else {
		if session.concurrencyLimit < 1 {
			session.concurrencyLimit = 1
		}

		if session.rateLimit < 1 {
			session.rateLimit = 1
		}

		var err error
		if client, err = funnel.New(session.ctx, &http.Client{}, session.lstream, session.rateLimit*time.Second, 1, session.concurrencyLimit); err == nil && client != nil {
			session.funnelMap[endpoint] = client
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
	req, err = http.NewRequest(http.MethodGet, session.Config.Address()+qsAppliance, nil)
	if err == nil {
		err = session.makeRequest(req, func(response *http.Response) (err error) {
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
