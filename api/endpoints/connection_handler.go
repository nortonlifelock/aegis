package endpoints

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/nortonlifelock/aegis/internal/integrations"
	"github.com/nortonlifelock/aegis/pkg/domain"
	"github.com/nortonlifelock/aegis/pkg/log"
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
					go gatherBurndownInfo(orgID)
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

func gatherBurndownInfo(orgID string) {
	for {
		burnDownPackage := &BurndownPackage{}
		burnDownPackage.MediumStats = &BurndownStats{}
		burnDownPackage.HighStats = &BurndownStats{}
		burnDownPackage.CritStats = &BurndownStats{}

		calculateBurndownForPriority(7, 4, 90, burnDownPackage.MediumStats, orgID)
		calculateBurndownForPriority(9, 7, 60, burnDownPackage.HighStats, orgID)
		calculateBurndownForPriority(10, 9, 30, burnDownPackage.CritStats, orgID)

		burnDownLock.Lock()
		burnDownCache[orgID] = burnDownPackage
		burnDownLock.Unlock()

		time.Sleep(time.Minute * 10)
	}
}
