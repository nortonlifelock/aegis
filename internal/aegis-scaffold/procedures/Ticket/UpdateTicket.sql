DROP PROCEDURE IF EXISTS `UpdateTicket`;

CREATE PROCEDURE `UpdateTicket`(_Title VARCHAR(36), _Status VARCHAR(100), _OrganizationID VARCHAR(100), _AssignmentGroup VARCHAR(200), _Assignee VARCHAR(200), _CreatedDate DATETIME, _UpdatedDate DATETIME, _ResolutionDate DATETIME, _DefaultTime DATETIME)

#BEGIN#
UPDATE Ticket T SET T.Status = _Status, T.AssignmentGroup = _AssignmentGroup, T.Assignee = _Assignee, T.Created = NULLIF(_CreatedDate, _DefaultTime), T.UpdatedDate = _UpdatedDate, T.ResolutionDate = NULLIF(_ResolutionDate, _DefaultTime) WHERE T.Title = _Title AND T.OrganizationID = _OrganizationID;