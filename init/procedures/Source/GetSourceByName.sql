/*
  RETURN Source SINGLE
  ID                VARCHAR(36)             NOT
  SourceTypeID      INT                     NOT
  Source            NVARCHAR(100)           NOT
  DBCreatedDate     DATETIME                NOT
  DBUpdatedDate     DATETIME                NULL
*/

DROP PROCEDURE IF EXISTS `GetSourceByName`;

CREATE PROCEDURE `GetSourceByName` (_Source NVARCHAR(100))
#BEGIN#
  SELECT
    S.Id,
    S.SourceTypeId,
    S.Source,
    S.Created,
    S.Updated
  FROM Source S
  WHERE S.Source = _Source;