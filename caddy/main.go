package main

import (
	"log"
	"net/http"

	"github.com/caddyserver/certmagic"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	def := certmagic.NewDefault()

	cache := certmagic.NewCache(certmagic.CacheOptions{
		GetConfigForCert: func(cert certmagic.Certificate) (*certmagic.Config, error) {
			// do whatever you need to do to get the right
			// configuration for this certificate; keep in
			// mind that this config value is used as a
			// template, and will be completed with any
			// defaults that are set in the Default config
			return &certmagic.Config{
				// ...
			}, nil
		},
	})

	magic := certmagic.New(cache, *def)

	myACME := certmagic.NewACMEManager(magic, certmagic.ACMEManager{
		CA:                   certmagic.LetsEncryptStagingCA,
		Email:                "dennis.kern@axelspringer.com",
		Agreed:               true,
		DisableHTTPChallenge: true,
	})

	magic.Issuer = myACME

	err := magic.ManageSync([]string{"digital15.swervetest.de"})
	if err != nil {
		log.Fatal(err)
	}

	tlsConfig := magic.TLSConfig()

	s := &http.Server{
		Addr:      ":443",
		TLSConfig: tlsConfig,
		Handler:   mux,
	}

	log.Fatal(s.ListenAndServeTLS("", ""))
}
