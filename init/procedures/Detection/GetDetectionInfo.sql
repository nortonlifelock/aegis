/*
  RETURN DetectionInfo SINGLE
  ID                  NVARCHAR(36)  NOT
  OrganizationID      NVARCHAR(36)  NOT
  SourceID            NVARCHAR(36)  NOT
  DeviceID            NVARCHAR(36)  NOT
  VulnerabilityID     NVARCHAR(36)  NOT
  IgnoreID            VARCHAR(36)   NULL
  AlertDate           DATETIME      NOT
  LastFound           DATETIME      NULL
  LastUpdated         DATETIME      NULL
  Proof               NVARCHAR(255) NOT
  Port                INT           NOT
  Protocol            VARCHAR(20)   NOT
  ActiveKernel        INT           NULL
  DetectionStatusID   INT           NOT
  TimesSeen           INT           NOT
  Updated             DATETIME      NOT
  ParentDetectionID   VARCHAR(36)   NULL
*/

DROP PROCEDURE IF EXISTS `GetDetectionInfo`;

CREATE PROCEDURE `GetDetectionInfo` (_DeviceID VARCHAR(360), _VulnerabilityID VARCHAR(36), _Port INT, _Protocol VARCHAR(360))
    #BEGIN#
SELECT
    D.ID,
    D.OrganizationID,
    D.SourceID,
    D.DeviceID,
    D.VulnerabilityID,
    D.IgnoreID,
    D.AlertDate,
    D.LastFound,
    D.LastUpdated,
    D.Proof,
    D.Port,
    D.Protocol,
    D.ActiveKernel,
    D.DetectionStatusID,
    D.TimesSeen,
    D.Updated,
    D.ParentDetectionId
FROM Detection D
WHERE D.DeviceID = _DeviceID AND D.VulnerabilityID = _VulnerabilityID AND D.Port = _Port AND D.Protocol = _Protocol;