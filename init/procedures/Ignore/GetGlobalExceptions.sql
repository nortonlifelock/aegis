/*
  RETURN Ignore
  ID              NVARCHAR(36)          NOT
  OrganizationID  NVARCHAR(36)          NOT
  VulnerabilityID NVARCHAR(36)          NOT
  OSRegex         VARCHAR(200)          NULL
  HostnameRegex   VARCHAR(100)          NULL
  DeviceIDRegex   VARCHAR(100)          NULL
  DueDate         DATETIME              NULL
*/

DROP PROCEDURE IF EXISTS `GetGlobalExceptions`;

CREATE PROCEDURE `GetGlobalExceptions` (_OrgID VARCHAR(36), _SourceID VARCHAR(100))
    #BEGIN#
SELECT
    ID,
    OrganizationID,
    VulnerabilityID,
    OSRegex,
    HostnameRegex,
    DeviceIDRegex,
    DueDate
FROM `Ignore` I
WHERE
      I.OrganizationID = _OrgID AND
      I.SourceID = _SourceID AND
      (I.OSRegex IS NOT NULL OR I.HostnameRegex IS NOT NULL or I.DeviceIDRegex IS NOT NULL) AND
      I.DeviceID = '' AND
      I.Active = b'1';