package connector

import (
	"fmt"
	"github.com/nortonlifelock/aegis/pkg/domain"
	"github.com/nortonlifelock/aegis/pkg/qualys"
	"strconv"
	"strings"
	"time"
)

type webAppFindingWrapper struct {
	f *qualys.WebAppFinding

	session *QsSession
	vuln    *vulnerabilityInfo
}

// ID returns the Aegis DB value which is not available from Qualys API
func (f *webAppFindingWrapper) ID() string {
	return ""
}

// VulnerabilityID returns the QID
func (f *webAppFindingWrapper) VulnerabilityID() string {
	// TODO
	var webPath string
	if len(f.f.Param) > 0 && len(f.f.ResultList.List.Result.AccessPath.List.URL) > 0 {
		webPath = f.f.ResultList.List.Result.AccessPath.List.URL[0]
	}
	return fmt.Sprintf("%s;%s;%s",f.f.Qid, f.f.Param, webPath)
}

func (f *webAppFindingWrapper) Status() string {
	var status = f.f.StatusVal // NEW, ACTIVE, REOPENED, PROTECTED, and FIXED
	detectionType := f.f.Type  // VULNERABILITY, SENSITIVE_CONTENT, or INFORMATION_GATHERED

	const (
		info = "INFORMATION_GATHERED"
	)

	if detectionType == info {
		status = domain.Informational
	} else {
		switch strings.ToLower(status) {
		case "new":
			status = domain.Vulnerable
		case "active":
			status = domain.Vulnerable
		case "reopened":
			status = domain.Vulnerable
		case "fixed":
			status = domain.Fixed
		case "protected":
			status = domain.Vulnerable
			// TODO what does protected mean precisely?
		default:
			// do nothing
		}
	}

	return status
}

func (f *webAppFindingWrapper) ActiveKernel() *int {
	return nil
}

const (
	webAppFindingTimeFormat = "2006-01-02T15:04:05Z"
)

// Detected returns the date the finding was first found
func (f *webAppFindingWrapper) Detected() (*time.Time, error) {
	timeVal, err := time.Parse(webAppFindingTimeFormat, f.f.FirstDetectedDate)
	return &timeVal, err
}

func (f *webAppFindingWrapper) TimesSeen() int {
	timesSeen, _ := strconv.Atoi(f.f.TimesDetected)
	return timesSeen
}

func (f *webAppFindingWrapper) Proof() string {
	var accessPath string
	for index, url := range f.f.ResultList.List.Result.AccessPath.List.URL {
		if len(accessPath) > 0 {
			accessPath = fmt.Sprintf("%s\n%d) %s", accessPath, index, url)
		} else {
			accessPath = fmt.Sprintf("%d) %s", index, url)
		}
	}

	var payload string
	for index, payloadInstance := range f.f.ResultList.List.Result.Payloads.List.PayloadInstance {

		payloadInfo := fmt.Sprintf("Payload: %s\nMethod: %s\nLink: %s\nHeaders: %s\nResponse:\n%s",
			payloadInstance.Payload,
			payloadInstance.Request.Method,
			payloadInstance.Request.Link,
			payloadInstance.Request.Headers,
			payloadInstance.Response,
		)
		if len(payload) > 0 {
			payload = fmt.Sprintf("%s\n%d) %s", payload, index, payloadInfo)
		} else {
			payload = fmt.Sprintf("%d) %s", index, payloadInfo)
		}
	}

	var cwes, owasps string
	for index, cwe := range f.f.Cwe.List.Long {
		if len(cwes) > 0 {
			cwes = fmt.Sprintf("%s\n%d) CWE-%s", cwes, index, cwe)
		} else {
			cwes = fmt.Sprintf("%d) CWE-%s", index, cwe)
		}
	}
	for index, owasp := range f.f.Owasp.List.OWASP {
		if len(owasps) > 0 {
			owasps = fmt.Sprintf("%s\n%d) [%s|%s]", owasps, index, owasp.Name, owasp.URL)
		} else {
			owasps = fmt.Sprintf("%d) [%s|%s]", index, owasp.Name, owasp.URL)
		}
	}
	var references = fmt.Sprintf("CWE\n%s\n\nOWASP\n%s", cwes, owasps)

	return fmt.Sprintf("Access Paths\n%s\nPayloads\n%s\n%s", accessPath, payload, references)
}

func (f *webAppFindingWrapper) Port() int {
	return 0
}

func (f *webAppFindingWrapper) Protocol() string {
	return ""
}

func (f *webAppFindingWrapper) IgnoreID() (*string, error) {
	return nil, nil
}

func (f *webAppFindingWrapper) LastFound() *time.Time {
	timeVal, _ := time.Parse(webAppFindingTimeFormat, f.f.LastDetectedDate)
	return &timeVal
}

func (f *webAppFindingWrapper) LastUpdated() *time.Time {
	timeVal, _ := time.Parse(webAppFindingTimeFormat, f.f.LastTestedDate)
	return &timeVal
}

func (f *webAppFindingWrapper) Device() (domain.Device, error) {
	return &WebAppWrapper{
		sourceID:  f.f.WebApp.ID,
		findingID: f.f.UniqueId,
		name:      f.f.WebApp.Name,
		url:       f.f.WebApp.URL,
	}, nil
}

func (f *webAppFindingWrapper) Vulnerability() (domain.Vulnerability, error) {
	var err error
	if f.vuln == nil {
		qidInt, _ := strconv.Atoi(f.f.Qid)
		f.vuln = lazyLoadVulnerabilityInfo(qidInt, f.session)
	}
	return f.vuln, err
}
