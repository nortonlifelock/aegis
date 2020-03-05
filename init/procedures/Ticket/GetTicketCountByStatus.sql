/*
  RETURN QueryData SINGLE
  Length                      INT           NOT
*/

DROP PROCEDURE IF EXISTS `GetTicketCountByStatus`;

CREATE PROCEDURE `GetTicketCountByStatus`(inStatus VARCHAR(100), inOrgID VARCHAR(36))
#BEGIN#
    SELECT
        count(*)
    FROM Ticket T WHERE T.Status = inStatus AND T.OrganizationID = inOrgID;