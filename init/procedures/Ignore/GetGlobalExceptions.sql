/*
  RETURN Ignore
  ID              NVARCHAR(36)          NOT
  OrganizationID  NVARCHAR(36)          NOT
  VulnerabilityID NVARCHAR(36)          NOT
  DeviceID        NVARCHAR(36)          NOT
  OSRegex         VARCHAR(200)          NULL
  DueDate         DATETIME              NULL
*/

DROP PROCEDURE IF EXISTS `GetGlobalExceptions`;

CREATE PROCEDURE `GetGlobalExceptions` (_OrgID VARCHAR(36))
    #BEGIN#
SELECT
    ID,
    OrganizationID,
    VulnerabilityID,
    DeviceID,
    OSRegex,
    DueDate
FROM `Ignore` I
WHERE I.OrganizationID = _OrgID AND I.OSRegex IS NOT NULL AND I.DeviceID IS NULL AND I.Active = b'1'=;