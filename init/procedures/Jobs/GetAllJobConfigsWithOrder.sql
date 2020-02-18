/*
  RETURN JobConfig
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
  Active             BIT      NOT
*/

DROP PROCEDURE IF EXISTS `GetAllJobConfigsWithOrder`;

CREATE PROCEDURE `GetAllJobConfigsWithOrder` (_offset INT, _limit INT, _configID VARCHAR(36), _jobid INT, _dataInSourceConfigID VARCHAR(36), _dataOutSourceConfigID VARCHAR(36),
                                          _priorityOverride INT, _continuous BIT, _Payload MEDIUMTEXT,_waitInSeconds INT, _maxInstances INT, _autoStart BIT, _OrgID VARCHAR(36), _updatedBy NVARCHAR(255), _createdBy NVARCHAR(255) , _sortField NVARCHAR(255),
                                          _sortOrder NVARCHAR(255), _updatedDate DATETIME ,_createdDate DATETIME,_lastJobStart DATETIME, _ID VARCHAR(36))
    #BEGIN#
BEGIN
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
        JC.LastJobStart,
        JC.Active
    FROM JobConfig JC
             JOIN SourceConfig ScIn on JC.DataInSourceConfigId=ScIn.Id
             JOIN SourceConfig ScOut on JC.DataOutSourceConfigId=ScOut.Id
             JOIN Job J on JC.JobId=J.Id
    WHERE JC.OrganizationId = _OrgId
      AND (JC.JobId = _jobid OR _jobid = '' OR _jobid is NULL)
      AND (JC.Id = _id OR _id = '' OR _id is NULL)
      AND (JC.Id = _configId OR _configId= 0 OR _configId is NULL)
      AND (JC.DataInSourceConfigId = _dataInSourceConfigId OR _dataInSourceConfigId = '' OR _dataInSourceConfigId is NULL)
      AND (JC.DataOutSourceConfigId = _dataOutSourceConfigId OR _dataOutSourceConfigId ='' OR _dataOutSourceConfigId is NULL)
      AND (JC.PriorityOverride = _priorityOverride OR _priorityOverride ='' OR _priorityOverride is NULL)
      AND (JC.Continuous = _continuous OR _continuous ='' OR _continuous is NULL)
      AND (JC.WaitInSeconds = _waitInSeconds OR _waitInSeconds ='' OR _waitInSeconds is NULL)
      AND (JC.MaxInstances = _maxInstances OR _maxInstances ='' OR _maxInstances is NULL)
      AND (JC.AutoStart = _autoStart OR _autoStart ='' OR _autoStart is NULL)
      AND (JC.Payload = _payload OR _payload ='' OR _payload is NULL)
      AND (JC.UpdatedBy = _updatedBy OR _updatedBy='' OR _updatedBy is NULL)
      AND (JC.CreatedBy = _createdBy OR _createdBy ='' OR _createdBy is NULL)
      AND (JC.CreatedDate = _createdDate OR _createdDate ='1970-01-02 00:00:00 +0000 UTC' OR _createdDate is NULL)
      AND (JC.UpdatedDate = _updatedDate OR _updatedDate ='1970-01-02 00:00:00 +0000 UTC' OR _updatedDate is NULL)
      AND (JC.LastJobStart = _lastJobStart OR _lastJobStart ='1970-01-02 00:00:00 +0000 UTC' OR _lastJobStart is NULL)
    ORDER BY
        CASE WHEN _sortField = 'job_id' AND _sortOrder='ASC' THEN J.Struct END,
        CASE WHEN _sortField = 'job_id' AND _sortOrder='DESC' THEN J.Struct END DESC,
        CASE WHEN _sortField = 'config_id' AND _sortOrder='ASC' THEN JC.Id END,
        CASE WHEN _sortField = 'config_id' AND _sortOrder='DESC' THEN JC.Id END DESC,
        CASE WHEN _sortField = 'priority_override' AND _sortOrder='ASC' THEN JC.PriorityOverride END,
        CASE WHEN _sortField = 'priority_override' AND _sortOrder='DESC' THEN JC.PriorityOverride END DESC,
        CASE WHEN _sortField = 'data_in_source_config_id' AND _sortOrder='ASC' THEN ScIn.Source END,
        CASE WHEN _sortField = 'data_in_source_config_id' AND _sortOrder='DESC' THEN ScIn.Source END DESC,
        CASE WHEN _sortField = 'data_out_source_config_id' AND _sortOrder='ASC' THEN ScOut.Source END,
        CASE WHEN _sortField = 'data_out_source_config_id' AND _sortOrder='DESC' THEN ScOut.Source END DESC,
        CASE WHEN _sortField = 'continuous' AND _sortOrder='ASC' THEN JC.Continuous END,
        CASE WHEN _sortField = 'continuous' AND _sortOrder='DESC' THEN JC.Continuous END DESC,
        CASE WHEN _sortField = 'payload' AND _sortOrder='ASC' THEN JC.Payload END,
        CASE WHEN _sortField = 'payload' AND _sortOrder='DESC' THEN JC.Payload END DESC,
        CASE WHEN _sortField = 'wait_in_seconds' AND _sortOrder='ASC' THEN JC.WaitInSeconds END,
        CASE WHEN _sortField = 'wait_in_seconds' AND _sortOrder='DESC' THEN JC.WaitInSeconds END DESC,
        CASE WHEN _sortField = 'max_instances' AND _sortOrder='ASC' THEN JC.MaxInstances END,
        CASE WHEN _sortField = 'max_instances' AND _sortOrder='DESC' THEN JC.MaxInstances END DESC,
        CASE WHEN _sortField = 'autostart' AND _sortOrder='ASC' THEN JC.AutoStart END,
        CASE WHEN _sortField = 'autostart' AND _sortOrder='DESC' THEN JC.AutoStart END DESC,
        CASE WHEN _sortField = 'created_by' AND _sortOrder='ASC' THEN JC.CreatedBy END,
        CASE WHEN _sortField = 'created_by' AND _sortOrder='DESC' THEN JC.CreatedBy END DESC,
        CASE WHEN _sortField = 'created_date' AND _sortOrder='ASC' THEN JC.CreatedDate END,
        CASE WHEN _sortField = 'created_date' AND _sortOrder='DESC' THEN JC.CreatedDate END DESC,
        CASE WHEN _sortField = 'updated_date' AND _sortOrder='ASC' THEN JC.UpdatedDate END,
        CASE WHEN _sortField = 'updated_date' AND _sortOrder='DESC' THEN JC.UpdatedDate END DESC,
        CASE WHEN _sortField = 'updated_by' AND _sortOrder='ASC' THEN JC.UpdatedBy END,
        CASE WHEN _sortField = 'updated_by' AND _sortOrder='DESC' THEN JC.UpdatedBy END DESC,
        CASE WHEN _sortField = 'last_job_start' AND _sortOrder='ASC' THEN JC.LastJobStart END,
        CASE WHEN _sortField = 'last_job_start' AND _sortOrder='DESC' THEN JC.LastJobStart END DESC
    LIMIT _offset,_limit;
END