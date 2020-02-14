DROP PROCEDURE IF EXISTS `CreateDetectionWithIgnoreActiveKernel`;

CREATE PROCEDURE `CreateDetectionWithIgnoreActiveKernel`(_OrgID NVARCHAR(36), _SourceID NVARCHAR(36), _DeviceID NVARCHAR(36), _VulnID NVARCHAR(36), _IgnoreID NVARCHAR(36), _AlertDate DATETIME, _Proof MEDIUMTEXT, _Port INT, _Protocol VARCHAR(20), _ActiveKernel INT, _DetectionStatusID INT, _TimesSeen INT)
    #BEGIN#
INSERT INTO Detection(OrganizationID, SourceID, DeviceID, VulnerabilityID, AlertDate, Proof, Port, Protocol, ActiveKernel, DetectionStatusID, TimesSeen, IgnoreID)
    VALUE (_OrgID, _SourceID, _DeviceID, _VulnID, _AlertDate, _Proof, _Port, _Protocol, _ActiveKernel, _DetectionStatusID, _TimesSeen, _IgnoreID);