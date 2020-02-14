/*
  RETURN QueryData SINGLE
  Length                      INT           NOT
*/
DROP PROCEDURE IF EXISTS `GetJobHistoryLength`;

CREATE PROCEDURE `GetJobHistoryLength` (_jobid INT, _jobconfig VARCHAR(36), _status INT, _Payload MEDIUMTEXT, _orgid VARCHAR(36))
  #BEGIN#

 SELECT
  count(*)
  FROM JobHistory JH
  JOIN JobConfig JC ON JH.ConfigId = JC.Id
where JC.OrganizationId = _OrgId
  AND (JH.JobId = _jobid OR _jobid = '' OR _jobid is NULL)
AND (JH.ConfigId = _jobconfig OR  _jobconfig = '' OR _jobconfig is NULL)
AND (JH.StatusId = _status OR _status ='' OR _status is NULL)
AND (JH.Payload = _payload OR _payload ='' OR _payload is NULL);