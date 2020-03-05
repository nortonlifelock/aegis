/*
    RETURN AssignmentRules
    AssignmentGroup    VARCHAR(100) NULL
    Assignee           VARCHAR(100) NULL
    OrganizationID     VARCHAR(36)  NOT NULL
    VulnTitleSubstring VARCHAR(200) NULL
    VulnTitleRegex     VARCHAR(100) NULL
    TagKeyID           INT          NULL
    TagKeyValue        VARCHAR(100) NULL
    Priority           INT          NOT NULL
*/

DROP PROCEDURE IF EXISTS `GetAssignmentRulesByOrg`;

CREATE PROCEDURE `GetAssignmentRulesByOrg` (_OrganizationID VARCHAR(36))
    #BEGIN#
SELECT
    AssignmentGroup,
    Assignee,
    OrganizationID,
    VulnTitleSubstring,
    VulnTitleRegex,
    TagKeyID,
    TagKeyValue,
    Priority
FROM AssignmentRules A
    WHERE A.OrganizationID = _OrganizationID order by A.Priority DESC;
