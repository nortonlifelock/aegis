package implementations

import (
	"context"
	"fmt"

	"github.com/nortonlifelock/aegis/internal/integrations"
	"github.com/nortonlifelock/aegis/pkg/domain"
	"github.com/nortonlifelock/aegis/pkg/log"
)

// TODO setup email log if asset is found in aws that doesn't already exist

/*
	If the asset is created first by the cloud sync, use the os info from the cloud sync
	the asset sync job needs to update that os info with what they find because it's more specific
	add the instance id and the ip to the device field
	the asset sync job needs to add the ip to the table
*/

// CloudSyncJob is the struct used to run the job, which is responsible for grabbing tag information from a cloud service provider
// and storing it in the database
type CloudSyncJob struct {
	id          string
	payloadJSON string
	ctx         context.Context
	db          domain.DatabaseConnection
	lstream     log.Logger
	appconfig   domain.Config
	config      domain.JobConfig
	insources   []domain.SourceConfig
	outsources  []domain.SourceConfig
}

const unknown = "Unknown"

// common tag keys
const (
	operatingSystem = "OperatingSystem"
	instanceID      = "InstanceId"
)

// Process pulls tag information associated with devices that are scanned in cloud service providers (e.g. AWS/Azure)
// the tags are used within the ticketing job to include additional information, or override information in a ticket
func (job *CloudSyncJob) Process(ctx context.Context, id string, appconfig domain.Config, db domain.DatabaseConnection, lstream log.Logger, payload string, jobConfig domain.JobConfig, inSource []domain.SourceConfig, outSource []domain.SourceConfig) (err error) {

	var ok bool
	if job.ctx, job.id, job.appconfig, job.db, job.lstream, job.payloadJSON, job.config, job.insources, job.outsources, ok = validInputsMultipleSources(ctx, id, appconfig, db, lstream, payload, jobConfig, inSource, outSource); ok {

		// TODO ensure the source of each cloud connection in the same
		var cloudConnections = make([]integrations.CloudServiceConnection, 0)

		for _, insource := range job.insources {
			var connection integrations.CloudServiceConnection
			if connection, err = integrations.GetCloudServiceConnection(job.ctx, job.db, insource.Source(), insource, job.appconfig, job.lstream); err == nil {
				cloudConnections = append(cloudConnections, connection)
			} else {
				job.lstream.Send(log.Error("error while establishing connection", err))
				break
			}
		}

		if err == nil {

			// first we handle IP tag mapping
			for _, connection := range cloudConnections {
				select {
				case <-job.ctx.Done():
					return err
				default:
				}
				job.lstream.Send(log.Debug("Gathering IP-Tag mappings"))

				// the first key for ipToKeyToValue is an ip which returns a map holding key-value pairs
				// the key is the name of the tag, and the value is the value within the tag
				var ipToKeyToValue map[domain.CloudIP]map[string]string
				ipToKeyToValue, err = connection.GetIPTagMapping()
				if err == nil {
					job.lstream.Send(log.Debug("IP-Tag mapping gathered"))

					// this method takes the tag information and creates entries in the Tag/TagKey tables
					job.processTagsForDB(connection, ipToKeyToValue)
				} else {
					job.lstream.Send(log.Error("error while gathering ip tag mapping", err))
				}
			}
		} else {
			job.lstream.Send(log.Errorf(err, "error while creating cloud connections"))
		}
	} else {
		err = fmt.Errorf("input validation failed for %v", job.id)
		job.lstream.Send(log.Errorf(err, "error during job startup"))
	}

	return err
}

func (job *CloudSyncJob) processTagsForDB(connection integrations.CloudServiceConnection, ipToKeyToValue map[domain.CloudIP]map[string]string) {
	var err error
	job.lstream.Send(log.Info("Processing tag keys in database"))

	// the key to the map is the tag key (the name of the tag) and the value is the primary key (the id) for that tag key in the database
	var tagKeyToDbID map[string]string
	tagKeyToDbID, err = job.createTagKeyDbEntriesAndIDMap(connection)
	if err == nil {
		job.lstream.Send(log.Info("Tag keys finished processing"))

		for ip, keyToValue := range ipToKeyToValue {
			select {
			case <-job.ctx.Done():
				return
			default:
			}

			// a cloud IP may be attached to more than one device in the DB (can happen if a scanner duplicates a device in a different asset group)
			var devices []domain.Device
			// createOrUpdateDevice creates an entry in the database for the device if one does not exist
			devices, err = job.createOrUpdateDevice(ip, keyToValue)
			if err == nil {
				for _, device := range devices {
					// add the tags returned from the cloud service provider to the database
					job.createOrUpdateTagsForDevice(keyToValue, tagKeyToDbID, device, ip)
				}
			} else {
				job.lstream.Send(log.Error("error while creating/updating device", err))
			}
		}
	} else {
		job.lstream.Send(log.Error("error while inserting tag keys into database", err))
	}
}

func (job *CloudSyncJob) createOrUpdateTagsForDevice(keyToValue map[string]string, tagKeyToDbID map[string]string, device domain.Device, ip domain.CloudIP) {
	job.lstream.Send(log.Infof("Processing device [%s|%s] with IP [%s]", sord(device.SourceID()), device.IP(), ip.IP()))
	var err error

	for tagKey, tagValue := range keyToValue {
		var tagKeyDbID = tagKeyToDbID[tagKey]
		if len(tagKeyDbID) > 0 {
			var tag domain.Tag
			tag, err = job.db.GetTagByDeviceAndTagKey(device.ID(), tagKeyDbID)
			if err == nil {
				if tag == nil {

					_, _, err = job.db.CreateTag(device.ID(), tagKeyDbID, tagValue)
					if err == nil {
						job.lstream.Send(log.Infof("Created tag [%s|%s] for device [%v]", tagKey, tagValue, sord(device.SourceID())))
					} else {
						job.lstream.Send(log.Errorf(err, "error while creating tag for device [%v]", sord(device.SourceID())))
					}

				} else {
					if tag.Value() != tagValue {
						_, _, err = job.db.UpdateTag(device.ID(), tagKeyDbID, tagValue)
						if err == nil {
							job.lstream.Send(log.Infof("Updated tag [%s|%s] for device [%v]", tagKey, tagValue, device.SourceID()))
						} else {
							job.lstream.Send(log.Errorf(err, "error while updating tag for device [%v]", sord(device.SourceID())))
						}
					}
				}

			} else {
				job.lstream.Send(log.Errorf(err, "error while finding tag for device and key [%v|%s]", device.SourceID(), tagKeyDbID))
			}
		} else {
			job.lstream.Send(log.Error("error while processing tag", fmt.Errorf("could not find id for [%s]", tagKey)))
		}

	}
}

func (job *CloudSyncJob) createTagKeyDbEntriesAndIDMap(connection integrations.CloudServiceConnection) (mapTagKeyToItsDbID map[string]string, err error) {
	mapTagKeyToItsDbID = make(map[string]string)
	var tagNames []string
	tagNames, err = connection.GetAllTagNames()

	if err == nil {
		for index := range tagNames {
			tagName := tagNames[index]

			if len(tagName) > 0 {
				var tagKey domain.TagKey
				tagKey, err = job.db.GetTagKeyByKey(tagName)
				if err == nil {
					if tagKey == nil {
						_, _, err = job.db.CreateTagKey(tagName)
						if err == nil {
							tagKey, err = job.db.GetTagKeyByKey(tagName)
						}
					}
				}

				if tagKey == nil && err == nil {
					err = fmt.Errorf("failed to create tag key in database")
				}

				if err == nil {
					mapTagKeyToItsDbID[tagName] = tagKey.ID()
				}
			}

			if err != nil {
				break
			}

		}
	}

	return mapTagKeyToItsDbID, err
}

// TODO how should we correlate cloud devices to scanner devices?
// The device for the IP will already exist if the asset sync job found it first. If the asset sync job did not find it first, this method will create
// an entry in the database
func (job *CloudSyncJob) createOrUpdateDevice(ip domain.CloudIP, keyToValue map[string]string) (devices []domain.Device, err error) {
	if devices, err = job.getDevicesForIP(ip, keyToValue); err == nil && len(devices) > 0 {

		for _, device := range devices {
			// device already exists but the instance id isn't set
			// (likely occurred because asset was created by asset sync job)
			if (len(sord(device.InstanceID())) == 0 && len(keyToValue[instanceID]) > 0) || (len(sord(device.Region())) == 0 && len(ip.Region()) > 0) {
				_, _, err = job.db.UpdateInstanceIDOfDevice(device.ID(), keyToValue[instanceID], job.insources[0].SourceID(), ip.State(), ip.Region(), job.config.OrganizationID())
				if err != nil {
					err = fmt.Errorf("error while updating instance information for device [%s] - %s", device.ID(), err.Error())
				}
			}

			switch ip.State() {
			case domain.DeviceRunning:
			case domain.DeviceStopped:
			case domain.DeviceDecommed:
			case domain.DeviceDeallocated:
			case domain.DeviceUnknown:
			default:
				job.lstream.Send(log.Errorf(nil, "unrecognized state %v", ip.State()))
			}

			_, _, err = job.db.UpdateStateOfDevice(device.ID(), ip.State(), job.config.OrganizationID()) // TODO is terminated a state?
			if err != nil {
				err = fmt.Errorf("error while updating state information for device [%s] - %s", device.ID(), err.Error())
			}
		}
	}

	return devices, err
}

func (job *CloudSyncJob) getDevicesForIP(ip domain.CloudIP, keyToValue map[string]string) (devices []domain.Device, err error) {
	devices = make([]domain.Device, 0)

	// GetDeviceByCloudSourceIDAndIP uses the AssetGroup mapping (links cloud subscriptions to vulnerability asset groups) to identify
	var devicesByCloudSource, devicesByInstanceID []domain.Device

	devicesByCloudSource, err = job.db.GetDeviceByCloudSourceIDAndIP(ip.IP(), job.insources[0].SourceID(), job.config.OrganizationID())
	if err == nil {
		devicesByInstanceID, err = job.db.GetDeviceByInstanceID(ip.InstanceID(), job.config.OrganizationID())
	}

	if err == nil {
		var seen = make(map[string]bool)
		for _, device := range devicesByCloudSource {
			if !seen[device.ID()] {
				seen[device.ID()] = true
				devices = append(devices, device)
			}
		}

		for _, device := range devicesByInstanceID {
			if !seen[device.ID()] {
				seen[device.ID()] = true
				devices = append(devices, device)
			}
		}

		if len(devices) == 0 {
			job.lstream.Send(log.Criticalf(err, "device with IP [%s] for org [%s] was not found in database by Cloud Sync Job", ip.IP(), job.config.OrganizationID()))
			var newDevice domain.Device
			newDevice, err = job.createAndReturnDevice(ip, keyToValue)
			devices = append(devices, newDevice)
		}
	}

	return devices, err
}

func (job *CloudSyncJob) createAndReturnDevice(ip domain.CloudIP, keyToValue map[string]string) (device domain.Device, err error) {
	var osType domain.OperatingSystemType
	var OS string
	osType, OS, err = job.createAndReturnOSType(keyToValue)

	if err == nil {
		if len(ip.IP()) > 0 {
			if len(keyToValue[instanceID]) == 0 {
				job.lstream.Send(log.Warningf(nil, "InstanceId of length 0 found for ip [%s]", ip.IP()))
			}

			_, _, err = job.db.CreateAssetWithIPInstanceID(ip.State(), ip.IP(), ip.MAC(), job.insources[0].SourceID(), keyToValue[instanceID], ip.Region(), job.config.OrganizationID(), OS, osType.ID())
			if err == nil {
				var devices []domain.Device
				devices, err = job.db.GetDeviceByInstanceID(ip.InstanceID(), job.config.OrganizationID())
				if err == nil {
					if len(devices) > 0 {
						device = devices[0]
					} else {
						err = fmt.Errorf("could not find recently created device for [%s] in database", ip.IP())
					}
				}
			} else {
				err = fmt.Errorf("error while creating asset for ip [%s] - %s", ip.IP(), err.Error())
			}
		} else {
			err = fmt.Errorf("ip of length 0 with associated tag info")
		}

	} else {
		err = fmt.Errorf("error while managing operating system type - %s", err.Error())
	}

	return device, err
}

func (job *CloudSyncJob) createAndReturnOSType(keyToValue map[string]string) (osType domain.OperatingSystemType, OS string, err error) {
	OS = keyToValue[operatingSystem]

	// check to see if the osType already exists in db
	osType, err = job.db.GetOperatingSystemType(OS)
	if err == nil {
		if osType == nil {
			err = fmt.Errorf("could not find operating system type for [%s]", OS)
		}
	} else {
		err = fmt.Errorf("error while gathering operating system type - %s", err.Error())
	}

	return osType, OS, err
}
