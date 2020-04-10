/*
    RETURN AssetGroup
    GroupID               VARCHAR(300) NOT NULL
    ScannerSourceID       VARCHAR(36)  NOT NULL
    CloudSourceID         VARCHAR(36)  NULL
    ScannerSourceConfigID VARCHAR(36)  NULL
    LastTicketing         DATETIME     NULL
*/

DROP PROCEDURE IF EXISTS `GetAssetGroupsForOrg`;

CREATE PROCEDURE `GetAssetGroupsForOrg` (inOrgID VARCHAR(36))
    #BEGIN#
SELECT
    AG.GroupID,
    AG.ScannerSourceID,
    AG.CloudSourceID,
    AG.ScannerSourceConfigID,
    AG.LastTicketing
FROM AssetGroup AG
WHERE AG.OrganizationID = inOrgID;