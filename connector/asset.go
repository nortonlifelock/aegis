package connector

import (
	"context"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/nexpose"
	"strconv"
	"sync"
)

type asset struct {
	conn  *Connection
	asset *nexpose.Asset
}

func (a *asset) OS() (param string) {
	return a.asset.OS
}

func (a *asset) HostName() (param string) {
	return a.asset.HostName
}

func (a *asset) MAC() (param string) {
	return a.asset.MAC
}

func (a *asset) IP() (param string) {
	return a.asset.IP
}

func (a *asset) SourceID() *string {
	id := strconv.Itoa(a.asset.ID)
	return &id
}

// Vulnerabilities returns a channel of detections for this asset from nexpose
func (a *asset) Vulnerabilities(ctx context.Context) (<-chan domain.Detection, error) {
	var detections = make(chan domain.Detection)
	var err error

	go func(detections chan<- domain.Detection) {
		defer handleRoutinePanic(a.conn.logger)
		defer close(detections)

		// Query the detections from the api
		var findings <-chan *nexpose.Finding
		if findings, err = a.conn.api.GetAssetDetections(ctx, sord(a.SourceID()), ""); err == nil {
			var wg = sync.WaitGroup{}

			for {
				select {
				case <-ctx.Done():
					return
				case finding, ok := <-findings:
					if ok {

						// Determine if the returned finding is a valid finding for listing detections
						if finding != nil && finding.Results != nil && len(finding.Results) > 0 {

							// Split off a go routine to loop through the results rather than
							// wait for them to be processed sequentially
							wg.Add(1)
							go func() {
								defer handleRoutinePanic(a.conn.logger)
								defer wg.Done()

								// Push the results onto the channel back to the requester
								for _, result := range finding.Results {

									// Push the wrapped detection onto the channel
									select {
									case <-ctx.Done():
										return
									case detections <- &detection{
										asset:           a,
										conn:            a.conn,
										vulnerabilityID: finding.ID,
										resultID:        result.ID,
										status:          finding.Status,
										detected:        finding.Since,
										proof:           result.Proof,
										port:            int(result.Port),
										protocol:        result.Protocol,
										updated:         finding.Since,
									}:
									}
								}
							}()
						}

					} else {
						// only return from this once the findings have all been processed
						wg.Wait()
						return
					}
				}
			}
		}
	}(detections)

	return detections, err
}

func (a *asset) ID() string {
	return ""
}

func (a *asset) Region() *string {
	return nil
}

func (a *asset) InstanceID() *string {
	return nil
}
