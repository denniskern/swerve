package config

import (
	"fmt"
	"net/url"
	"os"

	"github.com/axelspringer/swerve/acm"

	"github.com/axelspringer/swerve/database"

	"github.com/axelspringer/swerve/api"

	"github.com/sirupsen/logrus"

	"github.com/axelspringer/swerve/log"

	"github.com/jessevdk/go-flags"
)

/*
type DynamoDB struct {
	AwsKey           string `long:"dyno-aws-key" env:"SWERVE_DB_KEY" default:"0" required:"false" description:"AWS access key for dynamodb"`
	AwsSec           string `long:"dyno-aws-sec" env:"SWERVE_DB_SECRET" default:"0" required:"false" description:"AWS secret key for dynamodb"`
	AwsRegion        string `long:"dyno-aws-region" env:"SWERVE_DB_REGION" required:"false" description:"AWS region for dynamodb" default:"eu-central-1"`
	DynamoDBEndpoint string `long:"dyno-endpoint" env:"SWERVE_DB_ENDPOINT" required:"false" description:"Endpoint of dynamodb"`
	DefaultPW        string `long:"dyno-admin-pw" env:"SWERVE_DB_DEFAULT_PW" required:"false" description:"Default PW for the admin user"`
	TableRedirects   string `long:"dyno-tbl-redirects" env:"SWERVE_DYNO_TABLE_REDIRECTS" required:"false" description:"Table name for redirects" default:"Swerve_Redirects"`
	TableCertCache   string `long:"dyno-tbl-certcache" env:"SWERVE_DYNO_TABLE_CERTCACHE" required:"false" description:"Table name for cert cache" default:"Swerve_CertCache"`
	TableUsers       string `long:"dyno-tbl-users" env:"SWERVE_DYNO_TABLE_USERS" required:"false" description:"Table name for users" default:"Swerve_Users"`
	Bootstrap        bool   `long:"dyno-bootstrap" env:"SWERVE_DYNO_BOOTSTRAP" required:"false" description:"Create tables and default user on startup"`
}

type APISettings struct {
	JWTSecret string `long:"api-jwt-sec" env:"SWERVE_API_JWT_SECRET" description:"JWT token"`
	Version   string `long:"api-version" env:"SWERVE_API_VERSION" description:"api version in pattern of v1 or v2" default:"v1"`
	UIUrl     string `long:"api-ui-url" env:"SWERVE_API_UI_URL" description:"The url is needed for cors headers"`
	Listener  int    `long:"api-listener" env:"SWERVE_API_LISTENER" description:"Listener port for the api" default:"8082"`
}
*/

type Swerve struct {
	DynamoDB             database.Config `group:"DynamoDB settings"`
	API                  api.Config      `group:"API Settings"`
	ACM                  acm.Config      `group:"ACM settings"`
	LogLevel             string          `long:"log-level" env:"SWERVE_LOG_LEVEL" default:"info" description:"logging severity" choice:"debug" choice:"info" choice:"warn" choice:"error" `
	LogFormat            string          `long:"log-format" env:"SWERVE_LOG_FORMAT" default:"text" description:"logging format" choice:"text" choice:"json"`
	HttpListener         int             `long:"http-listener" env:"SWERVE_HTTP_LISTENER" default:"8080" description:"HTTP listener port"`
	HttpsListener        int             `long:"https-listener" env:"SWERVE_HTTPS_LISTENER" default:"8081" description:"HTTPS listener port"`
	DisableHTTPChallenge bool            `long:"enable-http-challenge" env:"SWERVE_DISABLE_HTTP_CHALLENGE" description:"Disable the challenge http-01"`
	CacheInterval        int             `long:"cache-interval" env:"SWERVE_CACHE_INTERVAL" description:"renew cache in minutes, has impact on redirects not on certificates" default:"5"`
}

func (s *Swerve) Validate() []error {
	var errors []error
	if s.API.Listener > 65535 || s.API.Listener < 0 {
		errors = append(errors, fmt.Errorf("API Listener port must be between 1 and 65535"))
	}
	if s.HttpListener > 65535 || s.HttpListener < 0 {
		errors = append(errors, fmt.Errorf("HTTP Listener port must be between 1 and 65535"))
	}
	if s.HttpsListener > 65535 || s.HttpsListener < 0 {
		errors = append(errors, fmt.Errorf("HTTPS Listener port must be between 1 and 65535"))
	}
	if hasDuplicates([]int{s.HttpsListener, s.HttpListener, s.API.Listener}) {
		errors = append(errors, fmt.Errorf("listener ports must be uniq"))
	}
	if _, err := url.ParseRequestURI(s.API.COR); err != nil {
		errors = append(errors, fmt.Errorf("api ui url is not a valid URL"))
	}
	if _, err := logrus.ParseLevel(s.LogLevel); err != nil {
		errors = append(errors, fmt.Errorf("invalid log level format"))
	}
	if s.CacheInterval < 1 {
		errors = append(errors, fmt.Errorf("cache interval must be greater than 1"))
	}
	if s.ACM.UsePebble {
		s.ACM.PebbleCA = defaultPebbleCACert
	}
	return nil
}

func Get() Swerve {
	s := Swerve{}
	p := flags.NewParser(&s, flags.IgnoreUnknown|flags.HelpFlag|flags.PrintErrors)
	if _, err := p.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			fmt.Printf("%v", err)
			p.WriteHelp(os.Stdout)
			os.Exit(1)
		}
	}
	err := s.Validate()
	if err != nil {
		log.Fatal(err)
	}

	return s
}
