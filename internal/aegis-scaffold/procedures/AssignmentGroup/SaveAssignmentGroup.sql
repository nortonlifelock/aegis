DROP PROCEDURE IF EXISTS `SaveAssignmentGroup`;

CREATE PROCEDURE `SaveAssignmentGroup` (_SourceID VARCHAR(36), _OrganizationID VARCHAR(36), _IpAddress NVARCHAR(20), _GroupName NVARCHAR(150))
  #BEGIN#

  INSERT INTO AssignmentGroup (SourceID, OrganizationID, IpAddress, GroupName)
  VALUE (_SourceID, _OrganizationID, _IpAddress, _GroupName)
  ON DUPLICATE KEY UPDATE GroupName = _GroupName;