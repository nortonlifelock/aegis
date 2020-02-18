DROP PROCEDURE IF EXISTS `UpdateSourceConfigToken`;

CREATE PROCEDURE `UpdateSourceConfigToken` (_ID NVARCHAR(36), _Token TEXT)
  #BEGIN#
#   UPDATE SourceConfig SET Token = _Token WHERE ID = _ID;
    UPDATE SourceConfig SET AuthInfo = JSON_SET(AuthInfo, "$.Token", _Token) Where ID = _ID;
