# swerve

Swerve is a redirection service that uses autocert to generate https certificates automatically. It 

## Build

    make

or

    make build/local

### Build docker container

    make build/docker

## Parameter

### Environment parameter

* SWERVE_DB_ENDPOINT - AWS endpoint for the DynamoDB
* SWERVE_DB_REGION - AWS region for the DynamoDB
* SWERVE_DB_KEY - AWS key for credential based access
* SWERVE_DB_SECRET - AWS secret for credential based access
* SWERVE_API - Address for the API listener
* SWERVE_HTTP - Address for the HTTP listener
* SWERVE_HTTPS - Address for the HTTPS listener
* BOOTSTRAP - DB table preparation
* LOG_LEVEL - Log level info, debug, warning, error, fatal and panic

### Application parameter

* db-endpoint - AWS endpoint for the DynamoDB
* db-region - AWS region for the DynamoDB
* db-key - AWS key for credential based access
* db-secret - AWS secret for credential based access
* bootstrap - DB table preparation
* api - Address for the API listener
* http - Address for the HTTP listener
* https - Address for the HTTPS listener

## Run

### run local dynamodb

    make run/dynamo

### run swerve (examples)

    SWERVE_DB_ENDPOINT=http://localhost:8000 SWERVE_DB_REGION=eu-west-1 SWERVE_HTTPS=:8081 ./bin/swerve

## API calls

### Get all domains

    curl -X GET http://localhost:8082/domain

### Get a single domain by name

    curl -X GET http://localhost:8082/domain/<name>

### Insert a new domain

    curl -X POST \
        http://localhost:8082/domain/ \
        -d '{
            "domain": "<name>",
            "redirect": "https://my.redirect.target.com",
            "code": 308,
            "description": "Example domain entry"
        }'

### Purge a domain by name

    curl -X DELETE http://localhost:8082/domain/<name>