DROP PROCEDURE IF EXISTS `SaveScanSummary`;

CREATE PROCEDURE `SaveScanSummary` (_ScanID NVARCHAR(100), _ScanStatus NVARCHAR(30))
  #BEGIN#

  UPDATE ScanSummary SET ScanStatus = _ScanStatus, UpdatedDate = NOW() WHERE SourceKey = _ScanID;