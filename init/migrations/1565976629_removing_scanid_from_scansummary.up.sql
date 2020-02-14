SET SQL_SAFE_UPDATES = 0;
UPDATE ScanSummary S SET S.SourceKey = S.ScanIdentifier WHERE S.SourceKey = '';
ALTER TABLE ScanSummary DROP PRIMARY KEY;
ALTER TABLE ScanSummary DROP COLUMN ScanIdentifier;
ALTER TABLE ScanSummary ADD PRIMARY KEY(SourceID, OrgID, SourceKey);
SET SQL_SAFE_UPDATES = 1;