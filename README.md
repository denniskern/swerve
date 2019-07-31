# :cyclone: SWERVE

[![Build Status](https://travis-ci.org/axelspringer/swerve.svg?branch=master)](https://travis-ci.org/axelspringer/swerve.svg?branch=master)

Swerve is a redirection service that uses autocert to generate https certificates automatically. The domain and certificate data are stored in a DynamoDB

## Setup the DB tables

Create the dynamodb tables. This example with AWS cli based on the assumption that the table prefix is empty

    aws dynamodb create-table --table-name Domains --attribute-definitions AttributeName=domain,AttributeType=S --key-schema AttributeName=domain,KeyType=HASH --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1
    aws dynamodb create-table --table-name DomainsTLSCache --attribute-definitions AttributeName=cacheKey,AttributeType=S --key-schema AttributeName=cacheKey,KeyType=HASH --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1

## Test

    make test/local

## Build

    make

or

    make build/local

or to build as linux binary

    make build/linux

### Build docker container

    make build/docker

### Start docker compose stack

    make compose/up

### Shutdown docker compose stack

    make compose/down

### Deploy local build to stack

    make deploy/local

## Parameter

### Environment parameter

* SWERVE_DB_ENDPOINT - AWS endpoint for the DynamoDB
* SWERVE_DB_REGION - AWS region for the DynamoDB
* SWERVE_DB_KEY - AWS key for credential based access
* SWERVE_DB_SECRET - AWS secret for credential based access
* SWERVE_DB_TABLE_PREFIX - DynamoDB table name prefix
* SWERVE_API - Address for the API listener
* SWERVE_HTTP - Address for the HTTP listener
* SWERVE_HTTPS - Address for the HTTPS listener
* SWERVE_BOOTSTRAP - DB table preparation
* SWERVE_LOG_LEVEL - Log level info, debug, warning, error, fatal and panic
* SWERVE_STAGING - Use letsencrypt staging api with much higher quota. Use this when you run tests
* SWERVE_API_SECRET - The bycrypt secret to check incoming pw against the pw in the database
* SWERVE_DOMAINS - The name of the domains table
* SWERVE_DOMAINS_TLS_CACHE - The name of the domains tls cache taböe
* SWERVE_USERS - The name of the table holding the user login data
* SWERVE_UI_DOMAIN - (https://swerve.tortuga.cloud) The url of the frontend (for CORS)

### Application parameter

* db-endpoint - AWS endpoint for the DynamoDB
* db-region - AWS region for the DynamoDB
* db-key - AWS key for credential based access
* db-secret - AWS secret for credential based access
* db-table-prefix - DynamoDB table name prefix
* bootstrap - DB table preparation
* api - Address for the API listener
* http - Address for the HTTP listener
* https - Address for the HTTPS listener
* client-static - Path to the API client static files
* log-level - Set the log level (info,debug,warning,error,fatal,panic)
* log-formatter - Set the log formatter (text,json)

## API

### User

You need at least 1 valid user to control the API. Passwords are stored as base64 encoded bcrypt with a uow of 12.
User login data has to be inserted manually in the login data table (env: SWERVE_USERS)
as 
name (string): testuser
password (string): JDJ5JDEyJFdNQUtzdk1ESmdyRE1sOWZ3NmJSb08xOTlIMTU3QjFCeEVXbUphd1YxSjhnUWVMY2VoNFRt

### Domain model

    {
        "id": "guid v4 will be generated",
        "domain": "my.domain.com",
        "paths": [
            {
                "from": "/match/path/prefix",
                "to": "/foo"
            },
            {
                "from": "/other/target",
                "to": "https://the.other.one/"
            }
        ],
        "redirect": "https://my.redirect.com"
        "promotable": false,
        "code": 301,
        "description": "Meanful description of this redirection",
        "created": "generated date",
        "modified": "generated date"
    }

#### id

Will be generated

#### domain

The domain name to keep track on. e.g. ```my.redirect.com```

#### paths

You can add an aditional path mapping conditional list. When defined the redirection based on the matching result of this list. Fallback is the default redirect

#### redirect

Redirection target

#### promotable

Promotable redirects are attaching the path of the request to the redirection location e.g.
with ```"promotable": true``` my.domain.com/foo/bar/baz.html will be redirected to https://my.redirect.com/foo/bar/baz.html
with ```"promotable": false``` my.domain.com/foo/bar/baz.html will be redirected to https://my.redirect.com

#### code

The redirection code. It has to be greater or equal 300 and less or equal than 399

#### description

Meanful description of the domain entry

## API calls

### Login
    curl -X POST \
        http://<api_host>:<api_port>/login \
        -d '{
            "name": "testuser",
            "password": "JDJ5JDEyJFdNQUtzdk1ESmdyRE1sOWZ3NmJSb08xOTlIMTU3QjFCeEVXbUphd1YxSjhnUWVMY2VoNFRt"
        }'

    (response returns a cookie)

### Get all domains

    curl -X GET http://<api_host>:<api_port>/api/domain

### Get a single domain by name

    curl -X GET http://<api_host>:<api_port>/api/domain/<name>

### Insert a new domain

    curl -X POST \
        http://<api_host>:<api_port>/api/domain/ \
        -d '{
            "domain": "<name>",
            "redirect": "https://my.redirect.target.com",
            "code": 308,
            "description": "Example domain entry"
        }'

### Purge a domain by name

    curl -X DELETE http://<api_host>:<api_port>/api/domain/<name>

### Export all domains

    curl -X GET http://<api_host>:<api_port>/api/export

### Import a export set

    curl -X POST \
    http://<api_host>:<api_port>/api/import \
    -H 'cache-control: no-cache' \
    -H 'content-type: application/json' \
    -d '{
        "domains": [
            {
                "domain": "my.domain.com",
                "redirect": "https://example.com",
                "promotable": false,
                "code": 308,
                "description": "example registration 2",
            }
        ]
    }'

## Example stack

Start the stack

    make compose/up

This should start a swerve, dynamodb and a test target nginx service

Lets add a target to swerve

    curl -X POST \
        http://127.0.0.1:8082/api/domain \
        -H 'cache-control: no-cache' \
        -H 'content-type: application/json' \
        -d '{"domain": "example.org", "redirect": "http://127.0.0.1:8090/", "code": 301, "description": "test", "promotable": true}'

Test the record

    curl -X GET 127.0.0.1:8082/api/domain

So lets see whether the http redirect works

    curl -X GET -H 'Host: example.org' -I 127.0.0.1:8080

The result should look like

    HTTP/1.1 301 Moved Permanently
    Location: http://127.0.0.1:8090/
    Date: Tue, 20 Nov 2018 13:00:44 GMT
    Content-Length: 57
    Content-Type: text/html; charset=utf-8

## Benchmark with vegeta (version > 12.0.0)

    printf "GET http://127.0.0.1:8080\nHost: example.org\n" > target.txt

    vegeta attack -rate 50 -duration 2m -targets target.txt | vegeta encode | \
        jaggr @count=rps \
            hist\[100,200,300,400,500\]:code \
            p25,p50,p95:latency \
            sum:bytes_in \
            sum:bytes_out | \
        jplot rps+code.hist.100+code.hist.200+code.hist.300+code.hist.400+code.hist.500 \
            latency.p95+latency.p50+latency.p25 \
            bytes_in.sum+bytes_out.sum
