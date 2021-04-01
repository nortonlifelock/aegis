ALTER TABLE Detection MODIFY COLUMN ParentDetectionId VARCHAR(36) NULL;
ALTER TABLE Detection ADD CONSTRAINT fk_vm_det_pdet FOREIGN KEY (ParentDetectionId) REFERENCES `Detection`(Id);