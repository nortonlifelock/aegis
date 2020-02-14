/*
  RETURN SourceConfig
  ID                VARCHAR(36)             NOT
  Address           NVARCHAR(100)           NOT
  Source            NVARCHAR(30)            NOT
  Port              NVARCHAR(30)            NOT

*/

DROP PROCEDURE IF EXISTS `GetSourceInsByJobID`;

CREATE PROCEDURE `GetSourceInsByJobID` (inJob INT, inOrgID VARCHAR(36))
    #BEGIN#
SELECT
    SC.Id,
    SC.Address,
    SC.Source,
    SC.Port
FROM Job J
     JOIN Source S on S.SourceTypeId = J.SourceTypeIn
     JOIN SourceConfig SC on SC.SourceId = S.Id
WHERE J.Id = inJob and SC.OrganizationId = inOrgID;