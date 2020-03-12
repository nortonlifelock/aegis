DROP TABLE IF EXISTS AssignmentRules;

CREATE TABLE AssignmentRules(
    AssignmentGroup    VARCHAR(100) NULL,
    Assignee           VARCHAR(100) NULL,
    OrganizationID     VARCHAR(36) NOT NULL,
    VulnTitleSubstring VARCHAR(200) NULL,
    VulnTitleRegex     VARCHAR(100) NULL,
    TagKeyID           INT NULL,
    TagKeyValue        VARCHAR(100) NULL,
    Priority           INT NOT NULL DEFAULT 0,
    KEY `fk_vm_ar_o` (`OrganizationID`),
    KEY `fk_vm_ar_tk` (`TagKeyID`),
    CONSTRAINT `fk_vm_ar_o` FOREIGN KEY (`OrganizationID`) REFERENCES `Organization` (`ID`),
    CONSTRAINT `fk_vm_ar_tk` FOREIGN KEY (`TagKeyID`) REFERENCES `TagKey` (`Id`)
);

ALTER TABLE AssignmentRules DROP COLUMN VulnTitleSubstring;
ALTER TABLE AssignmentRules CHANGE `TagKeyValue` `TagKeyRegex` VARCHAR(100) NULL;