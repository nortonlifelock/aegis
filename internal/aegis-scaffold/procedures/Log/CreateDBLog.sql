DROP PROCEDURE IF EXISTS `CreateDBLog`;

CREATE PROCEDURE `CreateDBLog` (_User TEXT, _Command TEXT, _Endpoint TEXT)
  #BEGIN#

  INSERT INTO DBLog (User, Command, Endpoint)
    VALUES(_User, _Command, _Endpoint);