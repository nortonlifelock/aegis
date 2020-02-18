/*
  RETURN Ignore SINGLE
  ID              NVARCHAR(36)          NOT
  OrganizationID  NVARCHAR(36)          NOT
  VulnerabilityID NVARCHAR(36)          NOT
  DeviceID        NVARCHAR(36)          NOT
  DueDate         DATETIME              NULL
*/

DROP PROCEDURE IF EXISTS `GetExceptionByVulnIDOrg`;

CREATE PROCEDURE `GetExceptionByVulnIDOrg` (_DeviceID VARCHAR(100), _VulnID NVARCHAR(255), _OrgID VARCHAR(36))
  #BEGIN#
  SELECT
    ID,
    OrganizationID,
    VulnerabilityID,
    DeviceID,
    DueDate
  FROM `Ignore` O
    WHERE O.VulnerabilityID = _VulnID AND O.OrganizationID = _OrgID AND O.DeviceId = _DeviceID;