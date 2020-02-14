ALTER TABLE `Organization` ADD COLUMN ParentOrgID VARCHAR(36) NULL AFTER ID;
ALTER TABLE `Organization` ADD CONSTRAINT `fk_vm_o_o` FOREIGN KEY (`ParentOrgID`) REFERENCES  `Organization` (`ID`);

ALTER TABLE `Permissions` DROP COLUMN `CanUpdateJob`;
ALTER TABLE `Permissions` DROP COLUMN `CanDeleteJob`;
ALTER TABLE `Permissions` DROP COLUMN `CanCreateJob`;
ALTER TABLE `Permissions` DROP COLUMN `CanUpdateConfig`;
ALTER TABLE `Permissions` DROP COLUMN `CanDeleteConfig`;
ALTER TABLE `Permissions` DROP COLUMN `CanCreateConfig`;
ALTER TABLE `Permissions` DROP COLUMN `CanUpdateSource`;
ALTER TABLE `Permissions` DROP COLUMN `CanDeleteSource`;
ALTER TABLE `Permissions` DROP COLUMN `CanCreateSource`;
ALTER TABLE `Permissions` DROP COLUMN `CanUpdateOrg`;
ALTER TABLE `Permissions` DROP COLUMN `CanDeleteOrg`;
ALTER TABLE `Permissions` DROP COLUMN `CanCreateOrg`;
ALTER TABLE `Permissions` DROP COLUMN `CanReadJobHistories`;
ALTER TABLE `Permissions` DROP COLUMN `CanRegisterUser`;
ALTER TABLE `Permissions` DROP COLUMN `CanUpdateUser`;
ALTER TABLE `Permissions` DROP COLUMN `CanDeleteUser`;
ALTER TABLE `Permissions` DROP COLUMN `CanReadUser`;
ALTER TABLE `Permissions` DROP COLUMN `CanBulkUpdate`;
ALTER TABLE `Permissions` DROP COLUMN `CanManageTags`;


ALTER TABLE `Permissions` ADD COLUMN Admin BIT NOT NULL DEFAULT b'0';
ALTER TABLE `Permissions` ADD COLUMN Manager BIT NOT NULL DEFAULT b'0';
ALTER TABLE `Permissions` ADD COLUMN Reader BIT NOT NULL DEFAULT b'0';
ALTER TABLE `Permissions` ADD COLUMN Reporter BIT NOT NULL DEFAULT b'0';

ALTER TABLE `Organization` ADD COLUMN EncryptionKey TEXT NULL AFTER Payload;