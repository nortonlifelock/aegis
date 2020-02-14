DROP PROCEDURE IF EXISTS `UpdateSourceConfigConcurrencyByID`;

CREATE PROCEDURE `UpdateSourceConfigConcurrencyByID` (_ID NVARCHAR(36), _Delay INT, _Retries INT, _Concurrency INT)
    #BEGIN#
    BEGIN
        UPDATE SourceConfig SET AuthInfo = JSON_SET(AuthInfo, "$.Delay", _Delay) Where ID = _ID;
        UPDATE SourceConfig SET AuthInfo = JSON_SET(AuthInfo, "$.Retries", _Retries) Where ID = _ID;
        UPDATE SourceConfig SET AuthInfo = JSON_SET(AuthInfo, "$.Concurrency", _Concurrency) Where ID = _ID;
    END;

