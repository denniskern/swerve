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

	"golang.org/x/crypto/acme"

	"github.com/TetsuyaXD/swerve/src/db"
	"github.com/TetsuyaXD/swerve/src/log"
	"golang.org/x/crypto/acme/autocert"
)

const (
	// LetsEncryptStagingURL uri for the LE staging environment with higher rate limits
	LetsEncryptStagingURL = "https://acme-staging.api.letsencrypt.org/directory"
)

var (
	errHostNotConfigured = errors.New("acme/autocert: host not configured")
)

// NewManager creates a new instance
func NewManager(d *db.DynamoDB, staging bool) *Manager {
	manager := &Manager{
		CertCache: NewPersistentCertCache(d),
	}

	directoryURL := acme.LetsEncryptURL
	if staging {
		directoryURL = LetsEncryptStagingURL
		log.Infof("Using CA staging environment")
	}
	log.Infof("CA URI %s", directoryURL)

	client := &acme.Client{
		DirectoryURL: directoryURL,
	}

	manager.AcmeManager = &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: manager.AllowHostPolicy,
		Cache:      manager.CertCache,
		Client:     client,
	}

	return manager
}

// AllowHostPolicy decides which host shall pass
func (m *Manager) AllowHostPolicy(_ context.Context, host string) error {
	if _, found := m.CertCache.IsDomainAcceptable(host); !found {
		return errHostNotConfigured
	}
	return nil
}

// GetCertificate wrapper for the cert getter
func (m *Manager) GetCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	return m.AcmeManager.GetCertificate(hello)
}

// Serve http.Handler bridge
func (m *Manager) Serve(fallback http.Handler, w http.ResponseWriter, r *http.Request) {
	m.AcmeManager.HTTPHandler(fallback).ServeHTTP(w, r)
}

// GetDomain by name
func (m *Manager) GetDomain(host string) (*db.Domain, error) {
	if domain, found := m.CertCache.IsDomainAcceptable(host); found {
		return domain, nil
	}
	return nil, errHostNotConfigured
}
