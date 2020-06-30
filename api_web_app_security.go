package qualys

import (
	"encoding/xml"
	"fmt"
	"github.com/nortonlifelock/log"
	"net/http"
	"strings"
	"time"
)

const (
	postLaunchScan      = "/qps/rest/3.0/launch/was/wasscan/"
	postGetSiteFindings = "/qps/rest/3.0/search/was/finding"
	getScanStatus       = "/qps/rest/3.0/status/was/wasscan/<id>"
)

func (session *Session) CreateWebAppVulnerabilityScan(webAppID string, webAppOptionProfileID string, scannerType string, scannerName string) (scanID string, err error) {
	reqBody := &createWebAppScanRequest{}

	// TODO configurable scan name
	reqBody.Data.WasScan.Name = fmt.Sprintf("aegis_webapp_vulnerability_scan_%s", time.Now().Format(time.RFC3339))
	reqBody.Data.WasScan.Type = "VULNERABILITY"

	reqBody.Data.WasScan.Target.WebApp.ID = webAppID
	reqBody.Data.WasScan.Target.WebAppAuthRecord.IsDefault = "true"

	reqBody.Data.WasScan.Target.ScannerAppliance.Type = scannerType

	// TODO not sure if text is the proper place to put this
	reqBody.Data.WasScan.Target.ScannerAppliance.Text = scannerName
	reqBody.Data.WasScan.Profile.ID = webAppOptionProfileID

	var reqBodyByte []byte
	if reqBodyByte, err = xml.Marshal(reqBody); err == nil {
		reqBodyString := string(reqBodyByte)

		resp := webAppScanResponse{}

		if err = session.httpCall(http.MethodPost, session.webAppBaseURL+postLaunchScan, nil, &reqBodyString, resp); err == nil {

			if len(resp.Data.WasScan.ID) > 0 {
				scanID = resp.Data.WasScan.ID
			} else {
				session.lstream.Send(log.Errorf(err, "could not find scan ID from [%s]", postLaunchScan))
			}
		} else {
			session.lstream.Send(log.Errorf(err, "nil response while calling api [%s]", postLaunchScan))
		}
	} else {
		session.lstream.Send(log.Errorf(err, "error while marshalling scan body"))
	}

	return scanID, err
}

func (session *Session) GetScanStatus(scanID string) (status string, err error) {
	url := strings.Replace(session.webAppBaseURL+getScanStatus, "<id>", scanID, 1)

	resp := &webAppScanResponse{}

	if err = session.httpCall(http.MethodGet, url, nil, nil, resp); err == nil {
		if len(resp.Data.WasScan.Status) > 0 {
			status = resp.Data.WasScan.Status
		} else {
			session.lstream.Send(log.Errorf(err, "could not find status from [%s]", url))
		}
	} else {
		session.lstream.Send(log.Errorf(err, "err while calling api [%s]", url))
	}

	return status, err
}

func (session *Session) GetVulnerabilitiesForSite(siteID string) (findings []*WebAppFinding, err error) {
	reqBody := &WebAppFindingsRequest{}

	reqBody.Filters.Criteria.Field = "webApp.id"
	reqBody.Filters.Criteria.Operator = "EQUALS"
	reqBody.Filters.Criteria.Text = siteID

	resp := &webAppFindingsResponse{}

	var reqBodyByte []byte
	if reqBodyByte, err = xml.Marshal(reqBody); err == nil {
		reqBodyString := string(reqBodyByte)

		if err = session.httpCall(http.MethodPost, session.webAppBaseURL+postGetSiteFindings, nil, &reqBodyString, resp); err == nil {
			if len(resp.Data.Finding) > 0 {
				findings = resp.Data.Finding
			} else {
				session.lstream.Send(log.Errorf(err, "could not find status from [%s]", postGetSiteFindings))
			}
		} else {
			session.lstream.Send(log.Errorf(err, "err while calling api [%s]", postGetSiteFindings))
		}
	} else {
		session.lstream.Send(log.Errorf(err, "error while marshalling GetVulnerabilitiesForSite body"))
	}

	return findings, err
}

type createWebAppScanRequest struct {
	XMLName xml.Name `xml:"ServiceRequest"`
	Text    string   `xml:",chardata"`
	Data    struct {
		Text    string `xml:",chardata"`
		WasScan struct {
			Text   string `xml:",chardata"`
			Name   string `xml:"name"` // NAME
			Type   string `xml:"type"` // VULNERABILITY/DISCOVERY
			Target struct {
				Text   string `xml:",chardata"`
				WebApp struct {
					Text string `xml:",chardata"`
					ID   string `xml:"id"` // id here
				} `xml:"webApp"`
				WebAppAuthRecord struct {
					Text      string `xml:",chardata"`
					IsDefault string `xml:"isDefault"` // true
				} `xml:"webAppAuthRecord"`
				ScannerAppliance struct {
					Text string `xml:",chardata"`
					Type string `xml:"type"` // EXTERNAL/what else?
				} `xml:"scannerAppliance"`
			} `xml:"target"`
			Profile struct {
				Text string `xml:",chardata"`
				ID   string `xml:"id"` // option profile ID
			} `xml:"profile"`
		} `xml:"WasScan"`
	} `xml:"data"`
}

type webAppScanResponse struct {
	XMLName                   xml.Name `xml:"ServiceResponse"`
	Text                      string   `xml:",chardata"`
	Xsi                       string   `xml:"xsi,attr"`
	NoNamespaceSchemaLocation string   `xml:"noNamespaceSchemaLocation,attr"`
	ResponseCode              string   `xml:"responseCode"`
	Count                     string   `xml:"count"`
	Data                      struct {
		Text    string `xml:",chardata"`
		WasScan struct {
			Text   string `xml:",chardata"`
			ID     string `xml:"id"`
			Status string `xml:"status"`
		} `xml:"WasScan"`
	} `xml:"data"`
}

type WebAppFindingsRequest struct {
	XMLName     xml.Name `xml:"ServiceRequest"`
	Text        string   `xml:",chardata"`
	Preferences struct {
		Text    string `xml:",chardata"`
		Verbose string `xml:"verbose"`
	} `xml:"preferences"`
	Filters struct {
		Text     string `xml:",chardata"`
		Criteria struct {
			Text     string `xml:",chardata"`
			Field    string `xml:"field,attr"`
			Operator string `xml:"operator,attr"`
		} `xml:"Criteria"`
	} `xml:"filters"`
}

type webAppFindingsResponse struct {
	XMLName                   xml.Name `xml:"ServiceResponse"`
	Text                      string   `xml:",chardata"`
	Xsi                       string   `xml:"xsi,attr"`
	NoNamespaceSchemaLocation string   `xml:"noNamespaceSchemaLocation,attr"`
	ResponseCode              string   `xml:"responseCode"`
	Count                     string   `xml:"count"`
	HasMoreRecords            string   `xml:"hasMoreRecords"`
	Data                      struct {
		Text    string           `xml:",chardata"`
		Finding []*WebAppFinding `xml:"Finding"`
	} `xml:"data"`
}

type WebAppFinding struct {
	Text        string `xml:",chardata"`
	IDVal       string `xml:"id"`
	UniqueId    string `xml:"uniqueId"`
	Qid         string `xml:"qid"`
	Name        string `xml:"name"`
	Type        string `xml:"type"`
	FindingType string `xml:"findingType"`
	Cwe         struct {
		Text  string `xml:",chardata"`
		Count string `xml:"count"`
		List  struct {
			Text string `xml:",chardata"`
			Long string `xml:"long"`
		} `xml:"list"`
	} `xml:"cwe"`
	Owasp struct {
		Text  string `xml:",chardata"`
		Count string `xml:"count"`
		List  struct {
			Text  string `xml:",chardata"`
			OWASP struct {
				Text string `xml:",chardata"`
				Name string `xml:"name"`
				URL  string `xml:"url"`
				Code string `xml:"code"`
			} `xml:"OWASP"`
		} `xml:"list"`
	} `xml:"owasp"`
	Wasc struct {
		Text  string `xml:",chardata"`
		Count string `xml:"count"`
		List  struct {
			Text string `xml:",chardata"`
			WASC struct {
				Text string `xml:",chardata"`
				Name string `xml:"name"`
				URL  string `xml:"url"`
				Code string `xml:"code"`
			} `xml:"WASC"`
		} `xml:"list"`
	} `xml:"wasc"`
	ResultList struct {
		Text  string `xml:",chardata"`
		Count string `xml:"count"`
		List  struct {
			Text   string `xml:",chardata"`
			Result struct {
				Text           string `xml:",chardata"`
				Authentication string `xml:"authentication"`
				Ajax           string `xml:"ajax"`
				Payloads       struct {
					Text  string `xml:",chardata"`
					Count string `xml:"count"`
					List  struct {
						Text            string `xml:",chardata"`
						PayloadInstance struct {
							Text    string `xml:",chardata"`
							Payload string `xml:"payload"`
							Request struct {
								Text    string `xml:",chardata"`
								Method  string `xml:"method"`
								Link    string `xml:"link"`
								Headers string `xml:"headers"`
							} `xml:"request"`
							Response string `xml:"response"`
						} `xml:"PayloadInstance"`
					} `xml:"list"`
				} `xml:"payloads"`
			} `xml:"Result"`
		} `xml:"list"`
	} `xml:"resultList"`
	Severity          string `xml:"severity"`
	URL               string `xml:"url"`
	StatusVal         string `xml:"status"`
	FirstDetectedDate string `xml:"firstDetectedDate"`
	LastDetectedDate  string `xml:"lastDetectedDate"`
	LastTestedDate    string `xml:"lastTestedDate"`
	TimesDetected     string `xml:"timesDetected"`
	WebApp            struct {
		Text string `xml:",chardata"`
		ID   string `xml:"id"`
		Name string `xml:"name"`
		URL  string `xml:"url"`
		Tags struct {
			Text  string `xml:",chardata"`
			Count string `xml:"count"`
			List  struct {
				Text string `xml:",chardata"`
				Tag  []struct {
					Text string `xml:",chardata"`
					ID   string `xml:"id"`
					Name string `xml:"name"`
				} `xml:"Tag"`
			} `xml:"list"`
		} `xml:"tags"`
	} `xml:"webApp"`
	IsIgnored     string `xml:"isIgnored"`
	IgnoredReason string `xml:"ignoredReason"`
	IgnoredBy     struct {
		Text      string `xml:",chardata"`
		ID        string `xml:"id"`
		Username  string `xml:"username"`
		FirstName string `xml:"firstName"`
		LastName  string `xml:"lastName"`
	} `xml:"ignoredBy"`
	IgnoredDate    string `xml:"ignoredDate"`
	IgnoredComment string `xml:"ignoredComment"`
	Retest         string `xml:"retest"`
}
