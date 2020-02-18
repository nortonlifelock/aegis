/*
  RETURN JobSchedule
  ID                      NVARCHAR(36)           NOT
  ConfigID                NVARCHAR(36)           NOT
  Payload                 TEXT                   NULL
*/

DROP PROCEDURE IF EXISTS `GetScheduledJobsToStart`;

CREATE PROCEDURE `GetScheduledJobsToStart` (_LastChecked DATETIME)
  #BEGIN#
  SELECT
    JS.Id,
    JS.ConfigId,
    JC.Payload
  FROM JobSchedule JS
    JOIN JobConfig JC ON JC.ID = JS.ConfigID
  WHERE JS.Active = 1
       AND JS.DaysOfWeek LIKE CONCAT('%',DAYOFWEEK(NOW()),'%')
       AND ((JS.TimeOfDay <= addtime(TIME(NOW()), '00:01:00')) -- Check next 1 minute interval for a scheduled job
       AND JS.TimeOfDay >= TIME(_LastChecked) -- Ensure the time of day is in the future
       AND (JS.LastRun IS NULL -- Make sure it's only running once in a 24 hour period of time
            OR JS.LastRun <= date_add(NOW(), INTERVAL -12 HOUR))
        OR (JS.LastRun IS NOT NULL AND JS.LastRun > date_add(JS.LastRun, INTERVAL 25 HOUR)));