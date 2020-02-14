-- MySQL dump 10.13  Distrib 8.0.15, for osx10.13 (x86_64)
--
-- Host: localhost    Database: Vulnerability_manager
-- ------------------------------------------------------
-- Server version	8.0.15

/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
SET NAMES utf8mb4 ;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `Application`
--

DROP TABLE IF EXISTS `Application`;

CREATE TABLE `Application` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `OrganizationId` int(11) NOT NULL,
  PRIMARY KEY (`Id`),
  KEY `fk_vm_a_o` (`OrganizationId`),
  CONSTRAINT `fk_vm_a_o` FOREIGN KEY (`OrganizationId`) REFERENCES `Organization` (`Id`)
);

--
-- Table structure for table `AssignmentGroup`
--

DROP TABLE IF EXISTS `AssignmentGroup`;

CREATE TABLE `AssignmentGroup` (
  `SourceId` int(11) NOT NULL,
  `OrganizationId` int(11) NOT NULL,
  `IpAddress` varchar(20) NOT NULL,
  `GroupName` varchar(150) NOT NULL,
  `DBCreatedDate` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `DBUpdatedDate` datetime DEFAULT NULL,
  PRIMARY KEY (`SourceId`,`OrganizationId`,`IpAddress`),
  KEY `fk_vm_ag_o` (`OrganizationId`),
  CONSTRAINT `fk_vm_ag_o` FOREIGN KEY (`OrganizationId`) REFERENCES `Organization` (`Id`),
  CONSTRAINT `fk_vm_ag_s` FOREIGN KEY (`SourceId`) REFERENCES `Source` (`Id`)
);

--
-- Table structure for table `Category`
--

DROP TABLE IF EXISTS `Category`;

CREATE TABLE `Category` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `Category` varchar(50) NOT NULL,
  `ParentCategoryId` int(11) DEFAULT NULL,
  PRIMARY KEY (`Id`),
  KEY `fk_vm_c_pc` (`ParentCategoryId`),
  CONSTRAINT `fk_vm_c_pc` FOREIGN KEY (`ParentCategoryId`) REFERENCES `Category` (`Id`)
);

--
-- Table structure for table `DataType`
--

DROP TABLE IF EXISTS `DataType`;

CREATE TABLE `DataType` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `Type` varchar(20) NOT NULL,
  `Name` varchar(100) NOT NULL,
  PRIMARY KEY (`Id`)
);

--
-- Table structure for table `DBLog`
--

DROP TABLE IF EXISTS `DBLog`;

CREATE TABLE `DBLog` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `User` text NOT NULL,
  `Command` text NOT NULL,
  `Endpoint` text NOT NULL,
  `TimeOfUpdate` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`Id`)
);

--
-- Table structure for table `Detection`
--

DROP TABLE IF EXISTS `Detection`;

CREATE TABLE `Detection` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `OrganizationId` int(11) NOT NULL,
  `SourceId` int(11) NOT NULL,
  `DeviceId` int(11) NOT NULL,
  `VulnerabilityId` int(11) NOT NULL,
  `IgnoreId` int(11) DEFAULT NULL,
  `ParentDetectionId` int(11) DEFAULT NULL,
  `AlertDate` datetime NOT NULL,
  `Proof` text NOT NULL,
  `DetectionStatusId` int(11) NOT NULL,
  `TimesSeen` int(11) NOT NULL,
  `Created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `Updated` datetime DEFAULT NULL,
  `CreatedBy` varchar(50) DEFAULT NULL,
  `UpdatedBy` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`Id`),
  KEY `fk_vm_det_ds` (`DetectionStatusId`),
  KEY `fk_vm_det_i` (`IgnoreId`),
  KEY `fk_vm_det_d` (`DeviceId`),
  KEY `fk_vm_det_s` (`SourceId`),
  KEY `fk_vm_det_o` (`OrganizationId`),
  KEY `fk_vm_det_v` (`VulnerabilityId`),
  CONSTRAINT `fk_vm_det_d` FOREIGN KEY (`DeviceId`) REFERENCES `Device` (`Id`),
  CONSTRAINT `fk_vm_det_ds` FOREIGN KEY (`DetectionStatusId`) REFERENCES `DetectionStatus` (`Id`),
  CONSTRAINT `fk_vm_det_i` FOREIGN KEY (`IgnoreId`) REFERENCES `Ignore` (`Id`),
  CONSTRAINT `fk_vm_det_o` FOREIGN KEY (`OrganizationId`) REFERENCES `Organization` (`Id`),
  CONSTRAINT `fk_vm_det_s` FOREIGN KEY (`SourceId`) REFERENCES `Source` (`Id`),
  CONSTRAINT `fk_vm_det_v` FOREIGN KEY (`VulnerabilityId`) REFERENCES `VulnerabilityInfo` (`Id`)
);

--
-- Table structure for table `DetectionMetadata`
--

DROP TABLE IF EXISTS `DetectionMetadata`;

CREATE TABLE `DetectionMetadata` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `OrganizationId` int(11) NOT NULL,
  `SourceId` int(11) NOT NULL,
  `DetectionId` int(11) NOT NULL,
  `DataTypeId` int(11) NOT NULL,
  `StartDate` datetime DEFAULT NULL,
  `EndDate` datetime DEFAULT NULL,
  `Created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `Updated` datetime DEFAULT NULL,
  `CreatedBy` varchar(50) DEFAULT NULL,
  `UpdatedBy` varchar(50) DEFAULT NULL,
  `Active` bit(1) NOT NULL DEFAULT b'1',
  `ArchiveDate` datetime DEFAULT NULL,
  `Value` text NOT NULL,
  PRIMARY KEY (`Id`),
  KEY `fk_vm_dm_dt` (`DataTypeId`),
  KEY `fk_vm_dm_d` (`DetectionId`),
  KEY `fk_vm_dm_s` (`SourceId`),
  KEY `fk_vm_dm_o` (`OrganizationId`),
  CONSTRAINT `fk_vm_dm_d` FOREIGN KEY (`DetectionId`) REFERENCES `Detection` (`Id`),
  CONSTRAINT `fk_vm_dm_dt` FOREIGN KEY (`DataTypeId`) REFERENCES `Datatype` (`Id`),
  CONSTRAINT `fk_vm_dm_o` FOREIGN KEY (`OrganizationId`) REFERENCES `Organization` (`Id`),
  CONSTRAINT `fk_vm_dm_s` FOREIGN KEY (`SourceId`) REFERENCES `Source` (`Id`)
);

--
-- Table structure for table `DetectionPorts`
--

DROP TABLE IF EXISTS `DetectionPorts`;

CREATE TABLE `DetectionPorts` (
  `DetectionId` int(11) NOT NULL,
  `Port` int(11) NOT NULL,
  `ProtocolId` int(11) NOT NULL,
  `Created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `Updated` datetime DEFAULT NULL,
  `CreatedBy` varchar(50) DEFAULT NULL,
  `UpdatedBy` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`DetectionId`,`Port`,`ProtocolId`),
  KEY `fk_vm_dp_p` (`ProtocolId`),
  CONSTRAINT `fk_vm_dp_d` FOREIGN KEY (`DetectionId`) REFERENCES `Detection` (`Id`),
  CONSTRAINT `fk_vm_dp_p` FOREIGN KEY (`ProtocolId`) REFERENCES `Protocol` (`Id`)
);

--
-- Table structure for table `DetectionStatus`
--

DROP TABLE IF EXISTS `DetectionStatus`;

CREATE TABLE `DetectionStatus` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `Status` varchar(25) NOT NULL,
  `Name` varchar(200) NOT NULL,
  `Created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `Updated` datetime DEFAULT NULL,
  `CreatedBy` varchar(50) DEFAULT NULL,
  `UpdatedBy` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`Id`)
);

--
-- Table structure for table `Device`
--

DROP TABLE IF EXISTS `Device`;

CREATE TABLE `Device` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `AssetId` int(11) DEFAULT NULL,
  `Ip` varchar(32) NOT NULL,
  `InstanceId` varchar(128) NOT NULL,
  `GroupId` int(11) DEFAULT NULL,
  `NetworkId` varchar(128) DEFAULT NULL,
  `OrganizationId` int(11) NOT NULL,
  `OSTypeId` int(11) NOT NULL,
  `ImageId` int(11) DEFAULT NULL,
  `IsVirtual` bit(1) NOT NULL,
  PRIMARY KEY (`Id`),
  KEY `fk_vm_d_o` (`OrganizationId`),
  KEY `fk_vm_d_osti` (`OSTypeId`),
  KEY `fk_vm_d_di` (`ImageId`),
  CONSTRAINT `fk_vm_d_di` FOREIGN KEY (`ImageId`) REFERENCES `DeviceImage` (`Id`),
  CONSTRAINT `fk_vm_d_o` FOREIGN KEY (`OrganizationId`) REFERENCES `Organization` (`Id`),
  CONSTRAINT `fk_vm_d_osti` FOREIGN KEY (`OSTypeId`) REFERENCES `OperatingSystemType` (`Id`)
);

--
-- Table structure for table `DeviceGroup`
--

DROP TABLE IF EXISTS `DeviceGroup`;

CREATE TABLE `DeviceGroup` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `OrganizationId` int(11) NOT NULL,
  `GroupName` varchar(50) NOT NULL,
  PRIMARY KEY (`Id`),
  KEY `fk_vm_dg_o` (`OrganizationId`),
  CONSTRAINT `fk_vm_dg_o` FOREIGN KEY (`OrganizationId`) REFERENCES `Organization` (`Id`)
);

--
-- Table structure for table `DeviceGroupSource`
--

DROP TABLE IF EXISTS `DeviceGroupSource`;

CREATE TABLE `DeviceGroupSource` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `DeviceGroupId` int(11) NOT NULL,
  `SourceId` int(11) NOT NULL,
  `SourceIdentifier` varchar(20) NOT NULL,
  PRIMARY KEY (`Id`),
  KEY `fk_vm_dgs_s` (`SourceId`),
  KEY `fk_vm_dgs_dg` (`DeviceGroupId`),
  CONSTRAINT `fk_vm_dgs_dg` FOREIGN KEY (`DeviceGroupId`) REFERENCES `DeviceGroup` (`Id`),
  CONSTRAINT `fk_vm_dgs_s` FOREIGN KEY (`SourceId`) REFERENCES `Source` (`Id`)
);

--
-- Table structure for table `DeviceImage`
--

DROP TABLE IF EXISTS `DeviceImage`;

CREATE TABLE `DeviceImage` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `Name` varchar(150) NOT NULL,
  PRIMARY KEY (`Id`)
);

--
-- Table structure for table `DeviceInfo`
--

DROP TABLE IF EXISTS `DeviceInfo`;

CREATE TABLE `DeviceInfo` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `OrganizationId` int(11) NOT NULL,
  `SourceId` int(11) NOT NULL,
  `DeviceId` int(11) NOT NULL,
  `Created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `Updated` datetime DEFAULT NULL,
  `CreatedBy` varchar(50) DEFAULT NULL,
  `UpdatedBy` varchar(50) DEFAULT NULL,
  `Active` bit(1) NOT NULL DEFAULT b'1',
  `FoundActive` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `OperatingSystem` varchar(200) NOT NULL,
  `SourceIdentifier` int(11) NOT NULL,
  PRIMARY KEY (`Id`),
  KEY `fk_vm_di_d` (`DeviceId`),
  KEY `fk_vm_di_s` (`SourceId`),
  KEY `fk_vm_di_o` (`OrganizationId`),
  CONSTRAINT `fk_vm_di_d` FOREIGN KEY (`DeviceId`) REFERENCES `Device` (`Id`),
  CONSTRAINT `fk_vm_di_o` FOREIGN KEY (`OrganizationId`) REFERENCES `Organization` (`Id`),
  CONSTRAINT `fk_vm_di_s` FOREIGN KEY (`SourceId`) REFERENCES `Source` (`Id`)
);

--
-- Table structure for table `DeviceMetadata`
--

DROP TABLE IF EXISTS `DeviceMetadata`;

CREATE TABLE `DeviceMetadata` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `OrganizationId` int(11) NOT NULL,
  `SourceId` int(11) DEFAULT NULL,
  `DeviceId` int(11) NOT NULL,
  `DeviceMetadataTypeId` int(11) NOT NULL,
  `DataTypeId` int(11) NOT NULL,
  `StartDate` datetime DEFAULT NULL,
  `EndDate` datetime DEFAULT NULL,
  `Created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `Updated` datetime DEFAULT NULL,
  `CreatedBy` varchar(50) DEFAULT NULL,
  `UpdatedBy` varchar(50) DEFAULT NULL,
  `ACTIVE` bit(1) NOT NULL DEFAULT b'1',
  PRIMARY KEY (`Id`),
  KEY `fk_vm_dmd_o` (`OrganizationId`),
  KEY `fk_vm_dmd_dmdt` (`DeviceMetadataTypeId`),
  KEY `fk_vm_dmd_dt` (`DataTypeId`),
  KEY `fk_vm_dmd_d` (`DeviceId`),
  KEY `fk_vm_dmd_s` (`SourceId`),
  CONSTRAINT `fk_vm_dmd_d` FOREIGN KEY (`DeviceId`) REFERENCES `Device` (`Id`),
  CONSTRAINT `fk_vm_dmd_dmdt` FOREIGN KEY (`DeviceMetadataTypeId`) REFERENCES `DeviceMetadataType` (`Id`),
  CONSTRAINT `fk_vm_dmd_dt` FOREIGN KEY (`DataTypeId`) REFERENCES `DataType` (`Id`),
  CONSTRAINT `fk_vm_dmd_o` FOREIGN KEY (`OrganizationId`) REFERENCES `Organization` (`Id`),
  CONSTRAINT `fk_vm_dmd_s` FOREIGN KEY (`SourceId`) REFERENCES `Source` (`Id`)
);

--
-- Table structure for table `DeviceMetadataType`
--

DROP TABLE IF EXISTS `DeviceMetadataType`;

CREATE TABLE `DeviceMetadataType` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `Type` varchar(20) NOT NULL,
  `Name` varchar(100) NOT NULL,
  PRIMARY KEY (`Id`)
);

--
-- Table structure for table `Ignore`
--

DROP TABLE IF EXISTS `Ignore`;

CREATE TABLE `Ignore` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `OrganizationId` int(11) NOT NULL,
  `VulnerabilityId` int(11) NOT NULL,
  `DeviceId` int(11) NOT NULL,
  `IgnoreTypeId` int(11) NOT NULL,
  `Due` datetime NOT NULL,
  `Active` bit(1) NOT NULL DEFAULT b'1',
  `ApprovalReference` varchar(150) NOT NULL,
  `ApprovalDate` datetime NOT NULL,
  `ApprovalSourceId` int(11) DEFAULT NULL,
  `Created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `Updated` datetime DEFAULT NULL,
  `CreatedBy` varchar(50) DEFAULT NULL,
  `UpdatedBy` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`Id`),
  KEY `fk_vm_i_it` (`IgnoreTypeId`),
  KEY `fk_vm_i_d` (`DeviceId`),
  KEY `fk_vm_i_v` (`VulnerabilityId`),
  KEY `fk_vm_i_o` (`OrganizationId`),
  CONSTRAINT `fk_vm_i_d` FOREIGN KEY (`DeviceId`) REFERENCES `Device` (`Id`),
  CONSTRAINT `fk_vm_i_it` FOREIGN KEY (`IgnoreTypeId`) REFERENCES `IgnoreType` (`Id`),
  CONSTRAINT `fk_vm_i_o` FOREIGN KEY (`OrganizationId`) REFERENCES `Organization` (`Id`),
  CONSTRAINT `fk_vm_i_v` FOREIGN KEY (`VulnerabilityId`) REFERENCES `Vulnerability` (`Id`)
);

--
-- Table structure for table `IgnoreType`
--

DROP TABLE IF EXISTS `IgnoreType`;

CREATE TABLE `IgnoreType` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `Type` varchar(25) NOT NULL,
  `Name` varchar(200) NOT NULL,
  `Created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `Updated` datetime DEFAULT NULL,
  `CreatedBy` varchar(50) DEFAULT NULL,
  `UpdatedBy` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`Id`)
);

--
-- Table structure for table `Job`
--

DROP TABLE IF EXISTS `Job`;

CREATE TABLE `Job` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `Struct` varchar(50) NOT NULL,
  `Priority` int(11) NOT NULL,
  `CreatedDate` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `CreatedBy` varchar(255) NOT NULL,
  `UpdatedDate` datetime DEFAULT NULL,
  `UpdatedBy` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`Id`)
);

--
-- Table structure for table `JobConfig`
--

DROP TABLE IF EXISTS `JobConfig`;

CREATE TABLE `JobConfig` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `JobId` int(11) NOT NULL,
  `OrganizationId` int(11) NOT NULL,
  `DataInSourceConfigId` int(11) DEFAULT NULL,
  `DataOutSourceConfigId` int(11) DEFAULT NULL,
  `Payload` text,
  `PriorityOverride` int(11) DEFAULT NULL,
  `Continuous` bit(1) NOT NULL DEFAULT b'0',
  `WaitInSeconds` int(11) NOT NULL DEFAULT '600',
  `MaxInstances` int(11) NOT NULL DEFAULT '1',
  `AutoStart` bit(1) NOT NULL DEFAULT b'0',
  `CreatedDate` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `CreatedBy` varchar(255) NOT NULL,
  `UpdatedDate` datetime DEFAULT NULL,
  `UpdatedBy` varchar(255) DEFAULT NULL,
  `Active` bit(1) NOT NULL DEFAULT b'1',
  `LastJobStart` datetime DEFAULT NULL,
  PRIMARY KEY (`Id`),
  KEY `fk_vm_jc_j` (`JobId`),
  KEY `fk_vm_jc_o` (`OrganizationId`),
  KEY `fk_vm_jc_sci` (`DataInSourceConfigId`),
  KEY `fk_vm_jc_sco` (`DataOutSourceConfigId`),
  CONSTRAINT `fk_vm_jc_j` FOREIGN KEY (`JobId`) REFERENCES `Job` (`Id`),
  CONSTRAINT `fk_vm_jc_o` FOREIGN KEY (`OrganizationId`) REFERENCES `Organization` (`Id`),
  CONSTRAINT `fk_vm_jc_sci` FOREIGN KEY (`DataInSourceConfigId`) REFERENCES `SourceConfig` (`Id`),
  CONSTRAINT `fk_vm_jc_sco` FOREIGN KEY (`DataOutSourceConfigId`) REFERENCES `SourceConfig` (`Id`)
);

--
-- Table structure for table `JobHistory`
--

DROP TABLE IF EXISTS `JobHistory`;

CREATE TABLE `JobHistory` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `JobId` int(11) NOT NULL,
  `ConfigId` int(11) NOT NULL,
  `StatusId` int(11) NOT NULL DEFAULT '0',
  `ParentJobId` int(11) DEFAULT NULL,
  `Identifier` varchar(100) DEFAULT NULL,
  `Priority` int(11) NOT NULL,
  `ManualStart` bit(1) NOT NULL DEFAULT b'0',
  `Continuous` bit(1) NOT NULL DEFAULT b'0',
  `WaitInSeconds` int(11) NOT NULL DEFAULT '600',
  `CurrentIteration` int(11) DEFAULT NULL,
  `Payload` text NOT NULL,
  `ThreadId` varchar(100) DEFAULT NULL,
  `PulseDate` datetime DEFAULT NULL,
  `CreatedDate` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `CreatedBy` varchar(255) DEFAULT NULL,
  `UpdatedDate` datetime DEFAULT NULL,
  `UpdatedBy` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`Id`),
  KEY `fk_vm_jh_j` (`JobId`),
  KEY `fk_vm_jh_jc` (`ConfigId`),
  KEY `fk_vm_jh_js` (`StatusId`),
  KEY `fk_vm_jh_pjh` (`ParentJobId`),
  CONSTRAINT `fk_vm_jh_j` FOREIGN KEY (`JobId`) REFERENCES `Job` (`Id`),
  CONSTRAINT `fk_vm_jh_jc` FOREIGN KEY (`ConfigId`) REFERENCES `JobConfig` (`Id`),
  CONSTRAINT `fk_vm_jh_js` FOREIGN KEY (`StatusId`) REFERENCES `JobStatus` (`Id`),
  CONSTRAINT `fk_vm_jh_pjh` FOREIGN KEY (`ParentJobId`) REFERENCES `JobHistory` (`Id`)
);

--
-- Table structure for table `JobSchedule`
--

DROP TABLE IF EXISTS `JobSchedule`;

CREATE TABLE `JobSchedule` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `ConfigId` int(11) NOT NULL,
  `DaysOfWeek` varchar(20) DEFAULT NULL,
  `TimeOfDay` time DEFAULT NULL,
  `Payload` text,
  `LastRun` datetime DEFAULT NULL,
  `CreatedDate` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `CreatedBy` varchar(255) NOT NULL,
  `UpdatedDate` datetime DEFAULT NULL,
  `UpdatedBy` varchar(255) DEFAULT NULL,
  `Active` bit(1) DEFAULT b'1',
  PRIMARY KEY (`Id`),
  KEY `fk_vm_js_jc` (`ConfigId`),
  CONSTRAINT `fk_vm_js_jc` FOREIGN KEY (`ConfigId`) REFERENCES `JobConfig` (`Id`)
);

--
-- Table structure for table `JobStatus`
--

DROP TABLE IF EXISTS `JobStatus`;

CREATE TABLE `JobStatus` (
  `Id` int(11) NOT NULL,
  `Status` varchar(25) NOT NULL,
  PRIMARY KEY (`Id`)
);

--
-- Table structure for table `Logs`
--

DROP TABLE IF EXISTS `Logs`;

CREATE TABLE `Logs` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `TypeId` int(11) NOT NULL,
  `Log` text NOT NULL,
  `Error` text NOT NULL,
  `CreateDate` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `JobHistoryId` int(11) DEFAULT NULL,
  PRIMARY KEY (`Id`),
  KEY `fk_vm_l_lt` (`TypeId`),
  CONSTRAINT `fk_vm_l_lt` FOREIGN KEY (`TypeId`) REFERENCES `LogType` (`Id`)
);

--
-- Table structure for table `LogType`
--

DROP TABLE IF EXISTS `LogType`;

CREATE TABLE `LogType` (
  `Id` int(11) NOT NULL,
  `Type` varchar(15) NOT NULL,
  `Name` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`Id`)
);

--
-- Table structure for table `OperatingSystemType`
--

DROP TABLE IF EXISTS `OperatingSystemType`;

CREATE TABLE `OperatingSystemType` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `Type` varchar(255) DEFAULT NULL,
  `Name` varchar(255) DEFAULT NULL,
  `Created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `Updated` datetime DEFAULT NULL,
  `CreatedBy` varchar(50) DEFAULT NULL,
  `UpdatedBy` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`Id`)
);

--
-- Table structure for table `Organization`
--

DROP TABLE IF EXISTS `Organization`;

CREATE TABLE `Organization` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `Code` varchar(20) NOT NULL,
  `Description` varchar(100) NOT NULL,
  `Payload` text,
  `TimeZoneOffset` float NOT NULL,
  `Created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `Updated` datetime DEFAULT NULL,
  `CreatedBy` varchar(50) DEFAULT NULL,
  `UpdatedBy` varchar(50) DEFAULT NULL,
  `PortDupl` bit(1) NOT NULL DEFAULT b'0',
  `Active` bit(1) NOT NULL DEFAULT b'1',
  PRIMARY KEY (`Id`)
);

--
-- Table structure for table `OwnerGroup`
--

DROP TABLE IF EXISTS `OwnerGroup`;

CREATE TABLE `OwnerGroup` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `OrganizationId` int(11) NOT NULL,
  PRIMARY KEY (`Id`),
  KEY `fk_vm_og_o` (`OrganizationId`),
  CONSTRAINT `fk_vm_og_o` FOREIGN KEY (`OrganizationId`) REFERENCES `Organization` (`Id`)
);

--
-- Table structure for table `Permissions`
--

DROP TABLE IF EXISTS `Permissions`;

CREATE TABLE `Permissions` (
  `UserId` int(11) NOT NULL,
  `OrgId` int(11) NOT NULL,
  `CanUpdateJob` bit(1) NOT NULL DEFAULT b'0',
  `CanDeleteJob` bit(1) NOT NULL DEFAULT b'0',
  `CanCreateJob` bit(1) NOT NULL DEFAULT b'0',
  `CanUpdateConfig` bit(1) NOT NULL DEFAULT b'0',
  `CanDeleteConfig` bit(1) NOT NULL DEFAULT b'0',
  `CanCreateConfig` bit(1) NOT NULL DEFAULT b'0',
  `CanUpdateSource` bit(1) NOT NULL DEFAULT b'0',
  `CanDeleteSource` bit(1) NOT NULL DEFAULT b'0',
  `CanCreateSource` bit(1) NOT NULL DEFAULT b'0',
  `CanUpdateOrg` bit(1) NOT NULL DEFAULT b'0',
  `CanDeleteOrg` bit(1) NOT NULL DEFAULT b'0',
  `CanCreateOrg` bit(1) NOT NULL DEFAULT b'0',
  `CanReadJobHistories` bit(1) NOT NULL DEFAULT b'0',
  `CanRegisterUser` bit(1) NOT NULL DEFAULT b'0',
  `CanUpdateUser` bit(1) NOT NULL DEFAULT b'0',
  `CanDeleteUser` bit(1) NOT NULL DEFAULT b'0',
  `CanReadUser` bit(1) NOT NULL DEFAULT b'0',
  `CanBulkUpdate` bit(1) NOT NULL DEFAULT b'0',
  `CanManageTags` bit(1) NOT NULL DEFAULT b'0',
  PRIMARY KEY (`UserId`,`OrgId`),
  CONSTRAINT `fk_vm_p_u` FOREIGN KEY (`UserId`) REFERENCES `Users` (`Id`)
);

--
-- Table structure for table `Protocol`
--

DROP TABLE IF EXISTS `Protocol`;

CREATE TABLE `Protocol` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `Protocol` varchar(20) NOT NULL,
  PRIMARY KEY (`Id`)
);

--
-- Table structure for table `ReferenceType`
--

DROP TABLE IF EXISTS `ReferenceType`;

CREATE TABLE `ReferenceType` (
  `Id` int(11) NOT NULL,
  `Type` varchar(30) NOT NULL,
  `DBCreatedDate` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `DBUpdatedDate` datetime DEFAULT NULL,
  PRIMARY KEY (`Id`)
);

--
-- Table structure for table `RefType`
--

DROP TABLE IF EXISTS `RefType`;

CREATE TABLE `RefType` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `Type` varchar(25) NOT NULL,
  `Name` varchar(200) NOT NULL,
  `Created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `Updated` datetime DEFAULT NULL,
  `CreatedBy` varchar(50) DEFAULT NULL,
  `UpdatedBy` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`Id`)
);

--
-- Table structure for table `ScanSummary`
--

DROP TABLE IF EXISTS `ScanSummary`;

CREATE TABLE `ScanSummary` (
  `SourceId` int(11) NOT NULL,
  `TemplateId` varchar(100) DEFAULT NULL,
  `OrgId` int(11) NOT NULL,
  `ScanIdentifier` int(11) NOT NULL,
  `SourceKey` varchar(100) NOT NULL,
  `ScanStatus` varchar(30) DEFAULT NULL,
  `ParentJobId` int(11) NOT NULL,
  `ScanClosePayload` text NOT NULL,
  `CreatedDate` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `UpdatedDate` datetime DEFAULT NULL,
  PRIMARY KEY (`SourceId`,`OrgId`,`ScanIdentifier`,`SourceKey`)
);

--
-- Table structure for table `Source`
--

DROP TABLE IF EXISTS `Source`;

CREATE TABLE `Source` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `SourceTypeId` int(11) NOT NULL,
  `Source` varchar(100) NOT NULL,
  `Created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `Updated` datetime DEFAULT NULL,
  `CreatedBy` varchar(50) DEFAULT NULL,
  `UpdatedBy` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`Id`),
  KEY `fk_vm_s_st` (`SourceTypeId`),
  CONSTRAINT `fk_vm_s_st` FOREIGN KEY (`SourceTypeId`) REFERENCES `SourceType` (`Id`)
);

--
-- Table structure for table `SourceConfig`
--

DROP TABLE IF EXISTS `SourceConfig`;

CREATE TABLE `SourceConfig` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `Source` varchar(100) DEFAULT NULL,
  `SourceId` int(11) NOT NULL,
  `OrganizationId` int(11) NOT NULL,
  `Address` varchar(100) DEFAULT NULL,
  `Port` varchar(30) DEFAULT NULL,
  `Username` text NOT NULL,
  `Password` text NOT NULL,
  `PrivateKey` text,
  `ConsumerKey` text,
  `Token` text,
  `DBCreatedDate` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `DBUpdatedDate` datetime DEFAULT NULL,
  `Payload` text NOT NULL,
  `Active` bit(1) NOT NULL DEFAULT b'1',
  `UpdatedBy` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`Id`),
  KEY `fk_vm_sc_o` (`OrganizationId`),
  KEY `fk_vm_sc_sc` (`SourceId`),
  CONSTRAINT `fk_vm_sc_o` FOREIGN KEY (`OrganizationId`) REFERENCES `Organization` (`Id`),
  CONSTRAINT `fk_vm_sc_sc` FOREIGN KEY (`SourceId`) REFERENCES `Source` (`Id`)
);

--
-- Table structure for table `SourceType`
--

DROP TABLE IF EXISTS `SourceType`;

CREATE TABLE `SourceType` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `Type` varchar(25) NOT NULL,
  `Created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `Updated` datetime DEFAULT NULL,
  `CreatedBy` varchar(50) DEFAULT NULL,
  `UpdatedBy` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`Id`)
);

--
-- Table structure for table `Subscription`
--

DROP TABLE IF EXISTS `Subscription`;

CREATE TABLE `Subscription` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `OwnerGroupId` int(11) NOT NULL,
  `SubscriptionTypeId` int(11) NOT NULL,
  PRIMARY KEY (`Id`),
  KEY `fk_vm_su_st` (`SubscriptionTypeId`),
  KEY `fk_vm_su_og` (`OwnerGroupId`),
  CONSTRAINT `fk_vm_su_og` FOREIGN KEY (`OwnerGroupId`) REFERENCES `OwnerGroup` (`Id`),
  CONSTRAINT `fk_vm_su_st` FOREIGN KEY (`SubscriptionTypeId`) REFERENCES `SubscriptionType` (`Id`)
);

--
-- Table structure for table `SubscriptionType`
--

DROP TABLE IF EXISTS `SubscriptionType`;

CREATE TABLE `SubscriptionType` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `Type` varchar(25) NOT NULL,
  `Name` varchar(200) NOT NULL,
  `Created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `Updated` datetime DEFAULT NULL,
  `CreatedBy` varchar(50) DEFAULT NULL,
  `UpdatedBy` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`Id`)
);

--
-- Table structure for table `Tag`
--

DROP TABLE IF EXISTS `Tag`;

CREATE TABLE `Tag` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `DeviceId` int(11) NOT NULL,
  `TagKeyId` int(11) NOT NULL,
  `Value` varchar(255) NOT NULL,
  PRIMARY KEY (`Id`),
  KEY `fk_vm_t_dev` (`DeviceId`),
  KEY `fk_vm_t_tagk` (`TagKeyId`),
  CONSTRAINT `fk_vm_t_dev` FOREIGN KEY (`DeviceId`) REFERENCES `Device` (`Id`),
  CONSTRAINT `fk_vm_t_tagk` FOREIGN KEY (`TagKeyId`) REFERENCES `TagKey` (`Id`)
);

--
-- Table structure for table `TagKey`
--

DROP TABLE IF EXISTS `TagKey`;

CREATE TABLE `TagKey` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `KeyValue` varchar(255) NOT NULL,
  PRIMARY KEY (`Id`)
);

--
-- Table structure for table `TagMap`
--

DROP TABLE IF EXISTS `TagMap`;

CREATE TABLE `TagMap` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `TicketingSourceId` int(11) NOT NULL,
  `TicketingTag` varchar(255) NOT NULL,
  `CloudSourceId` int(11) NOT NULL,
  `CloudTag` varchar(255) NOT NULL,
  `Options` varchar(255) NOT NULL,
  `Active` tinyint(1) NOT NULL DEFAULT '1',
  `OrganizationId` int(11) NOT NULL,
  PRIMARY KEY (`Id`),
  KEY `pk_vm_ts_s` (`TicketingSourceId`),
  KEY `pk_vm_cs_s` (`CloudSourceId`),
  KEY `fk_vm_ts_o` (`OrganizationId`),
  CONSTRAINT `fk_vm_ts_o` FOREIGN KEY (`OrganizationId`) REFERENCES `Organization` (`Id`),
  CONSTRAINT `pk_vm_cs_s` FOREIGN KEY (`CloudSourceId`) REFERENCES `Source` (`Id`),
  CONSTRAINT `pk_vm_ts_s` FOREIGN KEY (`TicketingSourceId`) REFERENCES `Source` (`Id`)
);

--
-- Table structure for table `Users`
--

DROP TABLE IF EXISTS `Users`;

CREATE TABLE `Users` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `UUID` varchar(255) NOT NULL,
  `Username` text,
  `FirstName` text NOT NULL,
  `LastName` text NOT NULL,
  `Email` text NOT NULL,
  `IsDisabled` bit(1) NOT NULL DEFAULT b'0',
  PRIMARY KEY (`Id`)
);

--
-- Table structure for table `UserSession`
--

DROP TABLE IF EXISTS `UserSession`;

CREATE TABLE `UserSession` (
  `UserId` int(11) NOT NULL,
  `OrgId` int(11) NOT NULL,
  `SessionKey` text NOT NULL,
  `IsDisabled` bit(1) NOT NULL DEFAULT b'0',
  `LoginTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `LastSeenTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`UserId`),
  KEY `fk_vm_us_o` (`OrgId`),
  CONSTRAINT `fk_vm_us_o` FOREIGN KEY (`OrgId`) REFERENCES `Organization` (`Id`),
  CONSTRAINT `fk_vm_us_u` FOREIGN KEY (`UserId`) REFERENCES `Users` (`Id`)
);

--
-- Table structure for table `Vulnerability`
--

DROP TABLE IF EXISTS `Vulnerability`;

CREATE TABLE `Vulnerability` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `Title` varchar(500) NOT NULL,
  `Summary` text NOT NULL,
  PRIMARY KEY (`Id`)
);

--
-- Table structure for table `Vulnerability_Reference`
--

DROP TABLE IF EXISTS `Vulnerability_Reference`;

CREATE TABLE `Vulnerability_Reference` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `VulnInfoId` int(11) NOT NULL,
  `SourceId` int(11) NOT NULL,
  `Reference` varchar(2083) DEFAULT NULL,
  `RefType` int(11) DEFAULT NULL,
  PRIMARY KEY (`Id`),
  KEY `fk_vm_vr_vi` (`VulnInfoId`),
  KEY `fk_vm_vr_s` (`SourceId`),
  CONSTRAINT `fk_vm_vr_s` FOREIGN KEY (`SourceId`) REFERENCES `Source` (`Id`),
  CONSTRAINT `fk_vm_vr_vi` FOREIGN KEY (`VulnInfoId`) REFERENCES `VulnerabilityInfo` (`Id`)
);

--
-- Table structure for table `VulnerabilityInfo`
--

DROP TABLE IF EXISTS `VulnerabilityInfo`;

CREATE TABLE `VulnerabilityInfo` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `SourceVulnId` varchar(255) NOT NULL,
  `Title` text,
  `VulnerabilityId` int(11) DEFAULT NULL,
  `SourceId` int(11) NOT NULL,
  `CVSS` float NOT NULL,
  `CVSS3` float NOT NULL,
  `Description` longtext NOT NULL,
  `Solution` mediumtext NOT NULL,
  `Severity` int(11) NOT NULL,
  `CVSSVector` text NOT NULL,
  `CVSS3Vector` text NOT NULL,
  `MatchConfidence` int(11) DEFAULT NULL,
  `MatchReasons` text,
  `Software` text,
  `DetectionInformation` text,
  `Updated` datetime DEFAULT NULL,
  `Created` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`Id`),
  KEY `fk_vm_vi_v` (`VulnerabilityId`),
  KEY `fk_vm_vi_s` (`SourceId`),
  CONSTRAINT `fk_vm_vi_s` FOREIGN KEY (`SourceId`) REFERENCES `Source` (`Id`),
  CONSTRAINT `fk_vm_vi_v` FOREIGN KEY (`VulnerabilityId`) REFERENCES `Vulnerability` (`Id`)
);

ALTER TABLE Detection DROP FOREIGN KEY fk_vm_det_i;

DROP TABLE IF EXISTS `Ignore`;

CREATE TABLE `Ignore`(
    Id              INT          NOT NULL AUTO_INCREMENT,
    SourceId        INT          NOT NULL,
    OrganizationId  INT          NOT NULL,
    TypeId          INT          NOT NULL,
    VulnerabilityId VARCHAR(120) NOT NULL,
    DeviceId        INT          NOT NULL,
    DueDate         DATETIME     NULL,
    Approval        VARCHAR(120) NOT NULL,
    Active          BIT(1)       NOT NULL,
    Port            VARCHAR(15)  NULL,
    DBCreatedDate   DATETIME     NOT NULL DEFAULT NOW(),
    DBUpdatedDate   DATETIME     NULL,

    CONSTRAINT pk_vm_Ignore PRIMARY KEY (Id),
    CONSTRAINT fk_vm_ig_s   FOREIGN KEY (SourceId)       REFERENCES Source (Id),
    CONSTRAINT fk_vm_ig_o   FOREIGN KEY (OrganizationId) REFERENCES Organization (Id),
    CONSTRAINT fk_vm_ig_it  FOREIGN KEY (TypeId)         REFERENCES IgnoreType (Id)
);

ALTER TABLE Detection ADD CONSTRAINT fk_vm_det_i FOREIGN KEY (IgnoreId) REFERENCES `Ignore`(Id);

/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;