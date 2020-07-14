#!/bin/sh

# create table Redirects
aws dynamodb create-table --table-name ${TABLE_REDIRECTS} --attribute-definitions AttributeName=redirect_from,AttributeType=S --key-schema AttributeName=redirect_from,KeyType=HASH --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1 --endpoint-url http://dynamodb:8000

# create table CertCache
aws dynamodb create-table --table-name ${TABLE_CERTCACHE} --attribute-definitions AttributeName=domain,AttributeType=S --key-schema AttributeName=domain,KeyType=HASH --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1 --endpoint-url http://dynamodb:8000

# create table CertCache
aws dynamodb create-table --table-name ${TABLE_USERS} --attribute-definitions AttributeName=username,AttributeType=S --key-schema AttributeName=username,KeyType=HASH --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1 --endpoint-url http://dynamodb:8000
aws dynamodb put-item --table-name ${TABLE_USERS} --item '{"username":{"S":"dkern"}, "password":{"S":"mytestpw"}}' --endpoint-url http://dynamodb:8000
