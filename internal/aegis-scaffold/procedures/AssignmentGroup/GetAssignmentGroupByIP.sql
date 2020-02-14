/*
  RETURN AssignmentGroup
  SourceID            INT                     NOT
  OrganizationID      VARCHAR(36)             NOT
  IPAddress           NVARCHAR(20)            NOT
  GroupName           NVARCHAR(150)           NOT
  DBCreatedDate       DATETIME                NOT
  DBUpdatedDate       DATETIME                NULL
*/

DROP PROCEDURE IF EXISTS `GetAssignmentGroupByIP`;

CREATE PROCEDURE `GetAssignmentGroupByIP` (_SourceID VARCHAR(36), _OrganizationID VARCHAR(36), _IP NVARCHAR(20))
#BEGIN#
  SELECT
    AG.SourceID,
    AG.OrganizationID,
    AG.IpAddress,
    AG.GroupName,
    AG.DBCreatedDate,
    AG.DBUpdatedDate
  FROM AssignmentGroup AG
  WHERE AG.SourceID = _SourceID
        AND AG.OrganizationID = _OrganizationID
        AND AG.IpAddress = _IP;
