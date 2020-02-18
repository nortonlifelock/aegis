/*
  RETURN DetectionStatus
  ID                  INT             NOT
  Status              NVARCHAR(128)   NOT
  Name                NVARCHAR(128)   NOT
*/

DROP PROCEDURE IF EXISTS `GetDetectionStatuses`;

CREATE PROCEDURE `GetDetectionStatuses` ()
    #BEGIN#
SELECT
    D.Id,
    D.Status,
    D.Name
FROM DetectionStatus D;