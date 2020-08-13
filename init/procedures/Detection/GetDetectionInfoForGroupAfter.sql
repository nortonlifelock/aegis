/*
  RETURN DetectionInfo
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
*/

DROP PROCEDURE IF EXISTS `GetDetectionInfoForGroupAfter`;

CREATE PROCEDURE `GetDetectionInfoForGroupAfter` (_LastUpdatedAfter DATETIME, _LastFoundAfter DATETIME, _OrgID VARCHAR(36), inGroupID VARCHAR(100), ticketInactiveKernels BOOL)
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
    D.Updated
FROM Detection D
         JOIN DetectionStatus DS on D.DetectionStatusId = DS.Id
         JOIN Device Dev ON Dev.AssetID = D.DeviceID
WHERE (ticketInactiveKernels OR D.ActiveKernel IS NULL OR D.ActiveKernel = 1) AND D.OrganizationID = _OrgID
  AND D.IgnoreID IS NULL AND (D.Updated > _LastUpdatedAfter OR D.Created > _LastUpdatedAfter)
  AND D.LastFound > _LastFoundAfter AND DS.Status != 'fixed' AND Dev.GroupID = inGroupID ORDER BY Dev.TrackingMethod, D.Created;
