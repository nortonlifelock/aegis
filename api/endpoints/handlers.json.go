package endpoints

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/nortonlifelock/aegis/internal/database/dal"
	"github.com/nortonlifelock/domain"
	"github.com/pkg/errors"
)

type (
	logRequest struct {
		MethodOfDiscovery string `json:"methodOfDiscovery"`
		JobType           int    `json:"jobType"`
		LogType           int    `json:"logType"`
		JobHistoryID      string `json:"jobHistoryId"`
		FromDate          string `json:"fromDate"`
		ToDate            string `json:"toDate"`
	}

	loginReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
		OrgID    string `json:"org"`
	}

	apiRequest struct {
		Source         *Source           `json:"source,omitempty"`
		Organization   *Organization     `json:"organization,omitempty"`
		Config         *JobConfig        `json:"job_config,omitempty"`
		History        *JobHistory       `json:"job_history,omitempty"`
		Histories      *Histories        `json:"histories,omitempty"`
		Users          *User             `json:"user,omitempty"`
		Tag            *TagMap           `json:"tag,omitempty"`
		BulkUpdateJob  *BulkUpdateJob    `json:"update,omitempty"`
		BulkUpdateFile *BulkUpdateFile   `json:"upload,omitempty"`
		Permission     []string          `json:"permissions,omitempty"`
		JobConfigs     *JobConfigRequest `json:"jobConfigs,omitempty"`
		Exceptions     *Exception        `json:"exceptions,omitempty"`
	}

	// BulkUpdateJob is a member of apiRequest and must be exported in order to be marshaled
	BulkUpdateJob struct {
		Filenames  []string `json:"files"`
		ServiceURL string   `json:"serviceURL"`
	}

	// BulkUpdateFile is a member of apiRequest and must be exported in order to be marshaled
	BulkUpdateFile struct {
		Filename string `json:"name"`
		Contents string `json:"contents"`
	}

	// JobHistory is a member of apiRequest and must be exported in order to be marshaled
	JobHistory struct {
		JobHistID string `json:"history_id"`
		JobID     int    `json:"job_id"`
		ConfigID  string `json:"config_id"`
		StatusID  int    `json:"status"`
		Payload   string `json:"payload"`
	}

	// JobConfig is a member of apiRequest and must be exported in order to be marshaled
	JobConfig struct {
		JobID                 int    `json:"job_id"`
		ConfigID              string `json:"config_id"`
		PriorityOverride      int    `json:"priority_override"`
		DataInSourceConfigID  string `json:"data_in_source_config_id"`
		DataOutSourceConfigID string `json:"data_out_source_config_id"`
		Continuous            bool   `json:"continuous"`
		WaitInSeconds         int    `json:"wait_in_seconds"`
		MaxInstances          int    `json:"max_instances"`
		AutoStart             bool   `json:"autostart"`
		Active                bool   `json:"active"`
		Payload               string `json:"payload"`
		Code                  int    `json:"code,omitempty"`
		Name                  string `json:"name,omitempty"`
		CreatedDate           string `json:"created_date"`
		CreatedBy             string `json:"created_by"`
		UpdatedDate           string `json:"updated_date,omitempty"`
		UpdatedBy             string `json:"updated_by,omitempty"`
		LastJobStart          string `json:"last_job_start,omitempty"`
		JobName               string `json:"job_name,omitempty"`
		SourceInName          string `json:"srcin_name,omitempty"`
	}

	// JobConfigRequest is a member of apiRequest and must be exported in order to be marshaled
	JobConfigRequest struct {
		JobConfig
		Offset      int    `json:"offset,omitempty"`
		Limit       int    `json:"limit,omitempty"`
		SortedField string `json:"sorted_field,omitempty"`
		SortOrder   string `json:"sort_order,omitempty"`
	}

	// Src returns the ... TODO:
	Src struct {
		Code string `json:"code"`
		Name string `json:"name,omitempty"`
	}

	// ExceptionType is the type of the exception in the database
	ExceptionType struct {
		Code int    `json:"code"`
		Name string `json:"name,omitempty"`
	}

	// Exception is the struct that maps to the exceptions in the database
	Exception struct {
		SourceID        string `json:"source_id"`
		OrganizationID  string `json:"org_id"`
		TypeID          int    `json:"type_id"`
		VulnerabilityID string `json:"vuln_id"`
		DeviceID        string `json:"device_id"`
		DueDate         string `json:"due_date,omitempty"`
		Approval        string `json:"approval,omitempty"`
		Active          bool   `json:"active"`
		Port            string `json:"port,omitempty"`
		DBCreatedDate   string `json:"created_date,omitempty"`
		DBUpdatedDate   string `json:"updated_date,omitempty"`
		CreatedBy       string `json:"created_by,omitempty"`
		UpdatedBy       string `json:"updated_by,omitempty"`
		Offset          int    `json:"offset,omitempty"`
		Limit           int    `json:"limit,omitempty"`
		SortedField     string `json:"sorted_field,omitempty"`
		SortOrder       string `json:"sort_order,omitempty"`
	}

	// SourceConfig maps to the source config db tables
	SourceConfig struct {
		Code string `json:"code,omitempty"`
		Name string `json:"name,omitempty"`
	}

	// Source is a member of apiRequest and must be exported in order to be marshaled
	Source struct {
		ID       string `json:"id,omitempty"`
		Source   string `json:"source,omitempty"`
		Address  string `json:"address,omitempty"`
		Port     int    `json:"port,omitempty"`
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
		Payload  string `json:"payload,omitempty"`

		// fields for Oauth (optional)
		PrivateKey  string `json:"private_key,omitempty"`
		ConsumerKey string `json:"consumer_key,omitempty"`
		Token       string `json:"token,omitempty"`
	}

	sourceDtoContainer struct {
		Sources           []*Source         `json:"sources"`
		UniqueSourceNames []string          `json:"unique_names,omitempty"`
		SourceToSkeleton  map[string]string `json:"source_to_payload"`
	}

	// Job is a member of apiRequest and must be exported in order to be marshaled
	Job struct {
		Code int    `json:"code,omitempty"`
		Name string `json:"name,omitempty"`
	}

	// LogType is a member of apiRequest and must be exported in order to be marshaled
	LogType struct {
		ID   int    `json:"id"`
		Type string `json:"type"`
		Name string `json:"name"`
	}

	// Log is a member of apiRequest and must be exported in order to be marshaled
	Log struct {
		ID           int    `json:"id"`
		TypeID       int    `json:"typeId"`
		Log          string `json:"log"`
		Error        string `json:"error"`
		JobHistoryID string `json:"jobHistoryId"`
		Date         string `json:"date"`
	}

	// VulnerabilityMatch is a member of apiRequest and must be exported in order to be marshaled
	VulnerabilityMatch struct {
		NexposeTitle    string `json:"nexpose_title"`
		NexposeID       string `json:"nexpose_id"`
		QualysTitle     string `json:"qualys_title"`
		QualysID        string `json:"qualys_id"`
		MatchConfidence int    `json:"match_confidence"`
		MatchReason     string `json:"match_reason"`
	}

	// Vulnerability is a member of apiRequest and must be exported in order to be marshaled
	Vulnerability struct {
		SourceID string  `json:"source_id"`
		Title    string  `json:"title"`
		CVSS     float32 `json:"cvss"`
	}

	// TagMap is a member of apiRequest and must be exported in order to be marshaled
	TagMap struct {
		CloudSource  string `json:"cloud_source"`
		CloudTag     string `json:"cloud_tag"`
		TicketSource string `json:"ticket_source"`
		TicketTag    string `json:"ticket_tag"`
		Option       string `json:"option"`
		Overwrite    bool   `json:"overwrite,omitempty"`
	}

	// Organization is a member of apiRequest and must be exported in order to be marshaled
	Organization struct {
		ID             string  `json:"org_id,omitempty"`
		Code           string  `json:"code,omitempty"`
		Description    string  `json:"description,omitempty"`
		TimeZoneOffset float32 `json:"timezone_offset,omitempty"`
	}

	// ScanSummary is a member of apiRequest and must be exported in order to be marshaled
	ScanSummary struct {
		ScanID    string `json:"scan_id"`
		Status    string `json:"status"`
		Tickets   string `json:"tickets"`
		Devices   string `json:"devices"`
		StartTime string `json:"start"`
		Duration  string `json:"duration"`
		GroupName string `json:"group_name"`
		GroupID   int    `json:"group_id"`
		Source    string `json:"source"`
	}

	// User is a member of apiRequest and must be exported in order to be marshaled
	User struct {
		ID         string `json:"id"`
		Username   string `json:"username"`
		FirstName  string `json:"firstname"`
		LastName   string `json:"lastname"`
		Email      string `json:"email"`
		IsDisabled bool   `json:"isdisabled"`
	}

	// Histories is a member of apiRequest and must be exported in order to be marshaled
	Histories struct {
		Offset    int    `json:"offset,omitempty"`
		Limit     int    `json:"limit,omitempty"`
		JobID     int    `json:"job_id"`
		JobHistID int    `json:"history_id"`
		ConfigID  string `json:"config_id"`
		StatusID  int    `json:"status"`
		Payload   string `json:"payload"`
	}

	// Permission is a member of apiRequest and must be exported in order to be marshaled
	Permission struct {
		Name  string `json:"name"`
		Value bool   `json:"value"`
	}

	// GeneralResp acts as a container for the response to a user's API apiRequest
	GeneralResp struct {
		Response interface{} `json:"response,omitempty"`
		Message  string      `json:"message,omitempty"`
		Success  bool        `json:"success,omitempty"`
		// TODO return permissions by marshalling dal object
		Permissions  *PublicPermission `json:"permissions,omitempty"`
		TotalRecords int               `json:"totalrecords,omitempty"`
	}
)

func newResponse(obj interface{}, totalRecords int) (resp *GeneralResp) {
	resp = &GeneralResp{
		Response:     obj,
		TotalRecords: totalRecords,
	}

	return resp
}

// PublicPermission is a member of GeneralResp and must be exported in order to be marshaled
type PublicPermission struct {
	CanDeleteJob        bool `json:"candeletejob,omitempty"`
	CanDeleteConfig     bool `json:"candeleteconfig,omitempty"`
	CanUpdateSource     bool `json:"canupdatesource,omitempty"`
	CanReadJobHistories bool `json:"canreadjobhistories,omitempty"`
	UserID              int  `json:"userid,omitempty"`
	CanCreateJob        bool `json:"cancreatejob,omitempty"`
	CanCreateSource     bool `json:"cancreatesource,omitempty"`
	CanDeleteOrg        bool `json:"candeleteorg,omitempty"`
	CanCreateConfig     bool `json:"cancreateconfig,omitempty"`
	CanUpdateConfig     bool `json:"canupdateconfig,omitempty"`
	CanUpdateJob        bool `json:"canupdatejob,omitempty"`
	CanDeleteSource     bool `json:"candeletesource,omitempty"`
	CanCreateOrg        bool `json:"cancreateorg,omitempty"`
	CanUpdateOrg        bool `json:"canupdateorg ,omitempty"`
}

type tokenPermission struct {
	*dal.Permission
}

func dalPermissionToPermissionArray(permission domain.Permission) (permissionArray []*Permission, err error) {
	if permission != nil {
		permissionArray = make([]*Permission, 0)
		permissionArray = append(permissionArray, &Permission{Name: "admin", Value: permission.Admin()})
		permissionArray = append(permissionArray, &Permission{Name: "manager", Value: permission.Manager()})
		permissionArray = append(permissionArray, &Permission{Name: "reader", Value: permission.Reader()})
		permissionArray = append(permissionArray, &Permission{Name: "reporter", Value: permission.Reporter()})
	} else {
		err = fmt.Errorf("permissions not found")
	}

	return permissionArray, err
}

// UnmarshalJSON implements the Unmarshaler interface on tokenPermission
func (perm *tokenPermission) UnmarshalJSON(b []byte) (err error) {
	perm.Permission = &dal.Permission{}
	permissionReflection := reflect.TypeOf(perm.Permission)

	var stringVal = string(b)
	//remove the surrounding brackets of the json string
	stringVal = strings.Trim(stringVal, "{")
	stringVal = strings.Trim(stringVal, "}")

	//take each json field with it's attached value and store it in a slice
	var fields = strings.Split(stringVal, ",")

	for _, field := range fields {

		if len(field) > 0 {
			//separate the field from it's value
			var fieldValuePair = strings.Split(field, ":")
			if len(fieldValuePair) == 2 {
				//the field name has quotes around it - remove them
				fieldValuePair[0] = strings.Replace(fieldValuePair[0], "\"", "", 2)

				//iterate through each method
				for methodIndex := 0; methodIndex < permissionReflection.NumMethod(); methodIndex++ {
					method := permissionReflection.Method(methodIndex)
					var methodName = strings.ToLower(method.Name)
					//find the set method that corresponds to the json field
					if strings.Index(methodName, "set") == 0 && strings.Index(methodName, fieldValuePair[0]) > 0 {
						var arguments []reflect.Value
						//pull the argument from the json pair
						if strings.ToLower(fieldValuePair[1]) == "true" {
							arguments = []reflect.Value{reflect.ValueOf(perm.Permission), reflect.ValueOf(true)}
						} else if strings.ToLower(fieldValuePair[1]) == "false" {
							arguments = []reflect.Value{reflect.ValueOf(perm.Permission), reflect.ValueOf(false)}
						} else { //int
							var intVal int
							intVal, err = strconv.Atoi(fieldValuePair[1])
							if err == nil {
								arguments = []reflect.Value{reflect.ValueOf(perm.Permission), reflect.ValueOf(int(intVal))}
							} else {
								err = errors.New("unable to unmarshal json string")
								break
							}
						}

						if err == nil {
							//execute the set method with the parsed argument
							method.Func.Call(arguments)
						}

						break
					}
				}

			} else {
				err = errors.New("unable to unmarshal json string")
			}
		} else {
			// TODO should this ever hit/should this be an error?
		}
	}

	return err
}

// MarshalJSON implements the Marshaler interface on tokenPermission
func (perm *tokenPermission) MarshalJSON() ([]byte, error) {
	var val []byte
	var err error
	var finalByte []byte

	val, err = json.Marshal(perm.Permission)

	if err == nil {
		var stringVal = string(val)
		stringVal = strings.Replace(stringVal, "{", "", -1)
		stringVal = strings.Replace(stringVal, "}", "", -1)
		elements := strings.Split(stringVal, ",")

		var permissions = make([]string, 0)

		for _, element := range elements {
			if strings.Index(element, "can") >= 0 {
				permissions = append(permissions, element)
			}
		}

		var jsonVal = "{" + strings.Join(permissions, ",") + "}"

		finalByte = []byte(jsonVal)
	}

	return finalByte, err
}
