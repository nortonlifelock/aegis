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

UPDATE Job SET SourceTypeIn = '2', SourceTypeOut = '1' WHERE ID = '1';
UPDATE Job SET SourceTypeIn = '1', SourceTypeOut = '2' WHERE ID = '2';
UPDATE Job SET SourceTypeIn = '2', SourceTypeOut = '1' WHERE ID = '3';
UPDATE Job SET SourceTypeIn = '2', SourceTypeOut = '1' WHERE ID = '4';
UPDATE Job SET SourceTypeIn = '1', SourceTypeOut = '1' WHERE ID = '5';
UPDATE Job SET SourceTypeIn = '1', SourceTypeOut = '2' WHERE ID = '6';
UPDATE Job SET SourceTypeIn = '2', SourceTypeOut = '2' WHERE ID = '7';
UPDATE Job SET SourceTypeIn = '1', SourceTypeOut = '1' WHERE ID = '8';
UPDATE Job SET SourceTypeIn = '3', SourceTypeOut = '3' WHERE ID = '9';
UPDATE Job SET SourceTypeIn = '1', SourceTypeOut = '1' WHERE ID = '10';
UPDATE Job SET SourceTypeIn = '4', SourceTypeOut = '2' WHERE ID = '11';
UPDATE Job SET SourceTypeIn = '2', SourceTypeOut = '2' WHERE ID = '12';
