INSERT INTO `Job` (Struct, Priority, CreatedBy) VALUE ('CodeRescanJob', 5, 'ryan_everhart@symantec.com');

INSERT INTO `SourceType` (`Type`) VALUE ('CodeScanner');

SELECT Id
INTO @st_id
FROM SourceType ST
WHERE ST.Type = 'CodeScanner' LIMIT 1;

INSERT INTO Source (SourceTypeId, Source) VALUE (@st_id, 'Black Duck');