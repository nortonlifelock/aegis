package qualys

import (
	"encoding/json"
	"fmt"
	"github.com/nortonlifelock/aegis/pkg/domain"
	"github.com/nortonlifelock/aegis/pkg/log"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type EvaluationResult struct {
	Content          []EvaluationResultContent `json:"content"`
	Pageable         Pageable                  `json:"pageable"`
	Last             bool                      `json:"last"`
	TotalPages       int                       `json:"totalPages"`
	TotalElements    int                       `json:"totalElements"`
	First            bool                      `json:"first"`
	Sort             Sort                      `json:"sort"`
	NumberOfElements int                       `json:"numberOfElements"`
	Size             int                       `json:"size"`
	Number           int                       `json:"number"`
}
type Evidences struct {
	SettingName string `json:"settingName"`
	ActualValue string `json:"actualValue"`
}
type EvaluationDates struct {
	FirstEvaluated string      `json:"firstEvaluated"`
	LastEvaluated  string      `json:"lastEvaluated"`
	DateReopen     interface{} `json:"dateReopen"`
	DateFixed      interface{} `json:"dateFixed"`
}
type EvaluationResultContent struct {
	ResourceID      string          `json:"resourceId"`
	Region          string          `json:"region"`
	AccountID       string          `json:"accountId"`
	EvaluatedOn     string          `json:"evaluatedOn"`
	Evidences       []Evidences     `json:"evidences"`
	ResourceType    string          `json:"resourceType"`
	ConnectorID     string          `json:"connectorId"`
	Result          string          `json:"result"`
	EvaluationDates EvaluationDates `json:"evaluationDates"`
}

type AccountEvaluationResponse struct {
	Content          []AccountEvaluationContent `json:"content"`
	Pageable         Pageable                   `json:"pageable"`
	Last             bool                       `json:"last"`
	TotalPages       int                        `json:"totalPages"`
	TotalElements    int                        `json:"totalElements"`
	First            bool                       `json:"first"`
	Sort             Sort                       `json:"sort"`
	NumberOfElements int                        `json:"numberOfElements"`
	Size             int                        `json:"size"`
	Number           int                        `json:"number"`
}
type AccountEvaluationContent struct {
	ControlName     string   `json:"controlName"`
	PolicyNames     []string `json:"policyNames"`
	Criticality     string   `json:"criticality"`
	Service         string   `json:"service"`
	Result          string   `json:"result"`
	ControlID       string   `json:"controlId"`
	PassedResources int      `json:"passedResources"`
	FailedResources int      `json:"failedResources"`
}

type Regions struct {
	UUID               interface{} `json:"uuid"`
	CodeName           string      `json:"codeName"`
	Name               interface{} `json:"name"`
	AssetsCount        int         `json:"assetsCount"`
	AssetsProtected    int         `json:"assetsProtected"`
	AssetsNotProtected int         `json:"assetsNotProtected"`
	Latitude           interface{} `json:"latitude"`
	Longitude          interface{} `json:"longitude"`
}
type RunFrequency struct {
	Hours   int `json:"hours"`
	Minutes int `json:"minutes"`
}

type Sort struct {
	Sorted   bool `json:"sorted"`
	Unsorted bool `json:"unsorted"`
}
type Pageable struct {
	Sort       Sort `json:"sort"`
	PageSize   int  `json:"pageSize"`
	PageNumber int  `json:"pageNumber"`
	Offset     int  `json:"offset"`
	Paged      bool `json:"paged"`
	Unpaged    bool `json:"unpaged"`
}

const (
	AWS_CLOUD_ACCOUNT    = "aws"
	AZURE_CLOUD_ACCOUNT  = "azure"
	GOOGLE_CLOUD_ACCOUNT = "gcp"
)

// getAccountType determines the account type. AccountID are determined through regexp filters and can be added to this func
func getAccountType(accountID string) string {
	matched := false
	// AWS
	if matched, _ = regexp.MatchString(`^\d{12}$`, accountID); matched{
		return AWS_CLOUD_ACCOUNT
	}
	// Azure
	if matched, _ = regexp.MatchString(`^[0-9A-Fa-f]{8}-([0-9A-Fa-f]{4}-){3}[0-9A-Fa-f]{12}$`, accountID); matched {
		return AZURE_CLOUD_ACCOUNT
	}
	//GCP
	return GOOGLE_CLOUD_ACCOUNT
}
// GetCloudAccountEvaluations was refactored with the ability to determine the account type via the getAccountType function. This removes unnecessary calls to the qualys api
func (session *Session) GetCloudAccountEvaluations(accountID string) (evaluations []AccountEvaluationContent, cloudAccountType string, err error) {
	accountType := getAccountType(accountID)
	var evals []AccountEvaluationContent

	if evals, err = session.GetCloudAccountEvaluationsWithCloudAccountType(accountID, accountType); err != nil {
		err = fmt.Errorf("error while gathering evaluations [%s|%s] - %s", accountID, accountType, err.Error())
	}

	return evals, accountType, err
}

func (session *Session) GetCloudAccountEvaluationsWithCloudAccountType(accountID string, cloudAccountType string) (evaluations []AccountEvaluationContent, err error) {
	evaluations = make([]AccountEvaluationContent, 0)
	accountEvaluation := &AccountEvaluationResponse{}

	var accPage = 0
	var lastAccPage bool

	for !lastAccPage {
		var req *http.Request
		// URL was split as percent signs in url affected Sprintf when running unit test
		req, err = http.NewRequest(http.MethodGet, session.Config.Address()+fmt.Sprintf("/cloudview-api/rest/v1/%s/evaluations/%s?pageNo=%d&sortOrder=asc&filter=evaluatedOn", cloudAccountType, accountID, accPage) +
			"%3A%5Bnow-24h%20..%20now-1s%5D%20and%20(policy.name%3ACIS%20Amazon%20Web%20Services%20Foundations%20Benchmark%20or%20policy.name%3AAegis%20AWS%20Benchmark%20or%20policy.name%3ACIS%20Microsoft%20Azure%20Foundations%20Benchmark%20or%20policy.name%3ACIS%20Google%20Cloud%20Platform%20Foundation%20Benchmark)",
			nil)

		if err != nil{
			err = fmt.Errorf("error creating url for request to Qualys - %s", err.Error())
			return nil, err
		}

		err = session.makeRequest(false, req, func(resp *http.Response) (err error) {
			var body []byte
			body, err = ioutil.ReadAll(resp.Body)
			if err == nil {
				err = json.Unmarshal(body, accountEvaluation)
			} else {
				err = fmt.Errorf("error while reading response body - %s", err.Error())
			}

			return  err
		})

		evaluations = append(evaluations, accountEvaluation.Content...)
		lastAccPage = accountEvaluation.Last
		accPage++
	}

	return evaluations, err
}

func (session *Session) GetCloudEvaluationFindings(accountID string, content AccountEvaluationContent, policyName string, cloudAccountType string) (findings []domain.Finding, err error) {
	findings = make([]domain.Finding, 0)
	var last bool
	var page int

	for !last {
		evaluationResult := &EvaluationResult{}
		var req *http.Request
		// URL was split as percent signs in url affected Sprintf when running unit test
		req, err = http.NewRequest(http.MethodGet, session.Config.Address()+fmt.Sprintf("/cloudview-api/rest/v1/%s/evaluations/%s/resources/%s?pageNo=%d", cloudAccountType, accountID, content.ControlID, page) +
			"&pageSize=1000&sortOrder=asc&filter=evaluatedOn%3A%5Bnow-24h%20..%20now-1s%5D", nil) // TODO qualys sorting does not seem to be working
		if err == nil {
			err = session.makeRequest(false, req, func(resp *http.Response) (err error) {
				var body []byte
				body, err = ioutil.ReadAll(resp.Body)
				if err == nil {
					err = json.Unmarshal(body, evaluationResult)
				} else {
					err = fmt.Errorf("error while reading response body - %s", err.Error())
				}

				return err
			})
		} else {
			err = fmt.Errorf("error while making request - %s", err.Error())
		}

		last = evaluationResult.Last
		page++

		const (
			fixedFinding = "PASS"
		)

		for _, finding := range evaluationResult.Content {
			if finding.Result != fixedFinding && !evidenceHasError(finding) && accountContentHasPolicy(policyName, content) && findingIsFresh(finding, session.lstream) {
				findings = append(findings, &cloudViewFinding{
					evaluationContent: finding,
					accountContent:    content,
					accountID:         accountID,
					policy:            policyName,
				})
			}
		}
	}

	return findings, err
}

// Sometimes the locations of the findings can be deleted (e.g. VM taken down) but the findings persist in CloudView
// Because evaluations should occur every couple hours, we consider a finding "stale" if it hasn't been found in 24 hrs
func findingIsFresh(finding EvaluationResultContent, lstream log.Logger) (fresh bool) {
	fresh = true
	const timeFormat = "2006-01-02T15:04:05-0700"
	val, err := time.Parse(timeFormat, finding.EvaluatedOn)
	if err == nil {
		if !val.IsZero() {
			if time.Since(val) > time.Hour*24 {
				fresh = false
			}
		} else {
			lstream.Send(log.Errorf(err, "empty evaluation date found for [%s|%s]", finding.ResourceID, finding.AccountID))
		}
	} else {
		lstream.Send(log.Errorf(err, "error parsing evaluation date found for [%s|%s]", finding.ResourceID, finding.AccountID))
	}

	return fresh
}

func accountContentHasPolicy(policyName string, content AccountEvaluationContent) (hasPolicy bool) {
	for _, policy := range content.PolicyNames {
		if strings.ToLower(policyName) == strings.ToLower(policy) {
			hasPolicy = true
			break
		}
	}

	return hasPolicy
}

func evidenceHasError(finding EvaluationResultContent) (hasError bool) {
	for _, evidence := range finding.Evidences {
		if strings.Contains(strings.ToLower(evidence.SettingName), "error") {
			hasError = true
			break
		}
	}

	return hasError
}

type cloudViewFinding struct {
	evaluationContent EvaluationResultContent
	accountContent    AccountEvaluationContent
	accountID         string
	policy            string
}

// ID corresponds to a vulnerability ID
func (f *cloudViewFinding) ID() string {
	return fmt.Sprintf("CV_%s", f.accountContent.ControlID)
}

// DeviceID corresponds to the entity violating the rule
func (f *cloudViewFinding) DeviceID() string {
	return f.evaluationContent.ResourceID
}

// AccountID corresponds to the cloud account that the entity lies within
func (f *cloudViewFinding) AccountID() string {
	if len(f.evaluationContent.AccountID) > 0 {
		return f.evaluationContent.AccountID
	} else {
		return f.accountID
	}
}

// ScanID corresponds to the assessment that found the finding
func (f *cloudViewFinding) ScanID() int {
	return 0
}

func (f *cloudViewFinding) Summary() string {
	return fmt.Sprintf("Aegis (%s-%s-%s)", f.accountID, f.DeviceID(), f.VulnerabilityTitle())
}
func (f *cloudViewFinding) VulnerabilityTitle() string {
	return strings.Replace(f.accountContent.ControlName, "\n", "", -1)
}
func (f *cloudViewFinding) Priority() string {
	return strings.Title(strings.ToLower(f.accountContent.Criticality))
}

func (f *cloudViewFinding) LastFound() time.Time {
	val, _ := time.Parse("2006-01-02T15:04:05+0000", f.evaluationContent.EvaluatedOn)
	return val
}

// String extracts relevant information from the finding
func (f *cloudViewFinding) String() string {
	var evidences string
	for index, evidence := range f.evaluationContent.Evidences {
		if index == 0 {
			evidences = fmt.Sprintf("%s: %s", evidence.SettingName, evidence.ActualValue)
		} else {
			evidences = fmt.Sprintf("%s\n%s: %s", evidences, evidence.SettingName, evidence.ActualValue)
		}
	}
	return fmt.Sprintf("Resource: %s\n\nRegion: %s\n\nEvidence\n%s\n\nResource Type: %s\n\nPolicy: %s\n\nControl ID: %s", f.evaluationContent.ResourceID, f.evaluationContent.Region, evidences, f.evaluationContent.ResourceType, strings.Join(f.accountContent.PolicyNames, ", "), f.accountContent.ControlID)
}

// not relevant to cloud view
func (f *cloudViewFinding) BundleID() string {
	return f.policy
}
