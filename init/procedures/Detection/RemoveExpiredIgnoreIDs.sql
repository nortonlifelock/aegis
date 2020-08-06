DROP PROCEDURE IF EXISTS `RemoveExpiredIgnoreIDs`;

CREATE PROCEDURE `RemoveExpiredIgnoreIDs` (_OrgID VARCHAR(36))
    #BEGIN#
UPDATE Detection SET IgnoreID = NULL
  WHERE OrganizationID = _OrgID
  AND IgnoreID IN (select ID from `Ignore` where TypeID = (select ID from IgnoreType where Type = 'EXCEPTION') and NOW() > DueDate)
  AND ID != '';