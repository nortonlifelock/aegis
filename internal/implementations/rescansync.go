package implementations

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nortonlifelock/aegis/internal/integrations"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/log"
	"strings"
	"time"
)

// ScanSyncJob is responsible for monitoring the ScanSummary table in the database, and updating the status using information from the scanner API
// the job should be marked as autostart and continuous
type ScanSyncJob struct {
	id          string
	payloadJSON string
	payload     *ScanSyncPayload
	ctx         context.Context
	db          domain.DatabaseConnection
	lstream     log.Logger
	appconfig   domain.Config
	config      domain.JobConfig
	insource    domain.SourceConfig
	outsource   domain.SourceConfig
}

type ScanSyncPayload struct {
	// these scheduled scans are pushed to Qualys AFTER the unfinished scans from the db
	// the scheduled scans in here may overlap with the unfinished scans from the db, but the driver should be able to filter out the additionals
	// using the scan titles
	ScheduledScanPayloads []string `json:"scheduled_scans"`
}

// Process monitors unfinished scans in the database, and queries the scanners to keep the status of the scans in the database up-to-date
// if the scanner reports the scan as finished, this job queues up a job history for a scan close job that will process the results of the scan
func (job *ScanSyncJob) Process(ctx context.Context, id string, appconfig domain.Config, db domain.DatabaseConnection, lstream log.Logger, payload string, jobConfig domain.JobConfig, inSource []domain.SourceConfig, outSource []domain.SourceConfig) (err error) {

	var ok bool
	if job.ctx, job.id, job.appconfig, job.db, job.lstream, job.payloadJSON, job.config, job.insource, job.outsource, ok = validInputs(ctx, id, appconfig, db, lstream, payload, jobConfig, inSource, outSource); ok {

		job.payload = &ScanSyncPayload{}
		if err = json.Unmarshal([]byte(job.payloadJSON), job.payload); err == nil {
			var scanner integrations.Vscanner

			var sourceName = job.insource.Source()

			job.lstream.Send(log.Debugf("Establishing %s session", sourceName))
			scanner, err = integrations.NewVulnScanner(job.ctx, sourceName, job.db, job.lstream, job.appconfig, job.insource)

			if err == nil {

				var baseJob domain.JobRegistration
				if baseJob, err = job.getBaseJob(); err == nil {
					var unfinishedScans []domain.ScanSummary
					unfinishedScans, err = job.db.GetUnfinishedScanSummariesBySourceConfigOrgID(job.insource.ID(), job.config.OrganizationID())

					if err == nil {
						if (unfinishedScans != nil && len(unfinishedScans) > 0) || len(job.payload.ScheduledScanPayloads) > 0 {
							job.lstream.Send(log.Infof("Found [%v] Scans to Sync", len(unfinishedScans)+len(job.payload.ScheduledScanPayloads)))

							var scanData <-chan domain.Scan
							if scanData = scanner.Scans(job.ctx, job.pushScansOntoChannel(unfinishedScans)); err == nil {

								for {
									select {
									case <-job.ctx.Done():
										return
									case scanToCheck, ok := <-scanData:
										if ok {
											var correspondingScanSummary domain.ScanSummary
											if correspondingScanSummary, err = job.getCorrespondingScanSummary(scanToCheck, unfinishedScans); err == nil {
												job.processSyncedScan(scanToCheck, correspondingScanSummary, baseJob)
											} else {
												job.lstream.Send(log.Errorf(err, "could not find scan summary"))
											}
										} else {
											return
										}
									}
								}
							} else {
								job.lstream.Send(log.Errorf(err, "Error while syncing scans for %s - %s", sourceName, err.Error()))
							}
						} else {
							job.lstream.Send(log.Debugf("No unfinished scans found for %s", sourceName))
						}
					} else {
						job.lstream.Send(log.Errorf(err, "Error while grabbing unfinished scans from the database [%s]", err.Error()))
					}
				} else {
					job.lstream.Send(log.Error("error while gathering base job", err))
				}
			} else {
				job.lstream.Send(log.Error("Error while creating vulnerability scanner", err))
			}
		} else {
			job.lstream.Send(log.Errorf(err, "error while unmarshalling job history payload"))
		}

	} else {
		err = fmt.Errorf("input validation failed")
	}

	return err
}

// Updates the status of the scan in the database, and creates a scan close job for the scan if its status was marked as finished by the scanner
func (job *ScanSyncJob) processSyncedScan(scan domain.Scan, correspondingScanSummary domain.ScanSummary, baseJob domain.JobRegistration) {
	if status, err := scan.Status(); err == nil {
		job.lstream.Send(log.Infof("Updating database scan [Id: %v] with status [%s]", scan.ID(), status))

		// Save the updated scan summary
		if _, _, err = job.db.SaveScanSummary(scan.ID(), status); err == nil {

			if strings.ToLower(status) == strings.ToLower(domain.ScanFINISHED) {
				job.createScanCloseJob(scan, correspondingScanSummary, baseJob)
			}

		} else {
			job.lstream.Send(log.Errorf(err, "Error while updating the scan summary [ID: %v] with status [%s]", scan.ID(), status))
		}
	} else {
		job.lstream.Send(log.Errorf(err, "error while gathering the scan status of [%v]", scan.ID()))
	}
}

// Builds and creates a job history for the scan close job that will be responsible for handling the scan results
func (job *ScanSyncJob) createScanCloseJob(scan domain.Scan, correspondingScanSummary domain.ScanSummary, baseJob domain.JobRegistration) {
	var err error

	if baseJob != nil {

		var configs []domain.JobConfig
		// here we pass the scanner source id so the spawned rescan job uses the same scanner (for the case of an organization using multiple scanners)
		if configs, err = job.db.GetJobConfigByOrgIDAndJobIDWithSC(job.config.OrganizationID(), baseJob.ID(), job.insource.ID()); err == nil {

			if configs != nil && len(configs) > 0 && configs[0] != nil {

				config := configs[0]

				_, _, err = job.db.CreateJobHistoryWithParentID(
					baseJob.ID(),
					config.ID(),
					domain.JobStatusPending,
					baseJob.Priority(),
					"",
					0,
					correspondingScanSummary.ScanClosePayload(),
					"",
					time.Now().UTC(),
					"SCAN SYNC JOB",
					correspondingScanSummary.ParentJobID())

				if err == nil {
					job.lstream.Send(log.Info("ScanCloseJob successfully queued"))
				} else {
					job.lstream.Send(log.Critical("ScanCloseJob could not be queued", err))
				}

			} else {
				job.lstream.Send(log.Errorf(err, "Could not find a config for the scan close job"))
			}
		} else {
			job.lstream.Send(log.Errorf(err, "Error while loading config from database for rescan close this, Error:  [%s]", err))
		}

	} else {
		job.lstream.Send(log.Errorf(err, "Base ScanClose Job returned null"))
	}
}

func (job *ScanSyncJob) pushScansOntoChannel(in []domain.ScanSummary) <-chan []byte {
	var out = make(chan []byte)

	go func() {
		defer handleRoutinePanic(job.lstream)
		defer close(out)

		for _, scan := range in {
			if payload, err := scanClosePayloadToScanPayload(scan.ScanClosePayload()); err == nil {
				select {
				case <-job.ctx.Done():
					return
				case out <- payload:
				}
			} else {
				job.lstream.Send(log.Errorf(err, "error while unmarshalling scan close payload"))
			}
		}

		for _, scan := range job.payload.ScheduledScanPayloads {
			select {
			case <-job.ctx.Done():
				return
			case out <- []byte(scan):
			}
		}
	}()

	return out
}

func (job *ScanSyncJob) getBaseJob() (baseJob domain.JobRegistration, err error) {
	if baseJob, err = job.db.GetJobsByStruct("ScanCloseJob"); err == nil {
		if baseJob == nil {
			job.lstream.Send(log.Errorf(err, "Empty list of base jobs returned"))
		}
	} else {
		job.lstream.Send(log.Errorf(err, "Error while loading ScanClose this struct from db"))
	}

	return baseJob, err
}

func (job *ScanSyncJob) getCorrespondingScanSummary(in domain.Scan, options []domain.ScanSummary) (out domain.ScanSummary, err error) {
	if len(in.ID()) > 0 {
		for _, option := range options {
			if sord(option.SourceKey()) == in.ID() {
				out = option
				break
			}
		}

		if out == nil {
			scanClosePayload := &ScanClosePayload{}
			scanClosePayload.Scan = in
			scanClosePayload.ScanID = in.ID()
			scanClosePayload.Type = domain.RescanScheduled

			var bytePayload []byte

			if bytePayload, err = json.Marshal(scanClosePayload); err == nil {
				_, _, err = job.db.CreateScanSummary(
					job.insource.SourceID(),
					job.insource.ID(),
					job.config.OrganizationID(),
					in.ID(),
					domain.ScanQUEUED,
					string(bytePayload),
					job.id,
				)

				if err == nil {
					job.lstream.Send(log.Infof("created ScanSummary db entry for [%s/%s]", in.ID(), in.Title()))

					if out, err = job.db.GetScanSummary(job.insource.SourceID(), job.config.OrganizationID(), in.ID()); err == nil {
						if out == nil {
							job.lstream.Send(log.Infof("could not find scheduled scan summary [%s]", in.ID()))
						}
					} else {
						job.lstream.Send(log.Infof("error while grabbing scheduled scan summary [%s]", in.ID()))
					}
				} else {
					job.lstream.Send(log.Infof("failed to create ScanSummary db entry for [%s/%s]", in.ID(), in.Title()))
				}
			}
		}
	} else {
		err = fmt.Errorf("empty scan ID")
	}

	return out, err
}
