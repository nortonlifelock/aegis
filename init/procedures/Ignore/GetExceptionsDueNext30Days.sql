/*
  RETURN CERF
  CERForm  NVARCHAR(120)   NOT
 */

DROP PROCEDURE IF EXISTS `GetExceptionsDueNext30Days`;

CREATE PROCEDURE `GetExceptionsDueNext30Days`()
  #BEGIN#
  SELECT DISTINCT
    Approval
  FROM `Ignore`
  WHERE DueDate <= DATE_ADD(NOW(), INTERVAL 30 DAY)
        AND TypeId = 0
        AND Active = 1;