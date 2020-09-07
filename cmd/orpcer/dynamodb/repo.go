package dynamodb

import (
	"encoding/base64"
	"net/url"

	"github.com/axelspringer/swerve/cmd/orpcer/config"
	"github.com/axelspringer/swerve/cmd/orpcer/orphan"
	"github.com/sirupsen/logrus"

	"github.com/axelspringer/swerve/cmd/orpcer/log"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Repo struct {
	Log *logrus.Logger
	DB  Dyno
}

func NewRepo(cfg config.Swerve) Repo {
	var err error

	r := Repo{}
	r.DB.Config.CertTable = cfg.TableCertCache
	r.DB.Client, err = NewDatabase(cfg)
	r.Log = log.SetupLogger(cfg.LogLevel, cfg.LogLevel)
	if err != nil {
		r.Log.Fatal(err)
	}

	return r
}

func (r Repo) GetCerts() ([]orphan.Cert, error) {
	var certs []orphan.Cert
	var startkey map[string]*dynamodb.AttributeValue

	for {
		res, err := r.DB.Client.Scan(&dynamodb.ScanInput{
			TableName:         aws.String(r.DB.Config.CertTable),
			Limit:             aws.Int64(25),
			ExclusiveStartKey: startkey,
		})
		if err != nil {
			return nil, err
		}

		var someCerts []orphan.Cert
		err = dynamodbattribute.UnmarshalListOfMaps(res.Items, &someCerts)
		if err != nil {
			return nil, err
		}

		certs = append(certs, someCerts...)
		value, ok := res.LastEvaluatedKey["domain"]

		if ok {
			lastkey := base64.StdEncoding.EncodeToString([]byte(*value.S))
			lastkey = url.QueryEscape(lastkey)
			startkey = map[string]*dynamodb.AttributeValue{
				"domain": {
					S: aws.String(*value.S),
				},
			}
			continue
		}
		break
	}

	return certs, nil
}
