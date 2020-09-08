package orphan

import (
	"fmt"
	"strconv"
	"time"

	"github.com/axelspringer/swerve/cmd/orpcer/config"
	"github.com/axelspringer/swerve/cmd/orpcer/log"
	"github.com/sirupsen/logrus"
)

type Client struct {
	MaxAge      int
	Repo        DBAdapter
	Log         *logrus.Logger
	AllCerts    []Cert
	OrphanCerts []Cert
}

func NewClient(repo DBAdapter, cfg config.Swerve) Client {
	logger := log.SetupLogger(cfg.LogLevel, cfg.LogFormat)
	return Client{cfg.MaxAge, repo, logger, nil, nil}
}

func (c *Client) GetOrphanCerts() error {

	var orphanCerts []Cert
	certs, err := c.Repo.GetCerts()
	if err != nil {
		return err
	}

	for _, cert := range certs {
		if (cert.CreatedAt == time.Time{}) {
			return fmt.Errorf("GetOrphanCerts: (%s) created_at is zero value, check if field is available. Seems the repo doesn't provide correct timestamp", cert.Domain)
		}
		duration := time.Since(cert.CreatedAt)
		days := int(duration.Hours() / 24)
		cert.Age = days
		if days > c.MaxAge {
			orphanCerts = append(orphanCerts, cert)
		}

	}
	c.OrphanCerts = orphanCerts
	c.AllCerts = certs

	return nil
}

func (c *Client) PrintCertsAsTable() {
	data := [][]string{}
	for _, cert := range c.AllCerts {
		age := strconv.Itoa(int(time.Now().Sub(cert.CreatedAt).Hours() / 24))
		data = append(data, []string{cert.Domain, age, cert.CreatedAt.Format("2006-01-02 15:04:05")})
	}
	printResults(sortSlice(data))
}
