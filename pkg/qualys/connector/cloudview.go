package connector

import (
	"github.com/nortonlifelock/aegis/pkg/domain"
	"github.com/nortonlifelock/aegis/pkg/qualys"
	"sync"
)

func (session *QsSession) RescanBundle(policyName string, cloudAccountID string) (findings []domain.Finding, err error) {
	findings = make([]domain.Finding, 0)

	var evaluations []qualys.AccountEvaluationContent
	var cloudAccountType string
	if evaluations, cloudAccountType, err = session.apiSession.GetCloudAccountEvaluations(cloudAccountID); err == nil {

		var wg sync.WaitGroup
		permit := getPermitThread(10)
		var lock sync.Mutex

		for index := range evaluations {
			wg.Add(1)
			<-permit
			go func(evaluation qualys.AccountEvaluationContent) {
				defer func() {
					permit <- true
					wg.Done()
				}()

				if evaluationFindings, threadErr := session.apiSession.GetCloudEvaluationFindings(cloudAccountID, evaluation, policyName, cloudAccountType); threadErr == nil {
					lock.Lock()
					findings = append(findings, evaluationFindings...)
					lock.Unlock()
				} else {
					err = threadErr // only one error needs to make it out, so overwrite isn't concerning
				}
			}(evaluations[index])
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
