ALTER TABLE `Ignore` MODIFY COLUMN DeviceID INT NULL;
ALTER TABLE `IgnoreAudit` MODIFY COLUMN DeviceID INT NULL;
ALTER TABLE `Ignore` ADD COLUMN OSRegex VARCHAR(100) NULL AFTER DeviceID;
ALTER TABLE `IgnoreAudit` ADD COLUMN OSRegex VARCHAR(100) NULL AFTER DeviceID;

DROP TRIGGER IgnoreAuditCreateTrigger;
DROP TRIGGER IgnoreAuditUpdateTrigger;
DROP TRIGGER IgnoreAuditDeleteTrigger;

CREATE TRIGGER IgnoreAuditCreateTrigger BEFORE INSERT ON `Ignore`
    FOR EACH ROW
    INSERT INTO `IgnoreAudit` select new.ID, new.SourceID, new.OrganizationID, new.TypeId, new.VulnerabilityId, new.DeviceId, new.OSRegex, new.DueDate, new.Approval, new.Active, new.Port, new.DBCreatedDate, new.DBUpdatedDate, 'CREATE', NOW();

CREATE TRIGGER IgnoreAuditUpdateTrigger AFTER UPDATE ON `Ignore`
    FOR EACH ROW
    INSERT INTO `IgnoreAudit` SELECT new.ID, new.SourceID, new.OrganizationID, new.TypeId, new.VulnerabilityId, new.DeviceId, new.OSRegex,  new.DueDate, new.Approval, new.Active, new.Port, new.DBCreatedDate, new.DBUpdatedDate, 'UPDATE', NOW();

CREATE TRIGGER IgnoreAuditDeleteTrigger BEFORE DELETE ON `Ignore`
    FOR EACH ROW
    INSERT INTO `IgnoreAudit` SELECT old.ID, old.SourceID, old.OrganizationID, old.TypeId, old.VulnerabilityId, old.DeviceId, old.OSRegex,  old.DueDate, old.Approval, old.Active, old.Port, old.DBCreatedDate, old.DBUpdatedDate, 'DELETE', NOW();