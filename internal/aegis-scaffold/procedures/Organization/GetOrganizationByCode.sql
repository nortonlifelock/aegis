/*
  RETURN Organization SINGLE
  ID              NVARCHAR(36)           NOT
  Code            NVARCHAR(150) NOT
  Description     NVARCHAR(500) NULL
  TimeZoneOffset  FLOAT         NOT
  Created         DATETIME      NOT
  Updated         DATETIME      NULL
  Payload         TEXT          NOT

*/

DROP PROCEDURE IF EXISTS `GetOrganizationByCode`;

CREATE PROCEDURE `GetOrganizationByCode` (Code NVARCHAR(20))
#BEGIN#
  SELECT
    O.Id,
    O.Code,
    O.Description,
    O.TimeZoneOffset,
    O.Created,
    O.Updated,
    O.Payload
  FROM Organization O
  WHERE O.Code = Code;
