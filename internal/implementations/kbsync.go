package implementations

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/nortonlifelock/aegis/internal/integrations"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/log"
)

// VulnSyncJob implements the IJob interface to sync vulnerabilities from a scanning engine
type VulnSyncJob struct {
	id          string
	payloadJSON string
	ctx         context.Context
	db          domain.DatabaseConnection
	lstream     log.Logger
	appconfig   domain.Config
	config      domain.JobConfig
	insource    domain.SourceConfig
	outsource   domain.SourceConfig
}

// Process downloads vulnerability information from a scanning engine, and then creates an entry in the VulnerabilityInfo table if one does not exist,
// and updates the entry in the VulnerabilityInfo table if one does not exist
func (job *VulnSyncJob) Process(ctx context.Context, id string, appconfig domain.Config, db domain.DatabaseConnection, lstream log.Logger, payload string, jobConfig domain.JobConfig, inSource []domain.SourceConfig, outSource []domain.SourceConfig) (err error) {

	var ok bool
	if job.ctx, job.id, job.appconfig, job.db, job.lstream, job.payloadJSON, job.config, job.insource, job.outsource, ok = validInputs(ctx, id, appconfig, db, lstream, payload, jobConfig, inSource, outSource); ok {

		job.lstream.Send(log.Infof("Establishing %s connection...", job.insource.Source()))

		var scanner integrations.Vscanner
		scanner, err = integrations.NewVulnScanner(job.ctx, job.insource.Source(), job.db, job.lstream, job.appconfig, job.insource)
		if err == nil {
			job.lstream.Send(log.Infof("%s connection established - loading vulnerabilities", job.insource.Source()))

			// read operations from the database are expensive (100ms+ each vuln for tables of len>100k) so we preload vulnerability information from the database.
			// The information gathered before the run starts does not need to be updated mid-run as each vulnerability only is updated a single time
			var preloadVulnerabilities []domain.VulnerabilityInfo
			preloadVulnerabilities, err = job.db.GetVulnInfoBySourceID(job.insource.SourceID())
			if err == nil {
				vulnerabilities := scanner.KnowledgeBase(job.ctx, job.config.LastJobStart())
				job.processVulnerabilities(job.ctx, vulnerabilities, preloadVulnerabilities)
			} else {
				job.lstream.Send(log.Error("error while preloading vulnerabilities", err))
			}

		} else {
			job.lstream.Send(log.Errorf(err, "Error while establishing %s session", job.insource.Source()))
		}

	} else {
		err = fmt.Errorf("input validation failed")
	}

	return err
}

func (job *VulnSyncJob) processVulnerabilities(ctx context.Context, vulnerabilities <-chan domain.Vulnerability, preloaded []domain.VulnerabilityInfo) {
	index := 0
	var wg sync.WaitGroup

	for {
		select {
		case <-ctx.Done():
			return
		case vulnerability, ok := <-vulnerabilities:
			if ok {
				index++
				if index%100 == 0 {
					wg.Wait()
				}

				wg.Add(1)
				go func(currentVuln domain.Vulnerability) {
					defer wg.Done()
					defer handleRoutinePanic(job.lstream)
					var err error

					if len(currentVuln.SourceID()) > 0 {

						var solution string
						if solutionChan, solutionErr := currentVuln.Solutions(ctx); solutionErr == nil {
							select {
							case <-ctx.Done():
								return
							case solutionInterface, ok := <-solutionChan:
								if ok {
									solution = solutionInterface.String()

									// kick off a thread to clear out the channel
									go func() {
										defer handleRoutinePanic(job.lstream)
										job.clearSolutionChannel(ctx, solutionChan)
									}()
								} else {
									break
								}
							}
						} else {
							job.lstream.Send(log.Errorf(err, "error while loading solution for vulnerability %v", currentVuln.SourceID()))
						}

						var vulnInDB = getPreloadedDbVuln(preloaded, currentVuln.SourceID())
						if vulnInDB == nil {
							job.processNewVulnerability(ctx, currentVuln, solution, vulnInDB)
						} else {
							job.processOldVulnerability(ctx, vulnInDB, currentVuln, solution)
						}
					} else {
						job.lstream.Send(log.Errorf(err, "vulnerability did not have an ID set"))
					}
				}(vulnerability)
			} else {
				wg.Wait()
				return
			}
		}
	}
}

func (job *VulnSyncJob) clearSolutionChannel(ctx context.Context, solutionChan <-chan domain.Solution) {
	for {
		select {
		case <-ctx.Done():
			return
		case _, ok := <-solutionChan:
			if !ok {
				return
			}
		}
	}
}

func (job *VulnSyncJob) processOldVulnerability(ctx context.Context, vulnInDB domain.VulnerabilityInfo, currentVuln domain.Vulnerability, solution string) {
	var err error
	if (vulnInDB.Updated() == nil && currentVuln.Updated().After(tord(vulnInDB.Created()))) || (vulnInDB.Updated() != nil && currentVuln.Updated().After(tord(vulnInDB.Updated()))) {

		// The vulnerability is already in the database, we can create it's references right away
		//if err = job.createVulnerabilityReferences(ctx, currentVuln, vulnInDB); err == nil {
		//	job.lstream.Send(log.Infof("Created vulnerability reference for VulnID: %v", currentVuln.SourceID()))
		//} else {
		//	job.lstream.Send(log.Errorf(err, "Error while creating vulnerability reference for VulnID: [%v]", currentVuln.SourceID()))
		//}

		// Update the information in the database
		if err = job.updateVulnerability(vulnInDB, currentVuln, solution); err == nil {
			job.lstream.Send(log.Infof("Updated vulnerability %v", currentVuln.ID()))
		} else {
			job.lstream.Send(log.Errorf(err, "Error while updating vulnerability [%v]", currentVuln.ID()))
		}
	} else {
		job.lstream.Send(log.Debugf("Did not need to update %v as it has not been modified recently", currentVuln.SourceID()))
	}
}

func (job *VulnSyncJob) processNewVulnerability(ctx context.Context, currentVuln domain.Vulnerability, solution string, vulnInDB domain.VulnerabilityInfo) {
	err := job.createVulnerability(currentVuln, solution)
	if err == nil {
		job.lstream.Send(log.Infof("Created vulnerability %v", currentVuln.SourceID()))

		// After we've created the vulnerability we can add it's vulnerability references to the database
		if err = job.createVulnerabilityReferences(ctx, currentVuln, vulnInDB); err == nil {
			job.lstream.Send(log.Infof("Created vulnerability reference for VulnID: %v", currentVuln.SourceID()))
		} else {
			job.lstream.Send(log.Errorf(err, "Error while creating vulnerability reference for VulnID: [%v]", currentVuln.SourceID()))
		}

	} else {
		job.lstream.Send(log.Errorf(err, "Error while creating vulnerability [%v]", currentVuln.SourceID()))
	}
}

func (job *VulnSyncJob) updateVulnerability(vulnInDB domain.VulnerabilityInfo, currentVuln domain.Vulnerability, solution string) (err error) {
	var cvss3Pointer = currentVuln.CVSS3()
	if cvss3Pointer != nil {
		_, _, err = job.db.UpdateVulnByID(
			vulnInDB.ID(),
			currentVuln.SourceID(),
			currentVuln.Name(),
			job.insource.SourceID(),
			currentVuln.CVSS2(),
			*cvss3Pointer,
			currentVuln.Description(),
			sord(currentVuln.Threat()),
			solution,
			currentVuln.Software(),
			currentVuln.DetectionInformation(),
		)
	} else {
		_, _, err = job.db.UpdateVulnByIDNoCVSS3(
			vulnInDB.ID(),
			currentVuln.SourceID(),
			currentVuln.Name(),
			job.insource.SourceID(),
			currentVuln.CVSS2(),
			currentVuln.Description(),
			sord(currentVuln.Threat()),
			solution,
			currentVuln.Software(),
			currentVuln.DetectionInformation(),
		)
	}

	return err
}

func (job *VulnSyncJob) createVulnerability(currentVuln domain.Vulnerability, solution string) (err error) {
	cvss3Pointer := currentVuln.CVSS3()

	if cvss3Pointer != nil {
		// Create the vulnerability if it hasn't been created yet
		_, _, err = job.db.CreateVulnInfo(
			currentVuln.SourceID(),
			currentVuln.Name(),
			job.insource.SourceID(),
			currentVuln.CVSS2(),
			*cvss3Pointer,
			currentVuln.Description(),
			sord(currentVuln.Threat()),
			solution,
			currentVuln.Software(),
			currentVuln.DetectionInformation(),
		)
	} else {
		_, _, err = job.db.CreateVulnInfoNoCVSS3(
			currentVuln.SourceID(),
			currentVuln.Name(),
			job.insource.SourceID(),
			currentVuln.CVSS2(),
			currentVuln.Description(),
			sord(currentVuln.Threat()),
			solution,
			currentVuln.Software(),
			currentVuln.DetectionInformation(),
		)
	}

	return err
}

// vulnSync method is called asynchronously, but only performs read operations
// therefore, mutexes are not necessary
func getPreloadedDbVuln(preloadVulnerabilities []domain.VulnerabilityInfo, sourceVulnID string) (vulnInDb domain.VulnerabilityInfo) {
	for index := range preloadVulnerabilities {
		if preloadVulnerabilities[index].SourceVulnID() == sourceVulnID {
			vulnInDb = preloadVulnerabilities[index]
			break
		}
	}

	return vulnInDb
}

func createVulnRef(db domain.DatabaseConnection, vulnInfoID string, vulnRef string, sourceID string, lstream log.Logger) (err error) {
	var existingRef domain.VulnerabilityReference
	if existingRef, err = db.GetVulnRef(vulnInfoID, sourceID, vulnRef); err == nil {
		if existingRef == nil {
			lstream.Send(log.Infof("Creating vulnerability reference %s", vulnRef))
			_, _, err = db.CreateVulnRef(vulnInfoID, sourceID, vulnRef, determineRefType(vulnRef))
		}
	}

	return err
}

func determineRefType(reference string) (refType int) {
	if strings.Contains(reference, "CVE") {
		refType = domain.CVE
	} else if strings.Index(reference, "MS") == 0 {
		refType = domain.MS
	} else {
		refType = domain.Vendor
	}

	return refType
}

// Vulnerability references include CVEs and Microsoft Ids (MSid)
func (job *VulnSyncJob) createVulnerabilityReferences(ctx context.Context, currentVuln domain.Vulnerability, vulnInDB domain.VulnerabilityInfo) (err error) {
	defer handleRoutinePanic(job.lstream)
	var vulnInfo domain.VulnerabilityInfo

	if vulnInDB == nil {
		vulnInfo, err = job.db.GetVulnInfoBySourceVulnIDSourceID(currentVuln.SourceID(), job.insource.SourceID(), currentVuln.Updated())
	} else {
		vulnInfo = vulnInDB
	}

	if err == nil {
		if vulnInfo != nil {
			err = job.createMSIDAndVendorRefFromVuln(ctx, currentVuln)
		} else {
			err = errors.New("no records were returned from GetVulnInfoBySourceVulnIDSourceId")
		}

	}
	return err
}

// createMSIDAndVendorRefFromVuln takes a list of vendor references and makes an entry in the database if one does not exist yet
func (job *VulnSyncJob) createMSIDAndVendorRefFromVuln(ctx context.Context, vuln domain.Vulnerability) error {
	refChan, err := vuln.References(ctx)
	if err == nil {
		for {
			select {
			case <-ctx.Done():
				return nil
			case ref, ok := <-refChan:
				if ok {
					err = createVulnRef(job.db, vuln.SourceID(), ref.Reference(), job.insource.SourceID(), job.lstream)
					if err != nil {
						return err
					}
				} else {
					return nil
				}
			}
		}
	}

	return err
}

// TODO add this in logic
func (job *VulnSyncJob) createCategoryIfNonExistent(vulnCategory string) (err error) {
	var category []domain.Category
	category, err = job.db.GetCategoryByName(vulnCategory)
	if err == nil {
		if category == nil || len(category) == 0 {
			_, _, err = job.db.CreateCategory(vulnCategory)
			if err == nil {
				job.lstream.Send(log.Info("Created new category " + vulnCategory))
			} else {
				err = fmt.Errorf("error while creating new category - %v", err)
			}
		}
	} else {
		err = fmt.Errorf("error while loading category from database - %v", err)
	}
	return err
}
