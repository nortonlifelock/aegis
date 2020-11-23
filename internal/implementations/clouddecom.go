package implementations

import (
	"context"
	"encoding/json"
	"fmt"

	"strings"
	"sync"
	"time"

	"github.com/nortonlifelock/aegis/internal/integrations"
	"github.com/nortonlifelock/aegis/pkg/domain"
	"github.com/nortonlifelock/aegis/pkg/log"
)

// CloudDecommissionJob pulls a history of tracked assets from the database and compares that to a list of live assets as reported
type CloudDecommissionJob struct {
	id          string
	payloadJSON string
	Payload     *CloudDecommissionPayload
	ctx         context.Context
	db          domain.DatabaseConnection
	lstream     log.Logger
	appconfig   domain.Config
	config      domain.JobConfig
	insources   []domain.SourceConfig
	outsources  []domain.SourceConfig
}

type CloudDecommissionPayload struct {
	// OnlyCheckIPs is an optional field. If it is not empty, a decommission check will only be done against these specific IPs as opposed to the entire cloud inventory
	OnlyCheckIPs []string `json:"only_check_ips"`

	DecommOnStoppedState bool `json:"decommission_on_stopped_state"`
}

// buildPayload loads the Payload from the job history into the Payload object
func (job *CloudDecommissionJob) buildPayload(pjson string) (err error) {
	job.Payload = &CloudDecommissionPayload{}

	if len(pjson) > 0 {
		err = json.Unmarshal([]byte(pjson), job.Payload)
	} else {
		err = fmt.Errorf("no Payload provided to job")
	}

	return err
}

// Process grabs a history of the devices tracked by the database. All devices belonging to a cloud service (AWS/Azure) are checked to see if they are still existent in the cloud inventory of that service. If they do not exist, the device is decommissioned in the database and its tickets are closed
// It also grabs the devices that were previously decommissioned, and verifies that they still no longer exist in the cloud inventory. If they are discovered to be alive again, their entry in the ignore table is deleted
func (job *CloudDecommissionJob) Process(ctx context.Context, id string, appconfig domain.Config, db domain.DatabaseConnection, lstream log.Logger, payload string, jobConfig domain.JobConfig, inSource []domain.SourceConfig, outSource []domain.SourceConfig) (err error) {
	var ok bool
	if job.ctx, job.id, job.appconfig, job.db, job.lstream, job.payloadJSON, job.config, job.insources, job.outsources, ok = validInputsMultipleSources(ctx, id, appconfig, db, lstream, payload, jobConfig, inSource, outSource); ok {
		if err = job.buildPayload(job.payloadJSON); err == nil {
			// TODO ensure the source of each cloud connection in the same
			var cloudConnections = make([]integrations.CloudServiceConnection, 0)

			for _, insource := range job.insources {
				var connection integrations.CloudServiceConnection
				if connection, err = integrations.GetCloudServiceConnection(job.db, insource.Source(), insource, job.appconfig, job.lstream); err == nil {
					cloudConnections = append(cloudConnections, connection)
				} else {
					job.lstream.Send(log.Error("error while establishing connection", err))
					break
				}
			}

			if err == nil {
				job.decommissionCloudAssets(cloudConnections)
			}
		} else {
			job.lstream.Send(log.Errorf(err, "error while building payload"))
		}
	} else {
		err = fmt.Errorf("failed validation")
		job.lstream.Send(log.Errorf(err, "input validation failed"))
	}

	return err
}

func (job *CloudDecommissionJob) decommissionCloudAssets(cloudConnections []integrations.CloudServiceConnection) {
	var err error

	var ticketingEngine integrations.TicketingEngine
	if ticketingEngine, err = integrations.GetEngine(job.ctx, job.outsources[0].Source(), job.db, job.lstream, job.appconfig, job.outsources[0]); err == nil {
		var allIPs = make([]domain.CloudIP, 0)
		// we query the cloud connection for the asset inventory so we can identify
		// which asset have been decommissioned
		for _, connection := range cloudConnections {
			var ipsForSubscription []domain.CloudIP
			// could map region to IP address in this method
			if ipsForSubscription, err = connection.IPAddresses(); err == nil {
				allIPs = append(allIPs, ipsForSubscription...)
			} else {
				break
			}
		}

		if err == nil {

			// Can use those to find which devices are missing
			var historyOfDevices []domain.Device
			if len(job.Payload.OnlyCheckIPs) == 0 {
				// no IPs specified, so we check the entire cloud inventory
				historyOfDevices, err = job.db.GetDevicesByCloudSourceID(job.insources[0].SourceID(), job.config.OrganizationID())
			} else {
				historyOfDevices = make([]domain.Device, 0)
				for _, ip := range job.Payload.OnlyCheckIPs {
					var devices []domain.Device
					if devices, err = job.db.GetDeviceByCloudSourceIDAndIP(ip, job.insources[0].SourceID(), job.config.OrganizationID()); err == nil && devices != nil {
						historyOfDevices = append(historyOfDevices, devices...)
					} else {
						job.lstream.Send(log.Warningf(err, "could not load device for IP and Cloud sources [%s|%s]", ip, job.insources[0].SourceID()))
					}
				}
			}

			if err == nil {
				deviceIDToDecommissionedDevice := job.findDecommissionedDevices(historyOfDevices, allIPs)
				job.markDevicesAsDecommissionedInDatabase(deviceIDToDecommissionedDevice)

				var orgInfo domain.Organization
				if orgInfo, err = job.db.GetOrganizationByID(job.config.OrganizationID()); err == nil {

					var assetGroups []domain.AssetGroup
					if assetGroups, err = job.db.GetAssetGroupsByCloudSource(job.config.OrganizationID(), job.insources[0].SourceID()); err == nil {

						var sourceIDToSource map[string]domain.Source
						if sourceIDToSource, err = job.getSourceMap(); err == nil {

							var tickets chan domain.Ticket
							if tickets, err = job.getTicketsForDecommCheck(assetGroups, sourceIDToSource, ticketingEngine, orgInfo); err == nil {
								job.closeTicketsForDecommissionedAssets(tickets, deviceIDToDecommissionedDevice, ticketingEngine, sourceIDToSource)
							}

							job.findIncorrectlyDecommissionedAssets(deviceIDToDecommissionedDevice)
						} else {
							job.lstream.Send(log.Errorf(err, "error while loading sources from database"))
						}
					} else {
						job.lstream.Send(log.Errorf(err, "error while loading asset groups from database"))
					}
				} else {
					job.lstream.Send(log.Errorf(err, "error while loading organization info for [%v]", job.config.OrganizationID()))
				}
			} else {
				job.lstream.Send(log.Errorf(err, "error while grabbing history of devices"))
			}
		} else {
			job.lstream.Send(log.Errorf(err, "error while gathering active IP addresses"))
		}
	} else {
		job.lstream.Send(log.Errorf(err, "error while creating ticketing connection"))
	}
}

func (job *CloudDecommissionJob) getTicketsForDecommCheck(assetGroups []domain.AssetGroup, sourceIDToSource map[string]domain.Source, ticketingEngine integrations.TicketingEngine, orgInfo domain.Organization) (tickets chan domain.Ticket, err error) {
	tickets = make(chan domain.Ticket)

	wg := &sync.WaitGroup{}
	for _, assetGroup := range assetGroups {

		if vulnSource := sourceIDToSource[assetGroup.ScannerSourceID()]; vulnSource != nil {

			var groupTickets <-chan domain.Ticket
			if groupTickets, err = ticketingEngine.GetOpenTicketsByGroupID(vulnSource.Source(), orgInfo.Code(), assetGroup.GroupID()); err == nil {

				wg.Add(1)
				go func(groupTickets <-chan domain.Ticket) {
					defer handleRoutinePanic(job.lstream)
					defer wg.Done()

					for {
						if ticket, ok := <-groupTickets; ok {
							select {
							case <-job.ctx.Done():
								return
							case tickets <- ticket:
							}
						} else {
							break
						}
					}
				}(groupTickets)
			} else {
				job.lstream.Send(log.Errorf(err, "error while loading tickets"))
				break
			}
		} else {
			err = fmt.Errorf("could not find source with ID [%v]", assetGroup.ScannerSourceID())
			job.lstream.Send(log.Errorf(err, "error while gathering source"))
			break
		}
	}

	go func() {
		defer handleRoutinePanic(job.lstream)
		defer close(tickets)
		wg.Wait()
	}()

	return tickets, err
}

func (job *CloudDecommissionJob) findIncorrectlyDecommissionedAssets(deviceIDToDecommDevice map[string]domain.Device) {
	var historyOfDeviceInfo []domain.DeviceInfo
	var err error
	if historyOfDeviceInfo, err = job.db.GetDevicesInfoByCloudSourceID(job.insources[0].Source(), job.config.OrganizationID()); err == nil {

		var devicesMarkedAsDecommedInDB = make([]domain.DeviceInfo, 0)
		for _, deviceInfo := range historyOfDeviceInfo {
			if sord(deviceInfo.State()) == domain.DeviceDecommed {
				devicesMarkedAsDecommedInDB = append(devicesMarkedAsDecommedInDB, deviceInfo)
			}
		}

		for _, decommDeviceInDB := range devicesMarkedAsDecommedInDB {

			deviceID := sord(decommDeviceInDB.SourceID())

			if len(deviceID) > 0 {
				// if a device ID is a valid key (returning a non-nil value) in the deviceIDToDecommDevice map, that means the device is decommissioned
				if deviceIDToDecommDevice[deviceID] == nil {
					// because the device ID was not a valid key, this device seen by the cloud service and does not appear to be decommissioned
					// the device that we checked against the map was marked as decommissioned in the database, so we must delete it's ignore entry in the db

					_, _, err = job.db.DeleteIgnoreForDevice(
						sord(decommDeviceInDB.ScannerSourceID()),
						deviceID,
						job.config.OrganizationID(),
					)

					if err == nil {
						job.lstream.Send(log.Warningf(err, "[%v] was marked as decommissioned but was found in the live asset inventory", deviceID))
					} else {
						job.lstream.Send(log.Errorf(err, "error while deleting decommission entry for falsely decommissioned asset [%v]", deviceID))
					}
				}
			}
		}

	} else {
		job.lstream.Send(log.Errorf(err, "error while loading device history while identifying incorrectly decommissioned assets"))
	}
}

func (job *CloudDecommissionJob) closeTicketsForDecommissionedAssets(tickets <-chan domain.Ticket, deviceIDToDecommDevice map[string]domain.Device, ticketingEngine integrations.TicketingEngine, sourceIDToSource map[string]domain.Source) {
	var deviceAlreadyDecommedInDB sync.Map

	wg := &sync.WaitGroup{}
	func() {
		for {
			select {
			case <-job.ctx.Done():
				return
			case tic, ok := <-tickets:
				if ok {

					// the ticket has a device id that was found to be decommissioned
					if deviceIDToDecommDevice[tic.DeviceID()] != nil {
						wg.Add(1)
						go func(tic domain.Ticket) {
							defer handleRoutinePanic(job.lstream)
							defer wg.Done()

							if err := ticketingEngine.Transition(
								tic,
								ticketingEngine.GetStatusMap(domain.StatusClosedDecommissioned),
								"Asset decommissioned as it was not located in the cloud asset inventory",
								sord(tic.AssignedTo())); err == nil {
								job.lstream.Send(log.Infof("%v marked as decommissioned as it's IP [%v] was not in the AWS inventory", tic.Title(), sord(tic.IPAddress())))
							} else {
								job.lstream.Send(log.Errorf(err, "error while marking %v as decommissioned", tic.Title()))
							}

							_, loaded := deviceAlreadyDecommedInDB.LoadOrStore(tic.DeviceID(), true)
							if !loaded { // if the val/ue wasn't loaded, that means this was the first time the device was processed
								var scannerSourceID string
								scannerSourceID, err := getSourceIDFromMethodOfDiscovery(sord(tic.MethodOfDiscovery()), sourceIDToSource)
								if err == nil {
									job.createIgnoreEntry(tic.DeviceID(), scannerSourceID, job.config.OrganizationID())
								} else {
									job.lstream.Send(log.Errorf(err, "error while creating ignore entry for [%s]", tic.DeviceID()))
								}
							}
						}(tic)
					} else if sord(tic.Status()) == ticketingEngine.GetStatusMap(domain.StatusResolvedDecom) {
						// this block hits if we have a ticket that was marked as resolved-decommissioned, but was
						// found in the cloud asset inventory

						wg.Add(1)
						go func(tic domain.Ticket) {
							defer handleRoutinePanic(job.lstream)
							defer wg.Done()

							if err := ticketingEngine.Transition(
								tic,
								ticketingEngine.GetStatusMap(domain.StatusReopened),
								"Ticket reopened as it was found in the cloud asset inventory",
								sord(tic.AssignedTo())); err == nil {
								job.lstream.Send(log.Infof("%v reopened as it's IP [%v] was in the AWS inventory", tic.Title(), sord(tic.IPAddress())))
							} else {
								job.lstream.Send(log.Errorf(err, "error while marking %v as reopened", tic.Title()))
							}
						}(tic)
					}
				} else {
					return
				}
			}
		}
	}()
	wg.Wait()
}

func getSourceIDFromMethodOfDiscovery(methodOfDiscovery string, sourceIDToSource map[string]domain.Source) (sourceID string, err error) {
	if len(methodOfDiscovery) > 0 {
		for mapSourceID, source := range sourceIDToSource {
			if strings.ToLower(source.Source()) == strings.ToLower(methodOfDiscovery) {
				sourceID = mapSourceID
			}
		}

		if len(sourceID) == 0 {
			err = fmt.Errorf("could not find source for [%v]", methodOfDiscovery)
		}
	} else {
		err = fmt.Errorf("empty method of discovery")
	}

	return sourceID, err
}

func (job *CloudDecommissionJob) markDevicesAsDecommissionedInDatabase(deviceIDToDecommissionedDevice map[string]domain.Device) {
	for deviceID := range deviceIDToDecommissionedDevice {
		_, _, err := job.db.UpdateStateOfDevice(deviceID, domain.DeviceDecommed, job.config.OrganizationID())
		if err != nil {
			job.lstream.Send(log.Errorf(err, "error while updating status of %s", deviceID))
		}
	}
}

// TODO do we want to take MAC into consideration here?
// maybe region could be used as a fallback if MAC isn't present
func (job *CloudDecommissionJob) findDecommissionedDevices(historyOfDevices []domain.Device, allIPs []domain.CloudIP) map[string]domain.Device {
	var dbInstanceIDToDevice = make(map[string]domain.Device)
	var cloudInstanceIDToDevice = make(map[string]domain.CloudIP)

	// TODO use instance id

	// map devices that we've previously seen and are stored in the database
	for _, device := range historyOfDevices {

		// the region being nil means that it has not been picked up by the CSJ yet
		if len(sord(device.InstanceID())) > 0 {
			dbInstanceIDToDevice[sord(device.InstanceID())] = device
		}
	}

	// map the devices that are in the inventory of the cloud services
	for _, IP := range allIPs {
		if len(IP.InstanceID()) > 0 {
			cloudInstanceIDToDevice[IP.InstanceID()] = IP
		} else {
			job.lstream.Send(log.Warningf(nil, "cloud service did not return an instance ID for [%s]", IP.IP()))
		}
	}

	// find which devices we had stored in the databases that are not in the inventory of the cloud services (and are assumed to be decommissioned)
	// as well as devices that were reported by the cloud service as decommissioned
	var deviceIDToDecommDevice = make(map[string]domain.Device, 0)
	for _, device := range historyOfDevices {

		// the asset sync job (meaning a vulnerability scanner) has also found the device
		if device.SourceID() != nil {
			// the db device is not in the cloud inventory
			if cloudInstanceIDToDevice[sord(device.InstanceID())] == nil {
				deviceIDToDecommDevice[sord(device.SourceID())] = device
			} else {
				var matchedCloudDevice = cloudInstanceIDToDevice[sord(device.InstanceID())]

				// the cloud service reported the device as decommissioned
				if matchedCloudDevice.State() == domain.DeviceDecommed {
					deviceIDToDecommDevice[sord(device.SourceID())] = device
				} else if job.Payload.DecommOnStoppedState && matchedCloudDevice.State() == domain.DeviceStopped {
					deviceIDToDecommDevice[sord(device.SourceID())] = device
				}
			}
		}
	}

	return deviceIDToDecommDevice
}

func (job *CloudDecommissionJob) getSourceMap() (sourceIDToSource map[string]domain.Source, err error) {
	sourceIDToSource = make(map[string]domain.Source)
	var sources []domain.Source
	if sources, err = job.db.GetSources(); err == nil {
		for _, source := range sources {
			sourceIDToSource[source.ID()] = source
		}
	}

	return sourceIDToSource, err
}

func (job *CloudDecommissionJob) createIgnoreEntry(assetID string, scannerSourceID string, orgID string) {
	_, _, err := job.db.DeleteIgnoreForDevice(
		scannerSourceID,
		assetID,
		orgID,
	)

	if err == nil {

		_, _, err = job.db.SaveIgnore(
			scannerSourceID,
			orgID,
			domain.DecommAsset,
			"", // decommissioned assets don't require a vulnerability ID
			assetID,
			time.Now(),
			"",
			true,
			"",
		)

		if err == nil {
			job.lstream.Send(log.Infof("%s marked as decommissioned in the database", assetID))
		} else {
			job.lstream.Send(log.Errorf(err, "Error while updating exception with asset Id  %s", assetID))
		}
	} else {
		job.lstream.Send(log.Errorf(err, "error while deleting old entries in the ignore table for the device [%s]", assetID))
	}
}
