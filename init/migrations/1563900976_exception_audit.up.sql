CREATE TABLE `IgnoreAudit` LIKE `Ignore`;
ALTER TABLE `IgnoreAudit` ADD COLUMN EventType VARCHAR(20) NOT NULL;
ALTER TABLE `IgnoreAudit` ADD COLUMN EventDate DATETIME NOT NULL;
ALTER TABLE `IgnoreAudit` DROP PRIMARY KEY;

CREATE TRIGGER IgnoreAuditPreventDeleteTrigger BEFORE DELETE ON `IgnoreAudit`
    FOR EACH ROW
    SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'Cannot delete from an auditing table';

CREATE TRIGGER IgnoreAuditCreateTrigger BEFORE INSERT ON `Ignore`
    FOR EACH ROW
    INSERT INTO `IgnoreAudit` select new.ID, new.SourceID, new.OrganizationID, new.TypeId, new.VulnerabilityId, new.DeviceId, new.DueDate, new.Approval, new.Active, new.Port, new.DBCreatedDate, new.DBUpdatedDate, 'CREATE', NOW();

CREATE TRIGGER IgnoreAuditUpdateTrigger AFTER UPDATE ON `Ignore`
    FOR EACH ROW
    INSERT INTO `IgnoreAudit` SELECT new.ID, new.SourceID, new.OrganizationID, new.TypeId, new.VulnerabilityId, new.DeviceId, new.DueDate, new.Approval, new.Active, new.Port, new.DBCreatedDate, new.DBUpdatedDate, 'UPDATE', NOW();

CREATE TRIGGER IgnoreAuditDeleteTrigger BEFORE DELETE ON `Ignore`
    FOR EACH ROW
    INSERT INTO `IgnoreAudit` SELECT old.ID, old.SourceID, old.OrganizationID, old.TypeId, old.VulnerabilityId, old.DeviceId, old.DueDate, old.Approval, old.Active, old.Port, old.DBCreatedDate, old.DBUpdatedDate, 'DELETE', NOW();