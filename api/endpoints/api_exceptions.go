package endpoints

import (
	"fmt"
	"github.com/nortonlifelock/domain"
	"github.com/pkg/errors"
	"net/http"
)

func (except *Exception) update(user domain.User, permission domain.Permission, originalBody string) (generalResp *GeneralResp, status int, err error) {
	generalResp = &GeneralResp{}
	status = http.StatusBadRequest
	return generalResp, status, err
}

func (except *Exception) verify() (verify bool) {
	verify = false
	//if len(except.DeviceID) > 0 {
	//	if len(except.VulnerabilityID) > 0 { // TODO what about sweeping exceptions?
	//		if except.TypeID > 0 {
	//			verify = true
	//		}
	//	}
	//}
	verify = true

	return verify
}

func (except *Exception) create(user domain.User, permission domain.Permission) (generalResp *GeneralResp, status int, err error) {
	return except.createOrUpdate(user, permission)
}

func (except *Exception) createOrUpdate(user domain.User, permission domain.Permission) (generalResp *GeneralResp, status int, err error) {
	generalResp = &GeneralResp{}

	status = http.StatusBadRequest
	if except != nil && permission != nil {
		//var dueDate time.Time
		//dueDate, err = formatStringToDate(except.DueDate)
		//if err == nil {
		//	_, _, err = Ms.CreateException(except.SourceID, permission.OrgID(), except.TypeID, except.VulnerabilityID,
		//		except.DeviceID, dueDate, except.Approval, true, except.Port, sord(user.Username()))
		//	if err == nil {
		//		status = http.StatusOK
		//		generalResp.Response = fmt.Sprintf("Exception created for vulnId: %v and deviceId: %v", except.VulnerabilityID, except.DeviceID)
		//	} else {
		//		err = errors.New("error while creating exception")
		//	}
		//} else {
		//	err = errors.Errorf("error while parsing the dates [%s]", err.Error())
		//}
	} else {
		err = errors.New("nil parameters while creating exception")
	}

	return generalResp, status, err
}

func (except *Exception) delete(user domain.User, permission domain.Permission) (generalResp *GeneralResp, status int, err error) {
	generalResp = &GeneralResp{}

	status = http.StatusBadRequest
	if except != nil && permission != nil {
		//_, _, err = Ms.DisableIgnore(except.SourceID, except.DeviceID, permission.OrgID(), except.VulnerabilityID, except.Port, sord(user.Username()))
		//if err == nil {
		//	status = http.StatusOK
		//	generalResp.Response = fmt.Sprintf("Exception deleting for vulnId: %v and deviceId: %v", except.VulnerabilityID, except.DeviceID)
		//} else {
		//	err = errors.New("error while deleting exception")
		//}

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

		var exceptions []domain.ExceptedDetection
		if exceptions, err = Ms.GetExceptionDetections(
			exception.Offset,
			exception.Limit,
			permission.OrgID(),
			exception.SortedField,
			exception.SortOrder,
			exception.Title,
			exception.IP,
			exception.Hostname,
			exception.VulnerabilityID,
			exception.Approval,
			exception.Expires,
			exception.AssignmentGroup,
			exception.OS,
			exception.OSRegex,
			exception.IgnoreTypeID,
		); err == nil {

			var queryData domain.QueryData
			queryData, _ = Ms.GetExceptionsLength(exception.Offset,
				exception.Limit,
				permission.OrgID(),
				exception.SortedField,
				exception.SortOrder,
				exception.Title,
				exception.IP,
				exception.Hostname,
				exception.VulnerabilityID,
				exception.Approval,
				exception.Expires,
				exception.AssignmentGroup,
				exception.OS,
				exception.OSRegex,
				exception.IgnoreTypeID)

			if queryData != nil {
				totalRecords = queryData.Length()
			}

			exceptionDTOs = toExceptionDtoSlice(exceptions)
			status = http.StatusOK
		} else {
			err = fmt.Errorf("error while loading excepted detections - %s", err.Error())
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
