DROP PROCEDURE IF EXISTS `UpdateTagMap`;

CREATE PROCEDURE `UpdateTagMap` (_TicketingSourceID NVARCHAR(36), _TicketingTag NVARCHAR(255),
                                                      _CloudSourceID NVARCHAR(36), _CloudTag NVARCHAR(255), _Options NVARCHAR(255), _OrganizationID NVARCHAR(36))
  #BEGIN#
  UPDATE TagMap T
  SET T.Options = _Options
  WHERE T.TicketingSourceID = _TicketingSourceID AND T.TicketingTag = _TicketingTag
        AND T.CloudSourceID = _CloudSourceID AND T.CloudTag = _CloudTag AND T.OrganizationID = _OrganizationID;