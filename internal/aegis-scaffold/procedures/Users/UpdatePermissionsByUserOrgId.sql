DROP PROCEDURE IF EXISTS `UpdatePermissionsByUserOrgID`;

CREATE PROCEDURE `UpdatePermissionsByUserOrgID`(_UserID VARCHAR(36), _OrgID VARCHAR(36), _Admin BIT, _Manager BIT, _Reader BIT, _Reporter BIT)
  #BEGIN#
  UPDATE Permissions
    SET Admin = _Admin, Manager = _Manager, Reader = _Reader, Reporter = _Reporter
  WHERE UserID = _UserID AND OrgID = _OrgID