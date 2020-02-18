DROP PROCEDURE IF EXISTS `CreateException`;

CREATE PROCEDURE `CreateException` (inSourceID VARCHAR(36), inOrganizationID VARCHAR(36), inTypeID INT, inVulnerabilityID NVARCHAR(120), inDeviceID VARCHAR(36), inDueDate DATETIME, inApproval NVARCHAR(120), inActive BIT, inPort NVARCHAR(15), inCreatedBy NVARCHAR(255))
    #BEGIN#

INSERT INTO `Ignore` (SourceId, OrganizationId, TypeId, VulnerabilityId, DeviceId, DueDate, Approval, Active, Port, CreatedBy, DBCreatedDate)
    VALUE (inSourceID, inOrganizationID, inTypeID, inVulnerabilityID, inDeviceID, inDueDate, inApproval, inActive, inPort, inCreatedBy, NOW())
ON DUPLICATE KEY UPDATE TypeId = inTypeID, DueDate = inDueDate, Approval = inApproval, Active = inActive, UpdatedBy = inCreatedBy,DBUpdatedDate = NOW();