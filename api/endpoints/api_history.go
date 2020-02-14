package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/nortonlifelock/domain"
	"github.com/pkg/errors"
)

func getJobHistoryByID(w http.ResponseWriter, r *http.Request) {

	executeTransaction(w, r, getHistoryByIDEndpoint, admin|manager|reporter|reader, func(trans *transaction) {
		params := mux.Vars(r)
		var id = params[idParam]
		if len(id) > 0 {
			trans.obj, trans.err = retrieveJobHistory(id)
			if trans.err != nil {
				trans.err = errors.Errorf("error while retrieving job history [%s]", trans.err.Error())
				(&trans.wrapper).addError(trans.err, processError)
			}
		} else {
			trans.err = errors.New("empty id passed")
			(&trans.wrapper).addError(trans.err, requestFormatError)
		}
	})
}

func getAllJobHistories(w http.ResponseWriter, r *http.Request) {

	executeTransaction(w, r, getHistoriesEndpoint, admin|manager|reporter|reader, func(trans *transaction) {
		if trans.endpoint.verify() {
			if histories, isHistory := trans.endpoint.(*Histories); isHistory {
				trans.obj, trans.status, trans.totalRecords, trans.err = readHistories(trans.permission, histories)
				(&trans.wrapper).addError(trans.err, processError)
			} else {
				trans.err = errors.Errorf("tried to pass a non-history as a history")
				(&trans.wrapper).addError(trans.err, requestFormatError)
			}
		} else {
			trans.err = errors.Errorf("invalid length of format of apiRequest fields")
			(&trans.wrapper).addError(trans.err, requestFormatError)
		}
	})
}

func deleteJobHistory(w http.ResponseWriter, r *http.Request) {
	executeTransaction(w, r, deleteHistoryEndpoint, admin, func(trans *transaction) {
		params := mux.Vars(r)
		var id = params[idParam]
		history := &JobHistory{}
		history.JobHistID = id
		trans.obj, trans.status, trans.err = history.delete(trans.user, trans.permission)
		(&trans.wrapper).addError(trans.err, processError)
	})
}

func updateJobHistory(w http.ResponseWriter, r *http.Request) {
	executeTransaction(w, r, updateHistoryEndpoint, admin|manager, func(trans *transaction) {
		if trans.endpoint.verify() {
			if history, isHistory := trans.endpoint.(*JobHistory); isHistory {
				trans.obj, trans.status, trans.err = history.update(trans.user, trans.permission, "")
				(&trans.wrapper).addError(trans.err, processError)
			} else {
				trans.err = errors.Errorf("tried to pass a non-history as a history")
				(&trans.wrapper).addError(trans.err, requestFormatError)
			}
		} else {
			trans.err = errors.Errorf("invalid length of format of apiRequest fields")
			(&trans.wrapper).addError(trans.err, requestFormatError)
		}
	})
}

func createJobHistory(w http.ResponseWriter, r *http.Request) {
	executeTransaction(w, r, createHistoryEndpoint, admin|manager, func(trans *transaction) {
		if trans.endpoint.verify() {
			if history, isHistory := trans.endpoint.(*JobHistory); isHistory {
				trans.obj, trans.status, trans.err = history.create(trans.user, trans.permission)
				(&trans.wrapper).addError(trans.err, processError)
			} else {
				trans.err = errors.Errorf("tried to pass a non-history as a history")
				(&trans.wrapper).addError(trans.err, requestFormatError)
			}
		} else {
			trans.err = errors.Errorf("invalid length of format of apiRequest fields")
			(&trans.wrapper).addError(trans.err, requestFormatError)
		}
	})
}

func (history *JobHistory) create(user domain.User, permission domain.Permission) (generalResp *GeneralResp, status int, err error) {
	generalResp = &GeneralResp{}

	status = http.StatusBadRequest

	var configs domain.JobConfig
	configs, err = Ms.GetJobConfigByID(history.ConfigID, permission.OrgID())

	if err == nil {
		if configs != nil {
			if len(history.ConfigID) > 0 {
				_, _, err = Ms.CreateJobHistory(
					configs.JobID(),
					configs.ID(),
					domain.JobStatusPending,
					iord(configs.PriorityOverride()),
					"",
					0,
					history.Payload,
					"",
					time.Now().UTC(),
					sord(user.Username()),
				)
				if err == nil {
					status = http.StatusOK
					generalResp.Response = fmt.Sprintf("Job history created using config id [%s]", history.ConfigID)
				} else {
					err = errors.New("error while creating job history")
				}
			} else {
				err = errors.New("both \"config_id\" and \"status\" must be nonzero")
			}
		} else {
			err = errors.Errorf("could not find config with id [%s]", history.ConfigID)
		}
	} else {
		err = errors.New("error while retreiving job config from database")
	}

	return generalResp, status, err
}

func (history *JobHistory) delete(user domain.User, permission domain.Permission) (generalResp *GeneralResp, status int, err error) {
	generalResp = &GeneralResp{}

	status = http.StatusBadRequest

	if len(history.JobHistID) > 0 {
		_, _, err = Ms.UpdateJobHistoryStatusDetailed(history.JobHistID, domain.JobStatusCancelled, sord(user.Username()))
		if err == nil {
			status = http.StatusOK
			generalResp.Response = fmt.Sprintf("job history %s cancelled", history.JobHistID)
		} else {
			err = errors.New("error while deleting job history")
		}
	} else {
		err = errors.New("a \"history_id\" greater than zero must be provided in the apiRequest body")
	}

	return generalResp, status, err
}

func (history *JobHistory) update(user domain.User, permission domain.Permission, originalBody string) (generalResp *GeneralResp, status int, err error) {
	generalResp = &GeneralResp{}
	status = http.StatusBadRequest
	generalResp.Response = "Modification of Job Histories not allowed. If you want to modify a job, cancel the old job and create a new one in it's place"
	return generalResp, status, err
}

func readHistories(permission domain.Permission, histories *Histories) (jobHistoriesDto []*JobHistory, status int, totalRecords int, err error) {
	status = http.StatusBadRequest

	var jobHistories []domain.JobHistory
	var queryData domain.QueryData
	queryData, err = Ms.GetJobHistoryLength(histories.JobID, histories.ConfigID, histories.StatusID, histories.Payload, permission.OrgID())
	if err == nil {
		totalRecords = queryData.Length()
		jobHistories, err = Ms.GetJobHistories(histories.Offset, histories.Limit, histories.JobID, histories.ConfigID, histories.StatusID, histories.Payload, permission.OrgID())
		if err == nil {
			status = http.StatusOK
			jobHistoriesDto = toHistoriesDtoSlice(jobHistories)
		} else {
			err = errors.Errorf("error while obtaining job histories from database [%s]", err.Error())
		}
	} else {
		err = errors.Errorf("error while obtaining job histories length from database [%s]", err.Error())
	}

	return jobHistoriesDto, status, totalRecords, err
}

func retrieveJobHistory(id string) (byteHistory []byte, err error) {
	var history domain.JobHistory
	history, err = Ms.GetJobHistoryByID(id)
	if err == nil {
		if history != nil {
			byteHistory, err = json.Marshal(history)
			if err == nil {
				//do nothing
			} else {
				err = errors.Errorf("error while marshalling history [%s]", err.Error())
			}
		} else {
			err = errors.Errorf("could not find a history with id [%s]", id)
		}
	} else {
		err = errors.Errorf("error while retrieving history from database [%s]", err.Error())
	}
	return byteHistory, err
}

func (history *Histories) create(user domain.User, permission domain.Permission) (generalResp *GeneralResp, status int, err error) {
	return
}

func (history *Histories) update(user domain.User, permission domain.Permission, body string) (generalResp *GeneralResp, status int, err error) {
	return
}
func (history *Histories) delete(user domain.User, permission domain.Permission) (generalResp *GeneralResp, status int, err error) {
	return
}

func (history *Histories) verify() (verify bool) {
	verify = false

	if history.Limit > 0 && history.Limit <= 500 {
		if history.Offset >= 0 {
			verify = true
		}
	}

	return verify
}

func (history *JobHistory) verify() (verify bool) {
	// TODO what sorts of constraints can be placed on the payload?
	verify = false

	if history.JobID > 0 {
		if len(history.ConfigID) > 0 {
			if len(history.JobHistID) >= 0 {
				if history.StatusID > 0 {
					verify = true
				}
			}
		}
	}

	return verify
}
