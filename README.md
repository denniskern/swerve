# swerve

swerve is a redirection service that uses autocert to generate https certificates automatically

## build

    make

## run local dynamodb

    make run/dynamo

## run swerve (examples)

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