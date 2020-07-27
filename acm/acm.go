package acm

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"

	"github.com/axelspringer/swerve/config"

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

// NewACM creates a new instance
func NewACM(hostPolicy autocert.HostPolicy, cache autocert.Cache, cfg *config.Configuration) *autocert.Manager {
	client := &acme.Client{
		HTTPClient: createHttpClient(cfg),
	}
	if !cfg.Prod {
		client.DirectoryURL = cfg.LetsencryptUrl
	} else {
		client.DirectoryURL = acme.LetsEncryptURL
	}
	return &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: hostPolicy,
		Cache:      cache,
		Client:     client,
	}
}

func createHttpClient(cfg *config.Configuration) *http.Client {
	if cfg.UsePebble {
		cpool := x509.NewCertPool()
		cpool.AppendCertsFromPEM([]byte(cfg.PebbleCA))

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: cpool,
			},
		}

		httpclient := &http.Client{
			Transport: tr,
		}
		return httpclient
	}
	return nil
}
