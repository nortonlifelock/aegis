package connector

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/nortonlifelock/log"
	"github.com/nortonlifelock/qualys"
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

func (session *QsSession) getProofForDeadHost(ip string, deadHostIPToProof map[string]string) (proof string) {
	proof = deadHostIPToProof[ip]
	if len(proof) == 0 {
		for deadHostIP, deadHostProof := range deadHostIPToProof {
			if session.ipIsInRange(ip, deadHostIP) {
				proof = deadHostProof
				break
			}
		}
	}

	return proof
}

// ipRange can either be a specific IP (100.0.0.0) or a range of IPs (100.0.0.0 - 100.0.0.100)
func (session *QsSession) ipIsInRange(ip, ipRange string) (match bool) {
	hyphenIndex := strings.Index(ipRange, "-")
	if hyphenIndex > 0 {

		firstRange := ipRange[0:hyphenIndex]
		secondRange := ipRange[hyphenIndex+1:]

		checkIP := net.ParseIP(ip)
		firstRangeIP := net.ParseIP(firstRange)
		secondRangeIP := net.ParseIP(secondRange)

		if firstRangeIP.To4() != nil && secondRangeIP.To4() != nil {
			if bytes.Compare(checkIP, firstRangeIP) >= 0 && bytes.Compare(checkIP, secondRangeIP) <= 0 {
				match = true
			}
		} else {
			session.lstream.Send(log.Errorf(nil, "either could not parse %s or %s as IPv4", firstRangeIP, secondRangeIP))
		}
	} else {
		if ip == ipRange {
			match = true
		}
	}

	return match
}
