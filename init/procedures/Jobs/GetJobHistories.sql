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

DROP PROCEDURE IF EXISTS `GetJobHistories`;

CREATE PROCEDURE `GetJobHistories` (_offset INT, _limit INT, _jobID INT, _jobconfig VARCHAR(36), _status INT, _Payload MEDIUMTEXT, _OrgID VARCHAR(36))
  #BEGIN#
BEGIN
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
  JOIN JobConfig JC ON JH.ConfigId = JC.Id
where JC.OrganizationId = _OrgID
and (JH.JobId = _jobID OR _jobID = '' OR _jobID is NULL)
AND (JH.ConfigId = _jobconfig OR  _jobconfig = '' OR _jobconfig is NULL)
AND (JH.StatusId = _status OR _status ='' OR _status is NULL)
AND (JH.Payload LIKE CONCAT('%', _payload , '%') OR _payload ='' OR _payload is NULL)
ORDER BY JH.Id DESC
LIMIT _offset,_limit;
END