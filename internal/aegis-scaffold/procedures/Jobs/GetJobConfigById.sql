/*
  RETURN JobConfig SINGLE
  ID                      NVARCHAR(36)           NOT
  JobID                   INT           NOT
  OrganizationID          NVARCHAR(36)           NOT
  DataInSourceConfigID    NVARCHAR(36)  NULL
  DataOutSourceConfigID   NVARCHAR(36)  NULL
  PriorityOverride        INT           NULL
  Continuous              BIT           NOT
  Payload                 TEXT          NULL
  WaitInSeconds           INT           NOT
  MaxInstances            INT           NOT
  AutoStart               BIT           NOT
  CreatedDate             DATETIME      NOT
  CreatedBy               NVARCHAR(255) NOT
  UpdatedDate             DATETIME      NULL
  UpdatedBy               NVARCHAR(255) NULL
  LastJobStart             DATETIME      NULL
*/

DROP PROCEDURE IF EXISTS `GetJobConfigByID`;

CREATE PROCEDURE `GetJobConfigByID` (_ID VARCHAR(36), _OrgID VARCHAR(36))
  #BEGIN#
  SELECT
    JC.Id,
    JC.JobId,
    JC.OrganizationId,
    JC.DataInSourceConfigId,
    JC.DataOutSourceConfigId,
    JC.PriorityOverride,
    JC.Continuous,
    JC.Payload,
    JC.WaitInSeconds,
    JC.MaxInstances,
    JC.AutoStart,
    JC.CreatedDate,
    JC.CreatedBy,
    JC.UpdatedDate,
    JC.UpdatedBy,
    JC.LastJobStart
  FROM JobConfig JC
  WHERE JC.Id = _ID AND JC.OrganizationId = _OrgID;