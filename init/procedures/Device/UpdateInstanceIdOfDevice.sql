DROP PROCEDURE IF EXISTS `UpdateInstanceIDOfDevice`;

CREATE PROCEDURE `UpdateInstanceIDOfDevice` (_ID NVARCHAR(36), _InstanceID NVARCHAR(255), _CloudSourceID VARCHAR(36), _State VARCHAR(100), _Region VARCHAR(100), _OrgID NVARCHAR(36))
  #BEGIN#
  UPDATE Device D
    SET D.InstanceID = _InstanceID, D.CloudSourceID = _CloudSourceID, D.State = _State, D.Region = _Region
  WHERE D.ID = _ID AND D.OrganizationID = _OrgID;