ALTER TABLE `VulnerabilityInfo` ADD COLUMN Severity int(11) NOT NULL;
ALTER TABLE `VulnerabilityInfo` ADD COLUMN CVSSVector text NOT NULL;
ALTER TABLE `VulnerabilityInfo` ADD COLUMN CVSS3Vector text NOT NULL;