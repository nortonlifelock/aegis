package qualys

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// CreateOptionProfile takes an option profile object outlining the configurable aspects of a scan, and creates that
// in Qualys, and returns its ID
func (session *Session) CreateOptionProfile(optionProfile *OptionProfiles) (optionProfileID string, err error) {
	var fields = make(map[string]string)
	fields["action"] = "import"

	if len(optionProfile.OptionProfile.Scan.Ports.TCPPorts.TCPPortsAdditional.HasAdditional) == 0 {
		optionProfile.OptionProfile.Scan.Ports.TCPPorts.TCPPortsAdditional.HasAdditional = "0"
	}

	if len(optionProfile.OptionProfile.Scan.Ports.UDPPorts.UDPPortsAdditional.HasAdditional) == 0 {
		optionProfile.OptionProfile.Scan.Ports.UDPPorts.UDPPortsAdditional.HasAdditional = "0"
	}

	if len(optionProfile.OptionProfile.Scan.DissolvableAgent.DissolvableAgentEnable) == 0 {
		optionProfile.OptionProfile.Scan.DissolvableAgent.DissolvableAgentEnable = "0"
	}

	if len(optionProfile.OptionProfile.Scan.DissolvableAgent.WindowsShareEnumerationEnable) == 0 {
		optionProfile.OptionProfile.Scan.DissolvableAgent.WindowsShareEnumerationEnable = "0"
	}

	if len(optionProfile.OptionProfile.MAP.UDPPorts.UDPPortsStandardScan) == 0 {
		optionProfile.OptionProfile.MAP.UDPPorts.UDPPortsStandardScan = "0"
	}

	var body []byte
	body, err = xml.MarshalIndent(&optionProfile, "", "\t")
	if err == nil {
		var resp = &simpleReturn{}
		var bodyString = string(body)
		err = session.httpCall(http.MethodPost, session.Config.Address()+qsOptionProfile, fields, &bodyString, resp)
		if err == nil {
			for _, item := range resp.Response.Items {
				if len(item.Key) > 0 {
					optionProfileID = item.Key
					break
				}
			}

			if len(optionProfileID) == 0 {
				err = fmt.Errorf("failed to grab the ID of the newly created option profile")
			}
		}
	}

	return optionProfileID, err
}

// GetOptionProfile returns the option profile in Qualys corresponding to the ID in the argument
func (session *Session) GetOptionProfile(optionProfileID int) (optionProfiles *OptionProfiles, err error) {
	var fields = make(map[string]string)
	fields["action"] = "export"
	fields["option_profile_id"] = strconv.Itoa(optionProfileID)
	optionProfiles = &OptionProfiles{}
	err = session.httpCall(http.MethodGet, session.Config.Address()+qsOptionProfile, fields, nil, &optionProfiles)

	return optionProfiles, err
}

// CreateSearchList creates a search list in Qualys which specifies the vulnerabilities for Qualys to scan
func (session *Session) CreateSearchList(qIDs []string, searchListFormatString string) (searchListID string, err error) {
	const idKey = "id"

	var fields = make(map[string]string)
	fields["action"] = "create"

	// TODO what do we want the search list title convention to be?
	// TODO should it be configurable?
	fields["title"] = fmt.Sprintf(searchListFormatString, time.Now().Nanosecond())

	fields["qids"] = strings.Join(qIDs, ",")
	var response = &simpleReturn{}

	err = session.post(session.Config.Address()+qsSearchList, fields, response)
	if err == nil {
		for _, item := range response.Response.Items {
			if strings.ToLower(item.Key) == idKey {
				searchListID = item.Value
				break
			}
		}
	}

	return searchListID, err
}

// DeleteSearchList calls the Qualys endpoint to delete a search list (which specifies the vulnerabilities to
// be scanned by Qualys)
func (session *Session) DeleteSearchList(searchListID string) (err error) {
	var fields = make(map[string]string)
	fields["action"] = "delete"
	fields["id"] = searchListID
	err = session.post(session.Config.Address()+qsSearchList, fields, nil)
	return err
}

// DeleteOptionProfile calls the Qualys endpoint to delete an option profile
func (session *Session) DeleteOptionProfile(optionProfileID string) (err error) {
	var fields = make(map[string]string)
	fields["action"] = "delete"
	fields["id"] = optionProfileID
	err = session.post(session.Config.Address()+qsOptionProfileDelete, fields, nil)
	return err
}

// GatherDeadHostsFoundSince returns a list of hosts that were found dead by a scan that started since a certain date
func (session *Session) GatherDeadHostsFoundSince(since time.Time) (output *ScanSummaryOutput, err error) {
	if !since.IsZero() {
		output = &ScanSummaryOutput{}

		fields := make(map[string]string)
		// Required fields
		fields["action"] = "list"
		fields["scan_date_since"] = since.Format("2006-01-02")

		// Optional fields
		fields["include_dead"] = "1"
		fields["include_excluded"] = "0"
		fields["include_unresolved"] = "0"
		fields["include_cancelled"] = "0"
		fields["include_notvuln"] = "0"
		fields["include_blocked"] = "0"
		fields["include_duplicate"] = "0"
		fields["include_aborted"] = "0"

		err = session.httpCall(http.MethodGet, session.Config.Address()+qsHostStatusFromScan, fields, nil, output)
	} else {
		err = fmt.Errorf("invalid since date")
	}

	return output, err
}
