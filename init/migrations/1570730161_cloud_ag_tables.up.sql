CREATE TABLE AssetGroup(
    GroupID         INT         NOT NULL,
    ScannerSourceID VARCHAR(36) NOT NULL,
    CloudSourceID   VARCHAR(36) NULL,

    PRIMARY KEY (GroupID, ScannerSourceID),
    CONSTRAINT `fk_vm_agss_s` FOREIGN KEY (`ScannerSourceID`) REFERENCES `Source` (`ID`),
    CONSTRAINT `fk_vm_agcs_s` FOREIGN KEY (`CloudSourceID`) REFERENCES `Source` (`ID`)
);

ALTER TABLE Device ADD CONSTRAINT `fk_vm_d_ag` FOREIGN KEY (`GroupID`) REFERENCES `AssetGroup` (`GroupID`);