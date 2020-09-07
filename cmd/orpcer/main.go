package main

import (
	"os"

	"github.com/axelspringer/swerve/cmd/orpcer/config"
	"github.com/axelspringer/swerve/cmd/orpcer/dynamodb"
	"github.com/axelspringer/swerve/cmd/orpcer/log"
	"github.com/axelspringer/swerve/cmd/orpcer/orphan"
)

func main() {
	cfg := config.Get()

	logger := log.SetupLogger(cfg.LogLevel, cfg.LogFormat)

	repo := dynamodb.NewRepo(cfg)
	client := orphan.NewClient(repo, cfg)

	err := client.GetOrphanCerts()
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	if len(client.OrphanCerts) > 0 {
		logger.Infof("Found cert(s) which older than %d days", cfg.MaxAge)
		for _, cert := range client.OrphanCerts {
			tmpl := "-> domain: %s, age: %d, created at: %s"
			logger.Infof(tmpl, cert.Domain, cert.Age, cert.CreatedAt)
		}
	} else {
		logger.Infof("no-orphan-certs-found")
		if cfg.Verbose {
			client.PrintCertsAsTable()
		}
	}
}
