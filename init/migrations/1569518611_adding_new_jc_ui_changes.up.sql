-- 2019-07-31 1621 INTEL-585 changing job table to include sources.sql
ALTER TABLE Job ADD COLUMN SourceTypeOut int AFTER Struct;
ALTER TABLE Job ADD COLUMN SourceTypeIn int AFTER Struct;
ALTER TABLE Job ADD CONSTRAINT `fk_i_j_dist` FOREIGN KEY (`SourceTypeIn`) REFERENCES `SourceType` (`Id`);
ALTER TABLE Job ADD CONSTRAINT `fk_i_j_dost` FOREIGN KEY (`SourceTypeOut`) REFERENCES `SourceType` (`Id`);

-- 2019-08-05 1018 INTEL-628 changing ignore table to include createdby.sql
ALTER TABLE `Ignore` ADD COLUMN UpdatedBy varchar(255) AFTER Active;
ALTER TABLE `Ignore` ADD COLUMN CreatedBy varchar(255) AFTER Active;