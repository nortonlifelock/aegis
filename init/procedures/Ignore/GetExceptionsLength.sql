/*
  RETURN QueryData SINGLE
  Length                      INT           NOT
*/
DROP PROCEDURE IF EXISTS `GetExceptionsLength`;

CREATE PROCEDURE `GetExceptionsLength`(_offset INT, _limit INT, _orgID VARCHAR(36), _sortField NVARCHAR(255), _sortOrder NVARCHAR(255),
                                       _Title VARCHAR(36), _IP VARCHAR(36), _Hostname VARCHAR(36), _VulnID VARCHAR(36), _Approval VARCHAR(100), _DueDate VARCHAR(100), _AssignmentGroup VARCHAR(100), _OS VARCHAR(100), _OSRegex VARCHAR(100),
                                       _TypeID INT)
#BEGIN#
BEGIN

    SELECT count(*) from `Ignore` I
        LEFT JOIN Detection Det ON Det.IgnoreID = I.ID
        LEFT JOIN Ticket T ON T.DetectionID = Det.ID
        JOIN Device Dev on Det.DeviceID = Dev.AssetID
        JOIN VulnerabilityInfo Vuln ON Det.VulnerabilityID = Vuln.ID
    WHERE I.OrganizationID = _orgID AND I.Active = b'1'
      AND (_TypeID = -1 OR _TypeID = I.TypeId)
      AND (_Title = '' OR T.Title LIKE CONCAT('%', _Title, '%'))
      AND (_IP = '' OR Dev.IP LIKE CONCAT('%', _IP, '%'))
      AND (_Hostname = '' OR Dev.HostName LIKE CONCAT('%', _Hostname, '%'))
      AND (_VulnID = '' OR Vuln.SourceVulnId LIKE CONCAT('%', _VulnID, '%'))
      AND (_Approval = '' OR I.Approval LIKE CONCAT('%', _Approval, '%'))
      AND (_AssignmentGroup = '' OR T.AssignmentGroup LIKE CONCAT('%', _AssignmentGroup, '%'))
      AND (_OS = '' OR Dev.OS LIKE CONCAT('%', _OS, '%'))
      AND (_OSRegex = '' OR I.OSRegex LIKE CONCAT('%', _OSRegex, '%'))
      AND (_DueDate = '' OR I.DueDate LIKE CONCAT('%', _DueDate, '%'));
END