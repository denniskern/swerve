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
	"errors"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/axelspringer/swerve/src/log"
)

const (
	dbDomainTableName = "Domains"
	dbCacheTableName  = "DomainsTLSCache"
)

var (
	dbDomainCacheTableCreate = &dynamodb.CreateTableInput{
		TableName: aws.String(dbCacheTableName),
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
	dbDomainCacheTableDescribe = &dynamodb.DescribeTableInput{
		TableName: aws.String(dbCacheTableName),
	}
	dbListAllDomains = &dynamodb.ScanInput{
		TableName: aws.String(dbDomainTableName),
	}
	dbDomainTableCreate = &dynamodb.CreateTableInput{
		TableName: aws.String(dbDomainTableName),
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
	dbDomainTableDescribe = &dynamodb.DescribeTableInput{
		TableName: aws.String(dbDomainTableName),
	}
)

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
}

// DeleteTLSCacheEntry deletes a chache entry
func (d *DynamoDB) DeleteTLSCacheEntry(key string) error {
	_, err := d.Service.DeleteItem(&dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"cacheKey": {
				S: aws.String(key),
			},
		},
		TableName: aws.String(dbCacheTableName),
	})

	return err
}

// GetTLSCache items from tls cache table
func (d *DynamoDB) GetTLSCache(key string) ([]byte, error) {
	res, err := d.Service.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(dbCacheTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"cacheKey": {
				S: aws.String(key),
			},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("Error while getting item. %v", err)
	}

	entryRes := &TLSCacheEntry{}
	if err = dynamodbattribute.UnmarshalMap(res.Item, &entryRes); err == nil {
		return []byte(entryRes.Value), nil
	}

	return nil, nil
}

// UpdateTLSCache updates the tls cache
func (d *DynamoDB) UpdateTLSCache(key string, data []byte) error {
	_, err := d.Service.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"cacheKey": {
				S: aws.String(key),
			},
			"cacheValue": {
				S: aws.String(string(data)),
			},
		},
		ReturnConsumedCapacity: aws.String("TOTAL"),
		TableName:              aws.String(dbCacheTableName),
	})

	return err
}

// UpdateCertificateData updates the cert data if a domain entry exist
func (d *DynamoDB) UpdateCertificateData(domain string, data []byte) error {
	_, err := d.Service.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(dbDomainTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"domain": {
				S: aws.String(domain),
			},
		},
		UpdateExpression: aws.String("set certificate = :c"),
		ReturnValues:     aws.String("UPDATED_NEW"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":c": {
				S: aws.String(string(data)),
			},
		},
	})

	return err
}

// DeleteByDomain items from domains table
func (d *DynamoDB) DeleteByDomain(domain string) (bool, error) {
	out, err := d.Service.DeleteItem(&dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"domain": {
				S: aws.String(domain),
			},
		},
		TableName: aws.String(dbDomainTableName),
	})

	return out != nil && err == nil, err
}

// FetchByDomain items from domains table
func (d *DynamoDB) FetchByDomain(domain string) (*Domain, error) {
	res, err := d.Service.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(dbDomainTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"domain": {
				S: aws.String(domain),
			},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("Error while getting item. %v", err)
	}

	domainRes := &Domain{}
	if err = dynamodbattribute.UnmarshalMap(res.Item, &domainRes); err == nil {
		return domainRes, nil
	}

	return nil, nil
}

// FetchAll items from domains table
func (d *DynamoDB) FetchAll() ([]Domain, error) {
	itemList, err := d.Service.Scan(dbListAllDomains)

	if err != nil {
		return nil, fmt.Errorf("Error while fetching domain items %v", err)
	}

	recs := []Domain{}
	err = dynamodbattribute.UnmarshalListOfMaps(itemList.Items, &recs)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal Dynamodb Scan Items, %v", err)
	}

	return recs, nil
}

// InsertDomain stores a domain
func (d *DynamoDB) InsertDomain(domain Domain) error {
	mm, err := dynamodbattribute.MarshalMap(domain)

	if err != nil {
		return err
	}

	_, err = d.Service.PutItem(&dynamodb.PutItemInput{
		Item:      mm,
		TableName: aws.String(dbDomainTableName),
	})

	return err
}

// DeleteAllDomains deletes all items from the domains table
func (d *DynamoDB) DeleteAllDomains() error {
	domains, err := d.FetchAll()

	if err != nil {
		return err
	}

	for _, do := range domains {
		_, err = d.DeleteByDomain(do.Name)
		if err != nil {
			return err
		}
	}

	return nil
}

// Import imports a export set
func (d *DynamoDB) Import(e *ExportDomains) error {
	for _, do := range e.Domains {
		mm, err := dynamodbattribute.MarshalMap(do)

		if err != nil {
			return err
		}

		_, err = d.Service.PutItem(&dynamodb.PutItemInput{
			Item:      mm,
			TableName: aws.String(dbDomainTableName),
		})
	}

	return nil
}

// Validate the domain
func (d *Domain) Validate() []error {
	res := []error{}

	if d.ID == "" {
		res = append(res, errors.New("Invalid id"))
	}

	validURL, err := url.Parse("//" + d.Name)
	if d.Name == "" || err != nil || validURL.Path != "" {
		res = append(res, errors.New("Invalid domain name"))
	}

	if d.Created == "" || d.Modified == "" {
		res = append(res, errors.New("Invalid domain date"))
	}

	if d.Redirect == "" {
		res = append(res, errors.New("Invalid domain redirect target"))
	}

	if d.RedirectCode < 300 || d.RedirectCode > 399 {
		res = append(res, errors.New("Invalid redirect http status code"))
	}

	return res
}

// GetRedirect returns calculated routes
func (d *Domain) GetRedirect(reqURL *url.URL) (string, int) {
	code := d.RedirectCode
	reURL := strings.TrimRight(d.Redirect, "/")
	rePath := ""
	reQuery := ""

	if d.Promotable == true {
		rePath = reqURL.Path

		if len(reqURL.RawQuery) > 0 {
			reQuery = "?" + reqURL.RawQuery
		}
	}

	if d.PathMapping != nil && len(*d.PathMapping) > 0 {
		for _, p := range *d.PathMapping {
			if p.To == "" {
				continue
			}
			// we match the path prefix
			if strings.HasPrefix(rePath, p.From) {
				rePath = rePath[len(p.From):]
				// path redirect
				if strings.HasPrefix(p.To, "http://") || strings.HasPrefix(p.To, "https://") {
					reURL = strings.TrimRight(p.To, "/")
				} else {
					rePath = path.Join(p.To, rePath)
				}
				break
			}
		}
	}

	return reURL + rePath + reQuery, code
}
