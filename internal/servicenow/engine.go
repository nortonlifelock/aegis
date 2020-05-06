package servicenow

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/funnel"
	"github.com/nortonlifelock/log"
	"net/http"
	"sync"
)

const (
	svcNowDateFormat = "2006-01-02T15:04:05.999-0700"
	//SVCNOWDATEONLY   = "2006-01-02"
)

type logger interface {
	Send(log log.Log)
}

// SvcNowConnector implements the TicketingEngine interface for Service Now
type SvcNowConnector struct {
	logger  logger
	db      domain.DatabaseConnection
	wg      *sync.WaitGroup
	config  domain.SourceConfig
	client  *SvcNowClient
	project string
}

type svcNowPayload struct {
	Project string `json:"project"`
}

// NewSvcNowConnector initializes a SvcNowConnector, and parses it's payload from the source config
func NewSvcNowConnector(ctx context.Context, db domain.DatabaseConnection, logger logger, config domain.SourceConfig) (connector *SvcNowConnector, err error) {

	if config.Payload() != nil && len(*config.Payload()) > 0 {

		var payload svcNowPayload
		if err = json.Unmarshal([]byte(*config.Payload()), &payload); err == nil {

			if len(payload.Project) > 0 {

				connector = &SvcNowConnector{
					project: payload.Project,
					db:      db,
					logger:  logger,
					config:  config,
					wg:      &sync.WaitGroup{},
				}

				var authInfo domain.BasicAuth
				if err = json.Unmarshal([]byte(config.AuthInfo()), &authInfo); err == nil {
					// Initialize Client
					err = connector.initClient(ctx, logger, config.Address(), authInfo.Username, authInfo.Password)
				} else {
					err = fmt.Errorf("error while parsing authentication information - %s", err.Error())
				}

			} else {
				err = fmt.Errorf("empty project specified for Source Config [SourceConfigId: %v]", config.ID())
			}
		} else {
			err = fmt.Errorf("error while parsing authentication information - %s", err.Error())
		}
	} else {
		err = fmt.Errorf("invalid Payload for Ticketing Job in Source Config [SourceConfigId: %v]", config.ID())
	}

	return connector, err
}

// Create a jira client instance
func (connector *SvcNowConnector) initClient(ctx context.Context, logger logger, svcNowURL, svcNowUser, svcNowPass string) (err error) {
	var funnelClient funnel.Client
	funnelClient, err = funnel.New(ctx, &http.Client{}, logger, 0, 0, 10)
	if err == nil {
		connector.client, err = NewClient(funnelClient, svcNowURL, svcNowUser, svcNowPass)
	}

	return err
}
