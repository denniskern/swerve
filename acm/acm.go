package acm

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/axelspringer/swerve/config"

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

// NewACM creates a new instance
func NewACM(hostPolicy autocert.HostPolicy, cache autocert.Cache, cfg *config.Configuration) *autocert.Manager {
	client := &acme.Client{
		HTTPClient: createHttpClient(cfg.UsePebble),
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

func createHttpClient(usePebble bool) *http.Client {
	if usePebble {
		cert, err := ioutil.ReadFile("acm/pebble/pebble.minica.pem")
		if err != nil {
			log.Fatal(err)
		}

		cpool := x509.NewCertPool()
		cpool.AppendCertsFromPEM(cert)

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				ClientCAs: cpool,
				RootCAs:   cpool,
			},
		}

		httpclient := &http.Client{
			Transport: tr,
		}
		return httpclient
	}
	return nil
}
