/*
  RETURN Organization
  ID              NVARCHAR(36)           NOT
  Code            NVARCHAR(150) NOT
  Description     NVARCHAR(500) NULL
  TimeZoneOffset  FLOAT         NOT
*/

DROP PROCEDURE IF EXISTS `GetLeafOrganizationsForUser`;

CREATE PROCEDURE `GetLeafOrganizationsForUser` (_UserID VARCHAR(36))
  #BEGIN#

SELECT
  O.Id,
  O.Code,
  O.Description,
  O.TimeZoneOffset
FROM Users U
  JOIN Permissions P ON U.Id = P.UserId
  JOIN Organization O ON P.OrgId = O.Id
  LEFT JOIN `Organization` Child ON Child.ParentOrgID = O.ID -- helps us only grab leaf organizations
    WHERE  Child.Code IS NULL AND UserId = _UserID;