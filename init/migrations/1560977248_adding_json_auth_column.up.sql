ALTER TABLE `SourceConfig` ADD COLUMN AuthInfo JSON NOT NULL AFTER `Port`;

UPDATE `SourceConfig` SET SourceConfig.AuthInfo = JSON_OBJECT(
        "PrivateKey",
        SourceConfig.PrivateKey,
        "ConsumerKey",
        SourceConfig.ConsumerKey,
        "Token",
        SourceConfig.Token
    ) WHERE SourceConfig.PrivateKey IS NOT NULL AND SourceConfig.ConsumerKey IS NOT NULL AND SourceConfig.Token IS NOT NULL;

UPDATE `SourceConfig` SET SourceConfig.AuthInfo = JSON_OBJECT(
        "Username",
        SourceConfig.Username,
        "Password",
        SourceConfig.Password
    ) WHERE SourceConfig.Username IS NOT NULL AND SourceConfig.Password IS NOT NULL AND SourceConfig.PrivateKey IS NULL;

ALTER TABLE `SourceConfig` DROP COLUMN `Username`;
ALTER TABLE `SourceConfig` DROP COLUMN `Password`;
ALTER TABLE `SourceConfig` DROP COLUMN `PrivateKey`;
ALTER TABLE `SourceConfig` DROP COLUMN `ConsumerKey`;
ALTER TABLE `SourceConfig` DROP COLUMN `Token`;