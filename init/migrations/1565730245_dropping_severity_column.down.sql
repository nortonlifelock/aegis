ALTER TABLE VulnerabilityInfo ADD COLUMN Severity INT NOT NULL AFTER Solution;
ALTER TABLE VulnerabilityInfo ADD COLUMN CVSSVector TEXT NOT NULL AFTER Severity;
ALTER TABLE VulnerabilityInfo ADD COLUMN CVSS3Vector TEXT NOT NULL AFTER CVSSVector;

ALTER TABLE VulnerabilityInfoAudit ADD COLUMN Severity INT NOT NULL AFTER Solution;
ALTER TABLE VulnerabilityInfoAudit ADD COLUMN CVSSVector TEXT NOT NULL AFTER Severity;
ALTER TABLE VulnerabilityInfoAudit ADD COLUMN CVSS3Vector TEXT NOT NULL AFTER CVSSVector;

DROP TRIGGER VulnerabilityInfoAuditCreateTrigger;
DROP TRIGGER VulnerabilityInfoAuditUpdateTrigger;

CREATE TRIGGER VulnerabilityInfoAuditCreateTrigger BEFORE INSERT ON `VulnerabilityInfo`
    FOR EACH ROW
    INSERT INTO `VulnerabilityInfoAudit` select new.ID, new.SourceVulnId, new.Title, new.VulnerabilityID, new.SourceID, new.CVSS, new.CVSS3, new.Description, new.Solution, new.Severity, new.CVSSVector, new.CVSS3Vector, new.MatchConfidence, new.MatchReasons, new.Software, new.DetectionInformation, new.Updated, new.Created, 'INSERT', NEW.Created;

CREATE TRIGGER VulnerabilityInfoAuditUpdateTrigger AFTER UPDATE ON `VulnerabilityInfo`
    FOR EACH ROW
    INSERT INTO `VulnerabilityInfoAudit` SELECT new.ID, new.SourceVulnId, new.Title, new.VulnerabilityID, new.SourceID, new.CVSS, new.CVSS3, new.Description, new.Solution, new.Severity, new.CVSSVector, new.CVSS3Vector, new.MatchConfidence, new.MatchReasons, new.Software, new.DetectionInformation, new.Updated, new.Created, 'UPDATE', NOW();

ALTER TABLE VulnerabilityInfo MODIFY COLUMN CVSS3 FLOAT NOT NULL;
ALTER TABLE VulnerabilityInfoAudit MODIFY COLUMN CVSS3 FLOAT NOT NULL;