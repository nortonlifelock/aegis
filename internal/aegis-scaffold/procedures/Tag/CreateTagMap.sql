DROP PROCEDURE IF EXISTS `CreateTagMap`;

CREATE PROCEDURE `CreateTagMap`(_TicketingSourceID NVARCHAR(36), _TicketingTag NVARCHAR(255),
  _CloudSourceID NVARCHAR(36), _CloudTag NVARCHAR(255), _Options NVARCHAR(255), _OrganizationID NVARCHAR(36))

  #BEGIN#
  INSERT INTO TagMap (TicketingSourceID, TicketingTag, CloudSourceID, CloudTag, Options, OrganizationID)
  VALUES (_TicketingSourceID, _TicketingTag, _CloudSourceID, _CloudTag, _Options, _OrganizationID);