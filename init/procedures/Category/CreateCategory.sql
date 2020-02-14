DROP PROCEDURE IF EXISTS `CreateCategory`;

CREATE PROCEDURE `CreateCategory`(_Category NVARCHAR(50))
  #BEGIN#
  INSERT INTO Category (Category) VALUE (_Category);