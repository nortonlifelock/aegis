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

DROP PROCEDURE IF EXISTS `GetTicketByIPGroupIDVulnID`;

CREATE PROCEDURE `GetTicketByIPGroupIDVulnID` (inIP VARCHAR(100), inGroupID VARCHAR(100), inVulnID VARCHAR(36), inPort INT, inProtocol VARCHAR(100), inOrgID NVARCHAR(36))
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
         JOIN Detection D on T.DetectionID = D.ID
         JOIN Device Dev ON Dev.AssetID = D.DeviceID
WHERE Dev.Ip = inIP AND Dev.GroupID = inGroupID AND D.VulnerabilityID = inVulnID AND (inPort = '' OR D.Port = inPort) and (inProtocol = '' OR D.Protocol = inProtocol) AND T.OrganizationID = inOrgID AND T.Status IN ('Open', 'In-Progress', 'Reopened', 'Resolved-Remediated', 'Resolved-FalsePositive', 'Resolved-Decommissioned', 'Resolved-Exception');