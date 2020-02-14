INSERT INTO DetectionStatus(Status, Name) VALUES ('vulnerable', 'Vulnerable'), ('exceptioned', 'Exceptioned'), ('dead host', 'Dead Host');

ALTER TABLE `Logs` MODIFY COLUMN JobHistoryID VARCHAR(40) NOT NULL;