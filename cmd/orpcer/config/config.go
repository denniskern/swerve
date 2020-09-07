package config

import (
	"fmt"
	"log"
	"os"

	flags "github.com/jessevdk/go-flags"
)

const (
	maxAgeMin = 0
	maxAgeMax = 90
)

type Swerve struct {
	AwsKey           string `long:"aws-key" env:"ORPCER_DYNO_AWS_KEY" default:"0" description:"AWS access key for dynamodb"`
	AwsSec           string `long:"aws-sec" env:"ORPCER_DYNO_AWS_SECRET" default:"0" description:"AWS secret key for dynamodb"`
	AwsRegion        string `long:"aws-region" env:"ORPCER_DYNO_AWS_REGION" required:"true" description:"AWS region for dynamodb"`
	DynamoDBEndpoint string `long:"db-endpoint" env:"ORPCER_DYNO_ENDPOINT" required:"true" description:"Endpoint of dynamodb"`
	TableCertCache   string `long:"table-certs" env:"ORPCER_DYNO_TABLE_CERTCACHE" required:"true" description:"dynamodb name of table certcache"`
	MaxAge           int    `long:"cert-maxage" env:"ORPCER_CERT_MAXAGE" required:"true" description:"Log an error if found a cert which is older than <cert-maxage>"`
	LogLevel         string `long:"log-level" env:"ORPCER_LOG_LEVEL" default:"info" description:"logging servety"`
	LogFormat        string `long:"log-format" env:"ORPCER_LOG_FORMAT" default:"test" description:"logging format"`
	Verbose          bool   `long:"verbose" env:"ORPCER_VERBOSE" description:"print verbose output"`
}

func Get() Swerve {
	s := Swerve{}
	p := flags.NewParser(&s, flags.Default)
	if _, err := p.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			fmt.Println()
			p.WriteHelp(os.Stdout)
			os.Exit(1)
		}
	}

	if s.MaxAge == maxAgeMin || s.MaxAge > maxAgeMax {
		log.Fatalf("--maxage / $SWERVE_CERT_MAXAGE must be greater than %d and less than %d", maxAgeMin, maxAgeMax)
	}

	return s
}
