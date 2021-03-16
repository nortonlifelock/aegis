package qualys

import (
	"encoding/xml"
	"fmt"
	"github.com/nortonlifelock/aegis/pkg/log"
	"net/http"
	"strings"
	"time"
)

const (
	postLaunchScan      = "/qps/rest/3.0/launch/was/wasscan/"
	postGetSiteFindings = "/qps/rest/3.0/search/was/finding"
	getScanStatus       = "/qps/rest/3.0/status/was/wasscan/<id>"
	getWebAppInfo       = "/qps/rest/3.0/get/was/webapp/<id>"
	postRetestFinding   = "/qps/rest/3.0/retest/was/finding/<uid>"
)

// findingUID in the form aaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa
func (session *Session) CreateRetestForWebAppVulnerabilityFinding(findingUID string) (count string, err error) {
	url := strings.Replace(session.webAppBaseURL+postRetestFinding, "<uid>", findingUID, 1)
	resp := &webAppRetestFindingResponse{}
	if err = session.httpCall(http.MethodPost, url, make(map[string]string), nil, resp); err == nil {
		count = resp.Count
	} else {
		session.lstream.Send(log.Errorf(err, "err while calling api [%s]", url))
	}

	return count, err
}

func (session *Session) CreateWebAppVulnerabilityScan(webAppID string, webAppOptionProfileID string, scannerType string, scannerName string) (scanID string, title string, err error) {
	if len(webAppID) > 0 && len(webAppOptionProfileID) > 0 && len(scannerType) > 0 { // scanner name can be empty for externel scans
		reqBody := &createWebAppScanRequest{}

		// TODO configurable scan name
		reqBody.Data.WasScan.Name = fmt.Sprintf("aegis_webapp_vulnerability_scan_%s", time.Now().Format(time.RFC3339))
		reqBody.Data.WasScan.Type = "VULNERABILITY"
		title = reqBody.Data.WasScan.Name

		reqBody.Data.WasScan.Target.WebApp.ID = webAppID
		reqBody.Data.WasScan.Target.WebAppAuthRecord.IsDefault = "true"

		reqBody.Data.WasScan.Target.ScannerAppliance.Type = scannerType

		// TODO not sure if text is the proper place to put this
		reqBody.Data.WasScan.Target.ScannerAppliance.FriendlyName = scannerName
		reqBody.Data.WasScan.Profile.ID = webAppOptionProfileID

		var reqBodyByte []byte
		if reqBodyByte, err = xml.Marshal(reqBody); err == nil {
			reqBodyString := string(reqBodyByte)

			resp := &webAppScanResponse{}

			if err = session.httpCall(http.MethodPost, session.webAppBaseURL+postLaunchScan, make(map[string]string), &reqBodyString, resp); err == nil {

				if len(resp.Data.WasScan.ID) > 0 {
					scanID = resp.Data.WasScan.ID
				} else {
					session.lstream.Send(log.Errorf(err, "could not find scan ID from [%s]", postLaunchScan))
					err = fmt.Errorf("could not find scan ID from [%s]", postLaunchScan)
				}
			} else {
				session.lstream.Send(log.Errorf(err, "nil response while calling api [%s]", postLaunchScan))
			}
		} else {
			session.lstream.Send(log.Errorf(err, "error while marshalling scan body"))
		}
	} else {
		err = fmt.Errorf("all input parameters not provided [%s|%s|%s]", webAppID, webAppOptionProfileID, scannerType)
	}

	return scanID, title, err
}

func (session *Session) GetWebAppScanStatus(scanID string) (status string, err error) {
	url := strings.Replace(session.webAppBaseURL+getScanStatus, "<id>", scanID, 1)

	resp := &webAppScanResponse{}

	if err = session.httpCall(http.MethodGet, url, make(map[string]string), nil, resp); err == nil {
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
	var hasMoreRecords = true
	var lastID = "0"

	// findings holds the parent and children findings together
	findings = make([]*WebAppFinding, 0)

	for hasMoreRecords {
		reqBody := &WebAppFindingsRequest{}

		webCriteria := SearchCriteria{
			Field:    "webApp.id",
			Operator: "EQUALS",
			Text:     siteID,
		}

		pagingCriteria := SearchCriteria{
			Field:    "id",
			Operator: "GREATER",
			Text:     lastID,
		}

		reqBody.Filters.Criteria = []SearchCriteria{webCriteria, pagingCriteria}

		reqBody.Preferences.Verbose = "true"

		resp := &webAppFindingsResponse{}

		var reqBodyByte []byte
		if reqBodyByte, err = xml.Marshal(reqBody); err == nil {
			reqBodyString := string(reqBodyByte)

			if err = session.httpCall(http.MethodPost, session.webAppBaseURL+postGetSiteFindings, make(map[string]string), &reqBodyString, resp); err == nil {
				if len(resp.Data.Finding) > 0 {
					findings = append(findings, resp.Data.Finding...)

					hasMoreRecords = strings.ToLower(resp.HasMoreRecords) == "true"
					lastID = resp.LastID
				} else {
					session.lstream.Send(log.Errorf(err, "could not find status from [%s]", postGetSiteFindings))
					break
				}
			} else {
				session.lstream.Send(log.Errorf(err, "err while calling api [%s]", postGetSiteFindings))
				break
			}
		} else {
			session.lstream.Send(log.Errorf(err, "error while marshalling GetVulnerabilitiesForSite body"))
			break
		}
	}

	return findings, err
}

func (session *Session) GetWebApplicationInfo(webAppID string) (defaultScannerName, defaultScannerType string, err error) {
	url := strings.Replace(session.webAppBaseURL+getWebAppInfo, "<id>", webAppID, 1)

	resp := &getWebAppResponse{}

	if err = session.httpCall(http.MethodGet, url, make(map[string]string), nil, resp); err == nil {
		defaultScannerType = resp.Data.WebApp.DefaultScanner.Type
		defaultScannerName = resp.Data.WebApp.DefaultScanner.FriendlyName
	} else {
		session.lstream.Send(log.Errorf(err, "err while calling api [%s]", url))
	}

	return defaultScannerName, defaultScannerType, err
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
					Text         string `xml:",chardata"`
					Type         string `xml:"type"` // EXTERNAL/what else?
					FriendlyName string `xml:"friendlyName"`
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
		Text     string           `xml:",chardata"`
		Criteria []SearchCriteria `xml:"Criteria"`
	} `xml:"filters"`
}

type SearchCriteria struct {
	Text     string `xml:",chardata"`
	Field    string `xml:"field,attr"`
	Operator string `xml:"operator,attr"`
}

type webAppFindingsResponse struct {
	XMLName                   xml.Name `xml:"ServiceResponse"`
	Text                      string   `xml:",chardata"`
	Xsi                       string   `xml:"xsi,attr"`
	NoNamespaceSchemaLocation string   `xml:"noNamespaceSchemaLocation,attr"`
	ResponseCode              string   `xml:"responseCode"`
	Count                     string   `xml:"count"`
	HasMoreRecords            string   `xml:"hasMoreRecords"`
	LastID                    string   `xml:"lastId"`
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
			Text string   `xml:",chardata"`
			Long []string `xml:"long"`
		} `xml:"list"`
	} `xml:"cwe"`
	Owasp struct {
		Text  string `xml:",chardata"`
		Count string `xml:"count"`
		List  struct {
			Text  string `xml:",chardata"`
			OWASP []struct {
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
				AccessPath     struct {
					XMLName xml.Name `xml:"accessPath"`
					Text    string   `xml:",chardata"`
					Count   string   `xml:"count"`
					List    struct {
						Text string   `xml:",chardata"`
						URL  []string `xml:"Url"`
					} `xml:"list"`
				} `xml:"accessPath"`
				Payloads struct {
					Text  string `xml:",chardata"`
					Count string `xml:"count"`
					List  struct {
						Text            string `xml:",chardata"`
						PayloadInstance []struct {
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
	Param             string `xml:"param"`
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

type getWebAppResponse struct {
	XMLName                   xml.Name `xml:"ServiceResponse"`
	Text                      string   `xml:",chardata"`
	Xsi                       string   `xml:"xsi,attr"`
	NoNamespaceSchemaLocation string   `xml:"noNamespaceSchemaLocation,attr"`
	ResponseCode              string   `xml:"responseCode"`
	Count                     string   `xml:"count"`
	Data                      struct {
		Text   string `xml:",chardata"`
		WebApp struct {
			Text  string `xml:",chardata"`
			ID    string `xml:"id"`
			Name  string `xml:"name"`
			URL   string `xml:"url"`
			Os    string `xml:"os"`
			Owner struct {
				Text      string `xml:",chardata"`
				ID        string `xml:"id"`
				Username  string `xml:"username"`
				FirstName string `xml:"firstName"`
				LastName  string `xml:"lastName"`
			} `xml:"owner"`
			Scope string `xml:"scope"`
			Uris  struct {
				Text  string `xml:",chardata"`
				Count string `xml:"count"`
				List  struct {
					Text string `xml:",chardata"`
					URL  string `xml:"Url"`
				} `xml:"list"`
			} `xml:"uris"`
			Attributes struct {
				Text  string `xml:",chardata"`
				Count string `xml:"count"`
			} `xml:"attributes"`
			DefaultProfile struct {
				Text string `xml:",chardata"`
				ID   string `xml:"id"`
				Name string `xml:"name"`
			} `xml:"defaultProfile"`
			DefaultScanner struct {
				Text         string `xml:",chardata"`
				Type         string `xml:"type"`
				FriendlyName string `xml:"friendlyName"`
			} `xml:"defaultScanner"`
			ScannerLocked string `xml:"scannerLocked"`
			UrlBlacklist  struct {
				Text  string `xml:",chardata"`
				Count string `xml:"count"`
			} `xml:"urlBlacklist"`
			UrlWhitelist struct {
				Text  string `xml:",chardata"`
				Count string `xml:"count"`
			} `xml:"urlWhitelist"`
			PostDataBlacklist struct {
				Text  string `xml:",chardata"`
				Count string `xml:"count"`
			} `xml:"postDataBlacklist"`
			LogoutRegexList struct {
				Text  string `xml:",chardata"`
				Count string `xml:"count"`
			} `xml:"logoutRegexList"`
			AuthRecords struct {
				Text  string `xml:",chardata"`
				Count string `xml:"count"`
			} `xml:"authRecords"`
			DnsOverrides struct {
				Text  string `xml:",chardata"`
				Count string `xml:"count"`
			} `xml:"dnsOverrides"`
			UseRobots           string `xml:"useRobots"`
			UseSitemap          string `xml:"useSitemap"`
			MalwareMonitoring   string `xml:"malwareMonitoring"`
			MalwareNotification string `xml:"malwareNotification"`
			Tags                struct {
				Text  string `xml:",chardata"`
				Count string `xml:"count"`
			} `xml:"tags"`
			Comments struct {
				Text  string `xml:",chardata"`
				Count string `xml:"count"`
			} `xml:"comments"`
			IsScheduled string `xml:"isScheduled"`
			LastScan    struct {
				Text string `xml:",chardata"`
				ID   string `xml:"id"`
				Name string `xml:"name"`
			} `xml:"lastScan"`
			CreatedBy struct {
				Text      string `xml:",chardata"`
				ID        string `xml:"id"`
				Username  string `xml:"username"`
				FirstName string `xml:"firstName"`
				LastName  string `xml:"lastName"`
			} `xml:"createdBy"`
			CreatedDate string `xml:"createdDate"`
			UpdatedBy   struct {
				Text      string `xml:",chardata"`
				ID        string `xml:"id"`
				Username  string `xml:"username"`
				FirstName string `xml:"firstName"`
				LastName  string `xml:"lastName"`
			} `xml:"updatedBy"`
			UpdatedDate     string `xml:"updatedDate"`
			Screenshot      string `xml:"screenshot"`
			Config          string `xml:"config"`
			CrawlingScripts struct {
				Text  string `xml:",chardata"`
				Count string `xml:"count"`
			} `xml:"crawlingScripts"`
		} `xml:"WebApp"`
	} `xml:"data"`
}

type webAppRetestFindingResponse struct {
	XMLName                   xml.Name `xml:"ServiceResponse"`
	Text                      string   `xml:",chardata"`
	Xsi                       string   `xml:"xsi,attr"`
	NoNamespaceSchemaLocation string   `xml:"noNamespaceSchemaLocation,attr"`
	ResponseCode              string   `xml:"responseCode"`
	Count                     string   `xml:"count"`
	Data                      struct {
		Text    string `xml:",chardata"`
		Finding struct {
			Text     string `xml:",chardata"`
			ID       string `xml:"id"`
			UniqueId string `xml:"uniqueId"`
		} `xml:"Finding"`
	} `xml:"data"`
}
