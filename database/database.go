package database

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/axelspringer/swerve/log"
	"github.com/pkg/errors"
)

// NewDatabase creates a new instance
func NewDatabase(c Config) (*Database, error) {
	db := &Database{}

	dynamoConfig := &aws.Config{
		Region: aws.String(c.Region),
	}

	if c.Endpoint != "" {
		dynamoConfig.Endpoint = aws.String(c.Endpoint)
	}

	dynamoConfig.Credentials = credentials.NewStaticCredentials(c.Key, c.Secret, "")

	sess, err := session.NewSession(dynamoConfig)
	if err != nil {
		return db, errors.WithMessage(err, ErrSessionCreate)
	}

	db.Service = dynamodb.New(sess)
	db.Config = c

	return db, nil
}

// Prepare prepares the database
func (d *Database) Prepare() error {
	log.Debug("Preparing tables")
	err := d.prepareTable(d.Config.TableRedirects, keyNameRedirectsTable)
	if err != nil {
		return errors.WithMessagef(err, ErrTableCreate, d.Config.TableRedirects)
	}
	err = d.prepareTable(d.Config.TableCertCache, keyNameCertCacheTable)
	if err != nil {
		return errors.WithMessagef(err, ErrTableCreate, d.Config.TableCertCache)
	}
	err = d.prepareTable(d.Config.TableUsers, keyNameUsersTable)
	if err != nil {
		return errors.WithMessagef(err, ErrTableCreate, d.Config.TableUsers)
	}
	err = d.createDefaultUser()
	if err != nil {
		return errors.WithMessage(err, ErrDefaultUserCreate)
	}

	return nil
}

func (d *Database) prepareTable(tableName string, keyName string) error {

	tableCreate := &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		KeySchema: []*dynamodb.KeySchemaElement{
			{AttributeName: aws.String(keyName), KeyType: aws.String("HASH")},
		},
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{AttributeName: aws.String(keyName), AttributeType: aws.String("S")},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
	}
	tableDescribe := &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}

	if _, err := d.Service.DescribeTable(tableDescribe); err != nil {
		log.Warn(errors.WithMessagef(err, "Table '%s' does not exist", tableName).Error())
		if _, cerr := d.Service.CreateTable(tableCreate); cerr != nil {
			return cerr
		}
		log.Infof("Table '%s' created", tableName)
	}
	return nil
}

func (d *Database) createDefaultUser() error {
	// The User is not required for swerve to work properly
	if d.Config.DefaultUserPW == "" && d.Config.DefaultUser == "" {
		log.Info("no dynamodb default user will be created because username and password are not provided")
		return nil
	}

	if d.Config.DefaultUserPW == "" || d.Config.DefaultUser == "" {
		return fmt.Errorf("can't create default user, you must provide a username and a password")
	}

	log.Infof("creating default user '%s'", defaultDynamoUser)
	_, err := d.Service.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(d.Config.TableUsers),
		Item: map[string]*dynamodb.AttributeValue{
			keyNameUsersTable: {
				S: aws.String(d.Config.DefaultUser),
			},
			attrNamePwd: {
				S: aws.String(d.Config.DefaultUserPW),
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}
