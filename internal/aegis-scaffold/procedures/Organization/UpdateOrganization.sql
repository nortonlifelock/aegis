DROP PROCEDURE IF EXISTS `UpdateOrganization`;

CREATE PROCEDURE `UpdateOrganization` (_ID NVARCHAR(36), _Description NVARCHAR(500), _TimezoneOffset FLOAT, _UpdatedBy TEXT)
  #BEGIN#
  BEGIN
    UPDATE Organization SET Description = _Description, TimezoneOffset = _TimezoneOffset,
      UpdatedBy = _UpdatedBy WHERE ID = _ID;
  END