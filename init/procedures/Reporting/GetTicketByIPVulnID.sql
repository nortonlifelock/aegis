/*
    RETURN TicketSummary
    Title          VARCHAR(36)       NOT NULL
    Status         VARCHAR(100)      NOT NULL
    DetectionID    VARCHAR(36)       NOT NULL
    OrganizationID VARCHAR(100)      NOT NULL
    UpdatedDate    DATETIME          NULL
    ResolutionDate DATETIME          NULL
    DueDate        DATETIME          NOT NULL
*/

DROP PROCEDURE IF EXISTS `GetTicketByIPVulnID`;

CREATE PROCEDURE `GetTicketByIPVulnID`(_IP VARCHAR(100), _VulnID VARCHAR(100))
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
WHERE
      T.Status != 'Closed-NA' and
      D.DeviceID IN (select AssetID from Device where IP = _IP) and
      D.VulnerabilityID = (select ID from VulnerabilityInfo where SourceVulnID = _VulnID)
order by T.Created desc;