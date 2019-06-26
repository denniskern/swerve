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
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

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
	reURL := d.Redirect
	rePath := ""
	reQuery := ""
	reqPath := reqURL.EscapedPath()

	if d.Promotable == true {
		rePath = reqURL.EscapedPath()

		if len(reqURL.RawQuery) > 0 {
			reQuery = "?" + reqURL.RawQuery
		}
	} else if reqURL.RawQuery != "" {
		reqPath += "?" + reqURL.RawQuery
	}

	if d.PathMapping != nil && len(*d.PathMapping) > 0 {
		for _, p := range *d.PathMapping {
			// skip empty path mapping
			if p.To == "" {
				continue
			}
			// we match the path prefix
			if strings.HasPrefix(reqPath, p.From) {
				rePath = reqPath[len(p.From):]
				// path redirect
				if strings.HasPrefix(p.To, "http://") || strings.HasPrefix(p.To, "https://") {
					reURL = p.To
				} else {
					if d.Promotable {
						rePath = path.Join(p.To, rePath)
					} else {
						rePath = p.To
					}
				}
				break
			}
		}
	}

	if strings.HasSuffix(reURL, "/") && strings.HasPrefix(rePath, "/") {
		rePath = strings.TrimLeft(rePath, "/")
	}

	return reURL + rePath + reQuery, code
}

// UpdateCertificateData updates the cert data if a domain entry exist
func (d *DynamoDB) UpdateCertificateData(domain string, data []byte) error {
	_, err := d.Service.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(DBTablePrefix + dbDomainTableName),
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
		TableName: aws.String(DBTablePrefix + dbDomainTableName),
	})

	return out != nil && err == nil, err
}

// FetchByDomain items from domains table
func (d *DynamoDB) FetchByDomain(domain string) (*Domain, error) {
	res, err := d.Service.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(DBTablePrefix + dbDomainTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"domain": {
				S: aws.String(domain),
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("Error while getting item. %v", err)
	}

	domainDBRes := &DomainDB{}
	if err = dynamodbattribute.UnmarshalMap(res.Item, &domainDBRes); err != nil {
		return nil, err
	}

	domainRes, err := domainDBRes.toDomain()
	if err != nil {
		return nil, err
	}

	return &domainRes, nil
}

// FetchAllSorted returns all items from table with a sorted paths (important for the redirects!)
func (d *DynamoDB) FetchAllSorted() ([]Domain, error) {
	domains, err := d.FetchAll()
	if err != nil {
		return nil, err
	}
	for _, domain := range domains {
		domain.sortPathMap()
	}
	return domains, nil
}

// FetchAll items from domains table
func (d *DynamoDB) FetchAll() ([]Domain, error) {
	domains := []Domain{}

	itemList, err := d.Service.Scan(&dynamodb.ScanInput{
		TableName: aws.String(DBTablePrefix + dbDomainTableName),
	})

	if err != nil {
		return nil, fmt.Errorf("Error while fetching domain items %v", err)
	}

	recs := []DomainDB{}
	err = dynamodbattribute.UnmarshalListOfMaps(itemList.Items, &recs)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal Dynamodb Scan Items, %v", err)
	}

	for _, domaindb := range recs {
		domain, err := domaindb.toDomain()
		if err != nil {
			return nil, err
		}
		domains = append(domains, domain)
	}

	return domains, nil
}

// FetchAllPaginated items from domains table
func (d *DynamoDB) FetchAllPaginated(cursor *string) ([]Domain, *string, error) {
	var startkey map[string]*dynamodb.AttributeValue
	var newCursor string
	domains := []Domain{}

	if cursor != nil {
		sk, err := base64.StdEncoding.DecodeString(*cursor)
		if err != nil {
			return nil, nil, fmt.Errorf("Error while decoding cursor, %v", err)
		}
		startkey = map[string]*dynamodb.AttributeValue{
			"domain": {
				S: aws.String(string(sk)),
			},
		}
	}

	itemList, err := d.Service.Scan(&dynamodb.ScanInput{
		TableName:         aws.String(DBTablePrefix + dbDomainTableName),
		ExclusiveStartKey: startkey,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("Error while fetching domain items, %v", err)
	}

	recs := []DomainDB{}
	err = dynamodbattribute.UnmarshalListOfMaps(itemList.Items, &recs)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to unmarshal Dynamodb Scan Items, %v", err)
	}

	val, ok := itemList.LastEvaluatedKey["domain"]
	if ok {
		newCursor = base64.StdEncoding.EncodeToString([]byte(*val.S))
	} else {
		newCursor = "EOF"
	}

	for _, domaindb := range recs {
		domain, err := domaindb.toDomain()
		if err != nil {
			return nil, nil, err
		}
		domains = append(domains, domain)
	}

	return domains, &newCursor, nil
}

// InsertDomain stores a domain
func (d *DynamoDB) InsertDomain(domain Domain) error {
	domaindb, err := domain.toDomainDB()
	if err != nil {
		return err
	}

	mm, err := dynamodbattribute.MarshalMap(domaindb)
	if err != nil {
		return err
	}

	_, err = d.Service.PutItem(&dynamodb.PutItemInput{
		Item:      mm,
		TableName: aws.String(DBTablePrefix + dbDomainTableName),
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
		ddb, err := do.toDomainDB()
		if err != nil {
			return err
		}

		mm, err := dynamodbattribute.MarshalMap(ddb)

		if err != nil {
			return err
		}

		_, err = d.Service.PutItem(&dynamodb.PutItemInput{
			Item:      mm,
			TableName: aws.String(DBTablePrefix + dbDomainTableName),
		})
	}

	return nil
}

func (d *Domain) toDomainDB() (DomainDB, error) {
	var buf bytes.Buffer
	domaindb := DomainDB{
		Domain: *d,
	}

	pm, _ := json.Marshal(d.PathMapping)
	if len(pm) > 200000 {
		domaindb.PathMapping = nil
		writer := gzip.NewWriter(&buf)
		if _, err := writer.Write(pm); err != nil {
			return domaindb, err
		}
		if err := writer.Flush(); err != nil {
			return domaindb, err
		}

		if err := writer.Close(); err != nil {
			return domaindb, err
		}

		bytes := buf.Bytes()
		domaindb.BinPathMapping = &bytes
	}

	return domaindb, nil
}

func (d *DomainDB) toDomain() (Domain, error) {
	var pl PathList
	domain := d.Domain

	if domain.PathMapping == nil && d.BinPathMapping != nil {
		b := bytes.NewBuffer(*d.BinPathMapping)
		reader, err := gzip.NewReader(b)
		if err != nil {
			return domain, err
		}

		s, err := ioutil.ReadAll(reader)
		if err != nil {
			return domain, err
		}

		if err := reader.Close(); err != nil {
			return domain, err
		}

		err = json.Unmarshal(s, &pl)
		if err != nil {
			return domain, err
		}

		domain.PathMapping = &pl
	}

	return domain, nil
}
