/*
  RETURN ExceptedDetection
  Title              VARCHAR(36)  NULL
  IP                 VARCHAR(36)  NULL
  Hostname           VARCHAR(36)  NULL
  VulnerabilityID    VARCHAR(36)  NULL
  VulnerabilityTitle VARCHAR(100) NULL
  Approval           VARCHAR(100) NULL
  DueDate            DATETIME     NULL
  AssignmentGroup    VARCHAR(100) NULL
  OS                 VARCHAR(100) NULL
  OSRegex            VARCHAR(100) NULL
  IgnoreID           VARCHAR(36)  NOT
  IgnoreType         INT          NOT
*/

DROP PROCEDURE IF EXISTS `GetExceptionDetections`;

CREATE PROCEDURE `GetExceptionDetections`(_offset INT, _limit INT, _orgID VARCHAR(36), _sortField NVARCHAR(255), _sortOrder NVARCHAR(255),
    _Title VARCHAR(36), _IP VARCHAR(36), _Hostname VARCHAR(36), _VulnID VARCHAR(36), _VulnTitle TEXT, _Approval VARCHAR(100), _DueDate VARCHAR(100), _AssignmentGroup VARCHAR(100), _OS VARCHAR(100), _OSRegex VARCHAR(100))
#BEGIN#
BEGIN

    SELECT T.Title, Dev.IP, Dev.Hostname, Vuln.SourceVulnID, Vuln.Title, I.Approval, I.DueDate, T.AssignmentGroup, Dev.OS, I.OSRegex, I.ID, I.TypeId from `Ignore` I
          LEFT JOIN Detection Det ON Det.IgnoreID = I.ID
          LEFT JOIN Ticket T ON T.DetectionID = Det.ID
          JOIN Device Dev on Det.DeviceID = Dev.AssetID
          JOIN VulnerabilityInfo Vuln ON Det.VulnerabilityID = Vuln.ID
    WHERE I.OrganizationID = _orgID
        AND (_Title = '' OR T.Title LIKE CONCAT('%', _Title, '%'))
        AND (_IP = '' OR Dev.IP LIKE CONCAT('%', _IP, '%'))
        AND (_Hostname = '' OR Dev.HostName LIKE CONCAT('%', _Hostname, '%'))
        AND (_VulnID = '' OR Vuln.SourceVulnId LIKE CONCAT('%', _VulnID, '%'))
#         AND (_VulnTitle = '' OR Vuln.Title LIKE CONCAT('%', _VulnTitle, '%'))
        AND (_Approval = '' OR I.Approval LIKE CONCAT('%', _Approval, '%'))
        AND (_AssignmentGroup = '' OR T.AssignmentGroup LIKE CONCAT('%', _AssignmentGroup, '%'))
        AND (_OS = '' OR Dev.OS LIKE CONCAT('%', _OS, '%'))
        AND (_OSRegex = '' OR I.OSRegex LIKE CONCAT('%', _OSRegex, '%'))
        AND (_DueDate = '' OR I.DueDate LIKE CONCAT('%', _DueDate, '%'))
    ORDER BY
        CASE WHEN _sortField = 'title' AND _sortOrder='ASC' THEN T.Title END,
        CASE WHEN _sortField = 'ip' AND _sortOrder='ASC' THEN Dev.IP END,
        CASE WHEN _sortField = 'hostname' AND _sortOrder='ASC' THEN Dev.Hostname END,
        CASE WHEN _sortField = 'approval' AND _sortOrder='ASC' THEN I.Approval END,
        CASE WHEN _sortField = 'vulnid' AND _sortOrder='ASC' THEN Vuln.SourceVulnID END,
        CASE WHEN _sortField = 'vulntitle' AND _sortOrder='ASC' THEN Vuln.Title END,
        CASE WHEN _sortField = 'expires' AND _sortOrder='ASC' THEN I.DueDate END,
        CASE WHEN _sortField = 'assignmentgroup' AND _sortOrder='ASC' THEN T.AssignmentGroup END,
        CASE WHEN _sortField = 'os' AND _sortOrder='ASC' THEN Dev.OS END,

        CASE WHEN _sortField = 'title' AND _sortOrder='DESC' THEN T.Title END DESC,
        CASE WHEN _sortField = 'ip' AND _sortOrder='DESC' THEN Dev.IP END DESC,
        CASE WHEN _sortField = 'hostname' AND _sortOrder='DESC' THEN Dev.Hostname END DESC,
        CASE WHEN _sortField = 'approval' AND _sortOrder='DESC' THEN I.Approval END DESC,
        CASE WHEN _sortField = 'vulnid' AND _sortOrder='DESC' THEN Vuln.SourceVulnID END DESC,
        CASE WHEN _sortField = 'vulntitle' AND _sortOrder='DESC' THEN Vuln.Title END DESC,
        CASE WHEN _sortField = 'expires' AND _sortOrder='DESC' THEN I.DueDate END DESC,
        CASE WHEN _sortField = 'assignmentgroup' AND _sortOrder='DESC' THEN T.AssignmentGroup END DESC,
        CASE WHEN _sortField = 'os' AND _sortOrder='DESC' THEN Dev.OS END DESC
    LIMIT _offset, _limit;
END