DROP PROCEDURE IF EXISTS `CreateTag`;

CREATE PROCEDURE `CreateTag`(_DeviceID NVARCHAR(36), _TagKeyID VARCHAR(36), _Value NVARCHAR(255))

  #BEGIN#
  INSERT INTO Tag (DeviceID, TagKeyID, Value) VALUES (_DeviceID, _TagKeyID, _Value);