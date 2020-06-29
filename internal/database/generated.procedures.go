package database

import (
	"database/sql"
	"github.com/nortonlifelock/aegis/internal/database/dal"
	"github.com/nortonlifelock/connection"
	"github.com/nortonlifelock/domain"
	"time"
)

//**********************************************************
// GENERATED CODE - DO NOT CHANGE
// This file is generated using scaffolding. Any changes to
// this file will be overwritten on the next build
//**********************************************************

// CleanUp executes the stored procedure CleanUp against the database
func (conn *dbconn) CleanUp() (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "CleanUp",
		Parameters: []interface{}{},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// CreateAssetGroup executes the stored procedure CreateAssetGroup against the database
func (conn *dbconn) CreateAssetGroup(inOrgID string, _GroupID string, _ScannerSourceID string, _ScannerSourceConfigID string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "CreateAssetGroup",
		Parameters: []interface{}{inOrgID, _GroupID, _ScannerSourceID, _ScannerSourceConfigID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// CreateAssetWithIPInstanceID executes the stored procedure CreateAssetWithIPInstanceID against the database
func (conn *dbconn) CreateAssetWithIPInstanceID(_State string, _IP string, _MAC string, _SourceID string, _InstanceID string, _Region string, _OrgID string, _OS string, _OsTypeID int) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "CreateAssetWithIPInstanceID",
		Parameters: []interface{}{_State, _IP, _MAC, _SourceID, _InstanceID, _Region, _OrgID, _OS, _OsTypeID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// CreateCategory executes the stored procedure CreateCategory against the database
func (conn *dbconn) CreateCategory(_Category string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "CreateCategory",
		Parameters: []interface{}{_Category},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// CreateDBLog executes the stored procedure CreateDBLog against the database
func (conn *dbconn) CreateDBLog(_User string, _Command string, _Endpoint string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "CreateDBLog",
		Parameters: []interface{}{_User, _Command, _Endpoint},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// CreateDetection executes the stored procedure CreateDetection against the database
func (conn *dbconn) CreateDetection(_OrgID string, _SourceID string, _DeviceID string, _VulnID string, _IgnoreID string, _AlertDate time.Time, _LastFound time.Time, _LastUpdated time.Time, _Proof string, _Port int, _Protocol string, _ActiveKernel int, _DetectionStatusID int, _TimesSeen int, _DefaultTime time.Time) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "CreateDetection",
		Parameters: []interface{}{_OrgID, _SourceID, _DeviceID, _VulnID, _IgnoreID, _AlertDate, _LastFound, _LastUpdated, _Proof, _Port, _Protocol, _ActiveKernel, _DetectionStatusID, _TimesSeen, _DefaultTime},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// CreateDevice executes the stored procedure CreateDevice against the database
func (conn *dbconn) CreateDevice(_AssetID string, _SourceID string, _Ip string, _Hostname string, inInstanceID string, _MAC string, _GroupID string, _OrgID string, _OS string, _OSTypeID int, inTrackingMethod string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "CreateDevice",
		Parameters: []interface{}{_AssetID, _SourceID, _Ip, _Hostname, inInstanceID, _MAC, _GroupID, _OrgID, _OS, _OSTypeID, inTrackingMethod},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// CreateException executes the stored procedure CreateException against the database
func (conn *dbconn) CreateException(inSourceID string, inOrganizationID string, inTypeID int, inVulnerabilityID string, inDeviceID string, inDueDate time.Time, inApproval string, inActive bool, inPort string, inCreatedBy string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "CreateException",
		Parameters: []interface{}{inSourceID, inOrganizationID, inTypeID, inVulnerabilityID, inDeviceID, inDueDate, inApproval, inActive, inPort, inCreatedBy},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// CreateJobConfig executes the stored procedure CreateJobConfig against the database
func (conn *dbconn) CreateJobConfig(_JobID int, _OrganizationID string, _PriorityOverride int, _Continuous bool, _WaitInSeconds int, _MaxInstances int, _AutoStart bool, _CreatedBy string, _DataInSourceID string, _DataOutSourceID string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "CreateJobConfig",
		Parameters: []interface{}{_JobID, _OrganizationID, _PriorityOverride, _Continuous, _WaitInSeconds, _MaxInstances, _AutoStart, _CreatedBy, _DataInSourceID, _DataOutSourceID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// CreateJobConfigWPayload executes the stored procedure CreateJobConfigWPayload against the database
func (conn *dbconn) CreateJobConfigWPayload(_JobID int, _OrganizationID string, _PriorityOverride int, _Continuous bool, _WaitInSeconds int, _MaxInstances int, _AutoStart bool, _CreatedBy string, _DataInSourceID string, _DataOutSourceID string, _Payload string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "CreateJobConfigWPayload",
		Parameters: []interface{}{_JobID, _OrganizationID, _PriorityOverride, _Continuous, _WaitInSeconds, _MaxInstances, _AutoStart, _CreatedBy, _DataInSourceID, _DataOutSourceID, _Payload},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// CreateJobHistory executes the stored procedure CreateJobHistory against the database
func (conn *dbconn) CreateJobHistory(_JobID int, _ConfigID string, _StatusID int, _Priority int, _Identifier string, _CurrentIteration int, _Payload string, _ThreadID string, _PulseDate time.Time, _CreatedBy string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "CreateJobHistory",
		Parameters: []interface{}{_JobID, _ConfigID, _StatusID, _Priority, _Identifier, _CurrentIteration, _Payload, _ThreadID, _PulseDate, _CreatedBy},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// CreateJobHistoryWithParentID executes the stored procedure CreateJobHistoryWithParentID against the database
func (conn *dbconn) CreateJobHistoryWithParentID(_JobID int, _ConfigID string, _StatusID int, _Priority int, _Identifier string, _CurrentIteration int, _Payload string, _ThreadID string, _PulseDate time.Time, _CreatedBy string, _ParentID string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "CreateJobHistoryWithParentID",
		Parameters: []interface{}{_JobID, _ConfigID, _StatusID, _Priority, _Identifier, _CurrentIteration, _Payload, _ThreadID, _PulseDate, _CreatedBy, _ParentID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// CreateLog executes the stored procedure CreateLog against the database
func (conn *dbconn) CreateLog(_TypeID int, _Log string, _Error string, _JobHistoryID string, _CreateDate time.Time) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "CreateLog",
		Parameters: []interface{}{_TypeID, _Log, _Error, _JobHistoryID, _CreateDate},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// CreateOrganization executes the stored procedure CreateOrganization against the database
func (conn *dbconn) CreateOrganization(_Code string, _Description string, _TimeZoneOffset float32, _UpdatedBy string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "CreateOrganization",
		Parameters: []interface{}{_Code, _Description, _TimeZoneOffset, _UpdatedBy},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// CreateOrganizationWithPayloadEkey executes the stored procedure CreateOrganizationWithPayloadEkey against the database
func (conn *dbconn) CreateOrganizationWithPayloadEkey(_Code string, _Description string, _TimeZoneOffset float32, _Payload string, _EKEY string, _UpdatedBy string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "CreateOrganizationWithPayloadEkey",
		Parameters: []interface{}{_Code, _Description, _TimeZoneOffset, _Payload, _EKEY, _UpdatedBy},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// CreateScanSummary executes the stored procedure CreateScanSummary against the database
func (conn *dbconn) CreateScanSummary(_SourceID string, _ScannerSourceConfigID string, _OrgID string, _ScanID string, _ScanStatus string, _ScanClosePayload string, _ParentJobID string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "CreateScanSummary",
		Parameters: []interface{}{_SourceID, _ScannerSourceConfigID, _OrgID, _ScanID, _ScanStatus, _ScanClosePayload, _ParentJobID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// CreateSourceConfig executes the stored procedure CreateSourceConfig against the database
func (conn *dbconn) CreateSourceConfig(_Source string, _SourceID string, _OrganizationID string, _Address string, _Port string, _Username string, _Password string, _PrivateKey string, _ConsumerKey string, _Token string, _Payload string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "CreateSourceConfig",
		Parameters: []interface{}{_Source, _SourceID, _OrganizationID, _Address, _Port, _Username, _Password, _PrivateKey, _ConsumerKey, _Token, _Payload},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// CreateTag executes the stored procedure CreateTag against the database
func (conn *dbconn) CreateTag(_DeviceID string, _TagKeyID string, _Value string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "CreateTag",
		Parameters: []interface{}{_DeviceID, _TagKeyID, _Value},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// CreateTagKey executes the stored procedure CreateTagKey against the database
func (conn *dbconn) CreateTagKey(_KeyValue string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "CreateTagKey",
		Parameters: []interface{}{_KeyValue},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// CreateTagMap executes the stored procedure CreateTagMap against the database
func (conn *dbconn) CreateTagMap(_TicketingSourceID string, _TicketingTag string, _CloudSourceID string, _CloudTag string, _Options string, _OrganizationID string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "CreateTagMap",
		Parameters: []interface{}{_TicketingSourceID, _TicketingTag, _CloudSourceID, _CloudTag, _Options, _OrganizationID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// CreateTicket executes the stored procedure CreateTicket against the database
func (conn *dbconn) CreateTicket(_Title string, _Status string, _DetectionID string, _OrganizationID string, _DueDate time.Time, _UpdatedDate time.Time, _ResolutionDate time.Time, _DefaultTime time.Time) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "CreateTicket",
		Parameters: []interface{}{_Title, _Status, _DetectionID, _OrganizationID, _DueDate, _UpdatedDate, _ResolutionDate, _DefaultTime},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// CreateTicketingJob executes the stored procedure CreateTicketingJob against the database
func (conn *dbconn) CreateTicketingJob(GroupID int, OrgID string, ScanStartDate string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "CreateTicketingJob",
		Parameters: []interface{}{GroupID, OrgID, ScanStartDate},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// CreateUser executes the stored procedure CreateUser against the database
func (conn *dbconn) CreateUser(_Username string, _FirstName string, _LastName string, _Email string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "CreateUser",
		Parameters: []interface{}{_Username, _FirstName, _LastName, _Email},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// CreateUserPermissions executes the stored procedure CreateUserPermissions against the database
func (conn *dbconn) CreateUserPermissions(_UserID string, _OrgID string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "CreateUserPermissions",
		Parameters: []interface{}{_UserID, _OrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// CreateUserSession executes the stored procedure CreateUserSession against the database
func (conn *dbconn) CreateUserSession(_UserID string, _OrgID string, _SessionKey string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "CreateUserSession",
		Parameters: []interface{}{_UserID, _OrgID, _SessionKey},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// CreateVulnInfo executes the stored procedure CreateVulnInfo against the database
func (conn *dbconn) CreateVulnInfo(_SourceVulnID string, _Title string, _SourceID string, _CVSSScore float32, _CVSS3Score float32, _Description string, _Threat string, _Solution string, _Software string, _Patchable string, _Category string, _DetectionInformation string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "CreateVulnInfo",
		Parameters: []interface{}{_SourceVulnID, _Title, _SourceID, _CVSSScore, _CVSS3Score, _Description, _Threat, _Solution, _Software, _Patchable, _Category, _DetectionInformation},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// CreateVulnRef executes the stored procedure CreateVulnRef against the database
func (conn *dbconn) CreateVulnRef(_VulnInfoID string, _SourceID string, _Reference string, _RefType int) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "CreateVulnRef",
		Parameters: []interface{}{_VulnInfoID, _SourceID, _Reference, _RefType},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// DeleteDecomIgnoreForDevice executes the stored procedure DeleteDecomIgnoreForDevice against the database
func (conn *dbconn) DeleteDecomIgnoreForDevice(_sourceID string, _devID string, _orgID string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "DeleteDecomIgnoreForDevice",
		Parameters: []interface{}{_sourceID, _devID, _orgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// DeleteIgnoreForDevice executes the stored procedure DeleteIgnoreForDevice against the database
func (conn *dbconn) DeleteIgnoreForDevice(_sourceID string, _devID string, _orgID string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "DeleteIgnoreForDevice",
		Parameters: []interface{}{_sourceID, _devID, _orgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// DeleteSessionByToken executes the stored procedure DeleteSessionByToken against the database
func (conn *dbconn) DeleteSessionByToken(_SessionKey string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "DeleteSessionByToken",
		Parameters: []interface{}{_SessionKey},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// DeleteTagMap executes the stored procedure DeleteTagMap against the database
func (conn *dbconn) DeleteTagMap(_TicketingSourceID string, _TicketingTag string, _CloudSourceID string, _CloudTag string, _OrganizationID string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "DeleteTagMap",
		Parameters: []interface{}{_TicketingSourceID, _TicketingTag, _CloudSourceID, _CloudTag, _OrganizationID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// DeleteUserByUsername executes the stored procedure DeleteUserByUsername against the database
func (conn *dbconn) DeleteUserByUsername(_Username string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "DeleteUserByUsername",
		Parameters: []interface{}{_Username},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// DisableIgnore executes the stored procedure DisableIgnore against the database
func (conn *dbconn) DisableIgnore(inSourceID string, inDevID string, inOrgID string, inVulnID string, inPortID string, inUpdatedBy string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "DisableIgnore",
		Parameters: []interface{}{inSourceID, inDevID, inOrgID, inVulnID, inPortID, inUpdatedBy},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// DisableJobConfig executes the stored procedure DisableJobConfig against the database
func (conn *dbconn) DisableJobConfig(_ID string, _UpdatedBy string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "DisableJobConfig",
		Parameters: []interface{}{_ID, _UpdatedBy},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// DisableOrganization executes the stored procedure DisableOrganization against the database
func (conn *dbconn) DisableOrganization(_ID string, _UpdatedBy string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "DisableOrganization",
		Parameters: []interface{}{_ID, _UpdatedBy},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// DisableSource executes the stored procedure DisableSource against the database
func (conn *dbconn) DisableSource(_ID string, _OrgID string, _UpdatedBy string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "DisableSource",
		Parameters: []interface{}{_ID, _OrgID, _UpdatedBy},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// GetAllDetectionInfo executes the stored procedure GetAllDetectionInfo against the database and returns the read results
func (conn *dbconn) GetAllDetectionInfo(_OrgID string) ([]domain.DetectionInfo, error) {
	var err error
	var retDetectionInfo = make([]domain.DetectionInfo, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetAllDetectionInfo",
		Parameters: []interface{}{_OrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myOrganizationID string
							var mySourceID string
							var myDeviceID string
							var myVulnerabilityID string
							var myIgnoreID *string
							var myAlertDate time.Time
							var myLastFound *time.Time
							var myLastUpdated *time.Time
							var myProof string
							var myPort int
							var myProtocol string
							var myActiveKernel *int
							var myDetectionStatusID int
							var myTimesSeen int
							var myUpdated time.Time

							if err = rows.Scan(

								&myID,
								&myOrganizationID,
								&mySourceID,
								&myDeviceID,
								&myVulnerabilityID,
								&myIgnoreID,
								&myAlertDate,
								&myLastFound,
								&myLastUpdated,
								&myProof,
								&myPort,
								&myProtocol,
								&myActiveKernel,
								&myDetectionStatusID,
								&myTimesSeen,
								&myUpdated,
							); err == nil {

								newDetectionInfo := &dal.DetectionInfo{
									IDvar:                myID,
									OrganizationIDvar:    myOrganizationID,
									SourceIDvar:          mySourceID,
									DeviceIDvar:          myDeviceID,
									VulnerabilityIDvar:   myVulnerabilityID,
									IgnoreIDvar:          myIgnoreID,
									AlertDatevar:         myAlertDate,
									LastFoundvar:         myLastFound,
									LastUpdatedvar:       myLastUpdated,
									Proofvar:             myProof,
									Portvar:              myPort,
									Protocolvar:          myProtocol,
									ActiveKernelvar:      myActiveKernel,
									DetectionStatusIDvar: myDetectionStatusID,
									TimesSeenvar:         myTimesSeen,
									Updatedvar:           myUpdated,
								}

								retDetectionInfo = append(retDetectionInfo, newDetectionInfo)
							}
						}

						return err
					})
			}
		},
	})

	return retDetectionInfo, err
}

// GetAllExceptions executes the stored procedure GetAllExceptions against the database and returns the read results
func (conn *dbconn) GetAllExceptions(_offset int, _limit int, _sourceID string, _orgID string, _typeID int, _vulnID string, _devID string, _dueDate time.Time, _port string, _approval string, _active bool, _dBCreatedDate time.Time, _dBUpdatedDate time.Time, _updatedBy string, _createdBy string, _sortField string, _sortOrder string) ([]domain.Ignore, error) {
	var err error
	var retIgnore = make([]domain.Ignore, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetAllExceptions",
		Parameters: []interface{}{_offset, _limit, _sourceID, _orgID, _typeID, _vulnID, _devID, _dueDate, _port, _approval, _active, _dBCreatedDate, _dBUpdatedDate, _updatedBy, _createdBy, _sortField, _sortOrder},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var mySourceID string
							var myOrganizationID string
							var myTypeID int
							var myVulnerabilityID string
							var myDeviceID string
							var myDueDate *time.Time
							var myApproval string
							var myActive []uint8
							var myPort string
							var myCreatedBy *string
							var myUpdatedBy *string
							var myDBCreatedDate time.Time
							var myDBUpdatedDate *time.Time

							if err = rows.Scan(

								&myID,
								&mySourceID,
								&myOrganizationID,
								&myTypeID,
								&myVulnerabilityID,
								&myDeviceID,
								&myDueDate,
								&myApproval,
								&myActive,
								&myPort,
								&myCreatedBy,
								&myUpdatedBy,
								&myDBCreatedDate,
								&myDBUpdatedDate,
							); err == nil {

								newIgnore := &dal.Ignore{
									IDvar:              myID,
									SourceIDvar:        mySourceID,
									OrganizationIDvar:  myOrganizationID,
									TypeIDvar:          myTypeID,
									VulnerabilityIDvar: myVulnerabilityID,
									DeviceIDvar:        myDeviceID,
									DueDatevar:         myDueDate,
									Approvalvar:        myApproval,
									Activevar:          myActive[0] > 0 && myActive[0] != 48, // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)
									Portvar:            myPort,
									CreatedByvar:       myCreatedBy,
									UpdatedByvar:       myUpdatedBy,
									DBCreatedDatevar:   myDBCreatedDate,
									DBUpdatedDatevar:   myDBUpdatedDate,
								}

								retIgnore = append(retIgnore, newIgnore)
							}
						}

						return err
					})
			}
		},
	})

	return retIgnore, err
}

// GetAllJobConfigs executes the stored procedure GetAllJobConfigs against the database and returns the read results
func (conn *dbconn) GetAllJobConfigs(_OrgID string) ([]domain.JobConfig, error) {
	var err error
	var retJobConfig = make([]domain.JobConfig, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetAllJobConfigs",
		Parameters: []interface{}{_OrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myJobID int
							var myOrganizationID string
							var myDataInSourceConfigID *string
							var myDataOutSourceConfigID *string
							var myPriorityOverride *int
							var myContinuous []uint8
							var myWaitInSeconds int
							var myMaxInstances int
							var myAutoStart []uint8
							var myCreatedDate time.Time
							var myCreatedBy string
							var myUpdatedDate *time.Time
							var myUpdatedBy *string
							var myPayload *string

							if err = rows.Scan(

								&myID,
								&myJobID,
								&myOrganizationID,
								&myDataInSourceConfigID,
								&myDataOutSourceConfigID,
								&myPriorityOverride,
								&myContinuous,
								&myWaitInSeconds,
								&myMaxInstances,
								&myAutoStart,
								&myCreatedDate,
								&myCreatedBy,
								&myUpdatedDate,
								&myUpdatedBy,
								&myPayload,
							); err == nil {

								newJobConfig := &dal.JobConfig{
									IDvar:                    myID,
									JobIDvar:                 myJobID,
									OrganizationIDvar:        myOrganizationID,
									DataInSourceConfigIDvar:  myDataInSourceConfigID,
									DataOutSourceConfigIDvar: myDataOutSourceConfigID,
									PriorityOverridevar:      myPriorityOverride,
									Continuousvar:            myContinuous[0] > 0 && myContinuous[0] != 48, // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)
									WaitInSecondsvar:         myWaitInSeconds,
									MaxInstancesvar:          myMaxInstances,
									AutoStartvar:             myAutoStart[0] > 0 && myAutoStart[0] != 48, // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)
									CreatedDatevar:           myCreatedDate,
									CreatedByvar:             myCreatedBy,
									UpdatedDatevar:           myUpdatedDate,
									UpdatedByvar:             myUpdatedBy,
									Payloadvar:               myPayload,
								}

								retJobConfig = append(retJobConfig, newJobConfig)
							}
						}

						return err
					})
			}
		},
	})

	return retJobConfig, err
}

// GetAllJobConfigsWithOrder executes the stored procedure GetAllJobConfigsWithOrder against the database and returns the read results
func (conn *dbconn) GetAllJobConfigsWithOrder(_offset int, _limit int, _configID string, _jobid int, _dataInSourceConfigID string, _dataOutSourceConfigID string, _priorityOverride int, _continuous bool, _Payload string, _waitInSeconds int, _maxInstances int, _autoStart bool, _OrgID string, _updatedBy string, _createdBy string, _sortField string, _sortOrder string, _updatedDate time.Time, _createdDate time.Time, _lastJobStart time.Time, _ID string) ([]domain.JobConfig, error) {
	var err error
	var retJobConfig = make([]domain.JobConfig, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetAllJobConfigsWithOrder",
		Parameters: []interface{}{_offset, _limit, _configID, _jobid, _dataInSourceConfigID, _dataOutSourceConfigID, _priorityOverride, _continuous, _Payload, _waitInSeconds, _maxInstances, _autoStart, _OrgID, _updatedBy, _createdBy, _sortField, _sortOrder, _updatedDate, _createdDate, _lastJobStart, _ID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myJobID int
							var myOrganizationID string
							var myDataInSourceConfigID *string
							var myDataOutSourceConfigID *string
							var myPriorityOverride *int
							var myContinuous []uint8
							var myPayload *string
							var myWaitInSeconds int
							var myMaxInstances int
							var myAutoStart []uint8
							var myCreatedDate time.Time
							var myCreatedBy string
							var myUpdatedDate *time.Time
							var myUpdatedBy *string
							var myLastJobStart *time.Time
							var myActive []uint8

							if err = rows.Scan(

								&myID,
								&myJobID,
								&myOrganizationID,
								&myDataInSourceConfigID,
								&myDataOutSourceConfigID,
								&myPriorityOverride,
								&myContinuous,
								&myPayload,
								&myWaitInSeconds,
								&myMaxInstances,
								&myAutoStart,
								&myCreatedDate,
								&myCreatedBy,
								&myUpdatedDate,
								&myUpdatedBy,
								&myLastJobStart,
								&myActive,
							); err == nil {

								newJobConfig := &dal.JobConfig{
									IDvar:                    myID,
									JobIDvar:                 myJobID,
									OrganizationIDvar:        myOrganizationID,
									DataInSourceConfigIDvar:  myDataInSourceConfigID,
									DataOutSourceConfigIDvar: myDataOutSourceConfigID,
									PriorityOverridevar:      myPriorityOverride,
									Continuousvar:            myContinuous[0] > 0 && myContinuous[0] != 48, // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)
									Payloadvar:               myPayload,
									WaitInSecondsvar:         myWaitInSeconds,
									MaxInstancesvar:          myMaxInstances,
									AutoStartvar:             myAutoStart[0] > 0 && myAutoStart[0] != 48, // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)
									CreatedDatevar:           myCreatedDate,
									CreatedByvar:             myCreatedBy,
									UpdatedDatevar:           myUpdatedDate,
									UpdatedByvar:             myUpdatedBy,
									LastJobStartvar:          myLastJobStart,
									Activevar:                myActive[0] > 0 && myActive[0] != 48, // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)

								}

								retJobConfig = append(retJobConfig, newJobConfig)
							}
						}

						return err
					})
			}
		},
	})

	return retJobConfig, err
}

// GetAssetGroup executes the stored procedure GetAssetGroup against the database and returns the read results
func (conn *dbconn) GetAssetGroup(inOrgID string, _GroupID string, _ScannerConfigSourceID string) (domain.AssetGroup, error) {
	var err error
	var retAssetGroup domain.AssetGroup

	conn.Read(&connection.Procedure{
		Proc:       "GetAssetGroup",
		Parameters: []interface{}{inOrgID, _GroupID, _ScannerConfigSourceID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myGroupID string
							var myOrganizationID string
							var myScannerSourceConfigID *string
							var myScannerSourceID string
							var myCloudSourceID *string
							var myLastTicketing *time.Time
							var myRescanQueueSkip []uint8

							if err = rows.Scan(

								&myGroupID,
								&myOrganizationID,
								&myScannerSourceConfigID,
								&myScannerSourceID,
								&myCloudSourceID,
								&myLastTicketing,
								&myRescanQueueSkip,
							); err == nil {

								newAssetGroup := &dal.AssetGroup{
									GroupIDvar:               myGroupID,
									OrganizationIDvar:        myOrganizationID,
									ScannerSourceConfigIDvar: myScannerSourceConfigID,
									ScannerSourceIDvar:       myScannerSourceID,
									CloudSourceIDvar:         myCloudSourceID,
									LastTicketingvar:         myLastTicketing,
									RescanQueueSkipvar:       myRescanQueueSkip[0] > 0 && myRescanQueueSkip[0] != 48, // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)

								}

								retAssetGroup = newAssetGroup
							}
						}

						return err
					})
			}
		},
	})

	return retAssetGroup, err
}

// GetAssetGroupForOrg executes the stored procedure GetAssetGroupForOrg against the database and returns the read results
func (conn *dbconn) GetAssetGroupForOrg(inScannerSourceConfigID string, inOrgID string) ([]domain.AssetGroup, error) {
	var err error
	var retAssetGroup = make([]domain.AssetGroup, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetAssetGroupForOrg",
		Parameters: []interface{}{inScannerSourceConfigID, inOrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myGroupID string
							var myOrganizationID string
							var myScannerSourceConfigID *string
							var myScannerSourceID string
							var myCloudSourceID *string
							var myLastTicketing *time.Time
							var myRescanQueueSkip []uint8

							if err = rows.Scan(

								&myGroupID,
								&myOrganizationID,
								&myScannerSourceConfigID,
								&myScannerSourceID,
								&myCloudSourceID,
								&myLastTicketing,
								&myRescanQueueSkip,
							); err == nil {

								newAssetGroup := &dal.AssetGroup{
									GroupIDvar:               myGroupID,
									OrganizationIDvar:        myOrganizationID,
									ScannerSourceConfigIDvar: myScannerSourceConfigID,
									ScannerSourceIDvar:       myScannerSourceID,
									CloudSourceIDvar:         myCloudSourceID,
									LastTicketingvar:         myLastTicketing,
									RescanQueueSkipvar:       myRescanQueueSkip[0] > 0 && myRescanQueueSkip[0] != 48, // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)

								}

								retAssetGroup = append(retAssetGroup, newAssetGroup)
							}
						}

						return err
					})
			}
		},
	})

	return retAssetGroup, err
}

// GetAssetGroupForOrgNoScanner executes the stored procedure GetAssetGroupForOrgNoScanner against the database and returns the read results
func (conn *dbconn) GetAssetGroupForOrgNoScanner(inOrgID string, inGroupID string) (domain.AssetGroup, error) {
	var err error
	var retAssetGroup domain.AssetGroup

	conn.Read(&connection.Procedure{
		Proc:       "GetAssetGroupForOrgNoScanner",
		Parameters: []interface{}{inOrgID, inGroupID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myGroupID string
							var myOrganizationID string
							var myScannerSourceConfigID *string
							var myScannerSourceID string
							var myCloudSourceID *string
							var myLastTicketing *time.Time
							var myRescanQueueSkip []uint8

							if err = rows.Scan(

								&myGroupID,
								&myOrganizationID,
								&myScannerSourceConfigID,
								&myScannerSourceID,
								&myCloudSourceID,
								&myLastTicketing,
								&myRescanQueueSkip,
							); err == nil {

								newAssetGroup := &dal.AssetGroup{
									GroupIDvar:               myGroupID,
									OrganizationIDvar:        myOrganizationID,
									ScannerSourceConfigIDvar: myScannerSourceConfigID,
									ScannerSourceIDvar:       myScannerSourceID,
									CloudSourceIDvar:         myCloudSourceID,
									LastTicketingvar:         myLastTicketing,
									RescanQueueSkipvar:       myRescanQueueSkip[0] > 0 && myRescanQueueSkip[0] != 48, // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)

								}

								retAssetGroup = newAssetGroup
							}
						}

						return err
					})
			}
		},
	})

	return retAssetGroup, err
}

// GetAssetGroupsByCloudSource executes the stored procedure GetAssetGroupsByCloudSource against the database and returns the read results
func (conn *dbconn) GetAssetGroupsByCloudSource(inOrgID string, inCloudSourceID string) ([]domain.AssetGroup, error) {
	var err error
	var retAssetGroup = make([]domain.AssetGroup, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetAssetGroupsByCloudSource",
		Parameters: []interface{}{inOrgID, inCloudSourceID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myGroupID string
							var myOrganizationID string
							var myScannerSourceConfigID *string
							var myScannerSourceID string
							var myCloudSourceID *string
							var myLastTicketing *time.Time
							var myRescanQueueSkip []uint8

							if err = rows.Scan(

								&myGroupID,
								&myOrganizationID,
								&myScannerSourceConfigID,
								&myScannerSourceID,
								&myCloudSourceID,
								&myLastTicketing,
								&myRescanQueueSkip,
							); err == nil {

								newAssetGroup := &dal.AssetGroup{
									GroupIDvar:               myGroupID,
									OrganizationIDvar:        myOrganizationID,
									ScannerSourceConfigIDvar: myScannerSourceConfigID,
									ScannerSourceIDvar:       myScannerSourceID,
									CloudSourceIDvar:         myCloudSourceID,
									LastTicketingvar:         myLastTicketing,
									RescanQueueSkipvar:       myRescanQueueSkip[0] > 0 && myRescanQueueSkip[0] != 48, // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)

								}

								retAssetGroup = append(retAssetGroup, newAssetGroup)
							}
						}

						return err
					})
			}
		},
	})

	return retAssetGroup, err
}

// GetAssetGroupsForOrg executes the stored procedure GetAssetGroupsForOrg against the database and returns the read results
func (conn *dbconn) GetAssetGroupsForOrg(inOrgID string) ([]domain.AssetGroup, error) {
	var err error
	var retAssetGroup = make([]domain.AssetGroup, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetAssetGroupsForOrg",
		Parameters: []interface{}{inOrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myGroupID string
							var myOrganizationID string
							var myScannerSourceConfigID *string
							var myScannerSourceID string
							var myCloudSourceID *string
							var myLastTicketing *time.Time
							var myRescanQueueSkip []uint8

							if err = rows.Scan(

								&myGroupID,
								&myOrganizationID,
								&myScannerSourceConfigID,
								&myScannerSourceID,
								&myCloudSourceID,
								&myLastTicketing,
								&myRescanQueueSkip,
							); err == nil {

								newAssetGroup := &dal.AssetGroup{
									GroupIDvar:               myGroupID,
									OrganizationIDvar:        myOrganizationID,
									ScannerSourceConfigIDvar: myScannerSourceConfigID,
									ScannerSourceIDvar:       myScannerSourceID,
									CloudSourceIDvar:         myCloudSourceID,
									LastTicketingvar:         myLastTicketing,
									RescanQueueSkipvar:       myRescanQueueSkip[0] > 0 && myRescanQueueSkip[0] != 48, // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)

								}

								retAssetGroup = append(retAssetGroup, newAssetGroup)
							}
						}

						return err
					})
			}
		},
	})

	return retAssetGroup, err
}

// GetAssignmentGroupByIP executes the stored procedure GetAssignmentGroupByIP against the database and returns the read results
func (conn *dbconn) GetAssignmentGroupByIP(_SourceID string, _OrganizationID string, _IP string) ([]domain.AssignmentGroup, error) {
	var err error
	var retAssignmentGroup = make([]domain.AssignmentGroup, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetAssignmentGroupByIP",
		Parameters: []interface{}{_SourceID, _OrganizationID, _IP},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var mySourceID int
							var myOrganizationID string
							var myIPAddress string
							var myGroupName string
							var myDBCreatedDate time.Time
							var myDBUpdatedDate *time.Time

							if err = rows.Scan(

								&mySourceID,
								&myOrganizationID,
								&myIPAddress,
								&myGroupName,
								&myDBCreatedDate,
								&myDBUpdatedDate,
							); err == nil {

								newAssignmentGroup := &dal.AssignmentGroup{
									SourceIDvar:       mySourceID,
									OrganizationIDvar: myOrganizationID,
									IPAddressvar:      myIPAddress,
									GroupNamevar:      myGroupName,
									DBCreatedDatevar:  myDBCreatedDate,
									DBUpdatedDatevar:  myDBUpdatedDate,
								}

								retAssignmentGroup = append(retAssignmentGroup, newAssignmentGroup)
							}
						}

						return err
					})
			}
		},
	})

	return retAssignmentGroup, err
}

// GetAssignmentGroupByOrgIP executes the stored procedure GetAssignmentGroupByOrgIP against the database and returns the read results
func (conn *dbconn) GetAssignmentGroupByOrgIP(_OrganizationID string, _IP string) ([]domain.AssignmentGroup, error) {
	var err error
	var retAssignmentGroup = make([]domain.AssignmentGroup, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetAssignmentGroupByOrgIP",
		Parameters: []interface{}{_OrganizationID, _IP},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var mySourceID int
							var myOrganizationID string
							var myIPAddress string
							var myGroupName string
							var myDBCreatedDate time.Time
							var myDBUpdatedDate *time.Time

							if err = rows.Scan(

								&mySourceID,
								&myOrganizationID,
								&myIPAddress,
								&myGroupName,
								&myDBCreatedDate,
								&myDBUpdatedDate,
							); err == nil {

								newAssignmentGroup := &dal.AssignmentGroup{
									SourceIDvar:       mySourceID,
									OrganizationIDvar: myOrganizationID,
									IPAddressvar:      myIPAddress,
									GroupNamevar:      myGroupName,
									DBCreatedDatevar:  myDBCreatedDate,
									DBUpdatedDatevar:  myDBUpdatedDate,
								}

								retAssignmentGroup = append(retAssignmentGroup, newAssignmentGroup)
							}
						}

						return err
					})
			}
		},
	})

	return retAssignmentGroup, err
}

// GetAssignmentRulesByOrg executes the stored procedure GetAssignmentRulesByOrg against the database and returns the read results
func (conn *dbconn) GetAssignmentRulesByOrg(_OrganizationID string) ([]domain.AssignmentRules, error) {
	var err error
	var retAssignmentRules = make([]domain.AssignmentRules, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetAssignmentRulesByOrg",
		Parameters: []interface{}{_OrganizationID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myAssignmentGroup *string
							var myAssignee *string
							var myOrganizationID string
							var myGroupID *string
							var myVulnTitleRegex *string
							var myExcludeVulnTitleRegex *string
							var myHostnameRegex *string
							var myOSRegex *string
							var myCategoryRegex *string
							var myTagKeyID *int
							var myTagKeyRegex *string
							var myPortCSV *string
							var myExcludePortCSV *string
							var myPriority int

							if err = rows.Scan(

								&myAssignmentGroup,
								&myAssignee,
								&myOrganizationID,
								&myGroupID,
								&myVulnTitleRegex,
								&myExcludeVulnTitleRegex,
								&myHostnameRegex,
								&myOSRegex,
								&myCategoryRegex,
								&myTagKeyID,
								&myTagKeyRegex,
								&myPortCSV,
								&myExcludePortCSV,
								&myPriority,
							); err == nil {

								newAssignmentRules := &dal.AssignmentRules{
									AssignmentGroupvar:       myAssignmentGroup,
									Assigneevar:              myAssignee,
									OrganizationIDvar:        myOrganizationID,
									GroupIDvar:               myGroupID,
									VulnTitleRegexvar:        myVulnTitleRegex,
									ExcludeVulnTitleRegexvar: myExcludeVulnTitleRegex,
									HostnameRegexvar:         myHostnameRegex,
									OSRegexvar:               myOSRegex,
									CategoryRegexvar:         myCategoryRegex,
									TagKeyIDvar:              myTagKeyID,
									TagKeyRegexvar:           myTagKeyRegex,
									PortCSVvar:               myPortCSV,
									ExcludePortCSVvar:        myExcludePortCSV,
									Priorityvar:              myPriority,
								}

								retAssignmentRules = append(retAssignmentRules, newAssignmentRules)
							}
						}

						return err
					})
			}
		},
	})

	return retAssignmentRules, err
}

// GetAutoStartJobs executes the stored procedure GetAutoStartJobs against the database and returns the read results
func (conn *dbconn) GetAutoStartJobs() ([]domain.JobConfig, error) {
	var err error
	var retJobConfig = make([]domain.JobConfig, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetAutoStartJobs",
		Parameters: []interface{}{},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myJobID int
							var myOrganizationID string
							var myDataInSourceConfigID string
							var myDataOutSourceConfigID string
							var myPriorityOverride *int
							var myContinuous []uint8
							var myWaitInSeconds int
							var myMaxInstances int
							var myAutoStart []uint8
							var myCreatedDate time.Time
							var myCreatedBy string
							var myUpdatedDate *time.Time
							var myUpdatedBy *string
							var myPayload *string

							if err = rows.Scan(

								&myID,
								&myJobID,
								&myOrganizationID,
								&myDataInSourceConfigID,
								&myDataOutSourceConfigID,
								&myPriorityOverride,
								&myContinuous,
								&myWaitInSeconds,
								&myMaxInstances,
								&myAutoStart,
								&myCreatedDate,
								&myCreatedBy,
								&myUpdatedDate,
								&myUpdatedBy,
								&myPayload,
							); err == nil {

								newJobConfig := &dal.JobConfig{
									IDvar:                    myID,
									JobIDvar:                 myJobID,
									OrganizationIDvar:        myOrganizationID,
									DataInSourceConfigIDvar:  &myDataInSourceConfigID,
									DataOutSourceConfigIDvar: &myDataOutSourceConfigID,
									PriorityOverridevar:      myPriorityOverride,
									Continuousvar:            myContinuous[0] > 0 && myContinuous[0] != 48, // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)
									WaitInSecondsvar:         myWaitInSeconds,
									MaxInstancesvar:          myMaxInstances,
									AutoStartvar:             myAutoStart[0] > 0 && myAutoStart[0] != 48, // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)
									CreatedDatevar:           myCreatedDate,
									CreatedByvar:             myCreatedBy,
									UpdatedDatevar:           myUpdatedDate,
									UpdatedByvar:             myUpdatedBy,
									Payloadvar:               myPayload,
								}

								retJobConfig = append(retJobConfig, newJobConfig)
							}
						}

						return err
					})
			}
		},
	})

	return retJobConfig, err
}

// GetCISAssignments executes the stored procedure GetCISAssignments against the database and returns the read results
func (conn *dbconn) GetCISAssignments(_OrganizationID string) ([]domain.CISAssignments, error) {
	var err error
	var retCISAssignments = make([]domain.CISAssignments, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetCISAssignments",
		Parameters: []interface{}{_OrganizationID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myOrganizationID string
							var myCloudAccountID *string
							var myBundleID *string
							var myRuleRegex *string
							var myRuleHash *string
							var myAssignmentGroup string

							if err = rows.Scan(

								&myOrganizationID,
								&myCloudAccountID,
								&myBundleID,
								&myRuleRegex,
								&myRuleHash,
								&myAssignmentGroup,
							); err == nil {

								newCISAssignments := &dal.CISAssignments{
									OrganizationIDvar:  myOrganizationID,
									CloudAccountIDvar:  myCloudAccountID,
									BundleIDvar:        myBundleID,
									RuleRegexvar:       myRuleRegex,
									RuleHashvar:        myRuleHash,
									AssignmentGroupvar: myAssignmentGroup,
								}

								retCISAssignments = append(retCISAssignments, newCISAssignments)
							}
						}

						return err
					})
			}
		},
	})

	return retCISAssignments, err
}

// GetCancelledJobs executes the stored procedure GetCancelledJobs against the database and returns the read results
func (conn *dbconn) GetCancelledJobs() ([]domain.JobHistory, error) {
	var err error
	var retJobHistory = make([]domain.JobHistory, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetCancelledJobs",
		Parameters: []interface{}{},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myJobID int
							var myConfigID string
							var myStatusID int
							var myParentJobID *string
							var myIdentifier *string
							var myPriority int
							var myCurrentIteration *int
							var myPayload string
							var myThreadID *string
							var myPulseDate *time.Time
							var myCreatedDate time.Time
							var myUpdatedDate *time.Time

							if err = rows.Scan(

								&myID,
								&myJobID,
								&myConfigID,
								&myStatusID,
								&myParentJobID,
								&myIdentifier,
								&myPriority,
								&myCurrentIteration,
								&myPayload,
								&myThreadID,
								&myPulseDate,
								&myCreatedDate,
								&myUpdatedDate,
							); err == nil {

								newJobHistory := &dal.JobHistory{
									IDvar:               myID,
									JobIDvar:            myJobID,
									ConfigIDvar:         myConfigID,
									StatusIDvar:         myStatusID,
									ParentJobIDvar:      myParentJobID,
									Identifiervar:       myIdentifier,
									Priorityvar:         myPriority,
									CurrentIterationvar: myCurrentIteration,
									Payloadvar:          myPayload,
									ThreadIDvar:         myThreadID,
									PulseDatevar:        myPulseDate,
									CreatedDatevar:      myCreatedDate,
									UpdatedDatevar:      myUpdatedDate,
								}

								retJobHistory = append(retJobHistory, newJobHistory)
							}
						}

						return err
					})
			}
		},
	})

	return retJobHistory, err
}

// GetCategoryByName executes the stored procedure GetCategoryByName against the database and returns the read results
func (conn *dbconn) GetCategoryByName(_Name string) ([]domain.Category, error) {
	var err error
	var retCategory = make([]domain.Category, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetCategoryByName",
		Parameters: []interface{}{_Name},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myCategory string
							var myParentCategoryID *string

							if err = rows.Scan(

								&myID,
								&myCategory,
								&myParentCategoryID,
							); err == nil {

								newCategory := &dal.Category{
									IDvar:               myID,
									Categoryvar:         myCategory,
									ParentCategoryIDvar: myParentCategoryID,
								}

								retCategory = append(retCategory, newCategory)
							}
						}

						return err
					})
			}
		},
	})

	return retCategory, err
}

// GetDetectionInfo executes the stored procedure GetDetectionInfo against the database and returns the read results
func (conn *dbconn) GetDetectionInfo(_DeviceID string, _VulnerabilityID string, _Port int, _Protocol string) (domain.DetectionInfo, error) {
	var err error
	var retDetectionInfo domain.DetectionInfo

	conn.Read(&connection.Procedure{
		Proc:       "GetDetectionInfo",
		Parameters: []interface{}{_DeviceID, _VulnerabilityID, _Port, _Protocol},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myOrganizationID string
							var mySourceID string
							var myDeviceID string
							var myVulnerabilityID string
							var myIgnoreID *string
							var myAlertDate time.Time
							var myLastFound *time.Time
							var myLastUpdated *time.Time
							var myProof string
							var myPort int
							var myProtocol string
							var myActiveKernel *int
							var myDetectionStatusID int
							var myTimesSeen int
							var myUpdated time.Time

							if err = rows.Scan(

								&myID,
								&myOrganizationID,
								&mySourceID,
								&myDeviceID,
								&myVulnerabilityID,
								&myIgnoreID,
								&myAlertDate,
								&myLastFound,
								&myLastUpdated,
								&myProof,
								&myPort,
								&myProtocol,
								&myActiveKernel,
								&myDetectionStatusID,
								&myTimesSeen,
								&myUpdated,
							); err == nil {

								newDetectionInfo := &dal.DetectionInfo{
									IDvar:                myID,
									OrganizationIDvar:    myOrganizationID,
									SourceIDvar:          mySourceID,
									DeviceIDvar:          myDeviceID,
									VulnerabilityIDvar:   myVulnerabilityID,
									IgnoreIDvar:          myIgnoreID,
									AlertDatevar:         myAlertDate,
									LastFoundvar:         myLastFound,
									LastUpdatedvar:       myLastUpdated,
									Proofvar:             myProof,
									Portvar:              myPort,
									Protocolvar:          myProtocol,
									ActiveKernelvar:      myActiveKernel,
									DetectionStatusIDvar: myDetectionStatusID,
									TimesSeenvar:         myTimesSeen,
									Updatedvar:           myUpdated,
								}

								retDetectionInfo = newDetectionInfo
							}
						}

						return err
					})
			}
		},
	})

	return retDetectionInfo, err
}

// GetDetectionInfoAfter executes the stored procedure GetDetectionInfoAfter against the database and returns the read results
func (conn *dbconn) GetDetectionInfoAfter(_After time.Time, _OrgID string) ([]domain.DetectionInfo, error) {
	var err error
	var retDetectionInfo = make([]domain.DetectionInfo, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetDetectionInfoAfter",
		Parameters: []interface{}{_After, _OrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myOrganizationID string
							var mySourceID string
							var myDeviceID string
							var myVulnerabilityID string
							var myIgnoreID *string
							var myAlertDate time.Time
							var myLastFound *time.Time
							var myLastUpdated *time.Time
							var myProof string
							var myPort int
							var myProtocol string
							var myActiveKernel *int
							var myDetectionStatusID int
							var myTimesSeen int
							var myUpdated time.Time

							if err = rows.Scan(

								&myID,
								&myOrganizationID,
								&mySourceID,
								&myDeviceID,
								&myVulnerabilityID,
								&myIgnoreID,
								&myAlertDate,
								&myLastFound,
								&myLastUpdated,
								&myProof,
								&myPort,
								&myProtocol,
								&myActiveKernel,
								&myDetectionStatusID,
								&myTimesSeen,
								&myUpdated,
							); err == nil {

								newDetectionInfo := &dal.DetectionInfo{
									IDvar:                myID,
									OrganizationIDvar:    myOrganizationID,
									SourceIDvar:          mySourceID,
									DeviceIDvar:          myDeviceID,
									VulnerabilityIDvar:   myVulnerabilityID,
									IgnoreIDvar:          myIgnoreID,
									AlertDatevar:         myAlertDate,
									LastFoundvar:         myLastFound,
									LastUpdatedvar:       myLastUpdated,
									Proofvar:             myProof,
									Portvar:              myPort,
									Protocolvar:          myProtocol,
									ActiveKernelvar:      myActiveKernel,
									DetectionStatusIDvar: myDetectionStatusID,
									TimesSeenvar:         myTimesSeen,
									Updatedvar:           myUpdated,
								}

								retDetectionInfo = append(retDetectionInfo, newDetectionInfo)
							}
						}

						return err
					})
			}
		},
	})

	return retDetectionInfo, err
}

// GetDetectionInfoByID executes the stored procedure GetDetectionInfoByID against the database and returns the read results
func (conn *dbconn) GetDetectionInfoByID(_ID string, _OrgID string) (domain.DetectionInfo, error) {
	var err error
	var retDetectionInfo domain.DetectionInfo

	conn.Read(&connection.Procedure{
		Proc:       "GetDetectionInfoByID",
		Parameters: []interface{}{_ID, _OrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myOrganizationID string
							var mySourceID string
							var myDeviceID string
							var myVulnerabilityID string
							var myIgnoreID *string
							var myAlertDate time.Time
							var myLastFound *time.Time
							var myLastUpdated *time.Time
							var myProof string
							var myPort int
							var myProtocol string
							var myActiveKernel *int
							var myDetectionStatusID int
							var myTimesSeen int
							var myUpdated time.Time

							if err = rows.Scan(

								&myID,
								&myOrganizationID,
								&mySourceID,
								&myDeviceID,
								&myVulnerabilityID,
								&myIgnoreID,
								&myAlertDate,
								&myLastFound,
								&myLastUpdated,
								&myProof,
								&myPort,
								&myProtocol,
								&myActiveKernel,
								&myDetectionStatusID,
								&myTimesSeen,
								&myUpdated,
							); err == nil {

								newDetectionInfo := &dal.DetectionInfo{
									IDvar:                myID,
									OrganizationIDvar:    myOrganizationID,
									SourceIDvar:          mySourceID,
									DeviceIDvar:          myDeviceID,
									VulnerabilityIDvar:   myVulnerabilityID,
									IgnoreIDvar:          myIgnoreID,
									AlertDatevar:         myAlertDate,
									LastFoundvar:         myLastFound,
									LastUpdatedvar:       myLastUpdated,
									Proofvar:             myProof,
									Portvar:              myPort,
									Protocolvar:          myProtocol,
									ActiveKernelvar:      myActiveKernel,
									DetectionStatusIDvar: myDetectionStatusID,
									TimesSeenvar:         myTimesSeen,
									Updatedvar:           myUpdated,
								}

								retDetectionInfo = newDetectionInfo
							}
						}

						return err
					})
			}
		},
	})

	return retDetectionInfo, err
}

// GetDetectionInfoBySourceVulnID executes the stored procedure GetDetectionInfoBySourceVulnID against the database and returns the read results
func (conn *dbconn) GetDetectionInfoBySourceVulnID(_SourceDeviceID string, _SourceVulnerabilityID string, _Port int, _Protocol string) (domain.DetectionInfo, error) {
	var err error
	var retDetectionInfo domain.DetectionInfo

	conn.Read(&connection.Procedure{
		Proc:       "GetDetectionInfoBySourceVulnID",
		Parameters: []interface{}{_SourceDeviceID, _SourceVulnerabilityID, _Port, _Protocol},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myOrganizationID string
							var mySourceID string
							var myDeviceID string
							var myVulnerabilityID string
							var myIgnoreID *string
							var myAlertDate time.Time
							var myLastFound *time.Time
							var myLastUpdated *time.Time
							var myProof string
							var myPort int
							var myProtocol string
							var myActiveKernel *int
							var myDetectionStatusID int
							var myTimesSeen int
							var myUpdated time.Time

							if err = rows.Scan(

								&myID,
								&myOrganizationID,
								&mySourceID,
								&myDeviceID,
								&myVulnerabilityID,
								&myIgnoreID,
								&myAlertDate,
								&myLastFound,
								&myLastUpdated,
								&myProof,
								&myPort,
								&myProtocol,
								&myActiveKernel,
								&myDetectionStatusID,
								&myTimesSeen,
								&myUpdated,
							); err == nil {

								newDetectionInfo := &dal.DetectionInfo{
									IDvar:                myID,
									OrganizationIDvar:    myOrganizationID,
									SourceIDvar:          mySourceID,
									DeviceIDvar:          myDeviceID,
									VulnerabilityIDvar:   myVulnerabilityID,
									IgnoreIDvar:          myIgnoreID,
									AlertDatevar:         myAlertDate,
									LastFoundvar:         myLastFound,
									LastUpdatedvar:       myLastUpdated,
									Proofvar:             myProof,
									Portvar:              myPort,
									Protocolvar:          myProtocol,
									ActiveKernelvar:      myActiveKernel,
									DetectionStatusIDvar: myDetectionStatusID,
									TimesSeenvar:         myTimesSeen,
									Updatedvar:           myUpdated,
								}

								retDetectionInfo = newDetectionInfo
							}
						}

						return err
					})
			}
		},
	})

	return retDetectionInfo, err
}

// GetDetectionInfoForGroupAfter executes the stored procedure GetDetectionInfoForGroupAfter against the database and returns the read results
func (conn *dbconn) GetDetectionInfoForGroupAfter(_After time.Time, _OrgID string, inGroupID string, ticketInactiveKernels bool) ([]domain.DetectionInfo, error) {
	var err error
	var retDetectionInfo = make([]domain.DetectionInfo, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetDetectionInfoForGroupAfter",
		Parameters: []interface{}{_After, _OrgID, inGroupID, ticketInactiveKernels},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myOrganizationID string
							var mySourceID string
							var myDeviceID string
							var myVulnerabilityID string
							var myIgnoreID *string
							var myAlertDate time.Time
							var myLastFound *time.Time
							var myLastUpdated *time.Time
							var myProof string
							var myPort int
							var myProtocol string
							var myActiveKernel *int
							var myDetectionStatusID int
							var myTimesSeen int
							var myUpdated time.Time

							if err = rows.Scan(

								&myID,
								&myOrganizationID,
								&mySourceID,
								&myDeviceID,
								&myVulnerabilityID,
								&myIgnoreID,
								&myAlertDate,
								&myLastFound,
								&myLastUpdated,
								&myProof,
								&myPort,
								&myProtocol,
								&myActiveKernel,
								&myDetectionStatusID,
								&myTimesSeen,
								&myUpdated,
							); err == nil {

								newDetectionInfo := &dal.DetectionInfo{
									IDvar:                myID,
									OrganizationIDvar:    myOrganizationID,
									SourceIDvar:          mySourceID,
									DeviceIDvar:          myDeviceID,
									VulnerabilityIDvar:   myVulnerabilityID,
									IgnoreIDvar:          myIgnoreID,
									AlertDatevar:         myAlertDate,
									LastFoundvar:         myLastFound,
									LastUpdatedvar:       myLastUpdated,
									Proofvar:             myProof,
									Portvar:              myPort,
									Protocolvar:          myProtocol,
									ActiveKernelvar:      myActiveKernel,
									DetectionStatusIDvar: myDetectionStatusID,
									TimesSeenvar:         myTimesSeen,
									Updatedvar:           myUpdated,
								}

								retDetectionInfo = append(retDetectionInfo, newDetectionInfo)
							}
						}

						return err
					})
			}
		},
	})

	return retDetectionInfo, err
}

// GetDetectionStatusByID executes the stored procedure GetDetectionStatusByID against the database and returns the read results
func (conn *dbconn) GetDetectionStatusByID(_ID int) (domain.DetectionStatus, error) {
	var err error
	var retDetectionStatus domain.DetectionStatus

	conn.Read(&connection.Procedure{
		Proc:       "GetDetectionStatusByID",
		Parameters: []interface{}{_ID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID int
							var myStatus string
							var myName string

							if err = rows.Scan(

								&myID,
								&myStatus,
								&myName,
							); err == nil {

								newDetectionStatus := &dal.DetectionStatus{
									IDvar:     myID,
									Statusvar: myStatus,
									Namevar:   myName,
								}

								retDetectionStatus = newDetectionStatus
							}
						}

						return err
					})
			}
		},
	})

	return retDetectionStatus, err
}

// GetDetectionStatusByName executes the stored procedure GetDetectionStatusByName against the database and returns the read results
func (conn *dbconn) GetDetectionStatusByName(_Name string) (domain.DetectionStatus, error) {
	var err error
	var retDetectionStatus domain.DetectionStatus

	conn.Read(&connection.Procedure{
		Proc:       "GetDetectionStatusByName",
		Parameters: []interface{}{_Name},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID int
							var myStatus string
							var myName string

							if err = rows.Scan(

								&myID,
								&myStatus,
								&myName,
							); err == nil {

								newDetectionStatus := &dal.DetectionStatus{
									IDvar:     myID,
									Statusvar: myStatus,
									Namevar:   myName,
								}

								retDetectionStatus = newDetectionStatus
							}
						}

						return err
					})
			}
		},
	})

	return retDetectionStatus, err
}

// GetDetectionStatuses executes the stored procedure GetDetectionStatuses against the database and returns the read results
func (conn *dbconn) GetDetectionStatuses() ([]domain.DetectionStatus, error) {
	var err error
	var retDetectionStatus = make([]domain.DetectionStatus, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetDetectionStatuses",
		Parameters: []interface{}{},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID int
							var myStatus string
							var myName string

							if err = rows.Scan(

								&myID,
								&myStatus,
								&myName,
							); err == nil {

								newDetectionStatus := &dal.DetectionStatus{
									IDvar:     myID,
									Statusvar: myStatus,
									Namevar:   myName,
								}

								retDetectionStatus = append(retDetectionStatus, newDetectionStatus)
							}
						}

						return err
					})
			}
		},
	})

	return retDetectionStatus, err
}

// GetDetectionsInfoForDevice executes the stored procedure GetDetectionsInfoForDevice against the database and returns the read results
func (conn *dbconn) GetDetectionsInfoForDevice(_DeviceID string) ([]domain.DetectionInfo, error) {
	var err error
	var retDetectionInfo = make([]domain.DetectionInfo, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetDetectionsInfoForDevice",
		Parameters: []interface{}{_DeviceID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myOrganizationID string
							var mySourceID string
							var myDeviceID string
							var myVulnerabilityID string
							var myIgnoreID *string
							var myAlertDate time.Time
							var myLastFound *time.Time
							var myLastUpdated *time.Time
							var myProof string
							var myPort int
							var myProtocol string
							var myActiveKernel *int
							var myDetectionStatusID int
							var myTimesSeen int
							var myUpdated time.Time

							if err = rows.Scan(

								&myID,
								&myOrganizationID,
								&mySourceID,
								&myDeviceID,
								&myVulnerabilityID,
								&myIgnoreID,
								&myAlertDate,
								&myLastFound,
								&myLastUpdated,
								&myProof,
								&myPort,
								&myProtocol,
								&myActiveKernel,
								&myDetectionStatusID,
								&myTimesSeen,
								&myUpdated,
							); err == nil {

								newDetectionInfo := &dal.DetectionInfo{
									IDvar:                myID,
									OrganizationIDvar:    myOrganizationID,
									SourceIDvar:          mySourceID,
									DeviceIDvar:          myDeviceID,
									VulnerabilityIDvar:   myVulnerabilityID,
									IgnoreIDvar:          myIgnoreID,
									AlertDatevar:         myAlertDate,
									LastFoundvar:         myLastFound,
									LastUpdatedvar:       myLastUpdated,
									Proofvar:             myProof,
									Portvar:              myPort,
									Protocolvar:          myProtocol,
									ActiveKernelvar:      myActiveKernel,
									DetectionStatusIDvar: myDetectionStatusID,
									TimesSeenvar:         myTimesSeen,
									Updatedvar:           myUpdated,
								}

								retDetectionInfo = append(retDetectionInfo, newDetectionInfo)
							}
						}

						return err
					})
			}
		},
	})

	return retDetectionInfo, err
}

// GetDeviceInfoByAssetOrgID executes the stored procedure GetDeviceInfoByAssetOrgID against the database and returns the read results
func (conn *dbconn) GetDeviceInfoByAssetOrgID(inAssetID string, inOrgID string) (domain.DeviceInfo, error) {
	var err error
	var retDeviceInfo domain.DeviceInfo

	conn.Read(&connection.Procedure{
		Proc:       "GetDeviceInfoByAssetOrgID",
		Parameters: []interface{}{inAssetID, inOrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var mySourceID *string
							var myOS string
							var myMAC string
							var myIP string
							var myHostName string
							var myRegion *string
							var myGroupID *string
							var myInstanceID *string
							var myTrackingMethod *string

							if err = rows.Scan(

								&myID,
								&mySourceID,
								&myOS,
								&myMAC,
								&myIP,
								&myHostName,
								&myRegion,
								&myGroupID,
								&myInstanceID,
								&myTrackingMethod,
							); err == nil {

								newDeviceInfo := &dal.DeviceInfo{
									IDvar:             myID,
									SourceIDvar:       mySourceID,
									OSvar:             myOS,
									MACvar:            myMAC,
									IPvar:             myIP,
									HostNamevar:       myHostName,
									Regionvar:         myRegion,
									GroupIDvar:        myGroupID,
									InstanceIDvar:     myInstanceID,
									TrackingMethodvar: myTrackingMethod,
								}

								retDeviceInfo = newDeviceInfo
							}
						}

						return err
					})
			}
		},
	})

	return retDeviceInfo, err
}

// GetDeviceInfoByCloudSourceIDAndIP executes the stored procedure GetDeviceInfoByCloudSourceIDAndIP against the database and returns the read results
func (conn *dbconn) GetDeviceInfoByCloudSourceIDAndIP(_IP string, _CloudSourceID string, _OrgID string) ([]domain.DeviceInfo, error) {
	var err error
	var retDeviceInfo = make([]domain.DeviceInfo, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetDeviceInfoByCloudSourceIDAndIP",
		Parameters: []interface{}{_IP, _CloudSourceID, _OrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var mySourceID *string
							var myOS string
							var myMAC string
							var myIP string
							var myHostName string
							var myState *string
							var myRegion *string
							var myInstanceID *string
							var myScannerSourceID string
							var myTrackingMethod *string

							if err = rows.Scan(

								&myID,
								&mySourceID,
								&myOS,
								&myMAC,
								&myIP,
								&myHostName,
								&myState,
								&myRegion,
								&myInstanceID,
								&myScannerSourceID,
								&myTrackingMethod,
							); err == nil {

								newDeviceInfo := &dal.DeviceInfo{
									IDvar:              myID,
									SourceIDvar:        mySourceID,
									OSvar:              myOS,
									MACvar:             myMAC,
									IPvar:              myIP,
									HostNamevar:        myHostName,
									Statevar:           myState,
									Regionvar:          myRegion,
									InstanceIDvar:      myInstanceID,
									ScannerSourceIDvar: &myScannerSourceID,
									TrackingMethodvar:  myTrackingMethod,
								}

								retDeviceInfo = append(retDeviceInfo, newDeviceInfo)
							}
						}

						return err
					})
			}
		},
	})

	return retDeviceInfo, err
}

// GetDeviceInfoByGroupIP executes the stored procedure GetDeviceInfoByGroupIP against the database and returns the read results
func (conn *dbconn) GetDeviceInfoByGroupIP(inIP string, inGroupID string, inOrgID string) (domain.DeviceInfo, error) {
	var err error
	var retDeviceInfo domain.DeviceInfo

	conn.Read(&connection.Procedure{
		Proc:       "GetDeviceInfoByGroupIP",
		Parameters: []interface{}{inIP, inGroupID, inOrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var mySourceID *string
							var myOS string
							var myMAC string
							var myIP string
							var myHostName string
							var myRegion *string
							var myInstanceID *string
							var myTrackingMethod *string

							if err = rows.Scan(

								&myID,
								&mySourceID,
								&myOS,
								&myMAC,
								&myIP,
								&myHostName,
								&myRegion,
								&myInstanceID,
								&myTrackingMethod,
							); err == nil {

								newDeviceInfo := &dal.DeviceInfo{
									IDvar:             myID,
									SourceIDvar:       mySourceID,
									OSvar:             myOS,
									MACvar:            myMAC,
									IPvar:             myIP,
									HostNamevar:       myHostName,
									Regionvar:         myRegion,
									InstanceIDvar:     myInstanceID,
									TrackingMethodvar: myTrackingMethod,
								}

								retDeviceInfo = newDeviceInfo
							}
						}

						return err
					})
			}
		},
	})

	return retDeviceInfo, err
}

// GetDeviceInfoByIP executes the stored procedure GetDeviceInfoByIP against the database and returns the read results
func (conn *dbconn) GetDeviceInfoByIP(_IP string, _OrgID string) (domain.DeviceInfo, error) {
	var err error
	var retDeviceInfo domain.DeviceInfo

	conn.Read(&connection.Procedure{
		Proc:       "GetDeviceInfoByIP",
		Parameters: []interface{}{_IP, _OrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var mySourceID *string
							var myOS string
							var myMAC string
							var myIP string
							var myHostName string
							var myRegion *string
							var myInstanceID *string
							var myTrackingMethod *string

							if err = rows.Scan(

								&myID,
								&mySourceID,
								&myOS,
								&myMAC,
								&myIP,
								&myHostName,
								&myRegion,
								&myInstanceID,
								&myTrackingMethod,
							); err == nil {

								newDeviceInfo := &dal.DeviceInfo{
									IDvar:             myID,
									SourceIDvar:       mySourceID,
									OSvar:             myOS,
									MACvar:            myMAC,
									IPvar:             myIP,
									HostNamevar:       myHostName,
									Regionvar:         myRegion,
									InstanceIDvar:     myInstanceID,
									TrackingMethodvar: myTrackingMethod,
								}

								retDeviceInfo = newDeviceInfo
							}
						}

						return err
					})
			}
		},
	})

	return retDeviceInfo, err
}

// GetDeviceInfoByIPMACAndRegion executes the stored procedure GetDeviceInfoByIPMACAndRegion against the database and returns the read results
func (conn *dbconn) GetDeviceInfoByIPMACAndRegion(_IP string, _MAC string, _Region string, _OrgID string) (domain.DeviceInfo, error) {
	var err error
	var retDeviceInfo domain.DeviceInfo

	conn.Read(&connection.Procedure{
		Proc:       "GetDeviceInfoByIPMACAndRegion",
		Parameters: []interface{}{_IP, _MAC, _Region, _OrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var mySourceID *string
							var myOS string
							var myMAC string
							var myIP string
							var myHostName string
							var myRegion *string
							var myInstanceID *string
							var myTrackingMethod *string

							if err = rows.Scan(

								&myID,
								&mySourceID,
								&myOS,
								&myMAC,
								&myIP,
								&myHostName,
								&myRegion,
								&myInstanceID,
								&myTrackingMethod,
							); err == nil {

								newDeviceInfo := &dal.DeviceInfo{
									IDvar:             myID,
									SourceIDvar:       mySourceID,
									OSvar:             myOS,
									MACvar:            myMAC,
									IPvar:             myIP,
									HostNamevar:       myHostName,
									Regionvar:         myRegion,
									InstanceIDvar:     myInstanceID,
									TrackingMethodvar: myTrackingMethod,
								}

								retDeviceInfo = newDeviceInfo
							}
						}

						return err
					})
			}
		},
	})

	return retDeviceInfo, err
}

// GetDeviceInfoByInstanceID executes the stored procedure GetDeviceInfoByInstanceID against the database and returns the read results
func (conn *dbconn) GetDeviceInfoByInstanceID(_InstanceID string, _OrgID string) ([]domain.DeviceInfo, error) {
	var err error
	var retDeviceInfo = make([]domain.DeviceInfo, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetDeviceInfoByInstanceID",
		Parameters: []interface{}{_InstanceID, _OrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var mySourceID *string
							var myOS string
							var myMAC string
							var myIP string
							var myHostName string
							var myRegion *string
							var myInstanceID *string
							var myTrackingMethod *string

							if err = rows.Scan(

								&myID,
								&mySourceID,
								&myOS,
								&myMAC,
								&myIP,
								&myHostName,
								&myRegion,
								&myInstanceID,
								&myTrackingMethod,
							); err == nil {

								newDeviceInfo := &dal.DeviceInfo{
									IDvar:             myID,
									SourceIDvar:       mySourceID,
									OSvar:             myOS,
									MACvar:            myMAC,
									IPvar:             myIP,
									HostNamevar:       myHostName,
									Regionvar:         myRegion,
									InstanceIDvar:     myInstanceID,
									TrackingMethodvar: myTrackingMethod,
								}

								retDeviceInfo = append(retDeviceInfo, newDeviceInfo)
							}
						}

						return err
					})
			}
		},
	})

	return retDeviceInfo, err
}

// GetDeviceInfoByScannerSourceID executes the stored procedure GetDeviceInfoByScannerSourceID against the database and returns the read results
func (conn *dbconn) GetDeviceInfoByScannerSourceID(_IP string, _GroupID string, _OrgID string) (domain.DeviceInfo, error) {
	var err error
	var retDeviceInfo domain.DeviceInfo

	conn.Read(&connection.Procedure{
		Proc:       "GetDeviceInfoByScannerSourceID",
		Parameters: []interface{}{_IP, _GroupID, _OrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var mySourceID *string
							var myOS string
							var myMAC string
							var myIP string
							var myHostName string
							var myState *string
							var myRegion *string
							var myInstanceID *string
							var myScannerSourceID *string
							var myTrackingMethod *string

							if err = rows.Scan(

								&myID,
								&mySourceID,
								&myOS,
								&myMAC,
								&myIP,
								&myHostName,
								&myState,
								&myRegion,
								&myInstanceID,
								&myScannerSourceID,
								&myTrackingMethod,
							); err == nil {

								newDeviceInfo := &dal.DeviceInfo{
									IDvar:              myID,
									SourceIDvar:        mySourceID,
									OSvar:              myOS,
									MACvar:             myMAC,
									IPvar:              myIP,
									HostNamevar:        myHostName,
									Statevar:           myState,
									Regionvar:          myRegion,
									InstanceIDvar:      myInstanceID,
									ScannerSourceIDvar: myScannerSourceID,
									TrackingMethodvar:  myTrackingMethod,
								}

								retDeviceInfo = newDeviceInfo
							}
						}

						return err
					})
			}
		},
	})

	return retDeviceInfo, err
}

// GetDevicesInfoByCloudSourceID executes the stored procedure GetDevicesInfoByCloudSourceID against the database and returns the read results
func (conn *dbconn) GetDevicesInfoByCloudSourceID(_CloudSourceID string, _OrgID string) ([]domain.DeviceInfo, error) {
	var err error
	var retDeviceInfo = make([]domain.DeviceInfo, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetDevicesInfoByCloudSourceID",
		Parameters: []interface{}{_CloudSourceID, _OrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var mySourceID *string
							var myOS string
							var myMAC string
							var myIP string
							var myHostName string
							var myState *string
							var myRegion *string
							var myInstanceID *string
							var myScannerSourceID string
							var myTrackingMethod *string

							if err = rows.Scan(

								&myID,
								&mySourceID,
								&myOS,
								&myMAC,
								&myIP,
								&myHostName,
								&myState,
								&myRegion,
								&myInstanceID,
								&myScannerSourceID,
								&myTrackingMethod,
							); err == nil {

								newDeviceInfo := &dal.DeviceInfo{
									IDvar:              myID,
									SourceIDvar:        mySourceID,
									OSvar:              myOS,
									MACvar:             myMAC,
									IPvar:              myIP,
									HostNamevar:        myHostName,
									Statevar:           myState,
									Regionvar:          myRegion,
									InstanceIDvar:      myInstanceID,
									ScannerSourceIDvar: &myScannerSourceID,
									TrackingMethodvar:  myTrackingMethod,
								}

								retDeviceInfo = append(retDeviceInfo, newDeviceInfo)
							}
						}

						return err
					})
			}
		},
	})

	return retDeviceInfo, err
}

// GetDevicesInfoBySourceID executes the stored procedure GetDevicesInfoBySourceID against the database and returns the read results
func (conn *dbconn) GetDevicesInfoBySourceID(_SourceID string, _OrgID string) ([]domain.DeviceInfo, error) {
	var err error
	var retDeviceInfo = make([]domain.DeviceInfo, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetDevicesInfoBySourceID",
		Parameters: []interface{}{_SourceID, _OrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var mySourceID *string
							var myOS string
							var myMAC string
							var myIP string
							var myHostName string
							var myRegion *string
							var myInstanceID *string
							var myTrackingMethod *string

							if err = rows.Scan(

								&myID,
								&mySourceID,
								&myOS,
								&myMAC,
								&myIP,
								&myHostName,
								&myRegion,
								&myInstanceID,
								&myTrackingMethod,
							); err == nil {

								newDeviceInfo := &dal.DeviceInfo{
									IDvar:             myID,
									SourceIDvar:       mySourceID,
									OSvar:             myOS,
									MACvar:            myMAC,
									IPvar:             myIP,
									HostNamevar:       myHostName,
									Regionvar:         myRegion,
									InstanceIDvar:     myInstanceID,
									TrackingMethodvar: myTrackingMethod,
								}

								retDeviceInfo = append(retDeviceInfo, newDeviceInfo)
							}
						}

						return err
					})
			}
		},
	})

	return retDeviceInfo, err
}

// GetExceptionByVulnIDOrg executes the stored procedure GetExceptionByVulnIDOrg against the database and returns the read results
func (conn *dbconn) GetExceptionByVulnIDOrg(_DeviceID string, _VulnID string, _OrgID string) (domain.Ignore, error) {
	var err error
	var retIgnore domain.Ignore

	conn.Read(&connection.Procedure{
		Proc:       "GetExceptionByVulnIDOrg",
		Parameters: []interface{}{_DeviceID, _VulnID, _OrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myOrganizationID string
							var myVulnerabilityID string
							var myDeviceID string
							var myDueDate *time.Time

							if err = rows.Scan(

								&myID,
								&myOrganizationID,
								&myVulnerabilityID,
								&myDeviceID,
								&myDueDate,
							); err == nil {

								newIgnore := &dal.Ignore{
									IDvar:              myID,
									OrganizationIDvar:  myOrganizationID,
									VulnerabilityIDvar: myVulnerabilityID,
									DeviceIDvar:        myDeviceID,
									DueDatevar:         myDueDate,
								}

								retIgnore = newIgnore
							}
						}

						return err
					})
			}
		},
	})

	return retIgnore, err
}

// GetExceptionDetections executes the stored procedure GetExceptionDetections against the database and returns the read results
func (conn *dbconn) GetExceptionDetections(_offset int, _limit int, _orgID string, _sortField string, _sortOrder string, _Title string, _IP string, _Hostname string, _VulnID string, _Approval string, _DueDate string, _AssignmentGroup string, _OS string, _OSRegex string, _TypeID int) ([]domain.ExceptedDetection, error) {
	var err error
	var retExceptedDetection = make([]domain.ExceptedDetection, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetExceptionDetections",
		Parameters: []interface{}{_offset, _limit, _orgID, _sortField, _sortOrder, _Title, _IP, _Hostname, _VulnID, _Approval, _DueDate, _AssignmentGroup, _OS, _OSRegex, _TypeID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myTitle *string
							var myIP *string
							var myHostname *string
							var myVulnerabilityID *string
							var myVulnerabilityTitle *string
							var myApproval *string
							var myDueDate *time.Time
							var myAssignmentGroup *string
							var myOS *string
							var myOSRegex *string
							var myIgnoreID string
							var myIgnoreType int

							if err = rows.Scan(

								&myTitle,
								&myIP,
								&myHostname,
								&myVulnerabilityID,
								&myVulnerabilityTitle,
								&myApproval,
								&myDueDate,
								&myAssignmentGroup,
								&myOS,
								&myOSRegex,
								&myIgnoreID,
								&myIgnoreType,
							); err == nil {

								newExceptedDetection := &dal.ExceptedDetection{
									Titlevar:              myTitle,
									IPvar:                 myIP,
									Hostnamevar:           myHostname,
									VulnerabilityIDvar:    myVulnerabilityID,
									VulnerabilityTitlevar: myVulnerabilityTitle,
									Approvalvar:           myApproval,
									DueDatevar:            myDueDate,
									AssignmentGroupvar:    myAssignmentGroup,
									OSvar:                 myOS,
									OSRegexvar:            myOSRegex,
									IgnoreIDvar:           myIgnoreID,
									IgnoreTypevar:         myIgnoreType,
								}

								retExceptedDetection = append(retExceptedDetection, newExceptedDetection)
							}
						}

						return err
					})
			}
		},
	})

	return retExceptedDetection, err
}

// GetExceptionTypes executes the stored procedure GetExceptionTypes against the database and returns the read results
func (conn *dbconn) GetExceptionTypes() ([]domain.ExceptionType, error) {
	var err error
	var retExceptionType = make([]domain.ExceptionType, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetExceptionTypes",
		Parameters: []interface{}{},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID int
							var myType string
							var myName string

							if err = rows.Scan(

								&myID,
								&myType,
								&myName,
							); err == nil {

								newExceptionType := &dal.ExceptionType{
									IDvar:   myID,
									Typevar: myType,
									Namevar: myName,
								}

								retExceptionType = append(retExceptionType, newExceptionType)
							}
						}

						return err
					})
			}
		},
	})

	return retExceptionType, err
}

// GetExceptionsByOrg executes the stored procedure GetExceptionsByOrg against the database and returns the read results
func (conn *dbconn) GetExceptionsByOrg(_OrgID string) ([]domain.Ignore, error) {
	var err error
	var retIgnore = make([]domain.Ignore, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetExceptionsByOrg",
		Parameters: []interface{}{_OrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myTypeID int
							var myApproval string
							var myOrganizationID string
							var myVulnerabilityID string
							var myDeviceID string
							var myPort string
							var myDueDate *time.Time

							if err = rows.Scan(

								&myID,
								&myTypeID,
								&myApproval,
								&myOrganizationID,
								&myVulnerabilityID,
								&myDeviceID,
								&myPort,
								&myDueDate,
							); err == nil {

								newIgnore := &dal.Ignore{
									IDvar:              myID,
									TypeIDvar:          myTypeID,
									Approvalvar:        myApproval,
									OrganizationIDvar:  myOrganizationID,
									VulnerabilityIDvar: myVulnerabilityID,
									DeviceIDvar:        myDeviceID,
									Portvar:            myPort,
									DueDatevar:         myDueDate,
								}

								retIgnore = append(retIgnore, newIgnore)
							}
						}

						return err
					})
			}
		},
	})

	return retIgnore, err
}

// GetExceptionsDueNext30Days executes the stored procedure GetExceptionsDueNext30Days against the database and returns the read results
func (conn *dbconn) GetExceptionsDueNext30Days() ([]domain.CERF, error) {
	var err error
	var retCERF = make([]domain.CERF, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetExceptionsDueNext30Days",
		Parameters: []interface{}{},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myCERForm string

							if err = rows.Scan(

								&myCERForm,
							); err == nil {

								newCERF := &dal.CERF{
									CERFormvar: myCERForm,
								}

								retCERF = append(retCERF, newCERF)
							}
						}

						return err
					})
			}
		},
	})

	return retCERF, err
}

// GetExceptionsLength executes the stored procedure GetExceptionsLength against the database and returns the read results
func (conn *dbconn) GetExceptionsLength(_offset int, _limit int, _orgID string, _sortField string, _sortOrder string, _Title string, _IP string, _Hostname string, _VulnID string, _Approval string, _DueDate string, _AssignmentGroup string, _OS string, _OSRegex string, _TypeID int) (domain.QueryData, error) {
	var err error
	var retQueryData domain.QueryData

	conn.Read(&connection.Procedure{
		Proc:       "GetExceptionsLength",
		Parameters: []interface{}{_offset, _limit, _orgID, _sortField, _sortOrder, _Title, _IP, _Hostname, _VulnID, _Approval, _DueDate, _AssignmentGroup, _OS, _OSRegex, _TypeID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myLength int

							if err = rows.Scan(

								&myLength,
							); err == nil {

								newQueryData := &dal.QueryData{
									Lengthvar: myLength,
								}

								retQueryData = newQueryData
							}
						}

						return err
					})
			}
		},
	})

	return retQueryData, err
}

// GetGlobalExceptions executes the stored procedure GetGlobalExceptions against the database and returns the read results
func (conn *dbconn) GetGlobalExceptions(_OrgID string) ([]domain.Ignore, error) {
	var err error
	var retIgnore = make([]domain.Ignore, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetGlobalExceptions",
		Parameters: []interface{}{_OrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myOrganizationID string
							var myVulnerabilityID string
							var myOSRegex *string
							var myHostnameRegex *string
							var myDueDate *time.Time

							if err = rows.Scan(

								&myID,
								&myOrganizationID,
								&myVulnerabilityID,
								&myOSRegex,
								&myHostnameRegex,
								&myDueDate,
							); err == nil {

								newIgnore := &dal.Ignore{
									IDvar:              myID,
									OrganizationIDvar:  myOrganizationID,
									VulnerabilityIDvar: myVulnerabilityID,
									OSRegexvar:         myOSRegex,
									HostnameRegexvar:   myHostnameRegex,
									DueDatevar:         myDueDate,
								}

								retIgnore = append(retIgnore, newIgnore)
							}
						}

						return err
					})
			}
		},
	})

	return retIgnore, err
}

// GetJobByID executes the stored procedure GetJobByID against the database and returns the read results
func (conn *dbconn) GetJobByID(_ID int) (domain.JobRegistration, error) {
	var err error
	var retJobRegistration domain.JobRegistration

	conn.Read(&connection.Procedure{
		Proc:       "GetJobByID",
		Parameters: []interface{}{_ID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID int
							var myGoStruct string
							var myPriority int
							var myCreatedDate time.Time
							var myCreatedBy string
							var myUpdatedDate *time.Time
							var myUpdatedBy *string

							if err = rows.Scan(

								&myID,
								&myGoStruct,
								&myPriority,
								&myCreatedDate,
								&myCreatedBy,
								&myUpdatedDate,
								&myUpdatedBy,
							); err == nil {

								newJobRegistration := &dal.JobRegistration{
									IDvar:          myID,
									GoStructvar:    myGoStruct,
									Priorityvar:    myPriority,
									CreatedDatevar: myCreatedDate,
									CreatedByvar:   myCreatedBy,
									UpdatedDatevar: myUpdatedDate,
									UpdatedByvar:   myUpdatedBy,
								}

								retJobRegistration = newJobRegistration
							}
						}

						return err
					})
			}
		},
	})

	return retJobRegistration, err
}

// GetJobConfig executes the stored procedure GetJobConfig against the database and returns the read results
func (conn *dbconn) GetJobConfig(_ID string) (domain.JobConfig, error) {
	var err error
	var retJobConfig domain.JobConfig

	conn.Read(&connection.Procedure{
		Proc:       "GetJobConfig",
		Parameters: []interface{}{_ID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myJobID int
							var myOrganizationID string
							var myDataInSourceConfigID *string
							var myDataOutSourceConfigID *string
							var myPayload string
							var myPriorityOverride *int
							var myContinuous []uint8
							var myWaitInSeconds int
							var myMaxInstances int
							var myAutoStart []uint8
							var myCreatedDate time.Time
							var myCreatedBy string
							var myUpdatedDate *time.Time
							var myUpdatedBy *string
							var myLastJobStart *time.Time

							if err = rows.Scan(

								&myID,
								&myJobID,
								&myOrganizationID,
								&myDataInSourceConfigID,
								&myDataOutSourceConfigID,
								&myPayload,
								&myPriorityOverride,
								&myContinuous,
								&myWaitInSeconds,
								&myMaxInstances,
								&myAutoStart,
								&myCreatedDate,
								&myCreatedBy,
								&myUpdatedDate,
								&myUpdatedBy,
								&myLastJobStart,
							); err == nil {

								newJobConfig := &dal.JobConfig{
									IDvar:                    myID,
									JobIDvar:                 myJobID,
									OrganizationIDvar:        myOrganizationID,
									DataInSourceConfigIDvar:  myDataInSourceConfigID,
									DataOutSourceConfigIDvar: myDataOutSourceConfigID,
									Payloadvar:               &myPayload,
									PriorityOverridevar:      myPriorityOverride,
									Continuousvar:            myContinuous[0] > 0 && myContinuous[0] != 48, // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)
									WaitInSecondsvar:         myWaitInSeconds,
									MaxInstancesvar:          myMaxInstances,
									AutoStartvar:             myAutoStart[0] > 0 && myAutoStart[0] != 48, // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)
									CreatedDatevar:           myCreatedDate,
									CreatedByvar:             myCreatedBy,
									UpdatedDatevar:           myUpdatedDate,
									UpdatedByvar:             myUpdatedBy,
									LastJobStartvar:          myLastJobStart,
								}

								retJobConfig = newJobConfig
							}
						}

						return err
					})
			}
		},
	})

	return retJobConfig, err
}

// GetJobConfigAudit executes the stored procedure GetJobConfigAudit against the database and returns the read results
func (conn *dbconn) GetJobConfigAudit(inJobConfigID string, inOrgID string) ([]domain.JobConfigAudit, error) {
	var err error
	var retJobConfigAudit = make([]domain.JobConfigAudit, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetJobConfigAudit",
		Parameters: []interface{}{inJobConfigID, inOrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myJobID int
							var myOrganizationID string
							var myDataInSourceConfigID *string
							var myDataOutSourceConfigID *string
							var myPayload *string
							var myPriorityOverride *int
							var myContinuous []uint8
							var myWaitInSeconds int
							var myMaxInstances int
							var myAutoStart []uint8
							var myCreatedDate time.Time
							var myCreatedBy string
							var myUpdatedDate *time.Time
							var myUpdatedBy *string
							var myActive []uint8
							var myLastJobStart *time.Time
							var myEventType string
							var myEventDate time.Time

							if err = rows.Scan(

								&myID,
								&myJobID,
								&myOrganizationID,
								&myDataInSourceConfigID,
								&myDataOutSourceConfigID,
								&myPayload,
								&myPriorityOverride,
								&myContinuous,
								&myWaitInSeconds,
								&myMaxInstances,
								&myAutoStart,
								&myCreatedDate,
								&myCreatedBy,
								&myUpdatedDate,
								&myUpdatedBy,
								&myActive,
								&myLastJobStart,
								&myEventType,
								&myEventDate,
							); err == nil {

								newJobConfigAudit := &dal.JobConfigAudit{
									IDvar:                    myID,
									JobIDvar:                 myJobID,
									OrganizationIDvar:        myOrganizationID,
									DataInSourceConfigIDvar:  myDataInSourceConfigID,
									DataOutSourceConfigIDvar: myDataOutSourceConfigID,
									Payloadvar:               myPayload,
									PriorityOverridevar:      myPriorityOverride,
									Continuousvar:            myContinuous[0] > 0 && myContinuous[0] != 48, // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)
									WaitInSecondsvar:         myWaitInSeconds,
									MaxInstancesvar:          myMaxInstances,
									AutoStartvar:             myAutoStart[0] > 0 && myAutoStart[0] != 48, // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)
									CreatedDatevar:           myCreatedDate,
									CreatedByvar:             myCreatedBy,
									UpdatedDatevar:           myUpdatedDate,
									UpdatedByvar:             myUpdatedBy,
									Activevar:                myActive[0] > 0 && myActive[0] != 48, // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)
									LastJobStartvar:          myLastJobStart,
									EventTypevar:             myEventType,
									EventDatevar:             myEventDate,
								}

								retJobConfigAudit = append(retJobConfigAudit, newJobConfigAudit)
							}
						}

						return err
					})
			}
		},
	})

	return retJobConfigAudit, err
}

// GetJobConfigByID executes the stored procedure GetJobConfigByID against the database and returns the read results
func (conn *dbconn) GetJobConfigByID(_ID string, _OrgID string) (domain.JobConfig, error) {
	var err error
	var retJobConfig domain.JobConfig

	conn.Read(&connection.Procedure{
		Proc:       "GetJobConfigByID",
		Parameters: []interface{}{_ID, _OrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myJobID int
							var myOrganizationID string
							var myDataInSourceConfigID *string
							var myDataOutSourceConfigID *string
							var myPriorityOverride *int
							var myContinuous []uint8
							var myPayload *string
							var myWaitInSeconds int
							var myMaxInstances int
							var myAutoStart []uint8
							var myCreatedDate time.Time
							var myCreatedBy string
							var myUpdatedDate *time.Time
							var myUpdatedBy *string
							var myLastJobStart *time.Time

							if err = rows.Scan(

								&myID,
								&myJobID,
								&myOrganizationID,
								&myDataInSourceConfigID,
								&myDataOutSourceConfigID,
								&myPriorityOverride,
								&myContinuous,
								&myPayload,
								&myWaitInSeconds,
								&myMaxInstances,
								&myAutoStart,
								&myCreatedDate,
								&myCreatedBy,
								&myUpdatedDate,
								&myUpdatedBy,
								&myLastJobStart,
							); err == nil {

								newJobConfig := &dal.JobConfig{
									IDvar:                    myID,
									JobIDvar:                 myJobID,
									OrganizationIDvar:        myOrganizationID,
									DataInSourceConfigIDvar:  myDataInSourceConfigID,
									DataOutSourceConfigIDvar: myDataOutSourceConfigID,
									PriorityOverridevar:      myPriorityOverride,
									Continuousvar:            myContinuous[0] > 0 && myContinuous[0] != 48, // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)
									Payloadvar:               myPayload,
									WaitInSecondsvar:         myWaitInSeconds,
									MaxInstancesvar:          myMaxInstances,
									AutoStartvar:             myAutoStart[0] > 0 && myAutoStart[0] != 48, // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)
									CreatedDatevar:           myCreatedDate,
									CreatedByvar:             myCreatedBy,
									UpdatedDatevar:           myUpdatedDate,
									UpdatedByvar:             myUpdatedBy,
									LastJobStartvar:          myLastJobStart,
								}

								retJobConfig = newJobConfig
							}
						}

						return err
					})
			}
		},
	})

	return retJobConfig, err
}

// GetJobConfigByJobHistoryID executes the stored procedure GetJobConfigByJobHistoryID against the database and returns the read results
func (conn *dbconn) GetJobConfigByJobHistoryID(_JobHistoryID string) (domain.JobConfig, error) {
	var err error
	var retJobConfig domain.JobConfig

	conn.Read(&connection.Procedure{
		Proc:       "GetJobConfigByJobHistoryID",
		Parameters: []interface{}{_JobHistoryID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myJobID int
							var myOrganizationID string
							var myDataInSourceConfigID *string
							var myDataOutSourceConfigID *string
							var myPriorityOverride *int
							var myContinuous []uint8
							var myWaitInSeconds int
							var myMaxInstances int
							var myAutoStart []uint8
							var myCreatedDate time.Time
							var myCreatedBy string
							var myUpdatedDate *time.Time
							var myUpdatedBy *string
							var myLastJobStart *time.Time

							if err = rows.Scan(

								&myID,
								&myJobID,
								&myOrganizationID,
								&myDataInSourceConfigID,
								&myDataOutSourceConfigID,
								&myPriorityOverride,
								&myContinuous,
								&myWaitInSeconds,
								&myMaxInstances,
								&myAutoStart,
								&myCreatedDate,
								&myCreatedBy,
								&myUpdatedDate,
								&myUpdatedBy,
								&myLastJobStart,
							); err == nil {

								newJobConfig := &dal.JobConfig{
									IDvar:                    myID,
									JobIDvar:                 myJobID,
									OrganizationIDvar:        myOrganizationID,
									DataInSourceConfigIDvar:  myDataInSourceConfigID,
									DataOutSourceConfigIDvar: myDataOutSourceConfigID,
									PriorityOverridevar:      myPriorityOverride,
									Continuousvar:            myContinuous[0] > 0 && myContinuous[0] != 48, // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)
									WaitInSecondsvar:         myWaitInSeconds,
									MaxInstancesvar:          myMaxInstances,
									AutoStartvar:             myAutoStart[0] > 0 && myAutoStart[0] != 48, // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)
									CreatedDatevar:           myCreatedDate,
									CreatedByvar:             myCreatedBy,
									UpdatedDatevar:           myUpdatedDate,
									UpdatedByvar:             myUpdatedBy,
									LastJobStartvar:          myLastJobStart,
								}

								retJobConfig = newJobConfig
							}
						}

						return err
					})
			}
		},
	})

	return retJobConfig, err
}

// GetJobConfigByOrgIDAndJobID executes the stored procedure GetJobConfigByOrgIDAndJobID against the database and returns the read results
func (conn *dbconn) GetJobConfigByOrgIDAndJobID(_OrgID string, _JobID int) ([]domain.JobConfig, error) {
	var err error
	var retJobConfig = make([]domain.JobConfig, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetJobConfigByOrgIDAndJobID",
		Parameters: []interface{}{_OrgID, _JobID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myJobID int
							var myOrganizationID string
							var myDataInSourceConfigID *string
							var myDataOutSourceConfigID *string
							var myPriorityOverride *int
							var myContinuous []uint8
							var myWaitInSeconds int
							var myMaxInstances int
							var myAutoStart []uint8
							var myCreatedDate time.Time
							var myCreatedBy string
							var myUpdatedDate *time.Time
							var myUpdatedBy *string
							var myLastJobStart *time.Time

							if err = rows.Scan(

								&myID,
								&myJobID,
								&myOrganizationID,
								&myDataInSourceConfigID,
								&myDataOutSourceConfigID,
								&myPriorityOverride,
								&myContinuous,
								&myWaitInSeconds,
								&myMaxInstances,
								&myAutoStart,
								&myCreatedDate,
								&myCreatedBy,
								&myUpdatedDate,
								&myUpdatedBy,
								&myLastJobStart,
							); err == nil {

								newJobConfig := &dal.JobConfig{
									IDvar:                    myID,
									JobIDvar:                 myJobID,
									OrganizationIDvar:        myOrganizationID,
									DataInSourceConfigIDvar:  myDataInSourceConfigID,
									DataOutSourceConfigIDvar: myDataOutSourceConfigID,
									PriorityOverridevar:      myPriorityOverride,
									Continuousvar:            myContinuous[0] > 0 && myContinuous[0] != 48, // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)
									WaitInSecondsvar:         myWaitInSeconds,
									MaxInstancesvar:          myMaxInstances,
									AutoStartvar:             myAutoStart[0] > 0 && myAutoStart[0] != 48, // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)
									CreatedDatevar:           myCreatedDate,
									CreatedByvar:             myCreatedBy,
									UpdatedDatevar:           myUpdatedDate,
									UpdatedByvar:             myUpdatedBy,
									LastJobStartvar:          myLastJobStart,
								}

								retJobConfig = append(retJobConfig, newJobConfig)
							}
						}

						return err
					})
			}
		},
	})

	return retJobConfig, err
}

// GetJobConfigByOrgIDAndJobIDWithSC executes the stored procedure GetJobConfigByOrgIDAndJobIDWithSC against the database and returns the read results
func (conn *dbconn) GetJobConfigByOrgIDAndJobIDWithSC(_OrgID string, _JobID int, _SourceConfigID string) ([]domain.JobConfig, error) {
	var err error
	var retJobConfig = make([]domain.JobConfig, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetJobConfigByOrgIDAndJobIDWithSC",
		Parameters: []interface{}{_OrgID, _JobID, _SourceConfigID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myJobID int
							var myOrganizationID string
							var myDataInSourceConfigID *string
							var myDataOutSourceConfigID *string
							var myPriorityOverride *int
							var myContinuous []uint8
							var myWaitInSeconds int
							var myMaxInstances int
							var myAutoStart []uint8
							var myCreatedDate time.Time
							var myCreatedBy string
							var myUpdatedDate *time.Time
							var myUpdatedBy *string
							var myLastJobStart *time.Time

							if err = rows.Scan(

								&myID,
								&myJobID,
								&myOrganizationID,
								&myDataInSourceConfigID,
								&myDataOutSourceConfigID,
								&myPriorityOverride,
								&myContinuous,
								&myWaitInSeconds,
								&myMaxInstances,
								&myAutoStart,
								&myCreatedDate,
								&myCreatedBy,
								&myUpdatedDate,
								&myUpdatedBy,
								&myLastJobStart,
							); err == nil {

								newJobConfig := &dal.JobConfig{
									IDvar:                    myID,
									JobIDvar:                 myJobID,
									OrganizationIDvar:        myOrganizationID,
									DataInSourceConfigIDvar:  myDataInSourceConfigID,
									DataOutSourceConfigIDvar: myDataOutSourceConfigID,
									PriorityOverridevar:      myPriorityOverride,
									Continuousvar:            myContinuous[0] > 0 && myContinuous[0] != 48, // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)
									WaitInSecondsvar:         myWaitInSeconds,
									MaxInstancesvar:          myMaxInstances,
									AutoStartvar:             myAutoStart[0] > 0 && myAutoStart[0] != 48, // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)
									CreatedDatevar:           myCreatedDate,
									CreatedByvar:             myCreatedBy,
									UpdatedDatevar:           myUpdatedDate,
									UpdatedByvar:             myUpdatedBy,
									LastJobStartvar:          myLastJobStart,
								}

								retJobConfig = append(retJobConfig, newJobConfig)
							}
						}

						return err
					})
			}
		},
	})

	return retJobConfig, err
}

// GetJobConfigLength executes the stored procedure GetJobConfigLength against the database and returns the read results
func (conn *dbconn) GetJobConfigLength(_configID string, _jobID int, _dataInSourceConfigID string, _dataOutSourceConfigID string, _priorityOverride int, _continuous bool, _Payload string, _waitInSeconds int, _maxInstances int, _autoStart bool, _OrgID string, _updatedBy string, _createdBy string, _updatedDate time.Time, _createdDate time.Time, _lastJobStart time.Time, _ID string) (domain.QueryData, error) {
	var err error
	var retQueryData domain.QueryData

	conn.Read(&connection.Procedure{
		Proc:       "GetJobConfigLength",
		Parameters: []interface{}{_configID, _jobID, _dataInSourceConfigID, _dataOutSourceConfigID, _priorityOverride, _continuous, _Payload, _waitInSeconds, _maxInstances, _autoStart, _OrgID, _updatedBy, _createdBy, _updatedDate, _createdDate, _lastJobStart, _ID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myLength int

							if err = rows.Scan(

								&myLength,
							); err == nil {

								newQueryData := &dal.QueryData{
									Lengthvar: myLength,
								}

								retQueryData = newQueryData
							}
						}

						return err
					})
			}
		},
	})

	return retQueryData, err
}

// GetJobHistories executes the stored procedure GetJobHistories against the database and returns the read results
func (conn *dbconn) GetJobHistories(_offset int, _limit int, _jobID int, _jobconfig string, _status int, _Payload string, _OrgID string) ([]domain.JobHistory, error) {
	var err error
	var retJobHistory = make([]domain.JobHistory, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetJobHistories",
		Parameters: []interface{}{_offset, _limit, _jobID, _jobconfig, _status, _Payload, _OrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myJobID int
							var myConfigID string
							var myStatusID int
							var myParentJobID *string
							var myIdentifier *string
							var myPriority int
							var myCurrentIteration *int
							var myPayload string
							var myThreadID *string
							var myPulseDate *time.Time
							var myCreatedDate time.Time
							var myUpdatedDate *time.Time

							if err = rows.Scan(

								&myID,
								&myJobID,
								&myConfigID,
								&myStatusID,
								&myParentJobID,
								&myIdentifier,
								&myPriority,
								&myCurrentIteration,
								&myPayload,
								&myThreadID,
								&myPulseDate,
								&myCreatedDate,
								&myUpdatedDate,
							); err == nil {

								newJobHistory := &dal.JobHistory{
									IDvar:               myID,
									JobIDvar:            myJobID,
									ConfigIDvar:         myConfigID,
									StatusIDvar:         myStatusID,
									ParentJobIDvar:      myParentJobID,
									Identifiervar:       myIdentifier,
									Priorityvar:         myPriority,
									CurrentIterationvar: myCurrentIteration,
									Payloadvar:          myPayload,
									ThreadIDvar:         myThreadID,
									PulseDatevar:        myPulseDate,
									CreatedDatevar:      myCreatedDate,
									UpdatedDatevar:      myUpdatedDate,
								}

								retJobHistory = append(retJobHistory, newJobHistory)
							}
						}

						return err
					})
			}
		},
	})

	return retJobHistory, err
}

// GetJobHistoryByID executes the stored procedure GetJobHistoryByID against the database and returns the read results
func (conn *dbconn) GetJobHistoryByID(_ID string) (domain.JobHistory, error) {
	var err error
	var retJobHistory domain.JobHistory

	conn.Read(&connection.Procedure{
		Proc:       "GetJobHistoryByID",
		Parameters: []interface{}{_ID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myJobID int
							var myConfigID string
							var myStatusID int
							var myParentJobID *string
							var myIdentifier *string
							var myPriority int
							var myCurrentIteration *int
							var myPayload string
							var myThreadID *string
							var myPulseDate *time.Time
							var myCreatedDate time.Time
							var myUpdatedDate *time.Time

							if err = rows.Scan(

								&myID,
								&myJobID,
								&myConfigID,
								&myStatusID,
								&myParentJobID,
								&myIdentifier,
								&myPriority,
								&myCurrentIteration,
								&myPayload,
								&myThreadID,
								&myPulseDate,
								&myCreatedDate,
								&myUpdatedDate,
							); err == nil {

								newJobHistory := &dal.JobHistory{
									IDvar:               myID,
									JobIDvar:            myJobID,
									ConfigIDvar:         myConfigID,
									StatusIDvar:         myStatusID,
									ParentJobIDvar:      myParentJobID,
									Identifiervar:       myIdentifier,
									Priorityvar:         myPriority,
									CurrentIterationvar: myCurrentIteration,
									Payloadvar:          myPayload,
									ThreadIDvar:         myThreadID,
									PulseDatevar:        myPulseDate,
									CreatedDatevar:      myCreatedDate,
									UpdatedDatevar:      myUpdatedDate,
								}

								retJobHistory = newJobHistory
							}
						}

						return err
					})
			}
		},
	})

	return retJobHistory, err
}

// GetJobHistoryLength executes the stored procedure GetJobHistoryLength against the database and returns the read results
func (conn *dbconn) GetJobHistoryLength(_jobid int, _jobconfig string, _status int, _Payload string, _orgid string) (domain.QueryData, error) {
	var err error
	var retQueryData domain.QueryData

	conn.Read(&connection.Procedure{
		Proc:       "GetJobHistoryLength",
		Parameters: []interface{}{_jobid, _jobconfig, _status, _Payload, _orgid},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myLength int

							if err = rows.Scan(

								&myLength,
							); err == nil {

								newQueryData := &dal.QueryData{
									Lengthvar: myLength,
								}

								retQueryData = newQueryData
							}
						}

						return err
					})
			}
		},
	})

	return retQueryData, err
}

// GetJobQueueByStatusID executes the stored procedure GetJobQueueByStatusID against the database and returns the read results
func (conn *dbconn) GetJobQueueByStatusID(_StatusID int) ([]domain.JobHistory, error) {
	var err error
	var retJobHistory = make([]domain.JobHistory, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetJobQueueByStatusID",
		Parameters: []interface{}{_StatusID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myJobID int
							var myConfigID string
							var myStatusID int
							var myParentJobID *string
							var myIdentifier *string
							var myPriority int
							var myPayload string
							var myThreadID *string
							var myPulseDate *time.Time
							var myCreatedDate time.Time
							var myUpdatedDate *time.Time
							var myMaxInstances int

							if err = rows.Scan(

								&myID,
								&myJobID,
								&myConfigID,
								&myStatusID,
								&myParentJobID,
								&myIdentifier,
								&myPriority,
								&myPayload,
								&myThreadID,
								&myPulseDate,
								&myCreatedDate,
								&myUpdatedDate,
								&myMaxInstances,
							); err == nil {

								newJobHistory := &dal.JobHistory{
									IDvar:           myID,
									JobIDvar:        myJobID,
									ConfigIDvar:     myConfigID,
									StatusIDvar:     myStatusID,
									ParentJobIDvar:  myParentJobID,
									Identifiervar:   myIdentifier,
									Priorityvar:     myPriority,
									Payloadvar:      myPayload,
									ThreadIDvar:     myThreadID,
									PulseDatevar:    myPulseDate,
									CreatedDatevar:  myCreatedDate,
									UpdatedDatevar:  myUpdatedDate,
									MaxInstancesvar: myMaxInstances,
								}

								retJobHistory = append(retJobHistory, newJobHistory)
							}
						}

						return err
					})
			}
		},
	})

	return retJobHistory, err
}

// GetJobs executes the stored procedure GetJobs against the database and returns the read results
func (conn *dbconn) GetJobs() ([]domain.JobRegistration, error) {
	var err error
	var retJobRegistration = make([]domain.JobRegistration, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetJobs",
		Parameters: []interface{}{},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID int
							var myGoStruct string
							var myPriority int

							if err = rows.Scan(

								&myID,
								&myGoStruct,
								&myPriority,
							); err == nil {

								newJobRegistration := &dal.JobRegistration{
									IDvar:       myID,
									GoStructvar: myGoStruct,
									Priorityvar: myPriority,
								}

								retJobRegistration = append(retJobRegistration, newJobRegistration)
							}
						}

						return err
					})
			}
		},
	})

	return retJobRegistration, err
}

// GetJobsByStruct executes the stored procedure GetJobsByStruct against the database and returns the read results
func (conn *dbconn) GetJobsByStruct(_Struct string) (domain.JobRegistration, error) {
	var err error
	var retJobRegistration domain.JobRegistration

	conn.Read(&connection.Procedure{
		Proc:       "GetJobsByStruct",
		Parameters: []interface{}{_Struct},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID int
							var myGoStruct string
							var myPriority int
							var myCreatedDate time.Time
							var myCreatedBy string
							var myUpdatedDate *time.Time
							var myUpdatedBy *string

							if err = rows.Scan(

								&myID,
								&myGoStruct,
								&myPriority,
								&myCreatedDate,
								&myCreatedBy,
								&myUpdatedDate,
								&myUpdatedBy,
							); err == nil {

								newJobRegistration := &dal.JobRegistration{
									IDvar:          myID,
									GoStructvar:    myGoStruct,
									Priorityvar:    myPriority,
									CreatedDatevar: myCreatedDate,
									CreatedByvar:   myCreatedBy,
									UpdatedDatevar: myUpdatedDate,
									UpdatedByvar:   myUpdatedBy,
								}

								retJobRegistration = newJobRegistration
							}
						}

						return err
					})
			}
		},
	})

	return retJobRegistration, err
}

// GetLeafOrganizationsForUser executes the stored procedure GetLeafOrganizationsForUser against the database and returns the read results
func (conn *dbconn) GetLeafOrganizationsForUser(_UserID string) ([]domain.Organization, error) {
	var err error
	var retOrganization = make([]domain.Organization, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetLeafOrganizationsForUser",
		Parameters: []interface{}{_UserID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myCode string
							var myDescription *string
							var myTimeZoneOffset float32

							if err = rows.Scan(

								&myID,
								&myCode,
								&myDescription,
								&myTimeZoneOffset,
							); err == nil {

								newOrganization := &dal.Organization{
									IDvar:             myID,
									Codevar:           myCode,
									Descriptionvar:    myDescription,
									TimeZoneOffsetvar: myTimeZoneOffset,
								}

								retOrganization = append(retOrganization, newOrganization)
							}
						}

						return err
					})
			}
		},
	})

	return retOrganization, err
}

// GetLogTypes executes the stored procedure GetLogTypes against the database and returns the read results
func (conn *dbconn) GetLogTypes() ([]domain.LogType, error) {
	var err error
	var retLogType = make([]domain.LogType, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetLogTypes",
		Parameters: []interface{}{},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID int
							var myLogType string
							var myName string

							if err = rows.Scan(

								&myID,
								&myLogType,
								&myName,
							); err == nil {

								newLogType := &dal.LogType{
									IDvar:      myID,
									LogTypevar: myLogType,
									Namevar:    myName,
								}

								retLogType = append(retLogType, newLogType)
							}
						}

						return err
					})
			}
		},
	})

	return retLogType, err
}

// GetLogsByParams executes the stored procedure GetLogsByParams against the database and returns the read results
func (conn *dbconn) GetLogsByParams(_MethodOfDiscovery string, _jobType int, _logType int, _jobHistoryID string, _fromDate time.Time, _toDate time.Time, _OrgID string) ([]domain.DBLog, error) {
	var err error
	var retDBLog = make([]domain.DBLog, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetLogsByParams",
		Parameters: []interface{}{_MethodOfDiscovery, _jobType, _logType, _jobHistoryID, _fromDate, _toDate, _OrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID int
							var myTypeID int
							var myLog string
							var myError string
							var myJobHistoryID string
							var myCreateDate time.Time

							if err = rows.Scan(

								&myID,
								&myTypeID,
								&myLog,
								&myError,
								&myJobHistoryID,
								&myCreateDate,
							); err == nil {

								newDBLog := &dal.DBLog{
									IDvar:           myID,
									TypeIDvar:       myTypeID,
									Logvar:          myLog,
									Errorvar:        myError,
									JobHistoryIDvar: myJobHistoryID,
									CreateDatevar:   myCreateDate,
								}

								retDBLog = append(retDBLog, newDBLog)
							}
						}

						return err
					})
			}
		},
	})

	return retDBLog, err
}

// GetMatchedVulns executes the stored procedure GetMatchedVulns against the database and returns the read results
func (conn *dbconn) GetMatchedVulns() ([]domain.VulnerabilityMatch, error) {
	var err error
	var retVulnerabilityMatch = make([]domain.VulnerabilityMatch, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetMatchedVulns",
		Parameters: []interface{}{},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myFirstID string
							var myFirstTitle string
							var mySecondID string
							var mySecondTitle string
							var myMatchConfidence int
							var myMatchReason string
							var myVulnerabilityID string

							if err = rows.Scan(

								&myFirstID,
								&myFirstTitle,
								&mySecondID,
								&mySecondTitle,
								&myMatchConfidence,
								&myMatchReason,
								&myVulnerabilityID,
							); err == nil {

								newVulnerabilityMatch := &dal.VulnerabilityMatch{
									FirstIDvar:         myFirstID,
									FirstTitlevar:      myFirstTitle,
									SecondIDvar:        mySecondID,
									SecondTitlevar:     mySecondTitle,
									MatchConfidencevar: myMatchConfidence,
									MatchReasonvar:     myMatchReason,
									VulnerabilityIDvar: myVulnerabilityID,
								}

								retVulnerabilityMatch = append(retVulnerabilityMatch, newVulnerabilityMatch)
							}
						}

						return err
					})
			}
		},
	})

	return retVulnerabilityMatch, err
}

// GetOperatingSystemType executes the stored procedure GetOperatingSystemType against the database and returns the read results
func (conn *dbconn) GetOperatingSystemType(_OS string) (domain.OperatingSystemType, error) {
	var err error
	var retOperatingSystemType domain.OperatingSystemType

	conn.Read(&connection.Procedure{
		Proc:       "GetOperatingSystemType",
		Parameters: []interface{}{_OS},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID int
							var myType string
							var myMatch string
							var myPriority int

							if err = rows.Scan(

								&myID,
								&myType,
								&myMatch,
								&myPriority,
							); err == nil {

								newOperatingSystemType := &dal.OperatingSystemType{
									IDvar:       myID,
									Typevar:     myType,
									Matchvar:    myMatch,
									Priorityvar: myPriority,
								}

								retOperatingSystemType = newOperatingSystemType
							}
						}

						return err
					})
			}
		},
	})

	return retOperatingSystemType, err
}

// GetOrganizationByCode executes the stored procedure GetOrganizationByCode against the database and returns the read results
func (conn *dbconn) GetOrganizationByCode(Code string) (domain.Organization, error) {
	var err error
	var retOrganization domain.Organization

	conn.Read(&connection.Procedure{
		Proc:       "GetOrganizationByCode",
		Parameters: []interface{}{Code},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myCode string
							var myDescription *string
							var myTimeZoneOffset float32
							var myCreated time.Time
							var myUpdated *time.Time
							var myPayload string

							if err = rows.Scan(

								&myID,
								&myCode,
								&myDescription,
								&myTimeZoneOffset,
								&myCreated,
								&myUpdated,
								&myPayload,
							); err == nil {

								newOrganization := &dal.Organization{
									IDvar:             myID,
									Codevar:           myCode,
									Descriptionvar:    myDescription,
									TimeZoneOffsetvar: myTimeZoneOffset,
									Createdvar:        myCreated,
									Updatedvar:        myUpdated,
									Payloadvar:        myPayload,
								}

								retOrganization = newOrganization
							}
						}

						return err
					})
			}
		},
	})

	return retOrganization, err
}

// GetOrganizationByID executes the stored procedure GetOrganizationByID against the database and returns the read results
func (conn *dbconn) GetOrganizationByID(ID string) (domain.Organization, error) {
	var err error
	var retOrganization domain.Organization

	conn.Read(&connection.Procedure{
		Proc:       "GetOrganizationByID",
		Parameters: []interface{}{ID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myParentOrgID *string
							var myCode string
							var myDescription *string
							var myTimeZoneOffset float32
							var myCreated time.Time
							var myUpdated *time.Time
							var myPayload string
							var myEncryptionKey *string

							if err = rows.Scan(

								&myID,
								&myParentOrgID,
								&myCode,
								&myDescription,
								&myTimeZoneOffset,
								&myCreated,
								&myUpdated,
								&myPayload,
								&myEncryptionKey,
							); err == nil {

								newOrganization := &dal.Organization{
									IDvar:             myID,
									ParentOrgIDvar:    myParentOrgID,
									Codevar:           myCode,
									Descriptionvar:    myDescription,
									TimeZoneOffsetvar: myTimeZoneOffset,
									Createdvar:        myCreated,
									Updatedvar:        myUpdated,
									Payloadvar:        myPayload,
									EncryptionKeyvar:  myEncryptionKey,
								}

								retOrganization = newOrganization
							}
						}

						return err
					})
			}
		},
	})

	return retOrganization, err
}

// GetOrganizations executes the stored procedure GetOrganizations against the database and returns the read results
func (conn *dbconn) GetOrganizations() ([]domain.Organization, error) {
	var err error
	var retOrganization = make([]domain.Organization, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetOrganizations",
		Parameters: []interface{}{},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myParentOrgID *string
							var myCode string
							var myDescription *string
							var myTimeZoneOffset float32
							var myCreated time.Time
							var myUpdated *time.Time
							var myPayload string

							if err = rows.Scan(

								&myID,
								&myParentOrgID,
								&myCode,
								&myDescription,
								&myTimeZoneOffset,
								&myCreated,
								&myUpdated,
								&myPayload,
							); err == nil {

								newOrganization := &dal.Organization{
									IDvar:             myID,
									ParentOrgIDvar:    myParentOrgID,
									Codevar:           myCode,
									Descriptionvar:    myDescription,
									TimeZoneOffsetvar: myTimeZoneOffset,
									Createdvar:        myCreated,
									Updatedvar:        myUpdated,
									Payloadvar:        myPayload,
								}

								retOrganization = append(retOrganization, newOrganization)
							}
						}

						return err
					})
			}
		},
	})

	return retOrganization, err
}

// GetPendingActiveRescanJob executes the stored procedure GetPendingActiveRescanJob against the database and returns the read results
func (conn *dbconn) GetPendingActiveRescanJob(_OrgID string) ([]domain.JobHistory, error) {
	var err error
	var retJobHistory = make([]domain.JobHistory, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetPendingActiveRescanJob",
		Parameters: []interface{}{_OrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myJobID int
							var myConfigID string
							var myStatusID int
							var myParentJobID *string
							var myIdentifier *string
							var myPriority int
							var myCurrentIteration *int
							var myPayload string
							var myThreadID *string
							var myPulseDate *time.Time
							var myCreatedDate time.Time
							var myUpdatedDate *time.Time

							if err = rows.Scan(

								&myID,
								&myJobID,
								&myConfigID,
								&myStatusID,
								&myParentJobID,
								&myIdentifier,
								&myPriority,
								&myCurrentIteration,
								&myPayload,
								&myThreadID,
								&myPulseDate,
								&myCreatedDate,
								&myUpdatedDate,
							); err == nil {

								newJobHistory := &dal.JobHistory{
									IDvar:               myID,
									JobIDvar:            myJobID,
									ConfigIDvar:         myConfigID,
									StatusIDvar:         myStatusID,
									ParentJobIDvar:      myParentJobID,
									Identifiervar:       myIdentifier,
									Priorityvar:         myPriority,
									CurrentIterationvar: myCurrentIteration,
									Payloadvar:          myPayload,
									ThreadIDvar:         myThreadID,
									PulseDatevar:        myPulseDate,
									CreatedDatevar:      myCreatedDate,
									UpdatedDatevar:      myUpdatedDate,
								}

								retJobHistory = append(retJobHistory, newJobHistory)
							}
						}

						return err
					})
			}
		},
	})

	return retJobHistory, err
}

// GetPermissionByUserOrgID executes the stored procedure GetPermissionByUserOrgID against the database and returns the read results
func (conn *dbconn) GetPermissionByUserOrgID(_UserID string, _OrgID string) (domain.Permission, error) {
	var err error
	var retPermission domain.Permission

	conn.Read(&connection.Procedure{
		Proc:       "GetPermissionByUserOrgID",
		Parameters: []interface{}{_UserID, _OrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myUserID string
							var myOrgID string
							var myAdmin []uint8
							var myManager []uint8
							var myReader []uint8
							var myReporter []uint8

							if err = rows.Scan(

								&myUserID,
								&myOrgID,
								&myAdmin,
								&myManager,
								&myReader,
								&myReporter,
							); err == nil {

								newPermission := &dal.Permission{
									UserIDvar:   myUserID,
									OrgIDvar:    myOrgID,
									Adminvar:    myAdmin[0] > 0 && myAdmin[0] != 48,       // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)
									Managervar:  myManager[0] > 0 && myManager[0] != 48,   // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)
									Readervar:   myReader[0] > 0 && myReader[0] != 48,     // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)
									Reportervar: myReporter[0] > 0 && myReporter[0] != 48, // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)

								}

								retPermission = newPermission
							}
						}

						return err
					})
			}
		},
	})

	return retPermission, err
}

// GetPermissionOfLeafOrgByUserID executes the stored procedure GetPermissionOfLeafOrgByUserID against the database and returns the read results
func (conn *dbconn) GetPermissionOfLeafOrgByUserID(_UserID string) (domain.Permission, error) {
	var err error
	var retPermission domain.Permission

	conn.Read(&connection.Procedure{
		Proc:       "GetPermissionOfLeafOrgByUserID",
		Parameters: []interface{}{_UserID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myUserID string
							var myOrgID string
							var myAdmin []uint8
							var myManager []uint8
							var myReader []uint8
							var myReporter []uint8

							if err = rows.Scan(

								&myUserID,
								&myOrgID,
								&myAdmin,
								&myManager,
								&myReader,
								&myReporter,
							); err == nil {

								newPermission := &dal.Permission{
									UserIDvar:   myUserID,
									OrgIDvar:    myOrgID,
									Adminvar:    myAdmin[0] > 0 && myAdmin[0] != 48,       // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)
									Managervar:  myManager[0] > 0 && myManager[0] != 48,   // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)
									Readervar:   myReader[0] > 0 && myReader[0] != 48,     // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)
									Reportervar: myReporter[0] > 0 && myReporter[0] != 48, // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)

								}

								retPermission = newPermission
							}
						}

						return err
					})
			}
		},
	})

	return retPermission, err
}

// GetRecentlyUpdatedScanSummaries executes the stored procedure GetRecentlyUpdatedScanSummaries against the database and returns the read results
func (conn *dbconn) GetRecentlyUpdatedScanSummaries(_OrgID string) ([]domain.ScanSummary, error) {
	var err error
	var retScanSummary = make([]domain.ScanSummary, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetRecentlyUpdatedScanSummaries",
		Parameters: []interface{}{_OrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var mySourceID string
							var myTemplateID *string
							var myOrgID string
							var mySourceKey *string
							var myScanStatus string
							var myScanClosePayload string
							var myParentJobID string
							var myCreatedDate time.Time
							var myUpdatedDate *time.Time
							var mySource string

							if err = rows.Scan(

								&mySourceID,
								&myTemplateID,
								&myOrgID,
								&mySourceKey,
								&myScanStatus,
								&myScanClosePayload,
								&myParentJobID,
								&myCreatedDate,
								&myUpdatedDate,
								&mySource,
							); err == nil {

								newScanSummary := &dal.ScanSummary{
									SourceIDvar:         mySourceID,
									TemplateIDvar:       myTemplateID,
									OrgIDvar:            myOrgID,
									SourceKeyvar:        mySourceKey,
									ScanStatusvar:       myScanStatus,
									ScanClosePayloadvar: myScanClosePayload,
									ParentJobIDvar:      myParentJobID,
									CreatedDatevar:      myCreatedDate,
									UpdatedDatevar:      myUpdatedDate,
									Sourcevar:           mySource,
								}

								retScanSummary = append(retScanSummary, newScanSummary)
							}
						}

						return err
					})
			}
		},
	})

	return retScanSummary, err
}

// GetScanSummariesBySourceName executes the stored procedure GetScanSummariesBySourceName against the database and returns the read results
func (conn *dbconn) GetScanSummariesBySourceName(_OrgID string, _SourceName string) ([]domain.ScanSummary, error) {
	var err error
	var retScanSummary = make([]domain.ScanSummary, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetScanSummariesBySourceName",
		Parameters: []interface{}{_OrgID, _SourceName},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var mySourceID string
							var myTemplateID *string
							var myOrgID string
							var mySourceKey *string
							var myScanStatus string
							var myParentJobID string
							var myScanClosePayload string
							var myCreatedDate time.Time
							var myUpdatedDate *time.Time

							if err = rows.Scan(

								&mySourceID,
								&myTemplateID,
								&myOrgID,
								&mySourceKey,
								&myScanStatus,
								&myParentJobID,
								&myScanClosePayload,
								&myCreatedDate,
								&myUpdatedDate,
							); err == nil {

								newScanSummary := &dal.ScanSummary{
									SourceIDvar:         mySourceID,
									TemplateIDvar:       myTemplateID,
									OrgIDvar:            myOrgID,
									SourceKeyvar:        mySourceKey,
									ScanStatusvar:       myScanStatus,
									ParentJobIDvar:      myParentJobID,
									ScanClosePayloadvar: myScanClosePayload,
									CreatedDatevar:      myCreatedDate,
									UpdatedDatevar:      myUpdatedDate,
								}

								retScanSummary = append(retScanSummary, newScanSummary)
							}
						}

						return err
					})
			}
		},
	})

	return retScanSummary, err
}

// GetScanSummary executes the stored procedure GetScanSummary against the database and returns the read results
func (conn *dbconn) GetScanSummary(_SourceID string, _OrgID string, _ScanID string) (domain.ScanSummary, error) {
	var err error
	var retScanSummary domain.ScanSummary

	conn.Read(&connection.Procedure{
		Proc:       "GetScanSummary",
		Parameters: []interface{}{_SourceID, _OrgID, _ScanID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var mySourceID string
							var myTemplateID *string
							var myOrgID string
							var mySourceKey *string
							var myScanStatus string
							var myParentJobID string
							var myScanClosePayload string
							var myCreatedDate time.Time
							var myUpdatedDate *time.Time

							if err = rows.Scan(

								&mySourceID,
								&myTemplateID,
								&myOrgID,
								&mySourceKey,
								&myScanStatus,
								&myParentJobID,
								&myScanClosePayload,
								&myCreatedDate,
								&myUpdatedDate,
							); err == nil {

								newScanSummary := &dal.ScanSummary{
									SourceIDvar:         mySourceID,
									TemplateIDvar:       myTemplateID,
									OrgIDvar:            myOrgID,
									SourceKeyvar:        mySourceKey,
									ScanStatusvar:       myScanStatus,
									ParentJobIDvar:      myParentJobID,
									ScanClosePayloadvar: myScanClosePayload,
									CreatedDatevar:      myCreatedDate,
									UpdatedDatevar:      myUpdatedDate,
								}

								retScanSummary = newScanSummary
							}
						}

						return err
					})
			}
		},
	})

	return retScanSummary, err
}

// GetScanSummaryBySourceKey executes the stored procedure GetScanSummaryBySourceKey against the database and returns the read results
func (conn *dbconn) GetScanSummaryBySourceKey(_SourceKey string) (domain.ScanSummary, error) {
	var err error
	var retScanSummary domain.ScanSummary

	conn.Read(&connection.Procedure{
		Proc:       "GetScanSummaryBySourceKey",
		Parameters: []interface{}{_SourceKey},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var mySourceID string
							var myTemplateID *string
							var myOrgID string
							var mySourceKey *string
							var myScanStatus string
							var myScanClosePayload string
							var myCreatedDate time.Time
							var myUpdatedDate *time.Time

							if err = rows.Scan(

								&mySourceID,
								&myTemplateID,
								&myOrgID,
								&mySourceKey,
								&myScanStatus,
								&myScanClosePayload,
								&myCreatedDate,
								&myUpdatedDate,
							); err == nil {

								newScanSummary := &dal.ScanSummary{
									SourceIDvar:         mySourceID,
									TemplateIDvar:       myTemplateID,
									OrgIDvar:            myOrgID,
									SourceKeyvar:        mySourceKey,
									ScanStatusvar:       myScanStatus,
									ScanClosePayloadvar: myScanClosePayload,
									CreatedDatevar:      myCreatedDate,
									UpdatedDatevar:      myUpdatedDate,
								}

								retScanSummary = newScanSummary
							}
						}

						return err
					})
			}
		},
	})

	return retScanSummary, err
}

// GetScheduledJobsToStart executes the stored procedure GetScheduledJobsToStart against the database and returns the read results
func (conn *dbconn) GetScheduledJobsToStart(_LastChecked time.Time) ([]domain.JobSchedule, error) {
	var err error
	var retJobSchedule = make([]domain.JobSchedule, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetScheduledJobsToStart",
		Parameters: []interface{}{_LastChecked},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myConfigID string
							var myPayload *string

							if err = rows.Scan(

								&myID,
								&myConfigID,
								&myPayload,
							); err == nil {

								newJobSchedule := &dal.JobSchedule{
									IDvar:       myID,
									ConfigIDvar: myConfigID,
									Payloadvar:  myPayload,
								}

								retJobSchedule = append(retJobSchedule, newJobSchedule)
							}
						}

						return err
					})
			}
		},
	})

	return retJobSchedule, err
}

// GetSessionByToken executes the stored procedure GetSessionByToken against the database and returns the read results
func (conn *dbconn) GetSessionByToken(_SessionKey string) (domain.Session, error) {
	var err error
	var retSession domain.Session

	conn.Read(&connection.Procedure{
		Proc:       "GetSessionByToken",
		Parameters: []interface{}{_SessionKey},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myUserID string
							var myOrgID string
							var mySessionKey string
							var myIsDisabled []uint8

							if err = rows.Scan(

								&myUserID,
								&myOrgID,
								&mySessionKey,
								&myIsDisabled,
							); err == nil {

								newSession := &dal.Session{
									UserIDvar:     myUserID,
									OrgIDvar:      myOrgID,
									SessionKeyvar: mySessionKey,
									IsDisabledvar: myIsDisabled[0] > 0 && myIsDisabled[0] != 48, // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)

								}

								retSession = newSession
							}
						}

						return err
					})
			}
		},
	})

	return retSession, err
}

// GetSourceByID executes the stored procedure GetSourceByID against the database and returns the read results
func (conn *dbconn) GetSourceByID(_ID string) (domain.Source, error) {
	var err error
	var retSource domain.Source

	conn.Read(&connection.Procedure{
		Proc:       "GetSourceByID",
		Parameters: []interface{}{_ID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var mySourceTypeID int
							var mySource string

							if err = rows.Scan(

								&myID,
								&mySourceTypeID,
								&mySource,
							); err == nil {

								newSource := &dal.Source{
									IDvar:           myID,
									SourceTypeIDvar: mySourceTypeID,
									Sourcevar:       mySource,
								}

								retSource = newSource
							}
						}

						return err
					})
			}
		},
	})

	return retSource, err
}

// GetSourceByName executes the stored procedure GetSourceByName against the database and returns the read results
func (conn *dbconn) GetSourceByName(_Source string) (domain.Source, error) {
	var err error
	var retSource domain.Source

	conn.Read(&connection.Procedure{
		Proc:       "GetSourceByName",
		Parameters: []interface{}{_Source},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var mySourceTypeID int
							var mySource string
							var myDBCreatedDate time.Time
							var myDBUpdatedDate *time.Time

							if err = rows.Scan(

								&myID,
								&mySourceTypeID,
								&mySource,
								&myDBCreatedDate,
								&myDBUpdatedDate,
							); err == nil {

								newSource := &dal.Source{
									IDvar:            myID,
									SourceTypeIDvar:  mySourceTypeID,
									Sourcevar:        mySource,
									DBCreatedDatevar: myDBCreatedDate,
									DBUpdatedDatevar: myDBUpdatedDate,
								}

								retSource = newSource
							}
						}

						return err
					})
			}
		},
	})

	return retSource, err
}

// GetSourceConfigByID executes the stored procedure GetSourceConfigByID against the database and returns the read results
func (conn *dbconn) GetSourceConfigByID(_ID string) (domain.SourceConfig, error) {
	var err error
	var retSourceConfig domain.SourceConfig

	conn.Read(&connection.Procedure{
		Proc:       "GetSourceConfigByID",
		Parameters: []interface{}{_ID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var mySourceID string
							var mySource string
							var myAddress string
							var myPort string
							var myAuthInfo string
							var myPayload *string
							var myOrganizationID string
							var myDBCreatedDate time.Time
							var myDBUpdatedDate *time.Time

							if err = rows.Scan(

								&myID,
								&mySourceID,
								&mySource,
								&myAddress,
								&myPort,
								&myAuthInfo,
								&myPayload,
								&myOrganizationID,
								&myDBCreatedDate,
								&myDBUpdatedDate,
							); err == nil {

								newSourceConfig := &dal.SourceConfig{
									IDvar:             myID,
									SourceIDvar:       mySourceID,
									Sourcevar:         mySource,
									Addressvar:        myAddress,
									Portvar:           myPort,
									AuthInfovar:       myAuthInfo,
									Payloadvar:        myPayload,
									OrganizationIDvar: myOrganizationID,
									DBCreatedDatevar:  myDBCreatedDate,
									DBUpdatedDatevar:  myDBUpdatedDate,
								}

								retSourceConfig = newSourceConfig
							}
						}

						return err
					})
			}
		},
	})

	return retSourceConfig, err
}

// GetSourceConfigByNameOrg executes the stored procedure GetSourceConfigByNameOrg against the database and returns the read results
func (conn *dbconn) GetSourceConfigByNameOrg(_Source string, _OrgID string) ([]domain.SourceConfig, error) {
	var err error
	var retSourceConfig = make([]domain.SourceConfig, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetSourceConfigByNameOrg",
		Parameters: []interface{}{_Source, _OrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var mySourceID string
							var mySource string
							var myOrganizationID string
							var myAddress string
							var myPort string
							var myAuthInfo string
							var myPayload *string
							var myDBCreatedDate time.Time
							var myDBUpdatedDate *time.Time

							if err = rows.Scan(

								&myID,
								&mySourceID,
								&mySource,
								&myOrganizationID,
								&myAddress,
								&myPort,
								&myAuthInfo,
								&myPayload,
								&myDBCreatedDate,
								&myDBUpdatedDate,
							); err == nil {

								newSourceConfig := &dal.SourceConfig{
									IDvar:             myID,
									SourceIDvar:       mySourceID,
									Sourcevar:         mySource,
									OrganizationIDvar: myOrganizationID,
									Addressvar:        myAddress,
									Portvar:           myPort,
									AuthInfovar:       myAuthInfo,
									Payloadvar:        myPayload,
									DBCreatedDatevar:  myDBCreatedDate,
									DBUpdatedDatevar:  myDBUpdatedDate,
								}

								retSourceConfig = append(retSourceConfig, newSourceConfig)
							}
						}

						return err
					})
			}
		},
	})

	return retSourceConfig, err
}

// GetSourceConfigByOrgID executes the stored procedure GetSourceConfigByOrgID against the database and returns the read results
func (conn *dbconn) GetSourceConfigByOrgID(_OrgID string) ([]domain.SourceConfig, error) {
	var err error
	var retSourceConfig = make([]domain.SourceConfig, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetSourceConfigByOrgID",
		Parameters: []interface{}{_OrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var mySourceID string
							var mySource string
							var myOrganizationID string
							var myAddress string
							var myPort string
							var myAuthInfo string
							var myPayload *string
							var myDBCreatedDate time.Time
							var myDBUpdatedDate *time.Time

							if err = rows.Scan(

								&myID,
								&mySourceID,
								&mySource,
								&myOrganizationID,
								&myAddress,
								&myPort,
								&myAuthInfo,
								&myPayload,
								&myDBCreatedDate,
								&myDBUpdatedDate,
							); err == nil {

								newSourceConfig := &dal.SourceConfig{
									IDvar:             myID,
									SourceIDvar:       mySourceID,
									Sourcevar:         mySource,
									OrganizationIDvar: myOrganizationID,
									Addressvar:        myAddress,
									Portvar:           myPort,
									AuthInfovar:       myAuthInfo,
									Payloadvar:        myPayload,
									DBCreatedDatevar:  myDBCreatedDate,
									DBUpdatedDatevar:  myDBUpdatedDate,
								}

								retSourceConfig = append(retSourceConfig, newSourceConfig)
							}
						}

						return err
					})
			}
		},
	})

	return retSourceConfig, err
}

// GetSourceConfigBySourceID executes the stored procedure GetSourceConfigBySourceID against the database and returns the read results
func (conn *dbconn) GetSourceConfigBySourceID(_OrgID string, _SourceID string) ([]domain.SourceConfig, error) {
	var err error
	var retSourceConfig = make([]domain.SourceConfig, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetSourceConfigBySourceID",
		Parameters: []interface{}{_OrgID, _SourceID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var mySourceID string
							var myOrganizationID string
							var myAddress string
							var myPort string
							var myAuthInfo string
							var myPayload *string
							var myDBCreatedDate time.Time
							var myDBUpdatedDate *time.Time

							if err = rows.Scan(

								&myID,
								&mySourceID,
								&myOrganizationID,
								&myAddress,
								&myPort,
								&myAuthInfo,
								&myPayload,
								&myDBCreatedDate,
								&myDBUpdatedDate,
							); err == nil {

								newSourceConfig := &dal.SourceConfig{
									IDvar:             myID,
									SourceIDvar:       mySourceID,
									OrganizationIDvar: myOrganizationID,
									Addressvar:        myAddress,
									Portvar:           myPort,
									AuthInfovar:       myAuthInfo,
									Payloadvar:        myPayload,
									DBCreatedDatevar:  myDBCreatedDate,
									DBUpdatedDatevar:  myDBUpdatedDate,
								}

								retSourceConfig = append(retSourceConfig, newSourceConfig)
							}
						}

						return err
					})
			}
		},
	})

	return retSourceConfig, err
}

// GetSourceInsByJobID executes the stored procedure GetSourceInsByJobID against the database and returns the read results
func (conn *dbconn) GetSourceInsByJobID(inJob int, inOrgID string) ([]domain.SourceConfig, error) {
	var err error
	var retSourceConfig = make([]domain.SourceConfig, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetSourceInsByJobID",
		Parameters: []interface{}{inJob, inOrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myAddress string
							var mySource string
							var myPort string

							if err = rows.Scan(

								&myID,
								&myAddress,
								&mySource,
								&myPort,
							); err == nil {

								newSourceConfig := &dal.SourceConfig{
									IDvar:      myID,
									Addressvar: myAddress,
									Sourcevar:  mySource,
									Portvar:    myPort,
								}

								retSourceConfig = append(retSourceConfig, newSourceConfig)
							}
						}

						return err
					})
			}
		},
	})

	return retSourceConfig, err
}

// GetSourceOauthByOrgURL executes the stored procedure GetSourceOauthByOrgURL against the database and returns the read results
func (conn *dbconn) GetSourceOauthByOrgURL(_URL string, _OrgID string) (domain.SourceConfig, error) {
	var err error
	var retSourceConfig domain.SourceConfig

	conn.Read(&connection.Procedure{
		Proc:       "GetSourceOauthByOrgURL",
		Parameters: []interface{}{_URL, _OrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var mySourceID string
							var mySource string
							var myAddress string
							var myPort string
							var myAuthInfo string
							var myPayload *string
							var myOrganizationID string
							var myDBCreatedDate time.Time
							var myDBUpdatedDate *time.Time

							if err = rows.Scan(

								&myID,
								&mySourceID,
								&mySource,
								&myAddress,
								&myPort,
								&myAuthInfo,
								&myPayload,
								&myOrganizationID,
								&myDBCreatedDate,
								&myDBUpdatedDate,
							); err == nil {

								newSourceConfig := &dal.SourceConfig{
									IDvar:             myID,
									SourceIDvar:       mySourceID,
									Sourcevar:         mySource,
									Addressvar:        myAddress,
									Portvar:           myPort,
									AuthInfovar:       myAuthInfo,
									Payloadvar:        myPayload,
									OrganizationIDvar: myOrganizationID,
									DBCreatedDatevar:  myDBCreatedDate,
									DBUpdatedDatevar:  myDBUpdatedDate,
								}

								retSourceConfig = newSourceConfig
							}
						}

						return err
					})
			}
		},
	})

	return retSourceConfig, err
}

// GetSourceOauthByURL executes the stored procedure GetSourceOauthByURL against the database and returns the read results
func (conn *dbconn) GetSourceOauthByURL(_URL string) (domain.SourceConfig, error) {
	var err error
	var retSourceConfig domain.SourceConfig

	conn.Read(&connection.Procedure{
		Proc:       "GetSourceOauthByURL",
		Parameters: []interface{}{_URL},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var mySourceID string
							var mySource string
							var myAddress string
							var myPort string
							var myAuthInfo string
							var myPayload *string
							var myOrganizationID string
							var myDBCreatedDate time.Time
							var myDBUpdatedDate *time.Time

							if err = rows.Scan(

								&myID,
								&mySourceID,
								&mySource,
								&myAddress,
								&myPort,
								&myAuthInfo,
								&myPayload,
								&myOrganizationID,
								&myDBCreatedDate,
								&myDBUpdatedDate,
							); err == nil {

								newSourceConfig := &dal.SourceConfig{
									IDvar:             myID,
									SourceIDvar:       mySourceID,
									Sourcevar:         mySource,
									Addressvar:        myAddress,
									Portvar:           myPort,
									AuthInfovar:       myAuthInfo,
									Payloadvar:        myPayload,
									OrganizationIDvar: myOrganizationID,
									DBCreatedDatevar:  myDBCreatedDate,
									DBUpdatedDatevar:  myDBUpdatedDate,
								}

								retSourceConfig = newSourceConfig
							}
						}

						return err
					})
			}
		},
	})

	return retSourceConfig, err
}

// GetSourceOutsByJobID executes the stored procedure GetSourceOutsByJobID against the database and returns the read results
func (conn *dbconn) GetSourceOutsByJobID(inJob int, inOrgID string) ([]domain.SourceConfig, error) {
	var err error
	var retSourceConfig = make([]domain.SourceConfig, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetSourceOutsByJobID",
		Parameters: []interface{}{inJob, inOrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myAddress string
							var mySource string
							var myPort string

							if err = rows.Scan(

								&myID,
								&myAddress,
								&mySource,
								&myPort,
							); err == nil {

								newSourceConfig := &dal.SourceConfig{
									IDvar:      myID,
									Addressvar: myAddress,
									Sourcevar:  mySource,
									Portvar:    myPort,
								}

								retSourceConfig = append(retSourceConfig, newSourceConfig)
							}
						}

						return err
					})
			}
		},
	})

	return retSourceConfig, err
}

// GetSources executes the stored procedure GetSources against the database and returns the read results
func (conn *dbconn) GetSources() ([]domain.Source, error) {
	var err error
	var retSource = make([]domain.Source, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetSources",
		Parameters: []interface{}{},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var mySourceTypeID int
							var mySource string
							var myDBCreatedDate time.Time
							var myDBUpdatedDate *time.Time

							if err = rows.Scan(

								&myID,
								&mySourceTypeID,
								&mySource,
								&myDBCreatedDate,
								&myDBUpdatedDate,
							); err == nil {

								newSource := &dal.Source{
									IDvar:            myID,
									SourceTypeIDvar:  mySourceTypeID,
									Sourcevar:        mySource,
									DBCreatedDatevar: myDBCreatedDate,
									DBUpdatedDatevar: myDBUpdatedDate,
								}

								retSource = append(retSource, newSource)
							}
						}

						return err
					})
			}
		},
	})

	return retSource, err
}

// GetTagByDeviceAndTagKey executes the stored procedure GetTagByDeviceAndTagKey against the database and returns the read results
func (conn *dbconn) GetTagByDeviceAndTagKey(_DeviceID string, _TagKeyID string) (domain.Tag, error) {
	var err error
	var retTag domain.Tag

	conn.Read(&connection.Procedure{
		Proc:       "GetTagByDeviceAndTagKey",
		Parameters: []interface{}{_DeviceID, _TagKeyID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myDeviceID string
							var myTagKeyID int
							var myValue string

							if err = rows.Scan(

								&myID,
								&myDeviceID,
								&myTagKeyID,
								&myValue,
							); err == nil {

								newTag := &dal.Tag{
									IDvar:       myID,
									DeviceIDvar: myDeviceID,
									TagKeyIDvar: myTagKeyID,
									Valuevar:    myValue,
								}

								retTag = newTag
							}
						}

						return err
					})
			}
		},
	})

	return retTag, err
}

// GetTagKeyByID executes the stored procedure GetTagKeyByID against the database and returns the read results
func (conn *dbconn) GetTagKeyByID(_ID string) (domain.TagKey, error) {
	var err error
	var retTagKey domain.TagKey

	conn.Read(&connection.Procedure{
		Proc:       "GetTagKeyByID",
		Parameters: []interface{}{_ID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myKeyValue string

							if err = rows.Scan(

								&myID,
								&myKeyValue,
							); err == nil {

								newTagKey := &dal.TagKey{
									IDvar:       myID,
									KeyValuevar: myKeyValue,
								}

								retTagKey = newTagKey
							}
						}

						return err
					})
			}
		},
	})

	return retTagKey, err
}

// GetTagKeyByKey executes the stored procedure GetTagKeyByKey against the database and returns the read results
func (conn *dbconn) GetTagKeyByKey(_KeyValue string) (domain.TagKey, error) {
	var err error
	var retTagKey domain.TagKey

	conn.Read(&connection.Procedure{
		Proc:       "GetTagKeyByKey",
		Parameters: []interface{}{_KeyValue},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myKeyValue string

							if err = rows.Scan(

								&myID,
								&myKeyValue,
							); err == nil {

								newTagKey := &dal.TagKey{
									IDvar:       myID,
									KeyValuevar: myKeyValue,
								}

								retTagKey = newTagKey
							}
						}

						return err
					})
			}
		},
	})

	return retTagKey, err
}

// GetTagMapsByOrg executes the stored procedure GetTagMapsByOrg against the database and returns the read results
func (conn *dbconn) GetTagMapsByOrg(_OrganizationID string) ([]domain.TagMap, error) {
	var err error
	var retTagMap = make([]domain.TagMap, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetTagMapsByOrg",
		Parameters: []interface{}{_OrganizationID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myTicketingSourceID string
							var myTicketingTag string
							var myCloudSourceID string
							var myCloudTag string
							var myOptions string

							if err = rows.Scan(

								&myID,
								&myTicketingSourceID,
								&myTicketingTag,
								&myCloudSourceID,
								&myCloudTag,
								&myOptions,
							); err == nil {

								newTagMap := &dal.TagMap{
									IDvar:                myID,
									TicketingSourceIDvar: myTicketingSourceID,
									TicketingTagvar:      myTicketingTag,
									CloudSourceIDvar:     myCloudSourceID,
									CloudTagvar:          myCloudTag,
									Optionsvar:           myOptions,
								}

								retTagMap = append(retTagMap, newTagMap)
							}
						}

						return err
					})
			}
		},
	})

	return retTagMap, err
}

// GetTagMapsByOrgCloudSourceID executes the stored procedure GetTagMapsByOrgCloudSourceID against the database and returns the read results
func (conn *dbconn) GetTagMapsByOrgCloudSourceID(_CloudID string, _OrganizationID string) ([]domain.TagMap, error) {
	var err error
	var retTagMap = make([]domain.TagMap, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetTagMapsByOrgCloudSourceID",
		Parameters: []interface{}{_CloudID, _OrganizationID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myTicketingSourceID string
							var myTicketingTag string
							var myCloudSourceID string
							var myCloudTag string
							var myOptions string

							if err = rows.Scan(

								&myID,
								&myTicketingSourceID,
								&myTicketingTag,
								&myCloudSourceID,
								&myCloudTag,
								&myOptions,
							); err == nil {

								newTagMap := &dal.TagMap{
									IDvar:                myID,
									TicketingSourceIDvar: myTicketingSourceID,
									TicketingTagvar:      myTicketingTag,
									CloudSourceIDvar:     myCloudSourceID,
									CloudTagvar:          myCloudTag,
									Optionsvar:           myOptions,
								}

								retTagMap = append(retTagMap, newTagMap)
							}
						}

						return err
					})
			}
		},
	})

	return retTagMap, err
}

// GetTagsForDevice executes the stored procedure GetTagsForDevice against the database and returns the read results
func (conn *dbconn) GetTagsForDevice(_DeviceID string) ([]domain.Tag, error) {
	var err error
	var retTag = make([]domain.Tag, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetTagsForDevice",
		Parameters: []interface{}{_DeviceID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myDeviceID string
							var myTagKeyID int
							var myValue string

							if err = rows.Scan(

								&myID,
								&myDeviceID,
								&myTagKeyID,
								&myValue,
							); err == nil {

								newTag := &dal.Tag{
									IDvar:       myID,
									DeviceIDvar: myDeviceID,
									TagKeyIDvar: myTagKeyID,
									Valuevar:    myValue,
								}

								retTag = append(retTag, newTag)
							}
						}

						return err
					})
			}
		},
	})

	return retTag, err
}

// GetTicketByDetectionID executes the stored procedure GetTicketByDetectionID against the database and returns the read results
func (conn *dbconn) GetTicketByDetectionID(inDetectionID string, _OrgID string) (domain.TicketSummary, error) {
	var err error
	var retTicketSummary domain.TicketSummary

	conn.Read(&connection.Procedure{
		Proc:       "GetTicketByDetectionID",
		Parameters: []interface{}{inDetectionID, _OrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myTitle string
							var myStatus string
							var myDetectionID string
							var myOrganizationID string
							var myUpdatedDate *time.Time
							var myResolutionDate *time.Time
							var myDueDate time.Time

							if err = rows.Scan(

								&myTitle,
								&myStatus,
								&myDetectionID,
								&myOrganizationID,
								&myUpdatedDate,
								&myResolutionDate,
								&myDueDate,
							); err == nil {

								newTicketSummary := &dal.TicketSummary{
									Titlevar:          myTitle,
									Statusvar:         myStatus,
									DetectionIDvar:    myDetectionID,
									OrganizationIDvar: myOrganizationID,
									UpdatedDatevar:    myUpdatedDate,
									ResolutionDatevar: myResolutionDate,
									DueDatevar:        myDueDate,
								}

								retTicketSummary = newTicketSummary
							}
						}

						return err
					})
			}
		},
	})

	return retTicketSummary, err
}

// GetTicketByDeviceIDVulnID executes the stored procedure GetTicketByDeviceIDVulnID against the database and returns the read results
func (conn *dbconn) GetTicketByDeviceIDVulnID(inDeviceID string, inVulnID string, inPort int, inProtocol string, inOrgID string) (domain.TicketSummary, error) {
	var err error
	var retTicketSummary domain.TicketSummary

	conn.Read(&connection.Procedure{
		Proc:       "GetTicketByDeviceIDVulnID",
		Parameters: []interface{}{inDeviceID, inVulnID, inPort, inProtocol, inOrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myTitle string
							var myStatus string
							var myDetectionID string
							var myOrganizationID string
							var myUpdatedDate *time.Time
							var myResolutionDate *time.Time
							var myDueDate time.Time

							if err = rows.Scan(

								&myTitle,
								&myStatus,
								&myDetectionID,
								&myOrganizationID,
								&myUpdatedDate,
								&myResolutionDate,
								&myDueDate,
							); err == nil {

								newTicketSummary := &dal.TicketSummary{
									Titlevar:          myTitle,
									Statusvar:         myStatus,
									DetectionIDvar:    myDetectionID,
									OrganizationIDvar: myOrganizationID,
									UpdatedDatevar:    myUpdatedDate,
									ResolutionDatevar: myResolutionDate,
									DueDatevar:        myDueDate,
								}

								retTicketSummary = newTicketSummary
							}
						}

						return err
					})
			}
		},
	})

	return retTicketSummary, err
}

// GetTicketByIPGroupIDVulnID executes the stored procedure GetTicketByIPGroupIDVulnID against the database and returns the read results
func (conn *dbconn) GetTicketByIPGroupIDVulnID(inIP string, inGroupID string, inVulnID string, inPort int, inProtocol string, inOrgID string) (domain.TicketSummary, error) {
	var err error
	var retTicketSummary domain.TicketSummary

	conn.Read(&connection.Procedure{
		Proc:       "GetTicketByIPGroupIDVulnID",
		Parameters: []interface{}{inIP, inGroupID, inVulnID, inPort, inProtocol, inOrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myTitle string
							var myStatus string
							var myDetectionID string
							var myOrganizationID string
							var myUpdatedDate *time.Time
							var myResolutionDate *time.Time
							var myDueDate time.Time

							if err = rows.Scan(

								&myTitle,
								&myStatus,
								&myDetectionID,
								&myOrganizationID,
								&myUpdatedDate,
								&myResolutionDate,
								&myDueDate,
							); err == nil {

								newTicketSummary := &dal.TicketSummary{
									Titlevar:          myTitle,
									Statusvar:         myStatus,
									DetectionIDvar:    myDetectionID,
									OrganizationIDvar: myOrganizationID,
									UpdatedDatevar:    myUpdatedDate,
									ResolutionDatevar: myResolutionDate,
									DueDatevar:        myDueDate,
								}

								retTicketSummary = newTicketSummary
							}
						}

						return err
					})
			}
		},
	})

	return retTicketSummary, err
}

// GetTicketByTitle executes the stored procedure GetTicketByTitle against the database and returns the read results
func (conn *dbconn) GetTicketByTitle(_Title string, _OrgID string) (domain.TicketSummary, error) {
	var err error
	var retTicketSummary domain.TicketSummary

	conn.Read(&connection.Procedure{
		Proc:       "GetTicketByTitle",
		Parameters: []interface{}{_Title, _OrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myTitle string
							var myStatus string
							var myDetectionID string
							var myOrganizationID string
							var myUpdatedDate *time.Time
							var myResolutionDate *time.Time
							var myDueDate time.Time

							if err = rows.Scan(

								&myTitle,
								&myStatus,
								&myDetectionID,
								&myOrganizationID,
								&myUpdatedDate,
								&myResolutionDate,
								&myDueDate,
							); err == nil {

								newTicketSummary := &dal.TicketSummary{
									Titlevar:          myTitle,
									Statusvar:         myStatus,
									DetectionIDvar:    myDetectionID,
									OrganizationIDvar: myOrganizationID,
									UpdatedDatevar:    myUpdatedDate,
									ResolutionDatevar: myResolutionDate,
									DueDatevar:        myDueDate,
								}

								retTicketSummary = newTicketSummary
							}
						}

						return err
					})
			}
		},
	})

	return retTicketSummary, err
}

// GetTicketCountByStatus executes the stored procedure GetTicketCountByStatus against the database and returns the read results
func (conn *dbconn) GetTicketCountByStatus(inStatus string, inOrgID string) (domain.QueryData, error) {
	var err error
	var retQueryData domain.QueryData

	conn.Read(&connection.Procedure{
		Proc:       "GetTicketCountByStatus",
		Parameters: []interface{}{inStatus, inOrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myLength int

							if err = rows.Scan(

								&myLength,
							); err == nil {

								newQueryData := &dal.QueryData{
									Lengthvar: myLength,
								}

								retQueryData = newQueryData
							}
						}

						return err
					})
			}
		},
	})

	return retQueryData, err
}

// GetTicketCreatedAfter executes the stored procedure GetTicketCreatedAfter against the database and returns the read results
func (conn *dbconn) GetTicketCreatedAfter(_UpperCVSS float32, _LowerCVSS float32, _CreatedAfter time.Time, _OrgID string) ([]domain.TicketSummary, error) {
	var err error
	var retTicketSummary = make([]domain.TicketSummary, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetTicketCreatedAfter",
		Parameters: []interface{}{_UpperCVSS, _LowerCVSS, _CreatedAfter, _OrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myTitle string
							var myStatus string
							var myDetectionID string
							var myOrganizationID string
							var myCreatedDate *time.Time
							var myUpdatedDate *time.Time
							var myResolutionDate *time.Time
							var myDueDate time.Time

							if err = rows.Scan(

								&myTitle,
								&myStatus,
								&myDetectionID,
								&myOrganizationID,
								&myCreatedDate,
								&myUpdatedDate,
								&myResolutionDate,
								&myDueDate,
							); err == nil {

								newTicketSummary := &dal.TicketSummary{
									Titlevar:          myTitle,
									Statusvar:         myStatus,
									DetectionIDvar:    myDetectionID,
									OrganizationIDvar: myOrganizationID,
									CreatedDatevar:    myCreatedDate,
									UpdatedDatevar:    myUpdatedDate,
									ResolutionDatevar: myResolutionDate,
									DueDatevar:        myDueDate,
								}

								retTicketSummary = append(retTicketSummary, newTicketSummary)
							}
						}

						return err
					})
			}
		},
	})

	return retTicketSummary, err
}

// GetTicketTrackingMethod executes the stored procedure GetTicketTrackingMethod against the database and returns the read results
func (conn *dbconn) GetTicketTrackingMethod(_Title string, _OrgID string) (domain.KeyValue, error) {
	var err error
	var retKeyValue domain.KeyValue

	conn.Read(&connection.Procedure{
		Proc:       "GetTicketTrackingMethod",
		Parameters: []interface{}{_Title, _OrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myKey string
							var myValue string

							if err = rows.Scan(

								&myKey,
								&myValue,
							); err == nil {

								newKeyValue := &dal.KeyValue{
									Keyvar:   myKey,
									Valuevar: myValue,
								}

								retKeyValue = newKeyValue
							}
						}

						return err
					})
			}
		},
	})

	return retKeyValue, err
}

// GetUnfinishedScanSummariesBySourceConfigOrgID executes the stored procedure GetUnfinishedScanSummariesBySourceConfigOrgID against the database and returns the read results
func (conn *dbconn) GetUnfinishedScanSummariesBySourceConfigOrgID(_ScannerSourceConfigID string, _OrgID string) ([]domain.ScanSummary, error) {
	var err error
	var retScanSummary = make([]domain.ScanSummary, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetUnfinishedScanSummariesBySourceConfigOrgID",
		Parameters: []interface{}{_ScannerSourceConfigID, _OrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var mySourceID string
							var myTemplateID *string
							var myOrgID string
							var mySourceKey *string
							var myScanStatus string
							var myScanClosePayload string
							var myParentJobID string
							var myCreatedDate time.Time
							var myUpdatedDate *time.Time

							if err = rows.Scan(

								&mySourceID,
								&myTemplateID,
								&myOrgID,
								&mySourceKey,
								&myScanStatus,
								&myScanClosePayload,
								&myParentJobID,
								&myCreatedDate,
								&myUpdatedDate,
							); err == nil {

								newScanSummary := &dal.ScanSummary{
									SourceIDvar:         mySourceID,
									TemplateIDvar:       myTemplateID,
									OrgIDvar:            myOrgID,
									SourceKeyvar:        mySourceKey,
									ScanStatusvar:       myScanStatus,
									ScanClosePayloadvar: myScanClosePayload,
									ParentJobIDvar:      myParentJobID,
									CreatedDatevar:      myCreatedDate,
									UpdatedDatevar:      myUpdatedDate,
								}

								retScanSummary = append(retScanSummary, newScanSummary)
							}
						}

						return err
					})
			}
		},
	})

	return retScanSummary, err
}

// GetUnfinishedScanSummariesBySourceOrgID executes the stored procedure GetUnfinishedScanSummariesBySourceOrgID against the database and returns the read results
func (conn *dbconn) GetUnfinishedScanSummariesBySourceOrgID(_SourceID string, _OrgID string) ([]domain.ScanSummary, error) {
	var err error
	var retScanSummary = make([]domain.ScanSummary, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetUnfinishedScanSummariesBySourceOrgID",
		Parameters: []interface{}{_SourceID, _OrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var mySourceID string
							var myTemplateID *string
							var myOrgID string
							var mySourceKey *string
							var myScanStatus string
							var myScanClosePayload string
							var myParentJobID string
							var myCreatedDate time.Time
							var myUpdatedDate *time.Time

							if err = rows.Scan(

								&mySourceID,
								&myTemplateID,
								&myOrgID,
								&mySourceKey,
								&myScanStatus,
								&myScanClosePayload,
								&myParentJobID,
								&myCreatedDate,
								&myUpdatedDate,
							); err == nil {

								newScanSummary := &dal.ScanSummary{
									SourceIDvar:         mySourceID,
									TemplateIDvar:       myTemplateID,
									OrgIDvar:            myOrgID,
									SourceKeyvar:        mySourceKey,
									ScanStatusvar:       myScanStatus,
									ScanClosePayloadvar: myScanClosePayload,
									ParentJobIDvar:      myParentJobID,
									CreatedDatevar:      myCreatedDate,
									UpdatedDatevar:      myUpdatedDate,
								}

								retScanSummary = append(retScanSummary, newScanSummary)
							}
						}

						return err
					})
			}
		},
	})

	return retScanSummary, err
}

// GetUnmatchedVulns executes the stored procedure GetUnmatchedVulns against the database and returns the read results
func (conn *dbconn) GetUnmatchedVulns(_SourceID int) ([]domain.VulnerabilityInfo, error) {
	var err error
	var retVulnerabilityInfo = make([]domain.VulnerabilityInfo, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetUnmatchedVulns",
		Parameters: []interface{}{_SourceID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var mySourceVulnID string
							var myTitle string
							var myVulnerabilityID *string
							var mySourceID string
							var myCVSSScore float32
							var myCVSS3Score *float32
							var myDescription string
							var mySolution string

							if err = rows.Scan(

								&myID,
								&mySourceVulnID,
								&myTitle,
								&myVulnerabilityID,
								&mySourceID,
								&myCVSSScore,
								&myCVSS3Score,
								&myDescription,
								&mySolution,
							); err == nil {

								newVulnerabilityInfo := &dal.VulnerabilityInfo{
									IDvar:              myID,
									SourceVulnIDvar:    mySourceVulnID,
									Titlevar:           myTitle,
									VulnerabilityIDvar: myVulnerabilityID,
									SourceIDvar:        mySourceID,
									CVSSScorevar:       myCVSSScore,
									CVSS3Scorevar:      myCVSS3Score,
									Descriptionvar:     myDescription,
									Solutionvar:        mySolution,
								}

								retVulnerabilityInfo = append(retVulnerabilityInfo, newVulnerabilityInfo)
							}
						}

						return err
					})
			}
		},
	})

	return retVulnerabilityInfo, err
}

// GetUserAnyOrg executes the stored procedure GetUserAnyOrg against the database and returns the read results
func (conn *dbconn) GetUserAnyOrg(_ID string) (domain.User, error) {
	var err error
	var retUser domain.User

	conn.Read(&connection.Procedure{
		Proc:       "GetUserAnyOrg",
		Parameters: []interface{}{_ID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myUsername *string
							var myFirstName string
							var myLastName string
							var myEmail string
							var myIsDisabled []uint8

							if err = rows.Scan(

								&myID,
								&myUsername,
								&myFirstName,
								&myLastName,
								&myEmail,
								&myIsDisabled,
							); err == nil {

								newUser := &dal.User{
									IDvar:         myID,
									Usernamevar:   myUsername,
									FirstNamevar:  myFirstName,
									LastNamevar:   myLastName,
									Emailvar:      myEmail,
									IsDisabledvar: myIsDisabled[0] > 0 && myIsDisabled[0] != 48, // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)

								}

								retUser = newUser
							}
						}

						return err
					})
			}
		},
	})

	return retUser, err
}

// GetUserByID executes the stored procedure GetUserByID against the database and returns the read results
func (conn *dbconn) GetUserByID(_ID string, _OrgID string) (domain.User, error) {
	var err error
	var retUser domain.User

	conn.Read(&connection.Procedure{
		Proc:       "GetUserByID",
		Parameters: []interface{}{_ID, _OrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myUsername *string
							var myFirstName string
							var myLastName string
							var myEmail string
							var myIsDisabled []uint8

							if err = rows.Scan(

								&myID,
								&myUsername,
								&myFirstName,
								&myLastName,
								&myEmail,
								&myIsDisabled,
							); err == nil {

								newUser := &dal.User{
									IDvar:         myID,
									Usernamevar:   myUsername,
									FirstNamevar:  myFirstName,
									LastNamevar:   myLastName,
									Emailvar:      myEmail,
									IsDisabledvar: myIsDisabled[0] > 0 && myIsDisabled[0] != 48, // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)

								}

								retUser = newUser
							}
						}

						return err
					})
			}
		},
	})

	return retUser, err
}

// GetUserByUsername executes the stored procedure GetUserByUsername against the database and returns the read results
func (conn *dbconn) GetUserByUsername(_Username string) (domain.User, error) {
	var err error
	var retUser domain.User

	conn.Read(&connection.Procedure{
		Proc:       "GetUserByUsername",
		Parameters: []interface{}{_Username},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myUsername *string
							var myFirstName string
							var myLastName string
							var myEmail string
							var myIsDisabled []uint8

							if err = rows.Scan(

								&myID,
								&myUsername,
								&myFirstName,
								&myLastName,
								&myEmail,
								&myIsDisabled,
							); err == nil {

								newUser := &dal.User{
									IDvar:         myID,
									Usernamevar:   myUsername,
									FirstNamevar:  myFirstName,
									LastNamevar:   myLastName,
									Emailvar:      myEmail,
									IsDisabledvar: myIsDisabled[0] > 0 && myIsDisabled[0] != 48, // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)

								}

								retUser = newUser
							}
						}

						return err
					})
			}
		},
	})

	return retUser, err
}

// GetUsersByOrg executes the stored procedure GetUsersByOrg against the database and returns the read results
func (conn *dbconn) GetUsersByOrg(_OrgID string) ([]domain.User, error) {
	var err error
	var retUser = make([]domain.User, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetUsersByOrg",
		Parameters: []interface{}{_OrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var myUsername *string
							var myFirstName string
							var myLastName string
							var myEmail string
							var myIsDisabled []uint8

							if err = rows.Scan(

								&myID,
								&myUsername,
								&myFirstName,
								&myLastName,
								&myEmail,
								&myIsDisabled,
							); err == nil {

								newUser := &dal.User{
									IDvar:         myID,
									Usernamevar:   myUsername,
									FirstNamevar:  myFirstName,
									LastNamevar:   myLastName,
									Emailvar:      myEmail,
									IsDisabledvar: myIsDisabled[0] > 0 && myIsDisabled[0] != 48, // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)

								}

								retUser = append(retUser, newUser)
							}
						}

						return err
					})
			}
		},
	})

	return retUser, err
}

// GetVulnInfoByID executes the stored procedure GetVulnInfoByID against the database and returns the read results
func (conn *dbconn) GetVulnInfoByID(_ID string) (domain.VulnerabilityInfo, error) {
	var err error
	var retVulnerabilityInfo domain.VulnerabilityInfo

	conn.Read(&connection.Procedure{
		Proc:       "GetVulnInfoByID",
		Parameters: []interface{}{_ID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var mySourceVulnID string
							var myTitle string
							var myVulnerabilityID *string
							var mySourceID string
							var myCVSSScore float32
							var myCVSS3Score *float32
							var myDescription string
							var myThreat *string
							var mySolution string
							var myDetectionInformation *string
							var myPatchable *string
							var myCategory *string
							var mySoftware *string
							var myCreated *time.Time
							var myUpdated *time.Time

							if err = rows.Scan(

								&myID,
								&mySourceVulnID,
								&myTitle,
								&myVulnerabilityID,
								&mySourceID,
								&myCVSSScore,
								&myCVSS3Score,
								&myDescription,
								&myThreat,
								&mySolution,
								&myDetectionInformation,
								&myPatchable,
								&myCategory,
								&mySoftware,
								&myCreated,
								&myUpdated,
							); err == nil {

								newVulnerabilityInfo := &dal.VulnerabilityInfo{
									IDvar:                   myID,
									SourceVulnIDvar:         mySourceVulnID,
									Titlevar:                myTitle,
									VulnerabilityIDvar:      myVulnerabilityID,
									SourceIDvar:             mySourceID,
									CVSSScorevar:            myCVSSScore,
									CVSS3Scorevar:           myCVSS3Score,
									Descriptionvar:          myDescription,
									Threatvar:               myThreat,
									Solutionvar:             mySolution,
									DetectionInformationvar: myDetectionInformation,
									Patchablevar:            myPatchable,
									Categoryvar:             myCategory,
									Softwarevar:             mySoftware,
									Createdvar:              myCreated,
									Updatedvar:              myUpdated,
								}

								retVulnerabilityInfo = newVulnerabilityInfo
							}
						}

						return err
					})
			}
		},
	})

	return retVulnerabilityInfo, err
}

// GetVulnInfoBySource executes the stored procedure GetVulnInfoBySource against the database and returns the read results
func (conn *dbconn) GetVulnInfoBySource(_Source string) ([]domain.VulnerabilityInfo, error) {
	var err error
	var retVulnerabilityInfo = make([]domain.VulnerabilityInfo, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetVulnInfoBySource",
		Parameters: []interface{}{_Source},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var mySourceVulnID string
							var myTitle string
							var myVulnerabilityID *string
							var mySourceID string
							var myCVSSScore float32
							var myCVSS3Score *float32
							var myDescription string
							var myThreat *string
							var mySolution string
							var myDetectionInformation *string
							var myPatchable *string
							var myCategory *string
							var mySoftware *string
							var myCreated *time.Time
							var myUpdated *time.Time

							if err = rows.Scan(

								&myID,
								&mySourceVulnID,
								&myTitle,
								&myVulnerabilityID,
								&mySourceID,
								&myCVSSScore,
								&myCVSS3Score,
								&myDescription,
								&myThreat,
								&mySolution,
								&myDetectionInformation,
								&myPatchable,
								&myCategory,
								&mySoftware,
								&myCreated,
								&myUpdated,
							); err == nil {

								newVulnerabilityInfo := &dal.VulnerabilityInfo{
									IDvar:                   myID,
									SourceVulnIDvar:         mySourceVulnID,
									Titlevar:                myTitle,
									VulnerabilityIDvar:      myVulnerabilityID,
									SourceIDvar:             mySourceID,
									CVSSScorevar:            myCVSSScore,
									CVSS3Scorevar:           myCVSS3Score,
									Descriptionvar:          myDescription,
									Threatvar:               myThreat,
									Solutionvar:             mySolution,
									DetectionInformationvar: myDetectionInformation,
									Patchablevar:            myPatchable,
									Categoryvar:             myCategory,
									Softwarevar:             mySoftware,
									Createdvar:              myCreated,
									Updatedvar:              myUpdated,
								}

								retVulnerabilityInfo = append(retVulnerabilityInfo, newVulnerabilityInfo)
							}
						}

						return err
					})
			}
		},
	})

	return retVulnerabilityInfo, err
}

// GetVulnInfoBySourceID executes the stored procedure GetVulnInfoBySourceID against the database and returns the read results
func (conn *dbconn) GetVulnInfoBySourceID(_SourceID string) ([]domain.VulnerabilityInfo, error) {
	var err error
	var retVulnerabilityInfo = make([]domain.VulnerabilityInfo, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetVulnInfoBySourceID",
		Parameters: []interface{}{_SourceID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var mySourceVulnID string
							var myTitle string
							var myVulnerabilityID *string
							var mySourceID string
							var myCVSSScore float32
							var myCVSS3Score *float32
							var myDescription string
							var myThreat *string
							var mySolution string
							var myDetectionInformation *string
							var myPatchable *string
							var myCategory *string
							var mySoftware *string
							var myCreated *time.Time
							var myUpdated *time.Time

							if err = rows.Scan(

								&myID,
								&mySourceVulnID,
								&myTitle,
								&myVulnerabilityID,
								&mySourceID,
								&myCVSSScore,
								&myCVSS3Score,
								&myDescription,
								&myThreat,
								&mySolution,
								&myDetectionInformation,
								&myPatchable,
								&myCategory,
								&mySoftware,
								&myCreated,
								&myUpdated,
							); err == nil {

								newVulnerabilityInfo := &dal.VulnerabilityInfo{
									IDvar:                   myID,
									SourceVulnIDvar:         mySourceVulnID,
									Titlevar:                myTitle,
									VulnerabilityIDvar:      myVulnerabilityID,
									SourceIDvar:             mySourceID,
									CVSSScorevar:            myCVSSScore,
									CVSS3Scorevar:           myCVSS3Score,
									Descriptionvar:          myDescription,
									Threatvar:               myThreat,
									Solutionvar:             mySolution,
									DetectionInformationvar: myDetectionInformation,
									Patchablevar:            myPatchable,
									Categoryvar:             myCategory,
									Softwarevar:             mySoftware,
									Createdvar:              myCreated,
									Updatedvar:              myUpdated,
								}

								retVulnerabilityInfo = append(retVulnerabilityInfo, newVulnerabilityInfo)
							}
						}

						return err
					})
			}
		},
	})

	return retVulnerabilityInfo, err
}

// GetVulnInfoBySourceVulnID executes the stored procedure GetVulnInfoBySourceVulnID against the database and returns the read results
func (conn *dbconn) GetVulnInfoBySourceVulnID(_SourceVulnID string) (domain.VulnerabilityInfo, error) {
	var err error
	var retVulnerabilityInfo domain.VulnerabilityInfo

	conn.Read(&connection.Procedure{
		Proc:       "GetVulnInfoBySourceVulnID",
		Parameters: []interface{}{_SourceVulnID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var mySourceVulnID string
							var myTitle string
							var myVulnerabilityID *string
							var mySourceID string
							var myCVSSScore float32
							var myCVSS3Score *float32
							var myDescription string
							var myThreat *string
							var mySolution string
							var myDetectionInformation *string
							var myPatchable *string
							var myCategory *string
							var mySoftware *string
							var myCreated *time.Time
							var myUpdated *time.Time

							if err = rows.Scan(

								&myID,
								&mySourceVulnID,
								&myTitle,
								&myVulnerabilityID,
								&mySourceID,
								&myCVSSScore,
								&myCVSS3Score,
								&myDescription,
								&myThreat,
								&mySolution,
								&myDetectionInformation,
								&myPatchable,
								&myCategory,
								&mySoftware,
								&myCreated,
								&myUpdated,
							); err == nil {

								newVulnerabilityInfo := &dal.VulnerabilityInfo{
									IDvar:                   myID,
									SourceVulnIDvar:         mySourceVulnID,
									Titlevar:                myTitle,
									VulnerabilityIDvar:      myVulnerabilityID,
									SourceIDvar:             mySourceID,
									CVSSScorevar:            myCVSSScore,
									CVSS3Scorevar:           myCVSS3Score,
									Descriptionvar:          myDescription,
									Threatvar:               myThreat,
									Solutionvar:             mySolution,
									DetectionInformationvar: myDetectionInformation,
									Patchablevar:            myPatchable,
									Categoryvar:             myCategory,
									Softwarevar:             mySoftware,
									Createdvar:              myCreated,
									Updatedvar:              myUpdated,
								}

								retVulnerabilityInfo = newVulnerabilityInfo
							}
						}

						return err
					})
			}
		},
	})

	return retVulnerabilityInfo, err
}

// GetVulnInfoBySourceVulnIDSourceID executes the stored procedure GetVulnInfoBySourceVulnIDSourceID against the database and returns the read results
func (conn *dbconn) GetVulnInfoBySourceVulnIDSourceID(_SourceVulnID string, _SourceID string, _Modified time.Time) (domain.VulnerabilityInfo, error) {
	var err error
	var retVulnerabilityInfo domain.VulnerabilityInfo

	conn.Read(&connection.Procedure{
		Proc:       "GetVulnInfoBySourceVulnIDSourceID",
		Parameters: []interface{}{_SourceVulnID, _SourceID, _Modified},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var mySourceVulnID string
							var myTitle string
							var myVulnerabilityID *string
							var mySourceID string
							var myCVSSScore float32
							var myCVSS3Score *float32
							var myDescription string
							var myThreat *string
							var mySolution string
							var myDetectionInformation *string
							var myPatchable *string
							var myCategory *string
							var mySoftware *string
							var myCreated *time.Time
							var myUpdated *time.Time

							if err = rows.Scan(

								&myID,
								&mySourceVulnID,
								&myTitle,
								&myVulnerabilityID,
								&mySourceID,
								&myCVSSScore,
								&myCVSS3Score,
								&myDescription,
								&myThreat,
								&mySolution,
								&myDetectionInformation,
								&myPatchable,
								&myCategory,
								&mySoftware,
								&myCreated,
								&myUpdated,
							); err == nil {

								newVulnerabilityInfo := &dal.VulnerabilityInfo{
									IDvar:                   myID,
									SourceVulnIDvar:         mySourceVulnID,
									Titlevar:                myTitle,
									VulnerabilityIDvar:      myVulnerabilityID,
									SourceIDvar:             mySourceID,
									CVSSScorevar:            myCVSSScore,
									CVSS3Scorevar:           myCVSS3Score,
									Descriptionvar:          myDescription,
									Threatvar:               myThreat,
									Solutionvar:             mySolution,
									DetectionInformationvar: myDetectionInformation,
									Patchablevar:            myPatchable,
									Categoryvar:             myCategory,
									Softwarevar:             mySoftware,
									Createdvar:              myCreated,
									Updatedvar:              myUpdated,
								}

								retVulnerabilityInfo = newVulnerabilityInfo
							}
						}

						return err
					})
			}
		},
	})

	return retVulnerabilityInfo, err
}

// GetVulnRefInfo executes the stored procedure GetVulnRefInfo against the database and returns the read results
func (conn *dbconn) GetVulnRefInfo(_VulnInfoID string, _SourceID string, _Reference string) (domain.VulnerabilityReferenceInfo, error) {
	var err error
	var retVulnerabilityReferenceInfo domain.VulnerabilityReferenceInfo

	conn.Read(&connection.Procedure{
		Proc:       "GetVulnRefInfo",
		Parameters: []interface{}{_VulnInfoID, _SourceID, _Reference},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string

							if err = rows.Scan(

								&myID,
							); err == nil {

								newVulnerabilityReferenceInfo := &dal.VulnerabilityReferenceInfo{
									IDvar: myID,
								}

								retVulnerabilityReferenceInfo = newVulnerabilityReferenceInfo
							}
						}

						return err
					})
			}
		},
	})

	return retVulnerabilityReferenceInfo, err
}

// GetVulnRefInfoVendor executes the stored procedure GetVulnRefInfoVendor against the database and returns the read results
func (conn *dbconn) GetVulnRefInfoVendor(_VulnInfoID string, _SourceID string) ([]domain.VulnerabilityReferenceInfo, error) {
	var err error
	var retVulnerabilityReferenceInfo = make([]domain.VulnerabilityReferenceInfo, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetVulnRefInfoVendor",
		Parameters: []interface{}{_VulnInfoID, _SourceID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myVulnInfoID string
							var mySourceID string
							var myReference string
							var myRefType int

							if err = rows.Scan(

								&myVulnInfoID,
								&mySourceID,
								&myReference,
								&myRefType,
							); err == nil {

								newVulnerabilityReferenceInfo := &dal.VulnerabilityReferenceInfo{
									VulnInfoIDvar: myVulnInfoID,
									SourceIDvar:   mySourceID,
									Referencevar:  myReference,
									RefTypevar:    myRefType,
								}

								retVulnerabilityReferenceInfo = append(retVulnerabilityReferenceInfo, newVulnerabilityReferenceInfo)
							}
						}

						return err
					})
			}
		},
	})

	return retVulnerabilityReferenceInfo, err
}

// GetVulnReferencesInfo executes the stored procedure GetVulnReferencesInfo against the database and returns the read results
func (conn *dbconn) GetVulnReferencesInfo(_VulnInfoID string, _SourceID string) ([]domain.VulnerabilityReferenceInfo, error) {
	var err error
	var retVulnerabilityReferenceInfo = make([]domain.VulnerabilityReferenceInfo, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetVulnReferencesInfo",
		Parameters: []interface{}{_VulnInfoID, _SourceID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myVulnInfoID string
							var mySourceID string
							var myReference string
							var myRefType int

							if err = rows.Scan(

								&myVulnInfoID,
								&mySourceID,
								&myReference,
								&myRefType,
							); err == nil {

								newVulnerabilityReferenceInfo := &dal.VulnerabilityReferenceInfo{
									VulnInfoIDvar: myVulnInfoID,
									SourceIDvar:   mySourceID,
									Referencevar:  myReference,
									RefTypevar:    myRefType,
								}

								retVulnerabilityReferenceInfo = append(retVulnerabilityReferenceInfo, newVulnerabilityReferenceInfo)
							}
						}

						return err
					})
			}
		},
	})

	return retVulnerabilityReferenceInfo, err
}

// GetVulnReferencesInfoBySourceAndRef executes the stored procedure GetVulnReferencesInfoBySourceAndRef against the database and returns the read results
func (conn *dbconn) GetVulnReferencesInfoBySourceAndRef(_SourceID string, _Reference string) ([]domain.VulnerabilityReferenceInfo, error) {
	var err error
	var retVulnerabilityReferenceInfo = make([]domain.VulnerabilityReferenceInfo, 0)

	conn.Read(&connection.Procedure{
		Proc:       "GetVulnReferencesInfoBySourceAndRef",
		Parameters: []interface{}{_SourceID, _Reference},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myVulnInfoID string
							var mySourceID string
							var myReference string

							if err = rows.Scan(

								&myVulnInfoID,
								&mySourceID,
								&myReference,
							); err == nil {

								newVulnerabilityReferenceInfo := &dal.VulnerabilityReferenceInfo{
									VulnInfoIDvar: myVulnInfoID,
									SourceIDvar:   mySourceID,
									Referencevar:  myReference,
								}

								retVulnerabilityReferenceInfo = append(retVulnerabilityReferenceInfo, newVulnerabilityReferenceInfo)
							}
						}

						return err
					})
			}
		},
	})

	return retVulnerabilityReferenceInfo, err
}

// HasDecommissioned executes the stored procedure HasDecommissioned against the database and returns the read results
func (conn *dbconn) HasDecommissioned(_devID string, _sourceID string, _orgID string) (domain.Ignore, error) {
	var err error
	var retIgnore domain.Ignore

	conn.Read(&connection.Procedure{
		Proc:       "HasDecommissioned",
		Parameters: []interface{}{_devID, _sourceID, _orgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var mySourceID string
							var myOrganizationID string
							var myTypeID int
							var myVulnerabilityID string
							var myDeviceID string
							var myDueDate *time.Time
							var myApproval string
							var myActive []uint8
							var myDBCreatedDate time.Time
							var myDBUpdatedDate *time.Time

							if err = rows.Scan(

								&myID,
								&mySourceID,
								&myOrganizationID,
								&myTypeID,
								&myVulnerabilityID,
								&myDeviceID,
								&myDueDate,
								&myApproval,
								&myActive,
								&myDBCreatedDate,
								&myDBUpdatedDate,
							); err == nil {

								newIgnore := &dal.Ignore{
									IDvar:              myID,
									SourceIDvar:        mySourceID,
									OrganizationIDvar:  myOrganizationID,
									TypeIDvar:          myTypeID,
									VulnerabilityIDvar: myVulnerabilityID,
									DeviceIDvar:        myDeviceID,
									DueDatevar:         myDueDate,
									Approvalvar:        myApproval,
									Activevar:          myActive[0] > 0 && myActive[0] != 48, // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)
									DBCreatedDatevar:   myDBCreatedDate,
									DBUpdatedDatevar:   myDBUpdatedDate,
								}

								retIgnore = newIgnore
							}
						}

						return err
					})
			}
		},
	})

	return retIgnore, err
}

// HasExceptionOrFalsePositive executes the stored procedure HasExceptionOrFalsePositive against the database and returns the read results
func (conn *dbconn) HasExceptionOrFalsePositive(_sourceID string, _vulnID string, _devID string, _orgID string, _port string, _OS string) ([]domain.Ignore, error) {
	var err error
	var retIgnore = make([]domain.Ignore, 0)

	conn.Read(&connection.Procedure{
		Proc:       "HasExceptionOrFalsePositive",
		Parameters: []interface{}{_sourceID, _vulnID, _devID, _orgID, _port, _OS},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var mySourceID string
							var myOrganizationID string
							var myTypeID int
							var myVulnerabilityID string
							var myDeviceID string
							var myDueDate *time.Time
							var myApproval string
							var myActive []uint8
							var myDBCreatedDate time.Time
							var myDBUpdatedDate *time.Time

							if err = rows.Scan(

								&myID,
								&mySourceID,
								&myOrganizationID,
								&myTypeID,
								&myVulnerabilityID,
								&myDeviceID,
								&myDueDate,
								&myApproval,
								&myActive,
								&myDBCreatedDate,
								&myDBUpdatedDate,
							); err == nil {

								newIgnore := &dal.Ignore{
									IDvar:              myID,
									SourceIDvar:        mySourceID,
									OrganizationIDvar:  myOrganizationID,
									TypeIDvar:          myTypeID,
									VulnerabilityIDvar: myVulnerabilityID,
									DeviceIDvar:        myDeviceID,
									DueDatevar:         myDueDate,
									Approvalvar:        myApproval,
									Activevar:          myActive[0] > 0 && myActive[0] != 48, // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)
									DBCreatedDatevar:   myDBCreatedDate,
									DBUpdatedDatevar:   myDBUpdatedDate,
								}

								retIgnore = append(retIgnore, newIgnore)
							}
						}

						return err
					})
			}
		},
	})

	return retIgnore, err
}

// HasIgnore executes the stored procedure HasIgnore against the database and returns the read results
func (conn *dbconn) HasIgnore(inSourceID string, inVulnID string, inDevID string, inOrgID string, inPort string, inMostCurrentDetection time.Time) (domain.Ignore, error) {
	var err error
	var retIgnore domain.Ignore

	conn.Read(&connection.Procedure{
		Proc:       "HasIgnore",
		Parameters: []interface{}{inSourceID, inVulnID, inDevID, inOrgID, inPort, inMostCurrentDetection},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if err == nil {

				err = conn.getRows(results,
					func(rows *sql.Rows) (err error) {
						if err = rows.Err(); err == nil {

							var myID string
							var mySourceID string
							var myOrganizationID string
							var myTypeID int
							var myVulnerabilityID string
							var myDeviceID string
							var myDueDate *time.Time
							var myApproval string
							var myActive []uint8
							var myDBCreatedDate time.Time
							var myDBUpdatedDate *time.Time

							if err = rows.Scan(

								&myID,
								&mySourceID,
								&myOrganizationID,
								&myTypeID,
								&myVulnerabilityID,
								&myDeviceID,
								&myDueDate,
								&myApproval,
								&myActive,
								&myDBCreatedDate,
								&myDBUpdatedDate,
							); err == nil {

								newIgnore := &dal.Ignore{
									IDvar:              myID,
									SourceIDvar:        mySourceID,
									OrganizationIDvar:  myOrganizationID,
									TypeIDvar:          myTypeID,
									VulnerabilityIDvar: myVulnerabilityID,
									DeviceIDvar:        myDeviceID,
									DueDatevar:         myDueDate,
									Approvalvar:        myApproval,
									Activevar:          myActive[0] > 0 && myActive[0] != 48, // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false)
									DBCreatedDatevar:   myDBCreatedDate,
									DBUpdatedDatevar:   myDBUpdatedDate,
								}

								retIgnore = newIgnore
							}
						}

						return err
					})
			}
		},
	})

	return retIgnore, err
}

// PulseJob executes the stored procedure PulseJob against the database
func (conn *dbconn) PulseJob(_JobHistoryID string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "PulseJob",
		Parameters: []interface{}{_JobHistoryID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// SaveAssignmentGroup executes the stored procedure SaveAssignmentGroup against the database
func (conn *dbconn) SaveAssignmentGroup(_SourceID string, _OrganizationID string, _IpAddress string, _GroupName string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "SaveAssignmentGroup",
		Parameters: []interface{}{_SourceID, _OrganizationID, _IpAddress, _GroupName},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// SaveIgnore executes the stored procedure SaveIgnore against the database
func (conn *dbconn) SaveIgnore(_SourceID string, _OrganizationID string, _TypeID int, _VulnerabilityID string, _DeviceID string, _DueDate time.Time, _Approval string, _Active bool, _port string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "SaveIgnore",
		Parameters: []interface{}{_SourceID, _OrganizationID, _TypeID, _VulnerabilityID, _DeviceID, _DueDate, _Approval, _Active, _port},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// SaveScanSummary executes the stored procedure SaveScanSummary against the database
func (conn *dbconn) SaveScanSummary(_ScanID string, _ScanStatus string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "SaveScanSummary",
		Parameters: []interface{}{_ScanID, _ScanStatus},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// SetScheduleLastRun executes the stored procedure SetScheduleLastRun against the database
func (conn *dbconn) SetScheduleLastRun(_ID string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "SetScheduleLastRun",
		Parameters: []interface{}{_ID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// UpdateAssetGroupLastTicket executes the stored procedure UpdateAssetGroupLastTicket against the database
func (conn *dbconn) UpdateAssetGroupLastTicket(inGroupID string, inOrgID string, inLastTicketTime time.Time) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "UpdateAssetGroupLastTicket",
		Parameters: []interface{}{inGroupID, inOrgID, inLastTicketTime},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// UpdateAssetIDOsTypeIDOfDevice executes the stored procedure UpdateAssetIDOsTypeIDOfDevice against the database
func (conn *dbconn) UpdateAssetIDOsTypeIDOfDevice(_ID string, _AssetID string, _ScannerSourceID string, _GroupID string, _OS string, _HostName string, _OsTypeID int, inTrackingMethod string, _OrgID string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "UpdateAssetIDOsTypeIDOfDevice",
		Parameters: []interface{}{_ID, _AssetID, _ScannerSourceID, _GroupID, _OS, _HostName, _OsTypeID, inTrackingMethod, _OrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// UpdateDetection executes the stored procedure UpdateDetection against the database
func (conn *dbconn) UpdateDetection(_ID string, _DeviceID string, _VulnID string, _Port int, _Protocol string, _ExceptionID string, _TimesSeen int, _StatusID int, _LastFound time.Time, _LastUpdated time.Time, _DefaultTime time.Time) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "UpdateDetection",
		Parameters: []interface{}{_ID, _DeviceID, _VulnID, _Port, _Protocol, _ExceptionID, _TimesSeen, _StatusID, _LastFound, _LastUpdated, _DefaultTime},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// UpdateDetectionIgnore executes the stored procedure UpdateDetectionIgnore against the database
func (conn *dbconn) UpdateDetectionIgnore(_DeviceID string, _VulnID string, _Port int, _Protocol string, _ExceptionID string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "UpdateDetectionIgnore",
		Parameters: []interface{}{_DeviceID, _VulnID, _Port, _Protocol, _ExceptionID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// UpdateExpirationDateByCERF executes the stored procedure UpdateExpirationDateByCERF against the database
func (conn *dbconn) UpdateExpirationDateByCERF(_CERForm string, _OrganizationID string, _DueDate time.Time) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "UpdateExpirationDateByCERF",
		Parameters: []interface{}{_CERForm, _OrganizationID, _DueDate},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// UpdateInstanceIDOfDevice executes the stored procedure UpdateInstanceIDOfDevice against the database
func (conn *dbconn) UpdateInstanceIDOfDevice(_ID string, _InstanceID string, _CloudSourceID string, _State string, _Region string, _OrgID string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "UpdateInstanceIDOfDevice",
		Parameters: []interface{}{_ID, _InstanceID, _CloudSourceID, _State, _Region, _OrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// UpdateJobConfig executes the stored procedure UpdateJobConfig against the database
func (conn *dbconn) UpdateJobConfig(_ID string, _DataInSourceID string, _DataOutSourceID string, _Autostart bool, _PriorityOverride int, _Continuous bool, _WaitInSeconds int, _MaxInstances int, _UpdatedBy string, _OrgID string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "UpdateJobConfig",
		Parameters: []interface{}{_ID, _DataInSourceID, _DataOutSourceID, _Autostart, _PriorityOverride, _Continuous, _WaitInSeconds, _MaxInstances, _UpdatedBy, _OrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// UpdateJobConfigLastRun executes the stored procedure UpdateJobConfigLastRun against the database
func (conn *dbconn) UpdateJobConfigLastRun(_ID string, _LastRun time.Time) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "UpdateJobConfigLastRun",
		Parameters: []interface{}{_ID, _LastRun},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// UpdateJobHistory executes the stored procedure UpdateJobHistory against the database
func (conn *dbconn) UpdateJobHistory(_ID string, _ConfigID string, _Payload string, _UpdatedBy string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "UpdateJobHistory",
		Parameters: []interface{}{_ID, _ConfigID, _Payload, _UpdatedBy},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// UpdateJobHistoryStatus executes the stored procedure UpdateJobHistoryStatus against the database
func (conn *dbconn) UpdateJobHistoryStatus(_ID string, _Status int) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "UpdateJobHistoryStatus",
		Parameters: []interface{}{_ID, _Status},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// UpdateJobHistoryStatusDetailed executes the stored procedure UpdateJobHistoryStatusDetailed against the database
func (conn *dbconn) UpdateJobHistoryStatusDetailed(_ID string, _Status int, _UpdatedBy string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "UpdateJobHistoryStatusDetailed",
		Parameters: []interface{}{_ID, _Status, _UpdatedBy},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// UpdateOrganization executes the stored procedure UpdateOrganization against the database
func (conn *dbconn) UpdateOrganization(_ID string, _Description string, _TimezoneOffset float32, _UpdatedBy string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "UpdateOrganization",
		Parameters: []interface{}{_ID, _Description, _TimezoneOffset, _UpdatedBy},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// UpdatePermissionsByUserOrgID executes the stored procedure UpdatePermissionsByUserOrgID against the database
func (conn *dbconn) UpdatePermissionsByUserOrgID(_UserID string, _OrgID string, _Admin bool, _Manager bool, _Reader bool, _Reporter bool) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "UpdatePermissionsByUserOrgID",
		Parameters: []interface{}{_UserID, _OrgID, _Admin, _Manager, _Reader, _Reporter},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// UpdateSourceConfig executes the stored procedure UpdateSourceConfig against the database
func (conn *dbconn) UpdateSourceConfig(_ID string, _OrgID string, _Address string, _Username string, _Password string, _PrivateKey string, _ConsumerKey string, _Token string, _Port string, _Payload string, _UpdatedBy string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "UpdateSourceConfig",
		Parameters: []interface{}{_ID, _OrgID, _Address, _Username, _Password, _PrivateKey, _ConsumerKey, _Token, _Port, _Payload, _UpdatedBy},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// UpdateSourceConfigConcurrencyByID executes the stored procedure UpdateSourceConfigConcurrencyByID against the database
func (conn *dbconn) UpdateSourceConfigConcurrencyByID(_ID string, _Delay int, _Retries int, _Concurrency int) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "UpdateSourceConfigConcurrencyByID",
		Parameters: []interface{}{_ID, _Delay, _Retries, _Concurrency},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// UpdateSourceConfigToken executes the stored procedure UpdateSourceConfigToken against the database
func (conn *dbconn) UpdateSourceConfigToken(_ID string, _Token string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "UpdateSourceConfigToken",
		Parameters: []interface{}{_ID, _Token},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// UpdateStateOfDevice executes the stored procedure UpdateStateOfDevice against the database
func (conn *dbconn) UpdateStateOfDevice(_ID string, _State string, _OrgID string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "UpdateStateOfDevice",
		Parameters: []interface{}{_ID, _State, _OrgID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// UpdateTag executes the stored procedure UpdateTag against the database
func (conn *dbconn) UpdateTag(_DeviceID string, _TagKeyID string, _Value string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "UpdateTag",
		Parameters: []interface{}{_DeviceID, _TagKeyID, _Value},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// UpdateTagMap executes the stored procedure UpdateTagMap against the database
func (conn *dbconn) UpdateTagMap(_TicketingSourceID string, _TicketingTag string, _CloudSourceID string, _CloudTag string, _Options string, _OrganizationID string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "UpdateTagMap",
		Parameters: []interface{}{_TicketingSourceID, _TicketingTag, _CloudSourceID, _CloudTag, _Options, _OrganizationID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// UpdateTicket executes the stored procedure UpdateTicket against the database
func (conn *dbconn) UpdateTicket(_Title string, _Status string, _OrganizationID string, _AssignmentGroup string, _Assignee string, _DueDate time.Time, _CreatedDate time.Time, _UpdatedDate time.Time, _ResolutionDate time.Time, _DefaultTime time.Time) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "UpdateTicket",
		Parameters: []interface{}{_Title, _Status, _OrganizationID, _AssignmentGroup, _Assignee, _DueDate, _CreatedDate, _UpdatedDate, _ResolutionDate, _DefaultTime},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// UpdateTicketDetectionID executes the stored procedure UpdateTicketDetectionID against the database
func (conn *dbconn) UpdateTicketDetectionID(_Title string, _DetectionID string, _OrganizationID string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "UpdateTicketDetectionID",
		Parameters: []interface{}{_Title, _DetectionID, _OrganizationID},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// UpdateUserByID executes the stored procedure UpdateUserByID against the database
func (conn *dbconn) UpdateUserByID(_ID string, _FirstName string, _LastName string, _Email string, _Disabled bool) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "UpdateUserByID",
		Parameters: []interface{}{_ID, _FirstName, _LastName, _Email, _Disabled},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// UpdateVulnByID executes the stored procedure UpdateVulnByID against the database
func (conn *dbconn) UpdateVulnByID(_ID string, _SourceVulnID string, _Title string, _SourceID string, _CVSSScore float32, _CVSS3Score float32, _Description string, _Threat string, _Solution string, _Software string, _Patchable string, _Category string, _DetectionInformation string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "UpdateVulnByID",
		Parameters: []interface{}{_ID, _SourceVulnID, _Title, _SourceID, _CVSSScore, _CVSS3Score, _Description, _Threat, _Solution, _Software, _Patchable, _Category, _DetectionInformation},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// UpdateVulnByIDNoCVSS3 executes the stored procedure UpdateVulnByIDNoCVSS3 against the database
func (conn *dbconn) UpdateVulnByIDNoCVSS3(_ID string, _SourceVulnID string, _Title string, _SourceID string, _CVSSScore float32, _CVSS3Score float32, _Description string, _Threat string, _Solution string, _Software string, _Patchable string, _DetectionInformation string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "UpdateVulnByIDNoCVSS3",
		Parameters: []interface{}{_ID, _SourceVulnID, _Title, _SourceID, _CVSSScore, _CVSS3Score, _Description, _Threat, _Solution, _Software, _Patchable, _DetectionInformation},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}

// UpdateVulnInfoID executes the stored procedure UpdateVulnInfoID against the database
func (conn *dbconn) UpdateVulnInfoID(_VulnInfoID string, _VulnID string, _MatchConfidence int, _MatchReasons string) (id int, affectedRows int, err error) {

	conn.Exec(&connection.Procedure{
		Proc:       "UpdateVulnInfoID",
		Parameters: []interface{}{_VulnInfoID, _VulnID, _MatchConfidence, _MatchReasons},
		Callback: func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}
