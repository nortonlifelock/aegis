package endpoints

import (
	"fmt"
	"net/http"
	"time"

	"github.com/nortonlifelock/domain"
	"github.com/pkg/errors"
)

func createException(w http.ResponseWriter, r *http.Request) {
	executeTransaction(w, r, createExceptionEndpoints, admin|manager, func(trans *transaction) {
		if trans.endpoint.verify() {
			if except, isAExcept := trans.endpoint.(*Exception); isAExcept {
				trans.obj, trans.status, trans.err = except.createOrUpdate(trans.user, trans.permission)
				(&trans.wrapper).addError(trans.err, processError)
			} else {
				trans.err = errors.Errorf("tried to pass a non-exception as a exception")
				(&trans.wrapper).addError(trans.err, requestFormatError)
			}
		} else {
			trans.err = errors.Errorf("invalid length of format of request fields")
			(&trans.wrapper).addError(trans.err, requestFormatError)
		}
	})
}

func updateException(w http.ResponseWriter, r *http.Request) {
	executeTransaction(w, r, updateExceptionEndpoint, admin|manager, func(trans *transaction) {
		if trans.endpoint.verify() {
			if except, isAExcept := trans.endpoint.(*Exception); isAExcept {
				trans.obj, trans.status, trans.err = except.createOrUpdate(trans.user, trans.permission)
				(&trans.wrapper).addError(trans.err, processError)
			} else {
				trans.err = errors.Errorf("tried to pass a non-exception as a exception")
				(&trans.wrapper).addError(trans.err, requestFormatError)
			}
		} else {
			trans.err = errors.Errorf("invalid length of format of request fields")
			(&trans.wrapper).addError(trans.err, requestFormatError)
		}
	})
}

func deleteException(w http.ResponseWriter, r *http.Request) {
	executeTransaction(w, r, deleteExceptionEndpoint, admin|manager, func(trans *transaction) {
		if except, isAExcept := trans.endpoint.(*Exception); isAExcept {
			trans.obj, trans.status, trans.err = except.delete(trans.user, trans.permission)
			(&trans.wrapper).addError(trans.err, processError)
		} else {
			trans.err = errors.Errorf("tried to pass a non-exception as a exception")
			(&trans.wrapper).addError(trans.err, requestFormatError)
		}
	})
}

func (except *Exception) update(user domain.User, permission domain.Permission, originalBody string) (generalResp *GeneralResp, status int, err error) {
	generalResp = &GeneralResp{}
	status = http.StatusBadRequest
	return generalResp, status, err
}

func (except *Exception) verify() (verify bool) {
	verify = false
	if len(except.DeviceID) > 0 {
		if len(except.VulnerabilityID) > 0 { // TODO what about sweeping exceptions?
			if except.TypeID > 0 {
				verify = true
			}
		}
	}

	return verify
}

func (except *Exception) create(user domain.User, permission domain.Permission) (generalResp *GeneralResp, status int, err error) {
	return except.createOrUpdate(user, permission)
}

func (except *Exception) createOrUpdate(user domain.User, permission domain.Permission) (generalResp *GeneralResp, status int, err error) {
	generalResp = &GeneralResp{}

	status = http.StatusBadRequest
	if except != nil && permission != nil {
		var dueDate time.Time
		dueDate, err = formatStringToDate(except.DueDate)
		if err == nil {
			_, _, err = Ms.CreateException(except.SourceID, permission.OrgID(), except.TypeID, except.VulnerabilityID,
				except.DeviceID, dueDate, except.Approval, true, except.Port, sord(user.Username()))
			if err == nil {
				status = http.StatusOK
				generalResp.Response = fmt.Sprintf("Exception created for vulnId: %v and deviceId: %v", except.VulnerabilityID, except.DeviceID)
			} else {
				err = errors.New("error while creating exception")
			}
		} else {
			err = errors.Errorf("error while parsing the dates [%s]", err.Error())
		}
	} else {
		err = errors.New("nil parameters while creating exception")
	}

	return generalResp, status, err
}

func (except *Exception) delete(user domain.User, permission domain.Permission) (generalResp *GeneralResp, status int, err error) {
	generalResp = &GeneralResp{}

	status = http.StatusBadRequest
	if except != nil && permission != nil {
		_, _, err = Ms.DisableIgnore(except.SourceID, except.DeviceID, permission.OrgID(), except.VulnerabilityID, except.Port, sord(user.Username()))
		if err == nil {
			status = http.StatusOK
			generalResp.Response = fmt.Sprintf("Exception deleting for vulnId: %v and deviceId: %v", except.VulnerabilityID, except.DeviceID)
		} else {
			err = errors.New("error while deleting exception")
		}

	} else {
		err = errors.New("nil parameters while deleting exception")
	}
	return generalResp, status, err
}

func getAllExceptions(w http.ResponseWriter, r *http.Request) {
	executeTransaction(w, r, getAllExceptionsEndpoint, admin|manager|reporter|reader, func(trans *transaction) {
		if exceptRequest, isExceptRequest := trans.endpoint.(*Exception); isExceptRequest {
			trans.obj, trans.status, trans.totalRecords, trans.err = readAllExceptions(trans.permission, exceptRequest)
			(&trans.wrapper).addError(trans.err, processError)
		} else {
			trans.err = errors.New("tried to pass a non-exception as Exception")
			(&trans.wrapper).addError(trans.err, requestFormatError)
		}
	})
}

func readAllExceptions(permission domain.Permission, exception *Exception) (exceptionDTOs []*Exception, status int, totalRecords int, err error) {
	status = http.StatusBadRequest
	if exception != nil && permission != nil && len(permission.OrgID()) > 0 {
		var updatedDate time.Time
		var createdDate time.Time
		var dueDate time.Time
		var exceptions []domain.Ignore
		var queryData domain.QueryData
		createdDate, err = formatStringToDate(exception.DBCreatedDate)
		dueDate, err = formatStringToDate(exception.DueDate)
		updatedDate, err = formatStringToDate(exception.DBUpdatedDate)
		if err == nil {
			queryData, err = Ms.GetExceptionsLength(exception.SourceID, permission.OrgID(), exception.TypeID, exception.VulnerabilityID, exception.DeviceID, dueDate, exception.Port, exception.Approval, exception.Active, createdDate, updatedDate, exception.UpdatedBy, exception.CreatedBy)
			if err == nil {
				totalRecords = queryData.Length()
				exceptions, err = Ms.GetAllExceptions(exception.Offset, exception.Limit, exception.SourceID, permission.OrgID(), exception.TypeID, exception.VulnerabilityID, exception.DeviceID, dueDate, exception.Port, exception.Approval,
					exception.Active, createdDate, updatedDate, exception.UpdatedBy, exception.CreatedBy, exception.SortedField, exception.SortOrder)
				if err == nil {
					status = http.StatusOK
					exceptionDTOs = toExceptionDtoSlice(exceptions)

				} else {
					err = errors.Errorf("error while obtaining exceptions from database [%s]", err.Error())
				}
			} else {
				err = errors.Errorf("error while obtaining exceptions records length from database [%s]", err.Error())
			}
		} else {
			err = errors.Errorf("error while parsing the dates [%s]", err.Error())
		}
	}
	return exceptionDTOs, status, totalRecords, err
}

func getAllExceptTypes(w http.ResponseWriter, r *http.Request) {
	executeTransaction(w, r, getAllExceptionTypeEndpoints, allAllowed, func(trans *transaction) {
		trans.obj, trans.status, trans.totalRecords, trans.err = readExceptTypes(trans.permission)
		if trans.err != nil {
			trans.err = fmt.Errorf("error while retrieving exception types [%s]", trans.err.Error())
			(&trans.wrapper).addError(trans.err, databaseError)
		}
	})
}

func readExceptTypes(permission domain.Permission) (exceptTypeDtos []*ExceptionType, status int, totalRecords int, err error) {
	status = http.StatusBadRequest
	if permission != nil && len(permission.OrgID()) > 0 {

		var exceptTypes []domain.ExceptionType

		exceptTypes, err = Ms.GetExceptionTypes()
		if err == nil {
			status = http.StatusOK
			if exceptTypes != nil {
				exceptTypeDtos = toExceptTypeDtoSlice(exceptTypes)
			}
		} else {
			err = fmt.Errorf("error while getting sexception types [%s]", err.Error())
		}
	}

	return exceptTypeDtos, status, totalRecords, err
}
