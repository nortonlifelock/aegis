CREATE TABLE `SourceConfigAudit` LIKE `SourceConfig`;
ALTER TABLE `SourceConfigAudit` ADD COLUMN EventType VARCHAR(20) NOT NULL;
ALTER TABLE `SourceConfigAudit` ADD COLUMN EventDate DATETIME NOT NULL;
ALTER TABLE `SourceConfigAudit` MODIFY COLUMN `AuthInfo` TEXT NOT NULL;
ALTER TABLE `SourceConfigAudit` DROP PRIMARY KEY;

CREATE TRIGGER SourceConfigAuditPreventDeleteTrigger BEFORE DELETE ON `SourceConfigAudit`
    FOR EACH ROW
    SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'Cannot delete from an auditing table';

CREATE TRIGGER SourceConfigAuditCreateTrigger BEFORE INSERT ON `SourceConfig`
    FOR EACH ROW
    INSERT INTO `SourceConfigAudit` select new.ID, new.Source, new.SourceID, new.OrganizationID, new.Address, new.Port, SHA2(new.AuthInfo, 256), new.DBCreatedDate, new.DBUpdatedDate, new.Payload, new.Active, new.UpdatedBy, 'INSERT', NEW.DBCreatedDate;

CREATE TRIGGER SourceConfigAuditUpdateTrigger AFTER UPDATE ON `SourceConfig`
    FOR EACH ROW
    INSERT INTO `SourceConfigAudit` SELECT new.ID, new.Source, new.SourceID, new.OrganizationID, new.Address, new.Port, SHA2(new.AuthInfo, 256), new.DBCreatedDate, new.DBUpdatedDate, new.Payload, new.Active, new.UpdatedBy, 'UPDATE', NOW();

CREATE TRIGGER SourceConfigAuditDeleteTrigger BEFORE DELETE ON `SourceConfig`
    FOR EACH ROW
    INSERT INTO `SourceConfigAudit` SELECT old.ID, old.Source, old.SourceID, old.OrganizationID, old.Address, old.Port, SHA2(old.AuthInfo, 256), old.DBCreatedDate, old.DBUpdatedDate, old.Payload, old.Active, old.UpdatedBy, 'DELETE', NOW();


CREATE TABLE `JobConfigAudit` LIKE `JobConfig`;
ALTER TABLE `JobConfigAudit` ADD COLUMN EventType VARCHAR(20) NOT NULL;
ALTER TABLE `JobConfigAudit` ADD COLUMN EventDate DATETIME NOT NULL;
ALTER TABLE `JobConfigAudit` DROP PRIMARY KEY;

CREATE TRIGGER JobConfigAuditPreventDeleteTrigger BEFORE DELETE ON `JobConfigAudit`
    FOR EACH ROW
    SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'Cannot delete from an auditing table';

CREATE TRIGGER JobConfigAuditCreateTrigger BEFORE INSERT ON `JobConfig`
    FOR EACH ROW
    INSERT INTO `JobConfigAudit` select new.ID, new.JobId, new.OrganizationID, new.DataInSourceConfigID, new.DataOutSourceConfigID, new.Payload, new.PriorityOverride, new.Continuous, new.WaitInSeconds, new.MaxInstances, new.AutoStart, new.CreatedDate, new.CreatedBy, new.UpdatedDate, new.UpdatedBy, new.Active, new.LastJobStart, 'INSERT', NEW.CreatedDate;

CREATE TRIGGER JobConfigAuditUpdateTrigger AFTER UPDATE ON `JobConfig`
    FOR EACH ROW
    INSERT INTO `JobConfigAudit` SELECT new.ID, new.JobId, new.OrganizationID, new.DataInSourceConfigID, new.DataOutSourceConfigID, new.Payload, new.PriorityOverride, new.Continuous, new.WaitInSeconds, new.MaxInstances, new.AutoStart, new.CreatedDate, new.CreatedBy, new.UpdatedDate, new.UpdatedBy, new.Active, new.LastJobStart, 'UPDATE', NOW();

CREATE TRIGGER JobConfigAuditDeleteTrigger BEFORE DELETE ON `JobConfig`
    FOR EACH ROW
    INSERT INTO `JobConfigAudit` SELECT old.ID, old.JobId, old.OrganizationID, old.DataInSourceConfigID, old.DataOutSourceConfigID, old.Payload, old.PriorityOverride, old.Continuous, old.WaitInSeconds, old.MaxInstances, old.AutoStart, old.CreatedDate, old.CreatedBy, old.UpdatedDate, old.UpdatedBy, old.Active, old.LastJobStart, 'DELETE', NOW();

CREATE TABLE `UsersAudit` LIKE `Users`;
ALTER TABLE `UsersAudit` ADD COLUMN EventType VARCHAR(20) NOT NULL;
ALTER TABLE `UsersAudit` ADD COLUMN EventDate DATETIME NOT NULL;
ALTER TABLE `UsersAudit` DROP PRIMARY KEY;

CREATE TRIGGER UsersAuditPreventDeleteTrigger BEFORE DELETE ON `UsersAudit`
    FOR EACH ROW
    SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'Cannot delete from an auditing table';

CREATE TRIGGER UsersAuditCreateTrigger BEFORE INSERT ON `Users`
    FOR EACH ROW
    INSERT INTO `UsersAudit` select new.ID, new.Username, new.FirstName, new.LastName, new.Email, new.IsDisabled, 'INSERT', NOW();

CREATE TRIGGER UsersAuditUpdateTrigger AFTER UPDATE ON `Users`
    FOR EACH ROW
    INSERT INTO `UsersAudit` SELECT new.ID, new.Username, new.FirstName, new.LastName, new.Email, new.IsDisabled, 'UPDATE', NOW();

CREATE TRIGGER UsersAuditDeleteTrigger BEFORE DELETE ON `Users`
    FOR EACH ROW
    INSERT INTO `UsersAudit` SELECT old.ID, old.Username, old.FirstName, old.LastName, old.Email, old.IsDisabled, 'DELETE', NOW();


CREATE TABLE `PermissionsAudit` LIKE `Permissions`;
ALTER TABLE `PermissionsAudit` ADD COLUMN EventType VARCHAR(20) NOT NULL;
ALTER TABLE `PermissionsAudit` ADD COLUMN EventDate DATETIME NOT NULL;
ALTER TABLE `PermissionsAudit` DROP PRIMARY KEY;

CREATE TRIGGER PermissionsAuditPreventDeleteTrigger BEFORE DELETE ON `PermissionsAudit`
    FOR EACH ROW
    SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'Cannot delete from an auditing table';

CREATE TRIGGER PermissionsAuditCreateTrigger BEFORE INSERT ON `Permissions`
    FOR EACH ROW
    INSERT INTO `PermissionsAudit` SELECT new.UserID, new.OrgID, new.Admin, new.Manager, new.Reader, new.Reporter, 'INSERT', NOW();

CREATE TRIGGER PermissionsAuditUpdateTrigger AFTER UPDATE ON `Permissions`
    FOR EACH ROW
    INSERT INTO `PermissionsAudit` SELECT new.UserID, new.OrgID, new.Admin, new.Manager, new.Reader, new.Reporter, 'UPDATE', NOW();

CREATE TRIGGER PermissionsAuditDeleteTrigger BEFORE DELETE ON `Permissions`
    FOR EACH ROW
    INSERT INTO `PermissionsAudit` SELECT old.UserID, old.OrgID, old.Admin, old.Manager, old.Reader, old.Reporter, 'DELETE', NOW();


CREATE TABLE `VulnerabilityInfoAudit` LIKE `VulnerabilityInfo`;
ALTER TABLE `VulnerabilityInfoAudit` ADD COLUMN EventType VARCHAR(20) NOT NULL;
ALTER TABLE `VulnerabilityInfoAudit` ADD COLUMN EventDate DATETIME NOT NULL;
ALTER TABLE `VulnerabilityInfoAudit` DROP PRIMARY KEY;

CREATE TRIGGER VulnerabilityInfoAuditPreventDeleteTrigger BEFORE DELETE ON `VulnerabilityInfoAudit`
    FOR EACH ROW
    SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'Cannot delete from an auditing table';

CREATE TRIGGER VulnerabilityInfoAuditCreateTrigger BEFORE INSERT ON `VulnerabilityInfo`
    FOR EACH ROW
    INSERT INTO `VulnerabilityInfoAudit` select new.ID, new.SourceVulnId, new.Title, new.VulnerabilityID, new.SourceID, new.CVSS, new.CVSS3, new.Description, new.Solution, new.Severity, new.CVSSVector, new.CVSS3Vector, new.MatchConfidence, new.MatchReasons, new.Software, new.DetectionInformation, new.Updated, new.Created, 'INSERT', NEW.Created;

CREATE TRIGGER VulnerabilityInfoAuditUpdateTrigger AFTER UPDATE ON `VulnerabilityInfo`
    FOR EACH ROW
    INSERT INTO `VulnerabilityInfoAudit` SELECT new.ID, new.SourceVulnId, new.Title, new.VulnerabilityID, new.SourceID, new.CVSS, new.CVSS3, new.Description, new.Solution, new.Severity, new.CVSSVector, new.CVSS3Vector, new.MatchConfidence, new.MatchReasons, new.Software, new.DetectionInformation, new.Updated, new.Created, 'UPDATE', NOW();

CREATE TRIGGER VulnerabilityInfoDeleteTrigger BEFORE DELETE ON `VulnerabilityInfo`
    FOR EACH ROW
    SIGNAL SQLSTATE '45000' -- "unhandled user-defined exception"
        SET MESSAGE_TEXT = 'Cannot delete from this table';