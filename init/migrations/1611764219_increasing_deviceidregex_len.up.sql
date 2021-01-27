alter table CISAssignmentRules Modify Column DeviceIDRegex VARCHAR(300) NULL;
alter table `Ignore` modify column DeviceIDRegex VARCHAR(200) NULL;
alter table `IgnoreAudit` modify column DeviceIDRegex VARCHAR(300) NULL;
