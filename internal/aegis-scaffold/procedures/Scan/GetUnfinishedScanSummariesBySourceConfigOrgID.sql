/*
  RETURN ScanSummary
  SourceID          NVARCHAR(36)              NOT
  TemplateID        NVARCHAR(100)    NULL
  OrgID             NVARCHAR(36)              NOT
  SourceKey         NVARCHAR(100)    NULL
  ScanStatus        NVARCHAR(30)     NOT
  ScanClosePayload  TEXT             NOT
  ParentJobID       VARCHAR(36)              NOT
  CreatedDate       DATETIME         NOT
  UpdatedDate       DATETIME         NULL
*/

DROP PROCEDURE IF EXISTS `GetUnfinishedScanSummariesBySourceConfigOrgID`;

CREATE PROCEDURE `GetUnfinishedScanSummariesBySourceConfigOrgID` (_ScannerSourceConfigID NVARCHAR(36), _OrgID NVARCHAR(36))
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
    SS.UpdatedDate
FROM ScanSummary SS
WHERE SS.ScanStatus NOT IN ('finished','canceled','stopped','error')
  AND _ScannerSourceConfigID = SS.ScannerSourceConfigID AND _OrgID = SS.OrgId;