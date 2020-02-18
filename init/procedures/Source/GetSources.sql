/*
  RETURN Source
  ID                VARCHAR(36)             NOT
  SourceTypeID      INT                     NOT
  Source            NVARCHAR(100)           NOT
  DBCreatedDate     DATETIME                NOT
  DBUpdatedDate     DATETIME                NULL
*/

DROP PROCEDURE IF EXISTS `GetSources`;

CREATE PROCEDURE `GetSources` ()
#BEGIN#
SELECT
    S.ID,
    S.SourceTypeId,
    S.Source,
    S.Created,
    S.Updated
FROM Source S;