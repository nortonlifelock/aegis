DROP PROCEDURE IF EXISTS `CreateUserPermissions`;

CREATE PROCEDURE `CreateUserPermissions` (_UserID VARCHAR(36), _OrgID VARCHAR(36))
  #BEGIN#
  INSERT INTO Permissions(UserID, OrgID) VALUE (_UserID, _OrgID);