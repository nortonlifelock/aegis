package database

import (
	"github.com/nortonlifelock/aegis/internal/database/dal"
	"github.com/nortonlifelock/domain"
	"time"
)

// procedures that are added here must be added to the DatabaseConnection interface in domain/interface.go

func (conn *dbconn) GetDeviceByAssetOrgID(_AssetID string, OrgID string) (domain.Device, error) {
	var device domain.Device
	info, err := conn.GetDeviceInfoByAssetOrgID(_AssetID, OrgID)
	if err == nil {
		if info != nil {
			device = &dal.Device{
				Conn: conn,
				Info: info,
			}
		}
	}

	return device, err
}

func (conn *dbconn) GetDeviceByIP(_IP string, _OrgID string) (domain.Device, error) {
	var device domain.Device
	info, err := conn.GetDeviceInfoByIP(_IP, _OrgID)
	if err == nil {
		if info != nil {
			device = &dal.Device{
				Conn: conn,
				Info: info,
			}
		}
	}

	return device, err
}

func (conn *dbconn) GetDeviceByInstanceID(_InstanceID string, _OrgID string) ([]domain.Device, error) {
	var devices = make([]domain.Device, 0)

	var infos []domain.DeviceInfo
	var err error
	infos, err = conn.GetDeviceInfoByInstanceID(_InstanceID, _OrgID)
	if err == nil {
		for _, info := range infos {
			if info != nil {
				devices = append(devices, &dal.Device{
					Conn: conn,
					Info: info,
				})
			}
		}
	}

	return devices, err
}

func (conn *dbconn) GetDeviceByScannerSourceID(_IP string, _GroupID string, _OrgID string) (domain.Device, error) {
	var device domain.Device
	info, err := conn.GetDeviceInfoByScannerSourceID(_IP, _GroupID, _OrgID)
	if err == nil {
		if info != nil {
			device = &dal.Device{
				Conn: conn,
				Info: info,
			}
		}
	}

	return device, err
}

func (conn *dbconn) GetDeviceByCloudSourceIDAndIP(_IP string, _CloudSourceID string, _OrgID string) ([]domain.Device, error) {
	var devices = make([]domain.Device, 0)

	var infos []domain.DeviceInfo
	var err error
	infos, err = conn.GetDeviceInfoByCloudSourceIDAndIP(_IP, _CloudSourceID, _OrgID)
	if err == nil {
		for _, info := range infos {
			if info != nil {
				devices = append(devices, &dal.Device{
					Conn: conn,
					Info: info,
				})
			}
		}
	}

	return devices, err
}

func (conn *dbconn) GetDevicesBySourceID(_SourceID string, _OrgID string) ([]domain.Device, error) {
	var devices = make([]domain.Device, 0)

	var infos []domain.DeviceInfo
	var err error
	infos, err = conn.GetDevicesInfoBySourceID(_SourceID, _OrgID)
	if err == nil {
		for _, info := range infos {
			if info != nil {
				devices = append(devices, &dal.Device{
					Conn: conn,
					Info: info,
				})
			}
		}
	}

	return devices, err
}

func (conn *dbconn) GetDevicesByCloudSourceID(_CloudSourceID string, _OrgID string) ([]domain.Device, error) {
	var devices = make([]domain.Device, 0)

	var infos []domain.DeviceInfo
	var err error
	infos, err = conn.GetDevicesInfoByCloudSourceID(_CloudSourceID, _OrgID)
	if err == nil {
		for _, info := range infos {
			if info != nil {
				devices = append(devices, &dal.Device{
					Conn: conn,
					Info: info,
				})
			}
		}
	}

	return devices, err
}

func (conn *dbconn) GetDetectionForGroupAfter(_After time.Time, _OrgID string, inGroupID string) ([]domain.Detection, error) {
	var detections = make([]domain.Detection, 0)

	var infos []domain.DetectionInfo
	var err error
	infos, err = conn.GetDetectionInfoForGroupAfter(_After, _OrgID, inGroupID)
	if err == nil {
		for _, info := range infos {
			if info != nil {
				detections = append(detections, &dal.Detection{
					Conn: conn,
					Info: info,
				})
			}
		}
	}

	return detections, err
}

func (conn *dbconn) GetDetection(_DeviceID string, _VulnerabilityID string, _Port int, _Protocol string) (domain.Detection, error) {
	var detection domain.Detection
	info, err := conn.GetDetectionInfo(_DeviceID, _VulnerabilityID, _Port, _Protocol)
	if err == nil {
		if info != nil {
			detection = &dal.Detection{
				Conn: conn,
				Info: info,
			}
		}
	}

	return detection, err
}

func (conn *dbconn) GetDetectionBySourceVulnID(_SourceDeviceID string, _SourceVulnerabilityID string, _Port int, _Protocol string) (domain.Detection, error) {
	var detection domain.Detection
	info, err := conn.GetDetectionInfoBySourceVulnID(_SourceDeviceID, _SourceVulnerabilityID, _Port, _Protocol)
	if err == nil {
		if info != nil {
			detection = &dal.Detection{
				Conn: conn,
				Info: info,
			}
		}
	}

	return detection, err
}

func (conn *dbconn) GetDetectionsForDevice(_DeviceID string) ([]domain.Detection, error) {
	var detections = make([]domain.Detection, 0)

	var infos []domain.DetectionInfo
	var err error
	infos, err = conn.GetDetectionsInfoForDevice(_DeviceID)
	if err == nil {
		for _, info := range infos {
			if info != nil {
				detections = append(detections, &dal.Detection{
					Conn: conn,
					Info: info,
				})
			}
		}
	}

	return detections, err
}

func (conn *dbconn) GetDetectionsAfter(after time.Time, orgID string) (detections []domain.Detection, err error) {
	detections = make([]domain.Detection, 0)

	var infos []domain.DetectionInfo
	infos, err = conn.GetDetectionInfoAfter(after, orgID)
	for _, info := range infos {
		if info != nil {
			detections = append(detections, &dal.Detection{
				Conn: conn,
				Info: info,
			})
		}
	}

	return detections, err
}

func (conn *dbconn) GetVulnReferences(vulnInfoID string, sourceID string) (references []domain.VulnerabilityReference, err error) {
	references = make([]domain.VulnerabilityReference, 0)

	var infos []domain.VulnerabilityReferenceInfo
	if infos, err = conn.GetVulnReferencesInfo(vulnInfoID, sourceID); err == nil {
		for _, info := range infos {
			if info != nil {
				references = append(references, &dal.VulnerablityReference{
					Conn: conn,
					Info: info,
				})
			}
		}
	}

	return references, err
}

func (conn *dbconn) GetVulnRef(vulnInfoID string, sourceID string, reference string) (existing domain.VulnerabilityReference, err error) {
	var info domain.VulnerabilityReferenceInfo
	info, err = conn.GetVulnRefInfo(vulnInfoID, sourceID, reference)
	if err == nil {
		if info != nil {
			existing = &dal.VulnerablityReference{
				Conn: conn,
				Info: info,
			}
		}
	}

	return existing, err
}

func (conn *dbconn) GetVulnBySourceVulnID(_SourceVulnID string) (vulnerability domain.Vulnerability, err error) {
	var info domain.VulnerabilityInfo
	if info, err = conn.GetVulnInfoBySourceVulnID(_SourceVulnID); err == nil {
		if info != nil {
			vulnerability = &dal.Vulnerability{
				Conn: conn,
				Info: info,
			}
		}
	}
	return vulnerability, err
}
