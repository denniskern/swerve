package db

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"golang.org/x/crypto/bcrypt"
)

// CheckPassword checks pwd hash on db against entered plain pwd
func (d *DynamoDB) CheckPassword(username string, plainPwd string) error {
	pwd, err := d.getPassword(username)
	if err != nil {
		return err
	}

	byteHash := []byte(pwd)
	bytePlain := []byte(plainPwd)

	err = bcrypt.CompareHashAndPassword(byteHash, bytePlain)
	if err != nil {
		return err
	}

	return nil
}

func (d *DynamoDB) getPassword(username string) (string, error) {
	res, err := d.Service.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(DBTablePrefix + dbUsersTable),
		Key: map[string]*dynamodb.AttributeValue{
			"name": {
				S: aws.String(username),
			},
		},
	})
	if err != nil || len(res.Item) == 0 {
		return "", fmt.Errorf("Error while getting item. %v", err)
	}

	userRes := &User{}
	dynamodbattribute.UnmarshalMap(res.Item, &userRes)
	if err != nil {
		return "", err
	}

	return userRes.Password, nil
}
