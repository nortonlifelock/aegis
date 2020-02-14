/*
  RETURN TagMap
  ID                  NVARCHAR(36)             NOT
  TicketingSourceID   NVARCHAR(36)             NOT
  TicketingTag        NVARCHAR(255)   NOT
  CloudSourceID       NVARCHAR(36)             NOT
  CloudTag            NVARCHAR(255)   NOT
  Options             NVARCHAR(255)   NOT
*/

DROP PROCEDURE IF EXISTS `GetTagMapsByOrg`;

CREATE PROCEDURE `GetTagMapsByOrg` (_OrganizationID NVARCHAR(36))
  #BEGIN#
  SELECT
    T.Id,
    T.TicketingSourceId,
    T.TicketingTag,
    T.CloudSourceId,
    T.CloudTag,
    T.Options
  FROM TagMap T
  WHERE T.OrganizationId = _OrganizationID AND T.Active = 1;