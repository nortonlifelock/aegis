package dome9

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// RunAssessmentOnBundle runs an assessment in Dome9 using the bundle and cloud account passed in
func (client *Client) RunAssessmentOnBundle(bundleID int, cloudAccountID string, cloudVendor string) (resp *AssessmentResult, err error) {
	req := &createAssessmentRequestBody{}
	req.CloudAccountID = cloudAccountID
	req.BundleID = bundleID
	req.CloudAccountType = cloudVendor

	var reqBody []byte
	if reqBody, err = json.Marshal(req); err == nil {

		var body []byte
		body, err = client.executeRequest(http.MethodPost, runAssessmentOnBundle, bytes.NewReader(reqBody))
		if err == nil {
			/*
				TODO do the following flags need to be taken into consideration?
				IsRelevant = {bool}
				IsValid = {bool}
				IsExcluded = {bool}
			*/
			resp = &AssessmentResult{}
			err = json.Unmarshal(body, resp)
		}
	}

	return resp, err
}
