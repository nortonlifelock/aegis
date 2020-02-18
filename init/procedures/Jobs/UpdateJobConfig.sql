DROP PROCEDURE IF EXISTS `UpdateJobConfig`;

CREATE PROCEDURE `UpdateJobConfig` (_ID NVARCHAR(36), _DataInSourceID NVARCHAR(36), _DataOutSourceID NVARCHAR(36), _Autostart BIT, _PriorityOverride INT,
                                            _Continuous BIT, _WaitInSeconds INT, _MaxInstances INT, _UpdatedBy TEXT, _OrgID VARCHAR(36))
  #BEGIN#
  BEGIN
    UPDATE JobConfig SET DataInSourceConfigID = _DataInSourceID, DataOutSourceConfigID = _DataOutSourceID,
      AutoStart = _Autostart, PriorityOverride = _PriorityOverride, Continuous = _Continuous, WaitInSeconds = _WaitInSeconds,
      MaxInstances = _MaxInstances, UpdatedDate = NOW(), UpdatedBy = _UpdatedBy WHERE ID = _ID AND OrganizationID = _OrgID;
  END