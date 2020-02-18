/*
  RETURN User SINGLE
  ID             VARCHAR(36)    NOT
  Username       TEXT           NULL
  FirstName      TEXT           NOT
  LastName       TEXT           NOT
  Email          NVARCHAR(255)  NOT
  IsDisabled     BIT            NOT
*/

DROP PROCEDURE IF EXISTS `GetUserByUsername`;

CREATE PROCEDURE `GetUserByUsername` (_Username TEXT)
  #BEGIN#
  Select
    U.Id,
    U.Username,
    U.FirstName,
    U.LastName,
    U.Email,
    U.IsDisabled
  FROM Users U
  WHERE U.Username = _Username;