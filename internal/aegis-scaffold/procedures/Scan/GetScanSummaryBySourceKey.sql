/*
  RETURN ScanSummary SINGLE
  SourceID          NVARCHAR(36)     NOT
  TemplateID        NVARCHAR(100)    NULL
  OrgID             NVARCHAR(36)     NOT
  SourceKey         NVARCHAR(100)    NULL
  ScanStatus        NVARCHAR(30)     NOT
  ScanClosePayload  TEXT             NOT
  CreatedDate       DATETIME         NOT
  UpdatedDate       DATETIME         NULL
*/

DROP PROCEDURE IF EXISTS `GetScanSummaryBySourceKey`;

CREATE PROCEDURE `GetScanSummaryBySourceKey` (_SourceKey NVARCHAR(100))
  #BEGIN#
  SELECT
    SS.SourceId,
    SS.TemplateId,
    SS.OrgId,
    SS.SourceKey,
    SS.ScanStatus,
    SS.ScanClosePayload,
    SS.CreatedDate,
    SS.UpdatedDate
  FROM ScanSummary SS
  WHERE SS.SourceKey = _SourceKey;
