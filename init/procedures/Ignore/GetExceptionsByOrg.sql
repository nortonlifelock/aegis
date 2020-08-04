/*
  RETURN Ignore
  ID              NVARCHAR(36)          NOT
  TypeID          INT                   NOT
  Approval        VARCHAR(100)          NOT
  OrganizationID  NVARCHAR(36)          NOT
  VulnerabilityID NVARCHAR(36)          NOT
  DeviceID        NVARCHAR(36)          NOT
  Port            VARCHAR(100)          NOT
  DueDate         DATETIME              NULL
*/

DROP PROCEDURE IF EXISTS `GetExceptionsByOrg`;

CREATE PROCEDURE `GetExceptionsByOrg` (_OrgID VARCHAR(36))
    #BEGIN#
SELECT
    ID,
    TypeID,
    Approval,
    OrganizationID,
    VulnerabilityID,
    DeviceID,
    Port,
    DueDate
FROM `Ignore` O
WHERE O.OrganizationID = _OrgID AND O.Active = b'1';