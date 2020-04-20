/*
    RETURN AssignmentRules
    AssignmentGroup    VARCHAR(100) NULL
    Assignee           VARCHAR(100) NULL
    OrganizationID     VARCHAR(36)  NOT NULL
    GroupID            VARCHAR(300) NULL
    VulnTitleRegex     VARCHAR(100) NULL
    HostnameRegex      VARCHAR(100) NULL
    TagKeyID           INT          NULL
    TagKeyRegex        VARCHAR(100) NULL
    Priority           INT          NOT NULL
*/

DROP PROCEDURE IF EXISTS `GetAssignmentRulesByOrg`;

CREATE PROCEDURE `GetAssignmentRulesByOrg` (_OrganizationID VARCHAR(36))
    #BEGIN#
SELECT
    AssignmentGroup,
    Assignee,
    OrganizationID,
    GroupID,
    VulnTitleRegex,
    HostnameRegex,
    TagKeyID,
    TagKeyRegex,
    Priority
FROM AssignmentRules A
    WHERE A.OrganizationID = _OrganizationID order by A.Priority DESC;
