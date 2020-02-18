/*
  RETURN QueryData SINGLE
  Length                      INT           NOT
*/
DROP PROCEDURE IF EXISTS `GetJobConfigLength`;

CREATE PROCEDURE `GetJobConfigLength` (_configID VARCHAR(36), _jobID INT, _dataInSourceConfigID VARCHAR(36), _dataOutSourceConfigID VARCHAR(36),
                                               _priorityOverride INT, _continuous BIT, _Payload MEDIUMTEXT, _waitInSeconds INT, _maxInstances INT, _autoStart BIT, _OrgID VARCHAR(36),
                                               _updatedBy NVARCHAR(255), _createdBy NVARCHAR(255) , _updatedDate DATETIME ,_createdDate DATETIME,_lastJobStart DATETIME, _ID VARCHAR(36))
    #BEGIN#
SELECT
    count(*)
FROM JobConfig JC
WHERE JC.OrganizationId = _OrgId
  AND JC.Active = 1
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