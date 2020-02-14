/*
  RETURN JobConfig SINGLE
  ID                      NVARCHAR(36)           NOT
  JobID                   INT           NOT
  OrganizationID          NVARCHAR(36)           NOT
  DataInSourceConfigID    NVARCHAR(36)  NULL
  DataOutSourceConfigID   NVARCHAR(36)  NULL
  Payload                 TEXT          NOT
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

DROP PROCEDURE IF EXISTS `GetJobConfig`;

CREATE PROCEDURE `GetJobConfig` (_ID VARCHAR(36))
  #BEGIN#
  SELECT
    JC.Id,
    JC.JobId,
    JC.OrganizationId,
    JC.DataInSourceConfigId,
    JC.DataOutSourceConfigId,
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
    JC.LastJobStart
  FROM JobConfig JC
  WHERE JC.Id = _ID;