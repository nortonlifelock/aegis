/*
  RETURN JobHistory
  ID                      NVARCHAR(36)  NOT
  JobID                   INT           NOT
  ConfigID                NVARCHAR(36)  NOT
  StatusID                INT           NOT
  ParentJobID             NVARCHAR(36)  NULL
  Identifier              NVARCHAR(100) NULL
  Priority                INT           NOT
  Payload                 TEXT          NOT
  ThreadID                NVARCHAR(100) NULL
  PulseDate               DATETIME      NULL
  CreatedDate             DATETIME      NOT
  UpdatedDate             DATETIME      NULL
  MaxInstances            INT           NOT
*/

DROP PROCEDURE IF EXISTS `GetJobQueueByStatusID`;

CREATE PROCEDURE `GetJobQueueByStatusID` (_StatusID INT)
  #BEGIN#
  SELECT
    JH.Id,
    JH.JobId,
    JH.ConfigId,
    JH.StatusId,
    JH.ParentJobId,
    JH.Identifier,
    JH.Priority,
    JH.Payload,
    JH.ThreadId,
    JH.PulseDate,
    JH.CreatedDate,
    JH.UpdatedDate,
    JC.MaxInstances
  FROM JobHistory JH
    JOIN JobConfig JC on JH.ConfigID = JC.ID
  WHERE JH.StatusId = _StatusID
  ORDER BY JH.Priority;