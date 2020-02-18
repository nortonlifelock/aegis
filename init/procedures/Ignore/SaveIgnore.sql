DROP PROCEDURE IF EXISTS `SaveIgnore`;

CREATE PROCEDURE `SaveIgnore` (_SourceID NVARCHAR(36), _OrganizationID NVARCHAR(36), _TypeID INT, _VulnerabilityID NVARCHAR(120), _DeviceID NVARCHAR(36), _DueDate DATETIME, _Approval NVARCHAR(120), _Active BIT, _port NVARCHAR(15))
  #BEGIN#

  INSERT INTO `Ignore` (SourceID, OrganizationID, TypeID, VulnerabilityID, DeviceID, DueDate, Approval, Active, Port)
  VALUE (_SourceID, _OrganizationID, _TypeID, _VulnerabilityID, _DeviceID, _DueDate, _Approval, _Active, _port)
  ON DUPLICATE KEY UPDATE TypeID = _TypeID, DueDate = _DueDate, Approval = _Approval, Active = _Active, DBUpdatedDate = NOW();