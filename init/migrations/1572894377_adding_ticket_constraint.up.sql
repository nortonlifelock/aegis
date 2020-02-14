ALTER TABLE Ticket ADD CONSTRAINT ticket_unique_det_org_id UNIQUE (DetectionID, OrganizationID);
ALTER TABLE Organization DROP COLUMN PortDupl;