DROP PROCEDURE IF EXISTS `UpdateJobHistoryStatusDetailed`;

CREATE PROCEDURE `UpdateJobHistoryStatusDetailed` (_ID NVARCHAR(36), _Status INT, _UpdatedBy TEXT)
  #BEGIN#
BEGIN
  UPDATE JobHistory SET UpdatedBy = _UpdatedBy WHERE ID = _ID;
  UPDATE JobHistory SET StatusID = _Status WHERE ID = _ID;
  UPDATE JobHistory SET UpdatedDate = NOW() WHERE ID = _ID;
END