package jira

import (
	"context"
	"encoding/json"
	"fmt"
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
	EncryptionType() string
	EncryptionKey() string
	KMSRegion() string
	KMSProfile() string
}

// funnelMap ensures all instances of JIRA don't create separate funnels
// the key is the URL of the JIRA instance and the value is the API funnel for that instance
var funnelMap sync.Map

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

	// TransitionMap has two keys - the first key being the starting status, the second key being the ending status, and the value being the series of transitions (in order) required to get from the starting status to the ending status
	TransitionMap map[string]map[string][]workflowTransition `json:"transition_map"`
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
		project:       payload.Project,
		config:        config,
		lstream:       lstream,
		payload:       payload,
		statusMap:     payload.StatusMap,
		TransitionMap: payload.TransitionMap,
	}

	if checkStatus := ensureAllStatusesExistInMap(connector.payload.StatusMap); checkStatus != nil {
		connector.lstream.Send(log.Warning("did not find all statuses in JIRA status map", checkStatus))
	}

	var authInfo = &domain.AllAuth{}
	authInfoJSON := config.AuthInfo()
	if err = json.Unmarshal([]byte(authInfoJSON), authInfo); err == nil {

		var authLayerClient funnel.Client
		if len(authInfo.PrivateKey) > 0 && len(authInfo.ConsumerKey) > 0 {
			authLayerClient, token, err = connector.getOauthClient(config.Address(), authInfo.PrivateKey, authInfo.ConsumerKey, authInfo.Token)
		} else {
			connector.lstream.Send(log.Warning("JIRA using basic authentication instead of Oauth", err))
			authLayerClient, err = connector.initBasicClient(config.Address(), authInfo.Username, authInfo.Password)
		}

		if err == nil {

			var driverFunnel funnel.Client
			if clientInterface, ok := funnelMap.Load(config.Address()); ok {
				if driverFunnel, ok = clientInterface.(funnel.Client); !ok {
					err = fmt.Errorf("error while loading client [%v] from cache", clientInterface)
				}
			} else {
				// driverFunnel uses a background context so a completing job doesn't destroy the shared funnel
				if driverFunnel, err = funnel.New(context.Background(), authLayerClient, lstream, authInfo.Delay(), authInfo.Retries(), authInfo.Concurrency()); err == nil {

					// perform a load or store so if multiple funnels are created due to concurrency, each jira driver still share a single client
					clientInterface, _ = funnelMap.LoadOrStore(config.Address(), driverFunnel)
					if driverFunnel, ok = clientInterface.(funnel.Client); !ok {
						err = fmt.Errorf("error while loading client [%v] from cache", clientInterface)
					}
				}
			}

			if err == nil {
				// instanceFunnel != driverFunnel as the driver funnel needs a context that cannot be cancelled
				// instanceFunnel uses the context that can be cancelled by a parent
				var instanceFunnel funnel.Client
				if instanceFunnel, err = funnel.New(ctx, driverFunnel, lstream, authInfo.Delay(), authInfo.Retries(), authInfo.Concurrency()); err == nil {
					connector.funnelClient = instanceFunnel

					if err = connector.initJIRALayerClient(connector.funnelClient, config.Address()); err == nil {
						// loads the fields, resolutions, statuses, issue types, and workflow
						err = connector.loadConnectorData()
					}
				}
			}
		}
	}

	return connector, token, err
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
					if len(connector.TransitionMap) == 0 {
						err = fmt.Errorf("transition map not present in JIRA source config payload")
					}
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

				var authLayerClient funnel.Client
				if authLayerClient, err = connector.initBasicClient(api, user, password); err == nil {
					if connector.funnelClient, err = funnel.New(context.Background(), authLayerClient, lstream, 5, 2, 10); err == nil {

						if err = connector.initJIRALayerClient(connector.funnelClient, api); err == nil {
							// Load the Fields for the connector object
							connector.Fields, err = connector.getFields()
						}
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
