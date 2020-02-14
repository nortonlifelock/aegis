DROP PROCEDURE IF EXISTS `DeleteUserByUsername`;

CREATE PROCEDURE `DeleteUserByUsername` (_Username TEXT)
  #BEGIN#
  UPDATE Users
    SET IsDisabled = b'1'
  WHERE Username = _Username and Id > 0;