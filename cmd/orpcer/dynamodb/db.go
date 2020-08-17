package dynamodb

import (
	"orpcer/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Dyno struct {
	Client DynamoDBAdapter
	Config struct {
		CertTable string
	}
}

// NewDatabase creates a new instance
func NewDatabase(c config.Swerve) (*dynamodb.DynamoDB, error) {
	dynamoConfig := &aws.Config{
		Region: aws.String(c.AwsRegion),
	}
	dynamoConfig.Endpoint = aws.String(c.DynamoDBEndpoint)
	dynamoConfig.Credentials = credentials.NewStaticCredentials(c.AwsKey, c.AwsSec, "")

	sess, err := session.NewSession(dynamoConfig)

	return dynamodb.New(sess), err
}
