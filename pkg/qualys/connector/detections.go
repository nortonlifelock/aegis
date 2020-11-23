package connector

import (
	"context"
	"fmt"
	"github.com/nortonlifelock/aegis/pkg/domain"
	"github.com/nortonlifelock/aegis/pkg/qualys"
	"sync"
)

func (session *QsSession) pushCombosForHost(ctx context.Context, h qualys.QHost, devVulnMutex *sync.Mutex, processedDevVulns map[string]bool, out chan<- domain.Detection) {
	for index := range h.Detections {
		v := h.Detections[index]

		// Only load detections that are NOT fixed, and are CONFIRMED as vulnerabilities that affect the specific host
		// Confirmed is important because in Qualys there are potential vulnerabilities in the KB and those vulnerabilities
		// can be ACTUAL vulnerabilities when scanned for but the KB always shows them as "Potential" so the actual
		// status on the host itself is what determines if it is actually a vulnerability or not
		if v.Type == "Confirmed" {

			// Read the port information from the detection if a port is specified
			var port = -1
			var protocol = ""
			if v.Port != nil {
				port = *(v.Port)

				if v.Protocol != nil {
					protocol = *(v.Protocol)
				}
			}

			// Lock threads accessing this code and update the processed vulns map so that the records aren't duplicated
			// along the channel for processing
			devVulnMutex.Lock()
			if processedDevVulns[fmt.Sprintf("%v-%v-%v%s", h.HostID, v.QualysID, port, protocol)] == false {
				processedDevVulns[fmt.Sprintf("%v-%v-%v%s", h.HostID, v.QualysID, port, protocol)] = true
				devVulnMutex.Unlock()

				select {
				case <-ctx.Done():
					return
				case out <- &hostDetectionCombo{
					host: &host{
						h: h,
					},
					detection: &detection{
						d:       v,
						session: session,
					},
				}:
				}
			} else {
				devVulnMutex.Unlock()
			}
		}
	}
}
