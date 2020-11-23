package endpoints

import (
	"encoding/json"
	"fmt"
	"github.com/nortonlifelock/aegis/pkg/domain"
	"net/http"
	"strings"
	"time"
)

func getLogs(w http.ResponseWriter, r *http.Request) {

	executeTransaction(w, r, getLogsEndpoint, allAllowed, func(trans *transaction) {
		trans.obj, trans.err = gatherLogs(trans.permission, trans.originalBody)
		if trans.err == nil {
			trans.status = http.StatusOK
		} else {
			(&trans.wrapper).addError(trans.err, processError)
		}
	})
}

func gatherLogs(permission domain.Permission, originalBody string) (interface{}, error) {
	var uiLogDateFormat = "2006-01-02T15:04"

	var obj interface{}
	var err error
	var lreq = &logRequest{}
	err = json.Unmarshal([]byte(originalBody), lreq)
	if err == nil {

		var skipLogType = strings.Index(originalBody, "logType") < 0
		var skipFromDate = strings.Index(originalBody, "fromDate") < 0
		var skipToDate = strings.Index(originalBody, "toDate") < 0

		var fromTime time.Time
		var toTime time.Time

		if skipLogType {
			lreq.LogType = -1
		}

		if !skipFromDate {
			fromTime, err = time.Parse(uiLogDateFormat, lreq.FromDate)
		} else {
			fromTime = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC) // set the fromTime to the zero time (aka unbounded)
		}

		// err check to avoid overwriting
		if !skipToDate && err == nil {
			toTime, err = time.Parse(uiLogDateFormat, lreq.ToDate)
		} else {
			toTime = time.Now().Add(time.Hour) // set the toTime to the future (aka unbounded)
		}

		if err == nil {
			var result []domain.DBLog

			result, err = Ms.GetLogsByParams(
				lreq.MethodOfDiscovery,
				lreq.JobType,
				lreq.LogType,
				lreq.JobHistoryID,
				fromTime,
				toTime,
				permission.OrgID(),
			)

			if err == nil {
				obj = toLogDtoSlice(result)
			}
		}

	}
	return obj, err
}

func getLogTypes(w http.ResponseWriter, r *http.Request) {
	executeTransaction(w, r, getAllLogTypesEndpoint, allAllowed, func(trans *transaction) {
		trans.obj, trans.status, trans.totalRecords, trans.err = readLogTypes()
		(&trans.wrapper).addError(trans.err, processError)
	})
}

func readLogTypes() (logsDto []*LogType, status int, totalRecords int, err error) {
	status = http.StatusBadRequest

	var logs []domain.LogType
	logs, err = Ms.GetLogTypes()
	if err == nil {
		status = http.StatusOK
		logsDto = toLogTypeDtoSlice(logs)
	} else {
		err = fmt.Errorf("error while marshalling job histories [%s]", err.Error())
	}

	return logsDto, status, totalRecords, err
}
