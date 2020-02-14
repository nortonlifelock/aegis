ALTER TABLE AssetGroup DROP FOREIGN KEY `fk_vm_ag_sc`;
ALTER TABLE AssetGroup DROP FOREIGN KEY `fk_vm_agcs_s`;
ALTER TABLE AssetGroup DROP FOREIGN KEY `fk_vm_agss_s`;
ALTER TABLE AssetGroup DROP FOREIGN KEY `fk_vm_assetg_o`;

ALTER TABLE AssetGroup DROP PRIMARY KEY;
ALTER TABLE AssetGroup ADD PRIMARY KEY (GroupID, OrganizationID, ScannerSourceConfigID);

ALTER TABLE AssetGroup ADD CONSTRAINT `fk_vm_ag_sc` FOREIGN KEY (`ScannerSourceConfigID`) REFERENCES `SourceConfig` (`ID`);
ALTER TABLE AssetGroup ADD CONSTRAINT `fk_vm_agcs_s` FOREIGN KEY (`CloudSourceID`) REFERENCES `Source` (`ID`);
ALTER TABLE AssetGroup ADD CONSTRAINT `fk_vm_agss_s` FOREIGN KEY (`ScannerSourceID`) REFERENCES `Source` (`ID`);
ALTER TABLE AssetGroup ADD CONSTRAINT `fk_vm_assetg_o` FOREIGN KEY (`OrganizationID`) REFERENCES `Organization` (`ID`);