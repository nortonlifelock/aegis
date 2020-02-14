/*
  GEN Ticket
  ID                  INT                     NOT
  OrganizationID      NVARCHAR(36)            NOT
  ScanID              INT                     NOT
  DeviceID            VARCHAR(36)             NOT
  VulnerabilityID     NVARCHAR(250)           NOT
  Title               NVARCHAR(500)           NOT
  Description         TEXT                    NULL
  TicketType          NVARCHAR(30)            NULL
  Status              NVARCHAR(30)            NULL
  ResolutionStatus    NVARCHAR(30)            NULL
  DueDate             DATETIME                NULL
  CreatedDate         DATETIME                NULL
  UpdatedDate         DATETIME                NULL
  AlertDate           DATETIME                NULL
  ResolutionDate      DATETIME                NULL
  Project             NVARCHAR(255)           NULL
  Summary             NVARCHAR(500)           NULL
  AssignedTo          NVARCHAR(255)           NULL
  ReportedBy          NVARCHAR(255)           NULL
  AssignmentGroup     NVARCHAR(255)           NULL
  GroupID             VARCHAR(100)            NOT
  Labels              NVARCHAR(500)           NULL
  MethodOfDiscovery   NVARCHAR(20)            NULL
  Priority            NVARCHAR(20)            NULL
  OperatingSystem     NVARCHAR(255)           NULL
  CVSS                FLOAT                   NULL
  CVEReferences       NVARCHAR(1000)          NULL
  IPAddress           NVARCHAR(255)           NULL
  MacAddress          NVARCHAR(255)           NULL
  HostName            NVARCHAR(255)           NULL
  ServicePorts        NVARCHAR(255)           NULL
  VulnerabilityTitle  NVARCHAR(500)           NULL
  CERF                NVARCHAR(30)            NOT
  CERFExpirationDate  DATETIME                NOT
  DBCreatedDate       DATETIME                NOT
  DBUpdatedDate       DATETIME                NULL
  CloudID             NVARCHAR(128)           NOT
  Configs             NVARCHAR(10)            NOT
  LastChecked         DATETIME                NULL
  Solution            TEXT                    NULL
  OrgCode             NVARCHAR(25)            NULL
  OSDetailed          NVARCHAR(255)           NULL
  VendorReferences    TEXT                    NULL
*/
#BEGIN#