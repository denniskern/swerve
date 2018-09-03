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
	Wildcard     bool   `json:"wildcard"`
	Certificate  string `json:"certificate"`
	RedirectCode int    `json:"code"`
	Description  string `json:"description"`
	Created      string `json:"created"`
	Modified     string `json:"modified"`
}

// ExportDomains model
type ExportDomains struct {
	Domains []Domain `json:"domains"`
}
