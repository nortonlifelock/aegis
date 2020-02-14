SELECT S.Id
INTO @cloud_id
FROM SourceType S
WHERE S.Type = 'Cloud' LIMIT 1;

SELECT S.Id
INTO @ticket_id
FROM SourceType S
WHERE S.Type = 'TicketEngine' LIMIT 1;

INSERT INTO Job(Struct, SourceTypeIn, SourceTypeOut, Priority, CreatedBy) VALUES ('CloudDecommissionJob', @cloud_id, @ticket_id, 4, 'Ryan');