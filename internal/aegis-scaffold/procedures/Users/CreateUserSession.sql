DROP PROCEDURE IF EXISTS `CreateUserSession`;

CREATE PROCEDURE `CreateUserSession` (_UserID VARCHAR(36), _OrgID VARCHAR(36), _SessionKey TEXT)
  #BEGIN#
BEGIN
  DELETE FROM UserSession WHERE UserID = _UserID;

  INSERT INTO UserSession (UserID, OrgID, SessionKey)
    VALUE (_UserID, _OrgID, _SessionKey);
END