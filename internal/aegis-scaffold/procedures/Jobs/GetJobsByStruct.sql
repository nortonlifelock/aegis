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

DROP PROCEDURE IF EXISTS `GetJobsByStruct`;

CREATE PROCEDURE `GetJobsByStruct` (_Struct NVARCHAR(50))
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
  WHERE J.Struct = _Struct;