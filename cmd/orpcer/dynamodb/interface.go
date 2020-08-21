package dynamodb

import "github.com/aws/aws-sdk-go/service/dynamodb"

type DynamoDBAdapter interface {
	Scan(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error)
}
