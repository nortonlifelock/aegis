package qualys

import (
	"encoding/json"
	"fmt"
	"github.com/nortonlifelock/aegis/pkg/domain"
	"io/ioutil"
	"net/http"
	"strings"
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

func (session *Session) GetCloudAccountEvaluations(accountID string) (evaluations []AccountEvaluationContent, cloudAccountType string, err error) {
	var possibleAccountTypes = []string{AWS_CLOUD_ACCOUNT, AZURE_CLOUD_ACCOUNT, GOOGLE_CLOUD_ACCOUNT}

	// from looking at the API documentation, I don't see a way to find the cloud account type by using the cloud account ID alone
	// so we just check all three and use one if it's present
	for _, possibleCloudAccountType := range possibleAccountTypes {
		var possibleEvals []AccountEvaluationContent

		if possibleEvals, err = session.GetCloudAccountEvaluationsWithCloudAccountType(accountID, possibleCloudAccountType); err == nil {
			for _, eval := range possibleEvals {
				if eval.FailedResources > 0 || eval.PassedResources > 0 {
					evaluations = possibleEvals
					cloudAccountType = possibleCloudAccountType
					break
				}
			}
		} else {
			err = fmt.Errorf("error while determining cloud account type for evaluation gathering [%s|%s]", accountID, possibleCloudAccountType)
			break
		}
	}

	if err == nil && len(cloudAccountType) == 0 {
		err = fmt.Errorf("could not determine cloud account type for accountID [%s]", accountID)
	}

	return evaluations, cloudAccountType, err
}

func (session *Session) GetCloudAccountEvaluationsWithCloudAccountType(accountID string, cloudAccountType string) (evaluations []AccountEvaluationContent, err error) {
	evaluations = make([]AccountEvaluationContent, 0)
	accountEvaluation := &AccountEvaluationResponse{}

	var accPage = 0
	var lastAccPage bool

	for !lastAccPage {
		var req *http.Request
		req, err = http.NewRequest(http.MethodGet, session.Config.Address()+fmt.Sprintf("/cloudview-api/rest/v1/%s/evaluations/%s?pageNo=%d&sortOrder=asc", cloudAccountType, accountID, accPage), nil)
		if err == nil {
			err = session.makeRequest(req, func(resp *http.Response) (err error) {
				var body []byte
				body, err = ioutil.ReadAll(resp.Body)
				if err == nil {
					err = json.Unmarshal(body, accountEvaluation)
				} else {
					err = fmt.Errorf("error while reading response body - %s", err.Error())
				}

				return err
			})
		} else {
			err = fmt.Errorf("error while making request - %s", err.Error())
		}

		if err == nil {
			evaluations = append(evaluations, accountEvaluation.Content...)
			lastAccPage = accountEvaluation.Last
			accPage++
		} else {
			break
		}
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
		req, err = http.NewRequest(http.MethodGet, session.Config.Address()+fmt.Sprintf("/cloudview-api/rest/v1/%s/evaluations/%s/resources/%s?pageNo=%d&sortOrder=asc", cloudAccountType, accountID, content.ControlID, page), nil)
		if err == nil {
			err = session.makeRequest(req, func(resp *http.Response) (err error) {
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
			if finding.Result != fixedFinding && !evidenceHasError(finding) && accountContentHasPolicy(policyName, content) {
				findings = append(findings, &cloudViewFinding{
					evaluationContent: finding,
					accountContent:    content,
					accountID:         accountID,
				})
			}
		}
	}

	return findings, err
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
	return fmt.Sprintf("Aegis (%s)", strings.Replace(f.accountContent.ControlName, "\n", "", -1))
}
func (f *cloudViewFinding) VulnerabilityTitle() string {
	return f.accountContent.ControlName
}
func (f *cloudViewFinding) Priority() string {
	return strings.Title(strings.ToLower(f.accountContent.Criticality))
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
	return ""
}
