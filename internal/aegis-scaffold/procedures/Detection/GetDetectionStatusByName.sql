/*
  RETURN DetectionStatus SINGLE
  ID                  INT             NOT
  Status              NVARCHAR(128)   NOT
  Name                NVARCHAR(128)   NOT
*/

DROP PROCEDURE IF EXISTS `GetDetectionStatusByName`;

CREATE PROCEDURE `GetDetectionStatusByName` (_Name NVARCHAR(128))
  #BEGIN#
  SELECT
    D.Id,
    D.Status,
    D.Name
  FROM DetectionStatus D
  WHERE D.Name = _Name;