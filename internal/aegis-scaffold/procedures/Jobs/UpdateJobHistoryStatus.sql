DROP PROCEDURE IF EXISTS `UpdateJobHistoryStatus`;

CREATE PROCEDURE `UpdateJobHistoryStatus` (_ID NVARCHAR(36), _Status INT)
  #BEGIN#

  UPDATE JobHistory SET StatusID = _Status WHERE ID = _ID;