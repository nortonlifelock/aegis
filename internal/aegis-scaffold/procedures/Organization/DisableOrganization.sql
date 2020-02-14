DROP PROCEDURE IF EXISTS `DisableOrganization`;

CREATE PROCEDURE `DisableOrganization` (_ID NVARCHAR(36), _UpdatedBy TEXT)
  #BEGIN#
  BEGIN
    UPDATE Organization SET Active = b'0' WHERE ID = _ID;
    UPDATE Organization SET Updated = NOW() WHERE ID = _ID;
    UPDATE Organization SET UpdatedBy = _UpdatedBy WHERE ID = _ID;
  END;