/*
  RETURN SourceConfig SINGLE
  ID                NVARCHAR(36)            NOT
  SourceID          NVARCHAR(36)            NOT
  Source            NVARCHAR(100)           NOT
  Address           NVARCHAR(100)           NOT
  Port              NVARCHAR(30)            NOT
  AuthInfo          TEXT                    NOT
  Payload           NVARCHAR(1000)          NULL
  OrganizationID    NVARCHAR(36)            NOT
  DBCreatedDate     DATETIME                NOT
  DBUpdatedDate     DATETIME                NULL
 */

DROP PROCEDURE IF EXISTS `GetSourceOauthByOrgURL`;

CREATE PROCEDURE `GetSourceOauthByOrgURL` (_URL NVARCHAR(255), _OrgID NVARCHAR(36))
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
  WHERE SC.Address = _URL
    AND SC.OrganizationId = _OrgID
    AND JSON_EXTRACT(AuthInfo, '$.PrivateKey') IS NOT NULL IS NOT NULL
    AND JSON_EXTRACT(AuthInfo, '$.ConsumerKey') IS NOT NULL IS NOT NULL
    AND JSON_EXTRACT(AuthInfo, '$.Token') IS NOT NULL IS NOT NULL;