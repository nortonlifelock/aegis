DROP PROCEDURE IF EXISTS `CreateDevice`;

CREATE PROCEDURE `CreateDevice`(_AssetID NVARCHAR(36), _SourceID VARCHAR(36), _Ip NVARCHAR(32), _Hostname VARCHAR(1000), _MAC VARCHAR(100), _GroupID INT, _OrgID NVARCHAR(36), _OS VARCHAR(300), _OSTypeID INT)
  #BEGIN#
  INSERT INTO Device(AssetID, SourceID, InstanceID, IP, Hostname, MAC, GroupID, OrganizationID, OS, OSTypeID, IsVirtual) VALUE (_AssetID, _SourceID, '', _Ip, _Hostname, _MAC, _GroupID, _OrgID, _OS, _OSTypeID, FALSE);