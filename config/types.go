package config

import (
	"github.com/TetsuyaXD/evade/api"
	"github.com/TetsuyaXD/evade/database"
)

// Configuration contains the config for the app
type Configuration struct {
	Database          database.Config
	API               api.Config
	HTTPListenerPort  int
	HTTPSListenerPort int
	LogLevel          string
	LogFormatter      string
	Prod              bool
	Bootstrap         bool
	CacheInterval     int
}
