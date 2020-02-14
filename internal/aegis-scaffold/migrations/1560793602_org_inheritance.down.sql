ALTER TABLE `Organization` DROP FOREIGN KEY `fk_vm_o_o`;
ALTER TABLE `Organization` DROP COLUMN `ParentOrgID`;

ALTER TABLE `Permissions` ADD COLUMN `CanUpdateJob` BIT NOT NULL DEFAULT b'0';
ALTER TABLE `Permissions` ADD COLUMN `CanDeleteJob` BIT NOT NULL DEFAULT b'0';
ALTER TABLE `Permissions` ADD COLUMN `CanCreateJob` BIT NOT NULL DEFAULT b'0';
ALTER TABLE `Permissions` ADD COLUMN `CanUpdateConfig` BIT NOT NULL DEFAULT b'0';
ALTER TABLE `Permissions` ADD COLUMN `CanDeleteConfig` BIT NOT NULL DEFAULT b'0';
ALTER TABLE `Permissions` ADD COLUMN `CanCreateConfig` BIT NOT NULL DEFAULT b'0';
ALTER TABLE `Permissions` ADD COLUMN `CanUpdateSource` BIT NOT NULL DEFAULT b'0';
ALTER TABLE `Permissions` ADD COLUMN `CanDeleteSource` BIT NOT NULL DEFAULT b'0';
ALTER TABLE `Permissions` ADD COLUMN `CanCreateSource` BIT NOT NULL DEFAULT b'0';
ALTER TABLE `Permissions` ADD COLUMN `CanUpdateOrg` BIT NOT NULL DEFAULT b'0';
ALTER TABLE `Permissions` ADD COLUMN `CanDeleteOrg` BIT NOT NULL DEFAULT b'0';
ALTER TABLE `Permissions` ADD COLUMN `CanCreateOrg` BIT NOT NULL DEFAULT b'0';
ALTER TABLE `Permissions` ADD COLUMN `CanReadJobHistories` BIT NOT NULL DEFAULT b'0';
ALTER TABLE `Permissions` ADD COLUMN `CanRegisterUser` BIT NOT NULL DEFAULT b'0';
ALTER TABLE `Permissions` ADD COLUMN `CanUpdateUser` BIT NOT NULL DEFAULT b'0';
ALTER TABLE `Permissions` ADD COLUMN `CanDeleteUser` BIT NOT NULL DEFAULT b'0';
ALTER TABLE `Permissions` ADD COLUMN `CanReadUser` BIT NOT NULL DEFAULT b'0';
ALTER TABLE `Permissions` ADD COLUMN `CanBulkUpdate` BIT NOT NULL DEFAULT b'0';
ALTER TABLE `Permissions` ADD COLUMN `CanManageTags` BIT NOT NULL DEFAULT b'0';

ALTER TABLE `Permissions` DROP COLUMN Admin;
ALTER TABLE `Permissions` DROP COLUMN Manager;
ALTER TABLE `Permissions` DROP COLUMN Reader;
ALTER TABLE `Permissions` DROP COLUMN Reporter;

ALTER TABLE `Organization` DROP COLUMN EncryptionKey;