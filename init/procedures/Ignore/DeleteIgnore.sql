DROP PROCEDURE IF EXISTS `DeleteIgnoreForDevice`;

CREATE PROCEDURE `DeleteIgnoreForDevice`(_sourceID VARCHAR(36), _devID VARCHAR(36), _orgID VARCHAR(36))
    #BEGIN#
    UPDATE `Ignore` SET Active = b'0' WHERE SourceID = _sourceID AND DeviceID = _devID AND OrganizationID = _orgID;