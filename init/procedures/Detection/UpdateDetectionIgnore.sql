DROP PROCEDURE IF EXISTS `UpdateDetectionIgnore`;

CREATE PROCEDURE `UpdateDetectionIgnore` (_DeviceID NVARCHAR(36), _VulnID NVARCHAR(36), _ExceptionID VARCHAR(36))
    #BEGIN#
UPDATE Detection D
SET IgnoreID = NULLIF(_ExceptionID, '')
WHERE D.DeviceID = _DeviceID AND D.VulnerabilityID = _VulnID;