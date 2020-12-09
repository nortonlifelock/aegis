package integrations

import (
	"context"
	"fmt"
	"github.com/nortonlifelock/aegis/pkg/aqua"
	"github.com/nortonlifelock/aegis/pkg/domain"
)

type IScanner interface {
	RescanImage(ctx context.Context, image string, registry string) (findings []domain.ImageFinding, err error)
	RescanContainer(ctx context.Context, namespace string) (findings []domain.ImageFinding, err error)
	CreateException(finding domain.ImageFinding, comment string) (err error)
}

const (
	Aqua = "Aqua"
)

func GetImageScanner(scannerID string, ms domain.DatabaseConnection, sourceConfig domain.SourceConfig, appConfig config, lstream logger) (client IScanner, err error) {
	var user, pass string
	user, pass, err = getUsernameAndPasswordFromEncryptedSourceConfig(ms, sourceConfig, appConfig)

	if err == nil {
		if len(scannerID) > 0 {
			switch scannerID {

			case Aqua:
				client, err = aqua.CreateClient(sourceConfig.Address(), user, pass, lstream)
			default:
				err = fmt.Errorf("unknown scanner type %s", scannerID)
			}
		} else {
			err = fmt.Errorf("empty scanner id passed to GetImageScanner")
		}
	}

	return client, err
}
