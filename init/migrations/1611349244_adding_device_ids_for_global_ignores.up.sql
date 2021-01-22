ALTER TABLE `Ignore` ADD COLUMN DeviceIDRegex VARCHAR(100) NULL AFTER HostnameRegex;
ALTER TABLE `IgnoreAudit` ADD COLUMN HostnameRegex VARCHAR(100) NULL AFTER OSRegex;
ALTER TABLE `IgnoreAudit` ADD COLUMN DeviceIDRegex VARCHAR(100) NULL AFTER HostnameRegex;

DROP TRIGGER IF EXISTS IgnoreAuditCreateTrigger;
DROP TRIGGER IF EXISTS IgnoreAuditUpdateTrigger;
DROP TRIGGER IF EXISTS IgnoreAuditDeleteTrigger;

CREATE TRIGGER IgnoreAuditCreateTrigger BEFORE INSERT ON `Ignore`
    FOR EACH ROW
    INSERT INTO `IgnoreAudit` select new.ID, new.SourceID, new.OrganizationID, new.TypeId, new.VulnerabilityId, new.DeviceId, new.OSRegex, new.HostnameRegex, new.DeviceIDRegex, new.DueDate, new.Approval, new.Active, new.Port, new.DBCreatedDate, new.DBUpdatedDate, 'CREATE', NOW();

CREATE TRIGGER IgnoreAuditUpdateTrigger AFTER UPDATE ON `Ignore`
    FOR EACH ROW
    INSERT INTO `IgnoreAudit` SELECT new.ID, new.SourceID, new.OrganizationID, new.TypeId, new.VulnerabilityId, new.DeviceId, new.OSRegex, new.HostnameRegex, new.DeviceIDRegex, new.DueDate, new.Approval, new.Active, new.Port, new.DBCreatedDate, new.DBUpdatedDate, 'UPDATE', NOW();

CREATE TRIGGER IgnoreAuditDeleteTrigger BEFORE DELETE ON `Ignore`
    FOR EACH ROW
    INSERT INTO `IgnoreAudit` SELECT old.ID, old.SourceID, old.OrganizationID, old.TypeId, old.VulnerabilityId, old.DeviceId, old.OSRegex, old.HostnameRegex, old.DeviceIDRegex, old.DueDate, old.Approval, old.Active, old.Port, old.DBCreatedDate, old.DBUpdatedDate, 'DELETE', NOW();