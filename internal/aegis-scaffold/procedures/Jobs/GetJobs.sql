/*
  RETURN JobRegistration
  ID                      INT             NOT
  GoStruct                TEXT            NOT
  Priority                INT             NOT
*/

DROP PROCEDURE IF EXISTS `GetJobs`;

CREATE PROCEDURE `GetJobs` ()
  #BEGIN#
  SELECT
    J.Id,
    J.Struct,
    J.Priority
  FROM Job J;