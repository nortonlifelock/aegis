package qualys

import (
	"fmt"
	"github.com/nortonlifelock/log"
	"net/http"
	"strconv"
	"strings"
)

// GetScanList loads the list of scans from Qualys
func (session *Session) GetScanList() (output QScanListOutput, err error) {
	output = QScanListOutput{}

	var fields = make(map[string]string)
	fields["action"] = "list"

	// TODO: INTEL-386 Using the Scan List Filter Options in Qualys API update this to look for scans between certain times so that the API call take less time
	// Options: scan_ref, state, processed, type, target, user_login,
	// launched_after_datetime, launched_before_datetime,
	// scan_type=certview, client_id and client_name (only for Consultant type subscriptions)

	err = session.post(session.Config.Address()+qsVMScan, fields, &output)

	return output, err
}

// GetScanByReference queries the Qualys API and recovers information for the scan corresponding to the scanReference argument
func (session *Session) GetScanByReference(scanReference string) (scan ScanQualys, err error) {
	var output = QScanListOutput{}
	var fields = make(map[string]string)
	fields["action"] = "list"
	fields["scan_ref"] = scanReference

	// Options: scan_ref, state, processed, type, target, user_login,
	// launched_after_datetime, launched_before_datetime,
	// scan_type=certview, client_id and client_name (only for Consultant type subscriptions)

	if err = session.post(session.Config.Address()+qsVMScan, fields, &output); err == nil {
		if len(output.Response.Scans) == 1 {
			scan = output.Response.Scans[0]
		} else {
			err = fmt.Errorf("unexpected scan count [%d] returned for reference [%s]", len(output.Response.Scans), scanReference)
		}
	}

	return scan, err
}

// CreateScan executes the API call to Qualys to create the scan with all of the information required by the endpoint
func (session *Session) CreateScan(scanTitle string, optionProfileID string, appliances []string, networkID int, ips []string, external bool) (scanID int, scanRef string, err error) {
	// TODO: Move this
	const externalScanner = "External"

	var fields = make(map[string]string)

	if networkID > 0 {
		fields["ip_network_id"] = strconv.Itoa(networkID)
	}

	if external {
		fields["iscanner_name"] = externalScanner
	} else {
		// Ensure we have engines to scan with
		if len(appliances) > 0 {
			fields["iscanner_id"] = strings.Join(appliances, ",") // concat the appliance values together in a comma separated list for passing to the API
		} else {
			err = fmt.Errorf("no scan appliances available to scan with for ips [%s]", strings.Join(ips, ","))
		}
	}

	if err == nil {
		fields["action"] = "launch"
		fields["scan_title"] = scanTitle      // Setup the scan title
		fields["option_id"] = optionProfileID // Option profile id that's configured in the database
		fields["ip"] = strings.Join(ips, ",") // concat the ips together in a comma separated list for the API
		// Execute the post call to the API to create the scan
		var ret = &simpleReturn{}
		if err = session.post(session.Config.Address()+qsVMScan, fields, ret); err == nil {

			// Determine if there were items returned in the response
			if len(ret.Response.Items) > 0 {

				// Pull the ID from the list of items as well as the reference and setup the return values for the method
				for _, item := range ret.Response.Items {

					if item.Key == "ID" {
						// Error if the ID of the scan cannot be parsed as an Integer
						if scanID, err = strconv.Atoi(item.Value); err != nil {
							err = fmt.Errorf("error occurred while converting Qualys scan Id to INT [%s]", err.Error())
						}
					} else if item.Key == "REFERENCE" {
						scanRef = item.Value
					}
				}
			} else {
				err = fmt.Errorf("invalid item list returned from create search list in Qualys")
			}
		} else {
			err = fmt.Errorf("error when executing call to Qualys to initialize scan | %s", err.Error())
		}
	}

	return scanID, scanRef, err
}

func (session *Session) GetAssetTagTargetOfScheduledScan(scheduleTitle string) (tagSetTarget string, err error) {
	var output = ScheduleScanListOutput{}
	var fields = make(map[string]string)
	fields["action"] = "list"

	if err = session.httpCall(http.MethodGet, session.Config.Address()+qsScheduledScan, fields, nil, &output); err == nil {
		for _, scheduledScan := range output.Response.ScheduleScanList.Scan {
			if scheduledScan.Title == scheduleTitle {
				tagSetTarget = scheduledScan.AssetTags.TagSetInclude
			}
		}
	}

	if len(tagSetTarget) == 0 {
		err = fmt.Errorf("could not find asset tag target of scheduled scan [%s]", scheduleTitle)
	}

	return tagSetTarget, err
}

func (session *Session) GetScheduledScan(scanTitle string) (scan *ScanQualys, err error) {
	var output = QScanListOutput{}
	var fields = make(map[string]string)
	fields["action"] = "list"
	fields["type"] = "Scheduled"
	fields["state"] = "Running,Paused,Queued,Loading,Finished" // TODO remove finished after testing

	if err = session.post(session.Config.Address()+qsVMScan, fields, &output); err == nil {

		var found bool
		for index, scheduledScan := range output.Response.Scans {
			if scheduledScan.Title == scanTitle {
				found = true
				scan = &output.Response.Scans[index]
				break
			}
		}

		if found {
			if scan != nil {
				session.lstream.Send(log.Debugf("found a scheduled scan with title [%s|%s]", scanTitle, scan.Reference))
			}
		} else {
			session.lstream.Send(log.Debugf("there is currently not a running scheduled with title [%s]", scanTitle))
		}
	}

	return scan, err
}
