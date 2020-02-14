DROP PROCEDURE IF EXISTS `UpdateAssetIDOsTypeIDOfDevice`;

CREATE PROCEDURE `UpdateAssetIDOsTypeIDOfDevice` (_ID NVARCHAR(36), _AssetID NVARCHAR(36), _ScannerSourceID VARCHAR(36), _GroupID INT, _OS VARCHAR(300), _HostName VARCHAR(1000), _OsTypeID INT, _OrgID NVARCHAR(36))
  #BEGIN#
  UPDATE Device D
  SET D.AssetID = _AssetID, D.OSTypeID = _OsTypeID, D.OS = _OS, D.GroupId = _GroupID, D.SourceID = _ScannerSourceID, D.HostName = _HostName
  WHERE D.ID = _ID AND D.OrganizationID = _OrgID;