DROP PROCEDURE IF EXISTS `CreateJobHistoryWithParentID`;

CREATE PROCEDURE `CreateJobHistoryWithParentID` (_JobID INT, _ConfigID NVARCHAR(36), _StatusID INT, _Priority INT, _Identifier NVARCHAR(100), _CurrentIteration INT, _Payload MEDIUMTEXT, _ThreadID NVARCHAR(100), _PulseDate DATETIME, _CreatedBy NVARCHAR(255), _ParentID VARCHAR(36))
  #BEGIN#

  INSERT INTO JobHistory (JobID, ConfigID, StatusID, ParentJobID, Priority, IDentifier, CurrentIteration, Payload, ThreadID, PulseDate, CreatedBy)
    VALUE (_JobID, _ConfigID, _StatusID, _ParentID, _Priority, _Identifier, _CurrentIteration, _Payload, _ThreadID, _PulseDate, _CreatedBy);
