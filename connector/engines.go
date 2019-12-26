package connector

import (
	"fmt"
	"github.com/nortonlifelock/nexpose"
	"strconv"
	"time"
)

// GetSiteEngine loads the engine information from cache if it exists. If it does not exist, it queries the Nexpose API for the information and then caches it
func (conn *Connection) GetSiteEngine(siteID string) (engineID string, err error) {
	var engine *nexpose.Engine

	if obj, ok := conn.engines.Load(siteID); ok {
		engine, _ = obj.(*nexpose.Engine)
	}

	if engine == nil {
		if engine, err = conn.api.GetEngineForSite(siteID); err == nil {
			engineID = strconv.Itoa(engine.ID)

			var cacheTime = time.Minute // default to a minute
			if conn.settings.EngineCacheTTL != nil {
				cacheTime = *conn.settings.EngineCacheTTL * time.Minute
			}

			conn.engines.Store(conn.ctx, siteID, engine, &cacheTime)
		} else {
			err = fmt.Errorf("error while loading engine for site [%v] - %s", siteID, err.Error())
		}
	}

	return engineID, err
}
