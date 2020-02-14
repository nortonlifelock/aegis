DROP PROCEDURE IF EXISTS `DeleteTagMap`;

CREATE PROCEDURE `DeleteTagMap` (_TicketingSourceID NVARCHAR(36), _TicketingTag NVARCHAR(255),
                                                      _CloudSourceID NVARCHAR(36), _CloudTag NVARCHAR(255), _OrganizationID NVARCHAR(36))
  #BEGIN#
  UPDATE TagMap T
    SET T.Active = b'0'
  WHERE T.TicketingSourceID = _TicketingSourceID AND T.TicketingTag = _TicketingTag
    AND T.CloudSourceID = _CloudSourceID AND T.CloudTag = _CloudTag AND T.OrganizationID = _OrganizationID;