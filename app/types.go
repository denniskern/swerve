package app

import (
	"github.com/TetsuyaXD/evade/api"
	"github.com/TetsuyaXD/evade/cache"
	"github.com/TetsuyaXD/evade/config"
	"github.com/TetsuyaXD/evade/http"
	"github.com/TetsuyaXD/evade/https"
)

// Application is evades app model
type Application struct {
	Config      *config.Configuration
	APIServer   *api.API
	HTTPServer  *http.HTTP
	HTTPSServer *https.HTTPS
	Cache       *cache.Cache
}
