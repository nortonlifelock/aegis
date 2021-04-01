package integrations

import (
	"context"
	"encoding/json"

	"github.com/nortonlifelock/aegis/pkg/crypto"
	"github.com/nortonlifelock/aegis/pkg/domain"
	"github.com/nortonlifelock/aegis/pkg/dome9"
	qualys "github.com/nortonlifelock/aegis/pkg/qualys/connector"
	"github.com/pkg/errors"
)

const (
	// Dome9 delineates that the CIS scanner connects to Dome9
	Dome9     = "Dome9"
	CloudView = "CloudView"
)

// CISScanner finds compliance violations within a cloud service
type CISScanner interface {
	RescanBundle(bundleID string, cloudAccountID string) (findings []domain.Finding, err error)
}

// GetCISScanner returns a struct that implements the TicketingEngine interface
func GetCISScanner(ctx context.Context, scannerID string, ms domain.DatabaseConnection, sourceConfig domain.SourceConfig, appConfig config, lstream logger) (client CISScanner, err error) {
	var user, pass string
	user, pass, err = getUsernameAndPasswordFromEncryptedSourceConfig(ms, sourceConfig, appConfig)

	if err == nil {
		if len(scannerID) > 0 {
			switch scannerID {

			case Dome9:
				client, err = dome9.CreateClient(user, pass, sourceConfig.Address(), lstream)
				break
			case Qualys, CloudView:
				var decryptedConfig domain.SourceConfig
				decryptedConfig, err = crypto.DecryptSourceConfig(ms, sourceConfig, appConfig)
				if err == nil {
					client, err = qualys.Connect(ctx, lstream, decryptedConfig)
				}
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
