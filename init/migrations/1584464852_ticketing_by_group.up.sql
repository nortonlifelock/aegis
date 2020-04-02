ALTER TABLE AssetGroup ADD COLUMN LastTicketing DATETIME NULL AFTER CloudSourceID;

UPDATE AssetGroup AG SET AG.LastTicketing = (SELECT LastJobStart from JobConfig JC WHERE JC.OrganizationID = AG.OrganizationID AND JC.JobID = (SELECT ID FROM Job WHERE Struct = 'TicketingJob') LIMIT 1);