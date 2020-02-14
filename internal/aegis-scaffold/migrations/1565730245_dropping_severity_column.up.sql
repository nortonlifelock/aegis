ALTER TABLE VulnerabilityInfoAudit DROP COLUMN Severity;
ALTER TABLE VulnerabilityInfoAudit DROP COLUMN CVSSVector;
ALTER TABLE VulnerabilityInfoAudit DROP COLUMN CVSS3Vector;

DROP TRIGGER VulnerabilityInfoAuditCreateTrigger;
DROP TRIGGER VulnerabilityInfoAuditUpdateTrigger;

CREATE TRIGGER VulnerabilityInfoAuditCreateTrigger BEFORE INSERT ON `VulnerabilityInfo`
    FOR EACH ROW
    INSERT INTO `VulnerabilityInfoAudit` select new.ID, new.SourceVulnId, new.Title, new.VulnerabilityID, new.SourceID, new.CVSS, new.CVSS3, new.Description, new.Solution, new.MatchConfidence, new.MatchReasons, new.Software, new.DetectionInformation, new.Updated, new.Created, 'INSERT', NEW.Created;

CREATE TRIGGER VulnerabilityInfoAuditUpdateTrigger AFTER UPDATE ON `VulnerabilityInfo`
    FOR EACH ROW
    INSERT INTO `VulnerabilityInfoAudit` SELECT new.ID, new.SourceVulnId, new.Title, new.VulnerabilityID, new.SourceID, new.CVSS, new.CVSS3, new.Description, new.Solution, new.MatchConfidence, new.MatchReasons, new.Software, new.DetectionInformation, new.Updated, new.Created, 'UPDATE', NOW();


ALTER TABLE VulnerabilityInfo MODIFY COLUMN CVSS3 FLOAT NULL;
ALTER TABLE VulnerabilityInfoAudit MODIFY COLUMN CVSS3 FLOAT NULL;

ALTER TABLE VulnerabilityInfo ADD KEY (SourceVulnID);
ALTER TABLE Vulnerability_Reference DROP FOREIGN KEY `fk_vm_vr_vi`;
ALTER TABLE Vulnerability_Reference MODIFY COLUMN VulnInfoID VARCHAR(255) NOT NULL;
ALTER TABLE `Vulnerability_Reference` ADD CONSTRAINT `fk_vm_vr_vi` FOREIGN KEY (`VulnInfoID`) REFERENCES `VulnerabilityInfo` (`SourceVulnId`);

ALTER TABLE Device ADD KEY(AssetID);
ALTER TABLE Detection DROP FOREIGN KEY `fk_vm_det_d`;
ALTER TABLE Device MODIFY COLUMN AssetID VARCHAR(36);
ALTER TABLE `Detection` ADD CONSTRAINT `fk_vm_det_d` FOREIGN KEY (`DeviceID`) REFERENCES `Device` (`AssetId`);