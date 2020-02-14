/*
  RETURN User SINGLE
  ID             VARCHAR(36)            NOT
  Username       TEXT           NULL
  FirstName      TEXT           NOT
  LastName       TEXT           NOT
  Email          NVARCHAR(255)  NOT
  IsDisabled     BIT            NOT
*/

DROP PROCEDURE IF EXISTS `GetUserAnyOrg`;

CREATE PROCEDURE `GetUserAnyOrg` (_ID VARCHAR(36))
  #BEGIN#
  Select
    U.Id,
    U.Username,
    U.FirstName,
    U.LastName,
    U.Email,
    U.IsDisabled
  FROM Users U
  WHERE U.Id = _ID;