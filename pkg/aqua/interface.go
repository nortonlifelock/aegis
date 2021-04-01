package aqua

import (
	"context"
	"fmt"
	"github.com/nortonlifelock/aegis/pkg/domain"
	"github.com/nortonlifelock/aegis/pkg/log"
)

const (
	finished = "finished"
	queued   = "queued"
	failed   = "failed"
)

func (cli *APIClient) RescanContainer(ctx context.Context, namespace string) (findings []domain.ImageFinding, err error) {
	var processedImage = make(map[string]bool)
	var containers []ContainerInfo
	if containers, err = cli.GetContainersForNamespace(ctx, namespace); err == nil {

		var exceptions []domain.ImageFinding
		if exceptions, err = cli.GetExceptions(ctx); err == nil {
			exceptionMap := mapFindingsByKey(exceptions)

			for _, container := range containers {
				select {
				case <-ctx.Done():
					return nil, fmt.Errorf("context closed")
				default:
				}

				if imageDigest := container.ImageID; !processedImage[imageDigest] {
					processedImage[imageDigest] = true

					var findingsForImage []*VulnerabilityResult
					findingsForImage, err = cli.GetVulnerabilitiesForImage(ctx, container.ImageName, "")

					for index := range findingsForImage {
						if findingsForImage[index].ImageDigest == imageDigest {
							if exceptionMap[getKeyForFinding(findingsForImage[index])] == nil {
								findings = append(findings, findingsForImage[index])
							} else {
								findings = append(findings, &exceptedFinding{findingsForImage[index]})
							}
						}
					}
				}
			}
		} else {
			err = fmt.Errorf("error while gathering exceptions - %s", err.Error())
		}

	} else {
		err = fmt.Errorf("error whiel gathering container info for [%s] - %s", namespace, err.Error())
	}

	return findings, err
}

func (cli *APIClient) RescanImage(ctx context.Context, repository string, registry string) (findings []domain.ImageFinding, err error) {
	findings = make([]domain.ImageFinding, 0)

	var images []ImageResult
	if images, err = cli.GetImagesForRepository(ctx, registry, repository); err == nil {
		mostRecentTag := cli.getMostRecentImageTag(images)
		if len(mostRecentTag) > 0 {

			cli.lstream.Send(log.Infof("loading exceptions"))
			var exceptions []domain.ImageFinding
			if exceptions, err = cli.GetExceptions(ctx); err == nil {
				exceptionMap := mapFindingsByKey(exceptions)

				cli.lstream.Send(log.Infof("loading vulnerabilities for [%s:%s|%s]", repository, mostRecentTag, registry))

				var findingsForImage []*VulnerabilityResult
				findingsForImage, err = cli.GetVulnerabilitiesForImage(ctx, fmt.Sprintf("%s:%s", repository, mostRecentTag), registry)

				for _, findingForImage := range findingsForImage {
					if exceptionMap[getKeyForFinding(findingForImage)] == nil {
						findings = append(findings, findingForImage)
					} else {
						findings = append(findings, &exceptedFinding{findingForImage})
					}
				}
			} else {
				err = fmt.Errorf("error while loading exceptions - %s", err.Error())
			}
		} else {
			err = fmt.Errorf("could not find any tags for [%s|%s]", repository, registry)
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

type exceptedFinding struct {
	domain.ImageFinding
}

func (e *exceptedFinding) Exception() bool {
	return true
}
