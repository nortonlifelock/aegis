DROP PROCEDURE IF EXISTS `UpdateTicket`;

CREATE PROCEDURE `UpdateTicket`(_Title VARCHAR(36), _Status VARCHAR(100), _OrganizationID VARCHAR(100), _AssignmentGroup VARCHAR(200), _Assignee VARCHAR(200), _DueDate DATETIME,_CreatedDate DATETIME, _UpdatedDate DATETIME, _ResolutionDate DATETIME, _ExceptionDate DATETIME, _DefaultTime DATETIME)

#BEGIN#
UPDATE Ticket T SET T.Status = _Status, T.AssignmentGroup = _AssignmentGroup, T.Assignee = _Assignee, T.DueDate = NULLIF(_DueDate, _DefaultTime), T.Created = NULLIF(_CreatedDate, _DefaultTime), T.UpdatedDate = _UpdatedDate, T.ResolutionDate = NULLIF(_ResolutionDate, _DefaultTime), T.ExceptionDate = NULLIF(_ExceptionDate, _DefaultTime) WHERE T.Title = _Title AND T.OrganizationID = _OrganizationID;