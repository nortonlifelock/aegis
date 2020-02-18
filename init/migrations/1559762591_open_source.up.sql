INSERT INTO Job
(Struct, Priority, CreatedBy)
    VALUE
    ('RescanQueueJob', 5, 'benjamin_vesterby@symantec.com'),
    ('RescanJob', 4, 'benjamin_vesterby@symantec.com'),
    ('ExceptionJob', 5, 'benjamin_vesterby@symantec.com'),
    ('TicketingJob', 4, 'benjamin_vesterby@symantec.com'),
    ('ScanSyncJob', 5, 'benjamin_vesterby@symantec.com'),
    ('ScanCloseJob', 5, 'ryan_everhart@symantec.com'),

    ('BulkUpdateJob', 5, 'ryan_everhart@symantec.com'),
    ('VulnSyncJob', 4, 'ryan_everhart@symantec.com'),
    ('CloudSyncJob', 5, 'ryan_everhart@symantec.com'),
    ('AssetSyncJob', 4, 'ryan_everhart@symantec.com');

INSERT INTO `SourceType`
    (`Type`)
    VALUE
    ('Scanner'),
    ('TicketEngine'),
    ('Cloud'),
    ('CISBench');

SELECT S.Id
       INTO @scanner_id
FROM SourceType S
WHERE S.Type = 'Scanner' LIMIT 1;

SELECT S.Id
       INTO @ticket_engine_id
FROM SourceType S
WHERE S.Type = 'TicketEngine' LIMIT 1;

SELECT Id
       INTO @cloud_id
FROM SourceType ST
WHERE ST.Type = 'Cloud' LIMIT 1;

SELECT Id
INTO @cis_id
FROM SourceType ST
WHERE ST.Type = 'CISBench' LIMIT 1;

INSERT INTO Source
(SourceTypeId, Source)
    VALUE
    (@ticket_engine_id, 'JIRA'),
    (@ticket_engine_id, 'ServiceNow'),
    (@scanner_id, 'Nexpose'),
    (@scanner_id, 'Qualys'),
    (@cis_id, 'Dome9'),
    (@cloud_id, 'AWS'),
    (@cloud_id, 'AZURE');

INSERT INTO LogType (Id, Type, Name) values (0, 'INFO', 'Information');
INSERT INTO LogType (Id, Type, Name) VALUES (1, 'WARN', 'Warning');
INSERT INTO LogType (Id, Type, Name) VALUES (2, 'ERR', 'Error');
INSERT INTO LogType (Id, Type, Name) VALUES (3, 'CRIT', 'Critical');
INSERT INTO LogType (Id, Type, Name) VALUES (4, 'FATL', 'Fatal');

INSERT INTO `JobStatus` (`Id`, `Status`)
    VALUES
        ('-1', 'CANCEL'),
        ('0', 'UNKNOWN'),
        ('1', 'PENDING'),
        ('2', 'IN PROGRESS'),
        ('3', 'COMPLETED'),
        ('4', 'ERROR'),
        ('5', 'CANCELED');

INSERT INTO DetectionStatus(Status, Name) VALUES ('new', 'New');
INSERT INTO DetectionStatus(Status, Name) VALUES ('active', 'Active');
INSERT INTO DetectionStatus(Status, Name) VALUES ('reopened', 'Re-Opened');
INSERT INTO DetectionStatus(Status, Name) VALUES ('fixed', 'Fixed');

INSERT INTO ReferenceType  (Id, Type, DBUpdatedDate)
VALUES (0, 'cve', NOW()),
       (1, 'ms', NOW()),
       (2, 'vendor', NOW());