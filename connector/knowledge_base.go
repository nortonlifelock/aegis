package connector

import (
	"github.com/nortonlifelock/qualys"
	"github.com/nortonlifelock/log"
	"time"
)

// loadAndCacheQualysKB loads ALL vulnerabilities from the Qualys KB into a map that is globally available on the session
func (session *QsSession) loadAndCacheQualysKB(since *time.Time) (err error) {
	var output *qualys.QKnowledgeBaseVulnOutput
	if output, err = session.apiSession.LoadVulnerabilities(since); err == nil {
		session.lstream.Send(log.Info("Vulnerabilities loaded. Beginning processing."))

		// NOTE: DO NOT FILTER OUT POTENTIAL VULNERABILITIES HERE!!! POTENTIAL VULNERABILITIES CAN STILL BE ACTUAL
		// VULNERABILITIES ON THE HOST WHEN DETECTED AS PART OF A SCAN
		for _, v := range output.Vulnerabilities {
			var vuln = v
			// Add the vulnerability to the map on the session object
			session.vulnerabilities[v.QualysID] = &vuln
		}
	}

	return err
}
