/*
  RETURN ScanSummary
  SourceID          NVARCHAR(36)              NOT
  TemplateID        NVARCHAR(100)    NULL
  OrgID             NVARCHAR(36)              NOT
  SourceKey         NVARCHAR(100)    NULL
  ScanStatus        NVARCHAR(30)     NOT
  ParentJobID       NVARCHAR(36)              NOT
  ScanClosePayload  TEXT             NOT
  CreatedDate       DATETIME         NOT
  UpdatedDate       DATETIME         NULL
*/

DROP PROCEDURE IF EXISTS `GetScanSummariesBySourceName`;

CREATE PROCEDURE `GetScanSummariesBySourceName` (_OrgID VARCHAR(36), _SourceName NVARCHAR(30))
  #BEGIN#
  SELECT
    SS.SourceId,
    SS.TemplateId,
    SS.OrgId,
    SS.SourceKey,
    SS.ScanStatus,
    SS.ParentJobId,
    SS.ScanClosePayload,
    SS.CreatedDate,
    SS.UpdatedDate
  FROM ScanSummary SS
    JOIN Source S ON SS.SourceId = S.Id
  WHERE SS.OrgId = _OrgID AND S.Source = _SourceName
  ORDER BY CreatedDate DESC;