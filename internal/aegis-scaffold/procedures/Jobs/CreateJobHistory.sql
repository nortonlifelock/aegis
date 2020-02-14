DROP PROCEDURE IF EXISTS `CreateJobHistory`;

CREATE PROCEDURE `CreateJobHistory` (_JobID INT, _ConfigID NVARCHAR(36), _StatusID INT, _Priority INT, _Identifier NVARCHAR(100), _CurrentIteration INT, _Payload MEDIUMTEXT, _ThreadID NVARCHAR(100), _PulseDate DATETIME, _CreatedBy NVARCHAR(255))
  #BEGIN#

  INSERT INTO JobHistory (JobID, ConfigID, StatusID, Priority, IDentifier, CurrentIteration, Payload, ThreadID, PulseDate, CreatedBy)
    VALUE (_JobID, _ConfigID, _StatusID, _Priority, _Identifier, _CurrentIteration, _Payload, _ThreadID, _PulseDate, _CreatedBy);
