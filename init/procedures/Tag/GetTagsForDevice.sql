/*
  RETURN Tag
  ID                  NVARCHAR(36)             NOT
  DeviceID            NVARCHAR(36)             NOT
  TagKeyID            INT             NOT
  Value               NVARCHAR(255)   NOT
*/

DROP PROCEDURE IF EXISTS `GetTagsForDevice`;

CREATE PROCEDURE `GetTagsForDevice` (_DeviceID NVARCHAR(36))
  #BEGIN#
  SELECT
    T.Id,
    T.DeviceId,
    T.TagKeyId,
    T.Value
  FROM Tag T
  WHERE T.DeviceId = _DeviceID;