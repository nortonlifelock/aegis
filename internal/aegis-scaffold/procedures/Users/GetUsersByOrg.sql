/*
  RETURN User
  ID             VARCHAR(36)            NOT
  Username       TEXT           NULL
  FirstName      TEXT           NOT
  LastName       TEXT           NOT
  Email          NVARCHAR(255)  NOT
  IsDisabled     BIT            NOT
*/

DROP PROCEDURE IF EXISTS `GetUsersByOrg`;

CREATE PROCEDURE `GetUsersByOrg` (_OrgID VARCHAR(36))
  #BEGIN#
  Select
    U.Id,
    U.Username,
    U.FirstName,
    U.LastName,
    U.Email,
    U.IsDisabled
  FROM Users U
    JOIN Permissions P ON P.UserId = U.Id
  WHERE P.OrgId = _OrgID;