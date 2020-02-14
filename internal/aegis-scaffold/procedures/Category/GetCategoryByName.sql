/*
  RETURN Category
  ID                      NVARCHAR(36)               NOT
  Category                NVARCHAR(50)               NOT
  ParentCategoryID        NVARCHAR(36)               NULL
*/

DROP PROCEDURE IF EXISTS `GetCategoryByName`;

CREATE PROCEDURE `GetCategoryByName` (_Name NVARCHAR(50))
  #BEGIN#
  SELECT
    C.Id,
    C.Category,
    C.ParentCategoryId
  FROM Category C
  WHERE C.Category = _Name;