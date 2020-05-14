INSERT INTO SourceType(Type) VALUE ('ImageScanner');
INSERT INTO Source(SourceTypeID, Source) SELECT ID, 'Aqua' from SourceType where Type = 'ImageScanner';

SELECT S.Id
INTO @scanner_id
FROM SourceType S
WHERE S.Type = 'ImageScanner' LIMIT 1;

SELECT S.Id
INTO @ticket_id
FROM SourceType S
WHERE S.Type = 'TicketEngine' LIMIT 1;

INSERT INTO Job(Struct, SourceTypeIn, SourceTypeOut, Priority, CreatedBy) VALUE ('ImageRescanJob', @scanner_Id, @ticket_id, 4, 'ryan.everhart@nortonlifelock.com');