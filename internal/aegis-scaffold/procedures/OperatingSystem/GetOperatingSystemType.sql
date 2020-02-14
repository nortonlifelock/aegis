/*
  RETURN OperatingSystemType SINGLE
  ID       INT               NOT
  Type     VARCHAR(300)      NOT
  Match    VARCHAR(300)      NOT
  Priority INT               NOT
*/

DROP PROCEDURE IF EXISTS `GetOperatingSystemType`;

CREATE PROCEDURE `GetOperatingSystemType` (_OS VARCHAR(300))
  #BEGIN#
  SELECT
    O.ID,
    O.Type,
    O.Match,
    O.Priority
  FROM OperatingSystemType O
  WHERE INSTR(_OS, O.`Match`) > 0 ORDER BY Priority DESC LIMIT 1;