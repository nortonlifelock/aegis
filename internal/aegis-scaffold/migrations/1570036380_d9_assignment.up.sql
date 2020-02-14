CREATE TABLE Dome9Assignment(
    OrganizationID    VARCHAR(36)  NOT NULL,
    CloudAccountID    VARCHAR(100) NOT NULL,
    BundleID          VARCHAR(100) NULL,
    RuleRegex         VARCHAR(200) NULL,
    RuleHash          VARCHAR(100) NULL,
    AssignmentGroup   VARCHAR(100) NOT NULL,

    CONSTRAINT `fk_vm_dn_o` FOREIGN KEY (`OrganizationID`) REFERENCES `Organization` (`ID`)
);