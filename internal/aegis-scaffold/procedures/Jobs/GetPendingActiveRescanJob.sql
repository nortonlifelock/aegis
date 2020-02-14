/*
  RETURN JobHistory
  ID                      NVARCHAR(36)           NOT
  JobID                   INT           NOT
  ConfigID                NVARCHAR(36)           NOT
  StatusID                INT           NOT
  ParentJobID             NVARCHAR(36)           NULL
  Identifier              NVARCHAR(100) NULL
  Priority                INT           NOT
  CurrentIteration        INT           NULL
  Payload                 TEXT          NOT
  ThreadID                NVARCHAR(100) NULL
  PulseDate               DATETIME      NULL
  CreatedDate             DATETIME      NOT
  UpdatedDate             DATETIME      NULL
*/

DROP PROCEDURE IF EXISTS `GetPendingActiveRescanJob`;

CREATE PROCEDURE `GetPendingActiveRescanJob` (_OrgID VARCHAR(36))
  #BEGIN#
  SELECT
    JH.Id,
    JH.JobId,
    JH.ConfigId,
    JH.StatusId,
    JH.ParentJobId,
    JH.Identifier,
    JH.Priority,
    JH.CurrentIteration,
    JH.Payload,
    JH.ThreadId,
    JH.PulseDate,
    JH.CreatedDate,
    JH.UpdatedDate
  FROM JobHistory JH
  JOIN JobConfig JC ON JC.Id = JH.ConfigId
  JOIN Job J ON J.Id = JC.JobId
  WHERE JC.OrganizationId = _OrgID
        and (J.Struct = 'RescanJob' OR J.Struct = 'ScanCloseJob')
        AND (JH.StatusId = 1 OR JH.StatusId = 2);