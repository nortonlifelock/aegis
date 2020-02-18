/*
  RETURN DBLog
  ID                INT                     NOT
  TypeID            INT                     NOT
  Log               NVARCHAR(255)           NOT
  Error             NVARCHAR(255)           NOT
  JobHistoryID      VARCHAR(36)             NOT
  CreateDate        DATETIME                NOT
*/

DROP PROCEDURE IF EXISTS `GetLogsByParams`;

CREATE PROCEDURE `GetLogsByParams` (_MethodOfDiscovery NVARCHAR(50), _jobType INT, _logType INT, _jobHistoryID VARCHAR(36), _fromDate DATETIME, _toDate DATETIME, _OrgID VARCHAR(36))
  #BEGIN#
  SELECT
    L.ID,
    L.TypeID,
    L.Log,
    L.Error,
    L.JobHistoryID,
    L.CreateDate
  FROM Logs L
    JOIN JobHistory JH ON JH.ID = L.JobHistoryID
    JOIN JobConfig JC ON JC.ID = JH.ConfigID
    JOIN SourceConfig SC ON SC.ID = JC.DataInSourceConfigID
    JOIN SourceConfig SCT ON SCT.ID = JC.DataOutSourceConfigID
  WHERE JC.OrganizationID = _OrgID
    AND (_MethodOfDiscovery = '' OR SC.Source = _MethodOfDiscovery OR SCT.Source = _MethodOfDiscovery)
    AND (_jobType = 0 OR JC.JobID = _jobType)
    AND (_logType = -1 OR L.TypeID = _logType)
    AND (_jobHistoryID = 0 OR L.JobHistoryID = _jobHistoryID)
    AND (L.CreateDate > _fromDate)
    AND (L.CreateDate < _toDate)
  ORDER BY L.ID DESC LIMIT 100; -- TODO should we limit at all? I could make this an input