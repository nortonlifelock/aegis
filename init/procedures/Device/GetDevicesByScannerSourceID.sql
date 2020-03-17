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
  ScannerSourceID VARCHAR(36)  NULL
*/

DROP PROCEDURE IF EXISTS `GetDeviceInfoByScannerSourceID`;

CREATE PROCEDURE `GetDeviceInfoByScannerSourceID` (_IP VARCHAR(36), _GroupID INT, _OrgID VARCHAR(36))
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
         JOIN AssetGroup AG on _GroupID = AG.GroupID AND D.CloudSourceID = AG.CloudSourceID
WHERE D.OrganizationID = _OrgID AND D.IP = _IP;