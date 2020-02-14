ALTER TABLE `Application` MODIFY ID VARCHAR(36); -- no fks
CREATE TRIGGER ApplicationTrigger BEFORE INSERT ON `Application`
    FOR EACH ROW
        SET new.ID = UUID();

ALTER TABLE `Category` DROP FOREIGN KEY `fk_vm_c_pc`;
ALTER TABLE `Category` MODIFY ID VARCHAR(36);
CREATE TRIGGER CategoryTrigger BEFORE INSERT ON `Category`
    FOR EACH ROW
        SET new.ID = UUID();
ALTER TABLE `Category` MODIFY ParentCategoryID VARCHAR(36);
ALTER TABLE `Category` ADD CONSTRAINT `fk_vm_c_pc` FOREIGN KEY (`ParentCategoryID`) REFERENCES `Category` (`ID`);


ALTER TABLE `DetectionPorts` DROP FOREIGN KEY `fk_vm_dp_d`;
ALTER TABLE `DetectionMetadata` DROP FOREIGN KEY `fk_vm_dm_d`;
ALTER TABLE `DetectionPorts` MODIFY DetectionID VARCHAR(36);
ALTER TABLE `DetectionMetadata` MODIFY DetectionID VARCHAR(36);
ALTER TABLE `Detection` MODIFY ID VARCHAR(36);
CREATE TRIGGER DetectionTrigger BEFORE INSERT ON `Detection`
    FOR EACH ROW
        SET new.ID = UUID();
ALTER TABLE `DetectionPorts` ADD CONSTRAINT `fk_vm_dp_d` FOREIGN KEY (`DetectionID`) REFERENCES `Detection` (`ID`);
ALTER TABLE `DetectionMetadata` ADD CONSTRAINT `fk_vm_dm_d` FOREIGN KEY (`DetectionID`) REFERENCES `Detection` (`ID`);

ALTER TABLE `DeviceInfo` DROP FOREIGN KEY `fk_vm_di_d`;
ALTER TABLE `DeviceMetadata` DROP FOREIGN KEY `fk_vm_dmd_d`;
ALTER TABLE `Tag` DROP FOREIGN KEY `fk_vm_t_dev`;
ALTER TABLE `Detection` DROP FOREIGN KEY `fk_vm_det_d`;
ALTER TABLE `DeviceInfo` MODIFY DeviceID VARCHAR(36);
ALTER TABLE `DeviceMetadata` MODIFY DeviceID VARCHAR(36);
ALTER TABLE `Tag` MODIFY DeviceID VARCHAR(36);
ALTER TABLE `Detection` MODIFY DeviceID VARCHAR(36);
ALTER TABLE `Device` MODIFY ID VARCHAR(36);
CREATE TRIGGER DeviceTrigger BEFORE INSERT ON `Device`
    FOR EACH ROW
        SET new.ID = UUID();
ALTER TABLE `DeviceInfo` ADD CONSTRAINT `fk_vm_di_d` FOREIGN KEY (`DeviceID`) REFERENCES `Device` (`ID`);
ALTER TABLE `DeviceMetadata` ADD CONSTRAINT `fk_vm_dmd_d` FOREIGN KEY (`DeviceID`) REFERENCES `Device` (`ID`);
ALTER TABLE `Tag` ADD CONSTRAINT `fk_vm_t_dev` FOREIGN KEY (`DeviceID`) REFERENCES `Device` (`ID`);
ALTER TABLE `Detection` ADD CONSTRAINT `fk_vm_det_d` FOREIGN KEY (`DeviceID`) REFERENCES `Device` (`ID`);

ALTER TABLE `DeviceGroupSource` DROP FOREIGN KEY `fk_vm_dgs_dg`;
ALTER TABLE `DeviceGroupSource` MODIFY DeviceGroupID VARCHAR(36);
ALTER TABLE `DeviceGroup` MODIFY ID VARCHAR(36);
CREATE TRIGGER DeviceGroupTrigger BEFORE INSERT ON `DeviceGroup`
    FOR EACH ROW
        SET new.ID = UUID();
ALTER TABLE `DeviceGroupSource` ADD CONSTRAINT `fk_vm_dgs_dg` FOREIGN KEY (`DeviceGroupID`) REFERENCES `DeviceGroup` (`ID`);


ALTER TABLE `Device` DROP FOREIGN KEY `fk_vm_d_di`;
ALTER TABLE `Device` MODIFY ImageID VARCHAR(36);
ALTER TABLE `DeviceImage` MODIFY ID VARCHAR(36);
CREATE TRIGGER DeviceImageTrigger BEFORE INSERT ON `DeviceImage`
    FOR EACH ROW
        SET new.ID = UUID();
ALTER TABLE `Device` ADD CONSTRAINT `fk_vm_d_di` FOREIGN KEY (`ImageID`) REFERENCES `DeviceImage` (`ID`);

ALTER TABLE `DeviceInfo` MODIFY ID VARCHAR(36); -- no fks
CREATE TRIGGER DeviceInfoTrigger BEFORE INSERT ON `DeviceInfo`
    FOR EACH ROW
        SET new.ID = UUID();


ALTER TABLE `Detection` DROP FOREIGN KEY `fk_vm_det_i`;
ALTER TABLE `Detection` MODIFY IgnoreID VARCHAR(36);
ALTER TABLE `Ignore` MODIFY ID VARCHAR(36);
CREATE TRIGGER IgnoreTrigger BEFORE INSERT ON `Ignore`
    FOR EACH ROW
        SET new.ID = UUID();
ALTER TABLE `Detection` ADD CONSTRAINT `fk_vm_det_i` FOREIGN KEY (`IgnoreID`) REFERENCES `Ignore` (`ID`);


ALTER TABLE `JobHistory` DROP FOREIGN KEY `fk_vm_jh_jc`;
ALTER TABLE `JobSchedule` DROP FOREIGN KEY `fk_vm_js_jc`;
ALTER TABLE `JobHistory` MODIFY ConfigID VARCHAR(36);
ALTER TABLE `JobSchedule` MODIFY ConfigID VARCHAR(36);
ALTER TABLE `JobConfig` MODIFY ID VARCHAR(36);
CREATE TRIGGER JobConfigTrigger BEFORE INSERT ON `JobConfig`
    FOR EACH ROW
        SET new.ID = UUID();
ALTER TABLE `JobHistory` ADD CONSTRAINT `fk_vm_jh_jc` FOREIGN KEY (`ConfigID`) REFERENCES `JobConfig` (`ID`);
ALTER TABLE `JobSchedule` ADD CONSTRAINT `fk_vm_js_jc` FOREIGN KEY (`ConfigID`) REFERENCES `JobConfig` (`ID`);


ALTER TABLE `JobHistory` DROP FOREIGN KEY `fk_vm_jh_pjh`;
ALTER TABLE `JobHistory` MODIFY ParentJobID VARCHAR(36);
ALTER TABLE `JobHistory` MODIFY ID VARCHAR(36);
CREATE TRIGGER JobHistoryTrigger BEFORE INSERT ON `JobHistory`
    FOR EACH ROW
        SET new.ID = UUID();
ALTER TABLE `JobHistory` ADD CONSTRAINT `fk_vm_jh_pjh` FOREIGN KEY (`ParentJobID`) REFERENCES `JobHistory` (`ID`);


ALTER TABLE `JobSchedule` MODIFY ID VARCHAR(36); -- no fks
CREATE TRIGGER JobScheduleTrigger BEFORE INSERT ON `JobSchedule`
    FOR EACH ROW
        SET new.ID = UUID();


ALTER TABLE `AssignmentGroup` DROP FOREIGN KEY `fk_vm_ag_o`;
ALTER TABLE `Device` DROP FOREIGN KEY `fk_vm_d_o`;
ALTER TABLE `DeviceGroup` DROP FOREIGN KEY `fk_vm_dg_o`;
ALTER TABLE `DeviceInfo` DROP FOREIGN KEY `fk_vm_di_o`;
ALTER TABLE `DeviceMetadata` DROP FOREIGN KEY `fk_vm_dmd_o`;
ALTER TABLE `JobConfig` DROP FOREIGN KEY `fk_vm_jc_o`;
ALTER TABLE `OwnerGroup` DROP FOREIGN KEY `fk_vm_og_o`;
ALTER TABLE `SourceConfig` DROP FOREIGN KEY `fk_vm_sc_o`;
ALTER TABLE `TagMap` DROP FOREIGN KEY `fk_vm_ts_o`;
ALTER TABLE `UserSession` DROP FOREIGN KEY `fk_vm_us_o`;
ALTER TABLE `Ignore` DROP FOREIGN KEY `fk_vm_ig_o`;
ALTER TABLE `Application` DROP FOREIGN KEY `fk_vm_a_o`;
ALTER TABLE `Detection` DROP FOREIGN KEY `fk_vm_det_o`;
ALTER TABLE `DetectionMetadata` DROP FOREIGN KEY `fk_vm_dm_o`;
ALTER TABLE `AssignmentGroup` MODIFY OrganizationID VARCHAR(36);
ALTER TABLE `Device` MODIFY OrganizationID VARCHAR(36);
ALTER TABLE `DeviceGroup` MODIFY OrganizationID VARCHAR(36);
ALTER TABLE `DeviceInfo` MODIFY OrganizationID VARCHAR(36);
ALTER TABLE `DeviceMetadata` MODIFY OrganizationID VARCHAR(36);
ALTER TABLE `JobConfig` MODIFY OrganizationID VARCHAR(36);
ALTER TABLE `OwnerGroup` MODIFY OrganizationID VARCHAR(36);
ALTER TABLE `SourceConfig` MODIFY OrganizationID VARCHAR(36);
ALTER TABLE `TagMap` MODIFY OrganizationID VARCHAR(36);
ALTER TABLE `UserSession` MODIFY OrgID VARCHAR(36);
ALTER TABLE `Ignore` MODIFY OrganizationID VARCHAR(36);
ALTER TABLE `Application` MODIFY OrganizationID VARCHAR(36);
ALTER TABLE `Detection` MODIFY OrganizationID VARCHAR(36);
ALTER TABLE `DetectionMetadata` MODIFY OrganizationID VARCHAR(36);
ALTER TABLE `Organization` MODIFY ID VARCHAR(36);
CREATE TRIGGER OrganizationTrigger BEFORE INSERT ON `Organization`
    FOR EACH ROW
        SET new.ID = UUID();
ALTER TABLE `AssignmentGroup` ADD CONSTRAINT `fk_vm_ag_o` FOREIGN KEY (`OrganizationID`) REFERENCES `Organization` (`ID`);
ALTER TABLE `Device` ADD CONSTRAINT `fk_vm_d_o` FOREIGN KEY (`OrganizationID`) REFERENCES `Organization` (`ID`);
ALTER TABLE `DeviceGroup` ADD CONSTRAINT `fk_vm_dg_o` FOREIGN KEY (`OrganizationID`) REFERENCES `Organization` (`ID`);
ALTER TABLE `DeviceInfo` ADD CONSTRAINT `fk_vm_di_o` FOREIGN KEY (`OrganizationID`) REFERENCES `Organization` (`ID`);
ALTER TABLE `DeviceMetadata` ADD CONSTRAINT `fk_vm_dmd_o` FOREIGN KEY (`OrganizationID`) REFERENCES `Organization` (`ID`);
ALTER TABLE `JobConfig` ADD CONSTRAINT `fk_vm_jc_o` FOREIGN KEY (`OrganizationID`) REFERENCES `Organization` (`ID`);
ALTER TABLE `OwnerGroup` ADD CONSTRAINT `fk_vm_og_o` FOREIGN KEY (`OrganizationID`) REFERENCES `Organization` (`ID`);
ALTER TABLE `SourceConfig` ADD CONSTRAINT `fk_vm_sc_o` FOREIGN KEY (`OrganizationID`) REFERENCES `Organization` (`ID`);
ALTER TABLE `TagMap` ADD CONSTRAINT `fk_vm_ts_o` FOREIGN KEY (`OrganizationID`) REFERENCES `Organization` (`ID`);
ALTER TABLE `UserSession` ADD CONSTRAINT `fk_vm_us_o` FOREIGN KEY (`OrgID`) REFERENCES `Organization` (`ID`);
ALTER TABLE `Ignore` ADD CONSTRAINT `fk_vm_ig_o` FOREIGN KEY (`OrganizationID`) REFERENCES `Organization` (`ID`);
ALTER TABLE `Application` ADD CONSTRAINT `fk_vm_a_o` FOREIGN KEY (`OrganizationID`) REFERENCES `Organization` (`ID`);
ALTER TABLE `Detection` ADD CONSTRAINT `fk_vm_det_o` FOREIGN KEY (`OrganizationID`) REFERENCES `Organization` (`ID`);
ALTER TABLE `DetectionMetadata` ADD CONSTRAINT `fk_vm_dm_o` FOREIGN KEY (`OrganizationID`) REFERENCES `Organization` (`ID`);


ALTER TABLE `Subscription` DROP FOREIGN KEY `fk_vm_su_og`;
ALTER TABLE `Subscription` MODIFY OwnerGroupID VARCHAR(36);
ALTER TABLE `OwnerGroup` MODIFY ID VARCHAR(36);
CREATE TRIGGER OwnerGroupTrigger BEFORE INSERT ON `OwnerGroup`
    FOR EACH ROW
        SET new.ID = UUID();
ALTER TABLE `Subscription` ADD CONSTRAINT `fk_vm_su_og` FOREIGN KEY (`OwnerGroupID`) REFERENCES `OwnerGroup` (`ID`);


ALTER TABLE `AssignmentGroup` DROP FOREIGN KEY `fk_vm_ag_s`;
ALTER TABLE `DeviceGroupSource` DROP FOREIGN KEY `fk_vm_dgs_s`;
ALTER TABLE `DeviceInfo` DROP FOREIGN KEY `fk_vm_di_s`;
ALTER TABLE `DeviceMetadata` DROP FOREIGN KEY `fk_vm_dmd_s`;
ALTER TABLE `SourceConfig` DROP FOREIGN KEY `fk_vm_sc_sc`;
ALTER TABLE `TagMap` DROP FOREIGN KEY `pk_vm_cs_s`;
ALTER TABLE `TagMap` DROP FOREIGN KEY `pk_vm_ts_s`;
ALTER TABLE `Ignore` DROP FOREIGN KEY `fk_vm_ig_s`;
ALTER TABLE `Detection` DROP FOREIGN KEY `fk_vm_det_s`;
ALTER TABLE `DetectionMetadata` DROP FOREIGN KEY `fk_vm_dm_s`;
ALTER TABLE `Vulnerability_Reference` DROP FOREIGN KEY `fk_vm_vr_s`;
ALTER TABLE `VulnerabilityInfo` DROP FOREIGN KEY `fk_vm_vi_s`;

ALTER TABLE `VulnerabilityInfo` MODIFY SourceID VARCHAR(36);
ALTER TABLE `Vulnerability_Reference` MODIFY SourceID VARCHAR(36);
ALTER TABLE `AssignmentGroup` MODIFY SourceID VARCHAR(36);
ALTER TABLE `DeviceGroupSource` MODIFY SourceID VARCHAR(36);
ALTER TABLE `DeviceInfo` MODIFY SourceID VARCHAR(36);
ALTER TABLE `DeviceMetadata` MODIFY SourceID VARCHAR(36);
ALTER TABLE `SourceConfig` MODIFY SourceID VARCHAR(36);
ALTER TABLE `TagMap` MODIFY CloudSourceID VARCHAR(36);
ALTER TABLE `TagMap` MODIFY TicketingSourceID VARCHAR(36);
ALTER TABLE `Ignore` MODIFY SourceID VARCHAR(36);
ALTER TABLE `Detection` MODIFY SourceID VARCHAR(36);
ALTER TABLE `DetectionMetadata` MODIFY SourceID VARCHAR(36);
ALTER TABLE `Source` MODIFY ID VARCHAR(36);
CREATE TRIGGER SourceTrigger BEFORE INSERT ON `Source`
    FOR EACH ROW
        SET new.ID = UUID();

ALTER TABLE `VulnerabilityInfo` ADD CONSTRAINT `fk_vm_vi_s` FOREIGN KEY (`SourceID`) REFERENCES `Source` (`ID`);
ALTER TABLE `Vulnerability_Reference` ADD CONSTRAINT `fk_vm_vr_s` FOREIGN KEY (`SourceID`) REFERENCES `Source` (`ID`);
ALTER TABLE `AssignmentGroup` ADD CONSTRAINT `fk_vm_ag_s` FOREIGN KEY (`SourceID`) REFERENCES `Source` (`ID`);
ALTER TABLE `DeviceGroupSource` ADD CONSTRAINT `fk_vm_dgs_s` FOREIGN KEY (`SourceID`) REFERENCES `Source` (`ID`);
ALTER TABLE `DeviceInfo` ADD CONSTRAINT `fk_vm_di_s` FOREIGN KEY (`SourceID`) REFERENCES `Source` (`ID`);
ALTER TABLE `DeviceMetadata` ADD CONSTRAINT `fk_vm_dmd_s` FOREIGN KEY (`SourceID`) REFERENCES `Source` (`ID`);
ALTER TABLE `SourceConfig` ADD CONSTRAINT `fk_vm_sc_sc` FOREIGN KEY (`SourceID`) REFERENCES `Source` (`ID`);
ALTER TABLE `TagMap` ADD CONSTRAINT `pk_vm_cs_s` FOREIGN KEY (`CloudSourceID`) REFERENCES `Source` (`ID`);
ALTER TABLE `TagMap` ADD CONSTRAINT `pk_vm_ts_s` FOREIGN KEY (`TicketingSourceID`) REFERENCES `Source` (`ID`);
ALTER TABLE `Ignore` ADD CONSTRAINT `fk_vm_ig_s` FOREIGN KEY (`SourceID`) REFERENCES `Source` (`ID`);
ALTER TABLE `Detection` ADD CONSTRAINT `fk_vm_det_s` FOREIGN KEY (`SourceID`) REFERENCES `Source` (`ID`);
ALTER TABLE `DetectionMetadata` ADD CONSTRAINT `fk_vm_dm_s` FOREIGN KEY (`SourceID`) REFERENCES `Source` (`ID`);


ALTER TABLE `JobConfig` DROP FOREIGN KEY `fk_vm_jc_sci`;
ALTER TABLE `JobConfig` DROP FOREIGN KEY `fk_vm_jc_sco`;
ALTER TABLE `JobConfig` MODIFY DataInSourceConfigID VARCHAR(36);
ALTER TABLE `JobConfig` MODIFY DataOutSourceConfigID VARCHAR(36);
ALTER TABLE `SourceConfig` MODIFY ID VARCHAR(36);
CREATE TRIGGER SourceConfigTrigger BEFORE INSERT ON `SourceConfig`
    FOR EACH ROW
        SET new.ID = UUID();
ALTER TABLE `JobConfig` ADD CONSTRAINT `fk_vm_jc_sci` FOREIGN KEY (`DataInSourceConfigID`) REFERENCES `SourceConfig` (`ID`);
ALTER TABLE `JobConfig` ADD CONSTRAINT `fk_vm_jc_sco` FOREIGN KEY (`DataOutSourceConfigID`) REFERENCES `SourceConfig` (`ID`);


ALTER TABLE `Subscription` MODIFY ID VARCHAR(36); -- no fks
CREATE TRIGGER SubscriptionTrigger BEFORE INSERT ON `Subscription`
    FOR EACH ROW
        SET new.ID = UUID();


ALTER TABLE `Tag` MODIFY ID VARCHAR(36); -- no fks
CREATE TRIGGER TagTrigger BEFORE INSERT ON `Tag`
    FOR EACH ROW
        SET new.ID = UUID();

ALTER TABLE `TagMap` MODIFY ID VARCHAR(36); -- no fks
CREATE TRIGGER TagMapTrigger BEFORE INSERT ON `TagMap`
    FOR EACH ROW
        SET new.ID = UUID();

ALTER TABLE `VulnerabilityInfo` DROP FOREIGN KEY `fk_vm_vi_v`;
ALTER TABLE `VulnerabilityInfo` MODIFY VulnerabilityID VARCHAR(36);
ALTER TABLE `Vulnerability` MODIFY ID VARCHAR(36);
CREATE TRIGGER VulnerabilityTrigger BEFORE INSERT ON `Vulnerability`
    FOR EACH ROW
        SET new.ID = UUID();
ALTER TABLE `VulnerabilityInfo` ADD CONSTRAINT `fk_vm_vi_v` FOREIGN KEY (`VulnerabilityID`) REFERENCES `Vulnerability` (`ID`);


ALTER TABLE `Vulnerability_Reference` MODIFY ID VARCHAR(36); -- no fks
CREATE TRIGGER Vulnerability_ReferenceTrigger BEFORE INSERT ON `Vulnerability_Reference`
    FOR EACH ROW
        SET new.ID = UUID();


ALTER TABLE `Vulnerability_Reference` DROP FOREIGN KEY `fk_vm_vr_vi`;
ALTER TABLE `Detection` DROP FOREIGN KEY `fk_vm_det_v`;
ALTER TABLE `Vulnerability_Reference` MODIFY VulnInfoID VARCHAR(36);
ALTER TABLE `Detection` MODIFY VulnerabilityID VARCHAR(36);
ALTER TABLE `VulnerabilityInfo` MODIFY ID VARCHAR(36);
CREATE TRIGGER VulnerabilityInfoTrigger BEFORE INSERT ON `VulnerabilityInfo`
    FOR EACH ROW
        SET new.ID = UUID();
ALTER TABLE `Vulnerability_Reference` ADD CONSTRAINT `fk_vm_vr_vi` FOREIGN KEY (`VulnInfoID`) REFERENCES `VulnerabilityInfo` (`ID`);
ALTER TABLE `Detection` ADD CONSTRAINT `fk_vm_det_v` FOREIGN KEY (`VulnerabilityID`) REFERENCES `VulnerabilityInfo` (`ID`);

ALTER TABLE `Permissions` DROP FOREIGN KEY `fk_vm_p_u`;
ALTER TABLE `UserSession` DROP FOREIGN KEY `fk_vm_us_u`;
ALTER TABLE `Permissions` MODIFY UserID VARCHAR(36);
ALTER TABLE `Permissions` MODIFY OrgID VARCHAR(36);
ALTER TABLE `UserSession` MODIFY UserID VARCHAR(36);
ALTER TABLE `Users` DROP COLUMN UUID;
ALTER TABLE `Users` MODIFY ID VARCHAR(36);
CREATE TRIGGER UserTrigger BEFORE INSERT ON `Users`
    FOR EACH ROW
    SET new.ID = UUID();
ALTER TABLE `Permissions` ADD CONSTRAINT `fk_vm_p_u` FOREIGN KEY (`UserId`) REFERENCES  `Users` (`ID`);
ALTER TABLE `UserSession` ADD CONSTRAINT `fk_vm_us_u` FOREIGN KEY (`UserId`) REFERENCES  `Users` (`ID`);