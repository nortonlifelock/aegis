package qualys

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nortonlifelock/log"
)

// getAvailableAppliances loads Appliances from the Qualys API returning the Network ID of the appliances and filters out
// the appliances that are NOT "Online"
func (session *Session) getAvailableAppliances(engine string) (networkID int, appliances []string, err error) {
	appliances = make([]string, 0)
	var output = &QAppliances{}

	// Map query fields for request
	var fields = make(map[string]string)
	fields["action"] = "list"
	fields["ids"] = engine

	// Execute POST call to Qualys API for appliance informatino
	if err = session.post(session.Config.Address()+qsAppliance, fields, output); err == nil {

		// Process the appliances returned from the API
		networkID, appliances = session.processApplianceResults(output)

		if len(appliances) == 0 {
			err = fmt.Errorf("engine %s did not appear to be online", engine)
		}
	}

	return networkID, appliances, err
}

// GetApplianceInformation loads the information for the appliance ids that are passed in
func (session *Session) GetApplianceInformation(appliances []string) (output *QAppliances, err error) {
	var fields = make(map[string]string)
	fields["action"] = "list"
	fields["ids"] = strings.Join(appliances, ",")

	output = &QAppliances{}
	err = session.post(session.Config.Address()+qsAppliance, fields, output)

	return output, err
}

// processApplianceResults reads the output from the Appliance Endpoint and creates slices of the available appliances
// and returns the network ID for those appliances
func (session *Session) processApplianceResults(output *QAppliances) (networkID int, appliances []string) {
	if output != nil && len(output.Appliances) > 0 {

		// Filter out appliances that are not "Online" since they cannot be used
		for applianceID := range output.Appliances {
			if output.Appliances[applianceID].Status == "Online" {

				// Return scan appliances based on the network Id
				if networkID != output.Appliances[applianceID].NetworkID {
					networkID = output.Appliances[applianceID].NetworkID
				}

				appliances = append(appliances, strconv.Itoa(output.Appliances[applianceID].ID))

			} else {
				session.lstream.Send(log.Warningf(nil, "Appliance [%s] has a status of [%s] and is not available for use", output.Appliances[applianceID].Name, output.Appliances[applianceID].Status))
			}
		}
	}

	return networkID, appliances
}
