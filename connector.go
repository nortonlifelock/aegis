package jira

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/andygrunwald/go-jira"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/funnel"
	"github.com/nortonlifelock/log"
)

type logger interface {
	Send(log log.Log)
}

// ConfigJira is an interface that defines the fields that JIRA requires from the app config
type ConfigJira interface {
	EncryptionKey() string
}

// ConnectorJira is the struct that is used to make API calls against JIRA
type ConnectorJira struct {
	client *jira.Client

	config  domain.SourceConfig
	lstream logger
	project string

	statusMap map[string]string
	CERFs     sync.Map
	//RequestChan chan *Request
	funnelClient funnel.Client

	payload PayloadJira

	Fields        map[string]*Field
	Resolutions   map[string]*jira.Resolution
	Statuses      map[string]*jira.Status
	IssueTypes    map[string]jira.IssueType
	TransitionMap map[string]map[string][]workflowTransition
}

// PayloadJira contains the structure of the JSON fields that need to be loaded from the JIRA source config
type PayloadJira struct {
	Project string `json:"project"`

	// Mappable fields is used by the UI to discern which JIRA fields can be mapped to a field supplied by a cloud provider
	MappableFields []string `json:"mappable_fields,omitempty"`

	// Maps the PDE name for a status to the JIRA-specific name for a status
	StatusMap map[string]string `json:"status_map"`

	// Maps the PDE name for a field to the JIRA-specific name for a field
	FieldMap map[string]string `json:"field_map"`
}

// NewJiraConnector creates a new JIRA connector for interacting with a JIRA system
// it attempts to use Oauth if the information is present in the source config. If not, basic authentication is used
func NewJiraConnector(ctx context.Context, lstream logger, config domain.SourceConfig) (connector *ConnectorJira, token string, err error) {

	if config.Payload() != nil && len(sord(config.Payload())) > 0 {

		var payload PayloadJira
		if err = json.Unmarshal([]byte(sord(config.Payload())), &payload); err == nil {

			if len(payload.Project) > 0 {
				connector, token, err = buildConnector(ctx, payload, config, lstream)
			} else {
				err = fmt.Errorf("empty project specified for Source Config [SourceConfigId: %v]", config.ID())
			}
		}
	} else {
		err = fmt.Errorf("invalid payload in jira source config [SourceConfigId: %v]", config.ID())
	}

	return connector, token, err
}

func buildConnector(ctx context.Context, payload PayloadJira, config domain.SourceConfig, lstream logger) (connector *ConnectorJira, token string, err error) {
	connector = &ConnectorJira{
		project:   payload.Project,
		config:    config,
		lstream:   lstream,
		payload:   payload,
		statusMap: payload.StatusMap,
	}

	if checkStatus := ensureAllStatusesExistInMap(connector.payload.StatusMap); checkStatus != nil {
		connector.lstream.Send(log.Warning("did not find all statuses in JIRA status map", checkStatus))
	}

	var authInfo = &domain.AllAuth{}
	authInfoJSON := config.AuthInfo()
	if err = json.Unmarshal([]byte(authInfoJSON), authInfo); err == nil {

		if len(authInfo.PrivateKey) > 0 && len(authInfo.ConsumerKey) > 0 {
			var client *http.Client
			client, token, err = connector.getOauthClient(config.Address(), authInfo.PrivateKey, authInfo.ConsumerKey, authInfo.Token)
			if err == nil {
				err = connector.initClient(client, config.Address())
			}
		} else {
			connector.lstream.Send(log.Warning("JIRA using basic authentication instead of Oauth", err))
			err = connector.initBasicClient(config.Address(), authInfo.Username, authInfo.Password)
		}

		if connector.client != nil && err == nil {

			var client funnel.Client
			if client, err = funnel.New(ctx, &jiraClientWrapper{connector.client}, lstream, authInfo.Delay(), authInfo.Retries(), authInfo.Concurrency()); err == nil {

				connector.funnelClient = client

				// loads the fields, resolutions, statuses, issue types, and workflow
				err = connector.loadConnectorData()
			}
		} else {
			if err == nil {
				err = fmt.Errorf("could not authenticate JIRA client")
			} else {
				err = fmt.Errorf("could not authenticate JIRA client - %s", err.Error())
			}
		}
	}

	return connector, token, err
}

type jiraClientWrapper struct {
	client *jira.Client
}

func (wrapper *jiraClientWrapper) Do(req *http.Request) (resp *http.Response, err error) {
	var jiraResp *jira.Response
	jiraResp, err = wrapper.client.Do(req, nil)
	if err == nil {
		if jiraResp != nil {
			resp = jiraResp.Response
		} else {
			err = fmt.Errorf("nil response")
		}
	}

	return resp, err
}

func ensureAllStatusesExistInMap(statusMap map[string]string) (err error) {
	var desiredStatuses = []string{
		StatusReopened, StatusClosedRemediated, StatusClosedFalsePositive, StatusClosedDecommissioned,
		StatusOpen, StatusInProgress, StatusResolvedException, StatusClosedException,
		StatusResolvedDecom, StatusResolvedRemediated,
	}

	for _, status := range desiredStatuses {
		if len(statusMap[strings.ToLower(status)]) == 0 {
			err = fmt.Errorf("could not find status map for [%s] in JIRA payload", status)
			break
		}
	}

	return err
}

func (connector *ConnectorJira) loadConnectorData() (err error) {
	// Load the Fields for the connector object
	if connector.Fields, err = connector.getFields(); err == nil {

		// Load the Resolutions for the connector object
		if connector.Resolutions, err = connector.getResolutions(); err == nil {

			// Load the Statuses for the connector object
			if connector.Statuses, err = connector.getStatuses(); err == nil {

				// Load the issue types for the connector object
				if connector.IssueTypes, err = connector.getIssueTypes(); err == nil {
					err = connector.getWorkflow()
				}
			}
		}
	}
	return err
}

// ConnectJira creates a JIRA connection using basic authentication
func ConnectJira(api string, user string, password string, lstream logger) (connector *ConnectorJira, err error) {

	if len(api) > 0 {
		if len(user) > 0 {
			if len(password) > 0 {
				connector = &ConnectorJira{
					lstream: lstream,
				}

				if err = connector.initBasicClient(api, user, password); err == nil {
					if connector.client != nil {
						if connector.funnelClient, err = funnel.New(context.Background(), &jiraClientWrapper{connector.client}, lstream, 5, 2, 10); err == nil {
							// Load the Fields for the connector object
							connector.Fields, err = connector.getFields()
						}
					} else {
						err = fmt.Errorf("could not authenticate JIRA client")
					}
				}
			} else {
				err = fmt.Errorf("password cannot be empty")
			}
		} else {
			err = fmt.Errorf("username cannot be empty")
		}
	} else {
		err = fmt.Errorf("API Path cannot be empty")
	}

	return connector, err
}

// GetProject returns the project that a JIRA connection is pointing to (defined in the JIRA source config payload)
func (connector *ConnectorJira) GetProject() string {
	return connector.project
}
