DROP PROCEDURE IF EXISTS `CreateLog`;

  CREATE PROCEDURE `CreateLog` (_TypeID INT, _Log TEXT, _Error TEXT, _JobHistoryID VARCHAR(100), _CreateDate DATETIME)
    #BEGIN#

  INSERT INTO Logs (TypeID, Log, Error, JobHistoryID, CreateDate)
  VALUES(_TypeID, _Log, _Error, _JobHistoryID, _CreateDate);