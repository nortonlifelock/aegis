/*
  RETURN Permission SINGLE
  UserID         VARCHAR(36)    NOT
  OrgID          VARCHAR(36)    NOT
  Admin          BIT            NOT
  Manager        BIT            NOT
  Reader         BIT            NOT
  Reporter       BIT            NOT
*/

DROP PROCEDURE IF EXISTS `GetPermissionOfLeafOrgByUserID`;

CREATE PROCEDURE `GetPermissionOfLeafOrgByUserID` (_UserID VARCHAR(36))
  #BEGIN#
  Select
    UserId,
    OrgId,
    Admin,
    Manager,
    Reader,
    Reporter
  FROM Permissions P
    LEFT JOIN `Organization` O ON O.ID = P.OrgID LEFT JOIN `Organization` Child ON Child.ParentOrgID = O.ID
  WHERE Child.Code IS NULL AND P.UserId = _UserID;