package main

import (
	"time"

	"orpcer/config"
	dynamodb "orpcer/dynamodb"
	"orpcer/orphan"

	"orpcer/log"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	cfg := config.Get()
	spew.Dump(cfg)

	logger := log.SetupLogger(cfg.LogLevel, cfg.LogFormat)

	repo := dynamodb.NewRepo(cfg)
	client := orphan.NewClient(repo, cfg)

	for {
		certs, err := client.GetOrphanCerts()
		if err != nil {
			logger.Error(err)

		}

		if len(certs) > 0 {
			logger.Infof("Found cert(s) which older than %d days", cfg.MaxAge)
			for _, cert := range certs {
				tmpl := "-> domain: %s, age: %d, created at: %s"
				logger.Infof(tmpl, cert.Domain, cert.Age, cert.CreatedAt)
			}
		} else {
			logger.Infof("no-orphan-certs-found")
		}
		time.Sleep(time.Minute * 60)
	}
}
