package integrations

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/nortonlifelock/crypto"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/dome9"
)

const (
	// Dome9 delineates that the CIS scanner connects to Dome9
	Dome9 = "Dome9"
)

// CISScanner finds compliance violations within a cloud service
type CISScanner interface {
	//GetAllFindingsForBundle(bundleID int) (retChan <-chan domain.Finding, err error)

	RescanBundle(bundleID int, cloudAccountID string) (findings []domain.Finding, err error)
}

// GetCISScanner returns a struct that implements the TicketingEngine interface
func GetCISScanner(scannerID string, ms domain.DatabaseConnection, sourceConfig domain.SourceConfig, appConfig config, lstream logger) (client CISScanner, err error) {
	var user, pass string
	user, pass, err = getUsernameAndPasswordFromEncryptedSourceConfig(ms, sourceConfig, appConfig)

	if err == nil {
		if len(scannerID) > 0 {
			switch scannerID {

			case Dome9:
				client, err = dome9.CreateClient(user, pass, sourceConfig.Address(), lstream)
				break
			default:
				err = errors.Errorf("Unknown scanner type %s", scannerID)
			}
		} else {
			err = errors.New("empty scanner id passed to GetCISScanner")
		}
	}

	return client, err
}

func getUsernameAndPasswordFromEncryptedSourceConfig(ms domain.DatabaseConnection, sourceConfig domain.SourceConfig, appConfig config) (user, password string, err error) {
	if sourceConfig, err = crypto.DecryptSourceConfig(ms, sourceConfig, appConfig); err == nil {

		var authInfo domain.BasicAuth
		if err = json.Unmarshal([]byte(sourceConfig.AuthInfo()), &authInfo); err == nil {
			user, password = authInfo.Username, authInfo.Password
		}
	}

	return user, password, err
}
