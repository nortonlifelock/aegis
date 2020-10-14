/*
    RETURN AssignmentRules
    AssignmentGroup       VARCHAR(100) NULL
    Assignee              VARCHAR(100) NULL
    ApplicationName       VARCHAR(300) NULL
    OrganizationID        VARCHAR(36)  NOT NULL
    GroupID               VARCHAR(300) NULL
    VulnTitleRegex        VARCHAR(100) NULL
    ExcludeVulnTitleRegex VARCHAR(100) NULL
    HostnameRegex         VARCHAR(100) NULL
    OSRegex               VARCHAR(100) NULL
    CategoryRegex         VARCHAR(300) NULL
    TagKeyID              INT          NULL
    TagKeyRegex           VARCHAR(100) NULL
    PortCSV               VARCHAR(100) NULL
    ExcludePortCSV        VARCHAR(100) NULL
    Priority              INT          NOT NULL
*/

DROP PROCEDURE IF EXISTS `GetAssignmentRulesByOrg`;

CREATE PROCEDURE `GetAssignmentRulesByOrg` (_OrganizationID VARCHAR(36))
    #BEGIN#
SELECT
    AssignmentGroup,
    Assignee,
    ApplicationName,
    OrganizationID,
    GroupID,
    VulnTitleRegex,
    ExcludeVulnTitleRegex,
    HostnameRegex,
    OSRegex,
    CategoryRegex,
    TagKeyID,
    TagKeyRegex,
    PortCSV,
    ExcludePortCSV,
    Priority
FROM AssignmentRules A
    WHERE A.OrganizationID = _OrganizationID order by A.Priority DESC;
