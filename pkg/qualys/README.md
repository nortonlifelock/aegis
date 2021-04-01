# Notes for developers
## There are multiple relevant modules of Qualys supported in this driver
### Vulnerability Management (VM)
- Accessed through the `integrations.Vscanner` interface
- Detections synced to database by AssetSyncJob
- Tickets created by TicketingJob
- Tickets rescanned by RescanQueueJob
### Web Application Scanning (WAS)
- Accessed through the `integrations.Vscanner` interface
- Detections synced to database by AssetSyncJob
- Tickets created by TicketingJob
- Tickets rescanned by RescanQueueJob
- Need to include URL for was using `host_web_app` in the Qualys SourceConfig payload
### CloudView (CV)
- Accessed through the `integrations.CISScanner` interface

If you would like to separate out the scans/ticketing of each one of these processes, separate source configs and job configs should be made for each