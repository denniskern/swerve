package db

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// DynamoDB model
type DynamoDB struct {
	Session *session.Session
	Service *dynamodb.DynamoDB
}

// DynamoConnection model
type DynamoConnection struct {
	Endpoint  string
	Key       string
	Secret    string
	TableName string
	Region    string
}

// DomainList db entry
type DomainList struct {
	Domains []Domain `json:"domains"`
}

// Domain db entry
type Domain struct {
	ID           string `json:"id"`
	Name         string `json:"domain"`
	PathPattern  string `json:"path"`
	Redirect     string `json:"redirect"`
	Promotable   bool   `json:"promotable"`
	Certificate  string `json:"certificate"`
	RedirectCode int    `json:"code"`
	Description  string `json:"description"`
	Created      string `json:"created"`
	Modified     string `json:"modified"`
}
