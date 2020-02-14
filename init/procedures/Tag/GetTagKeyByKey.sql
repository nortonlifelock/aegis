/*
  RETURN TagKey SINGLE
  ID                  NVARCHAR(36)             NOT
  KeyValue            NVARCHAR(255)   NOT
*/

DROP PROCEDURE IF EXISTS `GetTagKeyByKey`;

CREATE PROCEDURE `GetTagKeyByKey` (_KeyValue NVARCHAR(255))
  #BEGIN#
  SELECT
    T.Id,
    T.KeyValue
  FROM TagKey T
  WHERE T.KeyValue = _KeyValue;