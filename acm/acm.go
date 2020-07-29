package acm

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"

	"github.com/axelspringer/swerve/log"

	"github.com/axelspringer/swerve/config"

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

// NewACM creates a new instance
func NewACM(hostPolicy autocert.HostPolicy, cache autocert.Cache, cfg *config.Configuration) *autocert.Manager {
	client := &acme.Client{
		HTTPClient: createHttpClient(cfg),
	}
	if cfg.UsePebble {
		client.DirectoryURL = cfg.LetsencryptUrl
	} else if cfg.UseStage {
		client.DirectoryURL = LetsEncryptStagingURL
	} else {
		client.DirectoryURL = acme.LetsEncryptURL
	}
	return &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      cache,
		HostPolicy: hostPolicy,
		Client:     client,
	}
}

func createHttpClient(cfg *config.Configuration) *http.Client {
	if cfg.UsePebble {
		cpool := x509.NewCertPool()

		cpool.AppendCertsFromPEM([]byte(cfg.PebbleCA))

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		resp, err := client.Get(cfg.PebbleCAUrl + "/roots/0")
		if err != nil {
			log.Fatal(err)
		}
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		// cpool.AppendCertsFromPEM(cert)

		tr2 := &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: cpool,
				// InsecureSkipVerify: true,
			},
		}

		httpclient := &http.Client{
			Transport: tr2,
		}
		return httpclient
	}
	return nil
}
