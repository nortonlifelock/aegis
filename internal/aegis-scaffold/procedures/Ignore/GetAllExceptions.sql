/*
  RETURN Ignore
  ID              VARCHAR(30)     NOT
  SourceID        VARCHAR(36)     NOT
  OrganizationID  VARCHAR(36)     NOT
  TypeID          INT             NOT
  VulnerabilityID NVARCHAR(120)   NOT
  DeviceID        VARCHAR(36)     NOT
  DueDate         DATETIME        NULL
  Approval        NVARCHAR(120)   NOT
  Active          BIT             NOT
  Port            NVARCHAR(120)   NOT
  CreatedBy       NVARCHAR(255)   NULL
  UpdatedBy       NVARCHAR(255)   NULL
  DBCreatedDate   DATETIME        NOT
  DBUpdatedDate   DATETIME        NULL
 */

DROP PROCEDURE IF EXISTS ```GetAllExceptions`;

CREATE PROCEDURE ```GetAllExceptions`( _offset INT, _limit INT,_sourceID VARCHAR(36), _orgID VARCHAR(36), _typeID INT,_vulnID NVARCHAR(255), _devID VARCHAR(36), _dueDate DATETIME, _port NVARCHAR(15), _approval  NVARCHAR(120),
                                             _active BIT, _dBCreatedDate DATETIME,_dBUpdatedDate DATETIME,_updatedBy NVARCHAR(255), _createdBy NVARCHAR(255),_sortField NVARCHAR(255),
                                             _sortOrder NVARCHAR(255))
    #BEGIN#
BEGIN
    SELECT
        I.ID,
        I.SourceID,
        I.OrganizationID,
        I.TypeID,
        I.VulnerabilityID,
        I.DeviceID ,
        I.DueDate,
        I.Approval,
        I.Active,
        I.Port,
        I.CreatedBy,
        I.UpdatedBy,
        I.DBCreatedDate,
        I.DBUpdatedDate
    FROM `Ignore` I
             JOIN Source Sc on I.SourceID=Sc.ID
             JOIN IgnoreType It on I.TypeID=It.ID
    WHERE I.OrganizationID = _OrgID
      AND (I.SourceID = _sourceID OR _sourceID = '' OR _sourceID is NULL)
      AND (I.TypeID = _typeID OR _typeID= 0 OR _typeID is NULL)
      AND (I.VulnerabilityID = _vulnID OR _vulnID = '' OR _vulnID is NULL)
      AND (I.DeviceID = _devID OR _devID ='' OR _devID is NULL)
      AND (I.DueDate = _dueDate OR _dueDate ='1970-01-02 00:00:00 +0000 UTC' OR _dueDate is NULL)
      AND (I.Port = _port OR _port ='' OR _port is NULL)
      AND (I.Approval= _approval OR _approval ='' OR _approval is NULL)
      AND (I.Active = _active OR _active ='' OR _active is NULL)
      AND (I.UpdatedBy = _updatedBy OR _updatedBy ='' OR _updatedBy is NULL)
      AND (I.CreatedBy = _createdBy OR _createdBy ='' OR _createdBy is NULL)
      AND (I.DBCreatedDate = _dBCreatedDate OR _dBCreatedDate ='1970-01-02 00:00:00 +0000 UTC' OR _dBCreatedDate is NULL)
      AND (I.DBUpdatedDate = _dBUpdatedDate OR _dBUpdatedDate ='1970-01-02 00:00:00 +0000 UTC' OR _dBUpdatedDate is NULL)
    ORDER BY
        CASE WHEN _sortField = 'source_id' AND _sortOrder='ASC' THEN Sc.Source END,
        CASE WHEN _sortField = 'source_id' AND _sortOrder='DESC' THEN Sc.Source END DESC,
        CASE WHEN _sortField = 'type_id' AND _sortOrder='ASC' THEN It.Name END,
        CASE WHEN _sortField = 'type_id' AND _sortOrder='DESC' THEN It.Name END DESC,
        CASE WHEN _sortField = 'vuln_id' AND _sortOrder='ASC' THEN I.VulnerabilityID END,
        CASE WHEN _sortField = 'vuln_id' AND _sortOrder='DESC' THEN I.VulnerabilityID END DESC,
        CASE WHEN _sortField = 'device_id' AND _sortOrder='ASC' THEN I.DeviceID END,
        CASE WHEN _sortField = 'device_id' AND _sortOrder='DESC' THEN I.DeviceID END DESC,
        CASE WHEN _sortField = 'due_date' AND _sortOrder='ASC' THEN I.DueDate END,
        CASE WHEN _sortField = 'due_date' AND _sortOrder='DESC' THEN I.DueDate END DESC,
        CASE WHEN _sortField = 'approval' AND _sortOrder='ASC' THEN I.Approval END,
        CASE WHEN _sortField = 'approval' AND _sortOrder='DESC' THEN I.Approval END DESC,
        CASE WHEN _sortField = 'active' AND _sortOrder='ASC' THEN I.Active END,
        CASE WHEN _sortField = 'active' AND _sortOrder='DESC' THEN I.Active END DESC,
        CASE WHEN _sortField = 'port' AND _sortOrder='ASC' THEN I.Port END,
        CASE WHEN _sortField = 'port' AND _sortOrder='DESC' THEN I.Port END DESC,
        CASE WHEN _sortField = 'created_date' AND _sortOrder='ASC' THEN I.DBCreatedDate END,
        CASE WHEN _sortField = 'created_date' AND _sortOrder='DESC' THEN I.DBCreatedDate END DESC,
        CASE WHEN _sortField = 'updated_date' AND _sortOrder='ASC' THEN I.DBUpdatedDate END,
        CASE WHEN _sortField = 'updated_date' AND _sortOrder='DESC' THEN I.DBUpdatedDate END DESC,
        CASE WHEN _sortField = 'created_by' AND _sortOrder='ASC' THEN I.CreatedBy END,
        CASE WHEN _sortField = 'created_by' AND _sortOrder='DESC' THEN I.CreatedBy END DESC,
        CASE WHEN _sortField = 'updated_by' AND _sortOrder='ASC' THEN I.UpdatedBy END,
        CASE WHEN _sortField = 'updated_by' AND _sortOrder='DESC' THEN I.UpdatedBy END DESC
    LIMIT _offset,_limit;
END