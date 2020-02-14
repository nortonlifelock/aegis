/*
  RETURN Permission SINGLE
  UserID         VARCHAR(36)    NOT
  OrgID          VARCHAR(36)    NOT
  Admin          BIT            NOT
  Manager        BIT            NOT
  Reader         BIT            NOT
  Reporter       BIT            NOT
*/

DROP PROCEDURE IF EXISTS `GetPermissionByUserOrgID`;

CREATE PROCEDURE `GetPermissionByUserOrgID` (_UserID VARCHAR(36), _OrgID VARCHAR(36))
  #BEGIN#
  SELECT
    UserId,
    OrgId,
    Admin,
    Manager,
    Reader,
    Reporter
  FROM Permissions P
  WHERE P.UserId = _UserID AND P.OrgId = _OrgID