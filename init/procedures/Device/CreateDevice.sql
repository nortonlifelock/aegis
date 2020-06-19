DROP PROCEDURE IF EXISTS `CreateDevice`;

CREATE PROCEDURE `CreateDevice`(_AssetID NVARCHAR(36), _SourceID VARCHAR(36), _Ip NVARCHAR(32), _Hostname VARCHAR(1000), inInstanceID VARCHAR(200), _MAC VARCHAR(100), _GroupID VARCHAR(100), _OrgID NVARCHAR(36), _OS VARCHAR(300), _OSTypeID INT, inTrackingMethod VARCHAR(100))
  #BEGIN#
  INSERT INTO Device(AssetID, SourceID, IP, Hostname, InstanceId, MAC, GroupID, OrganizationID, OS, OSTypeID, TrackingMethod, IsVirtual) VALUE (_AssetID, _SourceID, _Ip, _Hostname, inInstanceID, _MAC, _GroupID, _OrgID, _OS, _OSTypeID, NULLIF(inTrackingMethod, ''), FALSE);