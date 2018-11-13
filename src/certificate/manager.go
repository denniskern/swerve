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

package certificate

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"

	"github.com/axelspringer/swerve/src/db"
	"github.com/axelspringer/swerve/src/log"
	"golang.org/x/crypto/acme/autocert"
)

var (
	errHostNotConfigured = errors.New("acme/autocert: host not configured")
)

// NewManager creates a new instance
func NewManager(d *db.DynamoDB) *Manager {
	manager := &Manager{
		certCache: newPersistentCertCache(d),
	}

	manager.acmeManager = &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: manager.allowHostPolicy,
		Cache:      manager.certCache,
	}

	return manager
}

// allowHostPolicy decides which host shall pass
func (m *Manager) allowHostPolicy(_ context.Context, host string) error {
	if _, found := m.certCache.IsDomainAcceptable(host); !found {
		log.Info("allowHostPolicy - errHostNotConfigured")
		return errHostNotConfigured
	}

	log.Infof("allowHostPolicy - Host accaptable %s", host)
	return nil
}

// GetCertificate wrapper for the cert getter
func (m *Manager) GetCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	log.Infof("GetCertificate - %#v", *hello)
	return m.acmeManager.GetCertificate(hello)
}

// Serve http.Handler bridge
func (m *Manager) Serve(fallback http.Handler, w http.ResponseWriter, r *http.Request) {
	m.acmeManager.HTTPHandler(fallback).ServeHTTP(w, r)
}

// GetDomain by name
func (m *Manager) GetDomain(host string) (*db.Domain, error) {
	log.Infof("GetDomain - %s", host)
	if domain, found := m.certCache.IsDomainAcceptable(host); found {
		log.Infof("GetDomain - %#v", domain)
		return domain, nil
	}
	log.Info("GetDomain - errHostNotConfigured")
	return nil, errHostNotConfigured
}
