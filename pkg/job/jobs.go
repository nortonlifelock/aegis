package job

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/benjivesterby/validator"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/log"
	"github.com/pkg/errors"
)

// GetDBJobs pulls pending jobs from the job queue, as kicks of the threads for running scheduled jobs and processing running jobs
func GetDBJobs(ctx context.Context, db domain.DatabaseConnection, lstream log.Logger, appconfig domain.Config, dispatcher Dispatcher, orgIDToOrg map[string]domain.Organization, sleepInSecondsPendingJobs int, sleepInSecondsScheduledJobs int) (err error) {
	// TODO: Set this up so that when it handles the panic it also closes the context for everything below then restarts it
	defer handleRoutinePanic(lstream)

	go manageActiveJobs(ctx, db, lstream, dispatcher)
	go runScheduledJobs(ctx, db, lstream, sleepInSecondsScheduledJobs)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			lstream.Send(log.Debug("Job Runner: Loading Job Queue from database"))

			var pendingJobs, runningJobs []domain.JobHistory
			if pendingJobs, err = db.GetJobQueueByStatusID(domain.JobStatusPending); err == nil {
				if runningJobs, err = db.GetJobQueueByStatusID(domain.JobStatusInProgress); err == nil {

					// If the queue contains data then build the jobs
					if pendingJobs != nil && len(pendingJobs) > 0 {
						lstream.Send(log.Infof("Job Runner: [%d] jobs found for processing", len(pendingJobs)))

						var configIDToInstanceCount = countRunningJobs(runningJobs)

						for index := range pendingJobs {
							if pendingJobs[index].MaxInstances() == 0 || configIDToInstanceCount[pendingJobs[index].ConfigID()] < pendingJobs[index].MaxInstances() {
								configIDToInstanceCount[pendingJobs[index].ConfigID()]++
								prepareJobForQueue(db, lstream, appconfig, pendingJobs[index], dispatcher, orgIDToOrg)
							} else {
								lstream.Send(log.Warningf(err, "Job Runner: waiting to start [%s] as its config currently has %d running instances which is its maximum", pendingJobs[index].ID(), configIDToInstanceCount[pendingJobs[index].ConfigID()]))
							}
						}

					} else {
						lstream.Send(log.Debugf("Job Runner: No Jobs to Process"))
					}
				} else {
					lstream.Send(log.Errorf(err, "Job Runner: error while gathering processing jobs"))
				}
			} else {
				lstream.Send(log.Errorf(err, "Job Runner: Error while retrieving job queue"))
			}
		}

		lstream.Send(log.Infof("Job Runner: Sleeping for %d seconds - [%d goroutines]", sleepInSecondsPendingJobs, runtime.NumGoroutine()))
		time.Sleep(time.Second * time.Duration(sleepInSecondsPendingJobs))
	}
}

func countRunningJobs(runningJobs []domain.JobHistory) map[string]int {
	var configIDToInstanceCount = make(map[string]int)

	for _, job := range runningJobs {
		configIDToInstanceCount[job.ConfigID()]++
	}

	return configIDToInstanceCount
}

func manageActiveJobs(ctx context.Context, db domain.DatabaseConnection, lstream log.Logger, spatcher Dispatcher) {
	defer func() {
		select {
		case <-ctx.Done():
			lstream.Send(log.Info("Job Runner: Closing Job Runner Context Manager for Manage Active Jobs"))
			return
		default:
			handleRoutinePanic(lstream)

			go manageActiveJobs(ctx, db, lstream, spatcher)
		}
	}()

	for {

		select {
		case <-ctx.Done():
			return
		default:
			// CallA DAL
			var jobs []domain.JobHistory
			var err error

			if jobs, err = db.GetCancelledJobs(); err == nil {
				for index := range jobs {

					// Send the id of the job to the dispatcher so it can be cancelled
					spatcher.Cancel(jobs[index].ID())
				}
			} else {
				lstream.Send(log.Error("Job Runner: Error getting jobs to cancel from the database", err))
			}

			time.Sleep(time.Second)
		}
	}
}

func runScheduledJobs(ctx context.Context, db domain.DatabaseConnection, lstream log.Logger, sleepInSeconds int) {
	defer func() {
		select {
		case <-ctx.Done():
			lstream.Send(log.Info("Job Runner: Closing Job Runner Context Manager for Scheduled Jobs"))
			return
		default:
			handleRoutinePanic(lstream)

			go runScheduledJobs(ctx, db, lstream, sleepInSeconds)
		}
	}()

	var lastChecked = time.Now()

	for {
		select {
		case <-ctx.Done():
			return
		default:

			lstream.Send(log.Info("Job Runner: Checking for scheduled jobs ready to execute"))

			// CallA DAL
			var schedules []domain.JobSchedule
			var err error
			if schedules, err = db.GetScheduledJobsToStart(lastChecked); err == nil {
				lastChecked = time.Now()

				var scheduleLen = len(schedules)

				if scheduleLen > 0 {

					lstream.Send(log.Infof("Job Runner: Queueing [%v] Scheduled Jobs", scheduleLen))

					for index := range schedules {
						createJobHistoryForSchedule(db, schedules, index, lstream)
					}
				}
			} else {
				lstream.Send(log.Error("Job Runner: Error getting job schedules from the database", err))
			}

			time.Sleep(time.Duration(sleepInSeconds) * time.Second)
		}
	}
}

func createJobHistoryForSchedule(db domain.DatabaseConnection, schedules []domain.JobSchedule, index int, lstream log.Logger) {
	var err error
	if _, _, err = db.SetScheduleLastRun(schedules[index].ID()); err == nil {

		var config domain.JobConfig
		if config, err = db.GetJobConfig(schedules[index].ConfigID()); err == nil {
			if config != nil {

				var baseJob domain.JobRegistration
				if baseJob, err = db.GetJobByID(config.JobID()); err == nil && baseJob != nil {
					var priority = baseJob.Priority()
					if config.PriorityOverride() != nil {
						priority = iord(config.PriorityOverride())
					}

					_, _, err = db.CreateJobHistory(
						baseJob.ID(),
						config.ID(),
						domain.JobStatusPending,
						priority,
						"",
						0,
						sord(config.Payload()),
						"",
						time.Now().UTC(),
						"RUNNER",
					)

					if err != nil {
						lstream.Send(log.Errorf(err, "error while creating job history for scheduled job with config id [%v]", config.ID()))
					}
				} else {
					if err != nil {
						err = fmt.Errorf("error while pulling job from database - %s", err.Error())
					}
					lstream.Send(log.Error("Error in the scheduler", err))
				}
			} else {
				err = fmt.Errorf("found nil Config for autostart job")
				lstream.Send(log.Error("Error in the scheduler", err))
			}
		} else {
			lstream.Send(log.Error("Error in the scheduler", err))
		}
	} else {
		lstream.Send(log.Error("Job Runner: Error getting jobs to cancel from the database", err))
	}
}

func prepareJobForQueue(db domain.DatabaseConnection, lstream log.Logger, appconfig domain.Config, jh domain.JobHistory, dispatcher Dispatcher, orgIDToOrg map[string]domain.Organization) {
	defer handleRoutinePanic(lstream)

	var err error

	// TODO: Add a re-queue flag to job histories so that we can force a re-queue of a job after it's been added to pending status

	if jh != nil {
		lstream.Send(log.Debugf("Job Runner: Loading Job Struct information for jobId [%s]", jh.ID()))

		var baseJob domain.JobRegistration
		if baseJob, err = db.GetJobByID(jh.JobID()); err == nil && baseJob != nil {
			lstream.Send(log.Debug("Job Runner: Loaded Job Struct"))

			var job domain.Job
			if job, err = getJob(baseJob.GoStruct()); err == nil && job != nil {
				lstream.Send(log.Debugf("Job Runner: Pushing Job [%s] onto the job channel for processing", baseJob.GoStruct()))

				// NOTE: the ctx is defined by the dispatcher so it shouldn't be defined here
				if err = dispatcher.Queue(&jobWrapper{
					id:        jh.ID(),
					name:      baseJob.GoStruct(),
					lstream:   lstream,
					appconfig: appconfig,
					db:        db,
					job:       job,
					payload:   jh.Payload(),
					orgMap:    orgIDToOrg,
				}); err != nil {
					lstream.Send(log.Error("Job Runner: Error while queuing", err))
				}
			} else {
				lstream.Send(log.Errorf(err, "Job Runner: Unable to find job [%s] in registry", baseJob.GoStruct()))
			}
		}
	} else {
		// TODO: nil job history
	}
}

type snsClient struct {
	SNSID string `json:"sns_id"`
}

func grabSNSIDFromOrgPayloadIfPresent(orgPayload string) (snsID string, exists bool) {
	parse := &snsClient{}
	_ = json.Unmarshal([]byte(orgPayload), parse)
	snsID = parse.SNSID
	exists = len(snsID) > 0
	return snsID, exists
}

// JOB REGISTRY

var (
	jobsMu      sync.RWMutex
	jobRegistry = make(map[string]domain.Job)
)

func getJob(name string) (job domain.Job, err error) {
	jobsMu.RLock()

	if registration, ok := jobRegistry[name]; ok {

		if job, ok = reflect.New(reflect.TypeOf(registration).Elem()).Interface().(domain.Job); ok {

			if !validator.IsValid(job) {
				err = errors.Errorf("Job Runner: Job [%s] is invalid", name)
			}
		}
	} else {
		err = errors.Errorf("Job [%s] is not registered with the job runner", name)
	}

	jobsMu.RUnlock()

	return job, err
}

// Register makes a job available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, job domain.Job) {
	jobsMu.Lock()
	defer jobsMu.Unlock()

	// Register the job
	if _, dup := jobRegistry[name]; !dup {
		jobRegistry[name] = job
	} else {
		panic(fmt.Sprintf("Job Runner: Register called twice for job [%s]", name))
	}
}
