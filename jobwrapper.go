package job

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/benjivesterby/validator"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/log"
	"github.com/pkg/errors"
)

type jobWrapper struct {
	id        string
	name      string
	payload   string
	orgCode   string
	orgMap    map[string]string
	job       domain.Job
	ctx       context.Context
	db        domain.DatabaseConnection
	lstream   log.Logger
	appconfig domain.Config
}

// Send implements the logger interface on the job wrapper so whenever the job creates a log, the job name is appended to the log
func (wrapper *jobWrapper) Send(inLog log.Log) {
	//inLog.Text = fmt.Sprintf("%s - %s", wrapper.name, inLog.Text)
	inLog.JobID = wrapper.id
	inLog.Job = wrapper.name
	inLog.OrgCode = wrapper.orgCode
	wrapper.lstream.Send(inLog)
}

// Execute initializes the process method of the job contained in the jobWrapper
func (wrapper *jobWrapper) Execute() (err error) {
	if validator.IsValid(wrapper) {

		var jobError error
		var cancelled = false

		defer func() {
			if jobError == nil {
				if cancelled {
					// Update the job to being CANCELED
					_, _, _ = wrapper.db.UpdateJobHistoryStatus(wrapper.id, domain.JobStatusCancelled)
					wrapper.Send(log.Info("Job cancelled"))
				} else {
					// Update the job to being COMPLETED
					_, _, _ = wrapper.db.UpdateJobHistoryStatus(wrapper.id, domain.JobStatusCompleted)
					wrapper.Send(log.Info("Job completed successfully"))
				}
			} else {
				// Update the job to being ERROR
				_, _, _ = wrapper.db.UpdateJobHistoryStatus(wrapper.id, domain.JobStatusError)
				wrapper.Send(log.Warning("Job failed to complete, exited with ERROR status", jobError))
			}
		}()

		func() {
			for {
				select {
				case <-wrapper.ctx.Done():
					cancelled = true
					return
				default:

					_, _, _ = wrapper.db.PulseJob(wrapper.id)

					// Update the job to being in progress
					if _, _, err = wrapper.db.UpdateJobHistoryStatus(wrapper.id, domain.JobStatusInProgress); err == nil {

						var jobConfig domain.JobConfig
						var inSources []domain.SourceConfig
						var outSources []domain.SourceConfig
						if jobConfig, inSources, outSources, err = wrapper.loadConfigs(); err == nil {
							wrapper.orgCode = wrapper.orgMap[jobConfig.OrganizationID()]
							// Determine job start time
							var startTime = time.Now().UTC()

							wrapper.lstream.Send(log.Debugf("Beginning Processing of Job [%v]", wrapper.id))

							cancelCtx, cancelFunc := context.WithCancel(wrapper.ctx)
							// Execute the process method using the jobWrapper
							jobError = wrapper.job.Process(cancelCtx, wrapper.id, wrapper.appconfig, wrapper.db, wrapper, wrapper.payload, jobConfig, inSources, outSources)
							cancelFunc()

							if jobError == nil {
								// Update the this config last run time so that it limits the tickets pulled from each
								// jira request to only those not previously seen
								_, _, _ = wrapper.db.UpdateJobConfigLastRun(jobConfig.ID(), startTime)
							}

							select {
							case <-wrapper.ctx.Done():
								cancelled = true
								return
							default:
								// Handle the continuous flag and the wait in seconds
								if jobConfig.Continuous() {

									if jobError != nil {
										wrapper.Send(log.Errorf(jobError, "error during job processing"))
									}

									wrapper.Send(log.Info(fmt.Sprintf("Sleeping for %v seconds", jobConfig.WaitInSeconds())))
									time.Sleep(time.Second * time.Duration(jobConfig.WaitInSeconds()))
								} else if !jobConfig.Continuous() {
									return
								}
							}
						} else {
							err = errors.Errorf("unable to load configurations for job [%v] - [%v]", wrapper.id, err)
							jobError = err
							return
						}
					}
				}
			}

		}()

	} else {
		err = errors.Errorf("invalid jobWrapper for job [%v]", wrapper.id)
	}

	return err
}

// Job returns the job object as the job interface for use by the caller
func (wrapper *jobWrapper) Job() (jb domain.Job) {
	return wrapper.job
}

// loadConfigs loads fresh source configs from the database. This method is only used in continuous jobs
func (wrapper *jobWrapper) loadConfigs() (config domain.JobConfig, inSources []domain.SourceConfig, outSources []domain.SourceConfig, err error) {
	inSources = make([]domain.SourceConfig, 0)
	outSources = make([]domain.SourceConfig, 0)

	// Update the job config so that each loop it's correct
	if config, err = wrapper.db.GetJobConfigByJobHistoryID(wrapper.id); err == nil {

		if len(sord(config.DataInSourceConfigID())) > 0 && len(sord(config.DataOutSourceConfigID())) > 0 {
			inSourceIDs := strings.Split(sord(config.DataInSourceConfigID()), ",")
			outSourceIDs := strings.Split(sord(config.DataOutSourceConfigID()), ",")

			for _, ID := range inSourceIDs {
				var inSourceConfig domain.SourceConfig
				if inSourceConfig, err = wrapper.db.GetSourceConfigByID(ID); err == nil {
					if inSourceConfig != nil {
						inSources = append(inSources, inSourceConfig)
					} else {
						err = fmt.Errorf("no source config found for [%v]", ID)
						break
					}
				} else {
					break
				}
			}

			if err == nil {
				for _, ID := range outSourceIDs {
					var outSourceConfig domain.SourceConfig
					if outSourceConfig, err = wrapper.db.GetSourceConfigByID(ID); err == nil {
						if outSourceConfig != nil {
							outSources = append(outSources, outSourceConfig)
						} else {
							err = fmt.Errorf("no source config found for [%v]", ID)
							break
						}
					} else {
						break
					}
				}
			}
		} else {
			err = fmt.Errorf("job history [%v] had an empty data in config or data out config", wrapper.id)
		}
	} else {
		err = errors.Errorf("error occurred while loading config for job [%v] | [error: %s]", wrapper.id, err.Error())
	}

	return config, inSources, outSources, err
}

// Validate checks the jobWrapper struct to ensure that everything is valid as well as the job itself
func (wrapper *jobWrapper) Validate() (valid bool) {

	if len(wrapper.id) > 0 {
		if len(wrapper.name) > 0 {
			if wrapper.ctx != nil {
				if wrapper.db != nil {
					if wrapper.lstream != nil {
						// Attempt type assertion for validity
						valid = validator.IsValid(wrapper.job)
					}
				}
			}
		}
	}

	return valid
}
