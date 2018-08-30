package api

import (
	"net/http"

	"github.com/axelspringer/swerve/src/db"
)

// Handler for API calls
type Handler struct {
}

// Server model
type Server struct {
	db       *db.DynamoDB
	Server   *http.Server
	Listener string
}
