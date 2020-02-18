/*
    RETURN AssetGroup
    GroupID         INT         NOT NULL
    OrganizationID  VARCHAR(36) NOT NULL
    ScannerSourceID VARCHAR(36) NOT NULL
    CloudSourceID   VARCHAR(36) NULL
    ScannerSourceConfigID VARCHAR(36) NULL
*/

DROP PROCEDURE IF EXISTS `GetAssetGroupsByCloudSource`;

CREATE PROCEDURE `GetAssetGroupsByCloudSource` (inOrgID VARCHAR(36), inCloudSourceID NVARCHAR(36))
    #BEGIN#
SELECT
    AG.GroupID,
    AG.OrganizationID,
    AG.ScannerSourceID,
    AG.CloudSourceID,
    AG.ScannerSourceConfigID
FROM AssetGroup AG
WHERE AG.CloudSourceID = inCloudSourceID AND AG.OrganizationID = inOrgID;