/*
  RETURN TagKey SINGLE
  ID                  NVARCHAR(36)             NOT
  KeyValue            NVARCHAR(255)   NOT
*/

DROP PROCEDURE IF EXISTS `GetTagKeyByID`;

CREATE PROCEDURE `GetTagKeyByID` (_ID NVARCHAR(36))
  #BEGIN#
  SELECT
    T.Id,
    T.KeyValue
  FROM TagKey T
  WHERE T.Id = _ID;