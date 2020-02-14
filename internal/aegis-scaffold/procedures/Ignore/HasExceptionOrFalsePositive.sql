/*
  RETURN Ignore
  ID              VARCHAR(30)     NOT
  SourceID        NVARCHAR(36)    NOT
  OrganizationID  NVARCHAR(36)    NOT
  TypeID          INT             NOT
  VulnerabilityID NVARCHAR(120)   NOT
  DeviceID        NVARCHAR(36)    NOT
  DueDate         DATETIME        NULL
  Approval        NVARCHAR(120)   NOT
  Active          BIT             NOT
  DBCreatedDate   DATETIME        NOT
  DBUpdatedDate   DATETIME        NULL
 */

DROP PROCEDURE IF EXISTS `HasExceptionOrFalsePositive`;

CREATE PROCEDURE `HasExceptionOrFalsePositive`(_sourceID VARCHAR(36), _vulnID NVARCHAR(255), _devID VARCHAR(36), _orgID VARCHAR(36), _port NVARCHAR(15), _OS VARCHAR(100))
  #BEGIN#
  SELECT
    ID,
    SourceId,
    OrganizationId,
    TypeId,
    VulnerabilityId,
    DeviceId ,
    DueDate,
    Approval,
    Active,
    DBCreatedDate,
    DBUpdatedDate
  FROM `Ignore` I
  WHERE I.SourceId = _sourceID
  and I.OrganizationId = _orgID
  and I.VulnerabilityId = _vulnID
  and (I.DeviceId = _devID)
  and ((_port <> '' and I.Port = _port) OR _port ='')
  and ((I.TypeId = 1)
       or (I.typeId = 0
          and I.DueDate > NOW())
      ) AND I.Active = b'1';