/*
  RETURN User SINGLE
  ID             VARCHAR(36)            NOT
  Username       TEXT           NULL
  FirstName      TEXT           NOT
  LastName       TEXT           NOT
  Email          NVARCHAR(255)  NOT
  IsDisabled     BIT            NOT
*/

DROP PROCEDURE IF EXISTS `GetUserByID`;

CREATE PROCEDURE `GetUserByID` (_ID NVARCHAR(255), _OrgID VARCHAR(36))
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
  WHERE U.ID = _ID AND P.OrgId = _OrgID;