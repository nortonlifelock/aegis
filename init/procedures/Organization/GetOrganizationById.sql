/*
  RETURN Organization SINGLE
  ID              NVARCHAR(36)  NOT
  ParentOrgID     VARCHAR(36)   NULL
  Code            NVARCHAR(150) NOT
  Description     NVARCHAR(500) NULL
  TimeZoneOffset  FLOAT         NOT
  Created         DATETIME      NOT
  Updated         DATETIME      NULL
  Payload         TEXT          NOT
  EncryptionKey   TEXT          NULL
*/

DROP PROCEDURE IF EXISTS `GetOrganizationByID`;

CREATE PROCEDURE `GetOrganizationByID` (ID NVARCHAR(36))
#BEGIN#
  SELECT
    O.Id,
    O.ParentOrgId,
    O.Code,
    O.Description,
    O.TimeZoneOffset,
    O.Created,
    O.Updated,
    O.Payload,
    O.EncryptionKey
  FROM Organization O
  WHERE O.Id = ID;
