/*
  RETURN QueryData SINGLE
  Length                      INT           NOT
*/
DROP PROCEDURE IF EXISTS `GetExceptionsLength`;

CREATE PROCEDURE `GetExceptionsLength`(inSourceID VARCHAR(36), inOrgID VARCHAR(36), inTypeID INT, inVulnID NVARCHAR(255), inDevID VARCHAR(36), inDueDate DATETIME, inPort NVARCHAR(15), inApproval  NVARCHAR(120),
                                               inActive BIT, inDBCreatedDate DATETIME, inDBUpdatedDate DATETIME, inUpdatedBy NVARCHAR(255), inCreatedBy NVARCHAR(255))
    #BEGIN#
SELECT
    count(*)
FROM `Ignore` I
WHERE I.OrganizationId = inOrgID
  AND (I.SourceId = inSourceID OR inSourceID = '' OR inSourceID is NULL)
  AND (I.TypeId = inTypeID OR inTypeID= 0 OR inTypeID is NULL)
  AND (I.VulnerabilityId = inVulnID OR inVulnID = '' OR inVulnID is NULL)
  AND (I.DeviceId = inDevID OR inDevID ='' OR inDevID is NULL)
  AND (I.DueDate = inDueDate OR inDueDate ='1970-01-02 00:00:00 +0000 UTC' OR inDueDate is NULL)
  AND (I.Port = inPort OR inPort ='' OR inPort is NULL)
  AND (I.Approval= inApproval OR inApproval ='' OR inApproval is NULL)
  AND (I.Active = inActive OR inActive ='' OR inActive is NULL)
  AND (I.UpdatedBy = inUpdatedBy OR inUpdatedBy ='' OR inUpdatedBy is NULL)
  AND (I.CreatedBy = inCreatedBy OR inCreatedBy ='' OR inCreatedBy is NULL)
  AND (I.DBCreatedDate = inDBCreatedDate OR inDBCreatedDate ='1970-01-02 00:00:00 +0000 UTC' OR inDBCreatedDate is NULL)
  AND (I.DBUpdatedDate = inDBUpdatedDate OR inDBUpdatedDate ='1970-01-02 00:00:00 +0000 UTC' OR inDBUpdatedDate is NULL)