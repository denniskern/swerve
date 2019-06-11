// Copyright 2018 Axel Springer SE
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package db

import (
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/axelspringer/swerve/src/log"
)

var (
	dbDomainTableName = getOSPrefixEnv("DOMAINS")
	dbCacheTableName  = getOSPrefixEnv("DOMAINS_TLS_CACHE")
	dbUsersTable      = getOSPrefixEnv("USERS")
)

var (
	// DBTablePrefix holds the db prefix
	DBTablePrefix = ""
)

const (
	envPrefix = "SWERVE_"
)

// getOSPrefixEnv get os env
func getOSPrefixEnv(s string) string {
	if e := strings.TrimSpace(os.Getenv(envPrefix + s)); len(e) > 0 {
		return e
	}

	return ""
}

// NewDynamoDB creates a new instance
func NewDynamoDB(c *DynamoConnection, bootstrap bool) (*DynamoDB, error) {
	ddb := &DynamoDB{}

	config := &aws.Config{
		Region: aws.String(c.Region),
	}

	if c.Endpoint != "" {
		config.Endpoint = aws.String(c.Endpoint)
	}

	if c.Key != "" && c.Secret != "" {
		config.Credentials = credentials.NewStaticCredentials(c.Key, c.Secret, "")
	}

	sess, err := session.NewSession(config)

	if err != nil {
		return nil, err
	}

	ddb.Session = sess
	ddb.Service = dynamodb.New(sess)

	if bootstrap {
		ddb.prepareTable()
	}

	return ddb, nil
}

// prepareTable checks for the main table
func (d *DynamoDB) prepareTable() {
	dbDomainCacheTableCreate := &dynamodb.CreateTableInput{
		TableName: aws.String(DBTablePrefix + dbCacheTableName),
		KeySchema: []*dynamodb.KeySchemaElement{
			{AttributeName: aws.String("cacheKey"), KeyType: aws.String("HASH")},
		},
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{AttributeName: aws.String("cacheKey"), AttributeType: aws.String("S")},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
	}
	dbDomainCacheTableDescribe := &dynamodb.DescribeTableInput{
		TableName: aws.String(DBTablePrefix + dbCacheTableName),
	}
	dbDomainTableCreate := &dynamodb.CreateTableInput{
		TableName: aws.String(DBTablePrefix + dbDomainTableName),
		KeySchema: []*dynamodb.KeySchemaElement{
			{AttributeName: aws.String("domain"), KeyType: aws.String("HASH")},
		},
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{AttributeName: aws.String("domain"), AttributeType: aws.String("S")},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
	}
	dbDomainTableDescribe := &dynamodb.DescribeTableInput{
		TableName: aws.String(DBTablePrefix + dbDomainTableName),
	}
	dbUsersTableCreate := &dynamodb.CreateTableInput{
		TableName: aws.String(DBTablePrefix + dbUsersTable),
		KeySchema: []*dynamodb.KeySchemaElement{
			{AttributeName: aws.String("name"), KeyType: aws.String("HASH")},
		},
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{AttributeName: aws.String("name"), AttributeType: aws.String("S")},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
	}
	dbUsersTableDescribe := &dynamodb.DescribeTableInput{
		TableName: aws.String(DBTablePrefix + dbUsersTable),
	}

	// setup the domain table by spec
	if _, err := d.Service.DescribeTable(dbDomainTableDescribe); err != nil {
		log.Error(err)
		log.Info("Table 'Domains' didn't exists. Creating ...")
		if _, cerr := d.Service.CreateTable(dbDomainTableCreate); cerr != nil {
			log.Fatal(cerr)
		}
		log.Info("Table 'Domains' created")
	}
	// setup the domain tls cache table by spec
	if _, err := d.Service.DescribeTable(dbDomainCacheTableDescribe); err != nil {
		log.Error(err)
		log.Info("Table 'DomainsTLSCache' didn't exists. Creating ...")
		if _, cerr := d.Service.CreateTable(dbDomainCacheTableCreate); cerr != nil {
			log.Fatal(cerr)
		}
		log.Info("Table 'DomainsTLSCache' created")
	}
	// setup the users table by spec
	if _, err := d.Service.DescribeTable(dbUsersTableDescribe); err != nil {
		log.Error(err)
		log.Info("Table 'SwerveUsers' didn't exists. Creating ...")
		if _, cerr := d.Service.CreateTable(dbUsersTableCreate); cerr != nil {
			log.Fatal(cerr)
		}
		log.Info("Table 'SwerveUsers' created")
	}
}
