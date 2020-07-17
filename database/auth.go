package database

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
)

// GetPwdHash returns a password hash from the database
func (d *Database) GetPwdHash(username string) (string, error) {
	tablePrefix := d.Config.TableNamePrefix
	tableName := d.Config.TableUsers
	res, err := d.Service.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tablePrefix + tableName),
		Key: map[string]*dynamodb.AttributeValue{
			keyNameUsersTable: {
				S: aws.String(username),
			},
		},
	})
	if err != nil {
		fmt.Printf("**1* ERR %v\n", err)
		return "", errors.WithMessage(err, ErrRedirectsFetch)
	}

	if len(res.Item) == 0 {
		fmt.Printf("**2* ERR %v\n", err)
		return "", errors.New(ErrUserNotFound)
	}

	userRes := &User{}
	if err = dynamodbattribute.UnmarshalMap(res.Item, &userRes); err != nil {
		fmt.Printf("**3* ERR %v\n", err)
		return "", errors.WithMessage(err, ErrRedirectMarshal)
	}

	return userRes.Pwd, nil
}
