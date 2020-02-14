DROP PROCEDURE IF EXISTS `CreateOrganization`;

CREATE PROCEDURE `CreateOrganization` (_Code NVARCHAR(150), _Description NVARCHAR(500), _TimeZoneOffset FLOAT, _UpdatedBy TEXT)
  #BEGIN#
  INSERT INTO Organization (Code, Description, TimeZoneOffset, UpdatedBy)
    VALUE (_Code, _Description, _TimeZoneOffset, _UpdatedBy);