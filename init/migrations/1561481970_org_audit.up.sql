CREATE TABLE `OrganizationAudit` LIKE `Organization`;
ALTER TABLE `OrganizationAudit` ADD COLUMN EventType VARCHAR(20) NOT NULL;
ALTER TABLE `OrganizationAudit` ADD COLUMN EventDate DATETIME NOT NULL;
ALTER TABLE `OrganizationAudit` DROP PRIMARY KEY;

CREATE TRIGGER OrganizationAuditPreventDeleteTrigger BEFORE DELETE ON `OrganizationAudit`
    FOR EACH ROW
    SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'Cannot delete from an auditing table';

CREATE TRIGGER OrganizationAuditCreateTrigger BEFORE INSERT ON `Organization`
    FOR EACH ROW
        INSERT INTO `OrganizationAudit` select new.ID, new.ParentOrgID, new.Code, new.Description, new.Payload, new.EncryptionKey, new.TimeZoneOffset, new.Created, new.Updated, new.CreatedBy, new.UpdatedBy, new.PortDupl, new.Active, 'INSERT', NEW.Created;

CREATE TRIGGER OrganizationAuditUpdateTrigger AFTER UPDATE ON `Organization`
    FOR EACH ROW
        INSERT INTO `OrganizationAudit` SELECT new.ID, new.ParentOrgID, new.Code, new.Description, new.Payload, new.EncryptionKey, new.TimeZoneOffset, new.Created, new.Updated, new.CreatedBy, new.UpdatedBy, new.PortDupl, new.Active, 'UPDATE', NOW();

CREATE TRIGGER OrganizationAuditDeleteTrigger BEFORE DELETE ON `Organization`
    FOR EACH ROW
        INSERT INTO `OrganizationAudit` SELECT old.ID, old.ParentOrgID, old.Code, old.Description, old.Payload, old.EncryptionKey, old.TimeZoneOffset, old.Created, old.Updated, old.CreatedBy, old.UpdatedBy, old.PortDupl, old.Active, 'DELETE', NOW();