DROP PROCEDURE IF EXISTS `PulseJob`;

CREATE PROCEDURE `PulseJob` (_JobHistoryID NVARCHAR(36))
  #BEGIN#
BEGIN
  SET @date = NOW();

  UPDATE JobHistory
  SET
    PulseDate = @date
  WHERE Id = _JobHistoryID;
  END;