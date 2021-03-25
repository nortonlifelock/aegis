package implementations

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nortonlifelock/aegis/internal/integrations"
	"github.com/nortonlifelock/aegis/pkg/domain"
	"github.com/nortonlifelock/aegis/pkg/log"
	"regexp"
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

	// a cache for exceptions that apply to an OS/vulnID combo (as opposed to device/vulnid combo)
	globalExceptions []compiledException

	// a cache for device/vuln specific exceptions
	deviceIDToVulnIDToException map[string]map[string]domain.Ignore

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
	GroupIDs                           []string `json:"groups"`
	GroupIDsExemptFromGlobalExceptions []string `json:"groups_exempt_from_global_exceptions"`
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

					if job.globalExceptions, job.deviceIDToVulnIDToException, err = job.preloadIgnores(); err == nil {

						var cloudTags = make([]string, 0)
						var groupIDs = make([]string, 0)
						const tagPrefix = "tag-"
						for _, id := range job.Payload.GroupIDs {
							if strings.Index(id, tagPrefix) >= 0 {
								cloudTags = append(cloudTags, id)
							} else {
								groupIDs = append(groupIDs, id)
							}
						}

						// it is important to pass the groupIDs one at a time so we know which group a returned asset belongs to
						for _, groupID := range groupIDs {
							if err = job.createAssetGroupInDB(groupID, job.insources.SourceID(), job.insources.ID()); err == nil {
								select {
								case <-job.ctx.Done():
									return
								default:
								}

								job.lstream.Send(log.Infof("started processing %v", groupID))
								job.processGroup(vscanner, []string{groupID})
								job.lstream.Send(log.Infof("finished processing %v", groupID))
							} else {
								job.lstream.Send(log.Error("error while creating asset group", err))
							}
						}

						for _, cloudTag := range cloudTags {
							if err = job.createAssetGroupInDB(cloudTag, job.insources.SourceID(), job.insources.ID()); err == nil {
								select {
								case <-job.ctx.Done():
									return
								default:
								}

								job.lstream.Send(log.Infof("started processing %v", cloudTag))
								job.processGroup(vscanner, strings.Split(cloudTag, ","))
								job.lstream.Send(log.Infof("finished processing %v", cloudTag))
							} else {
								job.lstream.Send(log.Error("error while creating asset group", err))
							}
						}
					} else {
						job.lstream.Send(log.Errorf(err, "error while loading global exceptions"))
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
func (job *AssetSyncJob) processGroup(vscanner integrations.Vscanner, groupIDs []string) {
	// gather the asset information
	detectionChan, err := vscanner.Detections(job.ctx, groupIDs)
	if err == nil {

		devIDToDetections := job.fanInDetections(detectionChan)
		job.lstream.Send(log.Infof("Finished loading detections for %d devices", len(devIDToDetections)))

		var permitThread = getPermitThread(10)

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

							decomIgnoreID, err := getDecommIgnoreEntryForAsset(job.db, job.lstream, job.config.OrganizationID(), deviceID, job.insources.ID(), detections)
							if err != nil {
								job.lstream.Send(log.Errorf(err, "error while loading decomm ignore entry"))
							}

							var deviceInDB domain.Device
							deviceInDB, err = addDeviceInformationToDB(job.db, job.lstream, job.insources.ID(), job.config.OrganizationID(), asset, strings.Join(groupIDs, ","))
							if err == nil {
								job.processAsset(asset, deviceInDB, detections, decomIgnoreID, strings.Join(groupIDs, ","))
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

func getDecommIgnoreEntryForAsset(db domain.DatabaseConnection, lstream log.Logger, orgID string, deviceID string, scannerSourceID string, detections []domain.Detection) (decomIgnoreID string, err error) {
	if len(deviceID) > 0 {
		var decommIgnoreEntry domain.Ignore
		if decommIgnoreEntry, err = db.HasDecommissioned(deviceID, scannerSourceID, orgID); err == nil {
			if decommIgnoreEntry != nil {
				var allDetectionsFoundBeforeDecommDate = true
				for _, detection := range detections {
					if detectedDate, err := detection.Detected(); err == nil && !detectedDate.IsZero() {
						dueDate := decommIgnoreEntry.DueDate()
						if !dueDate.IsZero() && detectedDate.After(*dueDate) {
							allDetectionsFoundBeforeDecommDate = false

							// we found a vulnerability after the device was marked as decommissioned in the database
							lstream.Send(log.Warningf(nil, "Device [%s] has a vulnerability found after it's decommission date [%s after %s], deleting it's ignore entry in the database", deviceID, detectedDate.Format(time.RFC822Z), dueDate.Format(time.RFC822Z)))

							_, _, err = db.DeleteDecomIgnoreForDevice(scannerSourceID, deviceID, orgID)
							if err != nil {
								lstream.Send(log.Errorf(err, "Error while deleting ignore entry to [%s]", deviceID))
							}

							break
						}
					}
				}

				if allDetectionsFoundBeforeDecommDate {
					decomIgnoreID = decommIgnoreEntry.ID()
				}
			}
		}
	} else {
		err = fmt.Errorf("malformed device ID - [%s]", deviceID)
	}

	return decomIgnoreID, err
}

// Only process the asset if it has not been processed by another group
func (job *AssetSyncJob) processAsset(asset domain.Device, deviceInDB domain.Device, detections []domain.Detection, decomIgnoreID string, groupID string) {
	var err error

	if detections != nil {

		var wg sync.WaitGroup
		for _, detection := range detections {

			wg.Add(1)
			go func(detection domain.Detection) {
				defer wg.Done()

				if detection != nil {
					_ = job.processAssetDetections(deviceInDB, sord(asset.SourceID()), detection, decomIgnoreID, groupID)
				} else {
					job.lstream.Send(log.Errorf(err, "nil detection found for [%s]", sord(asset.SourceID())))
				}
			}(detection)
		}
		wg.Wait()
	} else {
		job.lstream.Send(log.Errorf(err, "error while processing asset information in database"))
	}
}

// This method creates/gathers the entry for the OS Type as well as updates/creates the asset information in the database
func addDeviceInformationToDB(db domain.DatabaseConnection, lstream log.Logger, sourceID string, orgID string, asset domain.Device, groupID string) (deviceInDB domain.Device, err error) {
	var ostFromDb domain.OperatingSystemType
	if len(asset.OS()) > 0 {
		ostFromDb, err = grabAndCreateOsType(db, asset.OS())
	} else {
		ostFromDb, err = grabAndCreateOsType(db, unknown)
	}

	if err != nil {
		lstream.Send(log.Error("Couldn't gather database OS information", err))
		return nil, err
	}

	err = enterAssetInformationInDB(db, lstream, sourceID, orgID, asset, ostFromDb.ID(), groupID)
	if err != nil {
		lstream.Send(log.Error("error while processing asset", err))
		return nil, err
	}

	if len(sord(asset.SourceID())) == 0 {
		err = fmt.Errorf("empty asset ID gathered from scanner")
		return nil, err
	}

	deviceInDB, err = db.GetDeviceByAssetOrgID(sord(asset.SourceID()), orgID)
	if err != nil {
		err = fmt.Errorf("error while getting device [%s] - %s", sord(asset.SourceID()), err.Error())
	} else if deviceInDB == nil {
		err = fmt.Errorf("could not find recently created device [%s]", sord(asset.SourceID()))
	}

	return deviceInDB, err
}

// this method checks the database to see if an asset under that ip/org and creates an entry if one doesn't exist.
// if an entry exists but does not have an asset id set (which occurs when the CloudSync Job) finds the asset first,
// this method then enters the asset id for that entry
func enterAssetInformationInDB(db domain.DatabaseConnection, lstream log.Logger, sourceID string, orgID string, asset domain.Device, osTypeID int, groupID string) (err error) {
	if asset != nil {

		if len(sord(asset.SourceID())) > 0 {

			var ip = asset.IP()

			var deviceInDB domain.Device
			// first try to find the device in the database using the source asset id
			if deviceInDB, err = db.GetDeviceByAssetOrgID(sord(asset.SourceID()), orgID); err == nil {
				if deviceInDB == nil {

					if len(sord(asset.InstanceID())) > 0 {
						var devicesByInstanceID []domain.Device
						devicesByInstanceID, err = db.GetDeviceByInstanceID(sord(asset.InstanceID()), orgID)
						if len(devicesByInstanceID) > 0 {
							deviceInDB = devicesByInstanceID[0]
						}
					}

					// second we try to find the device in the database using the IP
					if err == nil && deviceInDB == nil && len(asset.IP()) > 0 {
						deviceInDB, err = db.GetDeviceByScannerSourceID(ip, groupID, orgID)
					}
				}

				if err == nil {
					var deviceIsAlsoBeingTrackedUnderDifferentAssetID, deviceFoundByCloudSyncJobFirst, osNeedsUpdate, trackingMethodNeedsUpdate bool
					if deviceInDB != nil {
						// the device appearing under a different asset ID could be a purged asset that was rediscovered after still being online, for example
						deviceIsAlsoBeingTrackedUnderDifferentAssetID = len(sord(deviceInDB.SourceID())) > 0 && sord(deviceInDB.SourceID()) != sord(asset.SourceID())
						deviceFoundByCloudSyncJobFirst = len(sord(deviceInDB.SourceID())) == 0 && len(sord(asset.SourceID())) > 0
						osNeedsUpdate = sord(deviceInDB.SourceID()) == sord(asset.SourceID()) && len(asset.OS()) > 0 && deviceInDB.OS() != asset.OS()
						trackingMethodNeedsUpdate = sord(deviceInDB.SourceID()) == sord(asset.SourceID()) && len(sord(asset.TrackingMethod())) > 0 && sord(deviceInDB.TrackingMethod()) != sord(asset.TrackingMethod())
					}

					if deviceInDB == nil || deviceIsAlsoBeingTrackedUnderDifferentAssetID {
						_, _, err = db.CreateDevice(
							sord(asset.SourceID()),
							sourceID,
							ip,
							asset.HostName(),
							sord(asset.InstanceID()),
							asset.MAC(),
							groupID,
							orgID,
							asset.OS(),
							osTypeID,
							sord(asset.TrackingMethod()),
						)
						if err == nil {
							lstream.Send(log.Infof("[+] Device [%v] created", sord(asset.SourceID())))
						} else {
							err = fmt.Errorf(fmt.Sprintf("[-] Error while creating device [%s] - %s", sord(asset.SourceID()), err.Error()))
						}
					} else if deviceFoundByCloudSyncJobFirst || osNeedsUpdate || trackingMethodNeedsUpdate {
						// this block of code is for when cloud sync job finds the asset before the ASJ does, as the CSJ doesn't set the asset id
						// we also update the os type id because the ASJ will have a more accurate os return
						_, _, err = db.UpdateAssetIDOsTypeIDOfDevice(
							deviceInDB.ID(),
							sord(asset.SourceID()),
							sourceID,
							groupID,
							asset.OS(),
							asset.HostName(),
							osTypeID,
							sord(asset.TrackingMethod()),
							orgID,
						)
						if err == nil {
							lstream.Send(log.Infof("Updated device info for asset [%v]", sord(asset.SourceID())))
						} else {
							err = fmt.Errorf(fmt.Sprintf("could not update the asset id for device with ip [%s] - %s", ip, err.Error()))
						}
					} else {
						lstream.Send(log.Debugf("DB entry for device [%v] exists, skipping...", sord(asset.SourceID())))
					}
				} else {
					lstream.Send(log.Errorf(err, "error while loading device from database"))
				}

			} else {
				lstream.Send(log.Errorf(err, "error while loading device from database"))
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
func (job *AssetSyncJob) processAssetDetections(deviceInDb domain.Device, assetID string, detectionFromScanner domain.Detection, decomIgnoreID string, groupID string) (err error) {
	// the result ID may be concatenated to the end of the vulnerability ID. we chop it off the result from the vulnerability ID with the following line
	var vulnID string
	var resultID string
	var vulnResult = strings.Split(detectionFromScanner.VulnerabilityID(), domain.VulnPathConcatenator)
	vulnID = vulnResult[0]
	if len(vulnResult) > 1 {
		resultID = vulnResult[1]
	}
	_ = resultID // TODO what to do with this?

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
			createOrUpdateDetection(
				job.db,
				job.lstream,
				job.detectionStatuses,
				deviceInDb,
				vulnInfo,
				detectionFromScanner,
				assetID,
				decomIgnoreID,
				job.config.OrganizationID(),
				job.insources.SourceID(),
				groupID,
				"", // this field is populated recursively when child detections are present
				job.getExceptionID,
			)
		} else {
			job.lstream.Send(log.Error("could not find vulnerability in database", fmt.Errorf("[%s] does not have an entry in the database", vulnID)))
		}

	} else {
		job.lstream.Send(log.Errorf(err, "Error while gathering vulnerability info for [%s]", vulnID))
	}

	return err
}

func (job *AssetSyncJob) getExceptionID(assetID string, deviceInDb domain.Device, port string, vulnInfo domain.VulnerabilityInfo, detectionFromScanner domain.Detection) (exceptionID string, dontUpdateExceptionID bool) {
	dontUpdateExceptionID = false // default value, but explicitly setting as to not imply it was unintentional. the asset sync should update the exception ID

	if job.deviceIDToVulnIDToException[assetID] != nil {
		if job.deviceIDToVulnIDToException[assetID][fmt.Sprintf("%s;%s", vulnInfo.SourceVulnID(), port)] != nil {

			var possibleMatch = job.deviceIDToVulnIDToException[assetID][fmt.Sprintf("%s;%s", vulnInfo.SourceVulnID(), port)]

			if possibleMatch.TypeID() != domain.Exception || tord1970(possibleMatch.DueDate()).After(time.Now()) { // only want to skip exceptions that have passed their due dates
				exceptionID = job.deviceIDToVulnIDToException[assetID][fmt.Sprintf("%s;%s", vulnInfo.SourceVulnID(), port)].ID()
			}
		}
	}

	if len(exceptionID) == 0 && !stringInSlice(job.Payload.GroupIDsExemptFromGlobalExceptions, sord(deviceInDb.GroupID())) {
		for _, globalException := range job.globalExceptions {

			if globalException.osRegex != nil || globalException.hostnameRegex != nil {
				if len(globalException.exception.VulnerabilityID()) == 0 || globalException.exception.VulnerabilityID() == vulnInfo.SourceVulnID() {

					var failedOSCheck, failedHostnameCheck bool

					if globalException.osRegex != nil {
						if !globalException.osRegex.Match([]byte(deviceInDb.OS())) {
							failedOSCheck = true
						}
					}

					if globalException.hostnameRegex != nil {
						if !globalException.hostnameRegex.Match([]byte(deviceInDb.HostName())) {
							failedHostnameCheck = true
						}
					}

					if !failedOSCheck && !failedHostnameCheck {
						exceptionID = globalException.exception.ID()
						break
					}
				}
			}

		}
	}

	return exceptionID, dontUpdateExceptionID
}

// This method creates a detection entry if one does not exist, and updates the entry if one does
func createOrUpdateDetection(db domain.DatabaseConnection, lstream log.Logger, detectionStatuses []domain.DetectionStatus, deviceInDb domain.Device, vulnInfo domain.VulnerabilityInfo, detectionFromScanner domain.Detection, assetID string, decomIgnoreID string, orgID string, sourceID string, groupID string, parentDetectionID string, getExceptionID func(string, domain.Device, string, domain.VulnerabilityInfo, domain.Detection) (string, bool)) {
	var err error

	var detectionInDB domain.DetectionInfo
	detectionInDB, err = db.GetDetectionInfo(sord(deviceInDb.SourceID()), vulnInfo.ID(), detectionFromScanner.Port(), detectionFromScanner.Protocol())
	if err == nil {
		var detectionStatus domain.DetectionStatus
		if detectionStatus = getDetectionStatus(detectionStatuses, detectionFromScanner.Status()); detectionStatus != nil {

			var port string
			if detectionFromScanner.Port() > 0 || len(detectionFromScanner.Protocol()) > 0 {
				port = fmt.Sprintf("%d %s", detectionFromScanner.Port(), detectionFromScanner.Protocol())
			}
			var exceptionID = decomIgnoreID
			var dontUpdateExceptionID bool
			if len(exceptionID) == 0 {
				exceptionID, dontUpdateExceptionID = getExceptionID(assetID, deviceInDb, port, vulnInfo, detectionFromScanner)
			}

			if detectionInDB == nil {
				createDetection(db, lstream, orgID, sourceID, detectionFromScanner, exceptionID, deviceInDb, vulnInfo, assetID, detectionStatus.ID(), parentDetectionID)
			} else {
				if len(exceptionID) == 0 && dontUpdateExceptionID {
					// we set the exceptionID to the value already in the database
					// this block hits when the caller didn't check for an Ignore entry while trying to update the detection data
					exceptionID = sord(detectionInDB.IgnoreID())
				}

				var canSkipUpdate bool
				if detectionFromScanner.LastUpdated() != nil && !detectionInDB.Updated().IsZero() && !detectionFromScanner.LastUpdated().IsZero() {
					if detectionInDB.Updated().After(*detectionFromScanner.LastUpdated()) {
						canSkipUpdate = true
					}
				}

				// even if the detection hasn't been updated in the scanner, if we've had a new exception added (or the exception was removed), we need to update the detection
				if canSkipUpdate && ((len(sord(detectionInDB.IgnoreID())) == 0 && len(exceptionID) > 0) || (len(sord(detectionInDB.IgnoreID())) > 0 && len(exceptionID) == 0)) {
					canSkipUpdate = false
				}

				if !canSkipUpdate {
					_, _, err = db.UpdateDetection(
						detectionInDB.ID(),
						sord(deviceInDb.SourceID()),
						vulnInfo.ID(),
						detectionFromScanner.Port(),
						detectionFromScanner.Protocol(),
						exceptionID,
						detectionFromScanner.TimesSeen(),
						detectionStatus.ID(),
						tord1970(detectionFromScanner.LastFound()),
						tord1970(detectionFromScanner.LastUpdated()),
						tord1970(nil),
					)

					if err == nil {
						lstream.Send(log.Infof("Updated detection for device/vuln [%v|%v]", assetID, vulnInfo.ID()))
					} else {
						lstream.Send(log.Errorf(err, "Error while updating detection for device/vuln [%v|%v]", assetID, vulnInfo.ID()))
					}
				} else {
					lstream.Send(log.Debugf("Skipping detection update for device/vuln [%v|%v] [%v after %v]", assetID, vulnInfo.ID(), detectionInDB.Updated(), tord1970(detectionFromScanner.LastUpdated())))
				}
			}
		} else {
			lstream.Send(log.Errorf(err, "could not find detection status with name [%s]", detectionFromScanner.Status()))
		}
	} else {
		lstream.Send(log.Debugf("Detection already exists for device/vuln [%v|%v]", assetID, vulnInfo.ID()))
	}

	processChildDetections(db, lstream, detectionStatuses, assetID, vulnInfo, detectionFromScanner, orgID, sourceID, groupID, getExceptionID)
}

func processChildDetections(db domain.DatabaseConnection, lstream log.Logger, detectionStatuses []domain.DetectionStatus, assetID string, vulnInfo domain.VulnerabilityInfo, detectionFromScanner domain.Detection, orgID string, sourceID string, groupID string, getExceptionID func(string, domain.Device, string, domain.VulnerabilityInfo, domain.Detection) (string, bool)) {
	if len(detectionFromScanner.ChildDetections()) == 0 {
		return
	}

	parentDetectionInDB, err := db.GetDetectionInfo(assetID, vulnInfo.ID(), detectionFromScanner.Port(), detectionFromScanner.Protocol())
	if err != nil || parentDetectionInDB == nil {
		lstream.Send(log.Errorf(err, "error while gathering parent detection"))
		return
	}

	for _, childDetection := range detectionFromScanner.ChildDetections() {
		device, err := childDetection.Device()
		if err != nil {
			lstream.Send(log.Errorf(err, "error while loading device entry from scanner for child detection [%s]"))
			continue
		}

		assetID := sord(device.SourceID())
		decomIgnoreID, err := getDecommIgnoreEntryForAsset(db, lstream, orgID, assetID, sourceID, []domain.Detection{childDetection})
		if err != nil {
			lstream.Send(log.Errorf(err, "error while loading decomm ignore entry"))
		}

		deviceInDB, err := addDeviceInformationToDB(db, lstream, sourceID, orgID, device, groupID)
		if err != nil {
			lstream.Send(log.Errorf(err, "error while loading device entry from db for child detection [%s]"))
			continue
		}

		createOrUpdateDetection(db, lstream, detectionStatuses, deviceInDB, vulnInfo, childDetection, assetID, decomIgnoreID, orgID, sourceID, groupID, parentDetectionInDB.ID(), getExceptionID)
	}
}

// This method creates the detection entry in the database
func createDetection(db domain.DatabaseConnection, lstream log.Logger, orgID string, sourceID string, vuln domain.Detection, exceptionID string, deviceInDb domain.Device, vulnInfo domain.VulnerabilityInfo, assetID string, detectionStatusID int, parentDetectionID string) {
	var err error

	var detected *time.Time
	if detected, err = vuln.Detected(); err == nil {
		if detected != nil {
			_, _, err = db.CreateDetection(
				orgID,
				sourceID,
				sord(deviceInDb.SourceID()),
				vulnInfo.ID(),
				exceptionID,
				*detected,
				tord1970(vuln.LastFound()),
				tord1970(vuln.LastUpdated()),
				vuln.Proof(),
				vuln.Port(),
				vuln.Protocol(),
				iord(vuln.ActiveKernel()),
				detectionStatusID,
				vuln.TimesSeen(),
				tord1970(nil),
				parentDetectionID,
			)
		} else {
			err = fmt.Errorf("could not find the time of the detection")
		}

	} else {
		err = fmt.Errorf("error while gathering date of detection - %v", err.Error())
	}

	if err == nil {
		lstream.Send(log.Infof("Created detection for device/vuln [%v|%v]", assetID, vulnInfo.ID()))
	} else {
		lstream.Send(log.Errorf(err, "Error while creating detection for device/vuln [%v|%v]", assetID, vulnInfo.ID()))
	}
}

func getDetectionStatus(detectionStatuses []domain.DetectionStatus, status string) (detectionStatus domain.DetectionStatus) {
	for _, potentialMatch := range detectionStatuses {
		if strings.ToLower(status) == strings.ToLower(potentialMatch.Status()) {
			detectionStatus = potentialMatch
			break
		}
	}

	return detectionStatus
}

// This method creates an entry in the database for the operating system type. It then returns the entry so that the id of the OST
// may be used for foreign key references
func grabAndCreateOsType(db domain.DatabaseConnection, operatingSystem string) (output domain.OperatingSystemType, err error) {
	if len(operatingSystem) > 0 {
		output, err = db.GetOperatingSystemType(operatingSystem)
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

func (job *AssetSyncJob) createAssetGroupInDB(groupID string, scannerSourceID string, scannerSourceConfigID string) (err error) {
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

type compiledException struct {
	exception     domain.Ignore
	osRegex       *regexp.Regexp
	hostnameRegex *regexp.Regexp
	deviceIDRegex *regexp.Regexp
}

func (job *AssetSyncJob) preloadIgnores() (globals []compiledException, deviceIDToVulnIDToException map[string]map[string]domain.Ignore, err error) {
	deviceIDToVulnIDToException = make(map[string]map[string]domain.Ignore)
	globals, err = loadGlobalExceptions(job.db, job.config.OrganizationID(), job.insources.SourceID())

	if err == nil {
		var specificExceptions []domain.Ignore
		if specificExceptions, err = job.db.GetExceptionsByOrg(job.config.OrganizationID()); err == nil {
			for _, exception := range specificExceptions {
				if len(exception.DeviceID()) > 0 && len(exception.VulnerabilityID()) > 0 {
					if deviceIDToVulnIDToException[exception.DeviceID()] == nil {
						deviceIDToVulnIDToException[exception.DeviceID()] = make(map[string]domain.Ignore)
					}

					deviceIDToVulnIDToException[exception.DeviceID()][fmt.Sprintf("%s;%s", exception.VulnerabilityID(), exception.Port())] = exception
				}
			}
		}
	}

	return globals, deviceIDToVulnIDToException, err
}

func loadGlobalExceptions(db domain.DatabaseConnection, orgID string, sourceID string) (globals []compiledException, err error) {
	globals = make([]compiledException, 0)
	var globalExceptions []domain.Ignore
	if globalExceptions, err = db.GetGlobalExceptions(orgID, sourceID); err == nil {
		for _, globalException := range globalExceptions {

			var osRegex *regexp.Regexp
			if len(sord(globalException.OSRegex())) > 0 {
				if osRegex, err = regexp.Compile(sord(globalException.OSRegex())); err != nil {
					err = fmt.Errorf("error while compiling regex for global ignore entry [%s]", globalException.ID())
					break
				}
			}

			var hostnameRegex *regexp.Regexp
			if len(sord(globalException.HostnameRegex())) > 0 {
				if hostnameRegex, err = regexp.Compile(sord(globalException.HostnameRegex())); err != nil {
					err = fmt.Errorf("error while compiling regex for global ignore entry [%s]", globalException.ID())
					break
				}
			}

			var deviceIDRegex *regexp.Regexp
			if len(sord(globalException.DeviceIDRegex())) > 0 {
				if deviceIDRegex, err = regexp.Compile(sord(globalException.DeviceIDRegex())); err != nil {
					err = fmt.Errorf("error while compiling regex for global ignore entry [%s]", globalException.ID())
					break
				}
			}

			if osRegex != nil || hostnameRegex != nil || deviceIDRegex != nil {
				globals = append(globals, compiledException{
					exception:     globalException,
					osRegex:       osRegex,
					hostnameRegex: hostnameRegex,
					deviceIDRegex: deviceIDRegex,
				})
			} else {
				err = fmt.Errorf("neither os regex nor hostname regex compiled for global ignore entry [%s]", globalException.ID())
				break
			}
		}
	} else {
		err = fmt.Errorf("error while loading global exceptions - %s", err.Error())
	}
	return globals, err
}

func stringInSlice(vals []string, lookFor string) (inSlice bool) {
	for _, val := range vals {
		if val == lookFor {
			inSlice = true
			break
		}
	}

	return inSlice
}
