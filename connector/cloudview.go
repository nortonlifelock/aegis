package connector

import (
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/qualys"
	"sync"
)

func (session *QsSession) RescanBundle(bundleID int, cloudAccountID string) (findings []domain.Finding, err error) {
	findings = make([]domain.Finding, 0)

	var evaluations []qualys.AccountEvaluationContent
	if evaluations, err = session.apiSession.GetCloudAccountEvaluations(cloudAccountID); err == nil {

		var wg sync.WaitGroup
		permit := getPermitThread(10)
		var lock sync.Mutex

		for _, evaluation := range evaluations {
			<-permit
			wg.Add(1)
			go func(evaluation qualys.AccountEvaluationContent) {
				defer wg.Done()
				defer func() {
					permit <- true
				}()

				if evaluationFindings, threadErr := session.apiSession.GetCloudEvaluationFindings(cloudAccountID, evaluation); threadErr == nil {
					lock.Lock()
					findings = append(findings, evaluationFindings...)
					lock.Unlock()
				} else {
					err = threadErr // only one error needs to make it out, so overwrite isn't concerning
				}
			}(evaluation)
		}
		wg.Wait()
	}

	return findings, err
}

func getPermitThread(simultaneousCount int) (permit chan bool) {
	permit = make(chan bool, simultaneousCount)
	for i := 0; i < simultaneousCount; i++ {
		permit <- true
	}

	return permit
}
