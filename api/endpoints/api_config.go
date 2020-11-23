package endpoints

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nortonlifelock/aegis/pkg/domain"
	"github.com/pkg/errors"
	"net/http"
	"strings"
	"time"
)

const (
	nilDate = "1970-01-02"
)

func getAuditHistoryForJobConfig(w http.ResponseWriter, r *http.Request) {
	executeTransaction(w, r, getJCAuditHistoryEndpoint, admin|manager, func(trans *transaction) {
		params := mux.Vars(r)
		var id = params[idParam]

		if len(id) > 0 {
			var jobConfigAudits []domain.JobConfigAudit
			if jobConfigAudits, trans.err = Ms.GetJobConfigAudit(id, trans.permission.OrgID()); trans.err == nil {
				trans.obj = jobConfigAudits
				trans.status = http.StatusOK
			} else {
				(&trans.wrapper).addError(trans.err, databaseError)
			}
		} else {
			trans.err = errors.Errorf("empty ID")
			(&trans.wrapper).addError(trans.err, requestFormatError)
		}
	})
}

func deleteJobConfig(w http.ResponseWriter, r *http.Request) {
	executeTransaction(w, r, deleteJConfigEndpoint, admin|manager, func(trans *transaction) {
		params := mux.Vars(r)
		var id = params[idParam]

		if len(id) > 0 {
			var config = &JobConfig{}
			config.ConfigID = id
			trans.obj, trans.status, trans.err = config.delete(trans.user, trans.permission)
			(&trans.wrapper).addError(trans.err, processError)
		} else {
			trans.err = errors.Errorf("empty ID")
			(&trans.wrapper).addError(trans.err, requestFormatError)
		}
	})
}

func getSourceInsByJobID(w http.ResponseWriter, r *http.Request) {
	executeTransaction(w, r, getSourceInByJobIDEndpoint, allAllowed, func(trans *transaction) {
		if trans != nil {
			if config, isAConfig := trans.endpoint.(*JobConfig); isAConfig {
				trans.obj, trans.status, trans.err = config.getSouceInsByJobName(trans.permission)
				(&trans.wrapper).addError(trans.err, processError)
			} else {
				trans.err = errors.Errorf("tried to pass a non-config as a config")
				(&trans.wrapper).addError(trans.err, requestFormatError)
			}
		}
	})
}

func getSourceOutsByJobIDAndSrcIn(w http.ResponseWriter, r *http.Request) {
	executeTransaction(w, r, getSourceOutByJobIDEndpoint, allAllowed, func(trans *transaction) {
		if config, isAConfig := trans.endpoint.(*JobConfig); isAConfig {
			trans.obj, trans.status, trans.err = config.getSourceOutsByJobName(trans.permission)
			(&trans.wrapper).addError(trans.err, processError)
		} else {
			trans.err = errors.Errorf("tried to pass a non-config as a config")
			(&trans.wrapper).addError(trans.err, requestFormatError)
		}
	})
}

func createJobConfig(w http.ResponseWriter, r *http.Request) {
	executeTransaction(w, r, createJConfigEndpoint, admin|manager, func(trans *transaction) {
		if trans.endpoint.verify() {
			if config, isAConfig := trans.endpoint.(*JobConfig); isAConfig {
				trans.obj, trans.status, trans.err = config.create(trans.user, trans.permission)
				(&trans.wrapper).addError(trans.err, processError)
			} else {
				trans.err = errors.Errorf("tried to pass a non-config as a config")
				(&trans.wrapper).addError(trans.err, requestFormatError)
			}
		} else {
			trans.err = errors.Errorf("invalid length of format of request fields")
			(&trans.wrapper).addError(trans.err, requestFormatError)
		}
	})
}

func updateJobConfig(w http.ResponseWriter, r *http.Request) {

	executeTransaction(w, r, updateJConfigEndpoint, admin|manager, func(trans *transaction) {
		if trans.endpoint.verify() {
			if config, isAConfig := trans.endpoint.(*JobConfig); isAConfig {
				trans.obj, trans.status, trans.err = config.update(trans.user, trans.permission, trans.originalBody)
				(&trans.wrapper).addError(trans.err, processError)
			} else {
				trans.err = errors.Errorf("tried to pass a non-config as a config")
				(&trans.wrapper).addError(trans.err, requestFormatError)
			}
		} else {
			trans.err = errors.Errorf("invalid length of format of request fields")
			(&trans.wrapper).addError(trans.err, requestFormatError)
		}
	})
}

func (config *JobConfig) update(user domain.User, permission domain.Permission, originalBody string) (generalResp *GeneralResp, status int, err error) {
	generalResp = &GeneralResp{}

	status = http.StatusBadRequest

	var currentConfig domain.JobConfig
	if config != nil {
		currentConfig, err = Ms.GetJobConfigByID(config.ConfigID, permission.OrgID())

		if err == nil {
			if currentConfig != nil {
				var bodyString = originalBody
				var dataInID = config.DataInSourceConfigID
				var dataOutID = config.DataOutSourceConfigID
				var autoStart = config.AutoStart
				var priority = config.PriorityOverride
				var continuous = config.Continuous
				var wait = config.WaitInSeconds
				var maxInstance = config.MaxInstances

				if strings.Index(bodyString, "data_in_source_config_id") < 0 {
					dataInID = sord(currentConfig.DataInSourceConfigID())
				}
				if strings.Index(bodyString, "data_out_source_config_id") < 0 {
					dataOutID = sord(currentConfig.DataOutSourceConfigID())
				}
				if strings.Index(bodyString, "autostart") < 0 {
					autoStart = currentConfig.AutoStart()
				}
				if strings.Index(bodyString, "priority_override") < 0 {
					priority = iord(currentConfig.PriorityOverride())
				}
				if strings.Index(bodyString, "continuous") < 0 {
					continuous = currentConfig.Continuous()
				}
				if strings.Index(bodyString, "wait_in_seconds") < 0 {
					wait = currentConfig.WaitInSeconds()
				}
				if strings.Index(bodyString, "max_instances") < 0 {
					maxInstance = currentConfig.MaxInstances()
				}

				_, _, err = Ms.UpdateJobConfig(
					config.ConfigID,
					dataInID,
					dataOutID,
					autoStart,
					priority,
					continuous,
					wait,
					maxInstance,
					sord(user.Username()),
					permission.OrgID())

				if err == nil {
					status = http.StatusOK
					generalResp.Response = fmt.Sprintf("job config updated for config id %s", config.ConfigID)
				} else {
					err = errors.New(fmt.Sprintf("error while updating job config"))
				}
			} else {
				err = errors.New(fmt.Sprintf("could not find a job config with id [%s]", config.ConfigID))
			}
		} else {
			err = errors.New(fmt.Sprintf("error while retreiving job config from db [%s]", err.Error()))
		}
	} else {
		err = errors.New("nil config object")
	}

	return generalResp, status, err
}

func (config *JobConfig) create(user domain.User, permission domain.Permission) (generalResp *GeneralResp, status int, err error) {
	generalResp = &GeneralResp{}

	status = http.StatusBadRequest
	if config != nil && permission != nil {
		_, _, err = Ms.CreateJobConfig(
			config.JobID,
			permission.OrgID(),
			config.PriorityOverride,
			config.Continuous,
			config.WaitInSeconds,
			config.MaxInstances,
			config.AutoStart,
			sord(user.Username()),
			config.DataInSourceConfigID,
			config.DataOutSourceConfigID)
		if err == nil {
			status = http.StatusOK
			generalResp.Response = fmt.Sprintf("job config created for job id %d", config.JobID)
		} else {
			err = errors.New(fmt.Sprintf("error while creating job config"))
		}
	} else {
		err = errors.Errorf("nil parameters while creating job config")
	}

	return generalResp, status, err
}

func (config *JobConfig) getSouceInsByJobName(permission domain.Permission) (sourcesDto []*SourceConfig, status int, err error) {
	status = http.StatusBadRequest

	var sources []domain.SourceConfig
	if config != nil && permission != nil {
		sources, err = Ms.GetSourceInsByJobID(config.JobID, permission.OrgID())
		if err == nil {
			status = http.StatusOK
			sourcesDto = toSourcesDtoSlice(sources)
		} else {
			err = fmt.Errorf("error while getting source types [%s]", err.Error())
		}
	} else {
		err = errors.New("nil parameters while getting a sourcein by jobname")
	}

	return sourcesDto, status, err
}

func (config *JobConfig) getSourceOutsByJobName(permission domain.Permission) (sourcesDto []*SourceConfig, status int, err error) {
	status = http.StatusBadRequest

	var sources []domain.SourceConfig
	if config != nil && permission != nil {
		sources, err = Ms.GetSourceOutsByJobID(config.JobID, permission.OrgID())
		if err == nil {
			status = http.StatusOK
			sourcesDto = toSourcesDtoSlice(sources)
		} else {
			err = fmt.Errorf("error while getting source types [%s]", err.Error())
		}
	} else {
		err = errors.New("nil parameters while getting source by jobname")
	}
	return sourcesDto, status, err
}

func (config *JobConfig) delete(user domain.User, permission domain.Permission) (generalResp *GeneralResp, status int, err error) {
	generalResp = &GeneralResp{}

	status = http.StatusBadRequest
	if config != nil && user != nil {
		_, _, err = Ms.DisableJobConfig(config.ConfigID, sord(user.Username()))
		if err == nil {
			status = http.StatusOK
			generalResp.Response = fmt.Sprintf("job config [%s] deleted", config.ConfigID)
		} else {
			err = errors.New(fmt.Sprintf("error while deleting job config"))
		}
	} else {
		err = errors.Errorf("nil parameters while deleting job config")
	}
	return generalResp, status, err
}

func getAllJobConfigs(w http.ResponseWriter, r *http.Request) {
	executeTransaction(w, r, getAllJobConfigsEndpoint, allAllowed, func(trans *transaction) {
		trans.obj, trans.status, trans.totalRecords, trans.err = readJobConfigs(trans.permission)
		if trans.err != nil {
			trans.err = fmt.Errorf("error while retrieving job configs [%s]", trans.err.Error())
			(&trans.wrapper).addError(trans.err, databaseError)
		}
	})
}

func readJobConfigs(permission domain.Permission) (configJobsDto []*JobConfig, status int, totalRecords int, err error) {
	status = http.StatusBadRequest

	var configJobs []domain.JobConfig
	if permission != nil {
		configJobs, err = getAllConfigMappings(Ms, permission.OrgID())
		if err == nil {
			status = http.StatusOK
			configJobsDto = toJobConfigDtoSlice(permission.OrgID(), configJobs, true)
		} else {
			err = fmt.Errorf("error while marshalling job configs [%s]", err.Error())
		}
	} else {
		err = errors.Errorf("nil parameters while reading configs")
	}
	return configJobsDto, status, totalRecords, err
}

type setDataInConfig struct {
	domain.JobConfig
	dataInConfig domain.SourceConfig
}

// DataInConfig overrides the data in source config field
func (jc *setDataInConfig) DataInConfig() domain.SourceConfig {
	return jc.dataInConfig
}

type setDataOutConfig struct {
	domain.JobConfig
	dataOutConfig domain.SourceConfig
}

// DataOutConfig overrides the data out source config field
func (jc *setDataOutConfig) DataOutConfig() domain.SourceConfig {
	return jc.dataOutConfig
}

type setOrganization struct {
	domain.JobConfig
	organization domain.Organization
}

// Organization overrides the organization field
func (jc *setOrganization) Organization() domain.Organization {
	return jc.organization
}

func getAllConfigMappings(ms domain.DatabaseConnection, orgID string) (configs []domain.JobConfig, err error) {
	if len(orgID) > 0 {
		configs, err = ms.GetAllJobConfigs(orgID)

		if err == nil {
			for _, jobConfig := range configs {
				if jobConfig != nil {

					inID := jobConfig.DataInSourceConfigID()
					if inID != nil {
						var sourceIn domain.SourceConfig
						sourceIn, err = ms.GetSourceConfigByID(sord(inID))
						if err == nil {
							if sourceIn != nil {
								jobConfig = &setDataInConfig{
									JobConfig:    jobConfig,
									dataInConfig: sourceIn,
								}
							}
						} else {
							err = errors.New(fmt.Sprintf("error while call the database for DataInConfig "))
							break
						}
					} else {
						err = fmt.Errorf("data in source config id not present for [%v]", jobConfig.ID())
					}

					outID := jobConfig.DataOutSourceConfigID()
					if outID != nil {
						var sourceOut domain.SourceConfig
						sourceOut, err = ms.GetSourceConfigByID(sord(outID))
						if err == nil {
							if sourceOut != nil {
								jobConfig = &setDataOutConfig{
									JobConfig:     jobConfig,
									dataOutConfig: sourceOut,
								}
							}
						} else {
							err = errors.New(fmt.Sprintf("error while call the database for DataOutConfig "))
							break
						}
					} else {
						err = fmt.Errorf("data out source config id not present for [%v]", jobConfig.ID())
					}

					var org domain.Organization
					org, err = ms.GetOrganizationByID(orgID)
					if err == nil {
						if org != nil {
							jobConfig = &setOrganization{
								JobConfig:    jobConfig,
								organization: org,
							}
						}
					} else {
						err = errors.New(fmt.Sprintf("error while call the database for Organization "))
						break
					}
				} else {
					err = errors.New(fmt.Sprintf("nil config job"))
				}
			}
		}

	} else {
		err = errors.New(fmt.Sprintf(" no organization Id is provided"))
	}

	return configs, err
}

func (config *JobConfig) verify() (verify bool) {
	verify = false
	if len(config.DataInSourceConfigID) > 0 {
		if len(config.DataOutSourceConfigID) > 0 {
			if config.JobID > 0 {
				verify = true
			}
		}
	}

	return verify
}

func (jcr *JobConfigRequest) verify() (verify bool) {
	verify = false

	if jcr.Limit > 0 && jcr.Limit <= 500 {
		if jcr.Offset >= 0 {
			verify = true
		}
	}

	return verify
}

func getAllJobConfigsWithOrder(w http.ResponseWriter, r *http.Request) {
	executeTransaction(w, r, postAllJobConfig, admin|manager|reader, func(trans *transaction) {
		if trans.endpoint.verify() {
			if configRequest, isJobConfigRequest := trans.endpoint.(*JobConfigRequest); isJobConfigRequest {
				trans.obj, trans.status, trans.totalRecords, trans.err = readAllJobConfigsWithOrder(trans.permission, configRequest)
				(&trans.wrapper).addError(trans.err, processError)
			} else {
				trans.err = errors.Errorf("tried to pass a non-jobConfigRequest as jobConfigRequest")
				(&trans.wrapper).addError(trans.err, requestFormatError)
			}
		} else {
			trans.err = errors.Errorf("invalid length of format of request fields")
			(&trans.wrapper).addError(trans.err, requestFormatError)
		}
	})
}

func readAllJobConfigsWithOrder(permission domain.Permission, configRequest *JobConfigRequest) (configJobsDto []*JobConfig, status int, totalRecords int, err error) {
	status = http.StatusBadRequest
	if configRequest != nil && permission != nil && len(permission.OrgID()) > 0 {
		var configJobs []domain.JobConfig
		var dates []time.Time
		var queryData domain.QueryData
		dates, err = formatStringToDates(configRequest.UpdatedDate, configRequest.CreatedDate, configRequest.LastJobStart)
		if err == nil {
			queryData, err = Ms.GetJobConfigLength(configRequest.ConfigID, configRequest.JobID, configRequest.DataInSourceConfigID, configRequest.DataInSourceConfigID, configRequest.PriorityOverride, configRequest.Continuous, configRequest.Payload, configRequest.WaitInSeconds, configRequest.MaxInstances, configRequest.AutoStart,
				permission.OrgID(), configRequest.UpdatedBy, configRequest.CreatedBy, dates[0], dates[1], dates[2], configRequest.ConfigID)
			if err == nil {
				totalRecords = queryData.Length()
				configJobs, err = Ms.GetAllJobConfigsWithOrder(configRequest.Offset, configRequest.Limit, configRequest.ConfigID, configRequest.JobID, configRequest.DataInSourceConfigID, configRequest.DataOutSourceConfigID, configRequest.PriorityOverride,
					configRequest.Continuous, configRequest.Payload, configRequest.WaitInSeconds, configRequest.MaxInstances, configRequest.AutoStart, permission.OrgID(), configRequest.UpdatedBy, configRequest.CreatedBy, configRequest.SortedField,
					configRequest.SortOrder, dates[0], dates[1], dates[2], configRequest.ConfigID)
				if err == nil {
					status = http.StatusOK
					configJobsDto = toJobConfigDtoSlice(permission.OrgID(), configJobs, false)
				} else {
					err = errors.Errorf("error while obtaining job configs from database [%s]", err.Error())
				}
			} else {
				err = errors.Errorf("error while obtaining job configs length from database [%s]", err.Error())
			}
		} else {
			err = errors.Errorf("error while parsing the dates [%s]", err.Error())
		}
	}
	return configJobsDto, status, totalRecords, err
}

func formatStringToDate(dateString string) (formattedDate time.Time, err error) {
	var dateFormat = "2006-01-02"
	if dateString == "" {
		formattedDate, _ = time.Parse(dateFormat, nilDate)
	} else {
		formattedDate, err = time.Parse(time.RFC3339, dateString)
	}
	return formattedDate, err
}

func formatStringToDates(dates ...string) (formattedDates []time.Time, err error) {
	var dateFormat = "2006-01-02"
	fd := make([]time.Time, len(dates))
	for index, dateString := range dates {

		if dateString == "" {
			fd[index], _ = time.Parse(dateFormat, nilDate)
		} else {
			fd[index], err = time.Parse(time.RFC3339, dateString)
			if err != nil {
				break
			}
		}
	}
	return fd, err
}
