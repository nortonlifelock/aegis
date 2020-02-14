package implementations

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nortonlifelock/aegis/backend/domain"
	"github.com/nortonlifelock/aegis/backend/integrations"
	"github.com/nortonlifelock/log"
	"strconv"
	"strings"
	"sync"
	"time"
)

// AssetSyncJob implements the Job interface required to run the job
type AssetSyncJob struct {
	Payload *AssetSyncPayload

	// the detection status must be queried for each detection, so we cache them
	detectionStatuses []domain.DetectionStatus

	// the vuln cache maps the vulnerability ID to the vulnerability information in the database
	vulnCache sync.Map

	id          string
	payloadJSON string
	ctx         context.Context
	db          domain.DatabaseConnection
	lstream     log.Logger
	appconfig   domain.Config
	config      domain.JobConfig
	insources   domain.SourceConfig
	outsource   domain.SourceConfig
}

// AssetSyncPayload holds the asset groups to be synced by the job. loaded from the job history Payload
type AssetSyncPayload struct {
	GroupIDs []int `json:"groups"`
}

// buildPayload loads the Payload from the job history into the Payload object
func (job *AssetSyncJob) buildPayload(pjson string) (err error) {
	job.Payload = &AssetSyncPayload{}

	if len(pjson) > 0 {
		err = json.Unmarshal([]byte(pjson), job.Payload)
		if err == nil {
			if len(job.Payload.GroupIDs) == 0 {
				err = fmt.Errorf("did not provide group in Payload")
			}
		}
	} else {
		err = fmt.Errorf("no Payload provided to job")
	}

	return err
}

// Process downloads asset information from a scanner (such as IP/vulnerability detections) and stores it in the database
func (job *AssetSyncJob) Process(ctx context.Context, id string, appconfig domain.Config, db domain.DatabaseConnection, lstream log.Logger, payload string, jobConfig domain.JobConfig, inSource []domain.SourceConfig, outSource []domain.SourceConfig) (err error) {

	var ok bool
	if job.ctx, job.id, job.appconfig, job.db, job.lstream, job.payloadJSON, job.config, job.insources, job.outsource, ok = validInputs(ctx, id, appconfig, db, lstream, payload, jobConfig, inSource, outSource); ok {

		if err = job.buildPayload(job.payloadJSON); err == nil {

			job.lstream.Send(log.Debug("Creating scanner connection..."))

			var vscanner integrations.Vscanner
			if vscanner, err = integrations.NewVulnScanner(job.ctx, job.insources.Source(), job.db, job.lstream, job.appconfig, job.insources); err == nil {

				if job.detectionStatuses, err = job.db.GetDetectionStatuses(); err == nil {
					job.lstream.Send(log.Debug("Scanner connection created, beginning processing..."))

					for _, groupID := range job.Payload.GroupIDs {
						if err = job.createAssetGroupInDB(groupID, job.insources.SourceID(), job.insources.ID()); err == nil {
							select {
							case <-job.ctx.Done():
								return
							default:
							}

							job.lstream.Send(log.Infof("started processing %v", groupID))
							job.processGroup(vscanner, groupID)
							job.lstream.Send(log.Infof("finished processing %v", groupID))
						} else {
							job.lstream.Send(log.Error("error while creating asset group", err))
						}
					}
				} else {
					job.lstream.Send(log.Error("error while preloading detection statuses", err))
				}
			} else {
				job.lstream.Send(log.Error("error while creating scanner connection", err))
			}
		} else {
			err = fmt.Errorf("error while building payload - %s", err.Error())
		}
	} else {
		err = fmt.Errorf("input validation failed")
	}

	return err
}

func (job *AssetSyncJob) fanInDetections(in <-chan domain.Detection) (devIDToDetection map[string][]domain.Detection) {
	devIDToDetection = make(map[string][]domain.Detection)

	for {
		select {
		case <-job.ctx.Done():
			return
		case detection, ok := <-in:
			if ok {

				if device, err := detection.Device(); err == nil && device != nil {
					if devIDToDetection[sord(device.SourceID())] == nil {
						devIDToDetection[sord(device.SourceID())] = make([]domain.Detection, 0)
					}
					devIDToDetection[sord(device.SourceID())] = append(devIDToDetection[sord(device.SourceID())], detection)
				} else {
					job.lstream.Send(log.Errorf(err, "error while loading device"))
				}

			} else {
				return devIDToDetection
			}
		}
	}
}

// This method is responsible for gathering the assets of the group, as well as kicking off the threads that process each asset
func (job *AssetSyncJob) processGroup(vscanner integrations.Vscanner, groupID int) {
	var groupIDString = strconv.Itoa(groupID)

	// gather the asset information
	detectionChan, err := vscanner.Detections(job.ctx, []string{groupIDString})
	if err == nil {

		devIDToDetections := job.fanInDetections(detectionChan)
		job.lstream.Send(log.Infof("Finished loading detections for %d devices", len(devIDToDetections)))

		const simultaneousCount = 10
		var permitThread = make(chan bool, simultaneousCount)
		for i := 0; i < simultaneousCount; i++ {
			permitThread <- true
		}

		var wg sync.WaitGroup
		for deviceID, detections := range devIDToDetections {
			select {
			case <-job.ctx.Done():
				return
			case <-permitThread:

				wg.Add(1)
				go func(deviceID string, detections []domain.Detection) {
					defer wg.Done()
					defer func() {
						permitThread <- true
					}()

					if len(detections) > 0 {

						job.lstream.Send(log.Infof("Working on %d detections for %s", len(detections), deviceID))

						if asset, err := detections[0].Device(); err == nil {
							err = job.addDeviceInformationToDB(asset, groupID)
							if err == nil {
								job.processAsset(deviceID, asset, detections, groupID)
							} else {
								job.lstream.Send(log.Errorf(err, "error while adding asset information to the database"))
							}
						} else {
							job.lstream.Send(log.Errorf(err, "error while loading device information for %s", deviceID))
						}
					}
				}(deviceID, detections)
			}
		}
		wg.Wait()
	} else {
		job.lstream.Send(log.Error("error while grabbing device and vulnerability information", err))
	}
}

// Only process the asset if it has not been processed by another group
func (job *AssetSyncJob) processAsset(deviceID string, asset domain.Device, detections []domain.Detection, groupID int) {
	var err error

	if len(sord(asset.SourceID())) > 0 {
		var existingDeviceInDb domain.Device
		if existingDeviceInDb, err = job.db.GetDeviceByAssetOrgID(sord(asset.SourceID()), job.config.OrganizationID()); err == nil && existingDeviceInDb != nil {

			if detections != nil {

				var wg sync.WaitGroup
				for _, detection := range detections {

					wg.Add(1)
					go func(detection domain.Detection) {
						defer wg.Done()

						if detection != nil {
							_ = job.processAssetDetections(existingDeviceInDb, sord(asset.SourceID()), detection)
						} else {
							job.lstream.Send(log.Errorf(err, "nil detection found for", sord(asset.SourceID())))
						}
					}(detection)
				}
				wg.Wait()
			} else {
				job.lstream.Send(log.Errorf(err, "error while processing asset information in database"))
			}
		} else {
			job.lstream.Send(log.Errorf(fmt.Errorf("could not find device in database for %s", sord(asset.SourceID())), "db error"))
		}
	} else {
		job.lstream.Send(log.Errorf(nil, "empty asset ID gathered from scanner"))
	}
}

// This method creates/gathers the entry for the OS Type as well as updates/creates the asset information in the database
func (job *AssetSyncJob) addDeviceInformationToDB(asset domain.Device, groupID int) (err error) {
	var ostFromDb domain.OperatingSystemType
	if len(asset.OS()) > 0 {
		ostFromDb, err = job.grabAndCreateOsType(asset.OS())
	} else {
		ostFromDb, err = job.grabAndCreateOsType(unknown)
	}

	// this updates asset's OST to the same OST but w/ populated db id
	if err == nil {
		err = job.enterAssetInformationInDB(asset, ostFromDb.ID(), groupID)
		if err != nil {
			job.lstream.Send(log.Error("error while processing asset", err))
		}
	} else {
		job.lstream.Send(log.Error("Couldn't gather database OS information", err))
	}

	return err
}

// this method checks the database to see if an asset under that ip/org and creates an entry if one doesn't exist.
// if an entry exists but does not have an asset id set (which occurs when the CloudSync Job) finds the asset first,
// this method then enters the asset id for that entry
func (job *AssetSyncJob) enterAssetInformationInDB(asset domain.Device, osTypeID int, groupID int) (err error) {
	if asset != nil {

		if len(sord(asset.SourceID())) > 0 {

			var ip = asset.IP()

			var deviceInDB domain.Device
			// first try to find the device in the database using the source asset id
			if deviceInDB, err = job.db.GetDeviceByAssetOrgID(sord(asset.SourceID()), job.config.OrganizationID()); err == nil { // TODO include org id parameter
				if deviceInDB == nil {

					// second we try to find the device in the database using the IP
					if len(asset.IP()) > 0 {
						deviceInDB, err = job.db.GetDeviceByScannerSourceID(ip, groupID, job.config.OrganizationID())
					}
				}

				if err == nil {
					if deviceInDB == nil {

						// TODO currently this procedure just sets IsVirtual to false - how do I find that value?
						_, _, err = job.db.CreateDevice(
							sord(asset.SourceID()),
							job.insources.SourceID(),
							ip,
							asset.HostName(),
							asset.MAC(),
							groupID,
							job.config.OrganizationID(),
							asset.OS(),
							osTypeID,
						)
						if err == nil {
							job.lstream.Send(log.Infof("[+] Device [%v] created", sord(asset.SourceID())))
						} else {
							err = fmt.Errorf(fmt.Sprintf("[-] Error while creating device [%s] - %s", sord(asset.SourceID()), err.Error()))
						}

					} else {

						// this block of code is for when cloud sync job finds the asset before the ASJ does, as the CSJ doesn't set the asset id
						// we also update the os type id because the ASJ will have a more accurate os return
						if len(sord(deviceInDB.SourceID())) == 0 && len(sord(asset.SourceID())) > 0 {
							_, _, err = job.db.UpdateAssetIDOsTypeIDOfDevice(deviceInDB.ID(), sord(asset.SourceID()), job.insources.SourceID(), groupID, asset.OS(), asset.HostName(), osTypeID, job.config.OrganizationID())
							if err == nil {
								job.lstream.Send(log.Infof("Updated device info for asset [%v]", sord(asset.SourceID())))
							} else {
								err = fmt.Errorf(fmt.Sprintf("could not update the asset id for device with ip [%s] - %s", ip, err.Error()))
							}
						} else {
							job.lstream.Send(log.Debugf("DB entry for device [%v] exists, skipping...", sord(asset.SourceID())))
						}
					}
				} else {
					job.lstream.Send(log.Errorf(err, "error while loading device from database"))
				}

			} else {
				job.lstream.Send(log.Errorf(err, "error while loading device from database"))
			}

		} else {
			err = fmt.Errorf("device with id [%s] did not have asset id returned from vuln scanner", sord(asset.SourceID()))
		}

	} else {
		err = fmt.Errorf("improper enterAssetInformationInDB input - nil device passed to process asset")
	}

	return err
}

// This method creates a detection entry in the database for the device/vulnerability combo
// If the detection entry already exists, it increments the amount of times it has been seen by this job by one
// This method is also responsible for gathering detections for the vulnerability
func (job *AssetSyncJob) processAssetDetections(deviceInDb domain.Device, assetID string, vuln domain.Detection) (err error) {
	// the result ID may be concatenated to the end of the vulnerability ID. we chop it off the result from the vulnerability ID with the following line
	vulnID := strings.Split(vuln.VulnerabilityID(), ";")[0]

	var vulnInfo domain.VulnerabilityInfo
	if vulnInfoInterface, ok := job.vulnCache.Load(vulnID); ok {
		if vulnInfo, ok = vulnInfoInterface.(domain.VulnerabilityInfo); !ok {
			err = fmt.Errorf("cache error while loading vulnerability info")
		}
	} else {
		vulnInfo, err = job.db.GetVulnInfoBySourceVulnID(vulnID)
		if err == nil && vulnInfo != nil {
			job.vulnCache.Store(vulnID, vulnInfo)
		}
	}

	if err == nil {
		if vulnInfo != nil {
			job.createOrUpdateDetection(deviceInDb, vulnInfo, vuln, assetID)
		} else {
			job.lstream.Send(log.Error("could not find vulnerability in database", fmt.Errorf("[%s] does not have an entry in the database", vulnID)))
		}

	} else {
		job.lstream.Send(log.Errorf(err, "Error while gathering vulnerability info for [%s]", vulnID))
	}

	return err
}

func (job *AssetSyncJob) getExceptionID(assetID string, vulnInfo domain.VulnerabilityInfo) (exceptionID string) {
	if exception, err := job.db.GetExceptionByVulnIDOrg(assetID, vulnInfo.SourceVulnID(), job.config.OrganizationID()); err == nil {
		if exception != nil {
			exceptionID = exception.ID()
		}
	} else {
		job.lstream.Send(log.Errorf(err, "Error while gathering exceptions for device [%v]", assetID))
	}

	return exceptionID
}

// This method creates a detection entry if one does not exist, and updates the entry if one does
func (job *AssetSyncJob) createOrUpdateDetection(deviceInDb domain.Device, vulnInfo domain.VulnerabilityInfo, detectionFromScanner domain.Detection, assetID string) {
	var err error

	var detectionInDB domain.Detection
	detectionInDB, err = job.db.GetDetection(sord(deviceInDb.SourceID()), vulnInfo.ID())
	if err == nil {
		var detectionStatus domain.DetectionStatus
		if detectionStatus = job.getDetectionStatus(detectionFromScanner.Status()); detectionStatus != nil {

			if detectionInDB == nil {
				job.createDetection(detectionFromScanner, job.getExceptionID(assetID, vulnInfo), deviceInDb, vulnInfo, assetID, detectionStatus.ID())
			} else {

				var canSkipUpdate bool
				if !detectionInDB.Updated().IsZero() && !detectionFromScanner.Updated().IsZero() {
					if detectionInDB.Updated().After(detectionFromScanner.Updated()) {
						canSkipUpdate = true
					}
				}

				if !canSkipUpdate {
					_, _, err = job.db.UpdateDetectionTimesSeen(
						sord(deviceInDb.SourceID()),
						vulnInfo.ID(),
						job.getExceptionID(assetID, vulnInfo),
						detectionFromScanner.TimesSeen(),
						detectionStatus.ID(),
					)

					if err == nil {
						job.lstream.Send(log.Infof("Updated detection for device/vuln [%v|%v]", assetID, vulnInfo.ID()))
					} else {
						job.lstream.Send(log.Errorf(err, "Error while updating detection for device/vuln [%v|%v]", assetID, vulnInfo.ID()))
					}
				} else {
					job.lstream.Send(log.Infof("Skipping detection update for device/vuln [%v|%v] [%v after %v]", assetID, vulnInfo.ID(), detectionInDB.Updated(), detectionFromScanner.Updated()))
				}
			}
		} else {
			job.lstream.Send(log.Errorf(err, "could not find detection status with name [%s]", detectionFromScanner.Status()))
		}
	} else {
		job.lstream.Send(log.Debugf("Detection already exists for device/vuln [%v|%v]", assetID, vulnInfo.ID()))
	}
}

// This method creates the detection entry in the database
func (job *AssetSyncJob) createDetection(vuln domain.Detection, exceptionID string, deviceInDb domain.Device, vulnInfo domain.VulnerabilityInfo, assetID string, detectionStatusID int) {
	var err error

	var detected *time.Time
	if detected, err = vuln.Detected(); err == nil {
		if detected != nil {
			if len(exceptionID) == 0 {

				if vuln.ActiveKernel() == nil {
					_, _, err = job.db.CreateDetection(
						job.config.OrganizationID(),
						job.insources.SourceID(),
						sord(deviceInDb.SourceID()),
						vulnInfo.ID(),
						*detected,
						vuln.Proof(),
						vuln.Port(),
						vuln.Protocol(),
						detectionStatusID,
						vuln.TimesSeen(),
					)
				} else {
					_, _, err = job.db.CreateDetectionActiveKernel(
						job.config.OrganizationID(),
						job.insources.SourceID(),
						sord(deviceInDb.SourceID()),
						vulnInfo.ID(),
						*detected,
						vuln.Proof(),
						vuln.Port(),
						vuln.Protocol(),
						iord(vuln.ActiveKernel()),
						detectionStatusID,
						vuln.TimesSeen(),
					)
				}

			} else {

				if vuln.ActiveKernel() == nil {
					_, _, err = job.db.CreateDetectionWithIgnore(
						job.config.OrganizationID(),
						job.insources.SourceID(),
						sord(deviceInDb.SourceID()),
						vulnInfo.ID(),
						exceptionID,
						*detected,
						vuln.Proof(),
						vuln.Port(),
						vuln.Protocol(),
						detectionStatusID,
						vuln.TimesSeen(),
					)
				} else {
					_, _, err = job.db.CreateDetectionWithIgnoreActiveKernel(
						job.config.OrganizationID(),
						job.insources.SourceID(),
						sord(deviceInDb.SourceID()),
						vulnInfo.ID(),
						exceptionID,
						*detected,
						vuln.Proof(),
						vuln.Port(),
						vuln.Protocol(),
						iord(vuln.ActiveKernel()),
						detectionStatusID,
						vuln.TimesSeen(),
					)
				}
			}
		} else {
			err = fmt.Errorf("could not find the time of the detection")
		}

	} else {
		err = fmt.Errorf("error while gathering date of detection - %v", err.Error())
	}

	if err == nil {
		job.lstream.Send(log.Infof("Created detection for device/vuln [%v|%v]", assetID, vulnInfo.ID()))
	} else {
		job.lstream.Send(log.Errorf(err, "Error while creating detection for device/vuln [%v|%v]", assetID, vulnInfo.ID()))
	}
}

func (job *AssetSyncJob) getDetectionStatus(status string) (detectionStatus domain.DetectionStatus) {
	for _, potentialMatch := range job.detectionStatuses {
		if strings.ToLower(status) == strings.ToLower(potentialMatch.Status()) {
			detectionStatus = potentialMatch
			break
		}
	}

	return detectionStatus
}

// This method creates an entry in the database for the operating system type. It then returns the entry so that the id of the OST
// may be used for foreign key references
func (job *AssetSyncJob) grabAndCreateOsType(operatingSystem string) (output domain.OperatingSystemType, err error) {
	if len(operatingSystem) > 0 {
		output, err = job.db.GetOperatingSystemType(operatingSystem)
		if err == nil {
			if output == nil {
				err = fmt.Errorf("could not discern operating system type of [%s]", operatingSystem)
			}
		} else {
			err = fmt.Errorf("(GetOST) %s - [%s]", err.Error(), operatingSystem)
		}
	} else {
		err = fmt.Errorf("operating system sent nil to grabAndCreateOsType")
	}

	return output, err
}

func (job *AssetSyncJob) createAssetGroupInDB(groupID int, scannerSourceID string, scannerSourceConfigID string) (err error) {
	var assetGroup domain.AssetGroup
	if assetGroup, err = job.db.GetAssetGroup(job.config.OrganizationID(), groupID, scannerSourceConfigID); err == nil {
		if assetGroup == nil {
			if _, _, err = job.db.CreateAssetGroup(job.config.OrganizationID(), groupID, scannerSourceID, scannerSourceConfigID); err == nil {

			} else {
				err = fmt.Errorf("error while creating asset group - %v", err.Error())
			}
		}
	} else {
		err = fmt.Errorf("error while grabbing asset group - %v", err.Error())
	}

	return err
}
