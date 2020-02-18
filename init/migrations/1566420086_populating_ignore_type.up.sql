ALTER TABLE `Ignore` DROP FOREIGN KEY fk_vm_ig_it;
ALTER TABLE IgnoreType MODIFY COLUMN ID INT NOT NULL;
ALTER TABLE `Ignore` ADD CONSTRAINT fk_vm_ig_it FOREIGN KEY (TypeID) REFERENCES IgnoreType(ID);

INSERT INTO IgnoreType(ID, Type, Name, CreatedBy)
VALUES
       ('0','EXCEPTION','Exception', 'Scaffold'),
       ('1','FALSEPOSITIVE','False Positive', 'Scaffold'),
       ('2','DECOMMASSET','Decommissioned Asset', 'Scaffold')
