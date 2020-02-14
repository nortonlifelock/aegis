ALTER TABLE ScanSummary MODIFY COLUMN ScanClosePayload MEDIUMTEXT NOT NULL;
set sql_safe_updates = 0;
Update Source Set Source = 'Azure' where Source = 'AZURE';
set sql_safe_updates = 1;