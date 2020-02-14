DROP PROCEDURE IF EXISTS `CreateTicketingJob`;

# The ScanStartDate is a string instead of a DATETIME, because the
# time.Now().Format("2006-01-02T00:00:00Z")
# the mysql driver cannot handle arrays, so only a single group is processed at a time
CREATE PROCEDURE `CreateTicketingJob`(GroupID INT, OrgID VARCHAR(36), ScanStartDate VARCHAR(100))
    #BEGIN#
    INSERT INTO JobHistory(
        JobId,
        ConfigID,
        StatusId,
        Priority,
        CurrentIteration,
        Payload
    ) VALUES (
        (select ID from Job where Struct = 'TicketingJob'),
        (select ID from JobConfig where JobID = (select ID from Job where Struct = 'TicketingJob') AND OrganizationID = OrgID),
        1,
        4,
        0,
        JSON_OBJECT(
            'mindate',
            ScanStartDate,
            'group',
            GroupID
        )
    );