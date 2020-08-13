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

			cli.lstream.Send(log.Infof("loading exceptions"))
			var exceptions []domain.ImageFinding
			if exceptions, err = cli.GetExceptions(ctx); err == nil {
				exceptionMap := mapFindingsByKey(exceptions)

				cli.lstream.Send(log.Infof("loading vulnerabilities for [%s:%s|%s]", repository, mostRecentTag, registry))

				var unfilteredFindings []domain.ImageFinding
				unfilteredFindings, err = cli.GetVulnerabilitiesForImage(ctx, fmt.Sprintf("%s:%s", repository, mostRecentTag), registry)

				for _, unfilteredFinding := range unfilteredFindings {
					if exceptionMap[getKeyForFinding(unfilteredFinding)] == nil {
						findings = append(findings, unfilteredFinding)
					} else {
						cli.lstream.Send(log.Infof("skipping [%s] as it has an exception in Aqua", getKeyForFinding(unfilteredFinding)))
					}
				}
			} else {
				err = fmt.Errorf("error while loading exceptions - %s", err.Error())
			}
		} else {
			err = fmt.Errorf("error while gathering the most recent tag for [%s|%s]", repository, registry)
		}
	}

	return findings, err
}

func getKeyForFinding(finding domain.ImageFinding) (key string) {
	key = fmt.Sprintf("%s;%s;%s;%s", finding.ImageName(), finding.Registry(), finding.VulnerabilityLocation(), finding.VulnerabilityID())
	return key
}

func mapFindingsByKey(findings []domain.ImageFinding) (keyToFinding map[string]domain.ImageFinding) {
	keyToFinding = make(map[string]domain.ImageFinding)

	for index := range findings {
		finding := findings[index]
		keyToFinding[getKeyForFinding(finding)] = finding
	}

	return keyToFinding
}
