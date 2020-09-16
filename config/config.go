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
	//if s.DynamoDB.Bootstrap && s.DynamoDB.DefaultUserPW == "" {
	//	errors = append(errors, fmt.Errorf("If dynamodb bootstrap is enabled, then SWERVE_DYNO_DEFAULT_PW or --dyno-default-user-pw is required"))
	//}

	return errors
}

func Get() Swerve {
	s := Swerve{}
	p := flags.NewParser(&s, flags.IgnoreUnknown|flags.HelpFlag|flags.PrintErrors)
	if _, err := p.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
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
