DROP PROCEDURE IF EXISTS `DisableJobConfig`;

CREATE PROCEDURE `DisableJobConfig` (_ID NVARCHAR(36), _UpdatedBy TEXT)
  #BEGIN#
  BEGIN
    UPDATE JobConfig SET Active = b'0', UpdatedDate = NOW(), UpdatedBy = _UpdatedBy WHERE ID = _ID;
  END