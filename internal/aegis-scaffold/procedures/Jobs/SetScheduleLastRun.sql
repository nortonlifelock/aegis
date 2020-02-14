DROP PROCEDURE IF EXISTS `SetScheduleLastRun`;

CREATE PROCEDURE `SetScheduleLastRun` (_ID NVARCHAR(36))
  #BEGIN#
  UPDATE JobSchedule SET LastRun = NOW() WHERE Id = _ID;