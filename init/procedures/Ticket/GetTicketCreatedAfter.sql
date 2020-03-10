/*
    RETURN TicketSummary
    Title          VARCHAR(36)       NOT NULL
    Status         VARCHAR(100)      NOT NULL
    DetectionID    VARCHAR(36)       NOT NULL
    OrganizationID VARCHAR(100)      NOT NULL
    CreatedDate    DATETIME          NULL
    UpdatedDate    DATETIME          NULL
    ResolutionDate DATETIME          NULL
    DueDate        DATETIME          NOT NULL
*/

DROP PROCEDURE IF EXISTS `GetTicketCreatedAfter`;

CREATE PROCEDURE `GetTicketCreatedAfter` (_UpperCVSS FLOAT, _LowerCVSS FLOAT, _CreatedAfter DATETIME, _OrgID NVARCHAR(36))
    #BEGIN#
SELECT
    T.Title,
    T.Status,
    T.DetectionID,
    T.OrganizationID,
    T.Created,
    T.UpdatedDate,
    T.ResolutionDate,
    T.DueDate
FROM Ticket T
     JOIN Detection D ON T.DetectionID = D.ID
     JOIN VulnerabilityInfo V ON V.ID = D.VulnerabilityID
WHERE T.Created > _CreatedAfter AND T.OrganizationID = _OrgID AND V.CVSS <= _UpperCVSS AND V.CVSS >= _LowerCVSS;