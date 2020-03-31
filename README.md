# Aegis System

[![Build Status](https://api.travis-ci.org/nortonlifelock/aegis.svg?branch=master)](https://travis-ci.org/nortonlifelock/aegis)
[![Go Report Card](https://goreportcard.com/badge/github.com/nortonlifelock/aegis)](https://goreportcard.com/report/github.com/nortonlifelock/aegis)
[![GoDoc](https://godoc.org/github.com/nortonlifelock/aegis?status.svg)](https://godoc.org/github.com/nortonlifelock/aegis)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0) [![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](http://makeapullrequest.com)

**NOTE: We're still building the Readme files for Aegis.**

Aegis is a security automation tool that is designed to coordinate between your vulnerability management solution and your vulnerability remediator while providing reliable and up to date information. No longer is the tedium of emailing spreadsheets and trusting that users are remediating vulnerabilities necessary as Aegis generates tickets that are tracked and confirmed remediated by the automation.

Users can trigger rescans of a vulnerability by marking a ticket as fixed, or letting Aegis passively rescan tickets as they become due. Aegis can automatically create dynamic scan templates and target one or many IPs with one or many scan engines to ensure the patch or configuration change was successfully implemented. If it fails, the ticket is reopened, and the proof can be added as a comment to the user. Watch your tickets close-out automatically and let users track open vulnerabilities in real time. Aegis can handle exception and false positive tracking, and can notice when a device is taken offline. Aegis tracks all scans that cover a ticket during its remediation, which is great to have when auditors come!

Installing Aegis:

- Go 1.11 or higher installed
- `$GOPATH` set
- Have `~/.aws/credentials` and `~/.aws/config` setup (or some other form of [AWS authentication](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-credentials))

## Have ready

- Credentials for all API accounts
- An enumeration of the asset groups you plan on rescanning (comma separated)
- AWS KMS encryption key
  - Must be `symmetric`

### Nexpose

- Templates for your vulnerability/discovery scans existing in your Nexpose instance and their IDs
- Nexpose site ID (integer) that you will be using for vulnerability rescans

### Qualys

- SearchList and OptionProfile for your vulnerability scans in your Qualys instance and their integer IDs
- OptionProfile for your discovery scans in your Qualys instance and it's integer ID
- A list of all the group IDs that require external scanners (comma separated)

### JIRA

- The project key for the project you plan on storing tickets
- If using your own ticket transition workflow/ticket schema, mapping must be done in the JIRA source config
    - The default ticket schema utilized by Aegis is outlines in setup/jira_ticket_schema.csv if you do not want to do the mapping
- Need your JIRA workflow XML
    - JIRA Administration (cog) -> Projects -> Workflows (on left) -> Select "actions" button next to desired workflow -> Export as XML
- (optional) a separate CERF project for tickets that remediators create for tracking exceptions and false-positives
- (optional) JIRA supports Oauth, if you have Oauth credentials, they can be set in the JIRA source config in the database after installation

### Database
- Must be MySQL
- Create a schema in your database for Aegis to utilize
---
## Installation (Mac/Linux)

```sh
git clone https://github.com/nortonlifelock/aegis

?.
cd ./aegis/cmd || exit
go install aegis.go

cd ../api/listener || exit
go install aegis-api.go

cd ../..

aegis -init-config -scaffold -init-org -cpath $PWD -sproc $PWD/init/procedures -migrate $PWD/init/migrations -tpath $PWD/init
```

## Aegis commands

```sh
# Run Aegis
# -config defaults to app.json if not specified
aegis -config app.json -cpath "path to configuration file"

# Create/update DB schema
aegis scaffold

# Create a new application config (same process that's done "aegis init")
aegis init-config

# Create a new organization (same process that's done during "aegis init")
aegis init-org
```

Base Json Configuration File: 
```json
{
  "db_path":"127.0.0.1",
  "db_port":"3306",
  "db_username":"",
  "db_password":"",
  "db_schema": "Aegis",

  "logpath":"",

  "key_id":"",
  "sns_id":"",
  "kms_region": "",
  "kms_profile": "",

  "log_to_console":true,
  "log_to_db":true,
  "log_to_file":true,
  "log_no_delete": true,
  "debug":true,

  "api_port":4040,
  "websocket_protocol":"wss",
  "transport_protocol":"https",
  "ui_location":"localhost:4200",

  "path_to_aegis": ""
}
```

The install-config script will generate this application config with your input and name the file app.json. The purpose of each field are as follows:
```
# Configurations used for the Aegis jobs
db_path
    The URL or IP address where your database resides
db_port
    The port on which your database listens
db_username
    The username that is used to authenticate against the database
db_password
    The (encrypted) password that is used to authenticate against the database
    The password is encrypted with the KMS key you provide by the install-config scirpt
db_schema
    The name of the schema that Aegis will utilize

logpath
    The path at which the logs are stored

key_id
    The ID of the (symmetric) KMS key that is used for the encryption
kms_region
    The region (e.g. us-east-1/us-west-1) that the KMS key lies in
kms_profile (optional)
    If using file authentication for AWS (~/.aws/credentials) that contains multiple profiles, you can specify which one here
sns_id (optional)
    The ID of the SNS topic that is used to publish critical/fatal logs

log_to_console
    Boolean controlling whether logs are printed to the console
log_to_db
    Boolean controlling whether logs are stored in the database 
log_to_file
    Boolean controlling whether logs are written to log files
log_no_delete
    Boolean controlling whether logs stored in files are deleted after a day
debug
    Boolean controlling whether debug logs are processed
 
# Configurations used by the Aegis API
api_port
    The port on which the Aegis API listens
websocket_protocol
    The protocol that is used for websocket connections [ws|wss]    
transport_protocol
    The protocol that is used for http connections [http|https]
ui_location
    The URL that the API lies at (used to check the origin header of HTTP requests)
path_to_aegis
    The path where the Aegis system is stored
```

- Step 4: Execute Aegis
```
aegis -config app.json -cpath "path to configuration file"
```

---
# Major processes
## Rescan Process
### Jobs involved:
1. RescanQueueJob

   Starts automatically and runs continuously. This job monitors JIRA and kicks off a RescanJob for tickets in particular statuses. There are four types of RSQ:
   1. Normal - looks tickets in a Resolved-Remediated status. These are for standard vulnerability rescans. Tickets are moved to Closed-Remediated once remediation has been confirmed by a scanner, or reopened if the scanner still detects the vulnerability.
   2. Decommission - looks for tickets in a Resolved-Decommission status. These are for confirming a device has been moved offline. These tickets are moved to Closed-Decommissioned once a scanner has confirmed they are a dead host, or reopened if the host is still alive.
   3. Passive - looks for tickets in an open, reopened, in-progress, resolved-exception status created within 20 days and due within 15 days
   4. Exception - looks for tickets in Closed-Exception status with an exception that expires within 30 days and verifies that the ticket has the vulnerability fixed before the exception expires. If the vulnerability is fixed the ticket is moved to Closed-Remediated. If the vulnerability is not fixed, the ticket is left in Closed-Exception.
2. RescanJob

    Created by the RescanQueueJob. This job is responsible for creating the scan in Qualys/Nexpose for the devices/vulnerabilities reported in the tickets.
3. ScanSyncJob

    Monitors scans created by the RescanJob. When it notices a scan has finished, the ScanSyncJob queues up a ScanCloseJob to close out the tickets based on the results of the scan
4. ScanCloseJob

    Pulls the detection information for the scanned hosts from Qualys/Nexpose and uses that detection information to close out the tickets. This job updates the LastFound/LastUpdated/LastChecked date of the ticket. If no detections were reported for the device and the device was not reported as dead, the ticket is moved to Scan-Error (likely host data purge)
    1. Normal
        1. If the device is reported as dead, it’s tickets are moved to resolved-decommission
        2. If the last updated date is before the scan date, the ticket is moved to scan-error (likely authentication issues)
        3. If the detection was marked as fixed, or no detection was returned for that device/vulnerability combo, the ticket is moved to closed-remediated
        4. If the device/vulnerability combo still exists and the ticket is in resolved-remediated, the ticket is reopened
    2. Decommission
        1. If the device is reported as dead, it’s tickets are moved to closed-decommission and the device is marked as decommissioned in the Ignore table
        2. If the device is not reported as dead and the ticket is in resolved-decommission, the ticket is reopened
    3. Passive
        1. If the device is reported as dead, it’s tickets are moved to resolved-decommission
        2. If the detection was marked as fixed, or no detection was returned for that device/vulnerability combo, the ticket is moved to closed-remediated
        3. If the device/vulnerability combo still exists and the ticket is in resolved-remediated, nothing happens
    4. Exception: same as passive

## Synchronization Processes
### Jobs involved:

1. ExceptionJob

    Creates entries in the Ignore table by pulling recently updated Closed-FalsePositive tickets and Closed-Exception tickets with non-empty CERFs
2. CloudSyncJob

    Sync cloud tag information into the Tag table of the database. We correlate could assets to scanner assets by correlating a scanner asset group to a cloud environment in the AssetGroup table (using ScannerSourceID and CloudSourceID fields)
3. VulnSyncJob

    Creates and updates entries in the VulnerabilityInfo table
4. TicketSyncJob

    Pulls recently updated tickets from JIRA and keeps the DB information on the ticket fresh
5. AssetSyncJob

    Pulls detection information from Qualys/Nexpose API and stores it in db. Uses the vulnerability information information pulled and stored during the VulnSyncJob. Also stores device information in DB. Looks at Ignore table and attaches it to detection if appropriate. Global exceptions are also checked here. Detections created/updated here are what are processed during a ticketing run

## Ticketing Process
### Jobs involved:

1. TicketingJob

    1. Unfixed detections with no ignore ID attached that were recently updated/created are pulled from the DB
    2. If they have a CVSS lower than what is specified in the CVSS min field in the organization payload, they are skipped.
    3. We first check the DB, and then JIRA for an existing ticket in the following statuses: Open, In-Progress, Reopened, Resolved-Remediated, Resolved-FalsePositive, Resolved-Decommissioned, Resolved-Exception
    4. We do mapping from cloud tags to ticket fields (if a tag mapping is specified)
    5. The ticket is assigned according to the org’s AssignmentRules table
    6. The ticket is created in JIRA and has an entry made in the Ticket table of the DB
