DROP PROCEDURE IF EXISTS `CreateTicket`;

CREATE PROCEDURE `CreateTicket`(_Title VARCHAR(36), _Status VARCHAR(100), _DetectionID VARCHAR(36), _OrganizationID VARCHAR(100), _DueDate DATETIME, _UpdatedDate DATETIME, _ResolutionDate DATETIME, _DefaultTime DATETIME)
#BEGIN#

INSERT INTO Ticket (Title, Status, DetectionID, OrganizationID, DueDate, UpdatedDate, ResolutionDate, Created)
    VALUES (_Title, _Status, _DetectionID, _OrganizationID, _DueDate, _UpdatedDate, NULLIF(_ResolutionDate, _DefaultTime), NOW());
