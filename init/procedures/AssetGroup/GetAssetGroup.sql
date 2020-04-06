/*
    RETURN AssetGroup SINGLE
    GroupID         VARCHAR(300)         NOT NULL
    ScannerSourceID VARCHAR(36) NOT NULL
    CloudSourceID   VARCHAR(36) NULL
    ScannerSourceConfigID VARCHAR(36) NULL
*/

DROP PROCEDURE IF EXISTS `GetAssetGroup`;

CREATE PROCEDURE `GetAssetGroup` (inOrgID VARCHAR(36), _GroupID VARCHAR(300), _ScannerConfigSourceID NVARCHAR(36))
    #BEGIN#
SELECT
    AG.GroupID,
    AG.ScannerSourceID,
    AG.CloudSourceID,
    AG.ScannerSourceConfigID
FROM AssetGroup AG
WHERE AG.GroupID = _GroupID AND AG.ScannerSourceConfigID = _ScannerConfigSourceID AND AG.OrganizationID = inOrgID;