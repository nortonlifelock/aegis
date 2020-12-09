/*
  RETURN DeviceInfo SINGLE
  ID          VARCHAR(36)  NOT
  SourceID    VARCHAR(36)  NULL
  OS          VARCHAR(36)  NOT
  MAC         VARCHAR(36)  NOT
  IP          VARCHAR(36)  NOT
  HostName    VARCHAR(300) NOT
  Region      VARCHAR(100) NULL
  GroupID     VARCHAR(200) NULL
  InstanceID  VARCHAR(100) NULL
  TrackingMethod VARCHAR(100) NULL
*/

DROP PROCEDURE IF EXISTS `GetDeviceInfoByAssetIDNoOrg`;

CREATE PROCEDURE `GetDeviceInfoByAssetIDNoOrg`(inAssetID VARCHAR(360))
    #BEGIN#
SELECT
    D.ID,
    D.AssetID,
    D.OS,
    D.MAC,
    D.IP,
    D.HostName,
    D.Region,
    D.GroupId,
    D.InstanceID,
    D.TrackingMethod
FROM Device D
WHERE D.AssetID = inAssetID;