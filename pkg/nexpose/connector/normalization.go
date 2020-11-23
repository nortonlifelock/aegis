package connector

import (
	"github.com/nortonlifelock/aegis/pkg/domain"
	"sync"
)

var scanStatusesMutex = sync.Mutex{}

// scanStatuses return the normalized scan statuses for nexpose scans
var scanStatuses = map[string]string{
	"integrating": domain.ScanPROCESSING,
	"running":     domain.ScanPROCESSING,
	"paused":      domain.ScanPAUSED,
	"stopped":     domain.ScanSTOPPED,
	"finished":    domain.ScanFINISHED,
	"error":       domain.ScanERRORED,
}

func scanStatus(status string) string {
	scanStatusesMutex.Lock()
	defer scanStatusesMutex.Unlock()

	return scanStatuses[status]
}

var detectionStatusesMutex = sync.Mutex{}

// Status list
//vulnerable-version-with-exception-applied
//vulnerable
//vulnerable-version
//vulnerable-potential
//vulnerable-potential-with-exception-applied
//invulnerable
//vulnerable-with-exception-applied
var detectionStatuses = map[string]string{
	"vulnerable":                                  domain.Vulnerable,
	"vulnerable-version":                          domain.Vulnerable,
	"vulnerable-potential":                        domain.Vulnerable,
	"vulnerable-version-with-exception-applied":   domain.Fixed,
	"vulnerable-potential-with-exception-applied": domain.Fixed,
	"invulnerable":                                domain.Fixed,
	"vulnerable-with-exception-applied":           domain.Fixed,
}

func detectionStatus(status string) string {
	detectionStatusesMutex.Lock()
	defer detectionStatusesMutex.Unlock()

	return detectionStatuses[status]
}
