DROP PROCEDURE IF EXISTS `CreateOrganizationWithPayloadEkey`;

CREATE PROCEDURE `CreateOrganizationWithPayloadEkey` (_Code NVARCHAR(150), _Description NVARCHAR(500), _TimeZoneOffset FLOAT, _Payload TEXT, _EKEY TEXT, _UpdatedBy TEXT)
    #BEGIN#
INSERT INTO Organization (Code, Description, TimeZoneOffset, Payload, EncryptionKey, UpdatedBy)
    VALUE (_Code, _Description, _TimeZoneOffset, _Payload, _EKEY, _UpdatedBy);