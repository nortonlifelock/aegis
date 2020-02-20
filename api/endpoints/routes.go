package endpoints

import (
	"net/http"

	"github.com/gorilla/mux"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

var routes = []route{
	{loginEndpoint, http.MethodPost, "/login", login},
	{loginEndpoint, http.MethodGet, "/login/{org}", updateSessionPermissions},
	{logoutEndpoint, http.MethodGet, "/logout", logout},

	{getHistoryByIDEndpoint, http.MethodGet, "/JobHistory/{" + idParam + "}", getJobHistoryByID},
	{createHistoryEndpoint, http.MethodPost, "/JobHistory", createJobHistory},
	{updateHistoryEndpoint, http.MethodPut, "/JobHistory", updateJobHistory},
	{deleteHistoryEndpoint, http.MethodDelete, "/JobHistory/{" + idParam + "}", deleteJobHistory},
	{getHistoriesEndpoint, http.MethodPost, "/AllHistories", getAllJobHistories},

	{createJConfigEndpoint, http.MethodPost, "/Config", createJobConfig},
	{updateJConfigEndpoint, http.MethodPut, "/Config", updateJobConfig},
	{deleteJConfigEndpoint, http.MethodDelete, "/Config/{" + idParam + "}", deleteJobConfig},
	{postAllJobConfig, http.MethodPost, "/AllJobConfigs", getAllJobConfigsWithOrder},
	{getJCAuditHistoryEndpoint, http.MethodGet, "/Config/Audit/{" + idParam + "}", getAuditHistoryForJobConfig},

	{createSourceEndpoint, http.MethodPost, "/Source", createSourceConfig},
	{updateSourceEndpoint, http.MethodPut, "/Source", updateSourceConfig},
	{deleteSourceEndpoint, http.MethodDelete, "/Source", deleteSourceConfig},
	{getAllSourceConfigsEndpoint, http.MethodGet, "/SourceConfigs", getAllSourceConfigs},
	{getAllSourcesEndpoint, http.MethodGet, "/Sources", getAllSources},

	{createOrgEndpoint, http.MethodPost, "/Org", createOrg},
	{updateOrgEndpoint, http.MethodPut, "/Org", updateOrg},
	{deleteOrgEndpoint, http.MethodDelete, "/Org", deleteOrg},
	{getAllOrgsEndpoint, http.MethodGet, "/Org", getAllOrganizations},
	{getOrgForUserEndpoint, http.MethodGet, "/Org/User", getOrgForUser},
	{getMyOrgEndpoint, http.MethodGet, "/Org/Me", getMyOrg},

	{createUserEndpoint, http.MethodPost, "/User", createUser},
	{updateUserEndpoint, http.MethodPut, "/User", updateUser},
	{deleteUserEndpoint, http.MethodDelete, "/User", deleteUser},

	{getAllUsersEndpoint, http.MethodGet, "/Users", getAllUsers},
	{getUserByIDEndpoint, http.MethodGet, "/User/{" + idParam + "}", getUserByID},

	{getUserPermEndpoint, http.MethodGet, "/User/{" + userParam + "}/{" + orgParam + "}/permission", getUserPermissionsByUserOrgID},
	{updateUserPermEndpoint, http.MethodPut, "/User/{" + userParam + "}/{" + orgParam + "}/permission", updateUserPermissionsByUserOrgID},
	{createUserPermEndpoint, http.MethodPost, "/User/{" + userParam + "}/{" + orgParam + "}/permission", createPermissionsForUser},
	{getUsersNameEndpoint, http.MethodGet, "/UserName", getMyName},

	{getPermListEndpoint, http.MethodGet, "/permission", getPermissionList},
	{createUserPermEndpoint, http.MethodPost, "/permission/{" + userParam + "}/{" + orgParam + "}", getPermissionList},

	{getAllJobsEndpoint, http.MethodGet, "/Jobs", getAllJobs},
	{getAllJobConfigsEndpoint, http.MethodGet, "/JobConfigs", getAllJobConfigs},

	{bulkUpdateEndpoint, http.MethodPost, "/BulkUpdate", createBulkUpdateJobHistory},
	{bulkUpdateEndpoint, http.MethodPost, "/BulkUpdateUpload", bulkUpdateUpload},

	{getAllLogTypesEndpoint, http.MethodGet, "/LogTypes", getLogTypes},
	{getLogsEndpoint, http.MethodPost, "/Logs", getLogs},

	{getScansEndpoint, http.MethodGet, "/Scans/{" + scannerParam + "}", getScansForScanner},
	{getAgsEndpoint, http.MethodGet, "/Scans/AG/{" + ticketParam + "}", getAGForTicket},

	{getVulnBySourceEndpoint, http.MethodGet, "/vuln/{" + sourceParam + "}", getVulnerabilitiesBySource},
	{getMatchedVulnsEndpoint, http.MethodGet, "/vulnMatch", getMatchedVulns},

	{getTicketStatusCountEndpoint, http.MethodGet, "/ticketCountByStatus", getCountOfJiraTicketsInStatus},

	{createWSConnectionEndpoint, http.MethodGet, "/connect/bulk", createWebsocketConnectionBulkUpdate},
	{burnDownWebSocketEndpoint, http.MethodGet, "/connect/burnDown", establishBurnDownWebsocket},
	{createScanWebsocketEndpoint, http.MethodGet, "/connect/scan", createWebsocketConnection},
	{internalWSEndpoint, http.MethodGet, "/connect/InternalWebSocket", listenForMessagesFromJobToSendToClient},

	{getAWSTagsEndpoint, http.MethodGet, "/tags/aws", getTagsFromAws},
	{getAzureTagsEndpoint, http.MethodGet, "/tags/azure", getTagsFromAzure},

	{createAWSTagsEndpoint, http.MethodPost, "/tags", createAwsTags},
	{updateAWSTagsEndpoint, http.MethodPut, "/tags", updateAwsTags},
	{deleteAWSTagsEndpoint, http.MethodPut, "/tags/delete", deleteAwsTags},
	{getTagsFromDBEndpoint, http.MethodGet, "/tags/existing/{" + sourceParam + "}", getTagsFromDb},

	{getFieldsForProjectEndpoint, http.MethodGet, "/tags/jira", getFieldsForJiraProject},

	{getJIRAURLsEndpoint, http.MethodGet, "/jira/urls", getJiraUrls},
	{getStatusMapsEndpoint, http.MethodGet, "/jira/statuses/{" + idParam + "}", getStatusMaps},
	{getFieldMapsEndpoint, http.MethodGet, "/jira/fields/{" + idParam + "}", getFieldMaps},

	{getSourceInByJobIDEndpoint, http.MethodPost, "/SrcInsByJobName", getSourceInsByJobID},
	{getSourceOutByJobIDEndpoint, http.MethodPost, "/SrcOutsByJobNameAndSrcIn", getSourceOutsByJobIDAndSrcIn},

	{getAllExceptionsEndpoint, http.MethodPost, "/Exceptions", getAllExceptions},
	//{getAllExceptionTypeEndpoints, http.MethodGet, "/ExceptTypes", getAllExceptTypes},
	//{createExceptionEndpoints, http.MethodPost, "/Exception", createException},
	//{updateExceptionEndpoints, http.MethodPut, "/Exception", updateException},
	//{deleteExceptionEndpoints, http.MethodPost, "/DeleteException", deleteException},
}

// NewRouter registers all the endpoints with their associated handler functions. It returns a mux router which can then be used to start
// the server
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, r := range routes {
		router.
			Methods(r.Method).
			Path(r.Pattern).
			Name(r.Name).
			Handler(r.HandlerFunc)
	}

	return router
}
