DROP PROCEDURE IF EXISTS `CreateAssetWithIPInstanceID`;

CREATE PROCEDURE `CreateAssetWithIPInstanceID`(_State VARCHAR(100), _IP NVARCHAR(16), _MAC VARCHAR(100), _SourceID VARCHAR(100), _InstanceID NVARCHAR(128), _Region VARCHAR(100), _OrgID NVARCHAR(36), _OS VARCHAR(300), _OsTypeID INT)
  #BEGIN#
  INSERT INTO Device(State, IP, MAC, CloudSourceID, InstanceID, HostName, Region, OrganizationID, OS, OSTypeID, IsVirtual) VALUE (_State, _IP, _MAC, _SourceID, _InstanceID, '', _Region, _OrgID, _OS, _OsTypeID, FALSE);