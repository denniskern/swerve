package acm

import (
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

// NewACM creates a new instance
func NewACM(hostPolicy autocert.HostPolicy, cache autocert.Cache, stage bool) *autocert.Manager {
	client := &acme.Client{}
	if stage {
		client.DirectoryURL = LetsEncryptStagingURL
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
