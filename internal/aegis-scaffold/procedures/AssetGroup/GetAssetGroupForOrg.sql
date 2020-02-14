/*
    RETURN AssetGroup
    GroupID         INT         NOT NULL
    ScannerSourceID VARCHAR(36) NOT NULL
    CloudSourceID   VARCHAR(36) NULL
    ScannerSourceConfigID VARCHAR(36) NULL
*/

DROP PROCEDURE IF EXISTS `GetAssetGroupForOrg`;

CREATE PROCEDURE `GetAssetGroupForOrg` (inScannerSourceConfigID VARCHAR(36), inOrgID VARCHAR(36))
    #BEGIN#
SELECT
    AG.GroupID,
    AG.ScannerSourceID,
    AG.CloudSourceID,
    AG.ScannerSourceConfigID
FROM AssetGroup AG
WHERE AG.OrganizationID = inOrgID AND AG.ScannerSourceConfigID = inScannerSourceConfigID;