package integrations

import (
	"context"
	"fmt"
	"github.com/nortonlifelock/aegis/pkg/blackduck"
	"github.com/nortonlifelock/aegis/pkg/domain"
)

type CodeScanner interface {
	GetProjectVulnerabilities(ctx context.Context, projectID string) (findings []domain.CodeFinding, err error)
}

const (
	BlackDuck = "Black Duck"
)

func GetCodeScanner(scannerID string, ms domain.DatabaseConnection, sourceConfig domain.SourceConfig, appConfig config, lstream logger) (client CodeScanner, err error) {
	var pass string
	_, pass, err = getUsernameAndPasswordFromEncryptedSourceConfig(ms, sourceConfig, appConfig)

	if err == nil {
		if len(scannerID) > 0 {
			switch scannerID {

			case BlackDuck:
				client, err = blackduck.NewBlackDuckClient(sourceConfig.Address(), pass, true)
			default:
				err = fmt.Errorf("unknown scanner type %s", scannerID)
			}
		} else {
			err = fmt.Errorf("empty scanner id passed to GetCodeScanner")
		}
	}

	return client, err
}
