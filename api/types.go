package api

import (
	"net/http"

	"github.com/axelspringer/swerve/config"
)

// API is an API server
type API struct {
	Model  ModelAdapter
	server *http.Server
	Config *config.Configuration
}
