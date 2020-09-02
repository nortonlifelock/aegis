CREATE TABLE CategoryRule(
    ID VARCHAR(36) NOT NULL,
    OrganizationID VARCHAR(36) NOT NULL,
    SourceID VARCHAR(36) NOT NULL,
    VulnerabilityTitle VARCHAR(300) NULL,
    VulnerabilityCategory VARCHAR(300) NULL,
    VulnerabilityType VARCHAR(100) NULL,
    Category VARCHAR(200) NOT NULL,

    PRIMARY KEY (`ID`),
    KEY `fk_vm_cr_o` (`OrganizationID`),
    KEY `fk_vm_cr_s` (`SourceID`),
    CONSTRAINT `fk_vm_cr_o` FOREIGN KEY (`OrganizationID`) REFERENCES `Organization` (`ID`),
    CONSTRAINT `fk_vm_cr_s` FOREIGN KEY (`SourceID`) REFERENCES `Source` (`ID`)
);

CREATE TRIGGER CategoryRuleTrigger BEFORE INSERT ON `CategoryRule`
    FOR EACH ROW
    SET new.ID = UUID();
