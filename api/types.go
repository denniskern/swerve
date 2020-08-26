package api

import (
	"net/http"
)

// API is an API server
type API struct {
	Model  ModelAdapter
	server *http.Server
	Config Config
}

// Config contains the API config
type Config struct {
	Secret   string `long:"api-jwt-sec" env:"SWERVE_API_JWT_SECRET" description:"JWT token"`
	Version  string `long:"api-version" env:"SWERVE_API_VERSION" description:"api version in pattern of v1 or v2" default:"v1"`
	COR      string `long:"api-ui-url" env:"SWERVE_API_UI_URL" description:"The url is needed for cors headers"`
	Listener int    `long:"api-listener" env:"SWERVE_API_LISTENER" description:"Listener port for the api" default:"8082"`
}
