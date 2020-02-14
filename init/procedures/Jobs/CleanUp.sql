DROP PROCEDURE IF EXISTS `CleanUp`;

CREATE PROCEDURE `CleanUp` ()
    #BEGIN#

UPDATE JobHistory
SET StatusId = 4
WHERE (StatusId = 2 OR StatusId < 0)
   OR (Id IN (SELECT Id FROM
    (
        SELECT JH.Id
        FROM JobHistory JH
                 INNER JOIN JobConfig JC ON JC.Id = JH.ConfigId
        WHERE JC.AutoStart = 1
    ) AS JHInfo
)
    AND StatusId = 1);