DROP PROCEDURE IF EXISTS `UpdateJobConfigLastRun`;

CREATE PROCEDURE `UpdateJobConfigLastRun` (_ID NVARCHAR(36), _LastRun DATETIME)
  #BEGIN#
  BEGIN
    UPDATE JobConfig SET LastJobStart = _LastRun WHERE ID = _ID;
  END