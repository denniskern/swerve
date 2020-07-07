package config

import (
	"github.com/TetsuyaXD/evade/api"
	"github.com/TetsuyaXD/evade/database"
)

// NewConfiguration creates a new instance
func NewConfiguration() *Configuration {
	return &Configuration{
		API: api.Config{
			Listener: defaultAPIListener,
		},
		HTTPListenerPort:  defaultHTTPListener,
		HTTPSListenerPort: defaultHTTPSListener,
		LogLevel:          defaultLogLevel,
		LogFormatter:      defaultLogFormatter,
		Prod:              defaultProd,
		Bootstrap:         defaultBootstrap,
		CacheInterval:     defaultCacheInterval,
		Database: database.Config{
			TableNamePrefix: defaultTableNamePrefix,
			Region:          defaultDBRegion,
			TableRedirects:  defaultTableRedirects,
			TableCertCache:  defaultTableCertCache,
		},
	}
}
