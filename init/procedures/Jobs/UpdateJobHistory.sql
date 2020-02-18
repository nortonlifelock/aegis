DROP PROCEDURE IF EXISTS `UpdateJobHistory`;

CREATE PROCEDURE `UpdateJobHistory` (_ID NVARCHAR(36), _ConfigID NVARCHAR(36), _Payload MEDIUMTEXT, _UpdatedBy TEXT)
  #BEGIN#

 BEGIN
  UPDATE JobHistory SET UpdatedDate = NOW(), UpdatedBy = _UpdatedBy,
    ConfigID = _ConfigID, Payload = _Payload WHERE ID = _ID;
  END