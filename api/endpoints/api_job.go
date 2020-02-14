package endpoints

import (
	"github.com/nortonlifelock/domain"
	"github.com/pkg/errors"
	"net/http"
)

func getAllJobs(w http.ResponseWriter, r *http.Request) {
	executeTransaction(w, r, getAllJobsEndpoint, allAllowed, func(trans *transaction) {
		trans.obj, trans.status, trans.totalRecords, trans.err = readJobs()
		(&trans.wrapper).addError(trans.err, processError)
	})
}

func readJobs() (jobsDto []*Job, status int, totalRecords int, err error) {
	status = http.StatusBadRequest

	var jobs []domain.JobRegistration
	jobs, err = Ms.GetJobs()
	if err == nil {
		status = http.StatusOK
		jobsDto = toJobDtoSlice(jobs)
	} else {
		err = errors.Errorf("error while marshalling job histories [%s]", err.Error())
	}

	return jobsDto, status, totalRecords, err
}

func (history *Job) create(user domain.User, permission domain.Permission) (generalResp *GeneralResp, status int, err error) {
	return
}

func (history *Job) update(user domain.User, permission domain.Permission, body string) (generalResp *GeneralResp, status int, err error) {
	return
}
func (history *Job) delete(user domain.User, permission domain.Permission) (generalResp *GeneralResp, status int, err error) {
	return
}
