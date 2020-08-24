/*
  RETURN Ignore
  ID              NVARCHAR(36)          NOT
  OrganizationID  NVARCHAR(36)          NOT
  VulnerabilityID NVARCHAR(36)          NOT
  OSRegex         VARCHAR(200)          NULL
  HostnameRegex   VARCHAR(100)          NULL
  DueDate         DATETIME              NULL
*/

DROP PROCEDURE IF EXISTS `GetGlobalExceptions`;

CREATE PROCEDURE `GetGlobalExceptions` (_OrgID VARCHAR(36))
    #BEGIN#
SELECT
    ID,
    OrganizationID,
    VulnerabilityID,
    OSRegex,
    HostnameRegex,
    DueDate
FROM `Ignore` I
WHERE
      I.OrganizationID = _OrgID AND
      (I.OSRegex IS NOT NULL OR I.HostnameRegex IS NOT NULL) AND
      I.DeviceID = '' AND
      I.Active = b'1';