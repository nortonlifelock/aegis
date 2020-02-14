DROP TABLE schema_migrations;

CREATE TABLE `schema_migrations` (
    `Name` VARCHAR(300) NOT NULL,
    `Date` DATETIME     NOT NULL DEFAULT NOW()
);

INSERT INTO schema_migrations (Name) VALUE ('1559762590_schema_initialize_single_file.up.sql');
INSERT INTO schema_migrations (Name) VALUE ('1559762591_open_source.up.sql');
INSERT INTO schema_migrations (Name) VALUE ('1560205281_addings_uuids.down.sql');
INSERT INTO schema_migrations (Name) VALUE ('1560205281_addings_uuids.up.sql');
INSERT INTO schema_migrations (Name) VALUE ('1560281326_removing_unused_tables.up.sql');
INSERT INTO schema_migrations (Name) VALUE ('1560793602_org_inheritance.down.sql');
INSERT INTO schema_migrations (Name) VALUE ('1560793602_org_inheritance.up.sql');
INSERT INTO schema_migrations (Name) VALUE ('1560973474_remove_unused_jh_fields.down.sql');
INSERT INTO schema_migrations (Name) VALUE ('1560973474_remove_unused_jh_fields.up.sql');
INSERT INTO schema_migrations (Name) VALUE ('1560977248_adding_json_auth_column.down.sql');
INSERT INTO schema_migrations (Name) VALUE ('1560977248_adding_json_auth_column.up.sql');
INSERT INTO schema_migrations (Name) VALUE ('1561402757_change_orgid_in_scan_summary_to_string.up.sql');
INSERT INTO schema_migrations (Name) VALUE ('1561481970_org_audit.down.sql');
INSERT INTO schema_migrations (Name) VALUE ('1561481970_org_audit.up.sql');
INSERT INTO schema_migrations (Name) VALUE ('1561562582_sc_jc_audit.down.sql');
INSERT INTO schema_migrations (Name) VALUE ('1561562582_sc_jc_audit.up.sql');
INSERT INTO schema_migrations (Name) VALUE ('1562772991_adding_new_cisrescan_job.down.sql');
INSERT INTO schema_migrations (Name) VALUE ('1562772991_adding_new_cisrescan_job.up.sql');
INSERT INTO schema_migrations (Name) VALUE ('1563900976_exception_audit.down.sql');
INSERT INTO schema_migrations (Name) VALUE ('1563900976_exception_audit.up.sql');
INSERT INTO schema_migrations (Name) VALUE ('1563901777_moving_js_payload_to_jc.down.sql');
INSERT INTO schema_migrations (Name) VALUE ('1563901777_moving_js_payload_to_jc.up.sql');
INSERT INTO schema_migrations (Name) VALUE ('1564511724_increasing_length_of_rs_payload.up.sql');
INSERT INTO schema_migrations (Name) VALUE ('1564610137_removing_kb_columns.up.sql');
INSERT INTO schema_migrations (Name) VALUE ('1565284419_changing_schema_migration.sql');