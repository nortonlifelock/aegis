package endpoints

import (
	"context"
	"fmt"
	"sync"

	"github.com/nortonlifelock/aegis/internal/integrations"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/jira"
	"github.com/nortonlifelock/log"
	"github.com/pkg/errors"
)

var ticketingConnectionMap map[string]map[string]integrations.TicketingEngine
var ticketingMapLock = &sync.Mutex{}

var cloudConnectionMap map[string]map[string]integrations.CloudServiceConnection
var cloudMapLock = &sync.Mutex{}

func getTicketingConnection(source string, orgID string) (engine integrations.TicketingEngine, err error) {
	ticketingMapLock.Lock()

	if ticketingConnectionMap == nil {
		ticketingConnectionMap = make(map[string]map[string]integrations.TicketingEngine)
	}

	if ticketingConnectionMap[source] == nil {
		ticketingConnectionMap[source] = make(map[string]integrations.TicketingEngine)
	}

	if ticketingConnectionMap[source][orgID] == nil {
		var sourceConfig []domain.SourceConfig

		sourceConfig, err = Ms.GetSourceConfigByNameOrg(source, orgID)
		if err == nil {
			if sourceConfig != nil && len(sourceConfig) > 0 {
				engine, err = integrations.GetEngine(context.Background(), source, Ms, logger{}, AppConfig, sourceConfig[0])
				if err == nil {
					// TODO what do we do if the connection closes? how do we check if the connection is still valid?
					ticketingConnectionMap[source][orgID] = engine
					gatherJiraData(engine, orgID)
				}
			} else {
				err = errors.Errorf("could not find a source config for the source %s and organization %s in the database", source, orgID)
			}

		} else {
			err = errors.Errorf("error while gathering source config for %s engine - %s", source, err.Error())
		}
	} else {
		engine = ticketingConnectionMap[source][orgID]
	}

	ticketingMapLock.Unlock()

	return engine, err
}

func getCloudConnection(source string, orgID string) (cloudConnection integrations.CloudServiceConnection, err error) {
	cloudMapLock.Lock()

	if cloudConnectionMap == nil {
		cloudConnectionMap = make(map[string]map[string]integrations.CloudServiceConnection)
	}

	if cloudConnectionMap[source] == nil {
		cloudConnectionMap[source] = make(map[string]integrations.CloudServiceConnection)
	}

	if cloudConnectionMap[source][orgID] == nil {
		var sourceConfig []domain.SourceConfig

		sourceConfig, err = Ms.GetSourceConfigByNameOrg(source, orgID)
		if err == nil {
			if sourceConfig != nil && len(sourceConfig) > 0 {

				cloudConnection, err = integrations.GetCloudServiceConnection(Ms, source, sourceConfig[0], AppConfig, logger{})
				if err == nil {
					// TODO what do we do if the connection closes? how do we check if the connection is still valid?
					cloudConnectionMap[source][orgID] = cloudConnection
				} else {
					err = fmt.Errorf("error creating cloud connection - %s", err.Error())
				}

			} else {
				err = errors.Errorf("could not find a source config for the source %s and organization %s in the database", source, orgID)
			}

		} else {
			err = errors.Errorf("error while gathering source config for %s engine - %s", source, err.Error())
		}
	} else {
		cloudConnection = cloudConnectionMap[source][orgID]
	}

	cloudMapLock.Unlock()

	return cloudConnection, err
}

type logger struct{}

// Send pushes a log onto the log receiver channel
func (logger logger) Send(log log.Log) {
	_ = logFunc("", log.Text, log.Error)
}

func logFunc(logType string, log string, logError error) (err error) {
	_ = logType
	if logError == nil {
		fmt.Println(log)
	} else {
		fmt.Println(fmt.Sprintf("%s : %s", log, logError.Error()))
	}

	return err
}

func gatherJiraData(engine integrations.TicketingEngine, orgID string) {
	if connector, ok := engine.(*jira.ConnectorJira); ok {
		gatherBurndownInfo(connector, orgID)
	}
}

func gatherBurndownInfo(connector *jira.ConnectorJira, orgID string) {
	var organization domain.Organization
	var err error
	organization, err = Ms.GetOrganizationByID(orgID)

	if err == nil {
		burnDownPackage := &BurndownPackage{}
		burnDownPackage.MediumStats = &BurndownStats{}
		burnDownPackage.HighStats = &BurndownStats{}
		burnDownPackage.CritStats = &BurndownStats{}

		burnDownLock.Lock()
		burnDownCache[orgID] = burnDownPackage
		burnDownLock.Unlock()

		go calculateBurndownForPriority(connector, "Medium", 90, burnDownPackage.MediumStats, organization.Code())
		go calculateBurndownForPriority(connector, "High", 60, burnDownPackage.HighStats, organization.Code())
		go calculateBurndownForPriority(connector, "Critical", 30, burnDownPackage.CritStats, organization.Code())
	}
}
