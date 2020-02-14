/*
  RETURN JobConfigAudit
  ID                      NVARCHAR(36)  NOT
  JobID                   INT           NOT
  OrganizationID          NVARCHAR(36)  NOT
  DataInSourceConfigID    NVARCHAR(36)  NULL
  DataOutSourceConfigID   NVARCHAR(36)  NULL
  Payload                 TEXT          NULL
  PriorityOverride        INT           NULL
  Continuous              BIT           NOT
  WaitInSeconds           INT           NOT
  MaxInstances            INT           NOT
  AutoStart               BIT           NOT
  CreatedDate             DATETIME      NOT
  CreatedBy               NVARCHAR(255) NOT
  UpdatedDate             DATETIME      NULL
  UpdatedBy               NVARCHAR(255) NULL
  Active                  BIT           NOT
  LastJobStart            DATETIME      NULL
  EventType               VARCHAR(100)  NOT
  EventDate               DATETIME      NOT
*/

DROP PROCEDURE IF EXISTS `GetJobConfigAudit`;

CREATE PROCEDURE `GetJobConfigAudit` (inJobConfigID VARCHAR(36), inOrgID VARCHAR(36))
    #BEGIN#
SELECT
    JC.ID,
    JC.JobId,
    JC.OrganizationID,
    JC.DataInSourceConfigID,
    JC.DataOutSourceConfigID,
    JC.Payload,
    JC.PriorityOverride,
    JC.Continuous,
    JC.WaitInSeconds,
    JC.MaxInstances,
    JC.AutoStart,
    JC.CreatedDate,
    JC.CreatedBy,
    JC.UpdatedDate,
    JC.UpdatedBy,
    JC.Active,
    JC.LastJobStart,
    JC.EventType,
    JC.EventDate
FROM JobConfigAudit JC
WHERE JC.OrganizationId = inOrgID AND JC.ID = inJobConfigID;
