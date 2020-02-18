/*
  RETURN ExceptionType
  ID              INT                    NOT
  Type            NVARCHAR(25)           NOT
  Name            NVARCHAR(50)           NOT
*/

DROP PROCEDURE IF EXISTS `GetExceptionTypes`;

CREATE PROCEDURE `GetExceptionTypes` ()
#BEGIN#
SELECT
    It.Id,
    It.Type,
    It.Name
FROM IgnoreType It;