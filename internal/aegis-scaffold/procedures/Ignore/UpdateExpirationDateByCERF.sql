DROP PROCEDURE IF EXISTS `UpdateExpirationDateByCERF`;

CREATE PROCEDURE `UpdateExpirationDateByCERF` (_CERForm NVARCHAR(25), _OrganizationID NVARCHAR(36), _DueDate DATETIME)
  #BEGIN#

  UPDATE `Ignore`
  SET DueDate = _DueDate,
      DBUpdatedDate = NOW()
  WHERE OrganizationId = _OrganizationId
        AND Approval = _CERForm
        AND DueDate <> _DueDate;