/*
    RETURN TicketSummary SINGLE
    Title          VARCHAR(36)       NOT NULL
    Status         VARCHAR(100)      NOT NULL
    DetectionID    VARCHAR(36)       NOT NULL
    OrganizationID VARCHAR(100)      NOT NULL
    UpdatedDate    DATETIME          NULL
    ResolutionDate DATETIME          NULL
    DueDate        DATETIME          NOT NULL
*/

DROP PROCEDURE IF EXISTS `GetTicketByTitle`;

CREATE PROCEDURE `GetTicketByTitle` (_Title VARCHAR(36), _OrgID NVARCHAR(36))
    #BEGIN#
SELECT
    T.Title,
    T.Status,
    T.DetectionID,
    T.OrganizationID,
    T.UpdatedDate,
    T.ResolutionDate,
    T.DueDate
FROM Ticket T
WHERE T.Title = _Title AND T.OrganizationID = _OrgID;