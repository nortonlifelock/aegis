/*
  RETURN Ignore SINGLE
  ID              VARCHAR(30)     NOT
  SourceID        VARCHAR(36)     NOT
  OrganizationID  VARCHAR(36)     NOT
  TypeID          INT             NOT
  VulnerabilityID NVARCHAR(120)   NOT
  DeviceID        VARCHAR(36)     NOT
  DueDate         DATETIME        NULL
  Approval        NVARCHAR(120)   NOT
  Active          BIT             NOT
  DBCreatedDate   DATETIME        NOT
  DBUpdatedDate   DATETIME        NULL
 */

DROP PROCEDURE IF EXISTS `HasDecommissioned`;

CREATE PROCEDURE `HasDecommissioned`( _devID VARCHAR(36), _sourceID VARCHAR(36), _orgID VARCHAR(36))
    #BEGIN#
SELECT
    ID,
    SourceID,
    OrganizationID,
    TypeID,
    VulnerabilityID,
    DeviceID ,
    DueDate,
    Approval,
    Active,
    DBCreatedDate,
    DBUpdatedDate
FROM `Ignore` I
WHERE  I.DeviceID = _devID AND I.SourceID = _sourceID AND I.OrganizationID = _orgID
  and I.TypeID = 2 AND I.Active = b'1';