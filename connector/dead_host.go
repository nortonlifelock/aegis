package connector

import (
	"encoding/xml"
	"fmt"
	"github.com/nortonlifelock/qualys"
	"time"
)

func (session *QsSession) getDeadHostsForScan(scanID string, created time.Time) (deadHostIPToProof map[string]string, err error) {
	deadHostIPToProof = make(map[string]string)

	var output *qualys.ScanSummaryOutput
	if output, err = session.apiSession.GatherDeadHostsFoundSince(created); err == nil {
		for _, scan := range output.Response.ScanSummaryList.ScanSummary {
			if scan.ScanRef == scanID {
				for _, host := range scan.HostSummary {
					var proofByte []byte
					if proofByte, err = xml.MarshalIndent(scan, "", "\t"); err == nil {
						// host.Text contains the IP of the host
						// proofByte contains the portion of the XML that displays the host as dead
						deadHostIPToProof[host.Text] = string(proofByte)
					} else {
						err = fmt.Errorf("error while marshalling proof [%v]", err)
						break
					}

				}
			}
		}
	}

	return deadHostIPToProof, err
}
