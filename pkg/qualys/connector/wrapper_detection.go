package connector

import (
	"context"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/log"
	"github.com/nortonlifelock/qualys"
	"strconv"
	"sync"
	"time"
)

type detection struct {
	d                 qualys.QDetection
	vulnerabilityInfo *vulnerabilityInfo

	session *QsSession
	lock    sync.Mutex
}

func (detection *detection) lazyLoadVulnerabilityInfoForDetection() {
	detection.lock.Lock()
	defer detection.lock.Unlock()

	needToLoad := false
	if detection.vulnerabilityInfo == nil {
		needToLoad = true
	} else if detection.vulnerabilityInfo.v == nil {
		needToLoad = true
	}

	if needToLoad {
		detection.vulnerabilityInfo = lazyLoadVulnerabilityInfo(detection.d.QualysID, detection.session)
	}
}

func lazyLoadVulnerabilityInfo(qid int, session *QsSession) (vi *vulnerabilityInfo) {
	session.vulnerabilityLock.Lock()
	if session.vulnerabilities[qid] == nil {
		session.vulnerabilityLock.Unlock()

		var err error
		vulnInfo := &vulnerabilityInfo{}
		vulnInfo.v, err = session.apiSession.LoadVulnerability(strconv.Itoa(qid))

		session.vulnerabilityLock.Lock()
		session.vulnerabilities[qid] = vulnInfo.v
		session.vulnerabilityLock.Unlock()

		if err == nil {
			vi = vulnInfo
		} else {
			session.lstream.Send(log.Errorf(err, "error while loading vulnerability information for detection [%v]", qid))
		}
	} else {
		vi = &vulnerabilityInfo{
			v: session.vulnerabilities[qid],
		}
		session.vulnerabilityLock.Unlock()
	}

	return vi
}

// TODO what to do if the vulnInfo is nil when we try and access it? should we continually attempt to load the vulnerability info?

func (detection *detection) ID() string {
	return ""
}

// SourceID returns the vulnerability ID of the detection
func (detection *detection) SourceID() string {
	return strconv.Itoa(detection.d.QualysID)
}

// Updated must return the time on the detection, not the time that the vulnerability was updated
// it is helpful to know the date of the detection in the context of the history of a device (e.g. a decommission date)
func (detection *detection) Updated() time.Time {
	return detection.d.LastFound
}

// Name returns the title of the vulnerability
func (detection *detection) Name() string {
	detection.lazyLoadVulnerabilityInfoForDetection()
	return detection.vulnerabilityInfo.Name()
}

// Description returns the consequence of the vulnerability
func (detection *detection) Description() string {
	detection.lazyLoadVulnerabilityInfoForDetection()
	return detection.vulnerabilityInfo.Description()
}

// CVSS2 returns the Common Vulnerability Scoring System version 2 score for the detection
func (detection *detection) CVSS2() float32 {
	detection.lazyLoadVulnerabilityInfoForDetection()
	return detection.vulnerabilityInfo.CVSS2()
}

// CVSS3 returns the Common Vulnerability Scoring System version 3 score for the detection
func (detection *detection) CVSS3() *float32 {
	detection.lazyLoadVulnerabilityInfoForDetection()
	return detection.vulnerabilityInfo.CVSS3()
}

// Solutions returns a channel containing
func (detection *detection) Solutions(ctx context.Context) (<-chan domain.Solution, error) {
	detection.lazyLoadVulnerabilityInfoForDetection()
	return detection.vulnerabilityInfo.Solutions(ctx)
}

func (detection *detection) References(ctx context.Context) (<-chan domain.VulnerabilityReference, error) {
	detection.lazyLoadVulnerabilityInfoForDetection()
	return detection.vulnerabilityInfo.References(ctx)
}

func (detection *detection) Severity() int {
	detection.lazyLoadVulnerabilityInfoForDetection()
	return detection.vulnerabilityInfo.Severity()
}

func (detection *detection) CVSS2Vector() string {
	detection.lazyLoadVulnerabilityInfoForDetection()
	return detection.vulnerabilityInfo.CVSS2Vector()
}

func (detection *detection) CVSS3Vector() string {
	detection.lazyLoadVulnerabilityInfoForDetection()
	return detection.vulnerabilityInfo.CVSS3Vector()
}

func (detection *detection) Software() string {
	detection.lazyLoadVulnerabilityInfoForDetection()
	return detection.vulnerabilityInfo.Software()
}

func (detection *detection) DetectionInformation() string {
	detection.lazyLoadVulnerabilityInfoForDetection()
	return detection.vulnerabilityInfo.DetectionInformation()
}
