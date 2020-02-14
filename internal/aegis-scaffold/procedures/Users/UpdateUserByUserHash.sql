DROP PROCEDURE IF EXISTS `UpdateUserByID`;

CREATE PROCEDURE `UpdateUserByID`(_ID VARCHAR(36), _FirstName TEXT, _LastName TEXT, _Email TEXT, _Disabled BIT)
  #BEGIN#
  BEGIN

    DECLARE FN TEXT;
    DECLARE LN TEXT;
    DECLARE EM TEXT;
#     DECLARE UN TEXT;

    IF _FirstName != ''
        THEN
          SET FN = _FirstName;
        ELSE
          SET FN = (SELECT FirstName from Users WHERE ID = _ID LIMIT 1);
    END IF;

    IF _LastName != ''
      THEN
        SET LN = _LastName;
      ELSE
        SET LN = (SELECT LastName from Users WHERE ID = _ID LIMIT 1);
    END IF;

    IF _Email != ''
      THEN
        SET EM = _Email;
      ELSE
        SET EM = (SELECT Email from Users WHERE ID = _ID LIMIT 1);
    END IF;

    UPDATE Users SET FirstName = FN, LastName = LN, Email = EM, IsDisabled = _Disabled WHERE ID = _ID;

  END