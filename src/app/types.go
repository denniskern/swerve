package app

import (
	"github.com/hammi85/swerve/src/certificate"
	"github.com/hammi85/swerve/src/configuration"
	"github.com/hammi85/swerve/src/db"
)

// Application model
type Application struct {
	Config       *configuration.Configuration
	DynamoDB     *db.DynamoDB
	Certificates *certificate.Manager
}
