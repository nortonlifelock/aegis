DROP PROCEDURE IF EXISTS `DeleteSessionByToken`;

CREATE PROCEDURE `DeleteSessionByToken` (_SessionKey TEXT)
  #BEGIN#
  DELETE FROM UserSession WHERE SessionKey = _SessionKey;