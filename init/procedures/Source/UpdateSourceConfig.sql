DROP PROCEDURE IF EXISTS `UpdateSourceConfig`;

CREATE PROCEDURE `UpdateSourceConfig` (_ID VARCHAR(36), _OrgID NVARCHAR(36), _Address TEXT, _Username TEXT, _Password TEXT, _PrivateKey TEXT, _ConsumerKey TEXT, _Token TEXT, _Port NVARCHAR(10), _Payload MEDIUMTEXT, _UpdatedBy TEXT)
  #BEGIN#
  UPDATE SourceConfig SET
    Address = _Address,
    AuthInfo = JSON_OBJECT(
        "Username",
        _Username,
        "Password",
        _Password,
        "PrivateKey",
        _PrivateKey,
        "ConsumerKey",
        _ConsumerKey,
        "Token",
        _Token
    ),
    Port = _Port,
    Payload = _Payload,
    DBUpdatedDate = NOW(),
    UpdatedBy = _UpdatedBy
  WHERE ID = _ID AND OrganizationID = _OrgID;
