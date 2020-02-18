DROP PROCEDURE IF EXISTS `CreateScanSummary`;

CREATE PROCEDURE `CreateScanSummary` (_SourceID NVARCHAR(36), _ScannerSourceConfigID VARCHAR(36), _OrgID NVARCHAR(36), _ScanID NVARCHAR(100), _ScanStatus NVARCHAR(30), _ScanClosePayload MEDIUMTEXT, _ParentJobID NVARCHAR(36))
    #BEGIN#

INSERT INTO ScanSummary (SourceID, ScannerSourceConfigID, OrgID, SourceKey, ScanStatus, ParentJobID, ScanClosePayload)
    VALUE (_SourceID, _ScannerSourceConfigID, _OrgID, _ScanID, _ScanStatus, _ParentJobID, _ScanClosePayload)
ON DUPLICATE KEY UPDATE ScanStatus = _ScanStatus, UpdatedDate = NOW();

