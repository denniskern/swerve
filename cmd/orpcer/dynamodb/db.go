package dynamodb

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/axelspringer/swerve/cmd/orpcer/config"
)

type Dyno struct {
	Client DynamoDBAdapter
	Config struct {
		CertTable string
	}
}

// NewDatabase creates a new instance
func NewDatabase(c config.Swerve) (*dynamodb.DynamoDB, error) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	} // #nosec
	httpClient := &http.Client{
		Timeout:   time.Second * 10,
		Transport: tr,
	}

	dynamoConfig := &aws.Config{
		Region: aws.String(c.AwsRegion),
	}
	dynamoConfig.Endpoint = aws.String(c.DynamoDBEndpoint)
	dynamoConfig.Credentials = credentials.NewStaticCredentials(c.AwsKey, c.AwsSec, "")
	dynamoConfig.HTTPClient = httpClient

	sess, err := session.NewSession(dynamoConfig)

	return dynamodb.New(sess), err
}
