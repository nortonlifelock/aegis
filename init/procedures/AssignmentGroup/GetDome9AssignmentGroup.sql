/*
    RETURN CISAssignments
    OrganizationID    VARCHAR(36)  NOT NULL
    CloudAccountID    VARCHAR(100) NULL
    BundleID          VARCHAR(100) NULL
    RuleRegex         VARCHAR(200) NULL
    RuleHash          VARCHAR(100) NULL
    AssignmentGroup   VARCHAR(100) NOT NULL
*/

DROP PROCEDURE IF EXISTS `GetCISAssignments`;

CREATE PROCEDURE `GetCISAssignments` (_OrganizationID VARCHAR(36))
#BEGIN#
SELECT
    D.OrganizationID,
    D.CloudAccountID,
    D.BundleID,
    D.RuleRegex,
    D.RuleHash,
    D.AssignmentGroup
FROM Dome9Assignment D
WHERE D.OrganizationID = _OrganizationID;
