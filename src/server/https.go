// Copyright 2018 Axel Springer SE
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"crypto/tls"
	"net/http"

	"github.com/axelspringer/swerve/src/certificate"
	"github.com/axelspringer/swerve/src/log"
)

// Listen to the https
func (h *HTTPS) Listen() error {
	log.Infof("HTTPS listening to %s", h.listener)
	return h.server.ListenAndServeTLS("", "")
}

// redirectHandler redirects the request to the domain redirect location
func (h *HTTPS) redirectHandler() http.Handler {
	return handlerWithLogging(func(w http.ResponseWriter, r *http.Request) {
		hostHeader := r.Host
		domain, err := h.certManager.GetDomain(hostHeader)
		msg := "Response with status code %d"

		// regular domain lookup
		if domain != nil && err == nil {
			redirectURL, redirectCode := domain.GetRedirect(r.URL)
			http.Redirect(w, r, redirectURL, redirectCode)
			log.Infof(msg, redirectCode)
			return
		}

		log.Infof(msg, http.StatusNotFound)
		http.NotFound(w, r)
	})
}

// NewHTTPSServer creates a new instance
func NewHTTPSServer(listener string, certManager *certificate.Manager) *HTTPS {
	server := &HTTPS{
		certManager: certManager,
		listener:    listener,
	}

	server.server = &http.Server{
		Addr: listener,
		TLSConfig: &tls.Config{
			GetCertificate: server.certManager.GetCertificate,
		},
		Handler: server.redirectHandler(),
		//		ReadTimeout:  DefaultRedirectRequestTimeout,
		//		WriteTimeout: DefaultRedirectRequestTimeout,
	}

	return server
}
