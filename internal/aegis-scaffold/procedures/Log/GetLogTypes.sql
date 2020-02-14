/*
  RETURN LogType
  ID                INT                     NOT
  LogType           NVARCHAR(15)            NOT
  Name              NVARCHAR(50)            NOT
*/

DROP PROCEDURE IF EXISTS `GetLogTypes`;

CREATE PROCEDURE `GetLogTypes` ()
  #BEGIN#
  SELECT
    LT.Id,
    LT.Type,
    LT.Name
  FROM LogType LT;