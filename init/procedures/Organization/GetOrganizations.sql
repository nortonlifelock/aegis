/*
  RETURN Organization
  ID              NVARCHAR(36)  NOT
  ParentOrgID     NVARCHAR(36)  NULL
  Code            NVARCHAR(150) NOT
  Description     NVARCHAR(500) NULL
  TimeZoneOffset  FLOAT         NOT
  Created         DATETIME      NOT
  Updated         DATETIME      NULL
  Payload         TEXT          NOT
*/

DROP PROCEDURE IF EXISTS `GetOrganizations`;

CREATE PROCEDURE `GetOrganizations` ()
#BEGIN#
  SELECT
    O.Id,
    O.ParentOrgID,
    O.Code,
    O.Description,
    O.TimeZoneOffset,
    O.Created,
    O.Updated,
       O.Payload
  FROM Organization O;
