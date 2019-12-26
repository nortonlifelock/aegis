package nexpose

import (
	"context"
	"fmt"
	"net/http"
	"github.com/nortonlifelock/log"
)

// GetAssetDetections loads all vulnerability detections for the asset ID that is passed.
// The returned data can be sorted by passing a sort string that uses the formatting available in the nexpose API.
func (a *Session) GetAssetDetections(ctx context.Context, assetID string, sort string) (<-chan *Finding, error) {
	var detections = make(chan *Finding)
	var err error

	go func(detections chan<- *Finding) {
		defer handleRoutinePanic(a.lstream)
		defer close(detections)

		fields := map[string]string{"sort": sort}

		var data <-chan interface{}
		if data, err = a.getPagedData(ctx, &PageOfFinding{}, fmt.Sprintf(apiGetAssetVulnerabilities, encode(assetID)), fields); err == nil {
			for {
				select {
				case <-ctx.Done():
					return
				case d, ok := <-data:
					if ok {
						var detection Finding
						if detection, ok = d.(Finding); ok {
							detections <- &detection
						} else {
							a.lstream.Send(log.Error("unable to cast paged return as detection", err))
						}
					} else {
						return
					}
				}
			}
		} else {
			a.lstream.Send(log.Error("error executing paged call against nexpose", err))
		}
	}(detections)

	return detections, err
}

// GetAssetDetection loads the specific detection results from nexpose api for the asset ID and vulnerability ID that
// are passed to the method
func (a *Session) GetAssetDetection(assetID string, vulnerabilityID string) (detection *Finding, err error) {

	detection = &Finding{}
	err = a.execute(http.MethodGet, fmt.Sprintf(apiGetAssetVulnerabilityDetail, encode(assetID), encode(vulnerabilityID)), nil, nil, detection)

	return detection, err
}
