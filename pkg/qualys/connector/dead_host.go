package connector

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/nortonlifelock/aegis/pkg/qualys"
	"net"
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

							var allIPsInRange = getAllIPsInRange(host)
							for _, ip := range allIPsInRange {
								// proofByte contains the portion of the XML that displays the host as dead
								deadHostIPToProof[ip] = string(proofByte)
							}
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

// returns all IPs within a range of IPs (e.g. 100.0.0.0 - 100.0.0.100). If it is not a range, returns only the input IP
func getAllIPsInRange(ipRange string) (allIPsInRange []string) {

	hyphenIndex := strings.Index(ipRange, "-")
	if hyphenIndex > 0 {

		firstRange := ipRange[0:hyphenIndex]
		secondRange := ipRange[hyphenIndex+1:]

		traverseRangeIP := net.ParseIP(firstRange)
		secondRangeIP := net.ParseIP(secondRange)

		const (
			leftmostTuple = iota + 12
			secondTuple
			thirdTuple
			rightmostTuple
		)

		for bytes.Compare(traverseRangeIP, secondRangeIP) <= 0 {
			allIPsInRange = append(allIPsInRange, net.ParseIP(fmt.Sprintf("%v.%v.%v.%v", traverseRangeIP[leftmostTuple], traverseRangeIP[secondTuple], traverseRangeIP[thirdTuple], traverseRangeIP[rightmostTuple])).String())

			traverseRangeIP[rightmostTuple]++
			if traverseRangeIP[rightmostTuple] == 0 {

				traverseRangeIP[thirdTuple]++
				if traverseRangeIP[thirdTuple] == 0 {

					traverseRangeIP[secondTuple]++
					if traverseRangeIP[secondTuple] == 0 {

						traverseRangeIP[leftmostTuple]++
					}
				}
			}
		}
	} else {
		allIPsInRange = append(allIPsInRange, ipRange)
	}

	return allIPsInRange
}
