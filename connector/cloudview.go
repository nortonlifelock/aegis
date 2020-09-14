package connector

import "github.com/nortonlifelock/domain"

func (session *QsSession) RescanBundle(bundleID int, cloudAccountID string) (findings []domain.Finding, err error) {
	err = session.apiSession.GetCloudViewFindings(cloudAccountID)
	return nil, err
}
