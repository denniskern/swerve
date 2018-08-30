package https

import (
	"net/http"

	"github.com/axelspringer/swerve/src/certificate"
)

// Server model
type Server struct {
	certManager *certificate.Manager
	Server      *http.Server
	Listener    string
}
