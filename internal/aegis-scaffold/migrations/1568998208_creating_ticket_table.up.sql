CREATE TABLE Ticket(
    Title          VARCHAR(36)       NOT NULL,
    Status         VARCHAR(100)      NOT NULL,
    DetectionID    VARCHAR(36)       NOT NULL,
    OrganizationID VARCHAR(100)      NOT NULL,
    UpdatedDate    DATETIME          NULL,
    ResolutionDate DATETIME          NULL,
    DueDate        DATETIME          NOT NULL,

    CONSTRAINT `fk_vm_t_d` FOREIGN KEY (`DetectionID`) REFERENCES `Detection` (`ID`),
    CONSTRAINT `fk_vm_t_o` FOREIGN KEY (`OrganizationID`) REFERENCES `Organization` (`ID`)
);

CREATE TABLE `TicketAudit` LIKE `Ticket`;
ALTER TABLE `TicketAudit` ADD COLUMN EventType VARCHAR(20) NOT NULL;
ALTER TABLE `TicketAudit` ADD COLUMN EventDate DATETIME NOT NULL;

CREATE TRIGGER TicketAuditPreventDeleteTrigger BEFORE DELETE ON `TicketAudit`
    FOR EACH ROW
    SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'Cannot delete from an auditing table';

CREATE TRIGGER TicketAuditCreateTrigger BEFORE INSERT ON `Ticket`
    FOR EACH ROW
    INSERT INTO `TicketAudit` select new.Title, new.Status, new.DetectionID, new.OrganizationID, new.UpdatedDate, new.ResolutionDate, new.DueDate, 'CREATE', NOW();

CREATE TRIGGER TicketAuditUpdateTrigger AFTER UPDATE ON `Ticket`
    FOR EACH ROW
    INSERT INTO `TicketAudit` SELECT new.Title, new.Status, new.DetectionID, new.OrganizationID, new.UpdatedDate, new.ResolutionDate, new.DueDate, 'UPDATE', NOW();

CREATE TRIGGER TicketAuditDeleteTrigger BEFORE DELETE ON `Ticket`
    FOR EACH ROW
    INSERT INTO `TicketAudit` SELECT old.Title, old.Status, old.DetectionID, old.OrganizationID, old.UpdatedDate, old.ResolutionDate, old.DueDate, 'DELETE', NOW();

CREATE UNIQUE INDEX ticket_title_org_index ON Ticket(Title, OrganizationID);
CREATE INDEX ticket_detection_id_index ON Ticket(DetectionID);

INSERT INTO Job (Struct, Priority, CreatedBy) VALUE ('TicketSyncJob', 5, 'ryan_everhart@symantec.com');