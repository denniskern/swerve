package https

import (
	"net/http"

	"github.com/hammi85/swerve/src/certificate"
)

// Server model
type Server struct {
	certManager *certificate.Manager
	Server      *http.Server
	Listener    string
}
