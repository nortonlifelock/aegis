# Aegis System

[![Build Status](https://api.travis-ci.org/nortonlifelock/aegis.svg?branch=master)](https://travis-ci.org/nortonlifelock/aegis)
[![Go Report Card](https://goreportcard.com/badge/github.com/nortonlifelock/aegis)](https://goreportcard.com/report/github.com/nortonlifelock/aegis)
[![GoDoc](https://godoc.org/github.com/nortonlifelock/aegis?status.svg)](https://godoc.org/github.com/nortonlifelock/aegis)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

aegis

Installing Aegis:

- Step 0: Create a JIRA project for vulnerability management by Aegis
- Step 1: Create the database schema for Aegis (*we use Aegis*)
- Step 2: Execute the Aegis Scaffold (github.com/nortonlifelock/aegis-scaffold)
- Step 3: Setup your Aegis Configuration File

Base Json Configuration File: 

```json
{
    "db_path":"127.0.0.1",
    "db_port":"3306",
    "db_username":"",
    "db_password":"",
    "db_schema": "",
    "logpath":"",
    "key_id":"",
    "sns_id":"",
    "log_to_console":true,
    "log_to_db":true,
    "log_to_file":true,
    "detailed_console_logging":false,
    "debug":true,
    "api_port":4040,
    "websocket_protocol":"ws",
    "transport_protocol":"http",
    "backend_location":"localhost:4040",
    "ui_location":"localhost:4200",
    "path_to_aegis": "",
    "ad_servers": ["", ""],
    "ad_ldap_tls_port": 636,
    "ad_base_dn": "OU=Corporate,OU=Standard,OU=People,DC=corp,DC=ad",
    "ad_skip_verify": true,
    "ad_member_of_attribute": "memberOf",
    "ad_search_string": "(sAMAccountName=%s)"
  }
```
- Step 4: Execute Aegis
```
aegis -config app.json -cpath "path to configuration file"
```