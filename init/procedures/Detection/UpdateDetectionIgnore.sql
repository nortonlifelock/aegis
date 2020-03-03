DROP PROCEDURE IF EXISTS `UpdateDetectionIgnore`;

CREATE PROCEDURE `UpdateDetectionIgnore` (_DeviceID NVARCHAR(36), _VulnID NVARCHAR(36), _Port INT, _Protocol VARCHAR(36), _ExceptionID VARCHAR(36))
    #BEGIN#
UPDATE Detection D JOIN VulnerabilityInfo VI on (D.VulnerabilityID = VI.ID)
SET IgnoreID = NULLIF(_ExceptionID, '')
WHERE D.DeviceID = _DeviceID AND VI.SourceVulnId = _VulnID AND D.Port = _Port AND D.Protocol = _Protocol;