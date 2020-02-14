DROP PROCEDURE IF EXISTS `UpdateStateOfDevice`;

CREATE PROCEDURE `UpdateStateOfDevice` (_ID NVARCHAR(36), _State NVARCHAR(100), _OrgID NVARCHAR(36))
    #BEGIN#
UPDATE Device D
SET D.State = _State
WHERE D.ID = _ID AND D.OrganizationID = _OrgID;