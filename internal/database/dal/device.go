package dal

import (
	"context"
	"fmt"
	"github.com/nortonlifelock/domain"
)

// Device implements the domain.Device interface and utilizes the DeviceInfo object (which is tied to the database table Device)
type Device struct {
	Conn domain.DatabaseConnection
	Info domain.DeviceInfo
}

// SourceID reports the ID of the device as according to the vulnerability scanner that found it
func (device *Device) SourceID() *string {
	return device.Info.SourceID()
}

// OS returns the operating system that is running on the device
func (device *Device) OS() string {
	return device.Info.OS()
}

// HostName returns the hostname of the device
func (device *Device) HostName() string {
	return device.Info.HostName()
}

// MAC returns the MAC address of the device
func (device *Device) MAC() string {
	return device.Info.MAC()
}

// IP returns the IP address of the device
func (device *Device) IP() string {
	return device.Info.IP()
}

// ID returns the ID of the device as according to the Aegis database
func (device *Device) ID() string {
	return device.Info.ID()
}

// Region returns the area in which the device exists (if the device exists in a cloud source like AWS/Azure)
func (device *Device) Region() *string {
	return device.Info.Region()
}

// InstanceID returns the ID of the device as according to the cloud environment that the device exists in
func (device *Device) InstanceID() *string {
	return device.Info.InstanceID()
}

func (device *Device) GroupID() *string {
	var groupID string
	if device.Info.GroupID() != nil {
		groupID = *device.Info.GroupID()
	}

	return &groupID
}

// Vulnerabilities returns a channel that contains all the detections that were found on the device
func (device *Device) Vulnerabilities(ctx context.Context) (param <-chan domain.Detection, err error) {
	out := make(chan domain.Detection)

	var id = device.SourceID()
	if id != nil && len(*id) > 0 {
		go func() {
			defer close(out)

			detections, err := device.Conn.GetDetectionsForDevice(*id)
			if err == nil {
				for _, detection := range detections {
					select {
					case <-ctx.Done():
						return
					case out <- detection:
					}
				}
			}
		}()
	} else {
		err = fmt.Errorf("empty device id")
	}

	return out, err
}
