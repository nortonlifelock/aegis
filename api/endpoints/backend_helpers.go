package endpoints

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/nortonlifelock/aegis/pkg/crypto"
	"github.com/nortonlifelock/aegis/pkg/domain"
	"github.com/pkg/errors"
)

// These consts tie an endpoint to it's method. It's used to discern which apiRequest landed at which endpoint in the DBLog table
const (
	unauthenticated = "unauthenticated"

	getAllJobsEndpoint     = "getAllJobs"
	getAllLogTypesEndpoint = "GetAllJobTypes"
	getLogsEndpoint        = "GetAllLogs"

	deleteJConfigEndpoint       = "DeleteConfig"
	createJConfigEndpoint       = "CreateConfig"
	updateJConfigEndpoint       = "UpdateConfig"
	getAllJobConfigsEndpoint    = "GetAllConfigJobs"
	getSourceInByJobIDEndpoint  = "GetSourceInForJob"
	getSourceOutByJobIDEndpoint = "GetSourceOutByJobID"
	getJCAuditHistoryEndpoint   = "GetJobConfigAuditHistory"

	deleteHistoryEndpoint  = "DeleteHistory"
	updateHistoryEndpoint  = "UpdateHistory"
	createHistoryEndpoint  = "CreateHistory"
	getHistoriesEndpoint   = "GetHistories"
	getHistoryByIDEndpoint = "getJobHistoryByID"

	updateOrgEndpoint     = "updateOrg"
	deleteOrgEndpoint     = "deleteOrg"
	createOrgEndpoint     = "createOrg"
	getAllOrgsEndpoint    = "GetAllOrgs"
	getOrgForUserEndpoint = "getOrgForUser"
	getMyOrgEndpoint      = "getMyOrg"

	updateSourceEndpoint = "UpdateSource"
	deleteSourceEndpoint = "DeleteSource"
	createSourceEndpoint = "CreateSource"
	getSourcesEndpoint   = "GetSources"

	updateUserEndpoint   = "updateUser"
	deleteUserEndpoint   = "deleteUser"
	createUserEndpoint   = "createUser"
	getAllUsersEndpoint  = "getAllUsers"
	getUserByIDEndpoint  = "getUserByID"
	getUsersNameEndpoint = "GetUsersName"

	getUserPermEndpoint    = "GetUserPermissionsByUserId"
	updateUserPermEndpoint = "UpdateUserPermissionsByUserId"
	getPermListEndpoint    = "getPermissionList"
	createUserPermEndpoint = "CreateUserPermission"

	getScansEndpoint = "GetScans"
	getAgsEndpoint   = "GetGroups"

	createScanWebsocketEndpoint = "ScanWebsocket"

	loginEndpoint  = "login"
	logoutEndpoint = "logout"

	bulkUpdateEndpoint         = "BulkUpdateJob"
	bulkUpdateUploadEndpoint   = "bulkUpdateUpload"
	createWSConnectionEndpoint = "connect"
	internalWSEndpoint         = "InternalWebSocket"

	getVulnBySourceEndpoint = "GetVulnBySource"
	getMatchedVulnsEndpoint = "getMatchedVulns"

	getTicketStatusCountEndpoint = "GetTicketStatusCount"

	burnDownWebSocketEndpoint   = "BurnDownWebSocket"
	getFieldsForProjectEndpoint = "GetFieldsForProject"
	getJIRAURLSEndpoint         = "getJiraUrls"
	getStatusMapsEndpoint       = "getStatusMaps"
	getFieldMapsEndpoint        = "getFieldMaps"
	attachCERFToTicketEndpoint  = "attachCERFToTicket"

	getAzureTagsEndpoint  = "GetAzureTags"
	getAWSTagsEndpoint    = "GetAWSTags"
	createAWSTagsEndpoint = "createAwsTags"
	updateAWSTagsEndpoint = "updateAwsTags"
	deleteAWSTagsEndpoint = "deleteAwsTags"
	getTagsFromDBEndpoint = "getTagsFromDb"

	deleteExceptionEndpoint = "DeleteException"
	updateExceptionEndpoint = "UpdateException"

	postAllJobConfig            = "PostAllConfigJobs"
	getAllSourceConfigsEndpoint = "GetAllSourceConfigs"
	getAllSourcesEndpoint       = "GetAllSources"

	getJIRAURLsEndpoint = "GetJiraUrls"

	getAllExceptionsEndpoint     = "GetAllExcepts"
	getAllExceptionTypeEndpoints = "GetAllExceptTypes"
	deleteExceptionEndpoints     = "DeleteException"
	updateExceptionEndpoints     = "UpdateException"
	createExceptionEndpoints     = "CreateException"
)

// handleRequest unmarshals the apiRequest body into the provided interface, and logs that fact to the db
func handleRequest(w http.ResponseWriter, r *http.Request, endpoint string) (bearerToken string, ep endpoint, originalBody string, err error) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var body []byte
	var bearerTokens []string
	var exists bool
	var user domain.User
	bearerTokens, exists = r.Header["Authorization"]

	if exists && len(bearerTokens) >= 1 {
		bearerToken = bearerTokens[0]
		bearerToken = strings.TrimPrefix(bearerToken, "Bearer ")

		user, _, err = validateToken(bearerToken)
		if err == nil {
			body, err = ioutil.ReadAll(io.LimitReader(r.Body, mbSize))

			if err == nil {
				err = r.Body.Close()
				if err == nil {

					if r.Method == http.MethodGet || r.Method == http.MethodDelete {

					} else {
						originalBody = string(body)

						var req = &apiRequest{}
						err = json.Unmarshal(body, req)
						if err == nil {
							ep, err = grabEndpointFromRequest(req, endpoint, originalBody)
						} else {
							err = fmt.Errorf("error while unmarshalling req [%s]", err.Error())
						}
					}
				} else {
					err = fmt.Errorf("error while processing apiRequest [%s]", err.Error())
				}
			} else {
				err = fmt.Errorf("error while processing apiRequest [%s]", err.Error())
			}

			if err == nil {
				// TODO have some sort of flag that distinguishes between errors an activity logs
				err = createDBLogForRequest(body, endpoint, r, user)
			}

		} else {
			err = fmt.Errorf("error while processing session bearerToken [%s]", err.Error())
		}
	} else {
		err = errors.New("must include JWT as a bearer token in authorization field")
	}

	return bearerToken, ep, originalBody, err
}

func grabEndpointFromRequest(request *apiRequest, endpoint string, originalBody string) (ep endpoint, err error) {
	if request.Organization != nil {
		ep = request.Organization
	} else if request.Source != nil {
		ep = request.Source
	} else if request.Config != nil {
		ep = request.Config
	} else if request.History != nil {
		ep = request.History
	} else if request.Histories != nil {
		ep = request.Histories
	} else if request.Users != nil {
		ep = request.Users
	} else if request.Tag != nil {
		ep = request.Tag
	} else if request.JobConfigs != nil {
		ep = request.JobConfigs
	} else if request.Exceptions != nil {
		ep = request.Exceptions
	} else if request.BulkUpdateJob != nil {
		// do nothing, we'll just use the body of the apiRequest
	} else if request.BulkUpdateFile != nil {
		// do nothing, we'll just use the body of the apiRequest
	} else if request.Permission != nil {
		// do nothing, we'll just use the body of the apiRequest
	} else if endpoint == getLogsEndpoint {
		// do nothing, we'll just use the body of the apiRequest
	} else {
		err = fmt.Errorf("could not identify apiRequest type - %s", originalBody)
	}

	return ep, err
}

func createDBLogForRequest(body []byte, endpoint string, r *http.Request, user domain.User) (err error) {
	var bodyString = string(body)
	// Bulk updates are too long to store in the DB, just take the description line
	if endpoint == bulkUpdateUploadEndpoint {
		bodyString, err = modifyDBEntryForBulkUpdate(body)
	}
	var messageToLog string
	messageToLog = createMessageToLog(bodyString, messageToLog, r)
	_, _, err = Ms.CreateDBLog(sord(user.Username()), messageToLog, endpoint)
	if err != nil {
		fmt.Printf("error while creating database log [%s]\n", err.Error())
	}
	return err
}

func createMessageToLog(bodyString string, messageToLog string, r *http.Request) string {
	if strings.Contains(strings.ToLower(bodyString), "\"password\":") {
		messageToLog = fmt.Sprintf("%s - apiRequest body contained sensitive information and was removed", r.Method)
	} else {
		messageToLog = fmt.Sprintf("%s - %s", r.Method, bodyString)
	}
	return messageToLog
}

func modifyDBEntryForBulkUpdate(body []byte) (bodyString string, err error) {
	var foundDescLine = false
	var req = &apiRequest{}
	err = json.Unmarshal(body, req)
	if err == nil {
		if req.BulkUpdateFile != nil {
			var descLineIndex = strings.Index(req.BulkUpdateFile.Contents, "\n")
			if descLineIndex > 0 {
				bodyString = req.BulkUpdateFile.Contents[0:descLineIndex]
				foundDescLine = true
			}
		}

	}
	if !foundDescLine {
		bodyString = "Could not grab description line from bulk update req"
	}
	return bodyString, err
}

func respondToUserWithStatusCode(user domain.User, resp *GeneralResp, w http.ResponseWriter, wrapper errorWrapper, endpoint string, status int) {
	if wrapper.detailedError != nil {
		var username string
		if user != nil {
			username = sord(user.Username())
		} else {
			username = unauthenticated
		}

		_, _, err := Ms.CreateDBLog(username, wrapper.detailedError.Error(), endpoint)
		if err == nil {
			if wrapper.friendlyError != nil {
				resp.Message = wrapper.friendlyError.Error()
			} else {
				resp.Message = "could not ascertain error information"
			}
		}
	}

	var respByte []byte
	var err error
	respByte, err = json.Marshal(resp)
	if err == nil {
		w.WriteHeader(status)
		_, _ = fmt.Fprintln(w, string(respByte))
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(err)
	}
}

func respondToUser(user domain.User, generalResp *GeneralResp, w http.ResponseWriter, wrapper errorWrapper, endpoint string) {
	respondToUserWithStatusCode(user, generalResp, w, wrapper, endpoint, http.StatusOK)
}

func validateToken(sessionToken string) (user domain.User, permission domain.Permission, err error) {
	var session domain.Session
	session, err = Ms.GetSessionByToken(sessionToken)
	if err == nil {
		if session != nil {
			if !session.IsDisabled() {
				user, err = Ms.GetUserAnyOrg(session.UserID())
				if err == nil {
					if !user.IsDisabled() {

						user, err = encryptOrDecryptUser(user, crypto.DecryptMode)
						if err == nil {
							//var claims *customClaims
							_, err = checkJWT(sessionToken)
							if err == nil {
								permission, err = gatherHierarchicalPermissions(user.ID(), session.OrgID())
							}
						} else {
							err = fmt.Errorf("error while decrypting user information - %s", err.Error())
						}

					} else {
						err = errors.New("user is marked as disabled")
					}
				} else {
					err = errors.New("error while retrieving user information [" + err.Error() + "]")
				}
			} else {
				err = errors.New("session marked as disabled")
			}
		} else {
			err = errors.New("could not find matching session token [" + sessionToken + "] in the database")
		}
	} else {
		err = errors.New("error while retrieving session from database [" + err.Error() + "]")
	}
	return user, permission, err
}

type permissionWithParent struct {
	domain.Permission
	parent domain.Permission
}

// ParentOrgPermission creates a linked list of organizations from parent to child
func (p *permissionWithParent) ParentOrgPermission() domain.Permission {
	return p.parent
}

func gatherHierarchicalPermissions(userID string, orgID string) (basePermission domain.Permission, err error) {
	if len(userID) > 0 && len(orgID) > 0 {

		var org domain.Organization
		org, err = Ms.GetOrganizationByID(orgID)

		var nextPermission domain.Permission
		for org != nil && err == nil {
			nextPermission, err = Ms.GetPermissionByUserOrgID(userID, org.ID())
			if err == nil {
				if basePermission == nil {
					basePermission = nextPermission
				} else {
					var traverse = basePermission
					for traverse.ParentOrgPermission() != nil {
						traverse = traverse.ParentOrgPermission()
					}
					// TODO not sure if this works - should test before pushing
					//traverse.SetParentOrgPermission(nextPermission)
					traverse = &permissionWithParent{traverse, nextPermission}
				}

				if len(sord(org.ParentOrgID())) > 0 {
					org, err = Ms.GetOrganizationByID(sord(org.ParentOrgID()))
				} else {
					break
				}
			}
		}

	} else {
		err = errors.New("could not properly retrieve user information off session token")
	}

	return basePermission, err
}
