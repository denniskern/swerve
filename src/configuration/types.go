package configuration

import "github.com/axelspringer/swerve/src/db"

// Configuration model
type Configuration struct {
	HTTPListener  string
	HTTPSListener string
	APIListener   string
	DynamoDB      db.DynamoConnection
	LogLevel      string
	LogFormatter  string
	Bootstrap     bool
}
