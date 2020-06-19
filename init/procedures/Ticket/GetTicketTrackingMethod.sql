/*
    RETURN KeyValue SINGLE
    Key   VARCHAR(100) NOT
    Value VARCHAR(100) NOT
*/

DROP PROCEDURE IF EXISTS `GetTicketTrackingMethod`;

CREATE PROCEDURE `GetTicketTrackingMethod` (_Title VARCHAR(36), _OrgID VARCHAR(36))
    #BEGIN#
SELECT
    'TrackingMethod',
    IFNULL(Dev.TrackingMethod, '')
FROM Ticket T
         JOIN Detection D ON T.DetectionID = D.ID
         JOIN Device Dev ON Dev.AssetID = D.DeviceID
WHERE T.Title = _Title and T.OrganizationID = _OrgID and D.OrganizationID = _OrgID and Dev.OrganizationID = _OrgID;