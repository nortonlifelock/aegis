ALTER TABLE `SourceConfig` DROP COLUMN `AuthInfo`;

ALTER TABLE `SourceConfig` ADD COLUMN `Username` TEXT NOT NULL AFTER Port;
ALTER TABLE `SourceConfig` ADD COLUMN `Password` TEXT NOT NULL AFTER Username;
ALTER TABLE `SourceConfig` ADD COLUMN `PrivateKey` TEXT NULL AFTER Password;
ALTER TABLE `SourceConfig` ADD COLUMN `ConsumerKey` TEXT NULL AFTER PrivateKey;
ALTER TABLE `SourceConfig` ADD COLUMN `Token` TEXT NULL AFTER ConsumerKey;