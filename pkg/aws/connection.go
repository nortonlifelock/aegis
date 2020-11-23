package awsclient

import (
	"encoding/json"
	"fmt"

	"github.com/nortonlifelock/domain"

	// Golang sometimes shows errors on this package because the files are too large for it to work, but it compiles and works
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// CloudConnection is a struct that implements the interface for gathering tags
type CloudConnection struct {
	config     *aws.Config
	ec2Service *ec2.EC2
	regions    []string
}

const (
	// State can be found using the tag API, but is not technically a tag, so it is included in the code
	State = "State"
	// InstanceID can be found using the tag API, but is not technically a tag, so it is included in the code
	InstanceID = "InstanceId"
)

// CloudSyncPayload parses the job history Payload for the cloud sync job
type CloudSyncPayload struct {
	Region string `json:"region"`
}

// buildPayload loads the Payload from the job history into the Payload object
func buildPayload(pjson string) (payload *CloudSyncPayload, err error) {
	payload = &CloudSyncPayload{}

	if len(pjson) > 0 {
		err = json.Unmarshal([]byte(pjson), payload)
		if err == nil {

			if len(payload.Region) <= 0 {
				err = fmt.Errorf("the Cloud Source Payload did not contain the aws region")
			}
		}
	} else {
		err = fmt.Errorf("no Payload provided to job")
	}

	return payload, err
}

// CreateConnection returns a struct that is used to interact with the AWS API
func CreateConnection(authInfoJSON string, payloadJSON string) (connection *CloudConnection, err error) {
	// accessKeyID corresponds to the username
	// secretAccessKey corresponds to the password
	var accessKeyID, secretAccessKey string
	var authInfo domain.BasicAuth
	if err = json.Unmarshal([]byte(authInfoJSON), &authInfo); err == nil {
		accessKeyID = authInfo.Username
		secretAccessKey = authInfo.Password

		if len(accessKeyID) > 0 && len(secretAccessKey) > 0 {
			var payload *CloudSyncPayload
			if payload, err = buildPayload(payloadJSON); err == nil {

				if connection, err = createCloudConnectionForRegion(credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""), payload.Region); err == nil {
					//connection.regions, err = connection.GetAllRegions()
					connection.regions = []string{payload.Region}
				}
			} else {
				err = fmt.Errorf("error while building the AWS payload - %s", err.Error())
			}
		} else {
			err = fmt.Errorf("both an access key id and a secret access key must be provided")
		}
	} else {
		err = fmt.Errorf("error while parsing authentication information - %s", err.Error())
	}

	return connection, err
}

func createCloudConnectionForRegion(creds *credentials.Credentials, region string) (connection *CloudConnection, err error) {
	config := &aws.Config{
		Region:      aws.String(region),
		Credentials: creds,
	}

	var sess *session.Session
	if sess, err = session.NewSession(config); err == nil {

		connection = &CloudConnection{
			config:     config,
			ec2Service: ec2.New(sess),
		}
	}

	return connection, err
}
