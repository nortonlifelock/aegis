package nexpose

import (
	"context"
	"fmt"
	"github.com/nortonlifelock/log"
	"net/http"
	"sync"
)

// GetSolutions loads all vulnerability solutions from the nexpose api
// The returned data can be sorted by passing a sort string that uses the formatting available in the nexpose API.
func (a *Session) GetSolutions(ctx context.Context, sort string) <-chan *Solution {
	var solutions = make(chan *Solution)

	go func(solutions chan<- *Solution) {
		defer handleRoutinePanic(a.lstream)
		defer close(solutions)
		var err error

		fields := map[string]string{"sort": sort}

		var data <-chan interface{}
		if data, err = a.getPagedData(ctx, &PageOfSolution{}, apiGetSolutions, fields); err == nil {
			for {
				select {
				case <-ctx.Done():
					return
				case d, ok := <-data:
					if ok {
						var solution Solution
						if solution, ok = d.(Solution); ok {
							solutions <- &solution
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
	}(solutions)

	return solutions
}

// GetSolution loads the solution details for the passed solution id
func (a *Session) GetSolution(vulnerabilityID string) (solution *Solution, err error) {
	solution = &Solution{}
	err = a.execute(http.MethodGet, fmt.Sprintf(apiGetSolution, encode(vulnerabilityID)), nil, nil, solution)
	return solution, err
}

// GetVulnerabilitySolutions loads the solutions for the vulnerability ID that is passed in
func (a *Session) GetVulnerabilitySolutions(ctx context.Context, vulnerabilityID string) (<-chan *Solution, error) {
	var solutions = make(chan *Solution)
	var err error

	go func(solutions chan<- *Solution) {
		defer handleRoutinePanic(a.lstream)
		defer close(solutions)
		var wg = sync.WaitGroup{}

		results := &Reference{}
		if err = a.execute(http.MethodGet, fmt.Sprintf(apiGetVulnerabilitySolutions, vulnerabilityID), nil, nil, results); err == nil {
			for _, result := range results.Resources {

				wg.Add(1)
				go func(result string) {
					defer handleRoutinePanic(a.lstream)
					defer wg.Done()

					var solution *Solution
					if solution, err = a.GetSolution(result); err == nil {
						select {
						case <-ctx.Done():
							return
						case solutions <- solution:
							return
						}
					} else {
						a.lstream.Send(log.Errorf(err, "Error getting solution for vulnerability [%s] by Id [%s]", vulnerabilityID, result))
					}
				}(result)
			}
		} else {
			a.lstream.Send(log.Errorf(err, "Error getting solutions for vulnerability [%s]", vulnerabilityID))
		}

		wg.Wait()

	}(solutions)

	return solutions, err
}

// TODO get solution for specific asset
