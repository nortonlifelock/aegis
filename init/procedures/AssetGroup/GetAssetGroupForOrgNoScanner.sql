/*
    RETURN AssetGroup
    GroupID               INT         NOT NULL
    ScannerSourceID       VARCHAR(36) NOT NULL
    CloudSourceID         VARCHAR(36) NULL
    ScannerSourceConfigID VARCHAR(36) NULL
*/

DROP PROCEDURE IF EXISTS `GetAssetGroupForOrgNoScanner`;

CREATE PROCEDURE `GetAssetGroupForOrgNoScanner` (inOrgID VARCHAR(36), inGroupID VARCHAR(100))
    #BEGIN#
SELECT
    AG.GroupID,
    AG.ScannerSourceID,
    AG.CloudSourceID,
    AG.ScannerSourceConfigID
FROM AssetGroup AG
WHERE AG.OrganizationID = inOrgID AND (inGroupID = '' OR AG.GroupID = inGroupID);