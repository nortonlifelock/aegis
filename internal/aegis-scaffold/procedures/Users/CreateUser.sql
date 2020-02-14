DROP PROCEDURE IF EXISTS `CreateUser`;

CREATE PROCEDURE `CreateUser` (_Username TEXT, _FirstName TEXT, _LastName TEXT, _Email TEXT)
  #BEGIN#
  INSERT INTO Users (Username, FirstName, LastName, Email)
    VALUE (_Username, _FirstName, _LastName, _Email);