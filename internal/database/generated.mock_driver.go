package database

//**********************************************************
// GENERATED CODE - DO NOT CHANGE
// This file is generated using scaffolding. Any changes to
// this file will be overwritten on the next build
//**********************************************************

import (
	"encoding/json"
	"github.com/nortonlifelock/aegis/pkg/domain"
	"time"
)

//**********************************************************
// Struct Declaration
//**********************************************************

// MockSQLDriver defines the struct that implements the MockSQLDriver interface
type MockSQLDriver struct {
	domain.DatabaseConnection
	FuncCleanUp                                       func() (id int, affectedRows int, err error)
	FuncCreateAssetGroup                              func(inOrgID string, _GroupID string, _ScannerSourceID string, _ScannerSourceConfigID string) (id int, affectedRows int, err error)
	FuncCreateAssetWithIPInstanceID                   func(_State string, _IP string, _MAC string, _SourceID string, _InstanceID string, _Region string, _OrgID string, _OS string, _OsTypeID int) (id int, affectedRows int, err error)
	FuncCreateCategory                                func(_Category string) (id int, affectedRows int, err error)
	FuncCreateDBLog                                   func(_User string, _Command string, _Endpoint string) (id int, affectedRows int, err error)
	FuncCreateDetection                               func(_OrgID string, _SourceID string, _DeviceID string, _VulnID string, _IgnoreID string, _AlertDate time.Time, _LastFound time.Time, _LastUpdated time.Time, _Proof string, _Port int, _Protocol string, _ActiveKernel int, _DetectionStatusID int, _TimesSeen int, _DefaultTime time.Time) (id int, affectedRows int, err error)
	FuncCreateDevice                                  func(_AssetID string, _SourceID string, _Ip string, _Hostname string, inInstanceID string, _MAC string, _GroupID string, _OrgID string, _OS string, _OSTypeID int, inTrackingMethod string) (id int, affectedRows int, err error)
	FuncCreateException                               func(inSourceID string, inOrganizationID string, inTypeID int, inVulnerabilityID string, inDeviceID string, inDueDate time.Time, inApproval string, inActive bool, inPort string, inCreatedBy string) (id int, affectedRows int, err error)
	FuncCreateJobConfig                               func(_JobID int, _OrganizationID string, _PriorityOverride int, _Continuous bool, _WaitInSeconds int, _MaxInstances int, _AutoStart bool, _CreatedBy string, _DataInSourceID string, _DataOutSourceID string) (id int, affectedRows int, err error)
	FuncCreateJobConfigWPayload                       func(_JobID int, _OrganizationID string, _PriorityOverride int, _Continuous bool, _WaitInSeconds int, _MaxInstances int, _AutoStart bool, _CreatedBy string, _DataInSourceID string, _DataOutSourceID string, _Payload string) (id int, affectedRows int, err error)
	FuncCreateJobHistory                              func(_JobID int, _ConfigID string, _StatusID int, _Priority int, _Identifier string, _CurrentIteration int, _Payload string, _ThreadID string, _PulseDate time.Time, _CreatedBy string) (id int, affectedRows int, err error)
	FuncCreateJobHistoryWithParentID                  func(_JobID int, _ConfigID string, _StatusID int, _Priority int, _Identifier string, _CurrentIteration int, _Payload string, _ThreadID string, _PulseDate time.Time, _CreatedBy string, _ParentID string) (id int, affectedRows int, err error)
	FuncCreateOrganization                            func(_Code string, _Description string, _TimeZoneOffset float32, _UpdatedBy string) (id int, affectedRows int, err error)
	FuncCreateOrganizationWithPayloadEkey             func(_Code string, _Description string, _TimeZoneOffset float32, _Payload string, _EKEY string, _UpdatedBy string) (id int, affectedRows int, err error)
	FuncCreateScanSummary                             func(_SourceID string, _ScannerSourceConfigID string, _OrgID string, _ScanID string, _ScanStatus string, _ScanClosePayload string, _ParentJobID string) (id int, affectedRows int, err error)
	FuncCreateSourceConfig                            func(_Source string, _SourceID string, _OrganizationID string, _Address string, _Port string, _Username string, _Password string, _PrivateKey string, _ConsumerKey string, _Token string, _Payload string) (id int, affectedRows int, err error)
	FuncCreateTag                                     func(_DeviceID string, _TagKeyID string, _Value string) (id int, affectedRows int, err error)
	FuncCreateTagKey                                  func(_KeyValue string) (id int, affectedRows int, err error)
	FuncCreateTagMap                                  func(_TicketingSourceID string, _TicketingTag string, _CloudSourceID string, _CloudTag string, _Options string, _OrganizationID string) (id int, affectedRows int, err error)
	FuncCreateTicket                                  func(_Title string, _Status string, _DetectionID string, _OrganizationID string, _DueDate time.Time, _UpdatedDate time.Time, _ResolutionDate time.Time, _ExceptionDate time.Time, _DefaultTime time.Time) (id int, affectedRows int, err error)
	FuncCreateTicketingJob                            func(GroupID int, OrgID string, ScanStartDate string) (id int, affectedRows int, err error)
	FuncCreateUser                                    func(_Username string, _FirstName string, _LastName string, _Email string) (id int, affectedRows int, err error)
	FuncCreateUserPermissions                         func(_UserID string, _OrgID string) (id int, affectedRows int, err error)
	FuncCreateUserSession                             func(_UserID string, _OrgID string, _SessionKey string) (id int, affectedRows int, err error)
	FuncCreateVulnInfo                                func(_SourceVulnID string, _Title string, _SourceID string, _CVSSScore float32, _CVSS3Score float32, _Description string, _Threat string, _Solution string, _Software string, _Patchable string, _Category string, _DetectionInformation string) (id int, affectedRows int, err error)
	FuncCreateVulnRef                                 func(_VulnInfoID string, _SourceID string, _Reference string, _RefType int) (id int, affectedRows int, err error)
	FuncDeleteDecomIgnoreForDevice                    func(_sourceID string, _devID string, _orgID string) (id int, affectedRows int, err error)
	FuncDeleteIgnoreForDevice                         func(_sourceID string, _devID string, _orgID string) (id int, affectedRows int, err error)
	FuncDeleteSessionByToken                          func(_SessionKey string) (id int, affectedRows int, err error)
	FuncDeleteTagMap                                  func(_TicketingSourceID string, _TicketingTag string, _CloudSourceID string, _CloudTag string, _OrganizationID string) (id int, affectedRows int, err error)
	FuncDeleteUserByUsername                          func(_Username string) (id int, affectedRows int, err error)
	FuncDisableIgnore                                 func(inSourceID string, inDevID string, inOrgID string, inVulnID string, inPortID string, inUpdatedBy string) (id int, affectedRows int, err error)
	FuncDisableJobConfig                              func(_ID string, _UpdatedBy string) (id int, affectedRows int, err error)
	FuncDisableOrganization                           func(_ID string, _UpdatedBy string) (id int, affectedRows int, err error)
	FuncDisableSource                                 func(_ID string, _OrgID string, _UpdatedBy string) (id int, affectedRows int, err error)
	FuncGetAllDetectionInfo                           func(_OrgID string) ([]domain.DetectionInfo, error)
	FuncGetAllDeviceInfo                              func() ([]domain.DeviceInfo, error)
	FuncGetAllExceptions                              func(_offset int, _limit int, _sourceID string, _orgID string, _typeID int, _vulnID string, _devID string, _dueDate time.Time, _port string, _approval string, _active bool, _dBCreatedDate time.Time, _dBUpdatedDate time.Time, _updatedBy string, _createdBy string, _sortField string, _sortOrder string) ([]domain.Ignore, error)
	FuncGetAllJobConfigs                              func(_OrgID string) ([]domain.JobConfig, error)
	FuncGetAllJobConfigsWithOrder                     func(_offset int, _limit int, _configID string, _jobid int, _dataInSourceConfigID string, _dataOutSourceConfigID string, _priorityOverride int, _continuous bool, _Payload string, _waitInSeconds int, _maxInstances int, _autoStart bool, _OrgID string, _updatedBy string, _createdBy string, _sortField string, _sortOrder string, _updatedDate time.Time, _createdDate time.Time, _lastJobStart time.Time, _ID string) ([]domain.JobConfig, error)
	FuncGetAssetGroup                                 func(inOrgID string, _GroupID string, _ScannerConfigSourceID string) (domain.AssetGroup, error)
	FuncGetAssetGroupForOrg                           func(inScannerSourceConfigID string, inOrgID string) ([]domain.AssetGroup, error)
	FuncGetAssetGroupForOrgNoScanner                  func(inOrgID string, inGroupID string) (domain.AssetGroup, error)
	FuncGetAssetGroupsByCloudSource                   func(inOrgID string, inCloudSourceID string) ([]domain.AssetGroup, error)
	FuncGetAssetGroupsForOrg                          func(inOrgID string) ([]domain.AssetGroup, error)
	FuncGetAssignmentGroupByIP                        func(_SourceID string, _OrganizationID string, _IP string) ([]domain.AssignmentGroup, error)
	FuncGetAssignmentGroupByOrgIP                     func(_OrganizationID string, _IP string) ([]domain.AssignmentGroup, error)
	FuncGetAssignmentRulesByOrg                       func(_OrganizationID string) ([]domain.AssignmentRules, error)
	FuncGetAutoStartJobs                              func() ([]domain.JobConfig, error)
	FuncGetCISAssignments                             func(_OrganizationID string) ([]domain.CISAssignments, error)
	FuncGetCancelledJobs                              func() ([]domain.JobHistory, error)
	FuncGetCategoryByName                             func(_Name string) ([]domain.Category, error)
	FuncGetCategoryRules                              func(_OrgID string, _SourceID string) ([]domain.CategoryRule, error)
	FuncGetDetectionInfo                              func(_DeviceID string, _VulnerabilityID string, _Port int, _Protocol string) (domain.DetectionInfo, error)
	FuncGetDetectionInfoAfter                         func(_After time.Time, _OrgID string) ([]domain.DetectionInfo, error)
	FuncGetDetectionInfoByID                          func(_ID string, _OrgID string) (domain.DetectionInfo, error)
	FuncGetDetectionInfoBySourceVulnID                func(_SourceDeviceID string, _SourceVulnerabilityID string, _Port int, _Protocol string) (domain.DetectionInfo, error)
	FuncGetDetectionInfoForDeviceID                   func(inDeviceID string, _OrgID string, ticketInactiveKernels bool) ([]domain.DetectionInfo, error)
	FuncGetDetectionInfoForGroupAfter                 func(_LastUpdatedAfter time.Time, _LastFoundAfter time.Time, _OrgID string, inGroupID string, ticketInactiveKernels bool) ([]domain.DetectionInfo, error)
	FuncGetDetectionStatusByID                        func(_ID int) (domain.DetectionStatus, error)
	FuncGetDetectionStatusByName                      func(_Name string) (domain.DetectionStatus, error)
	FuncGetDetectionStatuses                          func() ([]domain.DetectionStatus, error)
	FuncGetDetectionsInfoForDevice                    func(_DeviceID string) ([]domain.DetectionInfo, error)
	FuncGetDeviceInfoByAssetIDNoOrg                   func(inAssetID string) (domain.DeviceInfo, error)
	FuncGetDeviceInfoByAssetOrgID                     func(inAssetID string, inOrgID string) (domain.DeviceInfo, error)
	FuncGetDeviceInfoByCloudSourceIDAndIP             func(_IP string, _CloudSourceID string, _OrgID string) ([]domain.DeviceInfo, error)
	FuncGetDeviceInfoByGroupIP                        func(inIP string, inGroupID string, inOrgID string) (domain.DeviceInfo, error)
	FuncGetDeviceInfoByIP                             func(_IP string, _OrgID string) (domain.DeviceInfo, error)
	FuncGetDeviceInfoByIPMACAndRegion                 func(_IP string, _MAC string, _Region string, _OrgID string) (domain.DeviceInfo, error)
	FuncGetDeviceInfoByInstanceID                     func(_InstanceID string, _OrgID string) ([]domain.DeviceInfo, error)
	FuncGetDeviceInfoByScannerSourceID                func(_IP string, _GroupID string, _OrgID string) (domain.DeviceInfo, error)
	FuncGetDevicesInfoByCloudSourceID                 func(_CloudSourceID string, _OrgID string) ([]domain.DeviceInfo, error)
	FuncGetDevicesInfoBySourceID                      func(_SourceID string, _OrgID string) ([]domain.DeviceInfo, error)
	FuncGetExceptionByVulnIDOrg                       func(_DeviceID string, _VulnID string, _OrgID string) (domain.Ignore, error)
	FuncGetExceptionDetections                        func(_offset int, _limit int, _orgID string, _sortField string, _sortOrder string, _Title string, _IP string, _Hostname string, _VulnID string, _Approval string, _DueDate string, _AssignmentGroup string, _OS string, _OSRegex string, _TypeID int) ([]domain.ExceptedDetection, error)
	FuncGetExceptionTypes                             func() ([]domain.ExceptionType, error)
	FuncGetExceptionsByOrg                            func(_OrgID string) ([]domain.Ignore, error)
	FuncGetExceptionsDueNext30Days                    func() ([]domain.CERF, error)
	FuncGetExceptionsLength                           func(_offset int, _limit int, _orgID string, _sortField string, _sortOrder string, _Title string, _IP string, _Hostname string, _VulnID string, _Approval string, _DueDate string, _AssignmentGroup string, _OS string, _OSRegex string, _TypeID int) (domain.QueryData, error)
	FuncGetGlobalExceptions                           func(_OrgID string) ([]domain.Ignore, error)
	FuncGetJobByID                                    func(_ID int) (domain.JobRegistration, error)
	FuncGetJobConfig                                  func(_ID string) (domain.JobConfig, error)
	FuncGetJobConfigAudit                             func(inJobConfigID string, inOrgID string) ([]domain.JobConfigAudit, error)
	FuncGetJobConfigByID                              func(_ID string, _OrgID string) (domain.JobConfig, error)
	FuncGetJobConfigByJobHistoryID                    func(_JobHistoryID string) (domain.JobConfig, error)
	FuncGetJobConfigByOrgIDAndJobID                   func(_OrgID string, _JobID int) ([]domain.JobConfig, error)
	FuncGetJobConfigByOrgIDAndJobIDWithSC             func(_OrgID string, _JobID int, _SourceConfigID string) ([]domain.JobConfig, error)
	FuncGetJobConfigLength                            func(_configID string, _jobID int, _dataInSourceConfigID string, _dataOutSourceConfigID string, _priorityOverride int, _continuous bool, _Payload string, _waitInSeconds int, _maxInstances int, _autoStart bool, _OrgID string, _updatedBy string, _createdBy string, _updatedDate time.Time, _createdDate time.Time, _lastJobStart time.Time, _ID string) (domain.QueryData, error)
	FuncGetJobHistories                               func(_offset int, _limit int, _jobID int, _jobconfig string, _status int, _Payload string, _OrgID string) ([]domain.JobHistory, error)
	FuncGetJobHistoryByID                             func(_ID string) (domain.JobHistory, error)
	FuncGetJobHistoryLength                           func(_jobid int, _jobconfig string, _status int, _Payload string, _orgid string) (domain.QueryData, error)
	FuncGetJobQueueByStatusID                         func(_StatusID int) ([]domain.JobHistory, error)
	FuncGetJobs                                       func() ([]domain.JobRegistration, error)
	FuncGetJobsByStruct                               func(_Struct string) (domain.JobRegistration, error)
	FuncGetLeafOrganizationsForUser                   func(_UserID string) ([]domain.Organization, error)
	FuncGetLogTypes                                   func() ([]domain.LogType, error)
	FuncGetLogsByParams                               func(_MethodOfDiscovery string, _jobType int, _logType int, _jobHistoryID string, _fromDate time.Time, _toDate time.Time, _OrgID string) ([]domain.DBLog, error)
	FuncGetMatchedVulns                               func() ([]domain.VulnerabilityMatch, error)
	FuncGetOperatingSystemType                        func(_OS string) (domain.OperatingSystemType, error)
	FuncGetOrganizationByCode                         func(Code string) (domain.Organization, error)
	FuncGetOrganizationByID                           func(ID string) (domain.Organization, error)
	FuncGetOrganizations                              func() ([]domain.Organization, error)
	FuncGetPendingActiveCloudDecomJob                 func(_OrgID string) ([]domain.JobHistory, error)
	FuncGetPendingActiveRescanJob                     func(_OrgID string) ([]domain.JobHistory, error)
	FuncGetPermissionByUserOrgID                      func(_UserID string, _OrgID string) (domain.Permission, error)
	FuncGetPermissionOfLeafOrgByUserID                func(_UserID string) (domain.Permission, error)
	FuncGetRecentlyUpdatedScanSummaries               func(_OrgID string) ([]domain.ScanSummary, error)
	FuncGetScanSummariesBySourceName                  func(_OrgID string, _SourceName string) ([]domain.ScanSummary, error)
	FuncGetScanSummary                                func(_SourceID string, _OrgID string, _ScanID string) (domain.ScanSummary, error)
	FuncGetScanSummaryBySourceKey                     func(_SourceKey string) (domain.ScanSummary, error)
	FuncGetScheduledJobsToStart                       func(_LastChecked time.Time) ([]domain.JobSchedule, error)
	FuncGetSessionByToken                             func(_SessionKey string) (domain.Session, error)
	FuncGetSourceByID                                 func(_ID string) (domain.Source, error)
	FuncGetSourceByName                               func(_Source string) (domain.Source, error)
	FuncGetSourceConfigByID                           func(_ID string) (domain.SourceConfig, error)
	FuncGetSourceConfigByNameOrg                      func(_Source string, _OrgID string) ([]domain.SourceConfig, error)
	FuncGetSourceConfigByOrgID                        func(_OrgID string) ([]domain.SourceConfig, error)
	FuncGetSourceConfigBySourceID                     func(_OrgID string, _SourceID string) ([]domain.SourceConfig, error)
	FuncGetSourceInsByJobID                           func(inJob int, inOrgID string) ([]domain.SourceConfig, error)
	FuncGetSourceOauthByOrgURL                        func(_URL string, _OrgID string) (domain.SourceConfig, error)
	FuncGetSourceOauthByURL                           func(_URL string) (domain.SourceConfig, error)
	FuncGetSourceOutsByJobID                          func(inJob int, inOrgID string) ([]domain.SourceConfig, error)
	FuncGetSources                                    func() ([]domain.Source, error)
	FuncGetTagByDeviceAndTagKey                       func(_DeviceID string, _TagKeyID string) (domain.Tag, error)
	FuncGetTagKeyByID                                 func(_ID string) (domain.TagKey, error)
	FuncGetTagKeyByKey                                func(_KeyValue string) (domain.TagKey, error)
	FuncGetTagMapsByOrg                               func(_OrganizationID string) ([]domain.TagMap, error)
	FuncGetTagMapsByOrgCloudSourceID                  func(_CloudID string, _OrganizationID string) ([]domain.TagMap, error)
	FuncGetTagsForDevice                              func(_DeviceID string) ([]domain.Tag, error)
	FuncGetTicketByDetectionID                        func(inDetectionID string, _OrgID string) (domain.TicketSummary, error)
	FuncGetTicketByDeviceIDVulnID                     func(inDeviceID string, inVulnID string, inPort int, inProtocol string, inOrgID string) (domain.TicketSummary, error)
	FuncGetTicketByIPGroupIDVulnID                    func(inIP string, inGroupID string, inVulnID string, inPort int, inProtocol string, inOrgID string) (domain.TicketSummary, error)
	FuncGetTicketByTitle                              func(_Title string, _OrgID string) (domain.TicketSummary, error)
	FuncGetTicketCountByStatus                        func(inStatus string, inOrgID string) (domain.QueryData, error)
	FuncGetTicketCreatedAfter                         func(_UpperCVSS float32, _LowerCVSS float32, _CreatedAfter time.Time, _OrgID string) ([]domain.TicketSummary, error)
	FuncGetTicketTrackingMethod                       func(_Title string, _OrgID string) (domain.KeyValue, error)
	FuncGetUnfinishedScanSummariesBySourceConfigOrgID func(_ScannerSourceConfigID string, _OrgID string) ([]domain.ScanSummary, error)
	FuncGetUnfinishedScanSummariesBySourceOrgID       func(_SourceID string, _OrgID string) ([]domain.ScanSummary, error)
	FuncGetUnmatchedVulns                             func(_SourceID int) ([]domain.VulnerabilityInfo, error)
	FuncGetUserAnyOrg                                 func(_ID string) (domain.User, error)
	FuncGetUserByID                                   func(_ID string, _OrgID string) (domain.User, error)
	FuncGetUserByUsername                             func(_Username string) (domain.User, error)
	FuncGetUsersByOrg                                 func(_OrgID string) ([]domain.User, error)
	FuncGetVulnInfoByID                               func(_ID string) (domain.VulnerabilityInfo, error)
	FuncGetVulnInfoBySource                           func(_Source string) ([]domain.VulnerabilityInfo, error)
	FuncGetVulnInfoBySourceID                         func(_SourceID string) ([]domain.VulnerabilityInfo, error)
	FuncGetVulnInfoBySourceVulnID                     func(_SourceVulnID string) (domain.VulnerabilityInfo, error)
	FuncGetVulnInfoBySourceVulnIDSourceID             func(_SourceVulnID string, _SourceID string, _Modified time.Time) (domain.VulnerabilityInfo, error)
	FuncGetVulnRefInfo                                func(_VulnInfoID string, _SourceID string, _Reference string) (domain.VulnerabilityReferenceInfo, error)
	FuncGetVulnRefInfoVendor                          func(_VulnInfoID string, _SourceID string) ([]domain.VulnerabilityReferenceInfo, error)
	FuncGetVulnReferencesInfo                         func(_VulnInfoID string, _SourceID string) ([]domain.VulnerabilityReferenceInfo, error)
	FuncGetVulnReferencesInfoBySourceAndRef           func(_SourceID string, _Reference string) ([]domain.VulnerabilityReferenceInfo, error)
	FuncHasDecommissioned                             func(_devID string, _sourceID string, _orgID string) (domain.Ignore, error)
	FuncHasExceptionOrFalsePositive                   func(_sourceID string, _vulnID string, _devID string, _orgID string, _port string, _OS string) ([]domain.Ignore, error)
	FuncHasIgnore                                     func(inSourceID string, inVulnID string, inDevID string, inOrgID string, inPort string, inMostCurrentDetection time.Time) (domain.Ignore, error)
	FuncPulseJob                                      func(_JobHistoryID string) (id int, affectedRows int, err error)
	FuncRemoveExpiredIgnoreIDs                        func(_OrgID string) (id int, affectedRows int, err error)
	FuncSaveAssignmentGroup                           func(_SourceID string, _OrganizationID string, _IpAddress string, _GroupName string) (id int, affectedRows int, err error)
	FuncSaveIgnore                                    func(_SourceID string, _OrganizationID string, _TypeID int, _VulnerabilityID string, _DeviceID string, _DueDate time.Time, _Approval string, _Active bool, _port string) (id int, affectedRows int, err error)
	FuncSaveScanSummary                               func(_ScanID string, _ScanStatus string) (id int, affectedRows int, err error)
	FuncSetScheduleLastRun                            func(_ID string) (id int, affectedRows int, err error)
	FuncUpdateAssetGroupLastTicket                    func(inGroupID string, inOrgID string, inLastTicketTime time.Time) (id int, affectedRows int, err error)
	FuncUpdateAssetIDOsTypeIDOfDevice                 func(_ID string, _AssetID string, _ScannerSourceID string, _GroupID string, _OS string, _HostName string, _OsTypeID int, inTrackingMethod string, _OrgID string) (id int, affectedRows int, err error)
	FuncUpdateDetection                               func(_ID string, _DeviceID string, _VulnID string, _Port int, _Protocol string, _ExceptionID string, _TimesSeen int, _StatusID int, _LastFound time.Time, _LastUpdated time.Time, _DefaultTime time.Time) (id int, affectedRows int, err error)
	FuncUpdateDetectionIgnore                         func(_DeviceID string, _VulnID string, _Port int, _Protocol string, _ExceptionID string) (id int, affectedRows int, err error)
	FuncUpdateExpirationDateByCERF                    func(_CERForm string, _OrganizationID string, _DueDate time.Time) (id int, affectedRows int, err error)
	FuncUpdateInstanceIDOfDevice                      func(_ID string, _InstanceID string, _CloudSourceID string, _State string, _Region string, _OrgID string) (id int, affectedRows int, err error)
	FuncUpdateJobConfig                               func(_ID string, _DataInSourceID string, _DataOutSourceID string, _Autostart bool, _PriorityOverride int, _Continuous bool, _WaitInSeconds int, _MaxInstances int, _UpdatedBy string, _OrgID string) (id int, affectedRows int, err error)
	FuncUpdateJobConfigLastRun                        func(_ID string, _LastRun time.Time) (id int, affectedRows int, err error)
	FuncUpdateJobHistory                              func(_ID string, _ConfigID string, _Payload string, _UpdatedBy string) (id int, affectedRows int, err error)
	FuncUpdateJobHistoryStatus                        func(_ID string, _Status int) (id int, affectedRows int, err error)
	FuncUpdateJobHistoryStatusDetailed                func(_ID string, _Status int, _UpdatedBy string) (id int, affectedRows int, err error)
	FuncUpdateOrganization                            func(_ID string, _Description string, _TimezoneOffset float32, _UpdatedBy string) (id int, affectedRows int, err error)
	FuncUpdatePermissionsByUserOrgID                  func(_UserID string, _OrgID string, _Admin bool, _Manager bool, _Reader bool, _Reporter bool) (id int, affectedRows int, err error)
	FuncUpdateSourceConfig                            func(_ID string, _OrgID string, _Address string, _Username string, _Password string, _PrivateKey string, _ConsumerKey string, _Token string, _Port string, _Payload string, _UpdatedBy string) (id int, affectedRows int, err error)
	FuncUpdateSourceConfigConcurrencyByID             func(_ID string, _Delay int, _Retries int, _Concurrency int) (id int, affectedRows int, err error)
	FuncUpdateSourceConfigToken                       func(_ID string, _Token string) (id int, affectedRows int, err error)
	FuncUpdateStateOfDevice                           func(_ID string, _State string, _OrgID string) (id int, affectedRows int, err error)
	FuncUpdateTag                                     func(_DeviceID string, _TagKeyID string, _Value string) (id int, affectedRows int, err error)
	FuncUpdateTagMap                                  func(_TicketingSourceID string, _TicketingTag string, _CloudSourceID string, _CloudTag string, _Options string, _OrganizationID string) (id int, affectedRows int, err error)
	FuncUpdateTicket                                  func(_Title string, _Status string, _OrganizationID string, _AssignmentGroup string, _Assignee string, _DueDate time.Time, _CreatedDate time.Time, _UpdatedDate time.Time, _ResolutionDate time.Time, _ExceptionDate time.Time, _DefaultTime time.Time) (id int, affectedRows int, err error)
	FuncUpdateTicketDetectionID                       func(_Title string, _DetectionID string, _OrganizationID string) (id int, affectedRows int, err error)
	FuncUpdateUserByID                                func(_ID string, _FirstName string, _LastName string, _Email string, _Disabled bool) (id int, affectedRows int, err error)
	FuncUpdateVulnByID                                func(_ID string, _SourceVulnID string, _Title string, _SourceID string, _CVSSScore float32, _CVSS3Score float32, _Description string, _Threat string, _Solution string, _Software string, _Patchable string, _Category string, _DetectionInformation string) (id int, affectedRows int, err error)
	FuncUpdateVulnByIDNoCVSS3                         func(_ID string, _SourceVulnID string, _Title string, _SourceID string, _CVSSScore float32, _CVSS3Score float32, _Description string, _Threat string, _Solution string, _Software string, _Patchable string, _DetectionInformation string) (id int, affectedRows int, err error)
	FuncUpdateVulnInfoID                              func(_VulnInfoID string, _VulnID string, _MatchConfidence int, _MatchReasons string) (id int, affectedRows int, err error)
}

//**********************************************************
// Struct Methods
//**********************************************************

// MarshalJSON marshals the struct by converting it to a map
func (myMockSQLDriver MockSQLDriver) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"NA": "NA",
	})
}

func (myMockSQLDriver *MockSQLDriver) CleanUp() (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncCleanUp != nil {
		return myMockSQLDriver.FuncCleanUp()
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) CreateAssetGroup(inOrgID string, _GroupID string, _ScannerSourceID string, _ScannerSourceConfigID string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncCreateAssetGroup != nil {
		return myMockSQLDriver.FuncCreateAssetGroup(inOrgID, _GroupID, _ScannerSourceID, _ScannerSourceConfigID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) CreateAssetWithIPInstanceID(_State string, _IP string, _MAC string, _SourceID string, _InstanceID string, _Region string, _OrgID string, _OS string, _OsTypeID int) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncCreateAssetWithIPInstanceID != nil {
		return myMockSQLDriver.FuncCreateAssetWithIPInstanceID(_State, _IP, _MAC, _SourceID, _InstanceID, _Region, _OrgID, _OS, _OsTypeID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) CreateCategory(_Category string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncCreateCategory != nil {
		return myMockSQLDriver.FuncCreateCategory(_Category)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) CreateDBLog(_User string, _Command string, _Endpoint string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncCreateDBLog != nil {
		return myMockSQLDriver.FuncCreateDBLog(_User, _Command, _Endpoint)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) CreateDetection(_OrgID string, _SourceID string, _DeviceID string, _VulnID string, _IgnoreID string, _AlertDate time.Time, _LastFound time.Time, _LastUpdated time.Time, _Proof string, _Port int, _Protocol string, _ActiveKernel int, _DetectionStatusID int, _TimesSeen int, _DefaultTime time.Time) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncCreateDetection != nil {
		return myMockSQLDriver.FuncCreateDetection(_OrgID, _SourceID, _DeviceID, _VulnID, _IgnoreID, _AlertDate, _LastFound, _LastUpdated, _Proof, _Port, _Protocol, _ActiveKernel, _DetectionStatusID, _TimesSeen, _DefaultTime)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) CreateDevice(_AssetID string, _SourceID string, _Ip string, _Hostname string, inInstanceID string, _MAC string, _GroupID string, _OrgID string, _OS string, _OSTypeID int, inTrackingMethod string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncCreateDevice != nil {
		return myMockSQLDriver.FuncCreateDevice(_AssetID, _SourceID, _Ip, _Hostname, inInstanceID, _MAC, _GroupID, _OrgID, _OS, _OSTypeID, inTrackingMethod)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) CreateException(inSourceID string, inOrganizationID string, inTypeID int, inVulnerabilityID string, inDeviceID string, inDueDate time.Time, inApproval string, inActive bool, inPort string, inCreatedBy string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncCreateException != nil {
		return myMockSQLDriver.FuncCreateException(inSourceID, inOrganizationID, inTypeID, inVulnerabilityID, inDeviceID, inDueDate, inApproval, inActive, inPort, inCreatedBy)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) CreateJobConfig(_JobID int, _OrganizationID string, _PriorityOverride int, _Continuous bool, _WaitInSeconds int, _MaxInstances int, _AutoStart bool, _CreatedBy string, _DataInSourceID string, _DataOutSourceID string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncCreateJobConfig != nil {
		return myMockSQLDriver.FuncCreateJobConfig(_JobID, _OrganizationID, _PriorityOverride, _Continuous, _WaitInSeconds, _MaxInstances, _AutoStart, _CreatedBy, _DataInSourceID, _DataOutSourceID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) CreateJobConfigWPayload(_JobID int, _OrganizationID string, _PriorityOverride int, _Continuous bool, _WaitInSeconds int, _MaxInstances int, _AutoStart bool, _CreatedBy string, _DataInSourceID string, _DataOutSourceID string, _Payload string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncCreateJobConfigWPayload != nil {
		return myMockSQLDriver.FuncCreateJobConfigWPayload(_JobID, _OrganizationID, _PriorityOverride, _Continuous, _WaitInSeconds, _MaxInstances, _AutoStart, _CreatedBy, _DataInSourceID, _DataOutSourceID, _Payload)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) CreateJobHistory(_JobID int, _ConfigID string, _StatusID int, _Priority int, _Identifier string, _CurrentIteration int, _Payload string, _ThreadID string, _PulseDate time.Time, _CreatedBy string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncCreateJobHistory != nil {
		return myMockSQLDriver.FuncCreateJobHistory(_JobID, _ConfigID, _StatusID, _Priority, _Identifier, _CurrentIteration, _Payload, _ThreadID, _PulseDate, _CreatedBy)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) CreateJobHistoryWithParentID(_JobID int, _ConfigID string, _StatusID int, _Priority int, _Identifier string, _CurrentIteration int, _Payload string, _ThreadID string, _PulseDate time.Time, _CreatedBy string, _ParentID string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncCreateJobHistoryWithParentID != nil {
		return myMockSQLDriver.FuncCreateJobHistoryWithParentID(_JobID, _ConfigID, _StatusID, _Priority, _Identifier, _CurrentIteration, _Payload, _ThreadID, _PulseDate, _CreatedBy, _ParentID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) CreateOrganization(_Code string, _Description string, _TimeZoneOffset float32, _UpdatedBy string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncCreateOrganization != nil {
		return myMockSQLDriver.FuncCreateOrganization(_Code, _Description, _TimeZoneOffset, _UpdatedBy)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) CreateOrganizationWithPayloadEkey(_Code string, _Description string, _TimeZoneOffset float32, _Payload string, _EKEY string, _UpdatedBy string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncCreateOrganizationWithPayloadEkey != nil {
		return myMockSQLDriver.FuncCreateOrganizationWithPayloadEkey(_Code, _Description, _TimeZoneOffset, _Payload, _EKEY, _UpdatedBy)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) CreateScanSummary(_SourceID string, _ScannerSourceConfigID string, _OrgID string, _ScanID string, _ScanStatus string, _ScanClosePayload string, _ParentJobID string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncCreateScanSummary != nil {
		return myMockSQLDriver.FuncCreateScanSummary(_SourceID, _ScannerSourceConfigID, _OrgID, _ScanID, _ScanStatus, _ScanClosePayload, _ParentJobID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) CreateSourceConfig(_Source string, _SourceID string, _OrganizationID string, _Address string, _Port string, _Username string, _Password string, _PrivateKey string, _ConsumerKey string, _Token string, _Payload string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncCreateSourceConfig != nil {
		return myMockSQLDriver.FuncCreateSourceConfig(_Source, _SourceID, _OrganizationID, _Address, _Port, _Username, _Password, _PrivateKey, _ConsumerKey, _Token, _Payload)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) CreateTag(_DeviceID string, _TagKeyID string, _Value string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncCreateTag != nil {
		return myMockSQLDriver.FuncCreateTag(_DeviceID, _TagKeyID, _Value)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) CreateTagKey(_KeyValue string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncCreateTagKey != nil {
		return myMockSQLDriver.FuncCreateTagKey(_KeyValue)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) CreateTagMap(_TicketingSourceID string, _TicketingTag string, _CloudSourceID string, _CloudTag string, _Options string, _OrganizationID string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncCreateTagMap != nil {
		return myMockSQLDriver.FuncCreateTagMap(_TicketingSourceID, _TicketingTag, _CloudSourceID, _CloudTag, _Options, _OrganizationID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) CreateTicket(_Title string, _Status string, _DetectionID string, _OrganizationID string, _DueDate time.Time, _UpdatedDate time.Time, _ResolutionDate time.Time, _ExceptionDate time.Time, _DefaultTime time.Time) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncCreateTicket != nil {
		return myMockSQLDriver.FuncCreateTicket(_Title, _Status, _DetectionID, _OrganizationID, _DueDate, _UpdatedDate, _ResolutionDate, _ExceptionDate, _DefaultTime)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) CreateTicketingJob(GroupID int, OrgID string, ScanStartDate string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncCreateTicketingJob != nil {
		return myMockSQLDriver.FuncCreateTicketingJob(GroupID, OrgID, ScanStartDate)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) CreateUser(_Username string, _FirstName string, _LastName string, _Email string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncCreateUser != nil {
		return myMockSQLDriver.FuncCreateUser(_Username, _FirstName, _LastName, _Email)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) CreateUserPermissions(_UserID string, _OrgID string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncCreateUserPermissions != nil {
		return myMockSQLDriver.FuncCreateUserPermissions(_UserID, _OrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) CreateUserSession(_UserID string, _OrgID string, _SessionKey string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncCreateUserSession != nil {
		return myMockSQLDriver.FuncCreateUserSession(_UserID, _OrgID, _SessionKey)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) CreateVulnInfo(_SourceVulnID string, _Title string, _SourceID string, _CVSSScore float32, _CVSS3Score float32, _Description string, _Threat string, _Solution string, _Software string, _Patchable string, _Category string, _DetectionInformation string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncCreateVulnInfo != nil {
		return myMockSQLDriver.FuncCreateVulnInfo(_SourceVulnID, _Title, _SourceID, _CVSSScore, _CVSS3Score, _Description, _Threat, _Solution, _Software, _Patchable, _Category, _DetectionInformation)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) CreateVulnRef(_VulnInfoID string, _SourceID string, _Reference string, _RefType int) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncCreateVulnRef != nil {
		return myMockSQLDriver.FuncCreateVulnRef(_VulnInfoID, _SourceID, _Reference, _RefType)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) DeleteDecomIgnoreForDevice(_sourceID string, _devID string, _orgID string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncDeleteDecomIgnoreForDevice != nil {
		return myMockSQLDriver.FuncDeleteDecomIgnoreForDevice(_sourceID, _devID, _orgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) DeleteIgnoreForDevice(_sourceID string, _devID string, _orgID string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncDeleteIgnoreForDevice != nil {
		return myMockSQLDriver.FuncDeleteIgnoreForDevice(_sourceID, _devID, _orgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) DeleteSessionByToken(_SessionKey string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncDeleteSessionByToken != nil {
		return myMockSQLDriver.FuncDeleteSessionByToken(_SessionKey)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) DeleteTagMap(_TicketingSourceID string, _TicketingTag string, _CloudSourceID string, _CloudTag string, _OrganizationID string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncDeleteTagMap != nil {
		return myMockSQLDriver.FuncDeleteTagMap(_TicketingSourceID, _TicketingTag, _CloudSourceID, _CloudTag, _OrganizationID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) DeleteUserByUsername(_Username string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncDeleteUserByUsername != nil {
		return myMockSQLDriver.FuncDeleteUserByUsername(_Username)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) DisableIgnore(inSourceID string, inDevID string, inOrgID string, inVulnID string, inPortID string, inUpdatedBy string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncDisableIgnore != nil {
		return myMockSQLDriver.FuncDisableIgnore(inSourceID, inDevID, inOrgID, inVulnID, inPortID, inUpdatedBy)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) DisableJobConfig(_ID string, _UpdatedBy string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncDisableJobConfig != nil {
		return myMockSQLDriver.FuncDisableJobConfig(_ID, _UpdatedBy)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) DisableOrganization(_ID string, _UpdatedBy string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncDisableOrganization != nil {
		return myMockSQLDriver.FuncDisableOrganization(_ID, _UpdatedBy)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) DisableSource(_ID string, _OrgID string, _UpdatedBy string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncDisableSource != nil {
		return myMockSQLDriver.FuncDisableSource(_ID, _OrgID, _UpdatedBy)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetAllDetectionInfo(_OrgID string) ([]domain.DetectionInfo, error) {
	if myMockSQLDriver.FuncGetAllDetectionInfo != nil {
		return myMockSQLDriver.FuncGetAllDetectionInfo(_OrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetAllDeviceInfo() ([]domain.DeviceInfo, error) {
	if myMockSQLDriver.FuncGetAllDeviceInfo != nil {
		return myMockSQLDriver.FuncGetAllDeviceInfo()
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetAllExceptions(_offset int, _limit int, _sourceID string, _orgID string, _typeID int, _vulnID string, _devID string, _dueDate time.Time, _port string, _approval string, _active bool, _dBCreatedDate time.Time, _dBUpdatedDate time.Time, _updatedBy string, _createdBy string, _sortField string, _sortOrder string) ([]domain.Ignore, error) {
	if myMockSQLDriver.FuncGetAllExceptions != nil {
		return myMockSQLDriver.FuncGetAllExceptions(_offset, _limit, _sourceID, _orgID, _typeID, _vulnID, _devID, _dueDate, _port, _approval, _active, _dBCreatedDate, _dBUpdatedDate, _updatedBy, _createdBy, _sortField, _sortOrder)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetAllJobConfigs(_OrgID string) ([]domain.JobConfig, error) {
	if myMockSQLDriver.FuncGetAllJobConfigs != nil {
		return myMockSQLDriver.FuncGetAllJobConfigs(_OrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetAllJobConfigsWithOrder(_offset int, _limit int, _configID string, _jobid int, _dataInSourceConfigID string, _dataOutSourceConfigID string, _priorityOverride int, _continuous bool, _Payload string, _waitInSeconds int, _maxInstances int, _autoStart bool, _OrgID string, _updatedBy string, _createdBy string, _sortField string, _sortOrder string, _updatedDate time.Time, _createdDate time.Time, _lastJobStart time.Time, _ID string) ([]domain.JobConfig, error) {
	if myMockSQLDriver.FuncGetAllJobConfigsWithOrder != nil {
		return myMockSQLDriver.FuncGetAllJobConfigsWithOrder(_offset, _limit, _configID, _jobid, _dataInSourceConfigID, _dataOutSourceConfigID, _priorityOverride, _continuous, _Payload, _waitInSeconds, _maxInstances, _autoStart, _OrgID, _updatedBy, _createdBy, _sortField, _sortOrder, _updatedDate, _createdDate, _lastJobStart, _ID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetAssetGroup(inOrgID string, _GroupID string, _ScannerConfigSourceID string) (domain.AssetGroup, error) {
	if myMockSQLDriver.FuncGetAssetGroup != nil {
		return myMockSQLDriver.FuncGetAssetGroup(inOrgID, _GroupID, _ScannerConfigSourceID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetAssetGroupForOrg(inScannerSourceConfigID string, inOrgID string) ([]domain.AssetGroup, error) {
	if myMockSQLDriver.FuncGetAssetGroupForOrg != nil {
		return myMockSQLDriver.FuncGetAssetGroupForOrg(inScannerSourceConfigID, inOrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetAssetGroupForOrgNoScanner(inOrgID string, inGroupID string) (domain.AssetGroup, error) {
	if myMockSQLDriver.FuncGetAssetGroupForOrgNoScanner != nil {
		return myMockSQLDriver.FuncGetAssetGroupForOrgNoScanner(inOrgID, inGroupID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetAssetGroupsByCloudSource(inOrgID string, inCloudSourceID string) ([]domain.AssetGroup, error) {
	if myMockSQLDriver.FuncGetAssetGroupsByCloudSource != nil {
		return myMockSQLDriver.FuncGetAssetGroupsByCloudSource(inOrgID, inCloudSourceID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetAssetGroupsForOrg(inOrgID string) ([]domain.AssetGroup, error) {
	if myMockSQLDriver.FuncGetAssetGroupsForOrg != nil {
		return myMockSQLDriver.FuncGetAssetGroupsForOrg(inOrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetAssignmentGroupByIP(_SourceID string, _OrganizationID string, _IP string) ([]domain.AssignmentGroup, error) {
	if myMockSQLDriver.FuncGetAssignmentGroupByIP != nil {
		return myMockSQLDriver.FuncGetAssignmentGroupByIP(_SourceID, _OrganizationID, _IP)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetAssignmentGroupByOrgIP(_OrganizationID string, _IP string) ([]domain.AssignmentGroup, error) {
	if myMockSQLDriver.FuncGetAssignmentGroupByOrgIP != nil {
		return myMockSQLDriver.FuncGetAssignmentGroupByOrgIP(_OrganizationID, _IP)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetAssignmentRulesByOrg(_OrganizationID string) ([]domain.AssignmentRules, error) {
	if myMockSQLDriver.FuncGetAssignmentRulesByOrg != nil {
		return myMockSQLDriver.FuncGetAssignmentRulesByOrg(_OrganizationID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetAutoStartJobs() ([]domain.JobConfig, error) {
	if myMockSQLDriver.FuncGetAutoStartJobs != nil {
		return myMockSQLDriver.FuncGetAutoStartJobs()
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetCISAssignments(_OrganizationID string) ([]domain.CISAssignments, error) {
	if myMockSQLDriver.FuncGetCISAssignments != nil {
		return myMockSQLDriver.FuncGetCISAssignments(_OrganizationID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetCancelledJobs() ([]domain.JobHistory, error) {
	if myMockSQLDriver.FuncGetCancelledJobs != nil {
		return myMockSQLDriver.FuncGetCancelledJobs()
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetCategoryByName(_Name string) ([]domain.Category, error) {
	if myMockSQLDriver.FuncGetCategoryByName != nil {
		return myMockSQLDriver.FuncGetCategoryByName(_Name)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetCategoryRules(_OrgID string, _SourceID string) ([]domain.CategoryRule, error) {
	if myMockSQLDriver.FuncGetCategoryRules != nil {
		return myMockSQLDriver.FuncGetCategoryRules(_OrgID, _SourceID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetDetectionInfo(_DeviceID string, _VulnerabilityID string, _Port int, _Protocol string) (domain.DetectionInfo, error) {
	if myMockSQLDriver.FuncGetDetectionInfo != nil {
		return myMockSQLDriver.FuncGetDetectionInfo(_DeviceID, _VulnerabilityID, _Port, _Protocol)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetDetectionInfoAfter(_After time.Time, _OrgID string) ([]domain.DetectionInfo, error) {
	if myMockSQLDriver.FuncGetDetectionInfoAfter != nil {
		return myMockSQLDriver.FuncGetDetectionInfoAfter(_After, _OrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetDetectionInfoByID(_ID string, _OrgID string) (domain.DetectionInfo, error) {
	if myMockSQLDriver.FuncGetDetectionInfoByID != nil {
		return myMockSQLDriver.FuncGetDetectionInfoByID(_ID, _OrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetDetectionInfoBySourceVulnID(_SourceDeviceID string, _SourceVulnerabilityID string, _Port int, _Protocol string) (domain.DetectionInfo, error) {
	if myMockSQLDriver.FuncGetDetectionInfoBySourceVulnID != nil {
		return myMockSQLDriver.FuncGetDetectionInfoBySourceVulnID(_SourceDeviceID, _SourceVulnerabilityID, _Port, _Protocol)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetDetectionInfoForDeviceID(inDeviceID string, _OrgID string, ticketInactiveKernels bool) ([]domain.DetectionInfo, error) {
	if myMockSQLDriver.FuncGetDetectionInfoForDeviceID != nil {
		return myMockSQLDriver.FuncGetDetectionInfoForDeviceID(inDeviceID, _OrgID, ticketInactiveKernels)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetDetectionInfoForGroupAfter(_LastUpdatedAfter time.Time, _LastFoundAfter time.Time, _OrgID string, inGroupID string, ticketInactiveKernels bool) ([]domain.DetectionInfo, error) {
	if myMockSQLDriver.FuncGetDetectionInfoForGroupAfter != nil {
		return myMockSQLDriver.FuncGetDetectionInfoForGroupAfter(_LastUpdatedAfter, _LastFoundAfter, _OrgID, inGroupID, ticketInactiveKernels)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetDetectionStatusByID(_ID int) (domain.DetectionStatus, error) {
	if myMockSQLDriver.FuncGetDetectionStatusByID != nil {
		return myMockSQLDriver.FuncGetDetectionStatusByID(_ID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetDetectionStatusByName(_Name string) (domain.DetectionStatus, error) {
	if myMockSQLDriver.FuncGetDetectionStatusByName != nil {
		return myMockSQLDriver.FuncGetDetectionStatusByName(_Name)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetDetectionStatuses() ([]domain.DetectionStatus, error) {
	if myMockSQLDriver.FuncGetDetectionStatuses != nil {
		return myMockSQLDriver.FuncGetDetectionStatuses()
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetDetectionsInfoForDevice(_DeviceID string) ([]domain.DetectionInfo, error) {
	if myMockSQLDriver.FuncGetDetectionsInfoForDevice != nil {
		return myMockSQLDriver.FuncGetDetectionsInfoForDevice(_DeviceID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetDeviceInfoByAssetIDNoOrg(inAssetID string) (domain.DeviceInfo, error) {
	if myMockSQLDriver.FuncGetDeviceInfoByAssetIDNoOrg != nil {
		return myMockSQLDriver.FuncGetDeviceInfoByAssetIDNoOrg(inAssetID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetDeviceInfoByAssetOrgID(inAssetID string, inOrgID string) (domain.DeviceInfo, error) {
	if myMockSQLDriver.FuncGetDeviceInfoByAssetOrgID != nil {
		return myMockSQLDriver.FuncGetDeviceInfoByAssetOrgID(inAssetID, inOrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetDeviceInfoByCloudSourceIDAndIP(_IP string, _CloudSourceID string, _OrgID string) ([]domain.DeviceInfo, error) {
	if myMockSQLDriver.FuncGetDeviceInfoByCloudSourceIDAndIP != nil {
		return myMockSQLDriver.FuncGetDeviceInfoByCloudSourceIDAndIP(_IP, _CloudSourceID, _OrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetDeviceInfoByGroupIP(inIP string, inGroupID string, inOrgID string) (domain.DeviceInfo, error) {
	if myMockSQLDriver.FuncGetDeviceInfoByGroupIP != nil {
		return myMockSQLDriver.FuncGetDeviceInfoByGroupIP(inIP, inGroupID, inOrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetDeviceInfoByIP(_IP string, _OrgID string) (domain.DeviceInfo, error) {
	if myMockSQLDriver.FuncGetDeviceInfoByIP != nil {
		return myMockSQLDriver.FuncGetDeviceInfoByIP(_IP, _OrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetDeviceInfoByIPMACAndRegion(_IP string, _MAC string, _Region string, _OrgID string) (domain.DeviceInfo, error) {
	if myMockSQLDriver.FuncGetDeviceInfoByIPMACAndRegion != nil {
		return myMockSQLDriver.FuncGetDeviceInfoByIPMACAndRegion(_IP, _MAC, _Region, _OrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetDeviceInfoByInstanceID(_InstanceID string, _OrgID string) ([]domain.DeviceInfo, error) {
	if myMockSQLDriver.FuncGetDeviceInfoByInstanceID != nil {
		return myMockSQLDriver.FuncGetDeviceInfoByInstanceID(_InstanceID, _OrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetDeviceInfoByScannerSourceID(_IP string, _GroupID string, _OrgID string) (domain.DeviceInfo, error) {
	if myMockSQLDriver.FuncGetDeviceInfoByScannerSourceID != nil {
		return myMockSQLDriver.FuncGetDeviceInfoByScannerSourceID(_IP, _GroupID, _OrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetDevicesInfoByCloudSourceID(_CloudSourceID string, _OrgID string) ([]domain.DeviceInfo, error) {
	if myMockSQLDriver.FuncGetDevicesInfoByCloudSourceID != nil {
		return myMockSQLDriver.FuncGetDevicesInfoByCloudSourceID(_CloudSourceID, _OrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetDevicesInfoBySourceID(_SourceID string, _OrgID string) ([]domain.DeviceInfo, error) {
	if myMockSQLDriver.FuncGetDevicesInfoBySourceID != nil {
		return myMockSQLDriver.FuncGetDevicesInfoBySourceID(_SourceID, _OrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetExceptionByVulnIDOrg(_DeviceID string, _VulnID string, _OrgID string) (domain.Ignore, error) {
	if myMockSQLDriver.FuncGetExceptionByVulnIDOrg != nil {
		return myMockSQLDriver.FuncGetExceptionByVulnIDOrg(_DeviceID, _VulnID, _OrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetExceptionDetections(_offset int, _limit int, _orgID string, _sortField string, _sortOrder string, _Title string, _IP string, _Hostname string, _VulnID string, _Approval string, _DueDate string, _AssignmentGroup string, _OS string, _OSRegex string, _TypeID int) ([]domain.ExceptedDetection, error) {
	if myMockSQLDriver.FuncGetExceptionDetections != nil {
		return myMockSQLDriver.FuncGetExceptionDetections(_offset, _limit, _orgID, _sortField, _sortOrder, _Title, _IP, _Hostname, _VulnID, _Approval, _DueDate, _AssignmentGroup, _OS, _OSRegex, _TypeID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetExceptionTypes() ([]domain.ExceptionType, error) {
	if myMockSQLDriver.FuncGetExceptionTypes != nil {
		return myMockSQLDriver.FuncGetExceptionTypes()
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetExceptionsByOrg(_OrgID string) ([]domain.Ignore, error) {
	if myMockSQLDriver.FuncGetExceptionsByOrg != nil {
		return myMockSQLDriver.FuncGetExceptionsByOrg(_OrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetExceptionsDueNext30Days() ([]domain.CERF, error) {
	if myMockSQLDriver.FuncGetExceptionsDueNext30Days != nil {
		return myMockSQLDriver.FuncGetExceptionsDueNext30Days()
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetExceptionsLength(_offset int, _limit int, _orgID string, _sortField string, _sortOrder string, _Title string, _IP string, _Hostname string, _VulnID string, _Approval string, _DueDate string, _AssignmentGroup string, _OS string, _OSRegex string, _TypeID int) (domain.QueryData, error) {
	if myMockSQLDriver.FuncGetExceptionsLength != nil {
		return myMockSQLDriver.FuncGetExceptionsLength(_offset, _limit, _orgID, _sortField, _sortOrder, _Title, _IP, _Hostname, _VulnID, _Approval, _DueDate, _AssignmentGroup, _OS, _OSRegex, _TypeID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetGlobalExceptions(_OrgID string) ([]domain.Ignore, error) {
	if myMockSQLDriver.FuncGetGlobalExceptions != nil {
		return myMockSQLDriver.FuncGetGlobalExceptions(_OrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetJobByID(_ID int) (domain.JobRegistration, error) {
	if myMockSQLDriver.FuncGetJobByID != nil {
		return myMockSQLDriver.FuncGetJobByID(_ID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetJobConfig(_ID string) (domain.JobConfig, error) {
	if myMockSQLDriver.FuncGetJobConfig != nil {
		return myMockSQLDriver.FuncGetJobConfig(_ID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetJobConfigAudit(inJobConfigID string, inOrgID string) ([]domain.JobConfigAudit, error) {
	if myMockSQLDriver.FuncGetJobConfigAudit != nil {
		return myMockSQLDriver.FuncGetJobConfigAudit(inJobConfigID, inOrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetJobConfigByID(_ID string, _OrgID string) (domain.JobConfig, error) {
	if myMockSQLDriver.FuncGetJobConfigByID != nil {
		return myMockSQLDriver.FuncGetJobConfigByID(_ID, _OrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetJobConfigByJobHistoryID(_JobHistoryID string) (domain.JobConfig, error) {
	if myMockSQLDriver.FuncGetJobConfigByJobHistoryID != nil {
		return myMockSQLDriver.FuncGetJobConfigByJobHistoryID(_JobHistoryID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetJobConfigByOrgIDAndJobID(_OrgID string, _JobID int) ([]domain.JobConfig, error) {
	if myMockSQLDriver.FuncGetJobConfigByOrgIDAndJobID != nil {
		return myMockSQLDriver.FuncGetJobConfigByOrgIDAndJobID(_OrgID, _JobID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetJobConfigByOrgIDAndJobIDWithSC(_OrgID string, _JobID int, _SourceConfigID string) ([]domain.JobConfig, error) {
	if myMockSQLDriver.FuncGetJobConfigByOrgIDAndJobIDWithSC != nil {
		return myMockSQLDriver.FuncGetJobConfigByOrgIDAndJobIDWithSC(_OrgID, _JobID, _SourceConfigID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetJobConfigLength(_configID string, _jobID int, _dataInSourceConfigID string, _dataOutSourceConfigID string, _priorityOverride int, _continuous bool, _Payload string, _waitInSeconds int, _maxInstances int, _autoStart bool, _OrgID string, _updatedBy string, _createdBy string, _updatedDate time.Time, _createdDate time.Time, _lastJobStart time.Time, _ID string) (domain.QueryData, error) {
	if myMockSQLDriver.FuncGetJobConfigLength != nil {
		return myMockSQLDriver.FuncGetJobConfigLength(_configID, _jobID, _dataInSourceConfigID, _dataOutSourceConfigID, _priorityOverride, _continuous, _Payload, _waitInSeconds, _maxInstances, _autoStart, _OrgID, _updatedBy, _createdBy, _updatedDate, _createdDate, _lastJobStart, _ID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetJobHistories(_offset int, _limit int, _jobID int, _jobconfig string, _status int, _Payload string, _OrgID string) ([]domain.JobHistory, error) {
	if myMockSQLDriver.FuncGetJobHistories != nil {
		return myMockSQLDriver.FuncGetJobHistories(_offset, _limit, _jobID, _jobconfig, _status, _Payload, _OrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetJobHistoryByID(_ID string) (domain.JobHistory, error) {
	if myMockSQLDriver.FuncGetJobHistoryByID != nil {
		return myMockSQLDriver.FuncGetJobHistoryByID(_ID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetJobHistoryLength(_jobid int, _jobconfig string, _status int, _Payload string, _orgid string) (domain.QueryData, error) {
	if myMockSQLDriver.FuncGetJobHistoryLength != nil {
		return myMockSQLDriver.FuncGetJobHistoryLength(_jobid, _jobconfig, _status, _Payload, _orgid)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetJobQueueByStatusID(_StatusID int) ([]domain.JobHistory, error) {
	if myMockSQLDriver.FuncGetJobQueueByStatusID != nil {
		return myMockSQLDriver.FuncGetJobQueueByStatusID(_StatusID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetJobs() ([]domain.JobRegistration, error) {
	if myMockSQLDriver.FuncGetJobs != nil {
		return myMockSQLDriver.FuncGetJobs()
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetJobsByStruct(_Struct string) (domain.JobRegistration, error) {
	if myMockSQLDriver.FuncGetJobsByStruct != nil {
		return myMockSQLDriver.FuncGetJobsByStruct(_Struct)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetLeafOrganizationsForUser(_UserID string) ([]domain.Organization, error) {
	if myMockSQLDriver.FuncGetLeafOrganizationsForUser != nil {
		return myMockSQLDriver.FuncGetLeafOrganizationsForUser(_UserID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetLogTypes() ([]domain.LogType, error) {
	if myMockSQLDriver.FuncGetLogTypes != nil {
		return myMockSQLDriver.FuncGetLogTypes()
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetLogsByParams(_MethodOfDiscovery string, _jobType int, _logType int, _jobHistoryID string, _fromDate time.Time, _toDate time.Time, _OrgID string) ([]domain.DBLog, error) {
	if myMockSQLDriver.FuncGetLogsByParams != nil {
		return myMockSQLDriver.FuncGetLogsByParams(_MethodOfDiscovery, _jobType, _logType, _jobHistoryID, _fromDate, _toDate, _OrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetMatchedVulns() ([]domain.VulnerabilityMatch, error) {
	if myMockSQLDriver.FuncGetMatchedVulns != nil {
		return myMockSQLDriver.FuncGetMatchedVulns()
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetOperatingSystemType(_OS string) (domain.OperatingSystemType, error) {
	if myMockSQLDriver.FuncGetOperatingSystemType != nil {
		return myMockSQLDriver.FuncGetOperatingSystemType(_OS)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetOrganizationByCode(Code string) (domain.Organization, error) {
	if myMockSQLDriver.FuncGetOrganizationByCode != nil {
		return myMockSQLDriver.FuncGetOrganizationByCode(Code)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetOrganizationByID(ID string) (domain.Organization, error) {
	if myMockSQLDriver.FuncGetOrganizationByID != nil {
		return myMockSQLDriver.FuncGetOrganizationByID(ID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetOrganizations() ([]domain.Organization, error) {
	if myMockSQLDriver.FuncGetOrganizations != nil {
		return myMockSQLDriver.FuncGetOrganizations()
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetPendingActiveCloudDecomJob(_OrgID string) ([]domain.JobHistory, error) {
	if myMockSQLDriver.FuncGetPendingActiveCloudDecomJob != nil {
		return myMockSQLDriver.FuncGetPendingActiveCloudDecomJob(_OrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetPendingActiveRescanJob(_OrgID string) ([]domain.JobHistory, error) {
	if myMockSQLDriver.FuncGetPendingActiveRescanJob != nil {
		return myMockSQLDriver.FuncGetPendingActiveRescanJob(_OrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetPermissionByUserOrgID(_UserID string, _OrgID string) (domain.Permission, error) {
	if myMockSQLDriver.FuncGetPermissionByUserOrgID != nil {
		return myMockSQLDriver.FuncGetPermissionByUserOrgID(_UserID, _OrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetPermissionOfLeafOrgByUserID(_UserID string) (domain.Permission, error) {
	if myMockSQLDriver.FuncGetPermissionOfLeafOrgByUserID != nil {
		return myMockSQLDriver.FuncGetPermissionOfLeafOrgByUserID(_UserID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetRecentlyUpdatedScanSummaries(_OrgID string) ([]domain.ScanSummary, error) {
	if myMockSQLDriver.FuncGetRecentlyUpdatedScanSummaries != nil {
		return myMockSQLDriver.FuncGetRecentlyUpdatedScanSummaries(_OrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetScanSummariesBySourceName(_OrgID string, _SourceName string) ([]domain.ScanSummary, error) {
	if myMockSQLDriver.FuncGetScanSummariesBySourceName != nil {
		return myMockSQLDriver.FuncGetScanSummariesBySourceName(_OrgID, _SourceName)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetScanSummary(_SourceID string, _OrgID string, _ScanID string) (domain.ScanSummary, error) {
	if myMockSQLDriver.FuncGetScanSummary != nil {
		return myMockSQLDriver.FuncGetScanSummary(_SourceID, _OrgID, _ScanID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetScanSummaryBySourceKey(_SourceKey string) (domain.ScanSummary, error) {
	if myMockSQLDriver.FuncGetScanSummaryBySourceKey != nil {
		return myMockSQLDriver.FuncGetScanSummaryBySourceKey(_SourceKey)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetScheduledJobsToStart(_LastChecked time.Time) ([]domain.JobSchedule, error) {
	if myMockSQLDriver.FuncGetScheduledJobsToStart != nil {
		return myMockSQLDriver.FuncGetScheduledJobsToStart(_LastChecked)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetSessionByToken(_SessionKey string) (domain.Session, error) {
	if myMockSQLDriver.FuncGetSessionByToken != nil {
		return myMockSQLDriver.FuncGetSessionByToken(_SessionKey)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetSourceByID(_ID string) (domain.Source, error) {
	if myMockSQLDriver.FuncGetSourceByID != nil {
		return myMockSQLDriver.FuncGetSourceByID(_ID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetSourceByName(_Source string) (domain.Source, error) {
	if myMockSQLDriver.FuncGetSourceByName != nil {
		return myMockSQLDriver.FuncGetSourceByName(_Source)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetSourceConfigByID(_ID string) (domain.SourceConfig, error) {
	if myMockSQLDriver.FuncGetSourceConfigByID != nil {
		return myMockSQLDriver.FuncGetSourceConfigByID(_ID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetSourceConfigByNameOrg(_Source string, _OrgID string) ([]domain.SourceConfig, error) {
	if myMockSQLDriver.FuncGetSourceConfigByNameOrg != nil {
		return myMockSQLDriver.FuncGetSourceConfigByNameOrg(_Source, _OrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetSourceConfigByOrgID(_OrgID string) ([]domain.SourceConfig, error) {
	if myMockSQLDriver.FuncGetSourceConfigByOrgID != nil {
		return myMockSQLDriver.FuncGetSourceConfigByOrgID(_OrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetSourceConfigBySourceID(_OrgID string, _SourceID string) ([]domain.SourceConfig, error) {
	if myMockSQLDriver.FuncGetSourceConfigBySourceID != nil {
		return myMockSQLDriver.FuncGetSourceConfigBySourceID(_OrgID, _SourceID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetSourceInsByJobID(inJob int, inOrgID string) ([]domain.SourceConfig, error) {
	if myMockSQLDriver.FuncGetSourceInsByJobID != nil {
		return myMockSQLDriver.FuncGetSourceInsByJobID(inJob, inOrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetSourceOauthByOrgURL(_URL string, _OrgID string) (domain.SourceConfig, error) {
	if myMockSQLDriver.FuncGetSourceOauthByOrgURL != nil {
		return myMockSQLDriver.FuncGetSourceOauthByOrgURL(_URL, _OrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetSourceOauthByURL(_URL string) (domain.SourceConfig, error) {
	if myMockSQLDriver.FuncGetSourceOauthByURL != nil {
		return myMockSQLDriver.FuncGetSourceOauthByURL(_URL)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetSourceOutsByJobID(inJob int, inOrgID string) ([]domain.SourceConfig, error) {
	if myMockSQLDriver.FuncGetSourceOutsByJobID != nil {
		return myMockSQLDriver.FuncGetSourceOutsByJobID(inJob, inOrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetSources() ([]domain.Source, error) {
	if myMockSQLDriver.FuncGetSources != nil {
		return myMockSQLDriver.FuncGetSources()
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetTagByDeviceAndTagKey(_DeviceID string, _TagKeyID string) (domain.Tag, error) {
	if myMockSQLDriver.FuncGetTagByDeviceAndTagKey != nil {
		return myMockSQLDriver.FuncGetTagByDeviceAndTagKey(_DeviceID, _TagKeyID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetTagKeyByID(_ID string) (domain.TagKey, error) {
	if myMockSQLDriver.FuncGetTagKeyByID != nil {
		return myMockSQLDriver.FuncGetTagKeyByID(_ID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetTagKeyByKey(_KeyValue string) (domain.TagKey, error) {
	if myMockSQLDriver.FuncGetTagKeyByKey != nil {
		return myMockSQLDriver.FuncGetTagKeyByKey(_KeyValue)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetTagMapsByOrg(_OrganizationID string) ([]domain.TagMap, error) {
	if myMockSQLDriver.FuncGetTagMapsByOrg != nil {
		return myMockSQLDriver.FuncGetTagMapsByOrg(_OrganizationID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetTagMapsByOrgCloudSourceID(_CloudID string, _OrganizationID string) ([]domain.TagMap, error) {
	if myMockSQLDriver.FuncGetTagMapsByOrgCloudSourceID != nil {
		return myMockSQLDriver.FuncGetTagMapsByOrgCloudSourceID(_CloudID, _OrganizationID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetTagsForDevice(_DeviceID string) ([]domain.Tag, error) {
	if myMockSQLDriver.FuncGetTagsForDevice != nil {
		return myMockSQLDriver.FuncGetTagsForDevice(_DeviceID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetTicketByDetectionID(inDetectionID string, _OrgID string) (domain.TicketSummary, error) {
	if myMockSQLDriver.FuncGetTicketByDetectionID != nil {
		return myMockSQLDriver.FuncGetTicketByDetectionID(inDetectionID, _OrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetTicketByDeviceIDVulnID(inDeviceID string, inVulnID string, inPort int, inProtocol string, inOrgID string) (domain.TicketSummary, error) {
	if myMockSQLDriver.FuncGetTicketByDeviceIDVulnID != nil {
		return myMockSQLDriver.FuncGetTicketByDeviceIDVulnID(inDeviceID, inVulnID, inPort, inProtocol, inOrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetTicketByIPGroupIDVulnID(inIP string, inGroupID string, inVulnID string, inPort int, inProtocol string, inOrgID string) (domain.TicketSummary, error) {
	if myMockSQLDriver.FuncGetTicketByIPGroupIDVulnID != nil {
		return myMockSQLDriver.FuncGetTicketByIPGroupIDVulnID(inIP, inGroupID, inVulnID, inPort, inProtocol, inOrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetTicketByTitle(_Title string, _OrgID string) (domain.TicketSummary, error) {
	if myMockSQLDriver.FuncGetTicketByTitle != nil {
		return myMockSQLDriver.FuncGetTicketByTitle(_Title, _OrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetTicketCountByStatus(inStatus string, inOrgID string) (domain.QueryData, error) {
	if myMockSQLDriver.FuncGetTicketCountByStatus != nil {
		return myMockSQLDriver.FuncGetTicketCountByStatus(inStatus, inOrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetTicketCreatedAfter(_UpperCVSS float32, _LowerCVSS float32, _CreatedAfter time.Time, _OrgID string) ([]domain.TicketSummary, error) {
	if myMockSQLDriver.FuncGetTicketCreatedAfter != nil {
		return myMockSQLDriver.FuncGetTicketCreatedAfter(_UpperCVSS, _LowerCVSS, _CreatedAfter, _OrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetTicketTrackingMethod(_Title string, _OrgID string) (domain.KeyValue, error) {
	if myMockSQLDriver.FuncGetTicketTrackingMethod != nil {
		return myMockSQLDriver.FuncGetTicketTrackingMethod(_Title, _OrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetUnfinishedScanSummariesBySourceConfigOrgID(_ScannerSourceConfigID string, _OrgID string) ([]domain.ScanSummary, error) {
	if myMockSQLDriver.FuncGetUnfinishedScanSummariesBySourceConfigOrgID != nil {
		return myMockSQLDriver.FuncGetUnfinishedScanSummariesBySourceConfigOrgID(_ScannerSourceConfigID, _OrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetUnfinishedScanSummariesBySourceOrgID(_SourceID string, _OrgID string) ([]domain.ScanSummary, error) {
	if myMockSQLDriver.FuncGetUnfinishedScanSummariesBySourceOrgID != nil {
		return myMockSQLDriver.FuncGetUnfinishedScanSummariesBySourceOrgID(_SourceID, _OrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetUnmatchedVulns(_SourceID int) ([]domain.VulnerabilityInfo, error) {
	if myMockSQLDriver.FuncGetUnmatchedVulns != nil {
		return myMockSQLDriver.FuncGetUnmatchedVulns(_SourceID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetUserAnyOrg(_ID string) (domain.User, error) {
	if myMockSQLDriver.FuncGetUserAnyOrg != nil {
		return myMockSQLDriver.FuncGetUserAnyOrg(_ID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetUserByID(_ID string, _OrgID string) (domain.User, error) {
	if myMockSQLDriver.FuncGetUserByID != nil {
		return myMockSQLDriver.FuncGetUserByID(_ID, _OrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetUserByUsername(_Username string) (domain.User, error) {
	if myMockSQLDriver.FuncGetUserByUsername != nil {
		return myMockSQLDriver.FuncGetUserByUsername(_Username)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetUsersByOrg(_OrgID string) ([]domain.User, error) {
	if myMockSQLDriver.FuncGetUsersByOrg != nil {
		return myMockSQLDriver.FuncGetUsersByOrg(_OrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetVulnInfoByID(_ID string) (domain.VulnerabilityInfo, error) {
	if myMockSQLDriver.FuncGetVulnInfoByID != nil {
		return myMockSQLDriver.FuncGetVulnInfoByID(_ID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetVulnInfoBySource(_Source string) ([]domain.VulnerabilityInfo, error) {
	if myMockSQLDriver.FuncGetVulnInfoBySource != nil {
		return myMockSQLDriver.FuncGetVulnInfoBySource(_Source)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetVulnInfoBySourceID(_SourceID string) ([]domain.VulnerabilityInfo, error) {
	if myMockSQLDriver.FuncGetVulnInfoBySourceID != nil {
		return myMockSQLDriver.FuncGetVulnInfoBySourceID(_SourceID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetVulnInfoBySourceVulnID(_SourceVulnID string) (domain.VulnerabilityInfo, error) {
	if myMockSQLDriver.FuncGetVulnInfoBySourceVulnID != nil {
		return myMockSQLDriver.FuncGetVulnInfoBySourceVulnID(_SourceVulnID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetVulnInfoBySourceVulnIDSourceID(_SourceVulnID string, _SourceID string, _Modified time.Time) (domain.VulnerabilityInfo, error) {
	if myMockSQLDriver.FuncGetVulnInfoBySourceVulnIDSourceID != nil {
		return myMockSQLDriver.FuncGetVulnInfoBySourceVulnIDSourceID(_SourceVulnID, _SourceID, _Modified)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetVulnRefInfo(_VulnInfoID string, _SourceID string, _Reference string) (domain.VulnerabilityReferenceInfo, error) {
	if myMockSQLDriver.FuncGetVulnRefInfo != nil {
		return myMockSQLDriver.FuncGetVulnRefInfo(_VulnInfoID, _SourceID, _Reference)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetVulnRefInfoVendor(_VulnInfoID string, _SourceID string) ([]domain.VulnerabilityReferenceInfo, error) {
	if myMockSQLDriver.FuncGetVulnRefInfoVendor != nil {
		return myMockSQLDriver.FuncGetVulnRefInfoVendor(_VulnInfoID, _SourceID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetVulnReferencesInfo(_VulnInfoID string, _SourceID string) ([]domain.VulnerabilityReferenceInfo, error) {
	if myMockSQLDriver.FuncGetVulnReferencesInfo != nil {
		return myMockSQLDriver.FuncGetVulnReferencesInfo(_VulnInfoID, _SourceID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) GetVulnReferencesInfoBySourceAndRef(_SourceID string, _Reference string) ([]domain.VulnerabilityReferenceInfo, error) {
	if myMockSQLDriver.FuncGetVulnReferencesInfoBySourceAndRef != nil {
		return myMockSQLDriver.FuncGetVulnReferencesInfoBySourceAndRef(_SourceID, _Reference)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) HasDecommissioned(_devID string, _sourceID string, _orgID string) (domain.Ignore, error) {
	if myMockSQLDriver.FuncHasDecommissioned != nil {
		return myMockSQLDriver.FuncHasDecommissioned(_devID, _sourceID, _orgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) HasExceptionOrFalsePositive(_sourceID string, _vulnID string, _devID string, _orgID string, _port string, _OS string) ([]domain.Ignore, error) {
	if myMockSQLDriver.FuncHasExceptionOrFalsePositive != nil {
		return myMockSQLDriver.FuncHasExceptionOrFalsePositive(_sourceID, _vulnID, _devID, _orgID, _port, _OS)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) HasIgnore(inSourceID string, inVulnID string, inDevID string, inOrgID string, inPort string, inMostCurrentDetection time.Time) (domain.Ignore, error) {
	if myMockSQLDriver.FuncHasIgnore != nil {
		return myMockSQLDriver.FuncHasIgnore(inSourceID, inVulnID, inDevID, inOrgID, inPort, inMostCurrentDetection)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) PulseJob(_JobHistoryID string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncPulseJob != nil {
		return myMockSQLDriver.FuncPulseJob(_JobHistoryID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) RemoveExpiredIgnoreIDs(_OrgID string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncRemoveExpiredIgnoreIDs != nil {
		return myMockSQLDriver.FuncRemoveExpiredIgnoreIDs(_OrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) SaveAssignmentGroup(_SourceID string, _OrganizationID string, _IpAddress string, _GroupName string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncSaveAssignmentGroup != nil {
		return myMockSQLDriver.FuncSaveAssignmentGroup(_SourceID, _OrganizationID, _IpAddress, _GroupName)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) SaveIgnore(_SourceID string, _OrganizationID string, _TypeID int, _VulnerabilityID string, _DeviceID string, _DueDate time.Time, _Approval string, _Active bool, _port string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncSaveIgnore != nil {
		return myMockSQLDriver.FuncSaveIgnore(_SourceID, _OrganizationID, _TypeID, _VulnerabilityID, _DeviceID, _DueDate, _Approval, _Active, _port)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) SaveScanSummary(_ScanID string, _ScanStatus string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncSaveScanSummary != nil {
		return myMockSQLDriver.FuncSaveScanSummary(_ScanID, _ScanStatus)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) SetScheduleLastRun(_ID string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncSetScheduleLastRun != nil {
		return myMockSQLDriver.FuncSetScheduleLastRun(_ID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) UpdateAssetGroupLastTicket(inGroupID string, inOrgID string, inLastTicketTime time.Time) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncUpdateAssetGroupLastTicket != nil {
		return myMockSQLDriver.FuncUpdateAssetGroupLastTicket(inGroupID, inOrgID, inLastTicketTime)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) UpdateAssetIDOsTypeIDOfDevice(_ID string, _AssetID string, _ScannerSourceID string, _GroupID string, _OS string, _HostName string, _OsTypeID int, inTrackingMethod string, _OrgID string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncUpdateAssetIDOsTypeIDOfDevice != nil {
		return myMockSQLDriver.FuncUpdateAssetIDOsTypeIDOfDevice(_ID, _AssetID, _ScannerSourceID, _GroupID, _OS, _HostName, _OsTypeID, inTrackingMethod, _OrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) UpdateDetection(_ID string, _DeviceID string, _VulnID string, _Port int, _Protocol string, _ExceptionID string, _TimesSeen int, _StatusID int, _LastFound time.Time, _LastUpdated time.Time, _DefaultTime time.Time) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncUpdateDetection != nil {
		return myMockSQLDriver.FuncUpdateDetection(_ID, _DeviceID, _VulnID, _Port, _Protocol, _ExceptionID, _TimesSeen, _StatusID, _LastFound, _LastUpdated, _DefaultTime)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) UpdateDetectionIgnore(_DeviceID string, _VulnID string, _Port int, _Protocol string, _ExceptionID string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncUpdateDetectionIgnore != nil {
		return myMockSQLDriver.FuncUpdateDetectionIgnore(_DeviceID, _VulnID, _Port, _Protocol, _ExceptionID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) UpdateExpirationDateByCERF(_CERForm string, _OrganizationID string, _DueDate time.Time) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncUpdateExpirationDateByCERF != nil {
		return myMockSQLDriver.FuncUpdateExpirationDateByCERF(_CERForm, _OrganizationID, _DueDate)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) UpdateInstanceIDOfDevice(_ID string, _InstanceID string, _CloudSourceID string, _State string, _Region string, _OrgID string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncUpdateInstanceIDOfDevice != nil {
		return myMockSQLDriver.FuncUpdateInstanceIDOfDevice(_ID, _InstanceID, _CloudSourceID, _State, _Region, _OrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) UpdateJobConfig(_ID string, _DataInSourceID string, _DataOutSourceID string, _Autostart bool, _PriorityOverride int, _Continuous bool, _WaitInSeconds int, _MaxInstances int, _UpdatedBy string, _OrgID string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncUpdateJobConfig != nil {
		return myMockSQLDriver.FuncUpdateJobConfig(_ID, _DataInSourceID, _DataOutSourceID, _Autostart, _PriorityOverride, _Continuous, _WaitInSeconds, _MaxInstances, _UpdatedBy, _OrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) UpdateJobConfigLastRun(_ID string, _LastRun time.Time) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncUpdateJobConfigLastRun != nil {
		return myMockSQLDriver.FuncUpdateJobConfigLastRun(_ID, _LastRun)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) UpdateJobHistory(_ID string, _ConfigID string, _Payload string, _UpdatedBy string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncUpdateJobHistory != nil {
		return myMockSQLDriver.FuncUpdateJobHistory(_ID, _ConfigID, _Payload, _UpdatedBy)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) UpdateJobHistoryStatus(_ID string, _Status int) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncUpdateJobHistoryStatus != nil {
		return myMockSQLDriver.FuncUpdateJobHistoryStatus(_ID, _Status)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) UpdateJobHistoryStatusDetailed(_ID string, _Status int, _UpdatedBy string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncUpdateJobHistoryStatusDetailed != nil {
		return myMockSQLDriver.FuncUpdateJobHistoryStatusDetailed(_ID, _Status, _UpdatedBy)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) UpdateOrganization(_ID string, _Description string, _TimezoneOffset float32, _UpdatedBy string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncUpdateOrganization != nil {
		return myMockSQLDriver.FuncUpdateOrganization(_ID, _Description, _TimezoneOffset, _UpdatedBy)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) UpdatePermissionsByUserOrgID(_UserID string, _OrgID string, _Admin bool, _Manager bool, _Reader bool, _Reporter bool) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncUpdatePermissionsByUserOrgID != nil {
		return myMockSQLDriver.FuncUpdatePermissionsByUserOrgID(_UserID, _OrgID, _Admin, _Manager, _Reader, _Reporter)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) UpdateSourceConfig(_ID string, _OrgID string, _Address string, _Username string, _Password string, _PrivateKey string, _ConsumerKey string, _Token string, _Port string, _Payload string, _UpdatedBy string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncUpdateSourceConfig != nil {
		return myMockSQLDriver.FuncUpdateSourceConfig(_ID, _OrgID, _Address, _Username, _Password, _PrivateKey, _ConsumerKey, _Token, _Port, _Payload, _UpdatedBy)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) UpdateSourceConfigConcurrencyByID(_ID string, _Delay int, _Retries int, _Concurrency int) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncUpdateSourceConfigConcurrencyByID != nil {
		return myMockSQLDriver.FuncUpdateSourceConfigConcurrencyByID(_ID, _Delay, _Retries, _Concurrency)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) UpdateSourceConfigToken(_ID string, _Token string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncUpdateSourceConfigToken != nil {
		return myMockSQLDriver.FuncUpdateSourceConfigToken(_ID, _Token)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) UpdateStateOfDevice(_ID string, _State string, _OrgID string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncUpdateStateOfDevice != nil {
		return myMockSQLDriver.FuncUpdateStateOfDevice(_ID, _State, _OrgID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) UpdateTag(_DeviceID string, _TagKeyID string, _Value string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncUpdateTag != nil {
		return myMockSQLDriver.FuncUpdateTag(_DeviceID, _TagKeyID, _Value)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) UpdateTagMap(_TicketingSourceID string, _TicketingTag string, _CloudSourceID string, _CloudTag string, _Options string, _OrganizationID string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncUpdateTagMap != nil {
		return myMockSQLDriver.FuncUpdateTagMap(_TicketingSourceID, _TicketingTag, _CloudSourceID, _CloudTag, _Options, _OrganizationID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) UpdateTicket(_Title string, _Status string, _OrganizationID string, _AssignmentGroup string, _Assignee string, _DueDate time.Time, _CreatedDate time.Time, _UpdatedDate time.Time, _ResolutionDate time.Time, _ExceptionDate time.Time, _DefaultTime time.Time) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncUpdateTicket != nil {
		return myMockSQLDriver.FuncUpdateTicket(_Title, _Status, _OrganizationID, _AssignmentGroup, _Assignee, _DueDate, _CreatedDate, _UpdatedDate, _ResolutionDate, _ExceptionDate, _DefaultTime)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) UpdateTicketDetectionID(_Title string, _DetectionID string, _OrganizationID string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncUpdateTicketDetectionID != nil {
		return myMockSQLDriver.FuncUpdateTicketDetectionID(_Title, _DetectionID, _OrganizationID)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) UpdateUserByID(_ID string, _FirstName string, _LastName string, _Email string, _Disabled bool) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncUpdateUserByID != nil {
		return myMockSQLDriver.FuncUpdateUserByID(_ID, _FirstName, _LastName, _Email, _Disabled)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) UpdateVulnByID(_ID string, _SourceVulnID string, _Title string, _SourceID string, _CVSSScore float32, _CVSS3Score float32, _Description string, _Threat string, _Solution string, _Software string, _Patchable string, _Category string, _DetectionInformation string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncUpdateVulnByID != nil {
		return myMockSQLDriver.FuncUpdateVulnByID(_ID, _SourceVulnID, _Title, _SourceID, _CVSSScore, _CVSS3Score, _Description, _Threat, _Solution, _Software, _Patchable, _Category, _DetectionInformation)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) UpdateVulnByIDNoCVSS3(_ID string, _SourceVulnID string, _Title string, _SourceID string, _CVSSScore float32, _CVSS3Score float32, _Description string, _Threat string, _Solution string, _Software string, _Patchable string, _DetectionInformation string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncUpdateVulnByIDNoCVSS3 != nil {
		return myMockSQLDriver.FuncUpdateVulnByIDNoCVSS3(_ID, _SourceVulnID, _Title, _SourceID, _CVSSScore, _CVSS3Score, _Description, _Threat, _Solution, _Software, _Patchable, _DetectionInformation)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}

func (myMockSQLDriver *MockSQLDriver) UpdateVulnInfoID(_VulnInfoID string, _VulnID string, _MatchConfidence int, _MatchReasons string) (id int, affectedRows int, err error) {
	if myMockSQLDriver.FuncUpdateVulnInfoID != nil {
		return myMockSQLDriver.FuncUpdateVulnInfoID(_VulnInfoID, _VulnID, _MatchConfidence, _MatchReasons)
	} else {
		panic("method not implemented") // mock SQL drivers should only be used in testing
	}
}
