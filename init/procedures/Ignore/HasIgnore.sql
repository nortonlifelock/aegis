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

DROP PROCEDURE IF EXISTS `HasIgnore`;

CREATE PROCEDURE `HasIgnore`(inSourceID VARCHAR(36), inVulnID NVARCHAR(255), inDevID VARCHAR(36), inOrgID VARCHAR(36), inPort NVARCHAR(15), inMostCurrentDetection DATETIME)
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
  WHERE I.SourceID = inSourceID
  and I.OrganizationID = inOrgID
  and I.DeviceID = inDevID
  and I.VulnerabilityID = inVulnID
  and ((inPort <> '' and I.Port = inPort) OR inPort ='')
  and ((I.TypeID = 1 )
       or (I.typeID = 0
          and I.DueDate > inMostCurrentDetection)) AND I.Active = b'1';