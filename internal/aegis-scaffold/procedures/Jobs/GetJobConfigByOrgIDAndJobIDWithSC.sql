/*
  RETURN JobConfig
  ID                      NVARCHAR(36)           NOT
  JobID                   INT           NOT
  OrganizationID          NVARCHAR(36)           NOT
  DataInSourceConfigID    NVARCHAR(36)  NULL
  DataOutSourceConfigID   NVARCHAR(36)  NULL
  PriorityOverride        INT           NULL
  Continuous              BIT           NOT
  WaitInSeconds           INT           NOT
  MaxInstances            INT           NOT
  AutoStart               BIT           NOT
  CreatedDate             DATETIME      NOT
  CreatedBy               NVARCHAR(255) NOT
  UpdatedDate             DATETIME      NULL
  UpdatedBy               NVARCHAR(255) NULL
  LastJobStart            DATETIME      NULL
*/

DROP PROCEDURE IF EXISTS `GetJobConfigByOrgIDAndJobIDWithSC`;

CREATE PROCEDURE `GetJobConfigByOrgIDAndJobIDWithSC` (_OrgID VARCHAR(36), _JobID INT, _SourceConfigID VARCHAR(36))
    #BEGIN#
SELECT
    JC.Id,
    JC.JobId,
    JC.OrganizationId,
    JC.DataInSourceConfigId,
    JC.DataOutSourceConfigId,
    JC.PriorityOverride,
    JC.Continuous,
    JC.WaitInSeconds,
    JC.MaxInstances,
    JC.AutoStart,
    JC.CreatedDate,
    JC.CreatedBy,
    JC.UpdatedDate,
    JC.UpdatedBy,
    JC.LastJobStart
FROM JobConfig JC
WHERE JC.OrganizationId = _OrgID
  AND JC.JobId = _JobID AND (JC.DataInSourceConfigID = _SourceConfigID OR JC.DataOutSourceConfigID = _SourceConfigID);