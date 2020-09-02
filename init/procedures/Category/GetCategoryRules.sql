/*
  RETURN CategoryRule
  ID                      VARCHAR(36)  NOT NULL
  OrganizationID          VARCHAR(36)  NOT NULL
  SourceID                VARCHAR(36)  NOT NULL
  VulnerabilityTitle      VARCHAR(300) NULL
  VulnerabilityCategory   VARCHAR(300) NULL
  VulnerabilityType       VARCHAR(100) NULL
  Category                VARCHAR(200) NOT NULL
*/

DROP PROCEDURE IF EXISTS `GetCategoryRules`;

CREATE PROCEDURE `GetCategoryRules` (_OrgID VARCHAR(36), _SourceID VARCHAR(36))
    #BEGIN#
SELECT
    C.ID, OrganizationID, SourceID, VulnerabilityTitle, VulnerabilityCategory, VulnerabilityType, Category
FROM CategoryRule C
WHERE C.OrganizationID = _OrgID AND C.SourceID = _SourceID;