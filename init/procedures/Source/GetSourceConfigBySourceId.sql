/*
  RETURN SourceConfig
  ID                NVARCHAR(36)                     NOT
  SourceID          NVARCHAR(36)                     NOT
  OrganizationID    NVARCHAR(36)                     NOT
  Address           NVARCHAR(100)           NOT
  Port              NVARCHAR(30)            NOT
  AuthInfo          TEXT                    NOT
  Payload           NVARCHAR(1000)          NULL
  DBCreatedDate     DATETIME                NOT
  DBUpdatedDate     DATETIME                NULL
 */

DROP PROCEDURE IF EXISTS `GetSourceConfigBySourceID`;

CREATE PROCEDURE `GetSourceConfigBySourceID` (_OrgID NVARCHAR(36), _SourceID NVARCHAR(36))
  #BEGIN#
  SELECT
    SC.Id,
    SC.SourceId,
    SC.OrganizationId,
    SC.Address,
    SC.Port,
    SC.AuthInfo,
    SC.Payload,
    SC.DBCreatedDate,
    SC.DBUpdatedDate
  FROM SourceConfig SC
  WHERE SC.SourceId = _SourceID AND SC.OrganizationId = _OrgID;