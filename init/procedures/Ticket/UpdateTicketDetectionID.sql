DROP PROCEDURE IF EXISTS `UpdateTicketDetectionID`;

CREATE PROCEDURE `UpdateTicketDetectionID`(_Title VARCHAR(36), _DetectionID VARCHAR(36), _OrganizationID VARCHAR(100))

#BEGIN#
UPDATE Ticket T SET T.DetectionID = _DetectionID WHERE T.Title = _Title AND T.OrganizationID = _OrganizationID;