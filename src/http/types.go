package http

import (
	nethttp "net/http"

	"github.com/hammi85/swerve/src/certificate"
)

type fallbackHandler struct{}

// Server model
type Server struct {
	certManager *certificate.Manager
	Server      *nethttp.Server
	Listener    string
}
