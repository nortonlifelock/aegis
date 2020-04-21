/*
    RETURN AssetGroup
    GroupID               VARCHAR(300)   NOT NULL
    OrganizationID        VARCHAR(36)    NOT NULL
    ScannerSourceConfigID VARCHAR(36)    NULL
    ScannerSourceID       VARCHAR(36)    NOT NULL
    CloudSourceID         VARCHAR(36)    NULL
    LastTicketing         DATETIME       NULL
    RescanQueueSkip       BOOL           NOT NULL
*/

DROP PROCEDURE IF EXISTS `GetAssetGroupForOrg`;

CREATE PROCEDURE `GetAssetGroupForOrg` (inScannerSourceConfigID VARCHAR(36), inOrgID VARCHAR(36))
    #BEGIN#
SELECT
    AG.GroupID,
    AG.OrganizationID,
    AG.ScannerSourceConfigID,
    AG.ScannerSourceID,
    AG.CloudSourceID,
    AG.LastTicketing,
    AG.RescanQueueSkip
FROM AssetGroup AG
WHERE AG.OrganizationID = inOrgID AND AG.ScannerSourceConfigID = inScannerSourceConfigID;