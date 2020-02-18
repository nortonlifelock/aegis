/*
  RETURN JobRegistration SINGLE
  ID                      INT             NOT
  GoStruct                NVARCHAR(50)    NOT
  Priority                INT             NOT
  CreatedDate             DATETIME        NOT
  CreatedBy               NVARCHAR(255)   NOT
  UpdatedDate             DATETIME        NULL
  UpdatedBy               NVARCHAR(255)   NULL
*/

DROP PROCEDURE IF EXISTS `GetJobByID`;

CREATE PROCEDURE `GetJobByID` (_ID INT)
  #BEGIN#
  SELECT
    J.Id,
    J.Struct AS GoStruct,
    J.Priority,
    J.CreatedDate,
    J.CreatedBy,
    J.UpdatedDate,
    J.UpdatedBy
  FROM Job J
  WHERE J.Id = _ID;