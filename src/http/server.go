package http

import (
	nethttp "net/http"
	"net/url"
	"path"

	"github.com/hammi85/swerve/src/certificate"
	"github.com/hammi85/swerve/src/log"
)

// promoteRedirect append the path + querystring to the redirect host
func promoteRedirect(redirect string, reqURL *url.URL) string {
	newRedirect := path.Join(redirect, reqURL.Path)
	if len(reqURL.RawQuery) > 0 {
		newRedirect = newRedirect + "?" + reqURL.RawQuery
	}

	return newRedirect
}

// Listen to the http
func (s *Server) Listen() error {
	log.Infof("HTTP listening to %s", s.Listener)
	return s.Server.ListenAndServe()
}

// handle normal redirect request on http
func (s *Server) handleRedirect(w nethttp.ResponseWriter, r *nethttp.Request) {
	hostHeader := r.Host
	domain, err := s.certManager.GetDomain(hostHeader)

	// regular domain lookup
	if domain != nil && err == nil {
		redirect := domain.Redirect
		if domain.Promotable {
			redirect = promoteRedirect(redirect, r.URL)
		}

		log.Infof("http redirect %s => %s", r.URL.String(), redirect)
		nethttp.Redirect(w, r, redirect, domain.RedirectCode)
		return
	}

	nethttp.NotFound(w, r)
}

// Handler for requests
func (s *Server) Handler() nethttp.Handler {
	return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		s.certManager.Serve(nethttp.HandlerFunc(s.handleRedirect), w, r)
	})
}

// NewServer creates a new instance
func NewServer(listener string, certManager *certificate.Manager) *Server {
	server := &Server{
		Listener:    listener,
		certManager: certManager,
	}

	server.Server = &nethttp.Server{
		Addr:    listener,
		Handler: server.Handler(),
	}

	return server
}
