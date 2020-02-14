/*
  RETURN ScanSummary
  SourceID          NVARCHAR(36)              NOT
  TemplateID        NVARCHAR(100)    NULL
  OrgID             NVARCHAR(36)              NOT
  SourceKey         NVARCHAR(100)    NULL
  ScanStatus        NVARCHAR(30)     NOT
  ScanClosePayload  TEXT             NOT
  ParentJobID       NVARCHAR(36)              NOT
  CreatedDate       DATETIME         NOT
  UpdatedDate       DATETIME         NULL
  Source            NVARCHAR(30)     NOT
*/

DROP PROCEDURE IF EXISTS `GetRecentlyUpdatedScanSummaries`;

CREATE PROCEDURE `GetRecentlyUpdatedScanSummaries` (_OrgID VARCHAR(36))
  #BEGIN#
  SELECT
    SS.SourceId,
    SS.TemplateId,
    SS.OrgId,
    SS.SourceKey,
    SS.ScanStatus,
    SS.ScanClosePayload,
    SS.ParentJobId,
    SS.CreatedDate,
    SS.UpdatedDate,
    S.Source
  FROM ScanSummary SS
    JOIN Source S ON S.Id = SS.SourceId
    WHERE SS.OrgId = _OrgID AND (SS.UpdatedDate >= (NOW() - INTERVAL 5 MINUTE) OR SS.UpdatedDate IS NULL);