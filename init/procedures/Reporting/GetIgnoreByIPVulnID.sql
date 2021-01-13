/*
  RETURN Ignore
  ID              VARCHAR(30)     NOT
  SourceID        VARCHAR(36)     NOT
  OrganizationID  VARCHAR(36)     NOT
  TypeID          INT             NOT
  VulnerabilityID NVARCHAR(120)   NOT
  DeviceID        VARCHAR(36)     NOT
  DueDate         DATETIME        NULL
  Approval        NVARCHAR(120)   NOT
  Active          BIT             NOT
  Port            NVARCHAR(120)   NOT
  CreatedBy       NVARCHAR(255)   NULL
  UpdatedBy       NVARCHAR(255)   NULL
  DBCreatedDate   DATETIME        NOT
  DBUpdatedDate   DATETIME        NULL
 */

DROP PROCEDURE IF EXISTS `GetIgnoreByIPVulnID`;

CREATE PROCEDURE `GetIgnoreByIPVulnID`(_IP VARCHAR(100), _VulnID VARCHAR(100))
    #BEGIN#
SELECT
    I.ID,
    I.SourceID,
    I.OrganizationID,
    I.TypeID,
    I.VulnerabilityID,
    I.DeviceID ,
    I.DueDate,
    I.Approval,
    I.Active,
    I.Port,
    I.CreatedBy,
    I.UpdatedBy,
    I.DBCreatedDate,
    I.DBUpdatedDate
FROM `Ignore` I where
    I.DeviceID IN (select AssetID from Device where IP = _IP) and
    I.VulnerabilityId = _VulnID
order by I.DBUpdatedDate desc;