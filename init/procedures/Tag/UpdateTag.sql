DROP PROCEDURE IF EXISTS `UpdateTag`;

CREATE PROCEDURE `UpdateTag` (_DeviceID NVARCHAR(36), _TagKeyID VARCHAR(36), _Value NVARCHAR(255))
  #BEGIN#
  UPDATE Tag T
    SET T.Value = _Value
  WHERE T.DeviceID = _DeviceID AND T.TagKeyID = _TagKeyID;