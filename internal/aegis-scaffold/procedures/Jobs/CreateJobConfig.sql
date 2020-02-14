DROP PROCEDURE IF EXISTS `CreateJobConfig`;

CREATE PROCEDURE `CreateJobConfig` (_JobID INT, _OrganizationID NVARCHAR(36), _PriorityOverride INT, _Continuous BIT, _WaitInSeconds INT, _MaxInstances INT, _AutoStart BIT, _CreatedBy NVARCHAR(255), _DataInSourceID NVARCHAR(36), _DataOutSourceID NVARCHAR(36))
  #BEGIN#

  INSERT INTO JobConfig (JobID, OrganizationID, PriorityOverride, Continuous, WaitInSeconds, MaxInstances, AutoStart, CreatedBy, DataInSourceConfigID, DataOutSourceConfigID)
    VALUE (_JobID, _OrganizationID, _PriorityOverride, _Continuous, _WaitInSeconds, _MaxInstances, _AutoStart, _CreatedBy, _DataInSourceID, _DataOutSourceID);