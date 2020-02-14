DROP TRIGGER OrganizationAuditCreateTrigger;
DROP TRIGGER OrganizationAuditUpdateTrigger;
DROP TRIGGER OrganizationAuditDeleteTrigger;

ALTER TABLE OrganizationAudit DROP COLUMN PortDupl;

CREATE TRIGGER OrganizationAuditCreateTrigger BEFORE INSERT ON `Organization`
    FOR EACH ROW
    INSERT INTO `OrganizationAudit` select new.ID, new.ParentOrgID, new.Code, new.Description, new.Payload, new.EncryptionKey, new.TimeZoneOffset, new.Created, new.Updated, new.CreatedBy, new.UpdatedBy, new.Active, 'INSERT', NEW.Created;

CREATE TRIGGER OrganizationAuditUpdateTrigger AFTER UPDATE ON `Organization`
    FOR EACH ROW
    INSERT INTO `OrganizationAudit` SELECT new.ID, new.ParentOrgID, new.Code, new.Description, new.Payload, new.EncryptionKey, new.TimeZoneOffset, new.Created, new.Updated, new.CreatedBy, new.UpdatedBy, new.Active, 'UPDATE', NOW();

CREATE TRIGGER OrganizationAuditDeleteTrigger BEFORE DELETE ON `Organization`
    FOR EACH ROW
    INSERT INTO `OrganizationAudit` SELECT old.ID, old.ParentOrgID, old.Code, old.Description, old.Payload, old.EncryptionKey, old.TimeZoneOffset, old.Created, old.Updated, old.CreatedBy, old.UpdatedBy, old.Active, 'DELETE', NOW();

