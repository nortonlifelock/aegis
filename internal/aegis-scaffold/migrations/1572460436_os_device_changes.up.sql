ALTER TABLE Device ADD COLUMN OS VARCHAR(300) NULL AFTER OrganizationID;

ALTER TABLE OperatingSystemType DROP COLUMN CreatedBy;
ALTER TABLE OperatingSystemType DROP COLUMN UpdatedBy;
ALTER TABLE OperatingSystemType DROP COLUMN Name;

ALTER TABLE OperatingSystemType ADD COLUMN `Match` VARCHAR(100) NOT NULL AFTER `Type`;
ALTER TABLE OperatingSystemType ADD COLUMN `Priority` INT NOT NULL AFTER `Match`;

INSERT INTO OperatingSystemType(`Type`, `Match`, `Priority`) VALUE('Ubuntu', 'ubuntu', 3);
INSERT INTO OperatingSystemType(`Type`, `Match`, `Priority`) VALUE('AIX', 'aix', 5);
INSERT INTO OperatingSystemType(`Type`, `Match`, `Priority`) VALUE('Solaris', 'solaris', 5);
INSERT INTO OperatingSystemType(`Type`, `Match`, `Priority`) VALUE('Amazon Linux', 'amazon linux', 6);
INSERT INTO OperatingSystemType(`Type`, `Match`, `Priority`) VALUE('Redhat', 'red hat', 5);
INSERT INTO OperatingSystemType(`Type`, `Match`, `Priority`) VALUE('Redhat', 'redhat', 5);
INSERT INTO OperatingSystemType(`Type`, `Match`, `Priority`) VALUE('CentOS', 'centos', 5);
INSERT INTO OperatingSystemType(`Type`, `Match`, `Priority`) VALUE('IOS', 'cisco ios', 6);
INSERT INTO OperatingSystemType(`Type`, `Match`, `Priority`) VALUE('Cisco', 'cisco', 5);
INSERT INTO OperatingSystemType(`Type`, `Match`, `Priority`) VALUE('CoreOS', 'coreos', 5);
INSERT INTO OperatingSystemType(`Type`, `Match`, `Priority`) VALUE('Linux', 'fabric os', 5);
INSERT INTO OperatingSystemType(`Type`, `Match`, `Priority`) VALUE('NetApp', 'netapp', 5);
INSERT INTO OperatingSystemType(`Type`, `Match`, `Priority`) VALUE('Oracle', 'oracle', 5);
INSERT INTO OperatingSystemType(`Type`, `Match`, `Priority`) VALUE('UNIX', 'sunos', 1);
INSERT INTO OperatingSystemType(`Type`, `Match`, `Priority`) VALUE('Linux', 'linux', 1);
INSERT INTO OperatingSystemType(`Type`, `Match`, `Priority`) VALUE('Windows', 'windows', 1);
INSERT INTO OperatingSystemType(`Type`, `Match`, `Priority`) VALUE('VMware ESX/ESXi', 'vmware', 5);
INSERT INTO OperatingSystemType(`Type`, `Match`, `Priority`) VALUE('BIG-IP', 'big-ip', 5);
INSERT INTO OperatingSystemType(`Type`, `Match`, `Priority`) VALUE('Data ONTAP', 'data ontap', 5);
INSERT INTO OperatingSystemType(`Type`, `Match`, `Priority`) VALUE('z/OS', 'z/opsys', 5);
INSERT INTO OperatingSystemType(`Type`, `Match`, `Priority`) VALUE('BIP-IP', 'f5', 5);
INSERT INTO OperatingSystemType(`Type`, `Match`, `Priority`) VALUE('MAC OS X', 'os x', 1);
INSERT INTO OperatingSystemType(`Type`, `Match`, `Priority`) VALUE('MAC OS X', 'opsys x', 1);
INSERT INTO OperatingSystemType(`Type`, `Match`, `Priority`) VALUE('embedded', 'embedded', 5);
INSERT INTO OperatingSystemType(`Type`, `Match`, `Priority`) VALUE('embedded', 'vxworks', 5);

INSERT INTO OperatingSystemType(`Type`, `Match`, `Priority`) VALUE('Unknown', '', 0);