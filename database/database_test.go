package database

import (
	"encoding/base64"
	"log"
	"net/url"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var (
	d *Database
)

func TestMain(m *testing.M) {
	var err error
	cfg := Config{
		TableNamePrefix: "",
		Region:          "eu-west-1",
		TableRedirects:  "Redirects",
		TableCertCache:  "DomainsTLSCache",
		TableUsers:      "Users",
		Endpoint:        "http://localhost:8000",
	}
	d, err = NewDatabase(cfg)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}

func TestUserTable(t *testing.T) {
	var startkey map[string]*dynamodb.AttributeValue
	var newCursor string
	tablePrefix := d.Config.TableNamePrefix
	tableName := d.Config.TableUsers

	res, err := d.Service.Scan(&dynamodb.ScanInput{
		TableName:         aws.String(tablePrefix + tableName),
		Limit:             aws.Int64(25),
		ExclusiveStartKey: startkey,
	})
	if err != nil {
		t.Fatal(err)
	}

	users := []User{}
	err = dynamodbattribute.UnmarshalListOfMaps(res.Items, &users)
	if err != nil {
		t.Fatal(err)
	}

	if len(users) == 0 {
		t.Log("No Users found in dynamodb")
	}

	for _, v := range users {
		t.Logf("Found User %s with PW %s", v.Name, v.Pwd)
	}

	val, ok := res.LastEvaluatedKey[keyNameRedirectsTable]
	if ok {
		newCursor = base64.StdEncoding.EncodeToString([]byte(*val.S))
		newCursor = url.QueryEscape(newCursor)
	} else {
		newCursor = "EOF"
	}
}
