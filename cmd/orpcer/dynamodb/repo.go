package dynamodb

import (
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

	for {
		res, err := r.DB.Client.Scan(&dynamodb.ScanInput{
			TableName: aws.String(r.DB.Config.CertTable),
			Limit:     aws.Int64(25),
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
		if _, ok := res.LastEvaluatedKey["domain"]; !ok {
			break
		}
	}

	return certs, nil
}
