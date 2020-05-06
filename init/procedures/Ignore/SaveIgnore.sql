DROP PROCEDURE IF EXISTS `SaveIgnore`;

CREATE PROCEDURE `SaveIgnore` (_SourceID NVARCHAR(36), _OrganizationID NVARCHAR(36), _TypeID INT, _VulnerabilityID NVARCHAR(120), _DeviceID NVARCHAR(36), _DueDate DATETIME, _Approval NVARCHAR(120), _Active BIT, _port NVARCHAR(15))
  #BEGIN#
    IF EXISTS (SELECT * FROM `Ignore` WHERE DeviceID = _DeviceID AND VulnerabilityID = _VulnerabilityID AND Port = _Port) THEN
        UPDATE `Ignore` SET TypeID = _TypeID, DueDate = _DueDate, Approval = _Approval, Active = _Active, DBUpdatedDate = NOW() WHERE DeviceID = _DeviceID AND VulnerabilityID = _VulnerabilityID AND Port = _Port AND ID != '';
    ELSE
        INSERT INTO `Ignore` (SourceID, OrganizationID, TypeID, VulnerabilityID, DeviceID, DueDate, Approval, Active, Port)
            VALUE (_SourceID, _OrganizationID, _TypeID, _VulnerabilityID, _DeviceID, _DueDate, _Approval, _Active, _port);
    END IF;
