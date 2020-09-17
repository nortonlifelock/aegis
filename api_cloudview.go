package qualys

import (
	"encoding/json"
	"fmt"
	"github.com/nortonlifelock/domain"
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

type CloudConfigurationResp struct {
	Content          []Content `json:"content"`
	Pageable         Pageable  `json:"pageable"`
	Last             bool      `json:"last"`
	TotalPages       int       `json:"totalPages"`
	TotalElements    int       `json:"totalElements"`
	First            bool      `json:"first"`
	Sort             Sort      `json:"sort"`
	NumberOfElements int       `json:"numberOfElements"`
	Size             int       `json:"size"`
	Number           int       `json:"number"`
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
type Content struct {
	CloudType           string        `json:"cloudType"`
	UUID                string        `json:"uuid"`
	ScanUUID            string        `json:"scanUuid"`
	Name                string        `json:"name"`
	Description         interface{}   `json:"description"`
	IsGovCloud          bool          `json:"isGovCloud"`
	IsChinaRegion       bool          `json:"isChinaRegion"`
	Deleted             bool          `json:"deleted"`
	GlobalErrorMessage  string        `json:"globalErrorMessage"`
	ErrorDetails        interface{}   `json:"errorDetails"`
	Message             interface{}   `json:"message"`
	State               string        `json:"state"`
	LastSynch           string        `json:"lastSynch"`
	NextSynch           string        `json:"nextSynch"`
	RegionsNotProtected int           `json:"regionsNotProtected"`
	RegionsProtected    int           `json:"regionsProtected"`
	TotalAssets         int           `json:"totalAssets"`
	Modules             []string      `json:"modules"`
	Disabled            interface{}   `json:"disabled"`
	SynchType           interface{}   `json:"synchType"`
	SynchFrequency      interface{}   `json:"synchFrequency"`
	Regions             []Regions     `json:"regions"`
	TotalRegions        interface{}   `json:"totalRegions"`
	TotalErrors         interface{}   `json:"totalErrors"`
	AssetsCount         interface{}   `json:"assetsCount"`
	AssetsProtected     interface{}   `json:"assetsProtected"`
	AssetsNotProtected  interface{}   `json:"assetsNotProtected"`
	Groups              []interface{} `json:"groups"`
	RunFrequency        RunFrequency  `json:"runFrequency"`
	AwsAccountID        string        `json:"awsAccountId"`
	AccountAlias        interface{}   `json:"accountAlias"`
	BaseAccountID       interface{}   `json:"baseAccountId"`
	AwsExternalID       string        `json:"awsExternalId"`
	AwsArn              string        `json:"awsArn"`
	PortalUUID          string        `json:"portalUuid"`
	PortalConnector     bool          `json:"portalConnector"`
	CustomerBaseAccount bool          `json:"customerBaseAccount"`
	RegionUuids         interface{}   `json:"regionUuids"`
	ResponseCode        interface{}   `json:"responseCode"`
	TestResponse        interface{}   `json:"testResponse"`
	ResponseMessage     interface{}   `json:"responseMessage"`
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

func (session *Session) GetCloudAccountEvaluations(accountID string) (evaluations []AccountEvaluationContent, err error) {
	evaluations = make([]AccountEvaluationContent, 0)
	accountEvaluation := &AccountEvaluationResponse{}

	var accPage = 0
	var lastAccPage bool

	for !lastAccPage {
		var req *http.Request
		req, err = http.NewRequest(http.MethodGet, session.Config.Address()+fmt.Sprintf("/cloudview-api/rest/v1/aws/evaluations/%s?pageNo=%d&sortOrder=asc", accountID, accPage), nil)
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

func (session *Session) GetCloudEvaluationFindings(accountID string, content AccountEvaluationContent) (findings []domain.Finding, err error) {
	findings = make([]domain.Finding, 0)
	var last bool
	var page int

	for !last {
		evaluationResult := &EvaluationResult{}

		var req *http.Request
		req, err = http.NewRequest(http.MethodGet, session.Config.Address()+fmt.Sprintf("/cloudview-api/rest/v1/aws/evaluations/%s/resources/%s?pageNo=%d&sortOrder=asc", accountID, content.ControlID, page), nil)
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
			if finding.Result != fixedFinding {
				findings = append(findings, &cloudViewFinding{
					evaluationContent: finding,
					accountContent:    content,
				})
			}
		}
	}

	return findings, err
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
	return f.evaluationContent.AccountID
}

// ScanID corresponds to the assessment that found the finding
func (f *cloudViewFinding) ScanID() int {
	return 0
}

func (f *cloudViewFinding) Summary() string {
	return fmt.Sprintf("Aegis (%s) - %s", f.accountContent.ControlName, f.evaluationContent.ResourceID)
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
	return fmt.Sprintf("Region: %s\n\nEvidence\n%s\n\nResource Type: %s\n\nPolicy: %s\n\nControl ID: %s", f.evaluationContent.Region, evidences, f.evaluationContent.ResourceType, strings.Join(f.accountContent.PolicyNames, ", "), f.accountContent.ControlID)
}

// not relevant to cloud view
func (f *cloudViewFinding) BundleID() string {
	return ""
}
