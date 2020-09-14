package qualys

import (
	"fmt"
	"net/http"
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
	Content          []EvaluationContent `json:"content"`
	Pageable         Pageable            `json:"pageable"`
	Last             bool                `json:"last"`
	TotalPages       int                 `json:"totalPages"`
	TotalElements    int                 `json:"totalElements"`
	First            bool                `json:"first"`
	Sort             Sort                `json:"sort"`
	NumberOfElements int                 `json:"numberOfElements"`
	Size             int                 `json:"size"`
	Number           int                 `json:"number"`
}
type EvaluationContent struct {
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

func (session *Session) GetCloudViewFindings(accountID string) (err error) {
	// /cloudview/rest/1.0/aws/connectors/list?limit=50&offset=0&query=
	//

	/*
		/cloudview/rest/1.0/gcp/connectors/list
		/cloudview/rest/1.0/azure/connectors/list
		/cloudview/rest/1.0/aws/connectors/list
	*/

	gcpResp := &CloudConfigurationResp{}
	err = session.httpCall(
		http.MethodGet,
		session.Config.Address()+"/cloudview/rest/1.0/gcp/connectors/list",
		make(map[string]string),
		nil,
		gcpResp,
	)

	azureResp := &CloudConfigurationResp{}
	err = session.httpCall(
		http.MethodGet,
		session.Config.Address()+"/cloudview/rest/1.0/azure/connectors/list",
		make(map[string]string),
		nil,
		azureResp,
	)

	awsResp := &CloudConfigurationResp{}
	err = session.httpCall(
		http.MethodGet,
		session.Config.Address()+"/cloudview/rest/1.0/aws/connectors/list",
		make(map[string]string),
		nil,
		awsResp,
	)

	accountEvaluation := &AccountEvaluationResponse{}
	err = session.httpCall(
		http.MethodGet,
		session.Config.Address()+fmt.Sprintf("/cloudview-api/rest/v1/aws/evaluations/%s", accountID),
		make(map[string]string),
		nil,
		accountEvaluation,
	)

	for _, val := range accountEvaluation.Content {
		// TODO can page through results

		evaluationResult := &EvaluationResult{}
		err = session.httpCall(
			http.MethodGet,
			session.Config.Address()+fmt.Sprintf("/cloudview-api/rest/v1/aws/evaluations/%s/resources/%s", accountID, val.ControlID),
			make(map[string]string),
			nil,
			evaluationResult,
		)
	}

	/*
			GET /rest/v1/aws/evaluations/{accountId}/resources/{controlId}
		    GET /rest/v1/aws/evaluations/{accountID}
				from awsAccountId
	*/
	return err
}
