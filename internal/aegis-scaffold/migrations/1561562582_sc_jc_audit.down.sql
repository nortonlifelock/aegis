DROP TABLE IF EXISTS `SourceConfigAudit`;
DROP TRIGGER IF EXISTS SourceConfigAuditCreateTrigger;
DROP TRIGGER IF EXISTS SourceConfigAuditUpdateTrigger;
DROP TRIGGER IF EXISTS SourceConfigAuditDeleteTrigger;

DROP TABLE IF EXISTS `JobConfigAudit`;
DROP TRIGGER IF EXISTS JobConfigAuditCreateTrigger;
DROP TRIGGER IF EXISTS JobConfigAuditUpdateTrigger;
DROP TRIGGER IF EXISTS JobConfigAuditDeleteTrigger;

DROP TABLE IF EXISTS `UsersAudit`;
DROP TRIGGER IF EXISTS UsersAuditCreateTrigger;
DROP TRIGGER IF EXISTS UsersAuditUpdateTrigger;
DROP TRIGGER IF EXISTS UsersAuditDeleteTrigger;

DROP TABLE IF EXISTS `PermissionsAudit`;
DROP TRIGGER IF EXISTS PermissionsAuditCreateTrigger;
DROP TRIGGER IF EXISTS PermissionsAuditUpdateTrigger;
DROP TRIGGER IF EXISTS PermissionsAuditDeleteTrigger;

DROP TABLE IF EXISTS `VulnerabilityInfoAudit`;
DROP TRIGGER IF EXISTS VulnerabilityInfoAuditCreateTrigger;
DROP TRIGGER IF EXISTS VulnerabilityInfoAuditUpdateTrigger;
DROP TRIGGER IF EXISTS VulnerabilityInfoAuditDeleteTrigger;

DROP TRIGGER IF EXISTS SourceConfigAuditPreventDeleteTrigger;
DROP TRIGGER IF EXISTS JobConfigAuditPreventDeleteTrigger;
DROP TRIGGER IF EXISTS UsersAuditPreventDeleteTrigger;
DROP TRIGGER IF EXISTS PermissionsAuditPreventDeleteTrigger;
DROP TRIGGER IF EXISTS VulnerabilityInfoAuditPreventDeleteTrigger;