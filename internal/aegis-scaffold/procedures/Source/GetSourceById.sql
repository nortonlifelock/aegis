/*
  RETURN Source SINGLE
  ID                      VARCHAR(36)       NOT
  SourceTypeID            INT               NOT
  Source                  NVARCHAR(100)     NOT
*/

DROP PROCEDURE IF EXISTS `GetSourceByID`;

CREATE PROCEDURE `GetSourceByID` (_ID VARCHAR(36))
  #BEGIN#
  SELECT
    S.Id,
    S.SourceTypeId,
    S.Source

  From Source S
  WHERE S.Id = _ID;