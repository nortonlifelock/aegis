SET SQL_SAFE_UPDATES = 0;
ALTER TABLE ScanSummary DROP PRIMARY KEY;

ALTER TABLE ScanSummary ADD COLUMN ScanIdentifier INT NOT NULL AFTER OrgID;
UPDATE ScanSummary S SET S.ScanIdentifier = S.SourceKey WHERE S.SourceKey NOT LIKE 'scan%';
UPDATE ScanSummary S SET S.SourceKey = '' WHERE S.SourceKey NOT LIKE 'scan%';

ALTER TABLE ScanSummary ADD PRIMARY KEY(SourceID, OrgID, ScanIdentifier, SourceKey);
SET SQL_SAFE_UPDATES = 1;


