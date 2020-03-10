DROP PROCEDURE IF EXISTS `CreateDetection`;

CREATE PROCEDURE `CreateDetection`(_OrgID NVARCHAR(36), _SourceID NVARCHAR(36), _DeviceID NVARCHAR(36), _VulnID NVARCHAR(36), _IgnoreID NVARCHAR(36), _AlertDate DATETIME, _LastFound DATETIME, _LastUpdated DATETIME, _Proof MEDIUMTEXT, _Port INT, _Protocol VARCHAR(20), _ActiveKernel INT, _DetectionStatusID INT, _TimesSeen INT, _DefaultTime DATETIME)
    #BEGIN#
INSERT INTO Detection(OrganizationID, SourceID, DeviceID, VulnerabilityID, AlertDate, LastFound, LastUpdated, Proof, Port, Protocol, ActiveKernel, DetectionStatusID, TimesSeen, IgnoreID)
    VALUE (_OrgID, _SourceID, _DeviceID, _VulnID, _AlertDate, NULLIF(_LastFound, _DefaultTime), NULLIF(_LastUpdated, _DefaultTime), _Proof, _Port, _Protocol, NULLIF(_ActiveKernel, 0), _DetectionStatusID, _TimesSeen, NULLIF(_IgnoreID, ''));