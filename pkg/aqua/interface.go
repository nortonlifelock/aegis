package aqua

import (
	"context"
	"fmt"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/log"
	"time"
)

func (cli *apiClient) RescanImage(ctx context.Context, image string, groupID string) (findings []domain.ImageFinding, err error) {
	findings = make([]domain.ImageFinding, 0)

	const (
		finished = "finished"
	)

	if err = cli.StartFullImageRescan(groupID, image); err == nil {
		var status = "queued"
		for status != finished {
			select {
			case <-ctx.Done():
				return
			default:
			}

			var scan *Scan
			if scan, err = cli.GetScanSummaries(ctx, groupID, fmt.Sprintf("%s:latest", image)); scan != nil && err == nil {
				status = scan.StatusVar
				cli.lstream.Send(log.Infof("scan for [%s|%s] has a status [%s]", image, groupID, status))

				if status != finished { // TODO what's the error status?
					cli.lstream.Send(log.Debugf("waiting 30 seconds"))
					time.Sleep(time.Second * 30)
				}
			} else {
				err = fmt.Errorf("error while grabbing scan for [%s|%s]", image, groupID)
				break
			}
		}

		if err == nil && status == finished {
			cli.lstream.Send(log.Infof("loading vulnerabilities for [%s|%s]", image, groupID))
			findings, err = cli.GetVulnerabilitiesForImage(ctx, fmt.Sprintf("%s:latest", image), groupID)
		}
	} else {
		err = fmt.Errorf("error while creating rescan for [%s|%s]", image, groupID)
	}

	return findings, err
}
