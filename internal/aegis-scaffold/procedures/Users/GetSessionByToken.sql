/*
  RETURN Session SINGLE
  UserID        VARCHAR(36)        NOT
  OrgID         VARCHAR(36)        NOT
  SessionKey    TEXT       NOT
  IsDisabled    BIT        NOT
*/

DROP PROCEDURE IF EXISTS `GetSessionByToken`;

CREATE PROCEDURE `GetSessionByToken`(_SessionKey TEXT)
  #BEGIN#
  SELECT
    US.UserId,
    US.OrgId,
    US.SessionKey,
    US.IsDisabled
  FROM UserSession US
  WHERE US.SessionKey = _SessionKey;