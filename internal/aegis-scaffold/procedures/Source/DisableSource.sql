DROP PROCEDURE IF EXISTS `DisableSource`;

CREATE PROCEDURE `DisableSource` (_ID VARCHAR(36), _OrgID NVARCHAR(36), _UpdatedBy TEXT)
  #BEGIN#
  BEGIN
    UPDATE SourceConfig SC SET
        SC.Active = b'0',
        SC.DBUpdatedDate = NOW(),
        SC.UpdatedBy = _UpdatedBy
    WHERE SC.ID = _ID AND OrganizationID = _OrgID;
  END