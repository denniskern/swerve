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
	"net/http"

	"github.com/TetsuyaXD/swerve/src/certificate"
	"github.com/TetsuyaXD/swerve/src/log"
)

// Listen to the http
func (h *HTTP) Listen() error {
	log.Infof("HTTP listening to %s", h.listener)
	return h.server.ListenAndServe()
}

// handle normal redirect request on http
func (h *HTTP) handleRedirect(w http.ResponseWriter, r *http.Request) {
	hostHeader := r.Host
	domain, err := h.certManager.GetDomain(hostHeader)

	// regular domain lookup
	if domain != nil && err == nil {
		redirectURL, redirectCode := domain.GetRedirect(r.URL)
		http.Redirect(w, r, redirectURL, redirectCode)
		return
	}

	http.NotFound(w, r)
}

// Handler for requests
func (h *HTTP) handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.certManager.Serve(http.HandlerFunc(h.handleRedirect), w, r)
	})
}

// NewHTTPServer creates a new instance
func NewHTTPServer(listener string, certManager *certificate.Manager) *HTTP {
	server := &HTTP{
		listener:    listener,
		certManager: certManager,
	}

	server.server = &http.Server{
		Addr:    listener,
		Handler: server.handler(),
		//		ReadTimeout:  DefaultRedirectRequestTimeout,
		//		WriteTimeout: DefaultRedirectRequestTimeout,
	}

	return server
}
