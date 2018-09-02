# SWERVE

Swerve is a redirection service that uses autocert to generate https certificates automatically. The domain and certificate data are stored in a DynamoDB

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

## API

### Domain model

    {
        "id": "guid v4 will be generated",
        "domain": "my.domain.com",
        "path": ["/deep/link", "/foo"],
        "redirect": "https://my.redirect.com"
        "promotable": false,
        "wildcard": false,
        "certificate": "certificate data",
        "code": 308,
        "description": "Meanful description of this redirection",
        "created": "generated date",
        "modified": "generated date"
    }

#### id

Will be generated

#### domain

The domain name to keep track on. e.g. ```my.redirect.com```

#### path

Optional deep path match e.g. with "path": "/redirect/path" only my.domain.com/redirect/path matches

#### redirect

Redirection target

#### promotable

Promotable redirects are attaching the path of the request to the redirection location e.g.
with ```"promotable": true``` my.domain.com/foo/bar/baz.html will be redirected to https://my.redirect.com/foo/bar/baz.html
with ```"promotable": false``` my.domain.com/foo/bar/baz.html will be redirected to https://my.redirect.com

#### wildcard

Wildcard domains use their certificate on every subdomain. e.g. with ```"wildcard": true``` https://sub1.my.domain and https://supersub2.my.domain will use the certificate of https://my.domain

#### certificate

The certificate data

#### code

The redirection code. It has to be greater or equal 300 and less or equal than 399

#### description

Meanful description of the domain entry

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