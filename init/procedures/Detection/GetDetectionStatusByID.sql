/*
  RETURN DetectionStatus SINGLE
  ID                  INT             NOT
  Status              NVARCHAR(128)   NOT
  Name                NVARCHAR(128)   NOT
*/

DROP PROCEDURE IF EXISTS `GetDetectionStatusByID`;

CREATE PROCEDURE `GetDetectionStatusByID` (_ID INT)
    #BEGIN#
SELECT
    D.Id,
    D.Status,
    D.Name
FROM DetectionStatus D
WHERE D.ID = _ID;