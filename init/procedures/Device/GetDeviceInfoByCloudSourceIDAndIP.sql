/*
  RETURN DeviceInfo SINGLE
  ID              VARCHAR(36)  NOT
  SourceID        VARCHAR(36)  NULL
  OS              VARCHAR(36)  NOT
  MAC             VARCHAR(36)  NOT
  IP              VARCHAR(36)  NOT
  HostName        VARCHAR(300) NOT
  State           VARCHAR(100) NULL
  Region          VARCHAR(100) NULL
  InstanceID      VARCHAR(100) NULL
  ScannerSourceID VARCHAR(36)  NOT
*/

DROP PROCEDURE IF EXISTS `GetDeviceInfoByCloudSourceIDAndIP`;

CREATE PROCEDURE `GetDeviceInfoByCloudSourceIDAndIP` (_IP VARCHAR(36), _CloudSourceID VARCHAR(100), _OrgID VARCHAR(36))
    #BEGIN#
SELECT
    D.ID,
    D.AssetID,
    D.OS,
    D.MAC,
    D.IP,
    D.HostName,
    D.State,
    D.Region,
    D.InstanceID,
    D.SourceID
FROM Device D
         JOIN AssetGroup AG on D.GroupID = AG.GroupID AND D.SourceID = AG.ScannerSourceID
WHERE D.OrganizationID = _OrgID AND D.Ip = _IP AND AG.CloudSourceID = _CloudSourceID;