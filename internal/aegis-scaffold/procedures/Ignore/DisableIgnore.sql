DROP PROCEDURE IF EXISTS `DisableIgnore`;

CREATE PROCEDURE `DisableIgnore`(inSourceID VARCHAR(36), inDevID VARCHAR(36), inOrgID VARCHAR(36), inVulnID NVARCHAR(255), inPortID NVARCHAR(15), inUpdatedBy NVARCHAR(255))
    #BEGIN#
BEGIN
    UPDATE `Ignore` SET Active = b'0', DBUpdatedDate = NOW(), UpdatedBy = inUpdatedBy WHERE SourceId = inSourceID AND DeviceId = inDevID AND OrganizationId = inOrgID AND VulnerabilityId = inVulnID AND Port = inPortID;
END