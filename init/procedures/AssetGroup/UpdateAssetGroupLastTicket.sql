DROP PROCEDURE IF EXISTS `UpdateAssetGroupLastTicket`;

CREATE PROCEDURE `UpdateAssetGroupLastTicket` (inGroupID VARCHAR(36), inOrgID VARCHAR(36), inLastTicketTime DATETIME)
    #BEGIN#
UPDATE AssetGroup AG SET AG.LastTicketing = inLastTicketTime WHERE AG.OrganizationID = inOrgID AND (inGroupID = '' OR AG.GroupID = inGroupID);