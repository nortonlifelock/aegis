ALTER TABLE Ticket ADD COLUMN ExceptionDate DATETIME NULL AFTER UpdatedDate;
ALTER TABLE TicketAudit ADD COLUMN ExceptionDate DATETIME NULL AFTER UpdatedDate;

-- old vals I forgot to add
ALTER TABLE TicketAudit ADD COLUMN AssignmentGroup VARCHAR(200) NULL AFTER OrganizationID;
ALTER TABLE TicketAudit ADD COLUMN Assignee VARCHAR(100) NULL AFTER AssignmentGroup;
ALTER TABLE TicketAudit ADD COLUMN Created DATETIME NOT NULL DEFAULT NOW() AFTER Assignee;

DROP TRIGGER IF EXISTS TicketAuditCreateTrigger;
DROP TRIGGER IF EXISTS TicketAuditUpdateTrigger;
DROP TRIGGER IF EXISTS TicketAuditDeleteTrigger;

CREATE TRIGGER TicketAuditCreateTrigger BEFORE INSERT ON `Ticket`
    FOR EACH ROW
    INSERT INTO `TicketAudit` select new.Title, new.Status, new.DetectionID, new.OrganizationID, new.AssignmentGroup, new.Assignee, new.Created, new.ExceptionDate, new.UpdatedDate, new.ResolutionDate, new.DueDate, 'CREATE', NOW();

CREATE TRIGGER TicketAuditUpdateTrigger AFTER UPDATE ON `Ticket`
    FOR EACH ROW
    INSERT INTO `TicketAudit` select new.Title, new.Status, new.DetectionID, new.OrganizationID, new.AssignmentGroup, new.Assignee, new.Created, new.ExceptionDate, new.UpdatedDate, new.ResolutionDate, new.DueDate, 'UPDATE', NOW();

CREATE TRIGGER TicketAuditDeleteTrigger BEFORE DELETE ON `Ticket`
    FOR EACH ROW
    INSERT INTO `TicketAudit` SELECT old.Title, old.Status, old.DetectionID, old.OrganizationID, old.AssignmentGroup, old.Assignee, old.Created, old.ExceptionDate, old.UpdatedDate, old.ResolutionDate, old.DueDate, 'DELETE', NOW();