/*
  RETURN ScanSummary SINGLE
  SourceID          NVARCHAR(36)              NOT
  TemplateID        NVARCHAR(100)    NULL
  OrgID             NVARCHAR(36)              NOT
  SourceKey         NVARCHAR(100)    NULL
  ScanStatus        NVARCHAR(30)     NOT
  CreatedDate       DATETIME         NOT
  UpdatedDate       DATETIME         NULL
*/

DROP PROCEDURE IF EXISTS `GetScanSummary`;

CREATE PROCEDURE `GetScanSummary` (_SourceID VARCHAR(36), _OrgID VARCHAR(36), _ScanID NVARCHAR(100))
  #BEGIN#
  SELECT
    SS.SourceId,
    SS.TemplateId,
    SS.OrgId,
    SS.SourceKey,
    SS.ScanStatus,
    SS.CreatedDate,
    SS.UpdatedDate
  FROM ScanSummary SS
  WHERE SS.SourceId = _SourceID AND SS.OrgId = _OrgID AND SS.SourceKey = _ScanID
  LIMIT 1;
