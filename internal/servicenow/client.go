package servicenow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/nortonlifelock/funnel"
	"io"
	"net/http"
	"strconv"
)

// SvcNowClient is used to create a connection to ServiceNow and execute API calls
type SvcNowClient struct {
	// HTTP client used to communicate with the API.
	client funnel.Client

	// Base URL for API requests.
	baseURL string

	// Basic auth username
	username string
	// Basic auth password
	password string
}

// NewClient pulls together all the information required to create a SvcNowClient
func NewClient(httpClient funnel.Client, baseURL string, userName string, passWord string) (*SvcNowClient, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	c := &SvcNowClient{
		client:   httpClient,
		baseURL:  baseURL,
		username: userName,
		password: passWord,
	}

	return c, nil
}

func (client *SvcNowClient) newSvcNowRequest(table string, method string, id string, query string, fields []string, limit int, body interface{}) (req *http.Request, err error) {
	var sysIDWithSlash string
	//If sysId was provided pass it in a form of url/{sysId}
	if id != "" {
		sysIDWithSlash = fmt.Sprintf("/%s", id)

	}

	// specify the fields we want returned
	responseFields := "sysparm_fields=number,sys_id,priority,due_date,u_additional_information,state,substate,port,ip_address,dns,activity_due,qualys_severity,skills,additional_assignee_list,watch_list,qualys_ticket,group_list,user_input,description,close_notes,correlation_display,correlation_id,work_notes,closed_at,qualys_assignee_email,active"

	urlString := fmt.Sprintf("%s/%s%s?sysparm_display_value=true&%s&sysparm_limit=%s&sysparm_query=%s",
		client.baseURL,
		table,
		sysIDWithSlash,
		responseFields,
		strconv.Itoa(limit),
		query,
	)

	buf := &bytes.Buffer{}
	if body != nil {
		err = json.NewEncoder(buf).Encode(body)
	}

	if err == nil {
		if req, err = http.NewRequest(method, urlString, buf); err == nil {
			if body != nil {
				req.Header.Set("Content-Type", "application/json")
			}

			req.SetBasicAuth(client.username, client.password)

		}
	}
	return req, err
}

func (client *SvcNowClient) performFor(table string, method string, id string, query string, fields []string, limit int, body interface{}, out interface{}) (err error) {
	var stringID string
	var req *http.Request
	var res *http.Response
	u := fmt.Sprintf("%s/%s", client.baseURL, table)

	//If sysId was provided pass it in a form of url/{sysId}
	if id != "" {
		stringID = fmt.Sprintf("/%s", id)

	}

	resposefields := "sysparm_fields=number,sys_id,priority,due_date,u_additional_information,state,substate,port,ip_address,dns,activity_due,qualys_severity,skills,additional_assignee_list,watch_list,qualys_ticket,group_list,user_input,description,close_notes,correlation_display,correlation_id,work_notes,closed_at,qualys_assignee_email,active"
	//if method != http.MethodPost {
	//	//assignmentgroup to be added
	//	resposefields += ",assignment_group,assigned_to,requestor"
	//}
	//we are passing the size of the result and the query to the paramters
	urlString := fmt.Sprintf("%s%s?sysparm_display_value=true&%s&sysparm_limit=%s&sysparm_query=%s", u, stringID, resposefields, strconv.Itoa(limit), query)

	buf := &bytes.Buffer{}

	// if body is available pass it in for the PUT and Post methods
	if body != nil {
		err = json.NewEncoder(buf).Encode(body)
	}

	if err == nil {
		//Create Request
		if req, err = http.NewRequest(method, urlString, buf); err == nil {
			// set the contest type
			if body != nil {
				req.Header.Set("Content-Type", "application/json")
			}
			// set the authentication
			req.SetBasicAuth(client.username, client.password)

			// return an error id request failed
			if res, err = http.DefaultClient.Do(req); err == nil {

				defer res.Body.Close()

				buf.Reset()

				// Store JSON so we can do a preliminary error check
				var echeck Err
				//if method == http.MethodPut {
				//	var data []byte
				//	if data, err = ioutil.ReadAll(res.Body); err == nil {
				//
				//	}
				//	spew.Println(string(data))
				//}

				// If the request went through and you get unseccessful transaction grap that infromation from the response
				err = json.NewDecoder(io.TeeReader(res.Body, buf)).Decode(&echeck)
				if err == nil && echeck.ResError.Message != "" {
					err = fmt.Errorf("ServiceNow ERROR while caaling this api: %s,  Error: %s", urlString, echeck.Error())
				} else {
					err = json.NewDecoder(buf).Decode(out)
				}
			}
		}
	}
	return err

}
