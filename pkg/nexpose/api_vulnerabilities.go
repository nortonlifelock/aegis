package nexpose

import (
	"context"
	"fmt"
	"github.com/nortonlifelock/aegis/pkg/log"
	"net/http"
	"sync"
)

// TODO: add input validation

// GetVulnerabilities loads all vulnerabilities that nexpose scans for from the api and returns them
// over a channel back to the calling method. The returned data can be sorted by passing a sort string
// that uses the formatting available in the nexpose API.
func (a *Session) GetVulnerabilities(ctx context.Context, sort string) <-chan *Vulnerability {
	var vulnerabilities = make(chan *Vulnerability)

	go func(vulnerabilities chan<- *Vulnerability) {
		defer handleRoutinePanic(a.lstream)
		defer close(vulnerabilities)
		var err error

		var fields = make(map[string]string)
		if len(sort) > 0 {
			fields["sort"] = sort
		}

		var data <-chan interface{}
		if data, err = a.getPagedData(ctx, &PageOfVulnerability{}, apiGetVulnerabilities, fields); err == nil {
			for {
				select {
				case <-ctx.Done():
					return
				case d, ok := <-data:
					if ok {
						var vulnerability Vulnerability
						if vulnerability, ok = d.(Vulnerability); ok {
							vulnerabilities <- &vulnerability
						} else {
							a.lstream.Send(log.Error("unable to cast paged return as vulnerability", err))
						}
					} else {
						return
					}
				}
			}
		} else {
			a.lstream.Send(log.Error("error executing paged call against nexpose", err))
		}
	}(vulnerabilities)

	return vulnerabilities
}

// GetVulnerability pulls the vulnerability data from nexpose using the vulnerability ID passed in.
func (a *Session) GetVulnerability(vulnerabilityID string) (vulnerability *Vulnerability, err error) {

	vulnerability = &Vulnerability{}
	err = a.execute(http.MethodGet, fmt.Sprintf(apiGetVulnerability, encode(vulnerabilityID)), nil, nil, vulnerability)

	return vulnerability, err
}

// GetVulnerabilityReferences loads the references specific to the vulnerability ID passed in.
// The returned data can be sorted by passing a sort string that uses the formatting available in the nexpose API.
func (a *Session) GetVulnerabilityReferences(ctx context.Context, vulnerabilityID string, sort string) (<-chan *VulnerabilityReference, error) {
	var references = make(chan *VulnerabilityReference)
	var err error

	go func(references chan<- *VulnerabilityReference) {
		defer handleRoutinePanic(a.lstream)
		defer close(references)

		fields := map[string]string{"sort": sort}

		var data <-chan interface{}
		if data, err = a.getPagedData(ctx, &PageOfVulnerabilityReference{}, fmt.Sprintf(apiGetVulnerabilityReferences, vulnerabilityID), fields); err == nil {
			for {
				select {
				case <-ctx.Done():
					return
				case d, ok := <-data:
					if ok {
						var ref VulnerabilityReference
						if ref, ok = d.(VulnerabilityReference); ok {
							references <- &ref
						} else {
							a.lstream.Send(log.Error("unable to cast paged return as vulnerability reference", err))
						}
					} else {
						return
					}
				}
			}
		} else {
			a.lstream.Send(log.Error("error executing paged call against nexpose", err))
		}
	}(references)

	return references, err
}

// GetChecksForVulnerabilities loads all checks for the vulnerabilities passed in
func (a *Session) GetChecksForVulnerabilities(ctx context.Context, vulnerabilityIDs []string) (<-chan *Check, error) {
	var checks = make(chan *Check)
	var err error

	go func(checks chan<- *Check) {
		defer handleRoutinePanic(a.lstream)
		defer close(checks)
		var wg = sync.WaitGroup{}

		var seen = make(map[string]bool)
		var seenLock = sync.Mutex{}

		for _, ID := range vulnerabilityIDs {
			seenLock.Lock()
			isNew := !seen[ID]
			seenLock.Unlock()

			if isNew {
				select {
				case <-ctx.Done():
					return
				default:
					wg.Add(1)

					go func(ID string) {
						defer handleRoutinePanic(a.lstream)
						defer wg.Done()

						var vchecks <-chan *Check
						if vchecks, err = a.GetVulnerabilityChecks(ctx, ID); err == nil {

							// Push the results onto the fan in channel
							for {
								select {
								case <-ctx.Done():
									return
								case check, ok := <-vchecks:
									if ok {
										checks <- check
									} else {
										return
									}
								}
							}

						} else {
							a.lstream.Send(log.Errorf(err, "error encountered while pulling checks for vulnerability [%s]", ID))
						}
					}(ID)
				}
			}
		}

		wg.Wait()
	}(checks)

	return checks, err
}

// GetChecksIDsForVulnerabilities loads all checks for the vulnerabilities passed in
func (a *Session) GetChecksIDsForVulnerabilities(ctx context.Context, vulnerabilityIDs []string) (<-chan string, error) {
	var checks = make(chan string)
	var err error

	go func(checks chan<- string) {
		defer handleRoutinePanic(a.lstream)
		defer close(checks)
		var wg = sync.WaitGroup{}

		var seen = make(map[string]bool)
		var seenLock = sync.Mutex{}

		for _, ID := range vulnerabilityIDs {
			seenLock.Lock()
			isNew := !seen[ID]
			seenLock.Unlock()

			if isNew {
				select {
				case <-ctx.Done():
					return
				default:
					wg.Add(1)

					go func(ID string) {
						defer handleRoutinePanic(a.lstream)
						defer wg.Done()

						var vchecks []string
						if vchecks, err = a.GetVulnerabilityCheckIDs(ID); err == nil {

							// Push the results onto the fan in channel
							for _, check := range vchecks {
								select {
								case <-ctx.Done():
									return
								case checks <- check:
								}
							}

						} else {
							a.lstream.Send(log.Errorf(err, "error encountered while pulling checks for vulnerability [%s]", ID))
						}
					}(ID)
				}
			}
		}

		wg.Wait()
	}(checks)

	return checks, err
}

// GetVulnerabilityChecks loads the checks used to verify the existence of a vulnerability
func (a *Session) GetVulnerabilityChecks(ctx context.Context, vulnerabilityID string) (<-chan *Check, error) {
	var checks = make(chan *Check)
	var err error

	go func(checks chan<- *Check) {
		defer handleRoutinePanic(a.lstream)
		defer close(checks)
		var wg = sync.WaitGroup{}

		var ids []string
		if ids, err = a.GetVulnerabilityCheckIDs(encode(vulnerabilityID)); err == nil {

			for _, result := range ids {

				wg.Add(1)
				go func(result string) {
					defer handleRoutinePanic(a.lstream)
					defer wg.Done()

					var solution *Check
					if solution, err = a.GetVulnerabilityCheck(encode(result)); err == nil {
						select {
						case <-ctx.Done():
							return
						case checks <- solution:
						}
					}
				}(result)
			}
		}

		wg.Wait()

	}(checks)

	return checks, err
}

// GetVulnerabilityCheckIDs loads the check ids used to verify the existence of a vulnerability
func (a *Session) GetVulnerabilityCheckIDs(vulnerabilityID string) (ids []string, err error) {

	results := &Reference{}
	if err = a.execute(http.MethodGet, fmt.Sprintf(apiGetVulnerabilityChecks, encode(vulnerabilityID)), nil, nil, results); err == nil {
		ids = results.Resources
	}

	return ids, err
}

// GetVulnerabilityCheck loads a specific vulnerability check by the ID passed in
func (a *Session) GetVulnerabilityCheck(vulnerabilityID string) (check *Check, err error) {

	check = &Check{}
	err = a.execute(http.MethodGet, fmt.Sprintf(apiGetVulnerabilityCheck, encode(vulnerabilityID)), nil, nil, check)

	return check, err
}
