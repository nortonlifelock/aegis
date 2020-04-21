ALTER TABLE `AssetGroup` ADD COLUMN RescanQueueSkip BOOL NOT NULL DEFAULT b'0';

CREATE TABLE `AssetGroupAudit` LIKE `AssetGroup`;
ALTER TABLE `AssetGroupAudit` ADD COLUMN EventType VARCHAR(20) NOT NULL;
ALTER TABLE `AssetGroupAudit` ADD COLUMN EventDate DATETIME NOT NULL;
ALTER TABLE `AssetGroupAudit` DROP PRIMARY KEY;

CREATE TRIGGER AssetGroupAuditPreventDeleteTrigger BEFORE DELETE ON `AssetGroupAudit`
    FOR EACH ROW
    SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'Cannot delete from an auditing table';

CREATE TRIGGER AssetGroupAuditCreateTrigger BEFORE INSERT ON `AssetGroup`
    FOR EACH ROW
    INSERT INTO `AssetGroupAudit` select new.GroupID, new.OrganizationID, new.ScannerSourceConfigID, new.ScannerSourceID, new.CloudSourceID, new.LastTicketing, new.RescanQueueSkip, 'INSERT', NOW();

CREATE TRIGGER AssetGroupAuditUpdateTrigger AFTER UPDATE ON `AssetGroup`
    FOR EACH ROW
    INSERT INTO `AssetGroupAudit` select new.GroupID, new.OrganizationID, new.ScannerSourceConfigID, new.ScannerSourceID, new.CloudSourceID, new.LastTicketing, new.RescanQueueSkip, 'UPDATE', NOW();

CREATE TRIGGER AssetGroupAuditDeleteTrigger BEFORE DELETE ON `AssetGroup`
    FOR EACH ROW
    INSERT INTO `AssetGroupAudit` select old.GroupID, old.OrganizationID, old.ScannerSourceConfigID, old.ScannerSourceID, old.CloudSourceID, old.LastTicketing, old.RescanQueueSkip, 'DELETE', NOW();