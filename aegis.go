package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/pkg/errors"
	"github.com/nortonlifelock/config"
	"github.com/nortonlifelock/database"
	"github.com/nortonlifelock/domain"

	// Shadowed import for registering the jobs into the job runner. There is no need for a direct import here
	// as the jobs are used from a registry for the dispatcher.
	"time"

	_ "github.com/nortonlifelock/implementations"
	"github.com/nortonlifelock/job"
	"github.com/nortonlifelock/log"
)

func main() {
	var err error

	// Setting up config arguments for starting the jobObj runner
	configFile := flag.String("config", "app.json", "The filename of the config to load.")
	configPath := flag.String("cpath", "", "The directory path of the config to load.")
	maxWorkers := flag.Int("workers", 1250, "Number of workers to be used by the jobObj runner")

	jobRunnerWait := flag.Int("run_wait", 60, "The amount of seconds the job runner waits between checking the database for pending jobs")
	jobScheduleWait := flag.Int("schedule_wait", 30, "The amount of seconds the job runner waits between checking the database for scheduled jobs")

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if configFile != nil && configPath != nil && maxWorkers != nil {

		if appConfig, configErr := config.LoadConfig(*configPath, *configFile); configErr == nil {

			var dbconn domain.DatabaseConnection
			if dbconn, err = database.NewConnection(appConfig); err == nil {

				var orgIDToOrgCode map[string]string
				if orgIDToOrgCode, err = getOrgMap(dbconn); err == nil {
					var logger log.Logger
					if logger, err = log.NewLogStream(ctx, dbconn, appConfig); err == nil {

						// Start the dispatcher.
						var dispatcher job.Dispatcher
						if dispatcher, err = job.NewDispatcher(ctx, logger, *maxWorkers); err == nil {

							if err = dispatcher.Run(); err == nil {

								if err = populateAutoStartJobs(dbconn); err == nil {
									err = job.GetDBJobs(ctx, dbconn, logger, appConfig, dispatcher, orgIDToOrgCode, *jobRunnerWait, *jobScheduleWait)
								} else {
									err = errors.Errorf("error while queueing auto start jobs - %s", err.Error())
								}

							} else {
								err = errors.Errorf("error while instantiating worker threads")
							}
						} else {
							err = fmt.Errorf("error while initializing dispatcher - %s", err.Error())
						}
					}
				}
			} else {
				err = fmt.Errorf("error while loading database connection - %s", err.Error())
			}
		} else {
			err = configErr
		}
	} else {
		err = fmt.Errorf("must provide a -config, -cpath, and -cworkers flag")
	}

	if err != nil {
		fmt.Println(err)
	}
}

func getOrgMap(dbconn domain.DatabaseConnection) (orgIDToCode map[string]string, err error) {
	orgIDToCode = make(map[string]string)
	var orgs []domain.Organization
	if orgs, err = dbconn.GetOrganizations(); err == nil {
		for _, org := range orgs {
			orgIDToCode[org.ID()] = org.Code()
		}
	} else {
		err = fmt.Errorf("error while caching organizations - %s", err.Error())
	}

	return orgIDToCode, err
}

func populateAutoStartJobs(dbconn domain.DatabaseConnection) (err error) {
	if _, _, err = dbconn.CleanUp(); err == nil {

		var jobConfigs []domain.JobConfig
		if jobConfigs, err = dbconn.GetAutoStartJobs(); err == nil {

			for index := range jobConfigs {
				jobConfig := jobConfigs[index]

				if jobConfig != nil {

					var baseJob domain.JobRegistration
					if baseJob, err = dbconn.GetJobByID(jobConfig.JobID()); err == nil && baseJob != nil {

						var priority = baseJob.Priority()
						if jobConfig.PriorityOverride() != nil {
							priority = *jobConfig.PriorityOverride()
						}

						var payload string
						if jobConfig.Payload() != nil {
							payload = *jobConfig.Payload()
						}

						_, _, err = dbconn.CreateJobHistory(
							baseJob.ID(),
							jobConfig.ID(),
							domain.JobStatusPending,
							priority,
							"",
							0,
							payload,
							"",
							time.Now().UTC(),
							"RUNNER",
						)
					} else {
						if err == nil {
							err = fmt.Errorf("error while pulling job from database")
						} else {
							err = fmt.Errorf("error while pulling job from database - %s", err.Error())
						}
					}
				} else {
					err = fmt.Errorf("found nil config for autostart jobObj")
				}
			}
		} else {
			err = fmt.Errorf("error while grabbing autostart jobs - %s", err.Error())
		}
	} else {
		err = fmt.Errorf("error while cleaning up jobs - %s", err.Error())
	}

	return err
}
