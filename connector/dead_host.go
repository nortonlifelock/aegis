package connector

import (
	"encoding/xml"
	"fmt"
	"github.com/nortonlifelock/qualys"
	"strings"
	"time"
)

func (session *QsSession) getDeadHostsForScan(scanID string, created time.Time) (map[string]string, error) {
	var deadHostIPToProof = make(map[string]string)
	var output *qualys.ScanSummaryOutput
	var err error
	if output, err = session.apiSession.GatherDeadHostsFoundSince(created); err == nil {
		for _, scan := range output.Response.ScanSummaryList.ScanSummary {
			if scan.ScanRef == scanID {
				for _, host := range scan.HostSummary {
					var proofByte []byte
					if proofByte, err = xml.MarshalIndent(scan, "", "\t"); err == nil {

						// host.Text contains the IPs of the host in CSV
						var ipList = strings.Replace(host.Text, " ", "", -1) // remove spaces
						for _, host := range strings.Split(ipList, ",") {
							// proofByte contains the portion of the XML that displays the host as dead
							deadHostIPToProof[host] = string(proofByte)
						}

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
