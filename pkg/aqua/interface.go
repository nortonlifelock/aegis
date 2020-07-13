package aqua

import (
	"context"
	"fmt"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/log"
)

const (
	finished = "finished"
	queued   = "queued"
	failed   = "failed"
)

func (cli *APIClient) RescanImage(ctx context.Context, repository string, registry string) (findings []domain.ImageFinding, err error) {
	findings = make([]domain.ImageFinding, 0)

	var images []ImageResult
	if images, err = cli.GetImagesForRepository(registry, repository); err == nil {
		mostRecentTag := cli.getMostRecentImageTag(images)
		if len(mostRecentTag) > 0 {

			//cli.lstream.Send(log.Infof("Scanning %s:%s in registry %s", repository, mostRecentTag, registry))
			//
			//if err = cli.StartFullImageRescan(registry, repository); err == nil {
			//	var status = queued
			//	for status != finished {
			//		select {
			//		case <-ctx.Done():
			//			return
			//		default:
			//		}
			//
			//		var scan *Scan
			//		if scan, err = cli.GetScanSummaries(ctx, registry, fmt.Sprintf("%s:%s", repository, mostRecentTag)); scan != nil && err == nil {
			//			status = scan.StatusVar
			//			cli.lstream.Send(log.Infof("scan for [%s|%s] has a status [%s]", repository, registry, status))
			//
			//			if status == failed {
			//				err = fmt.Errorf("scan failed")
			//				break
			//			} else if status != finished { // TODO what's the error status?
			//				cli.lstream.Send(log.Debugf("waiting 30 seconds"))
			//				time.Sleep(time.Second * 30)
			//			}
			//		} else {
			//			err = fmt.Errorf("error while grabbing scan for [%s|%s] - %s", repository, registry, err.Error())
			//			break
			//		}
			//	}
			//
			//	if err == nil && status == finished {
			//		cli.lstream.Send(log.Infof("loading vulnerabilities for [%s:%s|%s]", repository, mostRecentTag, registry))
			//		findings, err = cli.GetVulnerabilitiesForImage(ctx, fmt.Sprintf("%s:%s", repository, mostRecentTag), registry)
			//	}
			//} else {
			//	err = fmt.Errorf("error while creating rescan for [%s|%s]", repository, registry)
			//}

			cli.lstream.Send(log.Infof("loading vulnerabilities for [%s:%s|%s]", repository, mostRecentTag, registry))
			findings, err = cli.GetVulnerabilitiesForImage(ctx, fmt.Sprintf("%s:%s", repository, mostRecentTag), registry)
		} else {
			err = fmt.Errorf("error while gathering the most recent tag for [%s|%s] - %s", repository, registry, err.Error())
		}
	}

	return findings, err
}
