/*
  RETURN SourceConfig SINGLE
  ID                NVARCHAR(36)                     NOT
  SourceID          NVARCHAR(36)                     NOT
  Source            NVARCHAR(100)           NOT
  Address           NVARCHAR(100)           NOT
  Port              NVARCHAR(30)            NOT
  AuthInfo          TEXT                    NOT
  Payload           NVARCHAR(1000)          NULL
  OrganizationID    NVARCHAR(36)                     NOT
  DBCreatedDate     DATETIME                NOT
  DBUpdatedDate     DATETIME                NULL
 */

DROP PROCEDURE IF EXISTS `GetSourceConfigByID`;

CREATE PROCEDURE `GetSourceConfigByID` (_ID NVARCHAR(36))
#BEGIN#
  SELECT
    SC.Id,
    SC.SourceId,
    S.Source,
    SC.Address,
    SC.Port,
    SC.AuthInfo,
    SC.Payload,
    SC.OrganizationId,
    SC.DBCreatedDate,
    SC.DBUpdatedDate
  FROM SourceConfig SC
  JOIN Source S on S.Id = SC.SourceId
  WHERE SC.Id = _ID;