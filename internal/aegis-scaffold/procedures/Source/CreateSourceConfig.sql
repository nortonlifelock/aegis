DROP PROCEDURE IF EXISTS `CreateSourceConfig`;

CREATE PROCEDURE `CreateSourceConfig` (_Source TEXT, _SourceID VARCHAR(36), _OrganizationID NVARCHAR(36), _Address TEXT,
                                                _Port NVARCHAR(10), _Username TEXT, _Password TEXT, _PrivateKey TEXT, _ConsumerKey TEXT, _Token TEXT, _Payload MEDIUMTEXT)
#BEGIN#
  INSERT INTO SourceConfig(Source, SourceID, OrganizationID, Address, Port, AuthInfo, Payload)
  VALUES (_Source,
          _SourceID,
          _OrganizationID,
          _Address,
          _Port,
          JSON_OBJECT(
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
          _Payload
  );