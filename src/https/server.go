package https

import (
	"crypto/tls"
	"net/http"
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

// Listen to the https
func (s *Server) Listen() error {
	log.Infof("HTTPS listening to %s", s.Listener)
	return s.Server.ListenAndServeTLS("", "")
}

// RedirectHandler redirects the request to the domain redirect location
func (s *Server) RedirectHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hostHeader := r.Host
		domain, err := s.certManager.GetDomain(hostHeader)

		if domain != nil && err == nil {
			redirect := domain.Redirect
			// promote redirect
			if domain.Promotable {
				redirect = promoteRedirect(redirect, r.URL)
			}

			log.Infof("https redirect %s => %s", r.URL.String(), redirect)
			http.Redirect(w, r, domain.Redirect, domain.RedirectCode)
			return
		}

		http.NotFound(w, r)
	})
}

// NewServer creates a new instance
func NewServer(listener string, certManager *certificate.Manager) *Server {
	server := &Server{
		certManager: certManager,
		Listener:    listener,
	}

	server.Server = &http.Server{
		Addr: listener,
		TLSConfig: &tls.Config{
			GetCertificate: server.certManager.GetCertificate,
		},
		Handler: server.RedirectHandler(),
	}

	return server
}
