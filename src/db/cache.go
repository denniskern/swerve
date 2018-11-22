package db

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// DeleteTLSCacheEntry deletes a chache entry
func (d *DynamoDB) DeleteTLSCacheEntry(key string) error {
	_, err := d.Service.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(DBTablePrefix + dbCacheTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"cacheKey": {
				S: aws.String(key),
			},
		},
	})

	return err
}

// GetTLSCache items from tls cache table
func (d *DynamoDB) GetTLSCache(key string) ([]byte, error) {
	res, err := d.Service.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(DBTablePrefix + dbCacheTableName),
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
		TableName: aws.String(DBTablePrefix + dbCacheTableName),
		Item: map[string]*dynamodb.AttributeValue{
			"cacheKey": {
				S: aws.String(key),
			},
			"cacheValue": {
				S: aws.String(string(data)),
			},
		},
		ReturnConsumedCapacity: aws.String("TOTAL"),
	})

	return err
}
