package config

import (
	"github.com/axelspringer/swerve/acm"
	"github.com/axelspringer/swerve/database"
)

// Configuration contains the config for the app
type Configuration struct {
	Database          database.Config
	API               Api
	ACM               acm.Config
	HTTPListenerPort  int
	HTTPSListenerPort int
	LogLevel          string
	LogFormatter      string
	Bootstrap         bool
	CacheInterval     int
}

type Api struct {
	COR      string
	Version  string
	Listener int
	Secret   string
}
