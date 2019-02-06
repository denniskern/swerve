#!/bin/sh

# create table Domains
aws dynamodb create-table --table-name Domains --attribute-definitions AttributeName=id,AttributeType=S --key-schema AttributeName=id,KeyType=HASH --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1 --endpoint-url http://dynamodb:8000

# create table DomainsTLSCache
aws dynamodb create-table --table-name DomainsTLSCache --attribute-definitions AttributeName=cacheKey,AttributeType=S --key-schema AttributeName=cacheKey,KeyType=HASH --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1 --endpoint-url http://dynamodb:8000

