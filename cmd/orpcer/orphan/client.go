package orphan

import (
	"orpcer/config"
	"orpcer/log"
	"time"

	"github.com/sirupsen/logrus"
)

type Client struct {
	MaxAge int
	Repo   DBAdapter
	Log    *logrus.Logger
}

func NewClient(repo DBAdapter, cfg config.Swerve) Client {
	logger := log.SetupLogger(cfg.LogLevel, cfg.LogFormat)
	return Client{cfg.MaxAge, repo, logger}
}

func (c *Client) GetOrphanCerts() ([]Cert, error) {

	var orphanCerts []Cert
	certs, err := c.Repo.GetCerts()
	if err != nil {
		return nil, err
	}

	for _, cert := range certs {
		duration := time.Since(cert.CreatedAt)
		days := int(duration.Hours() / 24)
		cert.Age = days
		if days > c.MaxAge {
			c.Log.Debugf("Found orphan cert for domain:%s (%d days old, maxage is %d)", cert.Domain, days, c.MaxAge)
			orphanCerts = append(orphanCerts, cert)
		}

	}

	return orphanCerts, nil
}
