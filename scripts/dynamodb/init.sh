#!/bin/sh

# create table Domains
aws dynamodb create-table --table-name Domains --attribute-definitions AttributeName=domain,AttributeType=S --key-schema AttributeName=domain,KeyType=HASH --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1 --endpoint-url http://dynamodb:8000

# create table DomainsTLSCache
aws dynamodb create-table --table-name DomainsTLSCache --attribute-definitions AttributeName=cacheKey,AttributeType=S --key-schema AttributeName=cacheKey,KeyType=HASH --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1 --endpoint-url http://dynamodb:8000
