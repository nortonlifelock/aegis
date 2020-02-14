/*
  RETURN JobConfig
  ID                      NVARCHAR(36)  NOT
  JobID                   INT           NOT
  OrganizationID          NVARCHAR(36)  NOT
  DataInSourceConfigID    NVARCHAR(36)  NOT
  DataOutSourceConfigID   NVARCHAR(36)  NOT
  PriorityOverride        INT           NULL
  Continuous              BIT           NOT
  WaitInSeconds           INT           NOT
  MaxInstances            INT           NOT
  AutoStart               BIT           NOT
  CreatedDate             DATETIME      NOT
  CreatedBy               NVARCHAR(255) NOT
  UpdatedDate             DATETIME      NULL
  UpdatedBy               NVARCHAR(255) NULL
  Payload                 TEXT          NULL
*/

DROP PROCEDURE IF EXISTS `GetAutoStartJobs`;

CREATE PROCEDURE `GetAutoStartJobs` ()
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
    JC.Payload
  FROM JobConfig JC
  WHERE JC.Continuous = TRUE
        AND JC.AutoStart = TRUE;